---
layout: post
title: 'Project: Grass Head [time lapse video]'
date: 2014-12-15 23:23:18.000000000 +01:00
type: post
parent_id: '0'
published: true
password: ''
status: publish
categories:
- AVR
tags:
- ATTINY45
- AVR
- C#
- grasshead
- quickdesign
- timelapse
meta:
  _edit_last: '1'
  _wpas_done_all: '1'
  _oembed_c08332ac40e076be73a83a0e6da00d5b: <iframe width="500" height="375" src="https://www.youtube.com/embed/6GFpo6b_j7E?feature=oembed"
    frameborder="0" allowfullscreen></iframe>
  _oembed_time_c08332ac40e076be73a83a0e6da00d5b: '1477220086'
author: Gergo Huszty
permalink: "/project-grass-head/"
---
## Story

Last week we started to foster a [grasshead](https://www.google.com/search?q=grass+head&source=lnms&tbm=isch&sa=X&ei=40ePVK39A4KrUb_vgeAM&ved=0CAgQ_AUoAQ&biw=1440&bih=763). After the first day I remembered a [cool project](http://www.doc-diy.net/photo/smatrig21/) which can remotely trigger any DSLR camera which have IR or wired remote shutter option. I have a [Pentax K200D](http://en.wikipedia.org/wiki/Pentax_K200D)which is a regular camera with such a [2.5mm jack based shutter trigger](http://www.doc-diy.net/photo/eos_wired_remote/pinout.png). This DIY trigger is great for a lot project. Some time ago I decided to build it to capture some lightning without observing the thunder during a whole night. Of course I've never built it :-) .

So I decided to create a time lapse video of the grasshead while the grass is growing. Needless to say that the linked project is an overkill to achieve this and I had only one day before the grass started to grow as&nbsp;it was already sprinkled :-) . So the quick and dirty solution was selected, which is a small AVR and a [solid state relay](http://www.vishay.com/docs/83806/lh1502ba.pdf), nothing more. The selected time interval was one picture in every half hour. Rendering into a 15fps movie, one day will take 3.3 sec.

<!--more-->

## "Hardware"

I had to achieve the maximum battery life for the camera so I set it to sleep after one minute of idle. To get fairly constant light, I used flash for the pictures and F8 aperture for the good depth of field. There was pretty cloudy days of this week, so the sun wasn't bother the pictures too much.

## [![ShutterTimer]({{ site.baseurl }}/assets/images/2014/12/ShutterTimer-300x141.png)](https://libesz.digitaltrip.hu/wp-content/uploads/ShutterTimer.png)

## [![SONY DSC]({{ site.baseurl }}/assets/images/2014/12/DSC07485-300x199.jpg)](https://libesz.digitaltrip.hu/wp-content/uploads/DSC07485.jpg) [![SONY DSC]({{ site.baseurl }}/assets/images/2014/12/DSC07488-300x199.jpg)](https://libesz.digitaltrip.hu/wp-content/uploads/DSC07488.jpg)

I installed one high brightness LED to help the autofocus when the camera is waking up from sleeping during the night.

## "Software"

For such a complex hardware, the software is also pretty huge :-) . The 1MHz downclocked ATTINY45 produces 3906 overflow with the 8 bit timers, which is almost 4 tick per a second if I use a 1024 prescaler. More than enough accurate for the quick design. I used only macroes to express the algorithm, so the program contains no function call at all and only need 2 bytes of memory :-) .

The elapsed time is counted in the timer interrupt. First the AF is triggered 4 seconds before the deadline. This was far enough for my camera to wake up from sleep mode and find the focus (if focus is not found when the shutter is triggered, it won't make the picture). One second before the selected time interval it triggers the shutter as well. Than it releases both trigger contacts.

```
/\* &nbsp;\* ShutterTimer.c &nbsp;\* &nbsp;\* Created: 12/7/2014 1:02:25 PM &nbsp;\*&nbsp; Author: libesz (huszty.gergo@digitaltrip.hu) &nbsp;\* &nbsp;\*&nbsp; Chip: ATTINY45@1MHz &nbsp;\*/ &nbsp; &nbsp; #include \<avr/interrupt.h\> #include \<avr/io.h\> &nbsp; //AutoFocus pin is driven from PORTB3 (the inverted half on the relay) #define AF\_OFF() PORTB|=1\<\<3; #define AF\_ON() PORTB&=~(1\<\<3); &nbsp; //Shutter pin is driven from PORTB4 (the non-inverted half on the relay) #define SHUTTER\_OFF() PORTB&=~(1\<\<4); #define SHUTTER\_ON() PORTB|=1\<\<4; &nbsp; //set the two needed pin as output and &nbsp; //turn off both output of the solid state relay (LH1502B) #define INIT\_GPIO() DDRB=(1\<\<3)|(1\<\<4);AF\_OFF();SHUTTER\_OFF(); &nbsp; //set TIMER0 prescaler to 1024; this will cause 3.8 tick/sec (~4HZ) #define INIT\_TIMER0() TCCR0B|=1\<\<CS02|1\<\<CS00;TIMSK|=1\<\<TOIE0; &nbsp; #define SECOND\_PRESCALER 4 #define MINUTE\_PRESCALER SECOND\_PRESCALER\*60 #define HOUR\_PRESCALER MINUTE\_PRESCALER\*60 volatile unsigned int overflow\_counter = 0; &nbsp; //set this to adjust the shooting interval #define SHOOT\_PRESCALER MINUTE\_PRESCALER\*30 &nbsp; ISR(TIMER0\_OVF\_vect) { &nbsp; ++overflow\_counter; &nbsp; switch(overflow\_counter) { &nbsp;&nbsp;&nbsp; case SHOOT\_PRESCALER-(5\*SECOND\_PRESCALER): &nbsp;&nbsp;&nbsp;&nbsp;&nbsp; AF\_ON(); &nbsp;&nbsp;&nbsp;&nbsp;&nbsp; break; &nbsp;&nbsp;&nbsp; case SHOOT\_PRESCALER-(1\*SECOND\_PRESCALER): &nbsp;&nbsp;&nbsp;&nbsp;&nbsp; SHUTTER\_ON(); &nbsp;&nbsp;&nbsp;&nbsp;&nbsp; break; &nbsp;&nbsp;&nbsp; case SHOOT\_PRESCALER: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp; SHUTTER\_OFF(); &nbsp;&nbsp;&nbsp;&nbsp;&nbsp; AF\_OFF(); &nbsp;&nbsp;&nbsp;&nbsp;&nbsp; overflow\_counter = 0; &nbsp;&nbsp;&nbsp;&nbsp;&nbsp; break; &nbsp; }&nbsp; &nbsp; } &nbsp; int main(void) { &nbsp; INIT\_GPIO(); &nbsp; INIT\_TIMER0(); &nbsp; sei(); &nbsp; while(1); }
```

## The result

After a week of taking the 340 pictures I rendered out the result with the help of [this article](http://ubuntuforums.org/showthread.php?t=2022316). Check it in 1080p, it is quite cool :-)

http://youtu.be/6GFpo6b\_j7E

You can build such a setup in an hour.

Happy timelapsing!

&nbsp;

