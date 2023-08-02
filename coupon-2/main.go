package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	cookies = map[string]string{
		//"6886": "userId=906082484;token=AgFnIc_G9WnCR-CCPYyQJl_CiHAYyc_dOiHo8yxhnGdVDrLTGKLyOtSPoBDUCp97sNmMiDI96fUcEgAAAADHGQAABg8OohMbfRiPzBgxgtTJdaVYnvlTZJb6BE4et1c4G85ze68Liv6XrsClwWdLhamC",
		"6408": "userId=3170637747;token=AgGJI9GQYEiBsRxSPWsIxdZK9RZcTvDSDPOxRqLWMn8PcX0AE_lrSGWwm573RMrwydTL6AvM_cqh5wAAAACpGQAAaBsGenSWPG4G7yk8Mn2Tv5VeavdF-L57DDaWJKpktMt7WCKO8oceddBorwIRLYG-",
	}
	wg sync.WaitGroup
)

const (
	Num        = 30
	SecAt      = "10:59:59"
	Id         = "419967B3A4064140BA78E6A046DF0FC1"
	GdId       = "379391"
	PageId     = "378925"
	InstanceId = "16619982800580.30892480633143027"
	SignUrl    = "http://127.0.0.1:9588/api/sign"
	//CouponInfoUrl  = "https://promotion.waimai.meituan.com/lottery/couponcomponent/info"
	//FetchCouponUrl = "https://promotion.waimai.meituan.com/lottery/couponcomponent/fetchcomponentcoupon"
	CouponInfoUrl  = "https://promotion.waimai.meituan.com/lottery/limitcouponcomponent/info"
	FetchCouponUrl = "https://promotion.waimai.meituan.com/lottery/limitcouponcomponent/fetchcoupon"
)

func main() {

	secTime := getSecTime()
	if secTime == nil {
		return
	}
	fmt.Println("在", secTime, "开始任务")

	go getCoupon(secTime)
	go fetchCoupon(secTime)

	for {

	}
}

type (
	fetchCouponRequest struct {
		CType         string `json:"cType"`
		FpPlatform    int    `json:"fpPlatform"`
		WxOpenId      string `json:"wxOpenId"`
		AppVersion    string `json:"appVersion"`
		MtFingerprint string `json:"mtFingerprint"`
	}
)

func fetchCoupon(secTime *time.Time) {
	var data []signData

	go func() {
		data = signDuration(secTime)
	}()

	d := secTime.Sub(time.Now()) + 850*time.Millisecond

	fmt.Println("在", d, "后抢券")
	t := time.NewTimer(d)
	<-t.C

	if len(data) == 0 {
		fmt.Println("没有签名")
		return
	}

	for _, sd := range data {
		go fc(sd)
	}
}

func signDuration(secTime *time.Time) []signData {
	d := secTime.Sub(time.Now()) - 20*time.Second
	fmt.Println("在", d, "后获取签名")
	t := time.NewTimer(d)
	<-t.C
	return sign(fetchCouponUrl())
}

func fc(sd signData) {
	wg.Add(Num)
	for i := 0; i < Num; i++ {
		go sd.fcReq()
	}
	wg.Wait()
}

func (sd *signData) fcReq() {
	defer wg.Done()
	payload, _ := json.Marshal(sd.sign.Data.Data)
	httpReq, err := http.NewRequest(http.MethodPost, sd.url, bytes.NewReader(payload))
	if err != nil {
		fmt.Println("err", err)
		return
	}
	httpReq.Header = baseHeader(sd.cookie, sd.sign.Data.Mtgsig)

	client := &http.Client{}
	response, err := client.Do(httpReq)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.0000"), fmt.Sprintf("====== data %s ========", sd.name), string(data))
}

func fetchCouponUrl() string {
	val := url.Values{}
	val.Set("couponReferId", Id)
	val.Set("geoType", "2")
	val.Set("gdPageId", GdId)
	val.Set("pageId", PageId)
	val.Set("instanceId", InstanceId)
	val.Set("componentId", InstanceId)
	return fmt.Sprintf("%s?%s", FetchCouponUrl, val.Encode())
	//return fmt.Sprintf("https://promotion.waimai.meituan.com/lottery/limitcouponcomponent/fetchcoupon?%s", val.Encode())
}

type (
	signData struct {
		cookie string
		name   string
		url    string
		sign   *signResponse
	}
	signRequest struct {
		Cookie string `json:"cookie"`
		Url    string `json:"url"`
	}
	signResponse struct {
		Code int `json:"code"`
		Data struct {
			Data   *fetchCouponRequest `json:"data"`
			Mtgsig string              `json:"mtgsig"`
		} `json:"data"`
	}
)

func sign(fetchCouponUrl string) []signData {
	var data []signData
	for name, cookie := range cookies {
		resp := getSign(cookie, fetchCouponUrl)
		if resp == nil {
			continue
		}
		data = append(data, signData{
			cookie: cookie,
			name:   name,
			sign:   resp,
			url:    fetchCouponUrl,
		})
	}
	return data
}

func getSign(cookie, url string) *signResponse {
	req := &signRequest{
		Cookie: cookie,
		Url:    url,
	}

	payload, _ := json.Marshal(req)

	httpReq, _ := http.NewRequest(http.MethodPost, SignUrl, bytes.NewReader(payload))
	httpReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(httpReq)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)

	resp := new(signResponse)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil
	}
	if resp.Code != 0 {
		return nil
	}
	return resp
}

func getCoupon(secTime *time.Time) {
	d := secTime.Sub(time.Now()) - 40*time.Second

	fmt.Println("在", d, "后执行couponInfo")
	timer := time.NewTimer(d)

	<-timer.C
	httpReq, _ := http.NewRequest(http.MethodGet, couponInfoUrl(), nil)
	client := &http.Client{}
	for _, cookie := range cookies {
		httpReq.Header = baseHeader(cookie, "")
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
	return fmt.Sprintf("%s?couponReferIds=%s", CouponInfoUrl, Id)
}

func getSecTime() *time.Time {
	today := time.Now().Format("2006-01-02")
	t, err := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s %s", today, SecAt), time.Local)
	if err != nil {
		return nil
	}
	return &t
}

func baseHeader(cookie, mtgsig string) http.Header {
	headers := http.Header{}
	headers.Set("Cookie", cookie)
	headers.Set("Origin", "https://market.waimai.meituan.com")
	headers.Set("Referer", "https://market.waimai.meituan.com/")
	headers.Set("Content-Type", "application/json")
	headers.Set("sec-ch-ua", `"Not_A Brand";v="99", "Google Chrome";v="109", "Chromium";v="109"`)
	headers.Set("sec-ch-ua-mobile", "?1")
	headers.Set("sec-ch-ua-platform", `"Android"`)
	headers.Set("Sec-Fetch-Dest", "empty")
	headers.Set("Sec-Fetch-Mode", "cors")
	headers.Set("Sec-Fetch-Site", "same-site")
	headers.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Mobile Safari/537.36")
	headers.Set("mtgsig", mtgsig)
	headers.Set("Accept", "application/json, text/plain, */*")
	return headers
}
