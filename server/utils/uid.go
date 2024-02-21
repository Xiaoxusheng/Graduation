package utils

import (
	"crypto/md5"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"server/dao"
	"server/global"
	"strconv"
	"strings"
	"time"
)

func GetUidV5(u string) string {
	return uuid.NewV5(uuid.NewV4(), u).String()
}

func GetUidV4() string {
	return uuid.NewV4().String()
}

func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// CreateUid 生成唯一id
func CreateUid(key string) int64 {
	t := time.Now()
	t1 := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
	//获取redis中的数据
	val := global.Global.Redis.IncrBy(global.Global.Ctx, key+t.Format(time.DateOnly), 1).Val()
	return t.Unix() - t1.Unix()<<32 | val
}

// GetUid 生成部门的员工编号
func GetUid(types int32) (int64, error) {
	val := global.Global.Redis.HGet(global.Global.Ctx, global.Uid, string(types)).Val()
	s := strings.Builder{}
	//不存在
	if val == "" {
		/*查一下数据库，有可能是redis导致的查不到数据*/
		employer, err := dao.GetEmployer(types)
		if err != nil {
			global.Global.Log.Error(err)
		}

		//Uif不存在
		if employer == nil {
			rand.NewSource(time.Now().UnixNano())
			s.WriteString(strconv.Itoa(int(types)))
			fmt.Println(s.String())
			for i := 0; i < 6; i++ {
				s.WriteString("0")
			}
			s.WriteString("1")
			_, err = global.Global.Redis.HSet(global.Global.Ctx, global.Uid, string(types), s.String()).Result()
			if err != nil {
				global.Global.Log.Error(err)
				return 0, err
			}
			fmt.Println(s.String())
			i, err := strconv.Atoi(s.String())
			if err != nil {
				global.Global.Log.Error(err)
				return 0, err
			}
			return int64(i), nil
		} else {
			//redis出错了，数据库中存在
			val = strconv.FormatInt(employer.Uid+1, 10)
			_, err = global.Global.Redis.HSet(global.Global.Ctx, global.Uid, string(types), employer.Uid+1).Result()
			if err != nil {
				global.Global.Log.Error(err)
				return 0, err
			}
			return employer.Uid + 1, nil
		}
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		global.Global.Log.Error(err)
		return 0, err
	}
	_, err = global.Global.Redis.HSet(global.Global.Ctx, global.Uid, string(types), i+1).Result()
	if err != nil {
		global.Global.Log.Error(err)
		return 0, err
	}
	return int64(i) + 1, nil
}
