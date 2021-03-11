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

1. provisioning ubuntu (update / upgrade)

    ```bash
    sudo du -sh /var/cache/apt/archives sudo rm -rf /var/cache/apt
    sudo mkdir -p /var/cache/apt/archives/partial
    sudo DEBIAN_FRONTEND=noninteractive apt -y clean
    sudo DEBIAN_FRONTEND=noninteractive apt -y autoremove
    sudo DEBIAN_FRONTEND=noninteractive apt -y update
    sudo DEBIAN_FRONTEND=noninteractive apt -y upgrade
    sudo DEBIAN_FRONTEND=noninteractive apt -y full-upgrade
    sudo DEBIAN_FRONTEND=noninteractive apt -y clean
    sudo DEBIAN_FRONTEND=noninteractive apt -y autoremove --purge
    sudo du -sh /var/cache/apt/archives
    sudo rm -rf /var/cache/apt
    sudo mkdir -p /var/cache/apt/archives/partial
    ```

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

1. create config file for using fixed IP addresses.

    - IP Address list
        - k8s-master: 192.168.13.101
        - k8s-node1:  192.168.13.102
        - k8s-node2:  192.168.13.103

    ```bash
    sudo vim /etc/netplan/99_configfile.yaml
    ===
    network:
      version: 2
      renderer: networkd
      ethernets:
        eth0:
          dhcp4: false
          dhcp6: false
          addresses: [<ip addresses>/24]
          gateway4: 192.168.13.1
          nameservers:
            addresses: [192.168.13.1, 8.8.8.8, 8.8.4.4]
    ```

1. apply netplan

    ```bash
    # apply /etc/netplan/99_config.yaml
    sudo netplan apply
    # Verify
    ip addr
    ```

1. set hostname

    ```bash
    hostnamectl set-hostname [k8s-master/k8s-node1/k8s-node2]
    ```

1. edit `/etc/hosts`

    ```bash
    sudo vim /etc/hosts
    ===
    # choose hostname from k8s-master, k8s-node1 or k8s-node2.
    127.0.0.1    [k8s-master/k8s-node1/k8s-node2]
    ::1        localhost ip6-localhost ip6-loopback
    ff02::1        ip6-allnodes
    ff02::2        ip6-allrouters

    127.0.1.1    raspberrypi
    192.168.13.101 k8s-master
    192.168.13.102 k8s-node1
    192.168.13.103 k8s-node2
    ```

## Install k3s

### setup for master node

1. Install

   Install k3s with disable flannel, k3s default network policy, and change the pod IP CIDR.

   ```bash
   curl -sfL https://get.k3s.io | INSTALL_K3S_CHANNEL="latest" INSTALL_K3S_EXEC="--flannel-backend=none --disable-network-policy --cluster-cidr=192.168.0.0/16" sh -
   ```

1. Copy configuration file

   ```bash
   sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
   sudo chmod 755 /home/ubuntu/.kube/config
   ```

1. Get Token

   ```bash
   sudo cat /var/lib/rancher/k3s/server/node-token
   ```

### setup for extra node

1. Install

   Install k3s and you should set correct server IP and token of master node

   ```bash
   curl -sfL https://get.k3s.io | K3S_URL={master node IP} K3S_TOKEN={master node TOKEN} sh -
   ```

1. Install calico operator

   ```bash
   kubectl create -f https://docs.projectcalico.org/manifests/tigera-operator.yaml
   ```

1. Install crd for calico operator

   ```bash
   kubectl create -f https://docs.projectcalico.org/manifests/custom-resources.yaml
   ```
