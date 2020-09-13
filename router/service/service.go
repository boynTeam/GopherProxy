package service

import (
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
