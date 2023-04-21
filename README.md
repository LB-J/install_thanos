# install_thanos

#### 介绍
###### 自动化部署thanos+prometheus+alertmanager，目前实现通过auto_install/install.py 脚本，可实现一键部署所有组件。
###### 通过hosts.yml 进行组件基础信息的配置，包括部署节点信息、alertmanager的 webhook等。之后会不断完善。
###### 大部分都来自于开源服务，个人进行了归纳整理，并提供了一些简单的工具，如有侵权，请及时告知，我将进行调整说明。 

#### 软件架构
###### 采用 docker-compose + python 脚本实现，thanos、alertmanager、prometheus 架构。
###### 版本信息：
##### 镜像版本
###### prometheus: bitnami/prometheus:2.43.0
###### thanos: thanosio/thanos:v0.31.0
###### alertmanager: bitnami/alertmanager:0.25.0
###### grafana: grafana/grafana:9.2.3


#### 安装教程
##### 依赖组件：docker-compose、docker、python


#### 使用说明
###### 主要引用 仓库地址如下：
https://github.com/prometheus
https://github.com/thanos-io/thanos
https://github.com/kubernetes/kube-state-metrics
2.  xxxx
3.  xxxx



#### 特技
