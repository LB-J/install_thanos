import sys
import subprocess
from ruamel import yaml
import os
import socket


# 获取当前节点IP，暂时无用
def get_host_ip():
    get_ip_command = "hostname -I | grep -P '(\d+\.){3}\d+'"
    res = subprocess.run(get_ip_command, timeout=30, shell=True,stdout=subprocess.PIPE).stdout
    return res.rstrip().split()

def parse_config(config_file):
    with open(config_file, 'r') as file:
        config = yaml.load(file, Loader=yaml.Loader)
        return config


def load_yaml(f):
    try:
        with open(f, 'r+') as f:
            conf = yaml.safe_load(f)
            return conf
    except Exception as e:
        print("Load %s file is filed: %s" % (f, e))
        sys.exit(1)


def dump_yaml(f, content):
    try:
        with open(f, 'w+') as f:
            yaml.dump(content, f, Dumper=yaml.RoundTripDumper, default_flow_style=False)
        return True
    except Exception as e:
        print("dump %s file failed:%s" % (f, e))
        sys.exit(1)


def docker_run(start_dir):
    command = "cd %s ; docker-compose up -d" % start_dir
    subprocess.run(command, timeout=600, shell=True)


def grafana(config):
    # run_dir = "../grafana/"
    # docker_run(run_dir)
    pass


def alert_manager(config):
    template_file = "../alertmanager/data/alertmanager/conf/alertmanager-template.yml"
    config_file = "../alertmanager/data/alertmanager/conf/alertmanager.yml"
    conf = load_yaml(template_file)
    conf["receivers"][0]["webhook_configs"][0]["url"] = config["monitor"]["alert_manager"]["web_hook"][0]['url']
    dump_yaml(config_file, conf)


def prometheus(config):
    template_file = "../prometheus/data/prometheus/conf/prometheus-template.yml"
    config_file = "../prometheus/data/prometheus/conf/prometheus.yml"
    conf = load_yaml(template_file)
    conf["global"]["external_labels"] = config["monitor"]["prometheus"]["external_labels"]
    dump_yaml(config_file, conf)


def thanos_query(config):
    template_file = "../thanos_query/data/thanos-query/conf/query-template.yml"
    config_file = "../thanos_query/data/thanos-query/conf/query.yml"
    conf = load_yaml(template_file)
    conf[0]["targets"] = config['monitor']['prometheus']['hosts']
    dump_yaml(config_file, conf)


def thanos_rule(config):
    query_template_file = "../thanos_rule/data/thanos-rule/conf/query-template.yml"
    query_config_file = "../thanos_rule/data/thanos-rule/conf/query.yml"
    alert_template_file = "../thanos_rule/data/thanos-rule/conf/alertmanager-template.yml"
    alert_config_file = "../thanos_rule/data/thanos-rule/conf/alertmanager.yml"
    query_conf = load_yaml(query_template_file)
    query_conf[0]["targets"] = config['monitor']['thanos_query']['hosts']
    dump_yaml(query_config_file, query_conf)
    alert_conf = load_yaml(alert_template_file)
    alert_conf["alertmanagers"][0]["static_configs"] = config['monitor']['alert_manager']['hosts']
    dump_yaml(alert_config_file, alert_conf)


def generate_config(conf, serve_list):
    for server_name in serve_list:
        eval(server_name)(conf)


def start_server(server_list):
    for serve_name in server_list:
        work_dir = "../%s" % serve_name
        docker_run(work_dir)


if __name__ == "__main__":
    host_config = parse_config(config_file=os.path.join(os.path.dirname(__file__), 'hosts.yml'))
    # 当前本机 IP 地址通过hosts.yml 文件进行配置
    local_ip = host_config["monitor"]["localhost"]["ip"]
    install_list = []
    query_hosts = [i.split(":")[0] for i in host_config["monitor"]["thanos_query"]["hosts"]]
    rule_hosts = host_config["monitor"]["thanos_rule"]["hosts"]
    prometheus_hosts = [i.split(":")[0] for i in host_config["monitor"]["prometheus"]["hosts"]]
    alert_hosts = [i.split(":")[0] for i in host_config["monitor"]["alertmanager"]["hosts"]]
    grafana_hosts = host_config["monitor"]["grafana"]["hosts"]
    if query_hosts in local_ip:
        install_list.append("thanos_query")
    if rule_hosts in local_ip:
        install_list.append("thanos_rule")
    if local_ip in prometheus_hosts:
        install_list.append("prometheus")
    if local_ip in alert_hosts:
        install_list.append("alertmanager")
    if local_ip in grafana_hosts:
        install_list.append("grafana")
    # 生成配置文件：
    generate_config(conf=host_config, serve_list=install_list)
    # 启动服务
    start_server(server_list=install_list)




