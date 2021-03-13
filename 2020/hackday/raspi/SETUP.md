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
    sudo du -sh /var/cache/apt/archives
    sudo rm -rf /var/cache/apt
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
   sudo vim /boot/firmware/cmdline.txt
   # add at the end of 1st line
   cgroup_enable=cpuset cgroup_memory=1 cgroup_enable=memory
   ```

1. ~~disable services~~

   ```bash
   systemctl status dphys-swapfile.service
   sudo systemctl stop dphys-swapfile.service
   systemctl status dphys-swapfile.service
   sudo systemctl disable dphys-swapfile.service
   ```

1. disable swap
   ```bash
   sudo swapoff -a
   ```
   
3. disable firewall

   ```bash
   sudo ufw disable
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
    sudo hostnamectl set-hostname [k8s-master/k8s-node1/k8s-node2]
    ```

1. edit `/etc/hosts`

    ```bash
    sudo vim /etc/hosts
    ===
    # add
    192.168.13.101 k8s-master
    192.168.13.102 k8s-node1
    192.168.13.103 k8s-node2
    ```
    
1. edit `/etc/sysctl.conf`

    ```bash
    sudo vim /etc/sysctl.conf
    ===
    # add
    net.ipv6.conf.all.disable_ipv6 = 1
    net.ipv6.conf.default.disable_ipv6 = 1
    net.ipv6.conf.eth0.disable_ipv6 = 1
    net.ipv6.conf.lo.disable_ipv6 = 1
    ```
1. reload sysctl

    ```bash
    sudo sysctl -p
    ```

1. Prevent `iptables` from using the `nftables` backend

    ```bash
    sudo apt-get install -y iptables arptables ebtables
    
    sudo update-alternatives --set iptables /usr/sbin/iptables-legacy
    sudo update-alternatives --set ip6tables /usr/sbin/ip6tables-legacy
    sudo update-alternatives --set arptables /usr/sbin/arptables-legacy
    sudo update-alternatives --set ebtables /usr/sbin/ebtables-legacy
    ```  
    
1. reboot

    ```bash
    sudo reboot
    ```

## Install k3s

### setup for master node

1. Install

   Install k3s with disable flannel, k3s default network policy, and change the pod IP CIDR.

   ```bash
   curl -sfL https://get.k3s.io | INSTALL_K3S_CHANNEL="latest" K3S_KUBECONFIG_MODE="644" INSTALL_K3S_EXEC="--flannel-backend=none --disable-network-policy --cluster-cidr=192.168.0.0/24" sh -
   ```

1. Copy configuration file

   ```bash
   sudo mkdir ~/.kube
   sudo cp /etc/rancher/k3s/k3s.yaml $HOME/.kube/config
   sudo chmod 755 $HOME/.kube/config
   sudo chmod 755 /etc/rancher/k3s/k3s.yaml
   ```

1. Get Token

   ```bash
   sudo cat /var/lib/rancher/k3s/server/node-token
   ```

### setup for extra node

1. Install

   Install k3s and you should set correct server IP and token of master node

   ```bash
   curl -sfL https://get.k3s.io | K3S_URL=https://192.168.13.101:6443 K3S_TOKEN={master node TOKEN} sh -
   ```

### setup cni plugin

1. Install calico from master node
   ```bash
   kubectl apply -f https://docs.projectcalico.org/manifests/calico.yaml
   ```
