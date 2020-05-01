#!/bin/bash

# Catch errors
trap "echo 'Failed to install BookPi. Please add an issue to https://github.com/akrantz01/bookpi with the command output.'; exit" ERR

echo "Installing BookPi to the system..."
echo "NOTICE: sudo is used to copy files to your PATH and to add systemd unit files. When prompted, please enter your password."

# Create user
adduser --system --shell /bin/bash --gecos 'BookPi' --group --disabled-password --home /opt/bookpi bookpi

# Create environment file
cat << EOF > /opt/bookpi/environment
HOST=0.0.0.0
PORT=80
DATABASE=/opt/bookpi/database.db
FILES_DIR=/opt/bookpi/files
RESET=no
EOF

# Move binaries to /usr/local/bin
sudo mv bookpi-*-server /usr/local/bin
sudo mv bookpi-*-display /usr/local/bin

# Install systemd units
sudo mv display.service /etc/systemd/system/bookpi-display.service
sudo mv server.service /etc/systemd/system/bookpi-server.service

# Reload systemd
sudo systemctl enable bookpi-display.service
sudo systemctl enable bookpi-server.service

# Start services on boot
sudo systemctl start bookpi-display.service
sudo systemctl start bookpi-server.service
