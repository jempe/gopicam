#!/bin/bash

GOPICAM_CONFIG_FOLDER=$( dirname "${BASH_SOURCE[0]}" )

openssl genrsa -out $GOPICAM_CONFIG_FOLDER/server.key 2048

openssl req -new -x509 -sha256 -key $GOPICAM_CONFIG_FOLDER/server.key -out $GOPICAM_CONFIG_FOLDER/server.crt -days 3650

