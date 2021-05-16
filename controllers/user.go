package controllers

import (
	"log"
	"time"

	"courseScheduling/models"

	beego "github.com/beego/beego/v2/server/web"
)

type UserController struct {
	beego.Controller
}

// @Title Login
// @Description Logs user into the system
// @Param	username		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /login [get]
func (this *UserController) Login() {
	username := this.GetString("username")
	password := this.GetString("password")
	if ok, u := models.CanLogin(username, password); ok {
		this.Data["json"] = map[string]interface{}{
			"code": 0,
			"msg":  "login success",
			"profile": map[string]interface{}{
				"username": u.Username,
				"role":     u.Role,
				"lastTime": u.LastLogin,
				"lastLoc":  u.LastLoc,
			},
		}
		err := models.UpdateLogin(u, time.Now(), getIPLoc(this.Ctx.Input.IP()))
		if err != nil {
			log.Println(err)
		}
	} else {
		this.Data["json"] = map[string]interface{}{
			"code": -10001,
			"msg":  "user not exist",
		}
	}
	err := this.ServeJSON()
	if err != nil {
		log.Println(err)
		return
	}
}

// @router /list [get]
func (this *UserController) GetAll() {
	var query struct {
		Search    string `form:"search"`
		PageIndex int    `form:"pageIndex"`
		PageSize  int    `form:"pageSize"`
	}
	if err := this.ParseForm(&query); err != nil {
		log.Println(err)
	}
	var (
		courses []*models.User
		total   int
	)
	if query.Search == "" {
		courses, total = models.ListUsers(getOffset(query.PageIndex, query.PageSize), query.PageSize)
	} else {
		courses, total = models.SearchUsers(getOffset(query.PageIndex, query.PageSize), query.PageSize, query.Search)
	}
	this.Data["json"] = map[string]interface{}{
		"list":      courses,
		"pageTotal": total,
	}
	err := this.ServeJSON()
	if err != nil {
		log.Println(err)
		return
	}
}
