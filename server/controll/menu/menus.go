package menu

import (
	"github.com/gin-gonic/gin"
	"server/dao"
	"server/global"
	"server/models"
	"server/result"
)

// GetMenuList 获取菜单接口
func GetMenuList(c *gin.Context) {
	//判断身份

	//下发对应的菜单
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
	global.Global.Log.Info(menuList)
	result.Ok(c, menuList)
}
