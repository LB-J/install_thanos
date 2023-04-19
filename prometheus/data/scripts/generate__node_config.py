import json


def generate_json():
    result = []

    #  hosts 文件中列按顺序分别为：主机名、内网IP、公网映射IP
    with open('hosts', 'r') as f:
        content = f.readlines()
        for i in content:
            x = i.split()
            result.append(
                {
                    "labels": {
                        "cluster": "product",
                        "hostname": x[0],
                        "ip_address": x[1],
                        "eip": x[2],
                        # group 从主机名中进行获取，如无法从主机名中获取，请根据实际情况调整hosts 文件格式，添加group 列
                        "group": "-".join(x[0].split('-')[1:-2]),
                        "type": "node_exporter",
                    },
                    "targets": [
                        "%s:9100" % x[1]
                    ]
                }
            )

    save_file_name = open("./hosts_file/prometheus_%s.json" % 'test', "w")
    json.dump(result, save_file_name, indent=4)


generate_json()
