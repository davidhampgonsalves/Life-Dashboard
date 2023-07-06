#!/bin/sh
ssh root@192.168.5.11 "/usr/sbin/mntroot rw"

scp start.sh root@192.168.5.11:/start.sh
scp target/armv7-unknown-linux-musleabi/release/life-dashboard root@192.168.5.11:/main
scp life-dashboard-init root@192.168.5.11:/etc/init.d/life-dashboard-init
scp pokemon.zip root@192.168.5.11:/pokemon.zip

ssh root@192.168.5.11 << EOF
	cd /
	mkdir pokemon
	unzip pokemon.zip -d pokemon
	chmod 777 start.sh
	chmod 777 main
	chmod 777 /etc/init.d/life-dashboard-init

	echo "Install complete, start script with: `/etc/init.d/life-dashboard-init start` and let it run till the ssh session hangs"
EOF
