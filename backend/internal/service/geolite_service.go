package service

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/netip"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/oschwald/maxminddb-golang/v2"

	"github.com/stonith404/pocket-id/backend/internal/common"
)

type GeoLiteService struct {
	mutex      sync.Mutex
	cityDBPath string
	asnDBPath  string
}

var localhostIPNets = []*net.IPNet{
	{IP: net.IPv4(127, 0, 0, 0), Mask: net.CIDRMask(8, 32)}, // 127.0.0.0/8
	{IP: net.IPv6loopback, Mask: net.CIDRMask(128, 128)},    // ::1/128
}

var privateLanIPNets = []*net.IPNet{
	{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},     // 10.0.0.0/8
	{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},  // 172.16.0.0/12
	{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)}, // 192.168.0.0/16
}

var tailscaleIPNets = []*net.IPNet{
	{IP: net.IPv4(100, 64, 0, 0), Mask: net.CIDRMask(10, 32)}, // 100.64.0.0/10
}

// NewGeoLiteService initializes a new GeoLiteService instance and starts a goroutine to update the GeoLite2 City database.
func NewGeoLiteService() *GeoLiteService {
	dbDir := common.EnvConfig.GeoLiteDBPath
	if filepath.Base(dbDir) == "GeoLite2-City.mmdb" {
		// backwards compatible with old meaning of the env var
		dbDir = filepath.Dir(dbDir)
	}
	service := &GeoLiteService{
		cityDBPath: filepath.Join(dbDir, "GeoLite2-City.mmdb"),
		asnDBPath:  filepath.Join(dbDir, "GeoLite2-ASN.mmdb"),
	}

	go func() {
		if err := service.updateDatabase("GeoLite2-City", service.cityDBPath); err != nil {
			log.Printf("Failed to update GeoLite2 City database: %v\n", err)
		}
		if err := service.updateDatabase("GeoLite2-ASN", service.asnDBPath); err != nil {
			log.Printf("Failed to update GeoLite2 ASN database: %v\n", err)
		}
	}()

	return service
}

// GetLocationByIP returns the country and city of the given IP address.
func (s *GeoLiteService) GetLocationByIP(ipAddress string) (country, city, isp string, as_number uint, err error) {
	// Check the IP address against known private IP ranges
	if ip := net.ParseIP(ipAddress); ip != nil {
		for _, ipNet := range tailscaleIPNets {
			if ipNet.Contains(ip) {
				return "Internal Network", "Tailscale", "", 0, nil
			}
		}
		for _, ipNet := range privateLanIPNets {
			if ipNet.Contains(ip) {
				return "Internal Network", "LAN/Docker/k8s", "", 0, nil
			}
		}
		for _, ipNet := range localhostIPNets {
			if ipNet.Contains(ip) {
				return "Internal Network", "localhost", "", 0, nil
			}
		}
	}

	// Race condition between reading and writing the database.
	s.mutex.Lock()
	defer s.mutex.Unlock()

	dbCity, err := maxminddb.Open(s.cityDBPath)
	if err != nil {
		return "", "", "", 0, err
	}
	defer dbCity.Close()

	addr := netip.MustParseAddr(ipAddress)

	var cityRecord struct {
		City struct {
			Names map[string]string `maxminddb:"names"`
		} `maxminddb:"city"`
		Country struct {
			Names map[string]string `maxminddb:"names"`
		} `maxminddb:"country"`
	}

	err = dbCity.Lookup(addr).Decode(&cityRecord)
	if err != nil {
		return "", "", "", 0, err
	}

	// Open the ASN database
	dbASN, err := maxminddb.Open(s.asnDBPath)
	if err != nil {
		return "", "", "", 0, err
	}
	defer dbASN.Close()

	var asnRecord struct {
		Name   string `maxminddb:"autonomous_system_organization"`
		Number uint   `maxminddb:"autonomous_system_number"`
	}

	err = dbASN.Lookup(addr).Decode(&asnRecord)
	if err != nil {
		return cityRecord.Country.Names["en"], cityRecord.City.Names["en"], "", 0, err
	}

	return cityRecord.Country.Names["en"], cityRecord.City.Names["en"], asnRecord.Name, asnRecord.Number, nil
}

// UpdateDatabase checks the age of the database and updates it if it's older than 14 days.
func (s *GeoLiteService) updateDatabase(editionID, dbPath string) error {
	if s.isDatabaseUpToDate(dbPath) {
		log.Printf("%s database is up-to-date.\n", editionID)
		return nil
	}

	log.Printf("Updating %s database...\n", editionID)

	// Download and extract the database
	downloadUrl := fmt.Sprintf(
		"https://download.maxmind.com/app/geoip_download?edition_id=%s&license_key=%s&suffix=tar.gz",
		editionID,
		common.EnvConfig.MaxMindLicenseKey,
	)
	// Download the database tar.gz file
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return fmt.Errorf("failed to download database: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download database, received HTTP %d", resp.StatusCode)
	}

	// Extract the database file directly to the target path
	if err := s.extractDatabase(resp.Body, dbPath); err != nil {
		return fmt.Errorf("failed to extract database: %w", err)
	}

	log.Printf("%s database successfully updated.\n", editionID)
	return nil
}

// isDatabaseUpToDate checks if the database file is older than 14 days.
func (s *GeoLiteService) isDatabaseUpToDate(dbPath string) bool {
	info, err := os.Stat(dbPath)
	if err != nil {
		// If the file doesn't exist, treat it as not up-to-date
		return false
	}
	return time.Since(info.ModTime()) < 14*24*time.Hour
}

// extractDatabase extracts the database file from the tar.gz archive directly to the target location.
func (s *GeoLiteService) extractDatabase(reader io.Reader, dbPath string) error {
	gzr, err := gzip.NewReader(reader)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	tarReader := tar.NewReader(gzr)

	// Iterate over the files in the tar archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar archive: %w", err)
		}

		if header.Typeflag == tar.TypeReg && filepath.Base(header.Name) == filepath.Base(dbPath) {
			// extract to a temporary file to avoid having a corrupted db in case of write failure.
			baseDir := filepath.Dir(common.EnvConfig.GeoLiteDBPath)
			tmpFile, err := os.CreateTemp(baseDir, "geolite.*.mmdb.tmp")
			if err != nil {
				return fmt.Errorf("failed to create temporary database file: %w", err)
			}
			tempName := tmpFile.Name()

			// Write the file contents directly to the target location
			if _, err := io.Copy(tmpFile, tarReader); err != nil {
				// if fails to write, then cleanup and throw an error
				tmpFile.Close()
				os.Remove(tempName)
				return fmt.Errorf("failed to write database file: %w", err)
			}
			tmpFile.Close()

			// ensure the database is not corrupted
			db, err := maxminddb.Open(tempName)
			if err != nil {
				// if fails to write, then cleanup and throw an error
				os.Remove(tempName)
				return fmt.Errorf("failed to open downloaded database file: %w", err)
			}
			db.Close()

			// ensure we lock the structure before we overwrite the database
			// to prevent race conditions between reading and writing the mmdb.
			s.mutex.Lock()
			// replace the old file with the new file
			err = os.Rename(tempName, common.EnvConfig.GeoLiteDBPath)
			s.mutex.Unlock()

			if err != nil {
				// if cannot overwrite via rename, then cleanup and throw an error
				os.Remove(tempName)
				return fmt.Errorf("failed to replace database file: %w", err)
			}
			return nil
		}
	}

	return fmt.Errorf("%s not found in archive", filepath.Base(dbPath))
}
