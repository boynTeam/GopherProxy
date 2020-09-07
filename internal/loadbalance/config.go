package loadbalance

import (
	"fmt"
	"strings"

	"github.com/BoynChan/GopherProxy/pkg"
	"github.com/sirupsen/logrus"
)

// Author:Boyn
// Date:2020/9/7

type Config interface {
	Attach(o Observer)
	GetConf() ([]ConfigValue, error)
	WatchConf()
	UpdateConf(conf []string)
}

type ZkConfig struct {
	observers  []Observer
	path       string
	zkHosts    []string
	activeList []string
	format     string
}

type Observer interface {
	Update([]ConfigValue) error
}

type ConfigValue struct {
	Key   string
	Value string
}

func (c *ZkConfig) Attach(o Observer) {
	c.observers = append(c.observers, o)
}

func (c *ZkConfig) GetConf() ([]ConfigValue, error) {
	zkManager := pkg.NewZkManager(c.zkHosts...)
	zkManager.GetConnect()
	defer zkManager.Close()
	configList := make([]ConfigValue, 0)
	for _, ip := range c.activeList {
		fullPath := fmt.Sprintf("%s/%s", c.path, fmt.Sprintf(c.format, ip))
		data, _, err := zkManager.GetPathData(fullPath)
		if err != nil {
			return nil, err
		}
		configList = append(configList, ConfigValue{
			Key:   fmt.Sprintf(c.format, ip),
			Value: string(data),
		})
	}
	return configList, nil
}

func (c *ZkConfig) WatchConf() {
	zkManager := pkg.NewZkManager(c.zkHosts...)
	zkManager.GetConnect()
	chanList, chanErr := zkManager.WatchServerListByPath(c.path)
	go func() {
		defer zkManager.Close()
		for {
			select {
			case changeErr := <-chanErr:
				logrus.Errorf("ChangeError %v", changeErr)
			case changeList := <-chanList:
				c.UpdateConf(changeList)
			}
		}
	}()
}

func (c *ZkConfig) UpdateConf(conf []string) {
	c.activeList = conf
	values, err := c.GetConf()
	if err != nil {
		logrus.Errorf("Update Conf -> get Conf failed %v", err)
		return
	}
	logrus.Infof("Update Service Config %+v", values)
	for _, obs := range c.observers {
		obs.Update(values)
	}
}

func NewZkConf(format, path string, zkHosts []string) (Config, error) {
	zkManager := pkg.NewZkManager(zkHosts...)
	err := zkManager.GetConnect()
	if err != nil {
		return nil, err
	}
	defer zkManager.Close()
	pathList := strings.Split(path[1:], "/")
	var pathBuffer string
	for _, existPath := range pathList {
		pathBuffer = fmt.Sprintf("%s/%s", pathBuffer, existPath)
		ex, err := zkManager.NodeExist(pathBuffer)
		if err != nil {
			logrus.Errorf("NewZkConf: Check path exists error %v", err)
			break
		}
		if !ex {
			err := zkManager.RegisterServerNode(pathBuffer, "")
			if err != nil {
				logrus.Errorf("NewZkConf: Register path %s error %v", pathBuffer, err)
			}
		}
	}
	zList, err := zkManager.GetServerListByPath(path)
	if err != nil {
		return nil, err
	}
	mConf := &ZkConfig{
		path:       path,
		zkHosts:    zkHosts,
		activeList: zList,
		format:     format,
	}
	mConf.WatchConf()
	return mConf, nil
}
