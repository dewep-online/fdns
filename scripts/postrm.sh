#!/bin/bash

if test -f "/lib/systemd/system/systemd-resolved.service"; then
    systemctl enable systemd-resolved
    systemctl start systemd-resolved

    systemctl daemon-reload
fi