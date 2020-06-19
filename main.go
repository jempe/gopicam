package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/jempe/gopicam/db"
	"github.com/jempe/gopicam/utils"
)

var configPathFlag = flag.String("config", "", "Define the path of config folder")

type CamServer struct {
	Db *db.DB
}

func main() {
	flag.Parse()

	var configPath string

	// Check if there is a config flag, unless use the default location
	if *configPathFlag == "" {
		configPath = os.Getenv("HOME") + "/.gopicam"
	} else {
		// use the config path from the argument and remove any trailing slash
		configPath = strings.TrimRight(*configPathFlag, "/")
	}

	if !utils.Exists(configPath) {
		// create a directory that is available for current user only
		createDirErr := os.MkdirAll(configPath, 0700)

		if createDirErr != nil {
			logAndExit("Error: Couldn't create the configuration folder " + configPath)
		}

	} else if utils.Exists(configPath) && !utils.IsDirectory(configPath) {
		logAndExit("Error: The config path " + configPath + " is not a folder")
	}

	dbPath := configPath + "/gopicam.db"

	database := &db.DB{Path: dbPath}

	err := database.InitDb()

	if err != nil {
		logAndExit("Couldn't create the DB")
	}
}

// Print error message and exit
func logAndExit(message string) {
	log.Println(message)
	os.Exit(1)
}
