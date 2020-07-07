#!/bin/bash

if [[ $EUID -ne 0 ]]; then
	echo "You must be root to run this script." 1>&2
	exit 100
fi

# install packages
apt install hostapd dnsmasq netfilter-persistent iptables-persistent

# enable hostapd service
systemctl unmask hostapd
systemctl enable hostapd

#interface wlan0
#    static ip_address=192.168.4.1/24
#    nohook wpa_supplicant

# /etc/sysctl.d/routed-ap.conf
# https://www.raspberrypi.org/documentation/configuration/wireless/access-point-routed.md
# Enable IPv4 routing
net.ipv4.ip_forward=1


iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE

netfilter-persistent save

# /etc/dnsmasq.conf
listen-address=192.168.12.1
bind-dynamic
dhcp-range=192.168.12.1,192.168.12.254,255.255.255.0,24h
dhcp-option-force=option:router,192.168.12.1
dhcp-option-force=option:dns-server,192.168.12.1
no-hosts

rfkill unblock wlan

#/etc/hostapd/hostapd.conf
beacon_int=100
ssid=MyAccessPoint
interface=wlan0
driver=nl80211
channel=1
ctrl_interface=/tmp/create_ap.wlan0.conf.BVxgAOuz/hostapd_ctrl
ctrl_interface_group=0
ignore_broadcast_ssid=0
ap_isolate=0
hw_mode=g
wpa=3
wpa_passphrase=MyPassPhrase
wpa_key_mgmt=WPA-PSK
wpa_pairwise=TKIP CCMP
rsn_pairwise=CCMP
