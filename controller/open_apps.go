package controller

import (
	"oauth/database/bean"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

func isAdmin(uid int64) bool {
	return uid == 0
}

type OpenApps struct {}

func (c *OpenApps) GetList(ctx iris.Context) {
	sess := sessions.Get(ctx)
	uid, _ := sess.GetInt64("uid")
	list, err := bean.GetOpenApplicationList()

	if err != nil {
		ctx.StatusCode(500)
		ctx.WriteString(err.Error())
		return
	}
	if isAdmin(uid) {
		ctx.JSON(list)
		return
	}

}

var OpenAppsCtrl OpenApps