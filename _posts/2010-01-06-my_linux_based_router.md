---
layout: post
title: My Linux based router
date: 2010-01-06 02:27:15.000000000 +01:00
type: post
parent_id: '0'
published: true
password: ''
status: publish
categories:
- Linux
tags:
- d945gclf
- debian
- Linux
- router
meta:
  _edit_last: '1'
  _wp_old_slug: '2'
author: Gergo Huszty
permalink: "/my_linux_based_router/"
excerpt: Some idea in the "how to build a home server with cheap and low power consuming
  hardware" topic :)
---
Hello,

This is my first real post, in which I will share my experiences about building my Intel Atom and Debian Linux based router (ok, not just a router), and the result itself :).  
I needed something:

- to store my stuff (movies, musics, photos, source codes etc.) locally, not on a remote server or something
- to get my stuff available on the internet if needed
- to download stuff to it, via bittorrent 24/7
- to handle routing and firewall roles
- to use less power than a desktop machine, with two ethernet cards

I was thinking about some ASUS or Linksys router with USB and an USB HDD, but it is not so efficient and powerful in p2p. Then I decided to make a serverlike Linux system, with a low power consumption x86 configuration, than I can do anything what Linux supports. Of course I wanted to solve it in a small place. Ok, let's get mini-ITX. However, VIA EPIAs are always expensive (maybe because they are for the industry), Intel had a new (in the middle of 2008) mini-ITX motherboard, called D201GLY2. This mobo had an embedded [Intel Celeron 220](http://ark.intel.com/Product.aspx?id=33102&processor=220&spec-codes=SLAF2) and SATA ports, which met my needs.

<!--more-->So the first hardware configuration was a bit different than the current, it was from:

- Intel D201GLY2A without CPU fan (it was very loud) and a replaced heatsink
- 1 GB DDR2 RAM
- Second ethernet card in the PCI slot
- [PicuPSU 90](http://www.mini-box.com/picoPSU-90) DC-DC converter
- Seagate Barracuda ST3320620AS SATAII/320GB HDD
- 5 port TP-Link 10/100 Mb ethernet switch built-in, without it's enclosure
- 12V 4A adapter
- Zalman ZM-F2 92mm fan
- a reused enclosure

It was quite cheap and it had only one quiet fan.

I added double resistor to the fan (100Ω) to make it quiet, because I couldn't adjust it enough down in the BIOS. Then I tried to reuse the holes on the enclosure (and drill the rest of course), because previously I wanted to make it to an amplifier (I've never finished it :D ). So I attached the mobo to the bottom, the fan and the HDD to the top (I haven't got enough place on the ground for the HDD, next to the mobo), and the 5 port switch to the back. I plugged the second network card to the switch, then it has 4 LAN ports still.

Then, a year later I started to outgrow from the 320 GBytes and as I had all my important data on that single HDD, if it crashes, I lost them all (I don't really like writing everything to DVDs). I also made a mistake with the Seagate hard drive, actually two. First was to buy it :). Seagates are very trustable and good hard drives, but this type is very loud when it writes to or reads from the disk (it hasn't got acoustic management). It was bad to sleep in one room with the router. The second mistake was to attach it to the top of the metal enclosure, because it added an extra noise to it, even when the hard drive was idle.

So I saved up some money and decided to upgrade my system to use two HDDs in [RAID1](http://en.wikipedia.org/wiki/RAID) to mirror my important data. I got two [1TB WD caviar green](http://www.wdc.com/en/products/products.asp?DriveID=336) hard drive for a good price (and with half power power consumption comparing to the regular drives, with acoustic management, etc.), and a new mobo: [Intel D945GCLF](http://www.intel.com/Products/Desktop/Motherboards/D945GCLF/D945GCLF-overview.htm). The mobo wasn't necessary to change, but I got it very cheap :) and it has [Intel Atom 230](http://ark.intel.com/Product.aspx?id=35635&processor=230&spec-codes=SLB6Z) with[hyperthreading](http://en.wikipedia.org/wiki/Hyper-threading).

I made some pictures about the reconstructing :)

[![alter]({{ site.baseurl }}/assets/images/2010/01/1-300x198.jpg "hdd and it's stand")]({{ site.baseurl }}/assets/images/2010/01/1.jpg)My colleague made me a little stand to get the drives attached to the bottom. The special in it is the 10mm distance what it holds on the left bottom. This is needed because when the enclosure is pieced together, it has a rail under the stand (this was in the way previously to get the hard drive to the bottom).

[![]({{ site.baseurl }}/assets/images/2010/01/2-300x196.jpg "hdd and it's stand 2")]({{ site.baseurl }}/assets/images/2010/01/2.jpg)Now I got the chance to say the distances of the screws on the 3,5 " hard drives to the whole world :D, because I coudn't find it in any standard or blog or nowhere on the internet. So every HDD has 3 screwhole per side, the lower distance is (which is between the closest hole to the connector to the center) is 42,5 mm, the other (the bigger, from the center to the orher side) is 60 mm.

[![almost ready]({{ site.baseurl }}/assets/images/2010/01/3-300x229.jpg "almost ready")]({{ site.baseurl }}/assets/images/2010/01/3.jpg)Every pieces in the house :) It has front USB as well.

[![]({{ site.baseurl }}/assets/images/2010/01/4-300x198.jpg "closed")]({{ site.baseurl }}/assets/images/2010/01/4.jpg)Ready to install.

[![]({{ site.baseurl }}/assets/images/2010/01/5-300x224.jpg "the back")]({{ site.baseurl }}/assets/images/2010/01/5.jpg)This is the connector side. Yes, the IO shield is still missing, I will make it there somehow. Here You can see the green UTP cable, connected to the second netword card and to the LAN switch.

So this is my current hardware configuration. My only need is to make the LAN side 1 Gigabit wide :)

Now we arrived to the software part, the really usable part of the post :) From now on, here will be some little tutorial.

The first thing was to install Debian on it. At this point I found out that, I have two already used SATA ports and I only have a SATA DVD-ROM to install, and no USB rack for it. Ok, I can install the system with one HDD, then plug the second and make the RAID mirror, but I wanted to do it at the install phase. So let's somehow install it from flash drive (I've never tried this). I tried to find some tutorial about this with no luck, it's probably my fault :), but finally I found that three lines what I needed. You only need another booted Linux machine, and follow the instructions on [debian-administration.org](http://www.debian-administration.org/article/Boot_Debian_from_an_USB_device). Be careful, for the latest stable version, you need to get the cd-image file from [here](http://www.debian.org/CD/netinst/) to copy to the flash drive. Another tip for the successful boot from USB: on the G945GCLF, there is about five stupid boot option about USB, turn everything on, the last option (USB Mass Storage Emulation Type) should be: "All Fixed Disc". Then it can boot from the flash drive. The process is the same from here, it will find the iso file on the drive, mount it as a cdrom, and start the installer...

## Configuring software RAID

during the install is pretty simple:

- Choose manual partitioning
- When it shows the available disks, you should first create the same size partitions for the mirror
- When creating, select the "Use as:" option to physical volume for RAID and don't forget to toogle the bootable flag to "on" to the boot partition (usually / or /boot if you separate it) on both side of the RAID1
- Now you can choose in the disk menu: "Configure software RAID" follow the instructions...
- Once you set up the RAIDed partition(s), you can see it as a new disk, there you have to select the real "Use as:" and "Mount point:" options as usual

The other difference during the installion is to select the primary network interface, which will be your WAN interface. After install is complete, you can turn off the USB things in the BIOS and set the boot order to the hard drives.

I made a 20 GB for the system and a 300 GB /home partition with RAID1 (for the irreplaceable like photos, source codes), 1+1GB swap and two other partitions without RAID for the not so important data (movies, temp data, etc) to maximise the space.

The next is:

## Setting up sudo (optional)

Now login to the brand new system as root. My practice is to install sudo, and let the regular user (you created one during the installion) in the sudoers file to do everything like in Ubuntu. To do this, add this line to the /etc/sudoers file:

regular\_username&nbsp; ALL=(ALL) ALL

Then you can set the root password to 20 char long and forget it :) because the regular user can become root with

sudo bash

and his own password.

## Setting up the routing itself

Set up the network interfaces in /etc/network/interfaces. It should look like:

> auto lo  
> iface lo inet loopback
> 
> #The WAN network interface  
> allow-hotplug eth0  
> iface eth0 inet dhcp
> 
> #The LAN network interface  
> allow-hotplug eth1  
> iface eth1 inet static  
> address 192.168.1.254  
> broadcast 192.168.1.255  
> netmask 255.255.255.0

I have cable modem connected to the WAN interface, due to this, eth0 just gets a simple IP from the ISP's DHCP. If you have DSL connection, you should configure it first.

The routing rules are basically just some iptables rules, which are executed once, during boot.

If I remember correctly, I found my sample [here](http://www.gentoo.org/doc/en/home-router-howto.xml).

I edited it to meet my needs, and now it look's like this:

```shell
#!/bin/sh

#Load kernel modules

modprobe ip_conntrack_ftp

#First we flush our current rules

iptables -F

iptables -t nat -F
#Setup default policies to handle unmatched traffic

iptables -P INPUT ACCEPT

iptables -P OUTPUT ACCEPT

iptables -P FORWARD DROP
#Interfaces

export LAN=eth1

export WAN=eth0
#detect WAN IP

export WANIP=`ifconfig ${WAN} | grep inet | cut -d: -f 2 | cut -d' ' -f 1`
#Shitlist - If you know some bad-bad guy's IP or IP range

#iptables -A INPUT -i ${WAN} -s xxx.xxx.xxx.xxx/y -j DROP
#Then we lock our services so they only work from the LAN

iptables -I INPUT 1 -i ${LAN} -j ACCEPT

iptables -I INPUT 1 -i lo -j ACCEPT

iptables -A INPUT -p UDP --dport bootps -i ! ${LAN} -j REJECT

iptables -A INPUT -p UDP --dport domain -i ! ${LAN} -j REJECT

#Allow access to our server from the WAN

#iptables -A INPUT -p TCP --dport ftp -i ${WAN} -j ACCEPT

#iptables -A INPUT -p TCP --dport 17654 -i ${WAN} -j ACCEPT

#etc.

#Drop TCP / UDP packets to privileged ports

iptables -A INPUT -p TCP -i ! ${LAN} -d 0/0 --dport 0:1023 -j DROP

iptables -A INPUT -p UDP -i ! ${LAN} -d 0/0 --dport 0:1023 -j DROP
#Port forwarding to the LAN machines

iptables -t nat -A PREROUTING -i ${WAN} -p tcp --dport 12345 -j DNAT --to 192.168.1.99:12345
#Finally we add the rules for NAT

iptables -I FORWARD -i ${LAN} -d 192.168.1.0/255.255.255.0 -j DROP

iptables -A FORWARD -i ${LAN} -s 192.168.1.0/255.255.255.0 -j ACCEPT

iptables -A FORWARD -i ${WAN} -d 192.168.1.0/255.255.255.0 -j ACCEPT

iptables -t nat -A POSTROUTING -o ${WAN} -j MASQUERADE
#DNAT 22 to 433 from some special places. It can happened that some proxy or firewall blocks the 22 port. From there you can not access your router. This is bad :) but if that firewall allows to connect trough 443 (https), and usually allows, you can specify IP ranges, from where the router accepts the connection on the TCP port 443, and forwards it to the 22 locally.

#iptables -t nat -A PREROUTING -i ${WAN} -p tcp --source xxx.xxx.xxx.xxx/y--dport 443 -j DNAT --to ${WANIP}:22

#Tell the kernel that ip forwarding is OK

echo 1 > /proc/sys/net/ipv4/ip_forward

for f in /proc/sys/net/ipv4/conf/*/rp_filter ; do echo 1 > $f ; done
```

You can place this file for example in /etc/init.d/router, make it executable and make a symlink for it with:

> ln -s /etc/init.d/router /etc/rc2.d/S95router

Then it will starts everytime you boot. Now you need to install some DHCP server to share IPs on the LAN. Install dnsmasq package! This is a lightweight DHCP server and DNS proxy, which is cool because you can make it to access your router by name.

Let's configure the new DHCP in /etc/dnsmasq.conf

> dhcp-range=192.168.1.1,192.168.1.199,72h  
> interface=eth1

Yes, this is enough :)

Set up your router's dns name on the LAN by adding /etc/hosts:

> 192.168.1.254&nbsp;&nbsp; router

Now install openssh-server package and reboot the system. If everything went well, you can connect a client to a LAN port and access the internet behind your router, access your router via ssh (instead of IP, just type "router") and disconnect the monitor and the keyboard from the router forever :)

And now you can add any service what you want. I also added:

- permanent free domain name for my dynamic WAN IP
- [munin](http://munin-monitoring.org/)to monitor the parameters of the system
- [rtorrent](http://libtorrent.rakshasa.no/) for download (with wtorrent GUI and apache2)
- samba to access every downloaded data on the clients, data can simply played on the clients without copying them, through the network
- svn version control system to store my source codes, accessible over http from outside with authentication

Well, thats it for first time, some more tutorials of this list will come soon. Do not hesitate to write comments if you have some question, or just found some mistake in the tutorial.

