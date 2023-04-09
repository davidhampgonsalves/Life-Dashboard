#!/bin/sh
# Run via /etc/init.d/life-dashboard-init start

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

eips -c 12 19 "Starting polling / sleep" && eips 15 20 "cycle in 30 seconds."
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

  echo "writing stats"
  let max_sleep=24*60*60
  let next_refresh=$(($(date -d 23:59:59 +%s) - $(date +%s))) # seconds till midnight
  let next_refresh=$(($next_refresh + (3*60*60)))
  if [ $next_refresh -le 0 ] || [ $next_refresh -ge $max_sleep ]; then next_refresh=$max_sleep; fi
  battery_level=$(lipc-get-prop com.lab126.powerd battLevel)
  let next_refresh_minutes=$next_refresh/60
  eips 2 37 "next $next_refresh_minutes (minutes) b.$battery_level"
  eips 2 38 "$(TZ=UTC+3 date -R "+%a %l:%M")"
  
  echo "entering rtc sleep"
  sleep 5
  echo $next_refresh > /sys/devices/platform/mxc_rtc.0/wakeup_enable
  echo "mem" > /sys/power/state
done