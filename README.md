# Life Dashboard
Heads up Display for every day life.

# Setup

## USB
transfer `pokemon` folder to kindle.

## Jailbreak and Setup SSH
See (https://wiki.mobileread.com/wiki/Kindle4NTHacking) and if bricked then use Kubrick in VM to restore.

connect to usb ssh
network settings, find RNDIS, change from DHCP to manual and ip: 192.168.15.201
```
# set ip of computers usb port
ifconfig # search for device with 192.168.15.201
sudo ifconfig en5 192.168.15.201
ssh root@192.168.15.244
/usr/sbin/mntroot rw
mv /mnt/base-us/pokemon/ /
```

## Cross Compiling to Kindle (ARM-7 Soft Float)
We need a statically compiled binary to run in the Kindle. There are many ways to do this but on OSX I use docker(via https://github.com/messense/rust-musl-cross) to avoid polluting my system with all the required bits and having to compile each requirement seperately.

```
docker pull messense/rust-musl-cross:arm-musleabi
alias rust-musl-builder='docker run --rm -it -v "$(pwd)":/home/rust/src messense/rust-musl-cross:arm-musleabi'
rust-musl-builder cargo build --release
```

# Transfer files to Kindle
```
install.sh
```

## Copy books to kindle vis SCP
```
scp book.mobi root@192.168.15.244:/mnt/base-us/documents/
dbus-send --system /default com.lab126.powerd.resuming int32:1
```
