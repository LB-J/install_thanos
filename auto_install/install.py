import sys
import subprocess
from ruamel import yaml
import os


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


def grafana():
    run_dir = "../grafana/"
    docker_run(run_dir)


def alert_manager(config):
    template_file = "../alertmanager/data/alertmanager/conf/alertmanager-template.yml"
    config_file = "../alertmanager/data/alertmanager/conf/alertmanager.yml"
    run_dir = "../alertmanager/"
    conf = load_yaml(template_file)
    conf["receivers"][0]["webhook_configs"][0]["url"] = config["monitor"]["alert_manager"]["web_hook"][0]['url']
    dump_yaml(config_file, conf)
    docker_run(run_dir)


def prometheus(config):
    template_file = "../prometheus/data/prometheus/conf/prometheus-template.yml"
    config_file = "../prometheus/data/prometheus/conf/prometheus.yml"
    run_dir = "../prometheus"
    conf = load_yaml(template_file)
    conf["global"]["external_labels"] = config["monitor"]["prometheus"]["external_labels"]
    dump_yaml(config_file, conf)
    docker_run(run_dir)


def thanos_query(config):
    template_file = "../thanos_query/data/thanos-query/conf/query-template.yml"
    config_file = "../thanos_query/data/thanos-query/conf/query.yml"
    run_dir = "../thanos_query"
    conf = load_yaml(template_file)
    conf[0]["targets"] = config['monitor']['prometheus']['hosts']
    dump_yaml(config_file, conf)
    docker_run(run_dir)


def thanos_rule(config):
    query_template_file = "../thanos_rule/data/thanos-rule/conf/query-template.yml"
    query_config_file = "../thanos_rule/data/thanos-rule/conf/query.yml"
    alert_template_file = "../thanos_rule/data/thanos-rule/conf/alertmanager-template.yml"
    alert_config_file = "../thanos_rule/data/thanos-rule/conf/alertmanager.yml"
    run_dir = "../thanos_rule"
    query_conf = load_yaml(query_template_file)
    query_conf[0]["targets"] = config['monitor']['thanos_query']['hosts']
    dump_yaml(query_config_file, query_conf)
    alert_conf = load_yaml(alert_template_file)
    alert_conf["alertmanagers"][0]["static_configs"] = config['monitor']['alert_manager']['hosts']
    dump_yaml(alert_config_file, alert_conf)
    docker_run(run_dir)


if __name__ == "__main__":
    host_config = parse_config(config_file=os.path.join(os.path.dirname(__file__), 'hosts.yml'))
    prometheus(host_config)
    thanos_query(host_config)
    thanos_rule(host_config)
    alert_manager(host_config)
