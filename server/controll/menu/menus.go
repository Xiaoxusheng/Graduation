package menu

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"server/dao"
	"server/global"
	"server/models"
	"server/result"
	"time"
)

// GetMenuList 获取菜单接口
func GetMenuList(c *gin.Context) {
	//判断身份

	//下发对应的菜单
	val := global.Global.Redis.Get(global.Global.Ctx, global.Menus).Val()
	if val != "" {
		menu := make([]global.Menu, 10)
		err := json.Unmarshal([]byte(val), &menu)
		if err != nil {
			global.Global.Log.Error(err)
			result.Fail(c, global.DataUnmarshal, global.DataUnmarshalError)
			return
		}
		result.Ok(c, menu)
		return
	}
	list, err := dao.GetMenuList()
	if err != nil {
		return
	}
	menuList := make([]models.Menu, 0)
	//处理
	for i := 0; i < len(list); i++ {
		//父节点
		if list[i].ParentPath == "" {
			menuList = append(menuList, list[i])
		}
	}
	global.Global.Log.Info(menuList)

	for i := 0; i < len(list); i++ {
		if list[i].ParentPath != "" {
			for j := 0; j < len(menuList); j++ {
				if menuList[j].Children == nil {
					menuList[j].Children = make([]models.Menu, 0)
				}
				if menuList[j].MenuUrl == list[i].ParentPath {
					menuList[j].Children = append(menuList[j].Children, list[i])
				}
			}
		}
	}
	//
	//发送处理好数据
	go func() {
		marshal, err := json.Marshal(menuList)
		if err != nil {
			global.Global.Log.Error(err)
			return
		}
		_, err = global.Global.Redis.Set(global.Global.Ctx, global.Menus, marshal, global.MenuTime*time.Second).Result()
		if err != nil {
			global.Global.Log.Error(err)
		}
	}()
	global.Global.Log.Info(menuList)
	result.Ok(c, menuList)
}

//修改菜单

// AddMenu 增加菜单
func AddMenu(c *gin.Context) {
	//
	menu := new(global.Menu)
	err := c.Bind(menu)
	if err != nil {
		result.Fail(c, global.BadRequest, global.QueryError)
		global.Global.Log.Error(err)
		return
	}

}

//删除菜单
