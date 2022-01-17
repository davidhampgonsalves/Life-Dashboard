#!/bin/sh

enable_wifi() {
  lipc-set-prop com.lab126.cmd wirelessEnable 1
  while ! lipc-get-prop com.lab126.wifid cmState | grep -q CONNECTED; do sleep 1; done
}
disable_wifi() { lipc-set-prop com.lab126.cmd wirelessEnable 0; }

RTC=/sys/devices/platform/mxc_rtc.0/wakeup_enable
rtc_sleep() {
  duration=$1
  [ "$(cat "$RTC")" -eq 0 ] && echo -n "$duration" >"$RTC"
  echo "mem" >/sys/power/state
}

cd /
/usr/sbin/mntroot rw

echo "setting up low power usage"
/etc/init.d/framework stop
echo powersave >/sys/devices/system/cpu/cpu0/cpufreq/scaling_governor
lipc-set-prop com.lab126.powerd preventScreenSaver 1

while true; do
  echo "enabling wifi"
  enable_wifi
  echo "running main"
  
  if ./main ; then
   echo "image generated" 
  else
    # dashboard isn't functioning, try and get back to sane configuration
    # lipc-set-prop com.lab126.powerd preventScreenSaver 0
    # echo ondemand >/sys/devices/system/cpu/cpu0/cpufreq/scaling_governor
    # /etc/init.d/framework start
    shutdown -r now
    exit
  fi

  echo "disabling wifi"
  disable_wifi

  echo "drawing image"
  eips -f -g /image.png

  sleep 10

  rtc_sleep 3600
done
