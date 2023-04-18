# -*- coding: utf-8 -*-
import yaml
import requests


def generate_config(data, group_name, rules_path):
    content = {
        "groups": [{
            "name": group_name,
            "rules": []
        }]}
    rule_list = []
    rule_level = {
        0: "INFO",
        1: "WARNING",
        2: "CRITICAL",
    }
    for rule in data:
        rule_list.append({
            "alert": rule['name'],
            "expr": rule['rule'],
            "for": "%dm" % rule['check_period'],
            "labels": {"rule_id": rule['id'],
               "severity": rule_level[rule['level']],
                "alert_host": "{{ $labels.instance }}"},
            "annotations": {
                "description": rule['description']
            }})
    content['groups'][0]['rules'] = rule_list

    with open('%s/%s.yml' % (rules_path, group_name), 'w') as f:
        yaml.dump(content, f, default_flow_style=False, encoding='utf-8', allow_unicode=True)
    return True


def reload_prom(url):
    req = requests.post("%s/-/reload" % url)
    content = req.status_code
    print(content)

# 获取所有规则分组
def get_group(base_url):
    url = "%s/rulemanager/api/rule/rule_group/" % base_url
    req = requests.get(url)
    result = req.json()
    return result

def get_rule():
    #p_address = "http://127.0.0.1:8000"
    p_address = "http://rulemanager.prod.ntsyzl.com"
    prom_address = "http://127.0.0.1:19193"
    rules_path = "/export/cloud/thanos-rule/rules"
    group_list = get_group(p_address)
    for i in group_list['results']:
        url = "%s/rulemanager/api/rule/common/rules/?group=%d&status=true" % (p_address, i['id'])
        res = requests.get(url)
        content = res.json()
        print(content)
        generate_config(data=content['results'], group_name=i['name'], rules_path=rules_path)
    reload_prom(prom_address)


get_rule()
