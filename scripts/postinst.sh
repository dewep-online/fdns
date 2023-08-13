#!/bin/bash


if [ -f "/etc/systemd/system/fdns.service" ]; then
    systemctl start fdns
    systemctl enable fdns
    systemctl daemon-reload
fi
