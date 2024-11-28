package service

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/netip"
	"os"
	"path/filepath"
	"time"

	"github.com/oschwald/maxminddb-golang/v2"

	"github.com/stonith404/pocket-id/backend/internal/common"
)

type GeoLiteService struct{}

// NewGeoLiteService initializes a new GeoLiteService instance and starts a goroutine to update the GeoLite2 City database.
func NewGeoLiteService() *GeoLiteService {
	service := &GeoLiteService{}

	go func() {
		if err := service.updateDatabase(); err != nil {
			log.Printf("Failed to update GeoLite2 City database: %v\n", err)
		}
	}()

	return service
}

// GetLocationByIP returns the country and city of the given IP address.
func (s *GeoLiteService) GetLocationByIP(ipAddress string) (country, city string, err error) {
	// Check if IP is in Tailscale's CGNAT range (100.64.0.0/10)
	if ip := net.ParseIP(ipAddress); ip != nil {
		if ip.To4() != nil && ip.To4()[0] == 100 && ip.To4()[1] >= 64 && ip.To4()[1] <= 127 {
			return "Internal Network", "Tailscale", nil
		}
	}

	db, err := maxminddb.Open(common.EnvConfig.GeoLiteDBPath)
	if err != nil {
		return "", "", err
	}
	defer db.Close()

	addr := netip.MustParseAddr(ipAddress)

	var record struct {
		City struct {
			Names map[string]string `maxminddb:"names"`
		} `maxminddb:"city"`
		Country struct {
			Names map[string]string `maxminddb:"names"`
		} `maxminddb:"country"`
	}

	err = db.Lookup(addr).Decode(&record)
	if err != nil {
		return "", "", err
	}

	return record.Country.Names["en"], record.City.Names["en"], nil
}

// UpdateDatabase checks the age of the database and updates it if it's older than 14 days.
func (s *GeoLiteService) updateDatabase() error {
	if s.isDatabaseUpToDate() {
		log.Println("GeoLite2 City database is up-to-date.")
		return nil
	}

	log.Println("Updating GeoLite2 City database...")

	// Download and extract the database
	downloadUrl := fmt.Sprintf(
		"https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=%s&suffix=tar.gz",
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
	if err := s.extractDatabase(resp.Body); err != nil {
		return fmt.Errorf("failed to extract database: %w", err)
	}

	log.Println("GeoLite2 City database successfully updated.")
	return nil
}

// isDatabaseUpToDate checks if the database file is older than 14 days.
func (s *GeoLiteService) isDatabaseUpToDate() bool {
	info, err := os.Stat(common.EnvConfig.GeoLiteDBPath)
	if err != nil {
		// If the file doesn't exist, treat it as not up-to-date
		return false
	}
	return time.Since(info.ModTime()) < 14*24*time.Hour
}

// extractDatabase extracts the database file from the tar.gz archive directly to the target location.
func (s *GeoLiteService) extractDatabase(reader io.Reader) error {
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

		// Check if the file is the GeoLite2-City.mmdb file
		if header.Typeflag == tar.TypeReg && filepath.Base(header.Name) == "GeoLite2-City.mmdb" {
			outFile, err := os.Create(common.EnvConfig.GeoLiteDBPath)
			if err != nil {
				return fmt.Errorf("failed to create target database file: %w", err)
			}
			defer outFile.Close()

			// Write the file contents directly to the target location
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("failed to write database file: %w", err)
			}
			return nil
		}
	}

	return errors.New("GeoLite2-City.mmdb not found in archive")
}
