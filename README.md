<img style="width:400px" align="right" src="https://github.com/davidhampgonsalves/life-dashboard/raw/master/example.jpg"/>

# Life Dashboard
Low power, heads up display for every day life running on a Kindle.

## Details
Second hand Kindles are waiting in drawers for someone to repurpose them into something great. Boasting large e-ink screens, wifi connectivity and ARM processors they are an amazing hacking platform.

(This is the second version of this project, see the post about the original [here](https://www.davidhampgonsalves.com/life-dashboard/))

## V2 Rewrite and Compromises
Ideally this dashboard would generate and display its image on its own. The issue with doing this originally was that the Kindles ability to display images (via eips) requires they be in a strange format. A few yars back when I started this it was hard to get GoLang to generate this format but it was easy to cross compile GoLang to the target ARM-7 softfloat arch. On the other hand Rust could generate the image but it was a pain to setup the cross compiler toolchain (also connecting to Google API's wasn't well supported).

In the 5 years that followed the dashboard was a great tool but as API services would die (Magicseaweed, Forecast.io, DarkSky, etc) it would break. Sometimes that would require changes to the Rust front end and I would have to setup the cross compiling toolchain on each new machine I was using and eventually this got annoying.

I found out about [FBInk](https://github.com/NiLuJe/FBInk) which has Go bindings and decided that using that I could use it to do text layout and print the resulting PNG's to the screen. Unfortunately I found that GoLangs OpenFont lib crashes when run on the kindle and that these old arm archetectures aren't well supported. This seemed like an unstable footing to build on.

This led me to my current compromise. I use FBInk on the kindle to display the images after curling them from a API Gateway/Lambda backend. This gives me a low friction way to update the API logic without needing to touch the kindle or cross compile anything. I also was able to use GoLangs [tdewolff/canvas](https://github.com/tdewolff/canvas/) which provides nice text setting and image generation tooling. I think is the right balance to keep this device productive for another 5+ years.

# Self Hosted
Start with `service lifedashboard start`

# Setup

## Jailbreak and Setup SSH
See (https://wiki.mobileread.com/wiki/Kindle4NTHacking) and if bricked then use Kubrick in VM to restore.

## SSH over wifi
Hold power button till light flashes, then press power button a few times to restart back to normal e-reader mode. SSH server will be running and wifi will auto connect.

## Install
* Install Fbink.
* Setup Wifi on Kindle and then run `install.sh`.

# Copy books to Kindle vis SCP
```
scp book.mobi root@192.168.15.244:/mnt/base-us/documents/
dbus-send --system /default com.lab126.powerd.resuming int32:1
```

# Frame
3D printed using wood filled filliment - https://www.thingiverse.com/thing:2536906

## Notes
The [mobileread forumn](https://www.mobileread.com/forums/) is the place for mobile reader hacking.

## My Device
Mac: F0:A2:25:04:37:2C

## Server
`service lifedashboard start`
But you have to kill the process to restart `ps -aux | grep go`
