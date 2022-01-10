#!/bin/sh

enable_wifi() {
  lipc-set-prop com.lab126.cmd wirelessEnable 1
  while ! lipc-get-prop com.lab126.wifid cmState | grep -q CONNECTED; do sleep 1; done
}
disable_wifi() { lipc-set-prop com.lab126.cmd wirelessEnable 0; }

rtc_sleep() {
  duration=$1
  [ "$(cat "$RTC")" -eq 0 ] && echo -n "$duration" >"$RTC"
  echo "mem" >/sys/power/state
}

cd /
/usr/sbin/mntroot rw

echo "setting up low power usage"
/etc/init.d/framework stop
initctl stop webreader >/dev/null 2>&1
echo powersave >/sys/devices/system/cpu/cpu0/cpufreq/scaling_governor
lipc-set-prop com.lab126.powerd preventScreenSaver 1

while true; do
  echo "enabling wifi"
  enable_wifi
  echo "running main"
  ./main
  echo "disabling wifi"
  disable_wifi

  echo "drawing image"
  eips -f -g /image.png

  sleep 10

  rtc_sleep 3600
done
