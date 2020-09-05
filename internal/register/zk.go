package register

// Author:Boyn
// Date:2020/9/4

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/sirupsen/logrus"
)

type ZkManager struct {
	hosts      []string
	conn       *zk.Conn
	pathPrefix string
}

func NewZkManager(hosts ...string) *ZkManager {
	return &ZkManager{hosts: hosts, pathPrefix: "/gateway_servers_"}
}

//连接zk服务器
func (z *ZkManager) GetConnect() error {
	conn, _, err := zk.Connect(z.hosts, 5*time.Second)
	if err != nil {
		return err
	}
	z.conn = conn
	return nil
}

//关闭服务
func (z *ZkManager) Close() {
	z.conn.Close()
	return
}

//获取配置
func (z *ZkManager) GetPathData(nodePath string) ([]byte, *zk.Stat, error) {
	return z.conn.Get(nodePath)
}

//更新配置
func (z *ZkManager) SetPathData(nodePath string, config []byte) (err error) {
	ex, _ := z.NodeExist(nodePath)
	if !ex {
		return errors.New("node does not exists")
	}
	_, dStat, err := z.GetPathData(nodePath)
	if err != nil {
		return
	}
	_, err = z.conn.Set(nodePath, config, dStat.Version)
	if err != nil {
		logrus.Errorf("Update node error %v", err)
		return err
	}
	return
}

func (z *ZkManager) RegisterServerNode(prefix, nodeName string, data ...byte) error {
	return z.doRegister(prefix, nodeName, false, data...)
}

//创建临时节点
func (z *ZkManager) RegistServerTmpNode(prefix, nodeName string, data ...byte) (err error) {
	return z.doRegister(prefix, nodeName, true, data...)
}

func (z *ZkManager) doRegister(prefix string, nodeName string, isTemp bool, data ...byte) error {
	var flag int32
	if isTemp {
		flag = zk.FlagEphemeral
	}
	ex, err := z.NodeExist(prefix)
	if err != nil {
		logrus.Errorf("Exists error %v", prefix)
		return err
	}
	if !ex {
		_, err = z.conn.Create(prefix, data, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			logrus.Errorf("Create error %v %s", err, prefix)
			return err
		}
	}
	subNodePath := fmt.Sprintf("%s/%s", prefix, nodeName)
	ex, err = z.NodeExist(subNodePath)
	if err != nil {
		logrus.Errorf("Exists error %v", subNodePath)
		return err
	}
	if !ex {
		_, err = z.conn.Create(subNodePath, data, flag, zk.WorldACL(zk.PermAll))
		if err != nil {
			logrus.Errorf("Create error %v", subNodePath)
			return err
		}
	}
	return nil
}

func (z *ZkManager) NodeExist(path string) (bool, error) {
	exist, _, err := z.conn.Exists(path)
	return exist, err
}

//获取服务列表
func (z *ZkManager) GetServerListByPath(path string) (list []string, err error) {
	list, _, err = z.conn.Children(path)
	return
}

//watch机制，服务器有断开或者重连，收到消息
func (z *ZkManager) WatchServerListByPath(path string) (chan []string, chan error) {
	conn := z.conn
	snapshots := make(chan []string)
	errors := make(chan error)
	go func() {
		for {
			snapshot, _, events, err := conn.ChildrenW(path)
			if err != nil {
				errors <- err
			}
			snapshots <- snapshot
			select {
			case evt := <-events:
				if evt.Err != nil {
					errors <- evt.Err
				}
				logrus.Infof("ChildrenW Event Path:%v, Type:%v", evt.Path, evt.Type)
			}
		}
	}()

	return snapshots, errors
}

//watch机制，监听节点值变化
func (z *ZkManager) WatchPathData(nodePath string) (chan []byte, chan error) {
	conn := z.conn
	snapshots := make(chan []byte)
	errors := make(chan error)

	go func() {
		for {
			dataBuf, _, events, err := conn.GetW(nodePath)
			if err != nil {
				errors <- err
				return
			}
			snapshots <- dataBuf
			select {
			case evt := <-events:
				if evt.Err != nil {
					errors <- evt.Err
					return
				}
				logrus.Errorf("GetW Event Path:%v, Type:%v", evt.Path, evt.Type)
			}
		}
	}()
	return snapshots, errors
}
