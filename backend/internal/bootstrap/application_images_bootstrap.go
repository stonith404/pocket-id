package bootstrap

import (
	"github.com/stonith404/pocket-id/backend/internal/common"
	"github.com/stonith404/pocket-id/backend/internal/utils"
	"log"
	"os"
	"strings"
)

// initApplicationImages copies the images from the images directory to the application-images directory
func initApplicationImages() {
	dirPath := common.EnvConfig.UploadPath + "/application-images"

	sourceFiles, err := os.ReadDir("./images")
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error reading directory: %v", err)
	}

	destinationFiles, err := os.ReadDir(dirPath)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error reading directory: %v", err)
	}

	// Copy images from the images directory to the application-images directory if they don't already exist
	for _, sourceFile := range sourceFiles {
		if sourceFile.IsDir() || imageAlreadyExists(sourceFile.Name(), destinationFiles) {
			continue
		}
		srcFilePath := "./images/" + sourceFile.Name()
		destFilePath := dirPath + "/" + sourceFile.Name()

		err := utils.CopyFile(srcFilePath, destFilePath)
		if err != nil {
			log.Fatalf("Error copying file: %v", err)
		}
	}

}

func imageAlreadyExists(fileName string, destinationFiles []os.DirEntry) bool {
	for _, destinationFile := range destinationFiles {
		sourceFileWithoutExtension := getImageNameWithoutExtension(fileName)
		destinationFileWithoutExtension := getImageNameWithoutExtension(destinationFile.Name())

		if sourceFileWithoutExtension == destinationFileWithoutExtension {
			return true
		}
	}

	return false
}

func getImageNameWithoutExtension(fileName string) string {
	splitted := strings.Split(fileName, ".")
	return strings.Join(splitted[:len(splitted)-1], ".")
}
