# GoPiCam

GoPiCam is a web-based interface for managing and controlling a Raspberry Pi camera using RaspiMJPEG, inspired by the RPi-Cam-Web-Interface project. Unlike the original, GoPiCam is implemented in Golang, which eliminates the need to install additional software like PHP or a database on the Raspberry Pi.

## Features

- Web-based interface for Raspberry Pi camera control
- Secure HTTP/HTTPS server support
- User authentication and session management
- Camera preview, start/stop recording, motion detection, and timelapse functionality
- Configuration management

## Installation

1. Clone the repository:

```sh
git clone https://github.com/jempe/gopicam.git
cd gopicam
```

2. Build the project:
```sh
go build -o bin/gopicam cmd/gopicam/main.go
```


3. Run the application:
```sh
./bin/gopicam
```

## Configuration

The configuration files are located in the default path `~/.gopicam`. You can specify a different path using the `-config` flag.

```sh
./bin/gopicam -config /path/to/config
```

## Usage

### Flags

- `-config`:  Define the path of the config folder
- `-reset`:  Reset admin username and password
- `-help`:  Show help
- `-insecure`:  Run web server without HTTPS
- `-port`:  Web server port (default: 443)
- `-debug`:  Print all debug messages

## Admin Account

If the admin username and password are not set, or if you use the `-reset` flag, you will be prompted to create an admin account.

## Running the Server

To run the server with HTTPS:

```sh
./bin/gopicam
```

To run the server without HTTPS:

```sh
./bin/gopicam -insecure
```

## Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss improvements or new features.


