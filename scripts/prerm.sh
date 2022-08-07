#!/bin/bash

if [ -f "/etc/systemd/system/fdns.service" ]; then
    systemctl stop fdns
    systemctl disable fdns

    systemctl daemon-reload
fi