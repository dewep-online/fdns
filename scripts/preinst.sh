#!/bin/bash

if test -f "/etc/systemd/system/fdns.service"; then
    systemctl stop fdns
    systemctl disable fdns

    systemctl daemon-reload
    systemctl reset-failed
fi

