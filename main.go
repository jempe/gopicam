package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs"
	"golang.org/x/crypto/bcrypt"

	"github.com/jempe/gopicam/db"
	"github.com/jempe/gopicam/handlers"
	"github.com/jempe/gopicam/utils"
	"github.com/jempe/gopicam/validator"
)

var configPathFlag = flag.String("config", "", "Define the path of config folder")
var resetAdmin = flag.Bool("reset", false, "Reset admin username and password")
var showHelp = flag.Bool("help", false, "Show Help")
var insecureServer = flag.Bool("insecure", false, "Run web server without HTTPS")
var port = flag.Int("port", 443, "Web Server Port")

func main() {
	flag.Parse()

	// if help argument is present show flag Defaults
	if *showHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

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

	sessionManager := scs.New()
	sessionManager.IdleTimeout = 20 * time.Minute
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Secure = !*insecureServer

	if err != nil {
		logAndExit("Couldn't create the DB")
	}

	// Check if configuration have the admin username and password or reset argument is present
	if database.GetConfigValue("username") == nil || database.GetConfigValue("password") == nil || *resetAdmin {
		//Ask Username
		var username string
		for {
			usernameInput, shellErr := simpleShell("Enter a username to create the admin account")

			if shellErr != nil {
				logAndExit(shellErr.Error())
			}

			validUsername, usernameError := validator.ValidateUsername(usernameInput)

			// check if username is valid to continue
			if validUsername {
				username = usernameInput
				break
			} else {
				fmt.Println("Error:", usernameError)
			}
		}

		//Ask Password
		password, shellErr := simpleShell("Enter the password of the admin account")

		if shellErr != nil {
			logAndExit(shellErr.Error())
		}

		//Save Username
		saveErr := database.SetConfigValue("username", []byte(username))
		if saveErr != nil {
			logAndExit(saveErr.Error())
		}

		//Hash and Save Password
		hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(password), 8)
		if hashedPasswordErr != nil {
			logAndExit(hashedPasswordErr.Error())
		}

		saveErr = database.SetConfigValue("password", hashedPassword)
		if saveErr != nil {
			logAndExit(saveErr.Error())
		}

		fmt.Println("Creating admin account", username, "with password", password)

		// Save Username and password
	}

	srv := &handlers.Server{Db: database, Sessions: sessionManager}

	// Handler to serve HTML Files
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("html/")))
	mux.HandleFunc("/api/login", srv.LoginHandler)

	// Setup Web Server

	serverPort := strconv.Itoa(*port)
	serverProtocol := "https"
	serverCertFile := configPath + "/server.crt"
	serverKeyFile := configPath + "/server.key"

	if *insecureServer {
		serverProtocol = "http"
	} else {
		// Check if certificate and key files exist

		if !utils.Exists(serverCertFile) || !utils.Exists(serverKeyFile) {

			// Create script to generate certificates and print message asking user to execute it
			certificateGenerator := configPath + "/generate_cert.sh"

			generateCertScript := []byte("#!/bin/bash\n\nGOPICAM_CONFIG_FOLDER=$( dirname \"${BASH_SOURCE[0]}\" )\n\nopenssl genrsa -out $GOPICAM_CONFIG_FOLDER/server.key 2048\n\nopenssl req -new -x509 -sha256 -key $GOPICAM_CONFIG_FOLDER/server.key -out $GOPICAM_CONFIG_FOLDER/server.crt -days 3650\n")
			generateScriptErr := ioutil.WriteFile(certificateGenerator, generateCertScript, 0700)

			if generateScriptErr != nil {
				logAndExit(generateScriptErr.Error())
			}

			logAndExit("Certificate or key files don't exist, run the script " + certificateGenerator + " to create them. You can also use the -insecure flag to run the server without HTTPS")
		}
	}

	showLocalIPs(serverPort, serverProtocol)

	//Start Web Server

	if *insecureServer {
		panic(http.ListenAndServe(":"+serverPort, sessionManager.LoadAndSave(mux)))
	} else {
		panic(http.ListenAndServeTLS(":"+serverPort, serverCertFile, serverKeyFile, sessionManager.LoadAndSave(mux)))
	}
}

// Shell Ask Question and return response
//
func simpleShell(question string) (response string, err error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(question, ":")
	input, inputErr := reader.ReadString('\n')

	if inputErr != nil {
		err = inputErr
		return
	}

	response = strings.Replace(input, "\n", "", -1)

	return
}

// Print error message and exit
func logAndExit(message string) {
	log.Println(message)
	os.Exit(1)
}

// Show the URLs where GoPiCam will run
//
func showLocalIPs(port string, protocol string) {
	urls := []string{protocol + "://localhost:" + port}

	localAddresses, err := net.InterfaceAddrs()
	if err == nil {
		for _, address := range localAddresses {
			ip := address.(*net.IPNet)
			if ip.IP.To4() != nil && !ip.IP.IsLoopback() {
				urls = append(urls, protocol+"://"+ip.IP.String()+":"+port)
			}
		}

	}

	log.Println("Running GoPiCam on", urls)
}
