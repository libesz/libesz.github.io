---
layout: post
title: Bicolor led matrix with atmega8
date: 2010-01-08 00:06:22.000000000 +01:00
type: post
parent_id: '0'
published: true
password: ''
status: publish
categories:
- AVR
tags:
- atmega8
- AVR
- led matrix
- php
- shift register
meta:
  _edit_last: '1'
author:
  login: admin
  email: huszty.gergo@digitaltrip.hu
  display_name: libesz
  first_name: Gergo
  last_name: Huszty
permalink: "/bicolor-led-matrix-with-atmega8/"
excerpt: A basic concept about how to control a bicolor 8x8 led matrix with AVR
---
Hello,

Today I'm going to show you, how could I draw and play some animation on a bicolor 8x8 led matrix. Approximately a year ago, I found a used, but cheap matrix (it was about 2€). I think this is the most common pinout for the bicolor led matrixes:

[![]({{ site.baseurl }}/assets/images/2010/01/led_matrix_guide2-300x134.jpg "led\_matrix\_inside")](https://libesz.digitaltrip.hu/wp-content/uploads/led_matrix_guide2.jpg)

So the schematic is is:

[![]({{ site.baseurl }}/assets/images/2010/01/ledmatrix-300x178.png "ledmatrix")](https://libesz.digitaltrip.hu/wp-content/uploads/ledmatrix.png)The control is rather simple. One line is drawed at once, but this line is changing very quickly. Thats why the human eye can see one still picture, not just blinking lines. Basically the matrix is driven by shift registers, vertically and horizontally as well. So it consumes only 4 IO pin (two for the data (hor. and vert.) and two for the clocks). The uC could be any attiny variant of the AVR family, but in this case, the AVR have to store all of the frames.

<!--more--> So one shift register controls the matrix vertically, followed by an UDN 2891 driver. This is needed because every line of the matrix has 16 leds (8 bicolor) and the current would be too high for one pin of the 74164 (I tried it, the matrix had almost no light). The first version of the project was build with a 3-8 demultiplexer to control the lines, but there was 6 IO pins needed, so this is more simple. OK, so we had one shift register to select the actual line, by sending one '1' in, and then it just steps over the shift register with the clock of the IC. Once we have a line selected, we have to load is with data. For this, there are another two shift registers, at the bottom for the first eigth LED, then it's last output supplies as the second register's input (for the second eigth LED). By this, we got a 16 bit long shift register (their clock is common). Well, now we can load the data into the horizontal shift registers in 16 steps. We can load the last bit first (it will be pushed to the back by the data behind it, so the first is the 16th LED's state), and so on. If it is loaded, the AVR waits for some time (a few msec) then goes to the next line. The frames practically stored in _unsigned int_ arrays, because they are 16 bits/item long.

A few months ago I assembled one to my colleague as a present:

[![]({{ site.baseurl }}/assets/images/2010/01/board-300x207.jpg "board")](https://libesz.digitaltrip.hu/wp-content/uploads/board.jpg)This board contained a MAX232 and a DB9 connector as well, to write custom programs on it, or just debugging, etc. So the hardware was already done and I had working code on it, but it was difficult to "draw" the frames by counting the bits in an integer array like this: {0x550,0x7ffd,0xfebf,0xfaaf,0xeaab,0xfaaf,0xfebf,0x7ffd} . I know this is not the best solution, I'm not so proud of it and anyway it is has dirty code, but I made a PHP script to, let's say design :) frames in browser. [Here](https://libesz.digitaltrip.hu/wp-content/uploads/ledmatrix.php) it is. It has several buttons, by pressing a button, it generates a C array definition, what you can insert to the code (this is also attached to the project). You have to attach these generated lines, attaching them with a comma to keep C syntax.

I made a video of the prototype:

<object classid="clsid:d27cdb6e-ae6d-11cf-96b8-444553540000" width="425" height="350" codebase="http://download.macromedia.com/pub/shockwave/cabs/flash/swflash.cab#version=6,0,40,0"><param name="src" value="http://www.youtube.com/v/yOodRXKPx60">
<embed type="application/x-shockwave-flash" width="425" height="350" src="http://www.youtube.com/v/yOodRXKPx60"></embed></object>So thats all for today, I know it is a very basic implementation. My plan is to write some low level function to manipulate the content of the frame, then it will able to draw custom things what comes at runtime (process some audio input with fft like in Winamp or little 'video games', etc.)

Here is the code: [led\_matrix.zip](https://libesz.digitaltrip.hu/downloads/led_matrix.zip)

