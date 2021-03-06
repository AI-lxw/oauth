package middleware

import (
	"oauth/database/bean"

	"github.com/kataras/iris/v12"
)

func AppHaveRole(ctx iris.Context) {
	appID, _ := ctx.Params().GetInt64("appID")
	roleID, _ := ctx.Params().GetInt64("roleID")
	exist, err := bean.Role.AppHaveRole(roleID, appID)
	if err != nil {
		ctx.StatusCode(500)
		ctx.WriteString(err.Error())
		return
	}
	if exist {
		ctx.Next()
	} else {
		ctx.StatusCode(403)
		return
	}
}
