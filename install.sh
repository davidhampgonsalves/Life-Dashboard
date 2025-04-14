#!/bin/sh
zip client.zip fbink life-dashboard-init start.sh life-dashaboard-init
ssh root@192.168.2.190 "/usr/sbin/mntroot rw"

scp client.zip root@192.168.2.190:client.zip

ssh root@192.168.2.190 << EOF
  unzip client.zip -d /
	cd /
  mv life-dashboard-init /etc/init.d/life-dashboard-init
	chmod 777 start.sh
	chmod 777 /etc/init.d/life-dashboard-init
  chmod u+x fbink

	echo "Install complete, start script with: `/etc/init.d/life-dashboard-init start` and let it run till the ssh session hangs"
EOF
