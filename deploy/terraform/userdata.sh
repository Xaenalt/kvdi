#!/bin/bash

yum -y update
yum -y install jq htop

amazon-linux-extras install -y docker
systemctl enable docker
systemctl start docker

export INSTALL_K3S_SKIP_START="true"
export INSTALL_K3S_EXEC="server --no-deploy traefik"
export PROMETHEUS_OPERATOR_VERSION="v0.41.0"
export K3S_MANIFEST_DIR="/var/lib/rancher/k3s/server/manifests"

mkdir -p "${K3S_MANIFEST_DIR}"

curl -sfL https://get.k3s.io | sh -

ln -s /usr/local/bin/k3s /usr/sbin/k3s

curl -JL -o "${K3S_MANIFEST_DIR}/prometheus-operator.yaml" \
    https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/${PROMETHEUS_OPERATOR_VERSION}/bundle.yaml

tee "${K3S_MANIFEST_DIR}/kvdi.yaml" << EOF
apiVersion: helm.cattle.io/v1
kind: HelmChart
metadata:
  name: kvdi
  namespace: kube-system
spec:
  chart: kvdi
  repo: https://tinyzimmer.github.io/kvdi/deploy/charts
  targetNamespace: default
  valuesContent: |-
    vdi:
      spec:
        app:
          auditLog: true
          replicas: 3
        desktops:
          maxSessionLength: 5m
        auth:
          tokenDuration: 4h
          allowAnonymous: true
        metrics:
          serviceMonitor:
            create: true
          prometheus:
            create: true
          grafana:
            enabled: true
      templates:
        - metadata:
            name: ubuntu-xfce4
          spec:
            image: quay.io/tinyzimmer/kvdi:ubuntu-xfce4-demo
            imagePullPolicy: IfNotPresent
            resources:
              requests:
                cpu: 500m
                memory: 512Mi
              limits:
                cpu: 1000m
                memory: 1024Mi
            config:
              allowRoot: false
              init: systemd
            tags:
              os: ubuntu
              desktop: xfce4
              applications: minimal
EOF

systemctl start k3s

## To get the admin password from a booted instance
# sudo k3s kubectl get secret kvdi-admin-secret -o json | jq -r .data.password | base64 -d && echo