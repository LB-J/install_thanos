# 采集节点 label 配置
# node_exporter label
# 环境、集群、大的分组含义（例如：product (生产环境)， hadoop （集群））。
cluster
# 节点的主机名，主机名的定义最好跟部署的环境、服务和节点的IP关联。例如：部署 hadoop datanode 服务的192.168.20.18，主机名定义为：prod-datanode-20-18。
# 此处只是示例，具体可根据实际情况，也可以直接使用主机的主机名。只是用于做为一个辨识。
hostname
# 主机的IP地址，部分用户环境可能存在多个IP地址，此处请填写主机实际的IP，并和prometheus 进行通信的IP。额外的IP可以增加label 来进行完善
ip_address
# 分组，表示再cluster 下再分。例如：product环境下的按照服务分（mysql，redis、middleware、application、k8s-master、k8s-node）等。
# 对于集群，hadoop 下面有datanode、namenode、journalnode、zookeeper等服务。
group
# 此处是我定义的指标采集的方式，例如：node_exporter,jmx_exporter,mysql_exporter
type
## 对于 cluster、hostname、ip_address、group label 为必选项。（如果使用该项目下grafana的template，如果不添加这些label，需要对导入的大屏进行调整）

node_exporter 示例：
[
{
    // cluster、hostname、ip_address、group 必须有
    // eip 为服务器的上层网络映射IP，（可根据实际情况进行定义或删除）
        "labels": {
            "cluster": "product",
            "hostname": "product-nginx-14-07",
            "ip_address": "192.168.14.7",
            "eip": "172.20.101.21",
            "group": "nginx",
            "type": "node_exporter"
        },
        "targets": [
            "192.168.14.7:9100"
        ]
    }
]
mysql_exporter 示例：
       {
        "labels": {
        // mysql 的用途。例如：某某应用库、某某配置库、某某专题库等，必填项
            "db_application": "application",
        // 指标采集的类型，mysql 为mysql_exporter ，可选项
            "type": "mysql_exporter",
        // 数据库的类型，主库：master、从库：slave、延迟从库：slow_slave。必填项
            "db_type": "slave",
        // 服务名称：mysql 必填项
            "server": "mysql",
        // 所属环境，例如：生产、测试、预发等。必填项
            "cluster": "product",
        // mysql 端口  必填项
            "port": "3307",
        // mysql 部署节点IP ，必填项
            "ip_address": "192.168.10.4"
        },
        "targets": [
            "192.168.10.4:9105"
        ]
    }
]
