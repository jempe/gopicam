#!/bin/bash

GOOS=linux GOARCH=arm GOARM=7 go build -o bin/gopicam_armv7l cmd/gopicam/main.go
GOOS=linux GOARCH=arm GOARM=6 go build -o bin/gopicam_armv6l cmd/gopicam/main.go

