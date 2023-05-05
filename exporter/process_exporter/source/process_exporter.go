package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

const (
	namespace = "bigdata_process"
)

var monitor = []string{"rss", "vms", "mem_locked", "mem_swap", "mem_stack", "cpu_percent", "num_fds", "read_bytes", "write_bytes", "num_threads",
	"total_ESTABLISHED", "total_TIME_WAIT", "total_CLOSE_WAIT", "total_LISTEN", "total_SYN_SENT", "total_SYN_RECV", "start_time", "up"}
var NetMonitor = []string{"ESTABLISHED", "TIME_WAIT", "CLOSE_WAIT", "LISTEN", "SYN_SENT", "SYN_RECV"}
var PortMonitor = []string{"port_state"}

type ProcessCollector struct {
	metrics map[string]ProcessMetrics
}

type ProcessMetrics struct {
	desc          *prometheus.Desc
	extract       func(string) float64
	extractLabels func(s string) []string
	valType       prometheus.ValueType
}
type Config struct {
	ProcessList []string `yaml:"check_processes"`
	BindPort    string   `yaml:"http_server_port"`
	PortList    []int    `yaml:"check_port"`
	Daemon      bool     `yaml:"Daemon"`
	OutputFile  string   `yaml:"OutputFile"`
}

func ReadYamlConfig(path string) (*Config, error) {
	conf := &Config{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		yaml.NewDecoder(f).Decode(conf)
	}
	return conf, nil
}

func StringToFloat(num string) float64 {
	ret, _ := strconv.ParseFloat(num, 64)
	return ret
}
func FloatToString(n float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(n, 'f', 5, 64)
}

//求交集
func intersect(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

//求差集 slice1-并集
func difference(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	inter := intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}
	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}

func CheckPortLive(p int) string {
	command := fmt.Sprintf("netstat -tnl |awk '{print $4}'|awk -F ':' '{print $NF}'|grep -w %d | wc -l", p)
	cmd := exec.Command("/bin/bash", "-c", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return "0"
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if len(errStr) > 0 {
		fmt.Printf(errStr)
		return "0"
	}
	return strings.TrimSpace(outStr)
}

func networkCollect(detail []net.ConnectionStat, ProcDetail map[string]string) (map[string]map[string]string, map[string]string) {
	everyCon := map[string]map[string]string{}
	var TotalCon = map[string]string{
		"pid":               ProcDetail["pid"],
		"name":              ProcDetail["name"],
		"total_ESTABLISHED": "0",
		"total_TIME_WAIT":   "0",
		"total_CLOSE_WAIT":  "0",
		"total_LISTEN":      "0",
		"total_SYN_RECV":    "0",
		"total_SYN_SENT":    "0",
	}
	for _, x := range detail {
		if x.Status == "LISTEN" {
			everyCon[strconv.Itoa(int(x.Laddr.Port))] = map[string]string{
				"pid":         ProcDetail["pid"],
				"name":        ProcDetail["name"],
				"ESTABLISHED": "0",
				"TIME_WAIT":   "0",
				"CLOSE_WAIT":  "0",
				"LISTEN":      "0",
				"SYN_RECV":    "0",
				"SYN_SENT":    "0",
			}
		}
	}
	for _, item := range detail {
		k := fmt.Sprintf("total_%v", item.Status)
		a, _ := strconv.Atoi(TotalCon[k])
		TotalCon[k] = strconv.Itoa(a + 1)
		if _, ok := everyCon[strconv.Itoa(int(item.Laddr.Port))]; ok {
			b, _ := strconv.Atoi(everyCon[strconv.Itoa(int(item.Laddr.Port))][item.Status])
			everyCon[strconv.Itoa(int(item.Laddr.Port))][item.Status] = strconv.Itoa(b + 1)
		}
	}
	return everyCon, TotalCon
}

func filterProcess(ProcList []string) ([]map[string]string, []map[string]map[string]string) {
	ret := make([]map[string]string, 0)
	netConDetail := make([]map[string]map[string]string, 0)
	var upProcess []string
	p1, _ := process.Processes()
	for _, p := range p1 {
		for _, ProcName := range ProcList {
			if name, _ := p.Cmdline(); strings.Contains(name, ProcName) {
				mem, _ := p.MemoryInfo()
				cpu, _ := p.CPUPercent()
				fd, _ := p.NumFDs()
				thread, _ := p.NumThreads()
				disk, _ := p.IOCounters()
				connect, _ := p.Connections()
				createTime, _ := p.CreateTime()
				singleDetail, singleTotal := networkCollect(connect, map[string]string{"pid": strconv.Itoa(int(p.Pid)), "name": ProcName})
				netConDetail = append(netConDetail, singleDetail)
				ret = append(ret, singleTotal)
				ret = append(ret, map[string]string{
					"rss": strconv.Itoa(int(mem.RSS)), "vms": strconv.Itoa(int(mem.VMS)),
					"mem_locked": strconv.Itoa(int(mem.Locked)), "mem_swap": strconv.Itoa(int(mem.Swap)),
					"mem_stack":   strconv.Itoa(int(mem.Stack)),
					"cpu_percent": FloatToString(cpu),
					"num_fds":     strconv.Itoa(int(fd)), "num_threads": strconv.Itoa(int(thread)),
					"read_bytes": strconv.Itoa(int(disk.ReadBytes)), "write_bytes": strconv.Itoa(int(disk.WriteBytes)),
					"pid": strconv.Itoa(int(p.Pid)), "name": ProcName, "start_time": strconv.FormatInt(createTime, 10), "up": "1"})
				upProcess = append(upProcess, ProcName)
			}
		}
	}
	upDown := difference(ProcList, upProcess)
	for _, d := range upDown {
		ret = append(ret, map[string]string{"pid": "none", "name": d, "up": "0"})
	}

	return ret, netConDetail
}

func NewProcessCollector() *ProcessCollector {
	metrics1 := map[string]ProcessMetrics{}
	for _, metric := range monitor {
		metrics1[metric] = ProcessMetrics{
			desc: prometheus.NewDesc(fmt.Sprintf("%v_%v", namespace, metric), metric, []string{"pid", "name"}, nil),
			extract: func(s string) float64 {
				return StringToFloat(s)
			},
			valType: prometheus.GaugeValue,
		}
	}
	for _, metric3 := range NetMonitor {
		metrics1[metric3] = ProcessMetrics{
			desc: prometheus.NewDesc(fmt.Sprintf("%v_%v", namespace, metric3), metric3, []string{"pid", "name", "port"}, nil),
			extract: func(s string) float64 {
				return StringToFloat(s)
			},
			valType: prometheus.GaugeValue,
		}
	}
	for _, metric4 := range PortMonitor {
		metrics1[metric4] = ProcessMetrics{
			desc: prometheus.NewDesc(fmt.Sprintf("%v_%v", namespace, metric4), metric4, []string{"port"}, nil),
			extract: func(s string) float64 {
				return StringToFloat(s)
			},
			valType: prometheus.GaugeValue,
		}
	}

	return &ProcessCollector{
		metrics: metrics1,
	}
}

func (c *ProcessCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, i := range c.metrics {
		ch <- i.desc
	}
}

func (c *ProcessCollector) Collect(ch chan<- prometheus.Metric) {
	conf, err := ReadYamlConfig(*config)
	if err != nil {
		log.Fatal(err)
	}
	data, netDetail := filterProcess(conf.ProcessList)
	for _, item := range data {
		for k, v := range item {
			if k != "pid" && k != "name" {
				metric, ok := c.metrics[k]
				if ok {
					if metric.extractLabels != nil {
						ch <- prometheus.MustNewConstMetric(metric.desc, metric.valType, metric.extract(v))
					} else {
						ch <- prometheus.MustNewConstMetric(metric.desc, metric.valType, metric.extract(v), item["pid"], item["name"])
					}
				}
			}
		}
	}
	for _, item := range netDetail {
		for port, v1 := range item {
			for metric2, value := range v1 {
				if metric2 != "name" && metric2 != "pid" {
					metric3, ok := c.metrics[metric2]
					if ok {
						if metric3.extractLabels != nil {
							ch <- prometheus.MustNewConstMetric(metric3.desc, metric3.valType, metric3.extract(value))
						} else {
							ch <- prometheus.MustNewConstMetric(metric3.desc, metric3.valType, metric3.extract(value), v1["pid"], v1["name"], port)
						}
					}
				}
			}
		}
	}
	for _, p := range conf.PortList {
		state := CheckPortLive(p)
		metrics4, ok := c.metrics["port_state"]
		if ok {
			ch <- prometheus.MustNewConstMetric(metrics4.desc, metrics4.valType, metrics4.extract(state), strconv.Itoa(p))
		}

	}
}

var (
	logLevel    log.Level = log.InfoLevel
	metricsPath           = flag.String("metrics-path", "/metrics", "path to metrics endpoint")
	rawLevel              = flag.String("log-level", "info", "log level")
	config                = flag.String("config-path", "/Users/asher/data/golang/common/monitor/process/process.yaml", "show version and exit")
)

func init() {
	use, _ := user.Current()
	if use.Name != "root" {
		fmt.Println("当前用户:", use.Name)
		fmt.Println("请切换root用户执行")
		//os.Exit(0)
	}
	flag.Parse()
	parsedLevel, err := log.ParseLevel(*rawLevel)
	if err != nil {
		log.Fatal(err)
	}
	logLevel = parsedLevel
	prometheus.MustRegister(NewProcessCollector())

}

func serveMetrics(bindAddr string) {
	listenAddr := ""
	if strings.Contains(bindAddr, string(':')) {
		listenAddr = bindAddr
	} else {
		listenAddr = fmt.Sprintf("0.0.0.0:%v", bindAddr)
	}
	log.Infof("Starting metric http endpoint on %s", listenAddr)
	http.Handle(*metricsPath, promhttp.Handler())
	http.ListenAndServe(listenAddr, nil)

}

func WriteToFile(OutFile string) {
	prometheus.WriteToTextfile(OutFile, prometheus.DefaultGatherer)
}

func main() {
	log.SetLevel(logLevel)
	log.Info("Starting process_exporter")
	conf, err := ReadYamlConfig(*config)
	if err != nil {
		log.Fatal(err)
	}
	if conf.Daemon {
		serveMetrics(conf.BindPort)
	} else {
		WriteToFile(conf.OutputFile)
	}
}
