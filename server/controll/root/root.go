package root

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global"
	"server/result"
	"strings"
)

/*RBAC*/

// AddResource 给所有的角色添加权限可以访问的资源
func AddResource(c *gin.Context) {
	user := c.PostForm("user")
	method := strings.ToUpper(c.PostForm("method"))
	path := c.PostForm("path")
	if method == "" || path == "" || user == "" {
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	//判断角色是否存在
	val := global.Global.Redis.SMembers(global.Global.Ctx, global.Role+user).Val()
	if len(val) == 0 {
		result.Fail(c, global.BadRequest, global.RoleNotfound)
		return
	}
	//分配资源
	for i := 0; i < len(val); i++ {
		_, err := global.Global.CasBin.AddPermissionForUser(val[i], path, method)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.DataValidationError, global.AddPermissionFail)
			return
		}
	}
	result.Ok(c, nil)
	return
}

// AddRolesForUser 为用户分配角色
func AddRolesForUser(c *gin.Context) {
	// role 为角色 identity  user 是角色 admin或root这种
	role := c.PostForm("role")
	user := c.PostForm("user")
	if role == "" || user == "" {
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	//限制
	//写入redis
	_, err := global.Global.Redis.SAdd(global.Global.Ctx, global.Role+user, role).Result()
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

// DeletePermissionForUser 删除用户能访问的资源
func DeletePermissionForUser(c *gin.Context) {
	role := c.PostForm("role")
	method := strings.ToUpper(c.PostForm("method"))
	path := strings.ReplaceAll(c.PostForm("path"), " ", "")
	if method == "" || path == "" || role == "" {
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}

	//判断是否存在
	if !global.Global.CasBin.HasPolicy(role, path, method) {
		result.Fail(c, global.BadRequest, global.PermissionNotFound)
		return
	}
	_, err := global.Global.CasBin.DeletePermissionForUser(role, path, method)
	if err != nil {
		global.Global.Log.Error(err)
		result.Fail(c, global.ServerError, global.DelPermissionFail)
		return
	}
	result.Ok(c, nil)
	return
}

// DeleteRoleForUser 删除用户的角色
func DeleteRoleForUser(c *gin.Context) {
	role := c.PostForm("role")
	user := c.PostForm("user")
	if role == "" || user == "" {
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	_, err := global.Global.Redis.SRem(global.Global.Ctx, global.Role+user, role).Result()
	if err != nil {
		global.Global.Log.Error(err)
		return
	}
	_, err = global.Global.CasBin.DeleteRoleForUser(user, role)
	if err != nil {
		result.Fail(c, global.ServerError, global.DelRoleFail)
		return
	}
	result.Ok(c, nil)
	return
}

// GetPermissionsForUser 查看用户能访问的资源
func GetPermissionsForUser(c *gin.Context) {
	role := c.PostForm("role")
	if role == "" {
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	list, err := global.Global.CasBin.GetPermissionsForUser(role)
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

// UpdatePolicy 修改用户权限
func UpdatePolicy(c *gin.Context) {
	//获取用户的权限
	role := c.PostForm("role")
	path := c.PostForm("path")
	newPath := strings.ReplaceAll(c.PostForm("newPath"), " ", "")
	if role == "" || path == "" || newPath == "" {
		result.Fail(c, global.BadRequest, global.QueryError)
		return
	}
	fmt.Println(role, path, newPath)
	list, _ := global.Global.CasBin.GetPermissionsForUser(role)
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
