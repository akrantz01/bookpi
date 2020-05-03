#!/bin/bash

# Catch errors
trap "echo 'Failed to install BookPi. Please add an issue to https://github.com/akrantz01/bookpi with the command output.'; exit" ERR

echo "Installing BookPi to the system..."

# Check if running as root
if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

# Upgrade the system
echo "Upgrading the system..."
apt-get update && apt-get upgrade -y

# Create user
echo "Creating bookpi user..."
adduser --system --shell /bin/bash --gecos 'BookPi' --group --disabled-password --home /opt/bookpi bookpi

# Change hostname
echo "Changing hostname to bookpi..."
raspi-config nonint do_hostname bookpi
sed -i "s/127.0.1.1.*raspberrypi//g" /etc/hosts

# Create environment file
echo "Setting up bookpi binaries..."
cat << EOF > /opt/bookpi/environment
HOST=0.0.0.0
PORT=80
DATABASE=/opt/bookpi/database.db
FILES_DIR=/opt/bookpi/files
RESET=no
EOF

# Install Pillow dependencies
echo "Installing dependencies for Pillow..."
apt-get install -y libtiff5 libopenjp2-7

# Move binaries to /usr/local/bin
mv bookpi-*-server /usr/local/bin/bookpi-server
mv bookpi-*-display /usr/local/bin/bookpi-display
mv bookpi.sh /usr/local/bin/bookpi
chmod +x /usr/local/bin/bookpi /usr/local/bin/bookpi-server /usr/local/bin/bookpi-display

# Install systemd units
echo "Installing services as systemd units..."
mv display.service /etc/systemd/system/bookpi-display.service
mv server.service /etc/systemd/system/bookpi-server.service

# Reload systemd
systemctl enable bookpi-display.service
systemctl enable bookpi-server.service

# Start services on boot
systemctl start bookpi-display.service
systemctl start bookpi-server.service

# Install hostapd
echo "Installing hostapd and dnsmasq..."
apt-get install hostapd dnsmasq -y

# Configure DHCP
echo "Configuring DHCP..."
cat << EOF >> /etc/dhcpcd.conf
interface wlan0
static ip_address=10.5.1.1
nohook wpa_supplicant
denyinterfaces wlan0
EOF
systemctl restart dhcpcd.service

# Configure access point
echo "Configuring access point..."
apPassword="bookpi-$(openssl rand -hex 10)"
cat << EOF > /etc/hostapd/hostapd.conf
interface=wlan0
hw_mode=g
channel=7
wmm_enabled=0
macaddr_acl=0
auth_algs=1
ignore_broadcast_ssid=0
wpa=2
wpa_key_mgmt=WPA-PSK
wpa_pairwise=TKIP
rsn_pairwise=CCMP
ssid=BookPi
wpa_passphrase=${apPassword}
EOF
echo "DAEMON_CONF=\"/etc/hostapd/hostapd.conf\""
systemctl unmask hostapd.service
systemctl enable hostapd.service
systemctl restart hostapd.service

# Configure DNS server
echo "Configuring DNS server..."
cat << EOF > /etc/dnsmasq.conf
interface=wlan0
dhcp-range=10.5.1.5,10.5.1.128,255.255.255.0,24h

no-resolv
domain-needed
bogus-priv

expand-hosts
domain=pi
EOF
printf "10.5.1.1\tbook" >> /etc/hosts
systemctl restart dnsmasq.service

# Done, reboot
echo "Successfully configured BookPi"
echo "    - A wireless network named 'BookPi' has been created"
echo "    - The password for the network is '${apPassword}'"
echo "    - It is recommended that you change your password"
echo "Automatically restarting to complete setup in 10 seconds..."
echo "Press Ctrl+c to cancel"
sleep 10
reboot
