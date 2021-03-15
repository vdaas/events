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

1. Create Job

    ```bash
    apiVersion: batch/v1
    kind: Job
    metadata:
      name: minio-make-bucket
    spec:
      template:
        spec:
          containers:
            - name: mc
              image: minio/mc
              imagePullPolicy: Always
              command:
                - /bin/sh
                - -c
                - |
                  mc alias set minio ${ENDPOINT} ${MINIO_ACCESS_KEY} ${MINIO_SECRET_KEY} --api S3v4
                  mc mb minio/vald-minio
              env:
                - name: ENDPOINT
                  value: http://minio:9000
                - name: MINIO_ACCESS_KEY
                  value: ACCESSKEY
                - name: MINIO_SECRET_KEY
                  value: SECRETKEY
          restartPolicy: Never
    ```
    
## Install Vald

1. Add vald charts

    ```bash
    helm repo add vald https://vald.vdaas.org/charts
    ```

1. Create values.yaml for helm

```
defaults:
  time_zone: Asia/Tokyo
  logging:
    format: raw
    level: debug
    logger: glg
  image:
    tag: v1.0.1
  observability:
    enabled: false
​
discoverer:
  clusterRole:
    enabled: true
  clusterRoleBinding:
    enabled: true
  discoverer:
    cache:
      enabled: true
      expire_duration: 2s
      expired_cache_check_duration: 200ms
    discovery_duration: 200ms
    name: ""
    namespace: _MY_POD_NAMESPACE_
  env:
    - name: MY_POD_NAMESPACE
      valueFrom:
        fieldRef:
          fieldPath: metadata.namespace
  hpa:
    enabled: false
​
meta:
  enabled: false
​
manager:
  backup:
    enabled: false
  compressor:
    enabled: false
  index:
    enabled: true
    maxReplicas: 1
    minReplicas: 1
​
gateway:
  lb:
    enabled: true
    maxReplicas: 3
    minReplicas: 3
    gateway_config:
      index_replica: 2
  backup:
    enabled: false
  meta:
    enabled: false
  vald:
    enabled: false
​
agent:
  minReplicas: 1
  podManagementPolicy: Parallel
  hpa:
    enabled: false
  resources:
    requests:
      cpu: 100m
      memory: 50Mi
  volumes:
    - name: ngt-index
      emptyDir: {}
  volumeMounts:
    - name: ngt-index
      mountPath: /var/ngt
  ngt:
    auto_index_duration_limit: 3m
    auto_index_check_duration: 1m
    auto_index_length: 1000
    dimension: 512
    index_path: /var/ngt/index
    enable_in_memory_mode: false
  sidecar:
    enabled: true
    initContainerEnabled: true
    env:
      - name: AWS_ACCESS_KEY
        value: ACCESSKEY
      - name: AWS_SECRET_ACCESS_KEY
        value: SECRETKEY
    resources:
      requests:
        cpu: 100m
        memory: 100Mi
    config:
      filename: vald-agent-ngt-index
      blob_storage:
        storage_type: "s3"
        bucket: "vald-minio"
        s3:
          endpoint: "http://minio.default.svc.cluster.local:9000"
          region: "us-east-1"
          force_path_style: true
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
