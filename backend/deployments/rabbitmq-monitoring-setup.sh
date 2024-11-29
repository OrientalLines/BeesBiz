cat << EOF | kubectl apply -f -
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: rabbitmq-monitor
  namespace: beesbiz-monitoring
  labels:
    release: monitoring
spec:
  endpoints:
  - interval: 30s
    port: prometheus
    path: /metrics
  selector:
    matchLabels:
      app.kubernetes.io/name: rabbitmq-cluster
  namespaceSelector:
    matchNames:
    - beesbiz-rabbitmq
EOF
