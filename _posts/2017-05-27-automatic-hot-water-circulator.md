---
layout: post
title: Automatic hot water circulator
date: 2017-05-27 23:37:19.000000000 +02:00
type: post
parent_id: '0'
published: true
password: ''
status: publish
categories:
- AVR
- C++
tags:
- AVR
- C#
- heating
- home automation
- motion sensor
- water
meta:
  _edit_last: '1'
  _wpas_done_all: '1'
  _jetpack_dont_email_post_to_subs: '1'
  _wp_old_slug: automatic-how-water-circulator
author: Gergo Huszty
permalink: "/automatic-hot-water-circulator/"
---
After moving to a new place, I've found new challenges to solve here and there. One such unresolved thing was the hot water circulation control. In this post I present one possible solution which is according to the user's behavior.

## The problem

This is a tiny project to control a water pump. In a typical domestic hot water system, there is a common problem that the water gets cold in the pipes if some faucets are too far. So you have to waste a lot of water (and time :) ) until you get the hot water, if is was not used recently.

The base of the solution is to build also a returning pipe between the last faucet and the boiler or hot water tank, where a pump can circulate the returning water, so that keeping the pipes hot.

<!--more-->

The big question is that based on what event are you going to circulate the water? After a bit of googling you can find these ideas:  
\* Timer switch, programmed for the typical hot water use times (morning, evening).  
\* Thermostat: you somehow measure the returning water's temperature and control the pump when needed.  
\* Detect water usage: there are solutions which can detect when you use any faucet in the house. So you can trigger the circulation by using the e.g. shower for a second, than wait until the hot water arrive before you really start to use it.

## Solution

I was not impressed by any of the solutions above, even not with any combinations. Fortunately I had the pump and the returning pipe built in, so I can control it how I want.  
My solution is based on the detection of motion in the house. I have a basic alarm system in the house with passive IR motion sensors in every rooms. The idea is that before every water usage, the people have to walk to the desired faucet. There might be additional time before you need the hot water, like you use the toilet first, for instance :). This time is most probably enough for the pump to get hot water everywhere.  
I inspected my alarm system and found a programmable output on the board which can follow any selection of the motion sensors status, so it can trigger this AVR project, which will control the pump for a predefined time.  
Also there is a hold timer, which will not switch on the pump ON again if the water is supposed to be still hot in the pipe.

My resulted timers are: 5 mins for every trigger and 25 mins hold-off. It worked out very well.

The code is written in C++ and using my base library [from here](https://github.com/libesz/AvrCppBaseLib).

## HW:

- ATTINY45@1MHz
- Opto couplers for both input from the alarm system and towards the output
- 12V relay to control the pump

The essence of the solution is obviously the trigger which comes from the alarm system. AVR only enlarges the trigger time from the movement interval to the necessary time which is needed for the water to take the loop.

[![]({{ site.baseurl }}/assets/images/2017/05/hwc.png)](https://libesz.digitaltrip.hu/wp-content/uploads/hwc.png)The schematic is crazy simple: 12V-5V conversion, input receiver part, output part.

The code is also very simple as you might think. It has to only maintain two time intervals, which are connected. The fancy part is that I used my AVR C++ library, presented in the [last post](https://libesz.digitaltrip.hu/cpp-on-avr/). There I already created the low-level parts for this application in the base library:

- _Soft timer set_: controls arbitrary counters, controlled by one real timer
- _Timed output_: timer handler implementation to set a GPIO to 1 until a predefined time

The missing part is the hold-off login, which became also a similar class to the _timed output_. This is the _Trigger_ class.

```
#include \<avr/io.h\> #include \<avr/interrupt.h\> #include "SoftTimerSet.h" #include "TimedOutput.h" #include "Trigger.h" //Config BEGIN #define OUTPUT\_ON\_TIME 300 #define TRIGGER\_HOLD\_OFF\_TIME 1500 //Config END #define OUTPUT\_PORT\_MASK 1\<\<4 #define INIT\_OUTPUT()&nbsp; DDRB|=OUTPUT\_PORT\_MASK;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; // PORTB3 -\> output #define INIT\_TRIGGER() PORTB|=1\<\<2;MCUCR|=1\<\<ISC01;GIMSK|=1\<\<INT0; // PORTB2 -\> input (by default), activate pull-up and falling edge interrupt #define INIT\_TIMER0() TCCR0B|=1\<\<CS02|1\<\<CS00;TIMSK|=1\<\<TOIE0; // Set TIMER0 prescaler to 1024; this will cause 3.8 tick/sec (~4HZ) #define SECOND\_PRESCALER 4 volatile unsigned int overflow\_counter = 0; SoftTimerSet\<2\> gSoftTimerSet; Trigger \*gTrigger = 0; ISR(TIMER0\_OVF\_vect) { &nbsp; if(++overflow\_counter == SECOND\_PRESCALER) { &nbsp;&nbsp;&nbsp; gSoftTimerSet.tick(); &nbsp;&nbsp;&nbsp; overflow\_counter = 0; &nbsp; } } ISR(INT0\_vect) { &nbsp; gTrigger-\>activate(); } int main(void) { &nbsp; INIT\_OUTPUT(); &nbsp; INIT\_TRIGGER(); &nbsp; INIT\_TIMER0(); &nbsp; TimedOutput output((volatile void \*)&PORTB, OUTPUT\_PORT\_MASK, OUTPUT\_ON\_TIME); &nbsp; gSoftTimerSet.add(&output); &nbsp; Trigger trigger(&output, TRIGGER\_HOLD\_OFF\_TIME); &nbsp; gSoftTimerSet.add(&trigger); &nbsp; gTrigger = &trigger; &nbsp; sei(); &nbsp; while(1); }
```

Pretty easy to consume I think. I became a C++ fan on AVRs :). As usual, my AVR projects are created with the official AVR studio, unfortunately available only on MS Windows. Anyway it would compile just fine under linux with pure avr-gcc. I was simply lazy to create Makefile for it.

## The result

Some visualization of the end result:

[![]({{ site.baseurl }}/assets/images/2017/05/IMG_20170527_162739.jpg)](https://libesz.digitaltrip.hu/wp-content/uploads/IMG_20170527_162739.jpg)

[![]({{ site.baseurl }}/assets/images/2017/05/IMG_20170527_192342.jpg)](https://libesz.digitaltrip.hu/wp-content/uploads/IMG_20170527_192342.jpg)

[![]({{ site.baseurl }}/assets/images/2017/05/IMG_20170527_192347.jpg)](https://libesz.digitaltrip.hu/wp-content/uploads/IMG_20170527_192347.jpg)

Resources are available in the following repo on [github](https://github.com/libesz/hwc).

You can directly download the hex file from [here](https://libesz.digitaltrip.hu/wp-content/uploads/hwc.zip).

