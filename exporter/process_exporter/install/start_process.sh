#!/bin/bash
# 定时任务方式
echo "*/1 * * * * /var/opt/process_monitor/bin/process_exporter --config-path /var/opt/process_monitor/config/process.yaml" >> /var/spool/cron/root
# 也可以修改配置中的daemon 参数，直接运行，源代码在 source