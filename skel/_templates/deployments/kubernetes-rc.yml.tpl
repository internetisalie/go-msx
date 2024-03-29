---
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: vms
  name: ${app.name}
  labels:
    app: ${app.name}
spec:
#if GENERATOR_SP
  replicas: "{{ deployment_mode_env[deployment_mode|lower]['replica_count']['servicepack_count'] }}"
#else GENERATOR_SP
  replicas: {{ deployment_mode_env[deployment_mode|lower]['replica_count']['${app.name}'] }}
#endif GENERATOR_SP
  selector:
    matchLabels:
      app: ${app.name}
      group: ${kubernetes.group}
      consul-gossip: allow
  template:
    metadata:
      name: ${app.name}
      labels:
        app: ${app.name}
        group: ${kubernetes.group}
        consul-gossip: allow
      annotations:
        tagprefix: logfmt
    spec:
      serviceAccountName: ${app.name}
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - ${app.name}
              topologyKey: kubernetes.io/hostname
{% if cloud == 'aws' %}
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - ${app.name}
              topologyKey: topology.kubernetes.io/zone
{% endif %}
      containers:
        - name: consul
          image: {{ consul_image }}:{{ consul_version }}
          command:
            - consul
            - agent
            - -bind=0.0.0.0
            - -client=0.0.0.0
            - -datacenter={{ consul_dc }}
            - -retry-join=consul.service.consul
            - -data-dir=/consul/data
            - -config-dir=/consul/config
          volumeMounts:
            - mountPath: /consul/config
              name: phi
        - name: ${app.name}
          image: {{ ${app.name}_image }}:{{ ${app.name}_version }}
          command:
            - "/usr/bin/${app.name}"
            - --profile
            - production
          livenessProbe:
            httpGet:
              path: ${server.contextPath}/admin/alive
              port: ${server.port}
            initialDelaySeconds: 300
            periodSeconds: 30
          readinessProbe:
            httpGet:
              path: ${server.contextPath}/admin/health
              port: ${server.port}
            initialDelaySeconds: 0
            periodSeconds: 30
          startupProbe:
            httpGet:
              path: ${server.contextPath}/admin/health
              port: ${server.port}
            failureThreshold: 45
            periodSeconds: 5
          resources:
            requests:
              memory: "64Mi"
              cpu: "500m"
            limits:
              memory: "256Mi"
              cpu: "2000m"
          env:
            - name: PROFILE
              value: production
            - name: SPRING_CLOUD_CONSUL_HOST
              value: "localhost"
            - name: SPRING_CLOUD_CONSUL_SCHEME
              value: "{{ vault_scheme }}"
            - name: SPRING_CLOUD_CONSUL_PORT
              value: "8500"
            - name: SPRING_CLOUD_CONSUL_CONFIG_ACLTOKEN
              valueFrom:
                secretKeyRef:
                  name: msxconsul
                  key: token
            - name: SPRING_CLOUD_VAULT_HOST
              value: "vault.service.consul"
            - name: SPRING_CLOUD_VAULT_PORT
              value: "8200"
            - name: SPRING_CLOUD_VAULT_SCHEME
              value: "{{ vault_scheme }}"
            - name: SPRING_CLOUD_VAULT_TOKEN-SOURCE_SOURCE
              value: "kubernetes"
            - name: SPRING_CLOUD_VAULT_TOKEN-SOURCE_KUBERNETES_ROLE
              value: "${app.name}"
            - name: SPRING_REDIS_SENTINEL_ENABLE
              value: "true"
          ports:
            - containerPort: ${server.port}
          volumeMounts:
            - mountPath: /etc/ssl/certs/ca-certificates.crt
              name: rootcert
            - mountPath: /keystore
              name: keystore
      volumes:
        - hostPath:
            path: /etc/ssl/certs/ca-bundle.crt
          name: rootcert
        - hostPath:
            path: /data/vms/keystore/
          name: keystore
        - configMap:
            name: phi
          name: phi
