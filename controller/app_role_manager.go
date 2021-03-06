package controller

import (
	"oauth/database/bean"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

type AppRoleManager struct {
	Session *sessions.Sessions
}

func (a *AppRoleManager) Get(ctx iris.Context) {
	appID, _ := ctx.Params().GetInt64("appID")
	list, err := bean.Role.GetRoleList(appID)
	if err != nil {
		ctx.StatusCode(500)
		ctx.WriteString(err.Error())
		return
	}
	ctx.JSON(list)
}

func (a *AppRoleManager) Post(ctx iris.Context) {
	appID, _ := ctx.Params().GetInt64("appID")
	form := map[string]string{}

	ctx.ReadJSON(&form)

	role, isExistrole := form["role"]

	if !isExistrole {
		ctx.StatusCode(406)
		return
	}

	role = strings.TrimSpace(role)
	if role == "" {
		ctx.StatusCode(406)
		return
	}

	if err := bean.Role.Add(role, appID); err != nil {
		if err != nil {
			ctx.StatusCode(500)
			ctx.WriteString(err.Error())
		}
	}
}

func (a *AppRoleManager) Delete(ctx iris.Context) {
	id, _ := ctx.Params().GetInt64("id")
	if err := bean.Role.Delete(id); err != nil {
		ctx.StatusCode(406)
		ctx.WriteString(err.Error())
	}
}

func (a *AppRoleManager) GetOne(ctx iris.Context) {
	id, _ := ctx.Params().GetInt64("id")
	appID, _ := ctx.Params().GetInt64("appID")
	role, err := bean.Role.Get(id)
	if err != nil {
		ctx.StatusCode(500)
		ctx.WriteString(err.Error())
		return
	}
	if role.AppID != appID {
		ctx.StatusCode(403)
		ctx.WriteString("无权限")
		return
	}
	ctx.JSON(role)
}
