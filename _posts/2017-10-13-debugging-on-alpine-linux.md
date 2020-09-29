---
layout: post
title: Debugging and contributing on Alpine Linux
date: 2017-10-13 23:54:44.000000000 +02:00
type: post
parent_id: '0'
published: true
password: ''
status: publish
categories:
- Linux
tags:
- alpine
- container
- dlna
- docker
- Linux
- minidlna
- router
meta:
  _edit_last: '1'
  _oembed_7df4e5b5c7d478969a18adb39d2c7082: "{{unknown}}"
  _oembed_4c9d35dc5c43766e899ecec6c0fffba8: "{{unknown}}"
  _oembed_b01a261f3549a8385fa618cf14d5a258: "{{unknown}}"
  _wpas_done_all: '1'
author: Gergo Huszty
permalink: "/debugging-on-alpine-linux/"
---
In one of my previous posts I [explained](https://libesz.digitaltrip.hu/linux-based-router-reloaded/) my renewed router / home server. One task of the box is to serve video/audio content on DLNA. One of the easy selection in this area is minidlna to do the streaming. So I grabbed the first working [minidlna docker container](https://hub.docker.com/r/vimagick/minidlna/), which in practice Alpine Linux based and started to [use](https://github.com/libesz/lhs-ansible/blob/master/roles/minidlna/tasks/main.yml) that.

Our happiness was not instant using the new configuration. Minidlna was never rock-solid, but in this setup it definitely crashed from time to time. As I inspected it crashed every single time when something added to the media library, in practice when the download completed. After checking minidlna issues, I have not found anything useful, so decided to locate the exact problem.

<!--more-->

# Get our hands dirty

The investigation started with installing the [_gdb_ debugger](https://www.gnu.org/software/gdb/) inside the minidlna container. Since the project is written in pure C, that is the default way of doing this. Let's create a container instance for the hacking and install minidlna + gdb.

```
$ docker run -ti --name minidlna\_debug alpine:3.5 sh / # apk update fetch http://dl-cdn.alpinelinux.org/alpine/v3.5/main/x86\_64/APKINDEX.tar.gz fetch http://dl-cdn.alpinelinux.org/alpine/v3.5/community/x86\_64/APKINDEX.tar.gz v3.5.2-90-g737768f35c [http://dl-cdn.alpinelinux.org/alpine/v3.5/main] v3.5.2-81-gf4d50b1370 [http://dl-cdn.alpinelinux.org/alpine/v3.5/community] OK: 7962 distinct packages available / # apk add minidlna gdb [...] Executing minidlna-1.1.5-r3.pre-install Executing busybox-1.25.1-r0.trigger OK: 98 MiB in 62 packages / #
```

With the new image, running minidlna with gdb I reproduced the issue. A lot of basic gdb usage and cheat sheet can be found out there, will not go into details here. As it is expected, it exactly shows nothing else than the type of the exception, which was naturally a Segmentation Fault. Something happened here what is not user error for sure :). Linux distributions are usually not shipping debugging symbols in the default production binary packages, so gdb showed only cryptic memory address in the stack trace. You can actually check this right from the binary with the _file_ tool:

```
# file /usr/sbin/minidlnad /usr/sbin/minidlnad: ELF 64-bit LSB shared object, x86-64, version 1 (SYSV), dynamically linked, interpreter /lib/ld-musl-x86\_64.so.1, stripped, with debug\_info
```

The main point is the notation that it is _stripped_. This means you have no symbol and source code information from compilation time.

That was the momentum when I had to get familiar with [Alpine Linux](https://alpinelinux.org/) build system, because the present minidlna container image used that as a base. Actually it was already familiar to me on some level. I use it @work, it is popular for it's extreme small size, etc. BTW size is the key why it is popular to build docker images on top of it. While mainstream full-blown distributions are consuming hundred(s) of MBs because of the lot usual tools in there by default, but Alpine base image is exactly 3.96MB today :).

```
$ docker images REPOSITORY&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; TAG&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; IMAGE ID&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; CREATED&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; SIZE ubuntu&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; latest&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 7b9b13f7b9c0&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 7 days ago&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 118MB alpine&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; latest&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; a41a7446062d&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 2 weeks ago&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 3.96MB
```

## Recompile an Alpine package

After spinning the Google, I found some intro page [like this](https://wiki.alpinelinux.org/wiki/Abuild_and_Helpers) on how to recompile things. So _abuild_ is the thing here. The whole compiling toolchain can be installed with a single command:

```
apk add alpine-sdk
```

Apk it the package manager for alpine. You might call also "_apk update_" before the very first apk . _Abuild_ is the tool for doing everything around package maintenance except that it won't tell you how to get the package :(. One example of my failures:

```
$ abuild fetch minidlna \>\>\> ERROR: : Could not find ./APKBUILD (PWD=/)
```

Basically that is the response for anything until you get the right _APKBUILD_, whatever it is. Actually it is the main descriptor of the packages, one for each. After further utilizing the most famous search engines, I found out that _Alpine_ guys are collecting their _APKBUILD_s and related files in a git repository. This is the _[aports](https://github.com/alpinelinux/aports/)_. The easiest way to get is to clone the whole repo with git. Let's do it after creating a regular user (+add it to the abuild group) for building stuff:

```
/ # adduser build Changing password for build New password: Retype password: passwd: password for build changed by root / # addgroup build abuild / # su build / $ cd ~ $ git clone https://github.com/alpinelinux/aports.git Cloning into 'aports'... remote: Counting objects: 344808, done. remote: Compressing objects: 100% (54/54), done. remote: Total 344808 (delta 29), reused 51 (delta 19), pack-reused 344732 Receiving objects: 100% (344808/344808), 165.25 MiB | 3.61 MiB/s, done. Resolving deltas: 100% (209266/209266), done. ~ $
```

Depending on the release of the alpine release you might change to some other branch in the repo (like 3.5). Also, add your user to the sudoers as the tools are trying to use sudo it when root privileges are needed.

Create a key to be able to generate signed packages and get into to the minidlna directory:

```
~ $ abuild-keygen -a -i ~ $ cd aports/community/minidlna/ ~/aports/community/minidlna $
```

Now we are finally ready to build the repackage the package. Before that you need to update the checksums in the APKBUILD file:

```
~/aports/community/minidlna $ abuild checksum && abuild -r \>\>\> minidlna: Updating the md5sums in APKBUILD... \>\>\> minidlna: Updating the sha256sums in APKBUILD... \>\>\> minidlna: Updating the sha512sums in APKBUILD... \>\>\> minidlna: Checking sanity of /home/build/aports/community/minidlna/APKBUILD... \>\>\> minidlna: Analyzing dependencies... \>\>\> minidlna: Installing for build: build-base bsd-compat-headers libvorbis-dev libogg-dev libid3tag-dev libexif-dev libjpeg-turbo-dev sqlite-dev ffmpeg-dev flac-dev WARNING: Ignoring /home/build/packages//community/x86\_64/APKINDEX.tar.gz: No such file or directory (1/49) Installing bsd-compat-headers (0.7-r1) (2/49) Installing libogg (1.3.2-r1) [...] (49/49) Purging xvidcore (1.3.4-r0) Executing busybox-1.25.1-r0.trigger OK: 192 MiB in 68 packages \>\>\> minidlna: Updating the community/x86\_64 repository index... \>\>\> minidlna: Signing the index... ~/aports/community/minidlna $
```

Somewhere in the middle you can inspect that the code is actually compiled. The final package can be found in ~/packages... .

Let's install it:

```
/ # apk add /home/build/packages/community/x86\_64/minidlna-1.1.5-r3.apk (1/40) Installing libogg (1.3.2-r1) [...] (40/40) Installing minidlna (1.1.5-r3) Executing minidlna-1.1.5-r3.pre-install Executing busybox-1.25.1-r0.trigger OK: 232 MiB in 108 packages / #
```

## Get debugging symbols in

Now we still have stripped end-result (you can check it with the file command after installing it with "_apk add file_"). We can change the compilation to let the debug symbols stay in the binary. It is about to add an option to the APKBUILD file of the package:

```
options="!strip"
```

If now you recompile and re-install the package, the file command confirms that we are getting closer:

```
/ # file /usr/sbin/minidlnad /usr/sbin/minidlnad: ELF 64-bit LSB shared object, x86-64, version 1 (SYSV), dynamically linked, interpreter /lib/ld-musl-x86\_64.so.1, not stripped
```

So we have the not stripped binary. As I reproduced the problem with this, I got the correct stack trace for real debugging.

## Contribution

Let's jump over the actual correction for now. The point is that the aports library actually does not contain the main source tree of the project, it only refers to the original source repo. What it really contains is:

- The build configuration for the alpine system (the APKBUILD file).
- Sample configuration or other files for the target system to be used after installation.
- Installation, upgrade hooks for the package management events.
- Finally the most important: the patches. If the project needs some source patches in order to run correctly on Alpine linux, those are added to this repo and applied during compilation.

Generally, if the problem happened to be a real bug in the software, you might end up with a patch for the software what fixes the issue. If you decide to contribute it back to the upstream, there can be two targets. If if is a generic bug in the software, you should send it to the author. In my case the correction was a special one, because the bug is alpine specific. It crashed because of the extreme small size of the child threads of a process. So my patch belongs to the aports repo. I sent the patch to their mail based patch handling system and after some mail exchange with the package maintainer it got [approved](http://patchwork.alpinelinux.org/patch/2940/). You can be much smarter by sending a patch directly to github, as I realized since than that they accept contribution [there](https://github.com/alpinelinux/aports/) also :).

My contribution is finally in the [stable 3.6 realease](https://pkgs.alpinelinux.org/package/v3.6/community/x86_64/minidlna) of the alpine distribution and I can use it without issues, hurray :).

### Some further references for the alpine package build system:

[https://wiki.alpinelinux.org/wiki/APKBUILD\_Reference](https://wiki.alpinelinux.org/wiki/APKBUILD_Reference)

[https://wiki.alpinelinux.org/wiki/Creating\_an\_Alpine\_package](https://wiki.alpinelinux.org/wiki/Creating_an_Alpine_package)

[https://wiki.alpinelinux.org/wiki/Abuild\_and\_Helpers](https://wiki.alpinelinux.org/wiki/Abuild_and_Helpers)

[https://wiki.alpinelinux.org/wiki/Include:Setup\_your\_system\_and\_account\_for\_building\_packages](https://wiki.alpinelinux.org/wiki/Include:Setup_your_system_and_account_for_building_packages)

&nbsp;

