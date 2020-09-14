package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/BoynChan/GopherProxy/domain"
	"github.com/BoynChan/GopherProxy/dto"
	"github.com/BoynChan/GopherProxy/pkg"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Author:Boyn
// Date:2020/9/10

func InitServiceRouter(r *gin.Engine) {
	serviceController := r.Group("/service")
	serviceController.GET("/list", serviceList)
	serviceController.GET("/detail/:id", serviceDetail)
	serviceController.POST("", createService)
	serviceController.POST("/:id", updateService)
	serviceController.DELETE("/:id", deleteService)

}

func serviceList(c *gin.Context) {
	var page dto.PageService
	if err := c.ShouldBindQuery(&page); err != nil {
		c.JSON(http.StatusOK, pkg.ErrorMessage(pkg.ParamErrorCode))
		return
	}
	var srv domain.ServiceInfo
	infos, total, err := srv.FindByPage(c, pkg.DefaultDB, page)
	if err != nil {
		c.JSON(http.StatusOK, pkg.ErrorMessage(pkg.DbErrorCode, err.Error()))
		return
	}
	result := domain.ServicePageList{
		ServiceList: infos,
		Total:       total,
		Page:        page.PageIndex,
		Size:        page.PageSize,
	}
	c.JSON(http.StatusOK, pkg.NewMessageBuilder().Data(result).Build())
}

func serviceDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, pkg.ErrorMessage(pkg.ParamErrorCode))
		return
	}
	info := &domain.ServiceInfo{
		Model: gorm.Model{ID: uint(id)},
	}
	detail, err := info.GetServiceDetail(c, pkg.DefaultDB)
	if err != nil {
		c.JSON(http.StatusOK, pkg.NewMessageBuilder().Code(pkg.DbErrorCode).Message(err).Build())
		return
	}
	c.JSON(http.StatusOK, pkg.NewMessageBuilder().Data(detail).Build())
}

func createService(c *gin.Context) {
	serviceType := c.Query("type")
	var input domain.ServiceInputAdapter
	switch serviceType {
	case fmt.Sprintf("%d", domain.HttpLoadType):
		input = &domain.HttpServiceJsonInput{}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusOK, pkg.ErrorMessage(pkg.ParamErrorCode))
			return
		}
	case fmt.Sprintf("%d", domain.GrpcLoadType):
		input = &domain.GrpcServiceJsonInput{}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusOK, pkg.ErrorMessage(pkg.ParamErrorCode))
			return
		}
	}
	tx := pkg.DefaultDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	info, rule := input.ExtractServiceAndRule()
	err := info.Save(c, tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, pkg.NewMessageBuilder().Code(pkg.DbErrorCode).Message(err).Build())
		return
	}
	search := &domain.ServiceInfo{
		ServiceName: info.ServiceName,
	}
	find, err := search.Find(c, tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, pkg.NewMessageBuilder().Code(pkg.DbErrorCode).Message(err).Build())
		return
	}
	rule.FillServiceId(find.ID)
	err = rule.Save(c, tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, pkg.NewMessageBuilder().Code(pkg.DbErrorCode).Message(err).Build())
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, pkg.NewMessageBuilder().Build())
	// TODO(boyn): 服务容器
}

func updateService(c *gin.Context) {
	serviceType := c.Query("type")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, pkg.ErrorMessage(pkg.ParamErrorCode))
		return
	}
	var input domain.ServiceInputAdapter
	switch serviceType {
	case fmt.Sprintf("%d", domain.HttpLoadType):
		input = &domain.HttpServiceJsonInput{}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusOK, pkg.ErrorMessage(pkg.ParamErrorCode))
			return
		}
	case fmt.Sprintf("%d", domain.GrpcLoadType):
		input = &domain.GrpcServiceJsonInput{}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusOK, pkg.ErrorMessage(pkg.ParamErrorCode))
			return
		}
	}
	tx := pkg.DefaultDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	info, rule := input.ExtractServiceAndRule()
	info.ID = uint(id)
	rule.FillServiceId(uint(id))
	err = info.Save(c, tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, pkg.NewMessageBuilder().Code(pkg.DbErrorCode).Message(err).Build())
		return
	}
	search := domain.HttpRule{
		ServiceId: uint(id),
	}
	find, err := search.Find(c, tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, pkg.NewMessageBuilder().Code(pkg.DbErrorCode).Message(err).Build())
		return
	}
	rule.FillId(find.GetId())
	err = rule.Save(c, tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, pkg.NewMessageBuilder().Code(pkg.DbErrorCode).Message(err).Build())
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, pkg.NewMessageBuilder().Build())
}

func deleteService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, pkg.ErrorMessage(pkg.ParamErrorCode))
		return
	}
	info := domain.ServiceInfo{
		Model: gorm.Model{
			ID: uint(id),
		},
	}
	tx := pkg.DefaultDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	find, err := info.Find(c, tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, pkg.NewMessageBuilder().Code(pkg.DbErrorCode).Message(err).Build())
		return
	}
	var rule domain.ServiceRule
	switch find.ServiceType {
	case domain.HttpLoadType:
		rule = &domain.HttpRule{
			ServiceId: find.ID,
		}
	case domain.GrpcLoadType:
		rule = &domain.GrpcRule{
			ServiceId: find.ID,
		}
	}
	err = rule.Delete(c, tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, pkg.NewMessageBuilder().Code(pkg.DbErrorCode).Message(err).Build())
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, pkg.NewMessageBuilder().Build())
}
