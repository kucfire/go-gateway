package controller

import (
	"encoding/json"
	"gatewayDemo/dao"
	"gatewayDemo/dto"
	"gatewayDemo/middleware"
	"gatewayDemo/public"
	"log"
	"time"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AdminLoginController struct{}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin)

}

// AdminLogin godoc
// @Summary 管理员登录
// @Description 管理员登录
// @Tags 管理员接口
// @ID /admin_login/login
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginInput} "success"
// @Router /admin_login/login [post]
func (adminligin *AdminLoginController) AdminLogin(c *gin.Context) {
	params := &dto.AdminLoginInput{}
	if err := params.BindingValidParams(c); err != nil {
		// log.F  atal("params.BindingValidParams err : %v", err)
		middleware.ResponseError(c, 2000, err)
		return
	}

	// 1. 从数据库中取得管理员信息 adminInfo
	// 2. 使用adminInfo + params.Password 进行sha256加密 => 得到saltPassword
	// 3. saltPassword == admininfo.password
	admin := &dao.AdminInfo{}
	// set mysql db connect
	tx, err := lib.GetGormPool("default")
	log.Println(tx)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// // 校验登录信息(密码)
	// admin, err = admin.LoginCheck(c, tx, params)
	// if err != nil {
	// 	middleware.ResponseError(c, 2002, err)
	// 	return
	// }

	// 设置session
	adminSessionInfo := &dto.AdminSessionInfo{
		ID:        admin.Id,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}

	//

	sessBts, errJson := json.Marshal(adminSessionInfo)
	if errJson != nil {
		middleware.ResponseError(c, 2003, errJson)
		return
	}
	sess := sessions.Default(c)
	sess.Set(public.AdminSessionInfoKey,
		string(sessBts))
	sess.Save()

	out := &dto.AdminLoginOutput{Token: admin.UserName}
	middleware.ResponseSuccess(c, out)

}
