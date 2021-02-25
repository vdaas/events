# Setup Vald on k3s with raspberry pi

- Spec:

  - Rspberry Pi 4(4GBi RAM)
  - OS: Ubuntsu 20.x.x

- Login User/Password

  - user: ubuntu
  - password: vdaas/vald

- TOC
  - Setup Raspberry Pi

## Preparation

1. enable cgroup

   ```bash
   sudo vim /boot/cmdline.txt
   # add at the end of 1st line
   cgroup_enable=cpuset cgroup_memory=1 cgroup_enable=memory swapaccount=1
   ```

1. disable services

   ```bash
   systemctl status dphys-swapfile.service
   sudo systemctl stop dphys-swapfile.service
   systemctl status dphys-swapfile.service
   sudo systemctl disable dphys-swapfile.service
   ```

## Install k3s

1. xxx

1. yyy

1. zzz
