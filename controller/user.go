package controller

import (
	"log"
	"strings"

	"oauth/config"
	"oauth/database/bean"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

type User struct {}


// 是否开放用户注册
//用户注册
func (c *User) IsOpenRegister(ctx iris.Context) {
	ctx.JSON(map[string]bool{
		"open": config.Get().OpenRegister,
	})
}

//用户注册
func (c *User) Post(ctx iris.Context) {
	sess := sessions.Get(ctx)
	//如果没有开放用户认证，用户不是管理员，那么就拒绝注册
	if !config.Get().OpenRegister && sess.GetString("username") != config.Get().Account.User {
		ctx.StatusCode(403)
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


type PassResetForm struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (c *User) GetList(ctx iris.Context) {
	sess := sessions.Get(ctx)
	// 用户是否已登陆
	uid, _ := sess.GetInt64("uid")

	if uid == 0 {
		if userList, err := bean.GetAllUser(); err != nil {
			log.Println(err)
			ctx.StatusCode(500)
		} else {
			ctx.JSON(userList)
		}
		return
	}
	ctx.JSON([]interface{}{})
}

func (c *User) GetAnyOneInfo(ctx iris.Context) {
	sess := sessions.Get(ctx)
	cid, _ := sess.GetInt64("uid")
	if cid != 0 {
		ctx.StatusCode(403)
		return
	}
	uid, _ := ctx.Params().GetInt64("id")

	info, err := bean.GetUserByID(uid)
	if err != nil {
		ctx.StatusCode(406)
		ctx.WriteString(err.Error())
		return
	}
	ctx.JSON(info)
}

func (c *User) ResetPassword(ctx iris.Context) {
	sess := sessions.Get(ctx)
	uid, _ := sess.GetInt64("uid")
	form := PassResetForm{}
	ctx.ReadJSON(&form)
	if uid == 0 {
		uid, _ = sess.GetInt64("adminID")
	}
	oldPassword := strings.TrimSpace(form.OldPassword)
	newPassword := strings.TrimSpace(form.NewPassword)

	if oldPassword == "" || newPassword == "" {
		ctx.StatusCode(406)
		ctx.WriteString("密码不能为空")
		return
	}

	err := bean.UpdateUserPassword(uid, oldPassword, newPassword)
	if err != nil {
		ctx.StatusCode(406)
		ctx.WriteString(err.Error())
		return
	}
	sess.Clear()
	sess.ClearFlashes()
	ctx.StatusCode(200)
}

func (c *User) ResetPassword4Admin(ctx iris.Context) {
	sess := sessions.Get(ctx)
	currendUID, _ := sess.GetInt64("uid")
	//非管理员
	if currendUID != 0 {
		ctx.StatusCode(403)
		return
	}
	uid, _ := ctx.Params().GetInt64("uid")

	form := map[string]string{}
	ctx.ReadJSON(&form)
	newPassword, ok := form["newPassword"]
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

func (c *User) DeleteUser(ctx iris.Context) {
	sess := sessions.Get(ctx)
	uid, _ := ctx.Params().GetInt64("uid")
	if currendUID, err := sess.GetInt64("adminID"); err != nil {
		ctx.StatusCode(403)
		return
	} else if currendUID == uid {
		ctx.StatusCode(406)
		ctx.WriteString("不允许删除管理员")
		return
	}

	if err := bean.DeleteUser(uid); err != nil {
		ctx.StatusCode(500)
		ctx.WriteString(err.Error())
		return
	}
	ctx.StatusCode(200)
}

var UserCtrl User