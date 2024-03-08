package root

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/dao"
	"server/global"
	"server/models"
	"server/result"
	"server/utils"
	"strings"
)

/*RBAC*/

// AddRolesForUser 为用户分配角色
func AddRolesForUser(c *gin.Context) {
	role := c.PostForm("role")
	user := c.PostForm("user")
	if role == "" || user == "" {
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	//限制
	//写入redis,记录用户所拥有的角色
	_, err := global.Global.Redis.SAdd(global.Global.Ctx, global.User+user, role).Result()
	if err != nil {
		result.Fail(c, global.ServerError, global.AddRoleFail)
		global.Global.Log.Error(err)
		return
	}
	_, err = global.Global.CasBin.AddRoleForUser(user, role)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.AddRoleFail)
		return
	}
	result.Ok(c, nil)
	return
}

// AddResource 给所有的角色添加权限可以访问的资源
//func AddResource(c *gin.Context) {
//	role := c.PostForm("role")
//	method := strings.ToUpper(c.PostForm("method"))
//	path := c.PostForm("path")
//	if method == "" || path == "" || role == "" {
//		result.Fail(c, global.BadRequest, global.QueryError)
//		return
//	}
//	//判断角色是否存在
//	val := global.Global.Redis.SMembers(global.Global.Ctx, global.Role+role).Val()
//	if len(val) == 0 {
//		result.Fail(c, global.BadRequest, global.RoleNotfound)
//		return
//	}
//	//分配资源
//	for i := 0; i < len(val); i++ {
//		_, err := global.Global.CasBin.AddPermissionForUser(val[i], path, method)
//		if err != nil {
//			global.Global.Log.Error(err)
//			result.Fail(c, global.DataValidationError, global.AddPermissionFail)
//			return
//		}
//	}
//	result.Ok(c, nil)
//	return
//}

// AddPermissionForUser 为角色添加单个资源
func AddPermissionForUser(c *gin.Context) {
	role := c.PostForm("role")
	method := strings.ToUpper(c.PostForm("method"))
	path := c.PostForm("path")
	if method == "" || path == "" || role == "" {
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	ok, err := global.Global.CasBin.AddPermissionForUser(role, path, method)
	if err != nil || !ok {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.AddPermissionFail)
		return
	}
	result.Ok(c, nil)
}

// DeletePermissionForUser 删除用户权限
func DeletePermissionForUser(c *gin.Context) {
	user := c.PostForm("user")
	method := strings.ToUpper(c.PostForm("method"))
	path := strings.ReplaceAll(c.PostForm("path"), " ", "")
	if method == "" || path == "" || user == "" {
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	//判断是否存在
	if !global.Global.CasBin.HasPolicy(user, path, method) {
		result.Fail(c, global.BadRequest, global.PermissionNotFound)
		return
	}
	_, err := global.Global.CasBin.DeletePermissionForUser(user, path, method)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.DelPermissionFail)
		return
	}
	result.Ok(c, nil)
	return
}

// DeletePermissionsForUser 删除用户所有的权限
func DeletePermissionsForUser(c *gin.Context) {
	user := c.PostForm("user")
	if user == "" {
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	ok, err := global.Global.CasBin.DeletePermissionsForUser(user)
	if err != nil || !ok {
		result.Fail(c, global.ServerError, global.DelPermissionFail)
		return
	}
	result.Ok(c, nil)
}

// UpdatePolicy 修改用户权限
func UpdatePolicy(c *gin.Context) {
	//获取用户的权限
	user := c.PostForm("user")
	path := c.PostForm("path")
	newPath := strings.ReplaceAll(c.PostForm("newPath"), " ", "")
	if user == "" || path == "" || newPath == "" {
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	fmt.Println(user, path, newPath)
	list, _ := global.Global.CasBin.GetPermissionsForUser(user)
	fmt.Println(list)
	for i := 0; i < len(list); i++ {
		old := list[i]
		if list[i][1] == path {
			list[i][1] = newPath
			_, err := global.Global.CasBin.UpdatePolicy(old, list[i])
			if err != nil {
				global.Global.Log.Error(err)
				result.Fail(c, global.ServerError, global.UpdatePermissionFail)
				return
			}
		}
	}
	err := global.Global.CasBin.SavePolicy()
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.UpdatePermissionFail)
		return
	}
	result.Ok(c, nil)
}

// DeleteRoleForUser 删除用户的角色
func DeleteRoleForUser(c *gin.Context) {
	//role 角色 user用户
	role := c.PostForm("role")
	user := c.PostForm("user")
	if role == "" || user == "" {
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	_, err := global.Global.Redis.SRem(global.Global.Ctx, global.User+role, user).Result()
	if err != nil {
		result.Fail(c, global.ServerError, global.RoleNotfound)
		global.Global.Log.Error(err)
		return
	}
	//删除角色
	ok, err := global.Global.CasBin.DeleteRoleForUser(user, role)
	if err != nil || !ok {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.DelRoleFail)
		return
	}
	result.Ok(c, nil)
	return
}

// GetPermissionsForUser 查看用户的权限
func GetPermissionsForUser(c *gin.Context) {
	user := c.PostForm("user")
	if user == "" {
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	list, err := global.Global.CasBin.GetPermissionsForUser(user)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.GetPermissionsForUserFail)
		return
	}
	result.Ok(c, list)
}

// GetAllNamedSubjects 查看所有用户权限
func GetAllNamedSubjects(c *gin.Context) {
	list := global.Global.CasBin.GetNamedPolicy("p")
	result.Ok(c, list)
	return
}

// DeleteRole 删除角色，同时也删除所有的权限
func DeleteRole(c *gin.Context) {
	role := c.PostForm("role")
	if role == "" {
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	ok, err := global.Global.CasBin.DeleteRole(role)
	if err != nil || !ok {
		global.Global.Log.Error(err)
		return
	}
	result.Ok(c, nil)
}

func Add(c *gin.Context) {
	//role := c.PostForm("role")
	//menu := c.PostForm("menu")
	//if menu == "" || role == "" {
	//	result.Fail(c, global.BadRequest, global.QueryError)
	//	return
	//}
	//ok, err := global.Global.CasBin.AddPermissionForUser(role, menu)
	//if err != nil || !ok {
	//	global.Global.Log.Error(err)
	//	result.Fail(c, global.ServerError, global.AddPermissionFail)
	//	return
	//}
	id := c.GetString("identity")
	fmt.Println(global.Global.CasBin.GetAllRoles())
	fmt.Println(global.Global.CasBin.GetRolesForUser(id))
	result.Ok(c, nil)
}

// AddRoleMenu 给角色分配菜单
func AddRoleMenu(c *gin.Context) {
	role := c.PostForm("role")
	menu := c.PostFormArray("menu")
	if len(menu) == 0 || role == "" {
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	//判断角色是否存在
	//var ok bool
	//for i := 0; i < len(global.Global.CasBin.GetAllRoles()); i++ {
	//	if global.Global.CasBin.GetAllRoles()[i] == role {
	//		ok = true
	//	}
	//}
	//if !ok {
	//	result.Fail(c, global.BadRequest, global.RoleNotfound)
	//	return
	//}
	//	操作数据库
	for i := 0; i < len(menu); i++ {
		err := dao.InsertRoleMenu(&models.RoleMenu{
			Identity: utils.GetUidV4(),
			Role:     role,
			Menu:     menu[i],
		})
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.ServerError, global.AddMenuRoleFail)
		}
	}
	result.Ok(c, nil)
}
