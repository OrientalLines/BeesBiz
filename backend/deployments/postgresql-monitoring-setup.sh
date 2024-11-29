
cat << EOF | kubectl apply -f -
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: postgresql-monitor
  namespace: beesbiz-monitoring
  labels:
    release: monitoring    # This must match Prometheus's serviceMonitorSelector
spec:
  endpoints:
  - interval: 30s
    port: metrics
    path: /metrics
  selector:
    matchLabels:
      cnpg.io/cluster: postgresql
  namespaceSelector:
    matchNames:
    - beesbiz-data
EOF


cat << EOF | kubectl apply -f -
apiVersion: v1
kind: Service
metadata:
  name: postgresql-metrics
  namespace: beesbiz-data
  labels:
    cnpg.io/cluster: postgresql
spec:
  ports:
  - name: metrics
    port: 9187
    targetPort: metrics
  selector:
    cnpg.io/cluster: postgresql
EOF

kubectl patch cluster postgresql -n beesbiz-data --type merge -p '
{
  "spec": {
    "monitoring": {
      "enablePodMonitor": true,
      "customQueriesConfigMap": [
        {
          "name": "cnpg-default-monitoring",
          "key": "queries"
        }
      ],
      "disableDefaultQueries": false
    }
  }
}'


cat << EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: postgres-metrics-config
  namespace: beesbiz-data
data:
  queries.yaml: |
    pg_database:
      query: "SELECT pg_database.datname, pg_database_size(pg_database.datname) as size_bytes FROM pg_database"
      metrics:
        - datname:
            usage: "LABEL"
            description: "Name of the database"
        - size_bytes:
            usage: "GAUGE"
            description: "Disk space used by the database"
EOF
