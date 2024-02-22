package utils

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"server/config"
	"server/global"
)

func Upload(file *multipart.FileHeader, c chan global.UrlList, index int) {
	//上传
	global.Global.Log.Info("第", index, "个文件上传")

	u, _ := url.Parse(config.Config.TencentCos.Url)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.Config.TencentCos.SecretId,  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
			SecretKey: config.Config.TencentCos.SecretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
		},
	})
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			// 设置默认的进度回调函数
			Listener: &cos.DefaultProgressListener{},
		},
	}

	// Case1 使用 Put 上传对象
	key := GetUidV4() + path.Ext(file.Filename)
	files, err := file.Open()
	defer func(files multipart.File) {
		err = files.Close()
		if err != nil {
			global.Global.Log.Error(err)
		}
	}(files)
	_, err = client.Object.Put(context.Background(), key, files, opt)
	if err != nil {
		c <- global.UrlList{
			Url:   "",
			Index: index,
		}
		return
	}
	c <- global.UrlList{
		Url:   config.Config.TencentCos.Url + key,
		Index: index,
	}
	global.Global.Log.Info("上传完成")
}
