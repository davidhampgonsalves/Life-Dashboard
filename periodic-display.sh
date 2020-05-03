#!/bin/sh

enable_wifi() {
  lipc-set-prop com.lab126.cmd wirelessEnable 1
  while ! lipc-get-prop com.lab126.wifid cmState | grep -q CONNECTED; do sleep 1; done
}
disable_wifi() { lipc-set-prop com.lab126.cmd wirelessEnable 0; }

cd /
/usr/sbin/mntroot rw
echo "disable powerd"
/etc/init.d/powerd stop

echo "enabling wifi"
enable_wifi
echo "running main"
./main
echo "disabling wifi"
disable_wifi

echo "drawing image"
eips -f -g /image.png
