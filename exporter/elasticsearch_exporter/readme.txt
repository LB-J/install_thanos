## 采用官方镜像部署
https://github.com/prometheus-community/elasticsearch_exporter
https://hub.docker.com/r/bitnami/elasticsearch-exporter


# redis-instances.json example
[
    {
        "labels": {
            "cluster": "product",
            "server": "elasticsearch"
        },
        "targets": [
            "192.168.10.2:9114"

        ]
    }

]