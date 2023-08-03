package utils

import (
	"encoding/json"
	"fetch-coupon/src/task"
	"fmt"
	"io"
	"net/http"
	"time"
)

type SecTimeData struct {
	Sec time.Time
	Mt  time.Time
}

type MtServerTime struct {
	Data    int64  `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func GetSecTime() *SecTimeData {
	today := time.Now().Format("2006-01-02")
	t, err := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s %s", today, task.SecAt), time.Local)
	if err != nil {
		return nil
	}

	return &SecTimeData{
		Sec: t,
		Mt:  getMtTime(),
	}
}

func getMtTime() time.Time {
	now := time.Now()
	resp, err := http.Get("https://cube.meituan.com/ipromotion/cube/toc/component/base/getServerCurrentTime")
	if err != nil {
		return now
	}
	if resp.StatusCode != http.StatusOK {
		return now
	}

	defer resp.Body.Close()
	all, _ := io.ReadAll(resp.Body)
	t := new(MtServerTime)
	if err := json.Unmarshal(all, t); err != nil {
		return now
	}

	if t.Status != 0 {
		return now
	}

	milli := time.UnixMilli(t.Data)

	if milli.Year() == 1970 {
		return now
	}
	return milli.Add(-50 * time.Millisecond)
}
