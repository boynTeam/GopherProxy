package domain

import (
	"github.com/gin-gonic/gin"
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
	ServiceDesc string `gorm:"column:service_desc" json:"service_desc"`
	RoundType   int    `json:"round_type" gorm:"column:round_type"`
}

func (t *ServiceInfo) TableName() string {
	return "gateway_service_info"
}

type ServiceRule interface {
	Find(c *gin.Context, db *gorm.DB, search ServiceRule) (ServiceRule, error)
	Save(c *gin.Context, db *gorm.DB) error
}

type HttpRule struct {
	gorm.Model
	ServiceId      uint   `json:"service_id" gorm:"column:service_id"`
	Prefix         string `json:"prefix" gorm:"column:prefix"` // URL前缀
	NeedHttps      int    `json:"need_https" gorm:"column:need_https"`
	NeedWebsocket  int    `json:"need_websocket" gorm:"column:need_websocket"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor"`
}

func (h *HttpRule) Find(c *gin.Context, db *gorm.DB, search ServiceRule) (ServiceRule, error) {
	var httpRule *HttpRule
	err := db.Where(search).Find(&httpRule).Error
	return httpRule, err
}

func (h *HttpRule) Save(c *gin.Context, db *gorm.DB) error {
	return db.Save(h).Error
}

func (h *HttpRule) TableName() string {
	return "gateway_service_http_rule"
}

type GrpcRule struct {
	gorm.Model
	ServiceId      uint   `json:"service_id" gorm:"column:service_id"`
	Port           int    `json:"port" gorm:"column:port"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor"`
}

func (g *GrpcRule) TableName() string {
	return "gateway_service_grpc_rule"
}
