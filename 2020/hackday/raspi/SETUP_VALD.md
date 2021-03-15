# How to install Vald on raspberry cluster


## Preparation

1. Install helm

    ```bash
    curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
    chmod 700 get_helm.sh
    ./get_helm.sh
    ```

1. Install k8s metrics server

    ```bash
    kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
    ```
1. Install minio db

    Create PersistentVolumeClaim
    
    ```bash
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      name: s3-pvc
    spec:
      storageClassName: local-path
      accessModes:
        - ReadWriteOnce
      resources:
        requests:
          storage: 10Gi
    ```
    
    Create Service
    
    ```bash
    apiVersion: v1
    kind: Service
    metadata:
      name: minio
    spec:
      ports:
        - port: 9000
      selector:
        app: minio
      clusterIP: None
    ```

    
    Create Deployment
    
    ```bash
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      labels:
        app: minio
      name: minio
    spec:
      progressDeadlineSeconds: 600
      replicas: 1
      selector:
        matchLabels:
          app: minio
      template:
        metadata:
          labels:
            app: minio
        spec:
          containers:
            - name: minio
              image: minio/minio
              imagePullPolicy: Always
              command:
                - "/usr/bin/docker-entrypoint.sh"
                - "server"
                - "/data"
              env:
                - name: MINIO_ACCESS_KEY
                  value: ACCESSKEY
                - name: MINIO_SECRET_KEY
                  value: SECRETKEY
              ports:
                - containerPort: 9000
                  name: minio
                  protocol: TCP
              resources:
                requests:
                  cpu: 100m
                  memory: 100Mi
          dnsPolicy: ClusterFirst
          restartPolicy: Always
          schedulerName: default-scheduler
          securityContext: {}
          terminationGracePeriodSeconds: 30
    ```

## Install Vald

1. Add vald charts

    ```bash
    helm repo add vald https://vald.vdaas.org/charts
    ```

1. Create values.yaml for helm

    ```bash
    cat << EOF > values.yaml
    ---
    
    ```
    
2. Generate k8s from helm template

    ```bash
    helm template vald-cluster vald/vald --values values.yaml --output-dir .
    ```
    
1. Install vald-agent-ngt
    
    ```bash
    kubectl apply -f vald/templates/agent
    ```
    
1. Install vald-discoverer
    
    ```bash
    kubectl apply -f vald/templates/discoverer
    ```

1. Install vald-manager-index
    
    ```bash
    kubectl apply -f vald/templates/manager/index
    ```

1. Install vald-lb-gateway
    
    ```bash
    kubectl apply -f vald/templates/gateway/lb
    ```
