package controller

import (
	"gatewayDemo/dto"
	"gatewayDemo/middleware"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
)

type DashBoardController struct{}

func DashBoardRegister(group *gin.RouterGroup) {
	dashboard := &DashBoardController{}

	group.GET("panel_group_data", dashboard.PanelGroupData)
}

// PanelGroupData godoc
// @Summary 指标统计
// @Description 指标统计
// @Tags 大盘
// @ID /dashboard/panel_group_data
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.PanelGroupDataOutput} "success"
// @Router /dashboard/panel_group_data [get]
func (dashboard *DashBoardController) PanelGroupData(c *gin.Context) {
	// 连接池
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	out := &dto.PanelGroupDataOutput{
		ServiceNum:      0,
		AppNum:          0,
		CurrentQPS:      0,
		TodayRequestNum: 0,
	}
	middleware.ResponseSuccess(c, out)
}
