#!/bin/sh

enable_wifi() {
  lipc-set-prop com.lab126.cmd wirelessEnable 1
  while ! lipc-get-prop com.lab126.wifid cmState | grep -q CONNECTED; do sleep 1; done
}
disable_wifi() { lipc-set-prop com.lab126.cmd wirelessEnable 0; }

cd /
/usr/sbin/mntroot rw
/etc/init.d/powerd stop

enable_wifi
./main
disable_wifi

eips -f -g /image.png

echo "sleeping for 30 minutes"
sleep 1850
