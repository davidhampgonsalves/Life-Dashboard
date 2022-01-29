#!/bin/sh

ssh root@192.168.2.12 "/usr/sbin/mntroot rw"

scp start.sh root@192.168.2.12:/start.sh
scp target/arm-unknown-linux-musleabi/release/life-dashboard root@192.168.2.12:/main
scp life-dashboard-init root@192.168.2.12:/etc/init.d/life-dashboard-init
scp pokemon.zip root@192.168.15.244:/pokemon.zip

ssh root@192.168.15.244 << EOF
	mkdir pokemon
	unzip pokemon.zip -d pokemon
	cd /
	chmod 777 start.sh
	chmod 777 main
	chmod 777 /etc/init.d/life-dashboard-init

	echo "Install complete, start script with: `/etc/init.d/life-dashboard-init start`."
EOF