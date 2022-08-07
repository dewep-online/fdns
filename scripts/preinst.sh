#!/bin/bash

if test -f "/lib/systemd/system/systemd-resolved.service"; then
    systemctl disable systemd-resolved
    systemctl stop systemd-resolved

    systemctl daemon-reload
fi

if ! [ -d /var/lib/fdns/ ]; then
    mkdir /var/lib/fdns
fi

if test -f "/etc/systemd/system/fdns.service"; then
    systemctl disable fdns
    systemctl stop fdns

    systemctl daemon-reload
fi

