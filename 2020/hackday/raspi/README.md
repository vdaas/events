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

### setup for master node

1. Install

Install k3s with disable flannel and k3s default network policy and change the pod ip CIDR.

```bash
curl -sfL https://get.k3s.io | INSTALL_K3S_CHANNEL="latest" INSTALL_K3S_EXEC="--flannel-backend=none --disable-network-policy --cluster-cidr=192.168.0.0/16" sh -
```

2. Copy configuration file

```bash
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chmod 755 /home/ubuntu/.kube/config
```

3. Get token

```
sudo cat /var/lib/rancher/k3s/server/node-token
```


### setup for extra node

1. Install

Install k3s and you should set correct server ip and token of master node

```bash
curl -sfL https://get.k3s.io | K3S_URL={master node ip} K3S_TOKEN={master node token} sh -
```

2. Install calico operator

```bash
kubectl create -f https://docs.projectcalico.org/manifests/tigera-operator.yaml
```

3. Install crd for calico operator

```bash
kubectl create -f https://docs.projectcalico.org/manifests/custom-resources.yaml
```

