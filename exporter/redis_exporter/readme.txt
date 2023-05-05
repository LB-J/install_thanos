## 采用官方镜像部署
https://github.com/oliver006/redis_exporter
https://hub.docker.com/r/oliver006/redis_exporter

# config example,add prometheus.yml
scrape_configs:
  - job_name: 'redis_exporter_targets'
    file_sd_configs:
      - files:
        - redis-instances.json
    metrics_path: /scrape
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 192.168.14.11:9121

# redis-instances.json example
[
  {
    "targets": [
     "redis://192.168.14.16:7001",
     "redis://192.168.14.16:7002",
     "redis://192.168.14.17:7003",
     "redis://192.168.14.17:7004",
     "redis://192.168.14.18:7005",
     "redis://192.168.14.18:7006"
     ],
    "labels": {
    "cluster": "test",
    "server": "redis"
     }
  }
]
