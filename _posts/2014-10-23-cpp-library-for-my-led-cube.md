---
layout: post
title: C++ library for my LED cube
date: 2014-10-23 18:28:34.000000000 +02:00
type: post
parent_id: '0'
published: true
password: ''
status: publish
categories:
- C++
- Linux
tags:
- C++11
- Cpp
- led
- led-cube
- ledcube
- Linux
meta:
  _edit_last: '1'
  _wpas_done_all: '1'
author: Gergo Huszty
permalink: "/cpp-library-for-my-led-cube/"
---
In the last two years I was mainly interrested in exploring C++. One result of my learning is the subject of this post. I started to reimplement the [snake]({{ site.baseurl }}/3d-snake-ledcube/) game using C++11 for the [LED cube]({{ site.baseurl }}/ledcube/). Later I divided the code into two sections. The general part became a base library for creating games or animations. The snake specific code is now only an application, which is based on the general library part. Later, the famous [2048 game](http://gabrielecirulli.github.io/2048/) was also implemented on top of the library.

<!--more-->

## Library

The basic classes and interfaces of the library are:

- _Point_: holds one LED's position and brightness level. The position is comparable with operator ==; collection of _Point_s passed through the higher level constructs
- _Display_: pure interface for the cube. Used to send the content to the LED cube driver. The real cube's driver is pluggable here (_CubeDisplay_), but there is also a debugger class to dump the content's coordinate's to the screen (_DisplayDumper_)

&nbsp;

- _Schedulable_: interface for the classes which produces variable content
- _ClockSource_: base class for scheduling. Holds a collection of _Schedulable_ objects. Two inherited solution is implemented:
  - _TimerClockSource_: notifies all Schedulable objects in predefined time intervals
  - _CondVarClockSource_: an event driven solution, may used for example to notify the _Schedulable_ objects on keypress
- _Renderable_: interface for content supplier objects; classes based on it has to provide _Point_ collections at any time representing their content
- _Renderer_: holds a collection of _Renderable_ objects and one _Display_; inherits _Schedulable_, so that the _ClockSource_ will trigger the redrawing of the cube's content as well; collects all _Renderable_ objects' contents and passes it to the _Display_ when scheduled
- _ClockDivider_: a _Schedulable_ helper class which holds a real _Schedulable_ object and a division rate; after passed into the _Renderer_, it will only pass a divided amount of update requests to the given real _Schedulable_ objects; this way the update rate can be fine tuned between the content suppliers

&nbsp;

- _Controllable_: pure interface for objects, which content is based on some expected physical _Direction_ (up,down,right,left,forward,backward)
- _KeyboardInput_: instructs a _Controllable_ object based on keypresses (wasd, / and '); keyboard events are waited in a separate thread

## 3D snake game

The snake game's classes are the following:

- _SnakeFood_: supplies one blinking LED on the cube, which the _Snake_ should catch; it is:
  - _Schedulable_ (has to replace itself on random place if the _Snake_ didn't reached it in a given amount of schedule cycles)
  - _Renderable_
- _Snake_: supplies the snake's body; it is:
  - _Renderable_
  - _Schedulable_ (to keep going in the current direction)
  - _Controllable_ (direction is changed based on keyboard presses)
- _SnakeController_: the main controller of the game; it counts the score, checks if the _Snake_ reached the _SnakeFood_, etc; it is:
  - _Schedulable_

The main function acts as a factory for the game. It creates a _Snake_, a _SnakeFood_, passes them (with the help of _ClockDivider_s) to a _Renderer_ together with a _CubeDisplay_. Creates a _KeyboardInput_ to control the _Snake_, gives the _Snake_ and the _SnakeFood_ to a _SnakeController_, etc etc. Then calls the _TimeClockSource_s loop method and the game starts. The gameover is determined by a helper class, namely a _SnakeExitConditionClass_, which passed into the _TimeClockSource_ and checked after every schedule round.

## Game 2048

An other example application is the LED cube adaption of the famous [2048 game](http://gabrielecirulli.github.io/2048/). The values are represented by lighting LED columns at the center of the cube. Each column's height is the number, which 2 has to powered to get the number value of the place. Hard to imagine, but really easy to play the game this way. Basically if two columns has equal height, the columns can be merged together and the result will be a one unit taller column, representing the twice bigger number in the grid. The classes are:

- _GameTable_: holds the gameboards' state and gives interface to alter that state. The core logic if borrowed from [clinew's C implementation of](https://github.com/clinew/2048/) the game, thanks for it
- _GameRenderable_: connects the _GameTable_ to the _Renderer_ by extracting the table content to column data displayed on the cube; it is _Renderable_
- _GameController_: core game logic; through the _Controllable_ interface it controls the _GameTable_
- _GameExitCondition_: similarly to the solution in the snake game, this class determines if the game is over (either by win of lose) of the user wants to quit

Just like in the snake game, the main function only wires the game parts together and let the _ClockSource_ to react to the user's keypresses. That means that here the other, _CondVarClockSource_ is used because the content need to be only rendered when some keypress occured. I've created an UML diagram to represent the 2048 game classes:

[![2048_uml]({{ site.baseurl }}/assets/images/2014/10/2048_uml-300x194.png)]({{ site.baseurl }}/assets/images/2014/10/2048_uml.png)

Code is on the [github](https://github.com/libesz/FadeCube_cpp). It is developed as Eclipse CDT projects under Linux.

You can also find some unittests in the repo written in [GoogleTest](https://code.google.com/p/googletest/). It was easier to reproduce corner cases with UT instead of eg. solve the 2048 on the cube just to test if the application determines the winning situation :).

Cheers!

&nbsp;

