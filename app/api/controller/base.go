package controller

import (
	"dj-api/tools/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"math"
	"strconv"
	"strings"
)

const (
	PageSize = 10
)

//获取分页起始偏移量
func GetPageOffset(c *gin.Context) int {
	offset := 0
	page := com.StrTo(c.Query("p")).MustInt()
	if page > 0 {
		offset = (page - 1) * PageSize
	}
	return offset
}

//绑定请求参数
func BindAndValid(c *gin.Context, form interface{}) bool {
	if err := c.ShouldBindJSON(form); err != nil {
		logger.Info("form bind err:", err.Error())
		return false
	}
	return true
}

/**
版本号比较
return 0 异常
return -1 已经是最新版本
return 1 可以更新
*/
func VersionCompare(curStr, newStr string) int8 {
	curStr = strings.TrimSpace(strings.Replace(curStr, ".", "", -1))
	newStr = strings.TrimSpace(strings.Replace(newStr, ".", "", -1))
	diff := len(curStr) - len(newStr) //版本号长度差
	if diff > 0 {
		newStr += strings.Repeat("0", diff)
	} else if diff < 0 {
		curStr += strings.Repeat("0", -diff)
	}
	cur, err := strconv.Atoi(curStr)
	if err != nil {
		return 0
	}
	news, err := strconv.Atoi(newStr)
	if err != nil {
		return 0
	}
	if cur >= news {
		return -1
	} else {
		return 1
	}
}

//格式化客户端输出
func FormatAmountToYUan(amount int64) float64 {
	if amount == 0 {
		return 0
	}
	a, err := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(amount)/100), 64)
	if err != nil {
		logger.Error(err)
		return 0
	}
	return a
}

//格式化Go服务端
func FormatAmountToFen(amount int) int {
	return int(math.Floor(float64(amount / 100)))
}
