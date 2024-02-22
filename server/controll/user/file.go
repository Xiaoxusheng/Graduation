package user

import (
	"github.com/gin-gonic/gin"
	"server/global"
	"server/result"
	"server/utils"
)

// Upload 文件上传
func Upload(c *gin.Context) {
	list := make([]global.UrlList, 0, 9)
	n := 0
	form, err := c.MultipartForm()
	if err != nil {
		global.Global.Log.Warn("上传失败" + err.Error())
		result.Fail(c, global.ServerError, global.QueryError)
		return
	}
	if len(form.File["file"]) > global.FileNumber {
		global.Global.Log.Warn("文件数量超过限制！")
		result.Fail(c, global.DataTooLarge, global.FileError)
		return
	}
	urlChan := make(chan global.UrlList, global.FileNumber)
	for i, r := range form.File["file"] {
		go utils.Upload(r, urlChan, i+1)
	}
	for n != len(form.File["file"]) {
		select {
		case s := <-urlChan:
			if s.Url != "" {
				n++
				list = append(list, s)
				continue
			}
			n++
		}
	}
	result.Ok(c, list)
}
