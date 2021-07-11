#!/bin/bash

sudo apt install supervisor

sudo adduser webcam

sudo adduser webcam video

echo 'net.ipv4.ip_unprivileged_port_start=443' > ~/50-unprivileged-ports.conf
sudo chown root:root ~/50-unprivileged-ports.conf
sudo mv ~/50-unprivileged-ports.conf /etc/sysctl.d/50-unprivileged-ports.conf

sudo sysctl --system


sudo cp shell_scripts/raspimjpeg /etc/

sudo su webcam -c 'mkdir /home/webcam/bin'

sudo su webcam -c 'cp bin/gopicam /home/webcam/bin'

sudo cp shell_scripts/gopicam.conf /etc/supervisor/conf.d/

sudo su webcam -c '/home/webcam/bin/gopicam'

sudo su webcam -c '/home/webcam/.gopicam/generate_cert.sh'

sudo su webcam -c 'mkdir /home/webcam/.gopicam/bin'

sudo su webcam -c 'cp bin/raspimjpeg /home/webcam/.gopicam/bin'

sudo supervisorctl reread

sudo supervisorctl reload
