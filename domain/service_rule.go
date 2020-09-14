package domain

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Author:Boyn
// Date:2020/9/12

type ServiceRule interface {
	Find(c *gin.Context, db *gorm.DB) (ServiceRule, error)
	Save(c *gin.Context, db *gorm.DB) error
	Update(c *gin.Context, db *gorm.DB) error
	Delete(c *gin.Context, db *gorm.DB) error
	FillServiceId(id uint)
	FillId(id uint)
	GetId() uint
}

type HttpRule struct {
	gorm.Model
	ServiceId      uint   `json:"service_id" gorm:"column:service_id"`
	Prefix         string `json:"prefix" gorm:"column:prefix"` // URL前缀
	NeedHttps      int    `json:"need_https" gorm:"column:need_https"`
	NeedWebsocket  int    `json:"need_websocket" gorm:"column:need_websocket"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor"`
}

func (h *HttpRule) Find(c *gin.Context, db *gorm.DB) (ServiceRule, error) {
	var httpRule HttpRule
	err := db.Where(h).First(&httpRule).Error
	return &httpRule, err
}

func (h *HttpRule) Save(c *gin.Context, db *gorm.DB) error {
	return db.Save(h).Error
}

func (h *HttpRule) Update(c *gin.Context, db *gorm.DB) error {
	return db.Save(h).Error
}

func (h *HttpRule) Delete(c *gin.Context, db *gorm.DB) error {
	return db.Where(h).Delete(h).Error
}

func (h *HttpRule) FillServiceId(id uint) {
	h.ServiceId = id
}

func (h *HttpRule) FillId(id uint) {
	h.ID = id
}

func (h *HttpRule) TableName() string {
	return "gateway_service_http_rule"
}

func (h *HttpRule) GetId() uint {
	return h.ID
}

type GrpcRule struct {
	gorm.Model
	ServiceId      uint   `json:"service_id" gorm:"column:service_id"`
	Port           int    `json:"port" gorm:"column:port"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor"`
}

func (g *GrpcRule) Find(c *gin.Context, db *gorm.DB) (ServiceRule, error) {
	var grpcRule GrpcRule
	err := db.Where(g).First(&grpcRule).Error
	return &grpcRule, err
}

func (g *GrpcRule) Save(c *gin.Context, db *gorm.DB) error {
	return db.Save(g).Error
}

func (g *GrpcRule) Update(c *gin.Context, db *gorm.DB) error {
	return db.Save(g).Error
}

func (g *GrpcRule) Delete(c *gin.Context, db *gorm.DB) error {
	return db.Where(g).Delete(g).Error
}

func (g *GrpcRule) FillServiceId(id uint) {
	g.ServiceId = id
}

func (g *GrpcRule) FillId(id uint) {
	g.ID = id
}

func (g *GrpcRule) GetId() uint {
	return g.ID
}

func (g *GrpcRule) TableName() string {
	return "gateway_service_grpc_rule"
}
