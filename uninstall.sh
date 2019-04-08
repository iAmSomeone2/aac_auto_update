#!/bin/bash

rm -f /usr/loca/bin/aacautoupdate
systemctl stop aacautoupdate && systemctl disable aacautoupdate
rm -r /etc/systemd/system/aacautoupdate.service
echo aacautoupdate disabled and removed.