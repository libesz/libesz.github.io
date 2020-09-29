---
layout: post
title: Six channel audio selector with remote control
date: 2011-04-05 22:43:22.000000000 +02:00
type: post
parent_id: '0'
published: true
password: ''
status: publish
categories:
- AVR
tags:
- atmega8
- audio
- audio selector
- audio switch
- AVR
- infrared
- ir
- music
- relay
- remote
- toner transfer method
- tsop31233
- udn2981
meta:
  _edit_last: '1'
author:
  login: admin
  email: huszty.gergo@digitaltrip.hu
  display_name: libesz
  first_name: Gergo
  last_name: Huszty
permalink: "/audio-selector-with-remote-control/"
---
Hi again!

Today I'm going to show You an audio selector, which I made for my friend and which has an extra feature: it can be controlled with a regular IR remote controller (which you probably have at home :)).

I could say it is a general remote control reciever library with a demo application: an audio selector :-). I could say this because I completely rewrote my infra red handling code, which I made for my [IR remote switch,](https://libesz.digitaltrip.hu/ir-remote-switch/) on a way that it can be used for any other applications easily. The decoding is the same as before, I have 4 IR remote controls around me, and it works with 3 of them. The biggest improvement is that you don't have to measure the startbit with an other application before you can use one of your controllers, it is automatic now. The code is under 2kBytes, so you can port it even to an ATtiny25 (I used ATmega8, cause (now at Hungary) it is for the same price as an ATtiny or an ATmega48). An other comfort feature is to remember the last selected source in case of the device is switched off.

<!--more-->

## Hardware

The requirement was the good sound quality of course, therefore I used relays to route the audio signals. Here is the schematic (VCC is 5V):

[![click to enlarge]({{ site.baseurl }}/assets/images/2011/04/audio_switch_6ch-803x1024.png "audio\_switch\_6ch")](https://libesz.digitaltrip.hu/wp-content/uploads/audio_switch_6ch.png)  
Originally it was planned for 8 channels, but it would be too expensive for to more unused relays :-). I designed the PCB and the code, the owner assembled it :-) The result:

&nbsp;

[![]({{ site.baseurl }}/assets/images/2011/04/IMGP8817-300x157.jpg "IMGP8817")](https://libesz.digitaltrip.hu/wp-content/uploads/IMGP8817.jpg)[![]({{ site.baseurl }}/assets/images/2011/04/IMGP8814-300x220.jpg "IMGP8814")](https://libesz.digitaltrip.hu/wp-content/uploads/IMGP8814.jpg)

[![]({{ site.baseurl }}/assets/images/2011/04/IMGP8821-e1301423117700-751x1024.jpg "IMGP8821")](https://libesz.digitaltrip.hu/wp-content/uploads/IMGP8821.jpg)[![]({{ site.baseurl }}/assets/images/2011/04/IMGP8823-650x1024.jpg "IMGP8823")](https://libesz.digitaltrip.hu/wp-content/uploads/IMGP8823.jpg)

You can see two strange things on the pictures.

The first is the supply connection, which was a mistake in the design to put it on the PCB with inverted polarity, thats why it is now connected on a cable :) If you download it, there is the corrected version.

The second is the wire, soldered directly to the AVR's 3th pin :-). It is the USART TX for the debugging, it will be removed.

The aim was to make it on two separate PCBs, because of the buttons has to put on the front panel of the enclosure, which will be the owner's task to make, here will be more pictures if it is ready.

## Software

The code is not that difficult, so you should just compile and load it, no finetuning is needed. What I will explain here, is the usage of the IR library, which should be now portable (I hope ;-) ).

There are several functions which acts as the IR library interface. The three basic ones are:

- ir\_handle\_timer(); It should be called in a timer interrupt. This will measure the time of the signal pieces. The ideal calling frequency is about 30kHz, which is the frequency of an 8-bit timer without prescaling the 8MHz clock.
- ir\_handle\_input( tsop\_in ); This should be called at least if the input level was changed (if it is connected to an external INT pin), but it can be called even in the main loop as many times as it can. The parameter is the level of the IR reciever's output (0 or 1).
- ir\_is\_recieve\_ready(); This procedure returns the recieved signal's number, which was learned previously. This should be called every time, when you are interested in what there is going on in the air ;-). Be aware that if this procedure returns back the number, it will delete it from the buffer, so you should use it after the first time you get it.

And two more for the learning:

- ir\_start\_learning( number\_of\_the\_signal ); This call will switch the IR stack to learning mode. You can define special conditions for the learning, for example I start it when the first button (on the board) is pressed when the device is switched on. While learning the ir\_handle\_input() still has to be called continously. You can define how much equal samples should recieved, before saving it to the EEPROM. The parameter is the number of the button, you will get this back from ir\_is\_recieve\_ready() whether it identifies the signal.
- ir\_get\_learning\_state(); This will return the state of the learning, you should call this also in a loop. When the learning is done, you can start the next one or go forward.

[caption id="attachment\_566" align="aligncenter" width="479" caption="Screenshot when the code is compiled to use the debugging"][![]({{ site.baseurl }}/assets/images/2011/04/audio_selector_debug.png "audio\_selector\_debug")](https://libesz.digitaltrip.hu/wp-content/uploads/audio_selector_debug.png)[/caption]

&nbsp;

&nbsp;

## Download

- [Eagle files](https://libesz.digitaltrip.hu/downloads/audio_switch_eagle_files)
- [Sourcecode](https://github.com/libesz/AVR_IR_audio_selector)(my new practice is to push everything to GitHub :) there you can also download the code, compressed in ZIP)

That's all for today. Suggestions and comments are welcome.

Cheers

&nbsp;

