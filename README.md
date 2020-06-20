## acupoint-portege-x30
Userspace driver for toshiba portege x30-f accupoint device 06cd cddc

## Readme
On new toshiba portege x30-f and some other new models 
touchpad + trackpoint (accoupoint) combo [not working properly](https://bugzilla.kernel.org/show_bug.cgi?id=205817):
touchpad is working, but the trackpoint and physical button aren't.

This programm reads raw input of this device, move your mouse and clicks your buttons.
So basicly it makes trackpoint work.

Programm supports middle button click (when two buttons pressed at once) and scrolling on middle button hold.


### Device file
You need to find where is device file, on my machine it's /dev/hidraw0. 
If on your machine it's located somewhere else, just pass it to promgramm as an argument:
```
./main -p /dev/hidraw#
```

### Modules in kernel

As I can understand in linux kernel exists modules for such devices. 
Most of the previous models with touchpad+trackpoint combo in one device is working. 
I suppose that you need just add couple of lines in suitable driver with device code (06cd cddc) and it will work flawlessly.
But I have no idea how it should be done properly, so if anyone got a clue or can give me advice I'll be very thankful.
