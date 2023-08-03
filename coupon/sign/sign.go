package sign

import (
	"bytes"
	"encoding/json"
	"fetch-coupon/task"
	"fetch-coupon/utils"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var(
	wg sync.WaitGroup
)

// const (
// 	Num            = 5
// 	Early          = 100 // 毫秒
// 	SecAt          = "10:30:00"
// 	Id             = "00B223429B424F7A910C0D4885957E99"
// 	GdId           = "379397"
// 	PageId         = "378931"
// 	InstanceId     = "16618616100670.97030510386642830"
// 	SignUrl        = "http://127.0.0.1:9588/api/sign"
// 	CouponInfoUrl  = "https://promotion.waimai.meituan.com/lottery/limitcouponcomponent/info"
// 	FetchCouponUrl = "https://promotion.waimai.meituan.com/lottery/limitcouponcomponent/fetchcoupon"

// 	// 经纬度
// 	ActualLng = "117.23344"
// 	ActualLat = "31.82658"
// )


type (
	SignData struct {
		Cookie string
		Name   string
		Url    string
		Sign   []signResponse
	}

	signRequest struct {
		Cookie string `json:"cookie"`
		Url    string `json:"url"`
	}

	signResponse struct {
		Code int `json:"code"`
		Data struct {
			Data   *FetchCouponRequest `json:"data"`
			Mtgsig string              `json:"mtgsig"`
		} `json:"data"`
	}
)

type (
	FetchCouponRequest struct {
		CType         string `json:"cType"`
		FpPlatform    int    `json:"fpPlatform"`
		WxOpenId      string `json:"wxOpenId"`
		AppVersion    string `json:"appVersion"`
		MtFingerprint string `json:"mtFingerprint"`
	}
)



func SignDuration(secTime *utils.SecTimeData) []SignData {
	d := secTime.Sec.Sub(secTime.Mt) - 30*time.Second
	fmt.Println("在", d, "后获取签名")
	t := time.NewTimer(d)
	<-t.C
	return sign(fetchCouponUrl())
}

func sign(fetchCouponUrl string) []SignData {
	var data []SignData
	for name, cookie := range task.Cookies {
		var signResp []signResponse
		for i := 0; i < task.Num; i++ {
			resp := getSign(cookie, fetchCouponUrl)
			if resp == nil {
				continue
			}
			signResp = append(signResp, *resp)
		}

		data = append(data, SignData{
			Cookie: cookie,
			Name:   name,
			Sign:   signResp,
			Url:    fetchCouponUrl,
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

	httpReq, _ := http.NewRequest(http.MethodPost, task.SignUrl, bytes.NewReader(payload))
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

func fetchCouponUrl() string {
	val := url.Values{}
	val.Set("couponReferId", task.Id)
	val.Set("geoType", "2")
	val.Set("gdPageId", task.GdId)
	val.Set("version", "1")
	val.Set("utmSource", "wxshare")
	val.Set("pageId", task.PageId)
	val.Set("instanceId", task.InstanceId)
	val.Set("componentId", task.InstanceId)
	val.Set("actualLng", task.ActualLng)
	val.Set("actualLat", task.ActualLat)
	return fmt.Sprintf("%s?%s", task.FetchCouponUrl, val.Encode())
	//return fmt.Sprintf("https://promotion.waimai.meituan.com/lottery/limitcouponcomponent/fetchcoupon?%s", val.Encode())
}



func Fc(sd SignData) {
	wg.Add(task.Num)
	sdSignLen := len(sd.Sign)
	for i := 0; i < task.Num; i++ {
		if i >= sdSignLen {
			return
		}
		go sd.fcReq(i)
	}
	wg.Wait()
}

func (sd *SignData) fcReq(i int) {
	defer wg.Done()

	s := sd.Sign[i]

	payload, _ := json.Marshal(s.Data.Data)
	httpReq, err := http.NewRequest(http.MethodPost, sd.Url, bytes.NewReader(payload))
	if err != nil {
		fmt.Println("err", err)
		return
	}

	httpReq.Header = utils.BaseHeader(sd.Cookie, s.Data.Mtgsig)
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
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.0000"), fmt.Sprintf("====== data %s ========", sd.Name), string(data))
}
