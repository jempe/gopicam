#!/bin/bash

GOARCH=arm GOARM=7 go build github.com/jempe/gopicam/cmd/gopicam

scp gopicam webcam@gopicam.local:/home/webcam/bin/
