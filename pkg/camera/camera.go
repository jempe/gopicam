package camera

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

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

// Start raspimjpeg
func (camController *CamController) StartRaspiMJPEG() {
	camController.KillRaspiMJPEG()

	cmd := exec.Command(camController.ConfigFolder + "/bin/raspimjpeg")
	raspiMJPEGOutput, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	for {
		scanner := bufio.NewScanner(raspiMJPEGOutput)
		for scanner.Scan() {
			fmt.Println(scanner.Text()) // Println will add back the final '\n'
		}
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

// Kill raspimjpeg
func (camController *CamController) KillRaspiMJPEG() {
	cmd := exec.Command("ps")
	log.Printf("Search RaspiMJPEG processes")
	stdoutStderr, err := cmd.CombinedOutput()

	if err == nil {
		psLineList := strings.Split(string(stdoutStderr), "\n")

		for _, psLine := range psLineList {
			if strings.Contains(psLine, "raspimjpeg") {
				pidRegex := regexp.MustCompile(`[0-9]*`)
				pid := pidRegex.FindString(strings.TrimLeft(psLine, " "))

				killCmd := exec.Command("kill", "-9", pid)
				log.Println("Killing process", pid)
				err := killCmd.Run()
				if err != nil {
					log.Println("Command finished with error:", err)
				}
			}
		}
	}
}

// Kill raspimjpeg
func (camController *CamController) ReadFIFO() {
	fifoMessage, err := os.OpenFile(camController.ConfigFolder+"/fifos/FIFO1", os.O_RDONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Reading FIFO")
	var fifoBuffer bytes.Buffer

	for {
		_, fifoErr := io.Copy(&fifoBuffer, fifoMessage)

		if fifoErr != nil {
			log.Fatal(fifoErr)
			return
		}

		if fifoBuffer.Len() > 0 {
			fmt.Println("FIFO Message:", fifoBuffer.String())
			fifoBuffer.Reset()
		}
		time.Sleep(100 * time.Millisecond)
	}

	fifoMessage.Close()
}
