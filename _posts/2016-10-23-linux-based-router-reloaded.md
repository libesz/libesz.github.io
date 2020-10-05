---
layout: post
title: linux based router reloaded
date: 2016-10-23 17:16:07.000000000 +02:00
type: post
parent_id: '0'
published: true
password: ''
status: publish
categories:
- Linux
tags:
- ansible
- debian
- docker
- Linux
- router
meta:
  _edit_last: '1'
author: Gergo Huszty
permalink: "/linux-based-router-reloaded/"
---
Hi,

It's been a long time since I set up my [mini-itx based router/download machine](https://libesz.digitaltrip.hu/my_linux_based_router/). It got a few unavoidable HW upgrade (power supply, new HDD) but the performance of that is remained. From SW point of view however, it has a lot more load and responsibility. The HD age with the smart TVs are here, so that the machine is not anymore a basic router/NAS, but something like a HTPC. I use [minidlna](https://sourceforge.net/projects/minidlna/) to stream the content towards our smart TVs. Now the performance of the famous [Intel Atom 230](http://ark.intel.com/products/35635/Intel-Atom-Processor-230-512K-Cache-1_60-GHz-533-MHz-FSB) is over unfortunately. The 100Mb/s internet is also a reality now and if the content is downloaded at max speed to the encrypted HDDs, it consumes the CPU entirely and the stream basically starts to lag. Of course it is possible to limit the CPU usage or the downloading bandwidth for rtorrent, but that would be a workaround only.<!--more-->

## The plan

Challenge accepted: reload the central part of our IT infrastructure at home. Every such project shall have some input requirements. My list:

- HW:
  - Upgrade to 1Gb/s ethernet on both LAN and WAN side
  - Get enough CPU power again for the next at least 5 years
  - Reuse as much component as possible
  - Should still consume only a few Watts of energy
- SW:
  - Reinstall with 64bit OS
  - Use uptodate, even over hyped technologies for the operation

## Hardware

From gigabit point of view we were almost already fine. On the LAN side the 100Mb/s switch died some years ago so that I bought a new gigabit version of that. We only need some gigabit NIC to go. To fulfill the remaining requirements, I started to check the ebay-equivalent hungarian market for some mint conditioned, used and cheap mini-itx motherboard. After a few weeks I got this: [Asrock N3700-ITX](http://www.asrock.com/mb/Intel/N3700-ITX/). It costed half as much as was the last price in the stock, like 60 EUR. The [Pentium N3700](http://ark.intel.com/products/87261/Intel-Pentium-Processor-N3700-2M-Cache-up-to-2_40-GHz) has built-in AES capabilities and 4 real cores, based on Braswell. I had a spare Intel 1000 PT for secondary ethernet card and I pulled out a 2GB DDR3 SO-DIMM from my laptop, so we are good to go from HW PoV.

[caption id="attachment\_762" align="aligncenter" width="300"][![The old retired HW.]({{ site.baseurl }}/assets/images/2016/10/IMG_20161020_183659-300x225.jpg)]({{ site.baseurl }}/assets/images/2016/10/IMG_20161020_183659.jpg) The old retired HW.[/caption]

[caption id="attachment\_761" align="aligncenter" width="300"][![The new mobo and the Intel NIC]({{ site.baseurl }}/assets/images/2016/10/IMG_20161020_183726-300x225.jpg)]({{ site.baseurl }}/assets/images/2016/10/IMG_20161020_183726.jpg) The new mobo and the Intel NIC[/caption]

[caption id="attachment\_760" align="aligncenter" width="300"][![The sad fate project box in the middle of the operation. I had to drill yet another hole to have the new NIC fixed.]({{ site.baseurl }}/assets/images/2016/10/IMG_20161020_183740-300x225.jpg)]({{ site.baseurl }}/assets/images/2016/10/IMG_20161020_183740.jpg) The sad fate project box in the middle of the operation. I had to drill yet another hole to have the new NIC fixed.[/caption]

## Software

The SW part was more entertaining than fixing the enclosure to fit for the new PCI-e network card. :)

As I am not a professional person in real linux operation, my installations usually are patched whenever needed because of new need of something failed. So this box represented my laziness. Basically I had no idea of how was some pieces configured, and why, and when :). People usually have backups from the critical data which holds the value, so did I. But no backup/list/doc of the sw components. I had to rediscover the things and create a plan how to apply everything on a new OS. Here comes the over-hyped technologies.

### Ansible

There is a buzzword nowadays, which is: [infrastructure as code](https://en.wikipedia.org/wiki/Infrastructure_as_Code). The point is to not have the components settings and relations in mind (or better: in doc), but in machine processable definition files. Something like what Makefile does for software, these definitions are for OS components configurations. One of the popular solutions for this is [Ansible](https://www.ansible.com/). If you are a hardcode shell script guy, you might think what is the advance in these specialized solutions. I found these:

- Definition files are reusable as components
- Idempotent nature of the ansible modules avoid unnecessary execution
- Ansible is totally agentless: you don't have to have anything on the target machine, uses pure SSH

Of course, the real advantage of ansible is if you have dozens of servers or VM instances in the cloud, so using it for my router is overkill. Indeed. But once I produced my playbooks for the router, the configuration is self documented. No need to reinvent the wheel next time. Well, hopefully :).

Let's have an example: we want to set up a dhcp service with dnsmasq on a debian system. You can do it manually: use apt-get/aptitude to install the package and the dependencies, tweak the config file manually, restart the service. Or you can put the steps into a shell script if you want some self-documentation. In this case you have a chance to reproduce the tasks in the future. Now, with configuration automation tools you have to only declare the desired state of the system instead of scripting everything manually. Of course you express the steps in some sense, but due to the idempotent nature, whenever you run the playbook again, it will make changes only where the actual state of the target system differs from the desired state. Let's see the dnsmasq's task file:

```
- name: Install dnsmasq
  become: yes
  apt: pkg=dnsmasq state=installed

- name: Configure dnsmasq for dhcp
  become: yes
  notify: Restart dnsmasq
  template: src=dnsmasq.conf.j2 dest=/etc/dnsmasq.conf
```

What do we have here?! The first step is more or less straightforward: install the dnsmasq package and the dependencies with apt package manager. The become statement means it must run with superuser, with a preselected method (e.g. sudo). The task statement is more interesting. It refers to a template configuration dnsmasq.conf.j2. The content of the template can be the following:

```
dhcp-range={{lanNicName}},192.168.1.20,192.168.1.98,72h
interface={{lanNicName}}
```

Ansible's native template filling mechanism is called [jinja2](http://jinja.pocoo.org/docs/dev/). You can define a variable set separated from your tasks and feed when you start the playbook. The other stranger line is the notify statement. You can notify special tasks called handlers. They are invoked only if the caller task made real changes. This handler looks like:

```
- name: Restart dnsmasq
  become: yes
  service: name=dnsmasq state=restarted
```

You can collect these task+handler+template packages into so-called [roles](http://docs.ansible.com/ansible/playbooks_roles.html). I used a sandbox VM to develop the playbooks prior to the motherboard change to minimize the maintenance window of the production system :).

### Docker

If you keep your eyes open in the IT world, you most probably heard about [docker](https://www.docker.com/), or at least the hype of the containerization. It made a fundamental change in the world in the last 1-2 years. Companies like Spotify, Uber, SoundCloud, etc. actually use it from the time when docker was not even in production state. But why is this a real turnaround, what is behind the buzzword "[cloud-native](http://www.informationweek.com/cloud/platform-as-a-service/cloud-native-what-it-means-why-it-matters/d/d-id/1321539)"? On the lowest level, a well crafted running container instance is nothing more than a linux process, in correct isolation. Something how a process shall be isolated from the beginning of the multitasking operating systems. Comparing to cloud computing, here the process is detached from the underlying OS like VM is detached from the real hardware. Containers utilizes [linux namespaces](http://man7.org/linux/man-pages/man7/namespaces.7.html). The only shared part of the host OS and the container (and also between the containers) is the host's kernel. Linux kernel has a pretty stable APIs, so that the compatibility is not an issue.

- Mount namespace ensures that the process it entirely enclosed in it's image. It can not examine other processes/services playground in the filesystem. You can bind-mount files/directories from the host or network if the container needs external content or for the sake of persistency.
- Network namespace ensures that you can restrict the app's connectivity without global finetuned firewall rules in the OS.
- With PID namespace every application think that it is the only process in the system (actually the PID 1).

There are other namespaces but these are the most relevant.

Docker can be used on it's own, with CLI. Containerization in fact evolved to a totally flexible utilization of resources. Docker is only one brick here: it has tools on the top of it for orchestration, scheduling, scaling, healing, etc. The good container is designed to be totally independent of it's infrastructure and have a very short list of requirements. The possibilities here are the must-have configuration parameters (in the form of env.variables), the resource limits, other required services (linking to other containers like a DB on the network), the intentionally exposed listening tcp or udp ports for the service and some persistent storage. You can supply these directly in docker CLI, or you can express them in a YAML file and feed it to a small tool called [docker-compose](https://docs.docker.com/compose/overview/). There are special topics how ideally a container should look like ([small size](http://blog.xebia.com/how-to-create-the-smallest-possible-docker-container-of-any-image/), [readonly nature](http://www.projectatomic.io/blog/2015/12/making-docker-images-write-only-in-production/), etc.).

Getting back to the present project, for me docker helped to get ready-made self-containing apps. Most probably other people also do things like me, so most probably others already created containers for my need as well :).

### Putting all together

Ansible has a lot of built-in modules to interact with services, not only installing packages, copying and altering config files. So it can interact directly with docker-compose to start my containerized services. I ended up with a role collection and a top-level playbook with which I can anytime force my box to the desired state. You can check the stuff on [github](https://github.com/libesz/lhs-ansible). Some migration related things are hard to do with ansible like initialization of the raid arrays and renaming the network interfaces. So the box needed some natural preparation after the minimal OS was there, but not that much.

There are a few missing parts as I would like to monitor my containers and maybe build my own containers in some cases, but the workflow is now written in stone at least :).

Cheers!

&nbsp;

