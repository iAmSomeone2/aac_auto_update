#!/bin/bash

cp ./aacautoupdate /usr/local/bin/
cp ./aacautoupdate.service /etc/systemd/system/
mkdir -p /var/www/.cache && chown -hR www-data:www-data /var/www/.cache
mkdir -p /var/www/cell.bdavidson.dev/html/data
chown -R www-data:www-data /var/www/cell.bdavidson.dev/html/data
chmod -R 764 /var/www/cell.bdavidson.dev/html/data
systemctl enable aacautoupdate && systemctl start aacautoupdate
echo aacautoupdate installed and active!