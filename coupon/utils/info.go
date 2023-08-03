package utils

import (
	"fetch-coupon/task"
	"fmt"
	"io"
	"net/http"
	"time"
)

func GetCouponInfo(secTime *SecTimeData) {
	d := secTime.Sec.Sub(secTime.Mt) - 40*time.Second

	fmt.Println("在", d, "后执行couponInfo")
	timer := time.NewTimer(d)

	<-timer.C
	httpReq, _ := http.NewRequest(http.MethodGet, couponInfoUrl(), nil)
	client := &http.Client{}
	for _, cookie := range task.Cookies {
		httpReq.Header = BaseHeader(cookie, "")
		response, err := client.Do(httpReq)
		if err != nil {
			fmt.Println(err)
			return
		}
		data, err := io.ReadAll(response.Body)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

func couponInfoUrl() string {
	return fmt.Sprintf("%s?couponReferIds=%s", task.CouponInfoUrl, task.Id)
}
