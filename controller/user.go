package controller

import (
	"oauth/config"
	"oauth/database/bean"
	"strings"

	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
)

type User struct {
	Session *sessions.Sessions
}

//用户注册
func (c *User) Post(ctx iris.Context) {
	sess := c.Session.Start(ctx)
	//如果没有开放用户认证，用户不是管理员，那么就拒绝注册
	if !config.Get().OpenRegister && sess.GetString("username") != config.Get().Account.User {
		ctx.StatusCode(401)
		return
	}

	account := accountForm{}
	ctx.ReadJSON(&account)
	username := strings.TrimSpace(account.Username)
	password := strings.TrimSpace(account.Password)
	if username == "" || password == "" {
		ctx.StatusCode(iris.StatusNotAcceptable)
		ctx.WriteString("用户名或密码不能为空")
		return
	}
	if err := bean.RegisterUser(username, password); err != nil {
		ctx.StatusCode(iris.StatusNotAcceptable)
		ctx.WriteString(err.Error())
	} else {
		ctx.StatusCode(200)
		ctx.WriteString("注册成功！")
	}
}

func (c *User) Logout(ctx iris.Context) {
	sess := c.Session.Start(ctx)
	sess.Clear()
	sess.ClearFlashes()
	ctx.StatusCode(200)
}

type accountForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *User) Login(ctx iris.Context) {
	account := accountForm{}
	ctx.ReadJSON(&account)
	username := strings.TrimSpace(account.Username)
	password := strings.TrimSpace(account.Password)
	if username == "" || password == "" {
		ctx.StatusCode(403)
		return
	}

	if u, exist, err := bean.FindUser(username, password); err != nil {
		ctx.StatusCode(500)
		ctx.WriteString(err.Error())
	} else if !exist {
		ctx.StatusCode(403)
	} else {
		sess := c.Session.Start(ctx)
		sess.Set("user-authorized", true)
		sess.Set("username", username)
		//管理员
		if username == config.Get().Account.User {
			sess.Set("uid", -1)
		} else {
			sess.Set("uid", u.ID)
		}
		ctx.StatusCode(200)
	}
}

type PassResetForm struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (c *User) ResetPassword(ctx iris.Context) {
	sess := c.Session.Start(ctx)
	uid, _ := sess.GetInt64("uid")
	form := PassResetForm{}
	ctx.ReadJSON(&form)
	err := bean.UpdateUserPassword(uid, form.OldPassword, form.NewPassword)
	if err != nil {
		ctx.StatusCode(406)
		ctx.WriteString(err.Error())
		return
	}
	sess.Clear()
	sess.ClearFlashes()
	ctx.StatusCode(200)
}

func (c *User) ResetPassword4AdminView(ctx iris.Context) {
	sess := c.Session.Start(ctx)
	currendUID, _ := sess.GetInt64("uid")
	if currendUID != -1 {
		ctx.StatusCode(403)
		return
	}
	uid, _ := ctx.Params().GetInt64("uid")

	user, err := bean.GetUserByID(uid)
	if err != nil {
		ctx.StatusCode(500)
		ctx.WriteString(err.Error())
		return
	}
	if user.ID == 0 {
		ctx.StatusCode(404)
		return
	}

	ctx.ViewData("User", user)
	ctx.View("password4admin.html")

}

func (c *User) ResetPassword4Admin(ctx iris.Context) {
	sess := c.Session.Start(ctx)
	currendUID, _ := sess.GetInt64("uid")
	if currendUID != -1 {
		ctx.StatusCode(403)
		return
	}
	uid, _ := ctx.Params().GetInt64("uid")

	form := map[string]string{}
	ctx.ReadJSON(&form)
	newPassword, ok := form["password"]
	if !ok || newPassword == "" {
		ctx.StatusCode(406)
		return
	}
	err := bean.UpdateUserPasswordNoOld(uid, newPassword)
	if err != nil {
		ctx.StatusCode(406)
		ctx.WriteString(err.Error())
		return
	}
	ctx.StatusCode(200)
}
