#!/bin/sh

ssh root@192.168.15.244 "/usr/sbin/mntroot rw"
scp periodic-display-init root@192.168.15.244:/etc/init.d/periodic-display-init
scp target/arm-unknown-linux-musleabi/release/life-dashboard root@192.168.15.244:/main
ssh root@192.168.15.244 << EOF
	chmod 777 /etc/init.d/periodic-display-init
	chmod 777 /main
	ln -sf /etc/init.d/periodic-display-init /etc/rc5.d/S97periodic-display
	shutdown -r now
	echo "Install complete, unplug Kindle now (or it will mess up the USB connection)."
EOF
