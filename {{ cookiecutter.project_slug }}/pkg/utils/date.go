package utils

import (
	"fmt"
	"time"

	"github.com/spf13/cast"
)

var parseLayouts = []string{
	"2006-01-02 15:04:05",
	"2006-1-2 15:04:05",
	"2006-1-2 15:4:5",
	"2006/01/02 15:04:05",
	"2006/1/2 15:04:05",
	"2006/1/2 15:4:5",
	"2006年01月02日 15时04分05秒",
	"2006年1月2日 15时04分05秒",
	"2006年1月2日 15时4分5秒",
	"2006年01月02日 15时04分",
	"2006年1月2日 15时04分",
	"2006年1月2日 15时4分",
	"2006年01月02日 15:04:05",
	"2006年1月2日 15:04:05",
	"2006年1月2日 15:4:5",
}

func ToDateE(v interface{}) (time.Time, error) {
	beijing, _ := time.LoadLocation("Asia/Shanghai")
	for _, layout := range parseLayouts {
		if t, err := time.ParseInLocation(layout, cast.ToString(v), beijing); err == nil {
			return t, nil
		}
	}
	if t, err := cast.ToTimeE(v); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("不能被识别的时间%s", cast.ToString(v))
}
