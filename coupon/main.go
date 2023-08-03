package main

import (
	"fetch-coupon/src/sign"
	"fetch-coupon/src/task"
	"fetch-coupon/src/utils"
	wgCommon "fetch-coupon/src/wg"
	"fmt"
	"time"
)


func main() {
	utils.StartServer()

	secTime := utils.GetSecTime()

	fetchCoupon(secTime)
	// for {} // 死循环不让进程结束
}

func fetchCoupon(secTime *utils.SecTimeData) {
	data := sign.SignDuration(secTime)
	// fmt.Printf("%+v\n",data)

	if len(data) == 0 {
		fmt.Println("没有签名")
		return
	}

	d := secTime.Sec.Sub(secTime.Mt) - task.Early*time.Millisecond
	fmt.Println("在", d, "后抢券")

	// t := time.NewTimer(d)
	// <-t.C

	for _, sd := range data {
		wgCommon.MainWg.Add(1)
		go sign.Fc(sd)
	}

	wgCommon.MainWg.Wait()
}


