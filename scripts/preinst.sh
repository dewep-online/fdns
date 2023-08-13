#!/bin/bash


if ! [ -d /var/lib/fdns/ ]; then
    mkdir /var/lib/fdns
fi

if [ -f "/etc/systemd/system/fdns.service" ]; then
    systemctl stop fdns
    systemctl disable fdns
    systemctl daemon-reload
fi
