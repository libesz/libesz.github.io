---
layout: post
title: 3D snake game for the LEDcube
date: 2011-06-30 23:01:34.000000000 +02:00
type: post
parent_id: '0'
published: true
password: ''
status: publish
categories:
- AVR
- Linux
tags:
- '1000'
- 10x10x10
- 3D
- '4017'
- atmega32
- AVR
- eagle
- enc28j60
- IRLIZ24
- Kingbright
- L829-1X1T-91
- led
- led-cube
- ledcube
- MagJack
- multiplexing
- snake
- TLC5947
meta:
  _edit_last: '1'
author: Gergo Huszty
permalink: "/3d-snake-ledcube/"
---
Hi!

I've just finished the next application for my [FadeCube](https://libesz.digitaltrip.hu/ledcube/), a 3D snake game! It is a very small pure C application, developed under Linux with CodeBlocks IDE. It uses a multithreaded solution. Actually it is my first try to use POSIX threads.

The first thread is listening for the keyboard. Another one is computing the snake pieces and the food. Tthe third is responsible to render the cube data and sends it to the cube. It's always blocked on a mutex, until the snake thread sends a signal.

The snake length is always incremented when it eats :-) . It is controlled by WASD and /' buttons.

The result:

<iframe width="480" height="360" src="http://www.youtube.com/embed/qmcie9ggPOk" frameborder="0" allowfullscreen></iframe>

It seems much better in live.

The code: [https://github.com/libesz/FadeCube\_3D\_snake](https://github.com/libesz/FadeCube_3D_snake)

Cheers!

&nbsp;

