package camera

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"

	"github.com/jempe/gopicam/pkg/utils"
)

const previewFolder = "/dev/shm/mjpeg"
const configFile = "/etc/raspimjpeg"

type CamController struct {
	ConfigFolder string
}

// Prepare everything to run raspimjpeg
func (camController *CamController) Init() error {
	if !utils.Exists(previewFolder) {
		// create a directory that is available for current user only
		createDirErr := os.MkdirAll(previewFolder, 0700)

		if createDirErr != nil {
			return errors.New("Error: Couldn't create preview folder: " + previewFolder)
		}

	} else if utils.Exists(previewFolder) && !utils.IsDirectory(previewFolder) {
		return errors.New("Error: The preview path is not a folder: " + previewFolder)
	}

	return nil
}

// Check the Preview Path and return the image as a base64 string
func (camController *CamController) GetPreview() (previewImage string, status string, err error) {
	previewImagePath := previewFolder + "/cam.jpg"
	statusPath := previewFolder + "/status_mjpeg.txt"

	status = "error"

	// check if Preview Image exists
	if !utils.Exists(previewImagePath) {
		err = errors.New("Error: The preview file doesn't exist: " + previewImagePath)
		return
	}

	// check if Status Text file exists
	if !utils.Exists(statusPath) {
		err = errors.New("Error: The status file doesn't exist: " + statusPath)
		return
	}

	// read image file content
	var imageBuffer bytes.Buffer

	imageFile, imageFileErr := os.Open(previewImagePath)
	if imageFileErr != nil {
		err = imageFileErr
		return
	}

	_, readImageErr := imageBuffer.ReadFrom(imageFile)
	if readImageErr != nil {
		err = readImageErr
		return
	}

	previewImage = utils.Base64Encode(imageBuffer.Bytes())

	imageFile.Close()

	// read status file content
	statusContent, statusFileErr := ioutil.ReadFile(statusPath)
	if statusFileErr != nil {
		err = statusFileErr
		return
	}

	status = string(statusContent)

	return
}
