package main

import (
	"fetch-coupon/src/sign"
	"fetch-coupon/src/task"
	"fetch-coupon/src/utils"
	"fmt"
	"time"
)

func main() {
	go utils.StartServer()


	secTime := utils.GetSecTime()
	if secTime == nil {
		return
	}

	go utils.GetCouponInfo(secTime)
	go fetchCoupon(secTime)

	for {}
}

func fetchCoupon(secTime *utils.SecTimeData) {
	var data []sign.SignData

	go func() {
		data = sign.SignDuration(secTime)
	}()

	d := secTime.Sec.Sub(secTime.Mt) - task.Early*time.Millisecond

	fmt.Println("在", d, "后抢券")
	t := time.NewTimer(d)
	<-t.C

	if len(data) == 0 {
		fmt.Println("没有签名")
		return
	}

	for _, sd := range data {
		go sign.Fc(sd)
	}
}


