package main

import (
	"oauth/config"
	"oauth/controller"
	"oauth/controller/middleware"
	_ "oauth/database"

	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
)

var (
	cookieNameForSessionID = "mgtv-oauth-sessionid"
	session                = sessions.New(sessions.Config{
		Cookie: cookieNameForSessionID,
		//Expires: 45 * time.Minute, // <=0 means unlimited life
	})
)

func GetApp() *iris.Application {
	app := iris.New()
	tmpl := iris.HTML("./static/template", ".html")
	tmpl.Reload(true)

	app.RegisterView(tmpl)

	app.StaticWeb("/static/", "./static/resource")

	//免登陆接口
	webIndexCtrl := controller.WebIndex{Session: session}
	app.Get("/", webIndexCtrl.Get)
	userCtrl := controller.User{Session: session}
	app.PartyFunc("/user", func(u iris.Party) {
		u.Get("/register", func(ctx iris.Context) { ctx.View("register.html") })
		//注册
		u.Post("/register", userCtrl.Post)
		//退出
		u.Delete("/logout", userCtrl.Logout)
		//提交登陆表单
		u.Post("/login", userCtrl.Login)
		//修改密码
		u.Get("/password", func(ctx iris.Context) { ctx.View("password.html") })
		u.Get("/password/{uid:long}", userCtrl.ResetPassword4AdminView)
	})

	appCtrl := controller.App{Session: session}
	app.PartyFunc("/app", func(u iris.Party) {
		u.Get("/register", func(ctx iris.Context) {
			ctx.View("app-register.html")
		})
		//注册
		u.Post("/register", appCtrl.Post)

		u.Get("/{appID:long}", appCtrl.EditPage)
	})

	//用户管理
	middle := middleware.MiddleWare{Session: session}
	appUserMangerCtrl := controller.AppUserManager{Session: session}
	app.PartyFunc("/app/{appID:long}/user_manager", func(u iris.Party) {
		//判定app是否归属于当前用户(该应用是否为该用户创建)
		u.Use(middle.UserHaveApp)
		u.Get("/", appUserMangerCtrl.GetView)
		u.Get("/role", appUserMangerCtrl.GetRoleView)
		u.Get("/role/{roleID:long}/permission", appUserMangerCtrl.GetRolePermissionView)
	})

	//以下为第三方调用接口
	authorizeCtrl := controller.Authorize{Session: session}
	app.PartyFunc("/authorize", func(u iris.Party) {
		u.Get("/", authorizeCtrl.Get)
		//权限校验
		u.Post("/", authorizeCtrl.Verify)
		//接口跳转
		u.Post("/jump", authorizeCtrl.Jump)
	})
	resourceCtrl := controller.Resource{Session: session}
	app.PartyFunc("/resource", func(u iris.Party) {
		u.Post("/account", resourceCtrl.GetAccount)
	})

	//数据接口

	API := app.Party("/api")
	API.PartyFunc("/user", func(u iris.Party) {
		u.Use(middle.UserAuth)
		u.Put("/password", userCtrl.ResetPassword)
		u.Put("/password/{uid:long}", userCtrl.ResetPassword4Admin)
		u.Delete("/{uid:long}", userCtrl.DeleteUser)
	})

	application := API.Party("/app/{appID:long}", middle.UserHaveApp)
	application.PartyFunc("/", func(u iris.Party) {
		u.Delete("/", appCtrl.Delete)
		u.Put("/", appCtrl.Put)
		u.Patch("/user_mode/{mode:string}", appCtrl.UpdateRunMode)
	})

	//黑白名单
	application.PartyFunc("/user_manager", func(u iris.Party) {
		u.Post("/", appUserMangerCtrl.Post)
		u.Delete("/{id: long}", appUserMangerCtrl.Delete)
	})

	//应用角色
	roleCtrl := controller.AppRoleManager{Session: session}
	Role := application.PartyFunc("/role", func(u iris.Party) {
		u.Post("/", roleCtrl.Post)
		u.Get("/", roleCtrl.Get)
		u.Delete("/{id:long}", roleCtrl.Delete)
	})

	permissionCtrl := controller.AppRolePermission{Session: session}
	Role.PartyFunc("/{roleID:long}/permission", func(u iris.Party) {
		u.Post("/", permissionCtrl.Post)
		u.Delete("/{id:long}", permissionCtrl.Delete)
	})

	return app
}

func main() {
	app := GetApp()
	app.Run(iris.Addr(":"+config.Get().Port), iris.WithoutVersionChecker)
}
