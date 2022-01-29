#!/bin/sh

enable_wifi() {
  lipc-set-prop com.lab126.cmd wirelessEnable 1
  while ! lipc-get-prop com.lab126.wifid cmState | grep -q CONNECTED; do sleep 1; done
}
disable_wifi() { lipc-set-prop com.lab126.cmd wirelessEnable 0; }

cd /
/usr/sbin/mntroot rw

echo "setting up low power usage"
/etc/init.d/framework stop
echo powersave >/sys/devices/system/cpu/cpu0/cpufreq/scaling_governor
lipc-set-prop com.lab126.powerd preventScreenSaver 1

eip -c 12 19 "Starting polling / sleep" && eips 15 20 "cycle in 30 seconds."
sleep 30

while true; do
  echo "enabling wifi"
  enable_wifi
  echo "running main"
  
  if ./main ; then
   echo "image generated" 
  else
    # dashboard isn't functioning, try and get back to a sane configuration
    shutdown -r now
    exit
  fi

  echo "disabling wifi"
  disable_wifi

  echo "drawing image"
  eips -f -g /image.png

  sleep 1
  batteryLevel=$(lipc-get-prop com.lab126.powerd battLevel)
  if [ $batteryLevel -le 10 ]; then
    eips 46 38 "$batteryLevel"
  fi

  echo "entering rtc sleep"
  sleep 5
  echo 86400 > /sys/devices/platform/mxc_rtc.0/wakeup_enable
  echo "mem" > /sys/power/state
done
