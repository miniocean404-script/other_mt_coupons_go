package task

const (
	Num            = 5
	Early          = 100 // 毫秒
	SecAt          = "17:52:00"
	Id             = "00B223429B424F7A910C0D4885957E99"
	GdId           = "379397"
	PageId         = "378931"
	InstanceId     = "16618616100670.97030510386642830"
	SignUrl        = "http://127.0.0.1:9588/api/sign" // lsof -i :9588 查看进程占用端口
	CouponInfoUrl  = "https://promotion.waimai.meituan.com/lottery/limitcouponcomponent/info"
	FetchCouponUrl = "https://promotion.waimai.meituan.com/lottery/limitcouponcomponent/fetchcoupon"

	// 经纬度
	ActualLng = "117.23344"
	ActualLat = "31.82658"
)

var (
	Cookies = map[string]string{
		// "6886": "userId=906082484;token=AgFnIc_G9WnCR-CCPYyQJl_CiHAYyc_dOiHo8yxhnGdVDrLTGKLyOtSPoBDUCp97sNmMiDI96fUcEgAAAADHGQAABg8OohMbfRiPzBgxgtTJdaVYnvlTZJb6BE4et1c4G85ze68Liv6XrsClwWdLhamC",
		// "6408": "userId=3170637747;token=AgGJI9GQYEiBsRxSPWsIxdZK9RZcTvDSDPOxRqLWMn8PcX0AE_lrSGWwm573RMrwydTL6AvM_cqh5wAAAACpGQAAaBsGenSWPG4G7yk8Mn2Tv5VeavdF-L57DDaWJKpktMt7WCKO8oceddBorwIRLYG-",
		"6408": "userId=3170637747;token=DPOxRqLWMn8PcX0AE_lrSGWwm573RMrwydTL6AvM_cqh5wAAAACpGQAAaBsGenSWPG4G7yk8Mn2Tv5VeavdF-L57DDaWJKpktMt7WCKO8oceddBorwIRLYG-",
	}
)
