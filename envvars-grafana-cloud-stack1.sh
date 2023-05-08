# Grafana Cloud Configuration

# Loki Remote Endpoint
# See: https://grafana.com/docs/loki/latest/api/
export GRAFANA_LOGS_HOST="logs-prod-us-central1.grafana.net"
export GRAFANA_LOGS_USERNAME="123456"
export GRAFANA_LOGS_API_KEY="eyJrIjoiYTkzYzM4OWIxYTQxMzM0OThiYTE4ZmZkMDc2N2Q4Y2Q5ZTMzOGQ5YyIsIm4iOiJTYXNoYSIsImlkIjo4NDkxNTV9"
#
export GRAFANA_LOGS_QUERY_URL="https://$GRAFANA_LOGS_HOST/loki/api/v1"
export GRAFANA_LOGS_WRITE_URL="https://$GRAFANA_LOGS_HOST/loki/api/v1/push"