#!/bin/sh
ssh root@192.168.5.11 "/usr/sbin/mntroot rw"

scp start.sh root@192.168.5.11:/start.sh
scp life-dashboard-init root@192.168.5.11:/etc/init.d/life-dashboard-init

ssh root@192.168.5.11 << EOF
	cd /
	chmod 777 start.sh
	chmod 777 /etc/init.d/life-dashboard-init

	echo "Install complete, start script with: `/etc/init.d/life-dashboard-init start` and let it run till the ssh session hangs"
EOF
