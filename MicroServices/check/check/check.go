package check

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/matscus/Hamster/Package/Services/service"
)

//InitGetResponseAllData - function to obtain information about all services from the database. all services are append to slice Services
func InitGetResponseAllData(project string) (responsealldata *[]service.Service, err error) {
	projects := strings.Split(project, ",")
	services, err := pgClient.GetServicesByProject(projects)
	if err != nil {
		return nil, err
	}
	l := len(*services)
	getResponceAllData := make([]service.Service, l, l)
	for i := 0; i < l; i++ {
		getResponceAllData[i] = service.Service{ID: (*services)[i].ID,
			Name:     (*services)[i].Name,
			Host:     (*services)[i].Host,
			Type:     (*services)[i].Type,
			Port:     (*services)[i].Port,
			Projects: (*services)[i].Projects,
		}
	}
	return &getResponceAllData, nil
}

//CheckStend - function to check stehd, checks the status of monitoring agents, memory utilization, hard disks and processors.
func CheckStend(getResponceAllData *[]service.Service) (res Result, err error) {
	var host string
	var port int
	var id int64
	prometheusstate := true
	prometheusUri := os.Getenv("PROMETHEUSURI")
	temp := make(map[string]Host)
	l := len(*getResponceAllData)
	checkhdd := CheckHDD{}
	checkcpu := CheckCPU{}
	checkmem := CheckMemory{}
	for i := 0; i < l; i++ {
		id = (*getResponceAllData)[i].ID
		host = (*getResponceAllData)[i].Host
		port = (*getResponceAllData)[i].Port
		conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
		if err != nil {
			res.ServiceRS = append(res.ServiceRS, ServerRS{ID: id, Status: false})
			if (*getResponceAllData)[i].Name == "prometheus" {
				prometheusstate = false
			}
		} else {
			res.ServiceRS = append(res.ServiceRS, ServerRS{ID: id, Status: true})
			conn.Close()
		}
	}
	if prometheusstate {
		res.Hosts.PrometheusState = true
		responsefs, err := http.Get(prometheusUri + "?query=node_filesystem_avail_bytes/node_filesystem_size_bytes*100")
		if err != nil {
			err = errors.New("error Get responsefs: %s" + err.Error())
		}
		defer responsefs.Body.Close()
		responsecpu, err := http.Get(prometheusUri + "?query=avg%20by(instance)(max_over_time(node_cpu_seconds_total{mode!=\"idle\"}[5m])-(min_over_time(node_cpu_seconds_total{mode!=\"idle\"}[5m])))")
		if err != nil {
			err = errors.New("error Get responsecpu: %s" + err.Error())
		}
		defer responsecpu.Body.Close()
		responsemem, err := http.Get(prometheusUri + "?query=node_memory_MemAvailable_bytes/node_memory_MemTotal_bytes*100")
		if err != nil {
			err = errors.New("error Get responsemem: %s" + err.Error())
		}
		defer responsemem.Body.Close()
		contentsfs, _ := ioutil.ReadAll(responsefs.Body)
		contentscpu, _ := ioutil.ReadAll(responsecpu.Body)
		contentsmem, _ := ioutil.ReadAll(responsemem.Body)
		err = json.Unmarshal(contentsfs, &checkhdd)
		if err != nil {
			err = errors.New("error unmarshal contentsfs: %s" + err.Error())
		}
		err = json.Unmarshal(contentscpu, &checkcpu)
		if err != nil {
			err = errors.New("error unmarshal contentscpu: %s" + err.Error())
		}
		err = json.Unmarshal(contentsmem, &checkmem)
		if err != nil {
			err = errors.New("error unmarshal contentsmem: %s" + err.Error())
		}
		for i := 0; i < len(checkcpu.Data.Result); i++ {
			tt := fmt.Sprint(checkcpu.Data.Result[i].Value[1])
			v, _ := strconv.ParseFloat(tt, 64)
			if v >= 70 {
				temp[checkcpu.Data.Result[i].Metric.Instance] = Host{
					Host: checkcpu.Data.Result[i].Metric.Instance,
					CPU:  "cpu is used over 70%\n",
				}
			}
		}
		for i := 0; i < len(checkhdd.Data.Result); i++ {
			tt := fmt.Sprint(checkhdd.Data.Result[i].Value[1])
			v, _ := strconv.ParseFloat(tt, 64)
			if v <= 10 {
				if v, ok := temp[checkhdd.Data.Result[i].Metric.Instance]; ok {
					temp[checkhdd.Data.Result[i].Metric.Instance] = Host{
						Host: v.Host,
						CPU:  v.CPU,
						HDD:  v.HDD + "free space in mountpoint " + checkhdd.Data.Result[i].Metric.Mountpoint + " <10%\n",
					}
				} else {
					temp[checkhdd.Data.Result[i].Metric.Instance] = Host{
						Host: checkhdd.Data.Result[i].Metric.Instance,
						HDD:  "free space in mountpoint " + checkhdd.Data.Result[i].Metric.Mountpoint + " <10%\n",
					}
				}
			}
		}
		for i := 0; i < len(checkmem.Data.Result); i++ {
			tt := fmt.Sprint(checkmem.Data.Result[i].Value[1])
			v, _ := strconv.ParseFloat(tt, 64)
			if v <= 30 {
				if v, ok := temp[checkmem.Data.Result[i].Metric.Instance]; ok {
					temp[checkmem.Data.Result[i].Metric.Instance] = Host{
						Host:   v.Host,
						CPU:    v.CPU,
						HDD:    v.HDD,
						Memory: v.Memory + "memory is used over 70%",
					}
				} else {
					temp[checkmem.Data.Result[i].Metric.Instance] = Host{
						Host:   checkmem.Data.Result[i].Metric.Instance,
						Memory: "memory is used over 70%",
					}
				}
			}
		}
		for _, v := range temp {
			res.Hosts.Нost = append(res.Hosts.Нost, v)
		}
	} else {
		res.Hosts.PrometheusState = false
		err = errors.New("Prometheus is not avalible")
	}
	return res, err
}
