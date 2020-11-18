package controller

import (
	"encoding/json"
	"gatewayDemo/dao"
	"gatewayDemo/dto"
	"gatewayDemo/middleware"
	"gatewayDemo/public"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
}

func AdminRegister(group *gin.RouterGroup) {
	admin := &AdminController{}
	group.GET("/info", admin.AdminInfo)
	group.POST("/changepwd", admin.ChangePWD)
}

// AdminInfo godoc
// @Summary 管理员信息
// @Description 管理员信息
// @Tags 管理员接口
// @ID /admin_info/info
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router /admin_info/info [get]
func (admininfo *AdminController) AdminInfo(c *gin.Context) {
	// 1. 读取sessionKey对应json，转换为结构体
	// 2. 取出数据然后封装输出结构体

	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(sessInfo.(string)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	out := &dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		Name:         adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avator:       "",
		Introduction: "I'm a super administrator!",
		Roles:        []string{"admin"},
	}
	middleware.ResponseSuccess(c, out)
}

// ChangePWD godoc
// @Summary 管理员密码修改
// @Description 管理员密码修改
// @Tags 管理员接口
// @ID /admin_info/changepwd
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin_info/info [post]
func (admininfo *AdminController) ChangePWD(c *gin.Context) {
	params := &dto.ChangePWDInput{}
	if err := params.BindingValidParams(c); err != nil {
		// log.F  atal("params.BindingValidParams err : %v", err)
		middleware.ResponseError(c, 2000, err)
		return
	}

	// 1. session读取用户信息到结构体 sessionInfo
	// 2. sessionInfo.ID 读取数据库信息 adminInfo
	// 3. originpassword + adminInfo.salt sha256 saltoriginpassword
	// 4. saltoriginpassword == adminInfo.password
	// 5. params.password + adminInfo.salt sha256 saltasswrod
	// 6. saltPassword ==> adminInfo.password 执行数据保存
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(sessInfo.(string)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	// 从数据库中读取adminInfo
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	adminInfo := &dao.AdminInfo{}
	adminInfo.Find(
		c,
		tx,
		(&dao.AdminInfo{UserName: adminSessionInfo.UserName}))

	middleware.ResponseSuccess(c, "")
}
