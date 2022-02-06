# Life Dashboard
Low power, heads up display for every day life running on a Kindle.

<img align="right" src="https://github.com/davidhampgonsalves/life-dashboard/raw/master/life-dashboard.jpg"/>

# Details
Second hand Kindles are waiting in drawers for someone to repurpose them into something great. Boasting large e-ink screens, wifi connectivity and ARM processors they are an amazing hacking platform.

In my case I created an information panel summarizing my day such as my calendar, surf and weather forecast, garbage schedule, school closures, etc. The extra space is filled by a random pokemon sprite.

The project uses a serverless backend to collate data from external services and on the Kindle itself [Rust](https://www.rust-lang.org/) code (cross compiled via docker) fetches and typesets the data into an image.

I built a stand rather then a more standard frame because the e-reader functionality of the Kindle is still present and can be used without modification. I also thought it was important to avoid obscuring the original device and celebrate its reuse.

More details can be found on my [blog](https://www.davidhampgonsalves.com/life-dashboard/).

# Setup

## USB
transfer `pokemon` folder to Kindle mounted as USB drive.

## Jailbreak and Setup SSH
See (https://wiki.mobileread.com/wiki/Kindle4NTHacking) and if bricked then use Kubrick in VM to restore.

## SSH over wifi
Hold power button till light flashes, then press power button a few times to restart back to normal e-reader mode. SSH server will be running and wifi will auto connect.

## SSH Over USB
network settings, find RNDIS, change from DHCP to manual and ip: 192.168.15.201.
```
# set ip of computers usb port
ifconfig # search for device with 192.168.15.201
sudo ifconfig en5 192.168.15.201

ssh root@192.168.15.244
/usr/sbin/mntroot rw
mv /mnt/base-us/pokemon/ /
```
## Install
Setup Wifi on Kindle and then run `install.sh` with Kindle connected via USB or wifi.

# Cross Compiling to Kindle (ARM-7 Soft Float)
We need a statically compiled binary to run in the Kindle. There are many ways to do this but on OSX I use docker(via https://github.com/messense/rust-musl-cross) to avoid polluting my system with all the required bits and having to compile each requirement separately.
```
docker pull messense/rust-musl-cross:armv7-musleabi && \
alias rust-musl-builder='docker run --rm -it -v "$(pwd)":/home/rust/src messense/rust-musl-cross:armv7-musleabi'
rust-musl-builder cargo build --release
```

# Cross
Cross doesn't support soft float for arm 7 yet.
```
cross build --target armv7-unknown-linux-musleabihf
```

# Copy books to Kindle vis SCP
```
scp book.mobi root@192.168.15.244:/mnt/base-us/documents/
dbus-send --system /default com.lab126.powerd.resuming int32:1
```

# Frame
3D printed using wood filled filliment - https://www.thingiverse.com/thing:2536906

## Notes
The [mobileread forumn](https://www.mobileread.com/forums/) is the place for mobile reader hacking.

I could have avoided the backend of this project and only used the Kindle but I had already created it for another project and saved time to reuse it.

