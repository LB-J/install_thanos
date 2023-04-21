# k8s 指标采集，使用了kube-state-metrics 项目。git地址如下：
https://github.com/kubernetes/kube-state-metrics/
# 安装方法如下：
git clone git@github.com:kubernetes/kube-state-metrics.git
cd kube-state-metrics/examples/standard/
kubectl apply -f ./*

# 注意事项：
如果默认镜像无法下载，可替换为hub.docker.com 仓库镜像，地址如下：
https://hub.docker.com/r/bitnami/kube-state-metrics
