package bootstrap

import (
	"github.com/stonith404/pocket-id/backend/internal/common"
	"github.com/stonith404/pocket-id/backend/internal/utils"
	"log"
	"os"
)

func initApplicationImages() {
	dirPath := common.EnvConfig.UploadPath + "/application-images"

	files, err := os.ReadDir(dirPath)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error reading directory: %v", err)
	}

	// Skip if files already exist
	if len(files) > 1 {
		return
	}

	// Copy files from source to destination
	err = utils.CopyDirectory("./images", dirPath)
	if err != nil {
		log.Fatalf("Error copying directory: %v", err)
	}
}
