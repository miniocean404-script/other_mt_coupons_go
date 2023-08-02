package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

var (
	cookies = map[string]string{
		"6886": "userId=906082484;token=AgFnIc_G9WnCR-CCPYyQJl_CiHAYyc_dOiHo8yxhnGdVDrLTGKLyOtSPoBDUCp97sNmMiDI96fUcEgAAAADHGQAABg8OohMbfRiPzBgxgtTJdaVYnvlTZJb6BE4et1c4G85ze68Liv6XrsClwWdLhamC",
		"6408": "userId=3170637747;token=AgGJI9GQYEiBsRxSPWsIxdZK9RZcTvDSDPOxRqLWMn8PcX0AE_lrSGWwm573RMrwydTL6AvM_cqh5wAAAACpGQAAaBsGenSWPG4G7yk8Mn2Tv5VeavdF-L57DDaWJKpktMt7WCKO8oceddBorwIRLYG-",
	}
	wg sync.WaitGroup
)

const (
	Num            = 5
	Early          = 100 // 毫秒
	SecAt          = "10:30:00"
	Id             = "00B223429B424F7A910C0D4885957E99"
	GdId           = "379397"
	PageId         = "378931"
	InstanceId     = "16618616100670.97030510386642830"
	SignUrl        = "http://127.0.0.1:9588/api/sign"
	CouponInfoUrl  = "https://promotion.waimai.meituan.com/lottery/limitcouponcomponent/info"
	FetchCouponUrl = "https://promotion.waimai.meituan.com/lottery/limitcouponcomponent/fetchcoupon"

	// 经纬度
	ActualLng = "117.23344"
	ActualLat = "31.82658"
)

func main() {
	go runSign()

	secTime := getSecTime()

	if secTime == nil {
		return
	}
	go getCoupon(secTime)
	go fetchCoupon(secTime)

	for {

	}
}

func runSign() {
	_, err := http.Get(SignUrl)
	if err == nil {
		return
	}

	switch runtime.GOOS {
	case "windows", "linux", "darwin":
		exec.Command("node", "../sign/index.js").Run()
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

func fetchCoupon(secTime *secTimeData) {
	var data []signData

	go func() {
		data = signDuration(secTime)
	}()

	d := secTime.sec.Sub(secTime.mt) - Early*time.Millisecond

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

func signDuration(secTime *secTimeData) []signData {
	d := secTime.sec.Sub(secTime.mt) - 30*time.Second
	fmt.Println("在", d, "后获取签名")
	t := time.NewTimer(d)
	<-t.C
	return sign(fetchCouponUrl())
}

func fc(sd signData) {
	wg.Add(Num)
	sdSignLen := len(sd.sign)
	for i := 0; i < Num; i++ {
		if i >= sdSignLen {
			return
		}
		go sd.fcReq(i)
	}
	wg.Wait()
}

func (sd *signData) fcReq(i int) {
	defer wg.Done()

	s := sd.sign[i]

	payload, _ := json.Marshal(s.Data.Data)
	httpReq, err := http.NewRequest(http.MethodPost, sd.url, bytes.NewReader(payload))
	if err != nil {
		fmt.Println("err", err)
		return
	}

	httpReq.Header = baseHeader(sd.cookie, s.Data.Mtgsig)
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
	val.Set("version", "1")
	val.Set("utmSource", "wxshare")
	val.Set("pageId", PageId)
	val.Set("instanceId", InstanceId)
	val.Set("componentId", InstanceId)
	val.Set("actualLng", ActualLng)
	val.Set("actualLat", ActualLat)
	return fmt.Sprintf("%s?%s", FetchCouponUrl, val.Encode())
	//return fmt.Sprintf("https://promotion.waimai.meituan.com/lottery/limitcouponcomponent/fetchcoupon?%s", val.Encode())
}

type (
	signData struct {
		cookie string
		name   string
		url    string
		sign   []signResponse
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
		var signResp []signResponse
		for i := 0; i < Num; i++ {
			resp := getSign(cookie, fetchCouponUrl)
			if resp == nil {
				continue
			}
			signResp = append(signResp, *resp)
		}

		data = append(data, signData{
			cookie: cookie,
			name:   name,
			sign:   signResp,
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

func getCoupon(secTime *secTimeData) {
	d := secTime.sec.Sub(secTime.mt) - 40*time.Second

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

type secTimeData struct {
	sec time.Time
	mt  time.Time
}

func getSecTime() *secTimeData {
	today := time.Now().Format("2006-01-02")
	t, err := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s %s", today, SecAt), time.Local)
	if err != nil {
		return nil
	}

	return &secTimeData{
		sec: t,
		mt:  getMtTime(),
	}
}

type mtServerTime struct {
	Data    int64  `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
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
	t := new(mtServerTime)
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
