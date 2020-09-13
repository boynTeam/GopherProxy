package domain

import (
	"errors"
	"fmt"

	"github.com/BoynChan/GopherProxy/dto"
	"github.com/BoynChan/GopherProxy/pkg"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// Author:Boyn
// Date:2020/9/10

const (
	HttpLoadType = iota + 1
	GrpcLoadType
)

const (
	HttpPrefixRuleType = iota + 1 //
	HttpDomainRuleType
)

const (
	RandomRound = iota + 1
	RobinRound
	ConsistencyHash
)

type ServiceInfo struct {
	gorm.Model
	ServiceType int    `gorm:"column:service_type" json:"service_type"`
	ServiceName string `gorm:"column:service_name" json:"service_name"`
	ServicePort int    `gorm:"column:service_port" json:"service_port"`
	ServiceDesc string `gorm:"column:service_desc" json:"service_desc"`
	RoundType   int    `json:"round_type" gorm:"column:round_type"`
}

func (t *ServiceInfo) TableName() string {
	return "gateway_service_info"
}

func (t *ServiceInfo) Save(c *gin.Context, db *gorm.DB) error {
	return db.Save(t).Error
}

func (t *ServiceInfo) Find(c *gin.Context, db *gorm.DB) (*ServiceInfo, error) {
	var serviceInfo ServiceInfo
	err := db.Where(t).First(&serviceInfo).Error
	return &serviceInfo, err
}

func (t *ServiceInfo) Delete(c *gin.Context, db *gorm.DB) error {
	return db.Unscoped().Where(t).Delete(t).Error
}

func (t *ServiceInfo) FindByPage(c *gin.Context, db *gorm.DB, page dto.PageService) (infos []ServiceInfo, total int64, err error) {
	err = db.Model(&ServiceInfo{}).Find(&infos).Limit((page.PageIndex - 1) * page.PageSize).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Model(&ServiceInfo{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return
}

func (t *ServiceInfo) GetServiceDetail(c *gin.Context, db *gorm.DB) (*ServiceDetail, error) {
	zkAddr := viper.GetString("Zk.Addr")

	baseInfo, err := t.Find(c, db)
	if err != nil {
		return nil, err
	}

	detail := &ServiceDetail{
		ServiceId:      baseInfo.ID,
		ServiceName:    baseInfo.ServiceName,
		ServiceAddress: fmt.Sprintf("127.0.0.1:%d", baseInfo.ServicePort),
		ServiceDesc:    baseInfo.ServiceDesc,
		RoundType:      baseInfo.RoundType,
	}
	var condition ServiceRule
	var zkRegisterPath string
	switch baseInfo.ServiceType {
	case HttpLoadType:
		condition = &HttpRule{
			ServiceId: baseInfo.ID,
		}
		zkRegisterPath = fmt.Sprintf("%s/%s", viper.GetString("Zk.HttpPrefix"), baseInfo.ServiceName)
	case GrpcLoadType:
		condition = &GrpcRule{
			ServiceId: baseInfo.ID,
		}
		zkRegisterPath = fmt.Sprintf("%s/%s", viper.GetString("Zk.GrpcPrefix"), baseInfo.ServiceName)
	default:
		return nil, errors.New("load type not support")
	}

	serviceRule, err := condition.Find(c, db)
	if err != nil {
		return nil, err
	}
	detail.Rule = serviceRule
	zkManager := pkg.NewZkManager([]string{zkAddr}...)
	err = zkManager.GetConnect()
	if err != nil {
		return nil, err
	}
	nodes, err := zkManager.GetServerListByPath(zkRegisterPath)
	if err != nil {
		return nil, err
	}
	detail.ServerList = nodes
	return detail, nil
}

type ServiceDetail struct {
	ServiceId      uint        `json:"service_id"`      // 服务ID
	ServiceName    string      `json:"service_name"`    // 服务名
	ServiceAddress string      `json:"service_address"` // 服务地址
	ServiceDesc    string      `json:"service_desc"`    // 服务描述
	RoundType      int         `json:"round_type"`      // 轮询方式
	Rule           ServiceRule `json:"rule"`            // 服务配置
	ServerList     []string    `json:"server_list"`     // 服务下游服务器地址
}
