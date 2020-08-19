package main

import (
	"errors"
	"fmt"
	"github.com/poemp/go-assign-manage/task-test/service"
	"os"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/util"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateFileOrDir(path string) error{

	b, err := PathExists(path)
	if err != nil {
		fmt.Printf("PathExists(%s),err(%v)\n", path, err)
	}
	if b {
		fmt.Printf("path %s 存在\n", path)
	} else {
		fmt.Printf("path %s 不存在\n", path)
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
			return errors.New(err.Error())
		} else {
			fmt.Printf("mkdir success!\n")
		}
	}
	return nil
}

func main() {
	logDir := "c:/tmp/nacos/log"
	cacheDir := "c:/tmp/nacos/cache"

	c :=CreateFileOrDir(logDir)
	if c != nil {
		panic(c)
	}
	ccc :=CreateFileOrDir(cacheDir)
	if ccc != nil {
		panic(ccc)
	}
	sc := []constant.ServerConfig{
		{
			IpAddr: "console.nacos.io",
			Port:   80,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         "e525eafa-f7d7-4029-83d9-008937f9d468", //namespace id
		TimeoutMs:           5000,
		ListenInterval:      10000,
		NotLoadCacheAtStart: true,
		LogDir:              logDir,
		CacheDir:            cacheDir,
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})

	if err != nil {
		panic(err)
	}

	//Register with default cluster and group
	//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
	service.ExampleserviceclientDeregisterserviceinstance(client, vo.DeregisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        8848,
		ServiceName: "demo.go",
		Ephemeral:   true,
	})

	//Register with cluster name
	//GroupName=DEFAULT_GROUP
	service.ExampleserviceclientRegisterserviceinstance(client, vo.RegisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        8848,
		ServiceName: "demo.go",
		Weight:      10,
		ClusterName: "cluster-a",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})

	//Register different cluster
	//GroupName=DEFAULT_GROUP
	service.ExampleserviceclientRegisterserviceinstance(client, vo.RegisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        8848,
		ServiceName: "demo.go",
		Weight:      10,
		ClusterName: "cluster-b",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})

	//Register different group
	service.ExampleserviceclientRegisterserviceinstance(client, vo.RegisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        8848,
		ServiceName: "demo.go",
		Weight:      10,
		ClusterName: "cluster-b",
		GroupName:   "group-a",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
	service.ExampleserviceclientRegisterserviceinstance(client, vo.RegisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        8848,
		ServiceName: "demo.go",
		Weight:      10,
		ClusterName: "cluster-b",
		GroupName:   "group-b",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})

	//DeRegister with ip,port,serviceName
	//ClusterName=DEFAULT, GroupName=DEFAULT_GROUP
	//Note:ip=127.0.0.1,port=8848 should belong to the cluster of DEFAULT and the group of DEFAULT_GROUP.
	service.ExampleserviceclientDeregisterserviceinstance(client, vo.DeregisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        8848,
		ServiceName: "demo.go",
		Ephemeral:   true, //it must be true
	})

	//DeRegister with ip,port,serviceName,cluster
	//GroupName=DEFAULT_GROUP
	//Note:ip=127.0.0.1,port=8848,cluster=cluster-a should belong to the group of DEFAULT_GROUP.
	service.ExampleserviceclientDeregisterserviceinstance(client, vo.DeregisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        8848,
		ServiceName: "demo.go",
		Cluster:     "cluster-a",
		Ephemeral:   true, //it must be true
	})

	//DeRegister with ip,port,serviceName,cluster,group
	service.ExampleserviceclientDeregisterserviceinstance(client, vo.DeregisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        8848,
		ServiceName: "demo.go",
		Cluster:     "cluster-b",
		GroupName:   "group-b",
		Ephemeral:   true, //it must be true
	})

	//Get service with serviceName
	//ClusterName=DEFAULT, GroupName=DEFAULT_GROUP
	service.ExampleserviceclientGetservice(client, vo.GetServiceParam{
		ServiceName: "demo.go",
	})
	//Get service with serviceName and cluster
	//GroupName=DEFAULT_GROUP
	service.ExampleserviceclientGetservice(client, vo.GetServiceParam{
		ServiceName: "demo.go",
		Clusters:    []string{"cluster-a", "cluster-b"},
	})
	//Get service with serviceName ,group
	//ClusterName=DEFAULT
	service.ExampleserviceclientGetservice(client, vo.GetServiceParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",
	})

	//SelectAllInstance return all instances,include healthy=false,enable=false,weight<=0
	//ClusterName=DEFAULT, GroupName=DEFAULT_GROUP
	service.ExampleserviceclientSelectallinstances(client, vo.SelectAllInstancesParam{
		ServiceName: "demo.go",
	})

	//SelectAllInstance
	//GroupName=DEFAULT_GROUP
	service.ExampleserviceclientSelectallinstances(client, vo.SelectAllInstancesParam{
		ServiceName: "demo.go",
		Clusters:    []string{"cluster-a", "cluster-b"},
	})

	//SelectAllInstance
	//ClusterName=DEFAULT
	service.ExampleserviceclientSelectallinstances(client, vo.SelectAllInstancesParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",
	})

	//SelectInstances only return the instances of healthy=${HealthyOnly},enable=true and weight>0
	//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
	service.ExampleserviceclientSelectinstances(client, vo.SelectInstancesParam{
		ServiceName: "demo.go",
	})

	//SelectOneHealthyInstance return one instance by WRR strategy for load balance
	//And the instance should be health=true,enable=true and weight>0
	//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
	service.ExampleserviceclientSelectonehealthyinstance(client, vo.SelectOneHealthInstanceParam{
		ServiceName: "demo.go",
	})

	//Subscribe key=serviceName+groupName+cluster
	//Note:We call add multiple SubscribeCallback with the same key.
	param := &vo.SubscribeParam{
		ServiceName: "demo.go",
		Clusters:    []string{"cluster-b"},
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			fmt.Printf("callback111 return services:%s \n\n", util.ToJsonString(services))
		},
	}
	service.ExampleserviceclientSubscribe(client, param)
	param2 := &vo.SubscribeParam{
		ServiceName: "demo.go",
		Clusters:    []string{"cluster-b"},
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			fmt.Printf("callback222 return services:%s \n\n", util.ToJsonString(services))
		},
	}
	service.ExampleserviceclientSubscribe(client, param2)
	service.ExampleserviceclientRegisterserviceinstance(client, vo.RegisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        8848,
		ServiceName: "demo.go",
		Weight:      10,
		ClusterName: "cluster-b",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
	//wait for client pull change from server
	time.Sleep(10 * time.Second)

	//Now we just unsubscribe callback1, and callback2 will still receive change event
	service.ExampleserviceclientUnsubscribe(client, param)
	service.ExampleserviceclientDeregisterserviceinstance(client, vo.DeregisterInstanceParam{
		Ip:          "127.0.0.1",
		Ephemeral:   true,
		Port:        8848,
		ServiceName: "demo.go",
		Cluster:     "cluster-b",
	})
	//wait for client pull change from server
	time.Sleep(10 * time.Second)

	//GeAllService will get the list of service name
	//NameSpace default value is public.If the client set the namespaceId, NameSpace will use it.
	//GroupName default value is DEFAULT_GROUP
	service.ExampleserviceclientGetallservice(client, vo.GetAllServiceInfoParam{
		PageNo:   1,
		PageSize: 10,
	})

	service.ExampleserviceclientGetallservice(client, vo.GetAllServiceInfoParam{
		NameSpace: "0e83cc81-9d8c-4bb8-a28a-ff703187543f",
		PageNo:    1,
		PageSize:  10,
	})
}
