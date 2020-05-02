#!/bin/bash

# Check if running as root
if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi


display_help() {
  read -r -d '' display << EOH
Usage: $0 [OPTIONS]

  -h|--help       Display this message
  -p|--password   Change the wireless password
  -s|--ssid       Change the broadcasted name
  --start         Start all the services
  --stop          Stop all the services
  --restart       Restart all the services
EOH
  echo "$display"
}

change_wifi_password() {
  sed -i "s/wpa_passphrase=.*/wpa_passphrase=$1/g" /etc/hostapd/hostapd.conf
  systemctl restart hostapd
}

change_wifi_ssid() {
  sed -i "s/ssid=.*/ssid=$1/g" /etc/hostapd/hostapd.conf
  systemctl restart hostapd
}

stop_services() {
  systemctl stop bookpi-display
  systemctl stop bookpi-server
  systemctl stop hostapd
  systemctl stop dnsmasq
  systemctl stop dhcpcd
}

start_services() {
  systemctl start bookpi-display
  systemctl start bookpi-server
  systemctl start hostapd
  systemctl start dnsmasq
  systemctl start dhcpcd
}

restart_services() {
  systemctl restart bookpi-display
  systemctl restart bookpi-server
  systemctl restart hostapd
  systemctl restart dnsmasq
  systemctl restart dhcpcd
}

if [ $# -eq 0 ]
then
	display_help
	exit 1
fi

while [[ $# -gt 0 ]]
do
case "$1" in
  -p|--password)
  change_wifi_password "$2"
  shift
  shift
	;;
  -s|--ssid)
  change_wifi_ssid "$2"
  shift
  shift
  ;;
  --start)
  start_services
  shift
  ;;
  --stop)
  stop_services
  shift
  ;;
  --restart)
  restart_services
  shift
  ;;
  -h|--help)
  display_help
  exit 0
  ;;
  *)
  echo "Unknown option: $1"
  display_help
  exit 1
  ;;
esac
done
