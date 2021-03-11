# How to use Vald on raspberry cluster

- Setup document: Click [here](./SETUP.md)

## List of each cluster info.

|hostname|ip address|login user| login pass|
|:---|:---|:---|:---|
|k8s-master|192.168.13.101|ubuntu|vdaas/vald|
|k8s-node1|192.168.13.102|ubuntu|vdaas/vald|
|k8s-node2|192.168.13.103|ubuntu|vdaas/vald|

## Setup Wi-Fi

1. Connect Wi-Fi router

  1. Turn on your Vald on raspberry cluster.
  1. Confirm your router's SSID.
  1. Enter password.

1. Connect Wi-Fi router and Internet.

  1. Select Wi-Fi which you'd like to connect.
  1. Confirm and Connect Wi-Fi.

## Login nodes via Wi-Fi router

- When you'd like to operate, please try to `ssh` at first.

```bash
# login k8s-master node
ssh ubuntu@192.168.13.101
```
