package c_user

import (
	"chess/api/components/auth"
	"chess/api/components/input"
	"github.com/gin-gonic/gin"
	"net/http"
	//"chess/api/components/sms"
	"chess/common/config"
	"chess/common/define"
	"chess/models"
)

type PasswordResetParams struct {
	MobileNumber string `json:"mobile_number" form:"mobile_number" binding:"required" description:"手机号"`
	Password     string `json:"password" form:"password" description:"密码"`
	Code         string `json:"code" form:"code" binding:"required" description:"验证码"`
	From         string `json:"from" form:"from" description:"请求来源"`
}

type PasswordResetResult struct {
	define.BaseResult
}

// @Title 重置密码
// @Description 重置密码
// @Summary 重置密码
// @Accept json
// @Param   body     body    c_user.PasswordResetParams  true        "post 数据"
// @Param   token     query    string   true        "token"
// @Param   user_id     path    int   true        "user_id"
// @Success 200 {object} c_user.PasswordResetResult
// @router /user/{user_id}/password/reset [post]
func PasswordReset(c *gin.Context) {
	var result PasswordResetResult
	var form PasswordResetParams

	_conf, ok1 := c.Get("config")
	cConf, ok2 := _conf.(*config.ApiConfig)
	if !ok1 || !ok2 {
		result.Msg = "Get config fail."
		c.JSON(http.StatusOK, result)
		return
	}

	var user = new(models.UsersModel)
	//var err error

	if input.BindJSON(c, &form, cConf) == nil {

		//result.Ret, result.Msg, err = sms.CheckCode(form.MobileNumber, form.Code, sms.SMS_PWD_RESET, cConf)
		//if err != nil {
		//	// 验证不通过
		//	c.JSON(http.StatusOK, result)
		//	return
		//}

		pwdRet := auth.Passwords.CheckPasswordStrong(form.Password)
		if pwdRet != 1 {
			result.Ret = pwdRet
			result.Msg = "password not strong engough"
			c.JSON(http.StatusOK, result)
			return
		}
		// hash password
		hp, err := auth.Passwords.Hash(form.Password)
		if err != nil {
			result.Ret = 0
			result.Msg = "server error"
			c.JSON(http.StatusOK, result)
			return
		}
		// get user
		err = models.Users.GetByMobileNumber(form.MobileNumber, user)
		if err != nil {
			result.Ret = 0
			result.Msg = "server error"
			c.JSON(http.StatusOK, result)
			return
		}
		err = user.UpdatePassword(hp)
		if err != nil {
			result.Ret = 0
			result.Msg = "server error"
			c.JSON(http.StatusOK, result)
			return
		}
		result.Ret = 1
		result.Msg = "ok"
		c.JSON(http.StatusOK, result)
		return

	}
	result.Ret = 0
	result.Msg = "Params invaild."
	c.JSON(http.StatusOK, result)
	return
}
