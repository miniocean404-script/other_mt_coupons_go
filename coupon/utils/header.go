package utils

import "net/http"

func BaseHeader(cookie, mtgsig string) http.Header {
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
