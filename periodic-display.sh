#!/bin/sh


enable_wifi() {
  lipc-set-prop com.lab126.cmd wirelessEnable 1
  while ! lipc-get-prop com.lab126.wifid cmState | grep -q CONNECTED; do sleep 1; done
}

disable_wifi() { lipc-set-prop com.lab126.cmd wirelessEnable 0; }

eips 'In 1 minute the kindle framework will be stopped and SSH will no longer be running. Act accordingly'
sleep 30
eips 'Tick-Tock, 30 seconds remaining.'
sleep 20
eips '10 seconds remaining.'
sleep 10

/etc/init.d/framework stop
/etc/init.d/powerd stop
/usr/sbin/mntroot rw

while true;
do
  enable_wifi
  ./main
  disable_wifi

  eips -f -g image.png

  echo "sleeping for an hour"
  sleep 3600
done

#suspend() {
  #log "WAITING TO SET ALARM"
  ## if we are active should we toggle power to get to screen saver?
  #if powerd_test -s | grep Active; then powerd_test -p; fi

  ##lipc-wait-event com.lab126.powerd readyToSuspend
  #while ! powerd_test -s | grep -q Screen; do sleep 1; done
  #sleep 5
  #log "DISPLAYING DATE"
  ##eips 10 20 "`date`"
  #eips -g tmp.png
  #sleep 2

  #while ! powerd_test -s | grep -q Ready; do sleep 1; done
  #log "READY TO SUSPEND"
  #lipc-set-prop -i com.lab126.powerd rtcWakeup 600
  #echo mem > /sys/power/state
#}

#log "STARTING INIT SLEEP"

##/etc/init.d/framework stop
#echo powersave > /sys/devices/system/cpu/cpu0/cpufreq/scaling_governor
#mntroot rw

##while true;
  ##do
  #log "STARTING LOOP"
  #./main
  ##./pngcrush -bit_depth 4 -c 0 tmp.png out.png
  #eips -f -g image.png


  #lipc-wait-event com.lab126.powerd goingToScreenSaver
  #log "Screen Saver Event"
  #sleep 5
  #./main
  #./pngcrush -bit_depth 4 -c 0 tmp.png out.png
  #eips -f -g out.png

  ##log "STARTING LOOP"

  ##rm /tmp.png
  ##enable_wifi
  ##ntpdate pool.ntp.org
  ##./main
  ##disable_wifi
  ##./pngcrush -bit_depth 4 -c 0 tmp.png out.png

  ##rm -f /opt/amazon/screen_saver/600x800/*

  ##lipc-wait-event com.lab126.powerd goingToScreenSaver

  ##log "DISPLAYING DATE"
  ###eips 10 20 "`date`"
  ##sleep 5
  ##eips -f -g out.png
  ##sleep 2

  ##lipc-wait-event com.lab126.powerd readyToSuspend

  ##log "READY TO SUSPEND"
  ##lipc-set-prop -i com.lab126.powerd rtcWakeup 120 2>> /var/log/periodic-display
  ##cat /proc/driver/rtc 2>> /var/log/periodic-display
  ###log "END OF LOG ATTEMPT"
  ###sleep 5
  ##WAKEUPTIMER=$(( `date +%s` + 60 ))
  ##echo 1 > /sys/class/rtc/rtc0/device/wakeup_enable
  ##echo 1 > /sys/class/rtc/rtc0/device/wakeup_from_halt
  ##echo 0 > /sys/class/rtc/rtc0/wakealarm
  ##echo $WAKEUPTIMER > /sys/class/rtc/rtc0/wakealarm
  ##echo 0 > /sys/devices/platform/pmic_rtc.1/rtc/rtc1/wakealarm
  ##echo $WAKEUPTIMER > /sys/devices/platform/pmic_rtc.1/rtc/rtc1/wakealarm
  ###echo "mem" > /sys/power/state

  ##lipc-wait-event -mt com.lab126.powerd '*'
##done



