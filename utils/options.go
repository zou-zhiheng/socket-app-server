package utils

import (
	"strconv"
	"time"
)

// GetLimit 获取本次操作条数的上限，优化数据库I/O
func GetLimit(startTime, endTime, interval, layout string) (int64, error) {

	if interval == "0" || interval == "" {
		interval = "1"
	}

	sTime, err := time.Parse(layout, startTime)
	if err != nil {
		return 0, err
	}
	eTime, err := time.Parse(layout, endTime)
	if err != nil {
		return 0, err
	}

	inter, err := strconv.Atoi(interval)

	//开始计算相差多少分钟
	limit := int64(int(eTime.Sub(sTime).Minutes()) / inter)

	return limit, nil

}
