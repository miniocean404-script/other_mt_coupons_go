//导入express
const express = require('express')
//创建web服务器
const app=express()
// 通过ap.listen进行服务器的配置，并启动服务器，接收两个配置参数，一个是对应的端口号，一个是启动成功的回调函数

// app.use(express.urlencoded());
app.use(express.json({limit:'1024mb'}));
app.post('/api/sign',(req,res)=>{
    const body = req.body
    let cookie = body.cookie
    let url = body.url
    if ( cookie === undefined || url === undefined) {
        res.send({
            "code":2000,
            "msg":"cookie或url不能为空"
        })
        return
    }
    (async (cookie,url) => {
        const data = await sign(cookie,url);
        res.send({data})

    })(cookie,url);

    // res.send({body})
})

app.listen(9588,()=>{
    console.log('服务器启动成功');
})



async function sign(cookie,url) {
    const H5guard = require('./mt.js');
    const data = {
        "cType": "mti", "fpPlatform": 3, "wxOpenId": "", "appVersion": ""
    };
    // const cookieStr = `userId=109669672; token=AgFmIZ71JVVOsboBJOfHTDScuNxxxxkie-Eo-3-xSY1LKxxxxxxqVKWI6w-cFJxkYC0nvr2kvQAAAADjFwAA188Rf9LCVaRyrH0e_CWM4626UUxKgSoNEJ8TsHfOIavuj5p1EKI9jRtcb13oDXqk; WEBDFPID=7u0xxy22201952x2z1xx8y040181624v812x965430297958672xuzvy-1997916869658-1682556867929AAQQYYSfd79fef3d01d5e9aadc18ccd4d0c95072564`;
// 安卓UA
    const userAgent = "Mozilla/5.0 (Linux; Android 9; MI 6 Build/PKQ1.190118.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/107.0.5304.141 Mobile Safari/537.36 XWEB/5075 MMWEBSDK/20230504 MMWEBID/5707 MicroMessenger/8.0.37.2380(0x28002598) WeChat/arm64 Weixin NetType/WIFI Language/zh_CN ABI/arm64";
    // const fullUrl = `https://promotion.waimai.meituan.com/lottery/limitcouponcomponent/fetchcoupon?couponReferId=F6CFF2A35BD94F49BDEE0CC6F7CF9FE4&geoType=2&gdPageId=306477&pageId=306004&version=1&utmSource=AppStore&utmCampaign=AgroupBgroupD0H0&instanceId=16620226080900.11717750606071209&componentId=16620226080900.11717750606071209`
    const h5guard = new H5guard(cookie, userAgent);
    const {mtgsig} = await h5guard.sign(url, data);

    return {'data':data,'mtgsig':mtgsig}
}