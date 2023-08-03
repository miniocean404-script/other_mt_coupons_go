package main

import (
	"fetch-coupon/src/sign"
	"fetch-coupon/src/task"
	"fetch-coupon/src/utils"
	"fmt"
	"time"
)

func main() {
	utils.StartServer()

	secTime := utils.GetSecTime()
	go fetchCoupon(secTime)

	for {}
}

func fetchCoupon(secTime *utils.SecTimeData) {
	data := sign.SignDuration(secTime)
	fmt.Printf("%+v\n",data)

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


