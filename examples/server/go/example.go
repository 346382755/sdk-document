package main

import (
	"./server"
	"fmt"
)

var (
	//公钥
	publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA/VweA5KME/PR7QIwe+Bo
Wf+yM5tRVpaXWOZC7S4SeLT5zyd1gNzKjLHCxAIhsxYUnXHRCdsC+cnKTVfCBuew
v7N2kZCKF+/gMiqSdfiJo3XE7lmrxpIKO6YnWt7Itq/VvMJoTO7g0KkjF/irzI+O
KUj5DSZORHwulA6OXNjUJbGeMdIGX7VHgbk7dv8oqx+FE0bAQ4APtOiSs5agBRdg
9De92tURcj15jztYtPvaRVn5O8ozTEui4Kh2Cmf1fPFKbv5yQyNhHaqMdI2tGPRE
s78wGIiHT1yOmCMeKHCVgElwwgmFnmxKYsBD9XZ9GM6wzt/95M53jh/aNp/+9Y1m
0QIDAQAB
-----END PUBLIC KEY-----`
	//示例login_key
	loginKey = "11111111"
)

func main()  {

	//登录验证
	loginEntity := `{"openid":"1-123123","account":"test","time":1482313093}`
	loginSign := "m7JTn/y3IpB084vyeqoR9ysZ5/GeowcsO3KDrDsaa8Fof2Xjq4gKdk5eDLh0nTRxZslfbet5AWs+p1M0rtF8Jan8T2VxCW/czoCkrj4o/xnYtb3wdenCXAT7LUoydjTmd+cf0I9kw0DjkPDurrT9kNqxzw6dq1A6EFZAts4f0/H5+7kn81rzPq1RkWOM6OGm8R2D2WW/jHfqZDo1mvfZvIzkA/F0M62z2VsUK821BUkgoT23dByCdeqgx4hauJTHlnvg2/MplDYMDONHai6gTFWb4FQgY2wJq1BIWCyXgDuid6n9Ck0m6paJEfOffbK7BgwH3ssaF7xYGD2sjEuCFw=="
	userInfo, err := server.LoginVerify(publicKey, loginEntity, loginSign, 36000)
	if err != nil {
		fmt.Printf("签名验证失败或登录已过期: %s", err.Error())
	} else {
		fmt.Printf("签名验证到成功, openid: %s", userInfo["openid"])
	}

	//============================================

	//支付验证
	//接收到的http post内容
	paymentPost := map[string]string{
		"account":"abcd",
		"amount":"6.00",
		"black_desc":"",
		"channel":"1",
		"currency":"",
		"extra":"123",
		"game_id":"GMG001",
		"is_black":"0",
		"is_cancel":"0",
		"is_recovery":"0",
		"is_test":"0",
		"memo":"",
		"openid":"1-1234",
		"order_id":"1399633295037630",
		"order_type":"0",
		"original_purchase_orderid":"",
		"product_id":"HWDPID0006",
		"time":"1404975144",
		"transaction_id":"1000000110081354",
		"version":"8.0",
		"zone_id":"1",
		"sign":"MT/Wp3jz8l8HV62/OrpXdDIdEF9BGw8zCtwu8VZNaxCi52QXCC7RO2Q5yLRsQy8SmZMk0h0hsxUOXSZ8zZysUL+r/ATVknUs76zveaHROAfqhwo4sBqNA3jvnkCtbV/x9vF46vmNB7tDJcG6UDrJM5a/9Bn3Xq0ctlfSJpqJskr4dqyIlvCbQyi2YtnW2f5aNWATyNUntWtzakXYxZ7mPVT7hhspREo/3v/YcygmUHJJsJ2NnqkyWH7kV3nTTtvTQirjbzgeqnFd+KZBBHuBxOXDZInNLcaTnOQLskRddt0sYJ9/L8aDx9AvZtMe8kgQJhlNq13JOLFbCoMAEHD4Gg==",
	}
	//示例游戏商品数据
	gameOrder := map[string]string{
		"product_id":"HWDPID0006",
		"amount":"6.00",
	}
	//验证签名一致
	if server.PaymentVerify(publicKey, paymentPost) {
		//游戏服务端根据order_id判断订单是否已经处理过，游戏服务端自行实现，并响应`{"code":0,"msg":"成功"}`

		//测试订单，如果游戏不允许测试订单通知，则返回`{"code":2,"msg":"不允许测试订单通知"}`
		if paymentPost["is_test"] == "1" {
			fmt.Println(`{"code":2,"msg":"不允许测试订单通知"}`)
			return
		}

		//风控拦截订单，游戏不能发货，游戏接收到此类订单通知时，可以通过邮件或其他方式告诉用户充值不到账原因和如何解决问题;
		if paymentPost["is_black"] == "1" {
			fmt.Println(`{"code":0,"msg":"风控拦截订单通知"}`)
			return
		}

		//恢复购买订单，仅苹果、谷歌支付非消耗型商品才能恢复购买操作，游戏接收到此类订单通知时，请特殊标识，以与正常支付订单区分
		if paymentPost["is_recovery"] == "1" {
			fmt.Println(`{"code":0,"msg":"恢复购买订单通知"}`)
			return
		}

		//取消购买订单，仅自动续订订阅型商品才有取消购买操作，游戏接收到此类订单通知时，请撤消之前的订单购买状态，并标识原订阅商品不可用
		if paymentPost["is_cancel"] == "1" {
			fmt.Println(`{"code":0,"msg":"取消购买订单通知"}`)
			return
		}

		//验证product_id和金额是否匹配
		if gameOrder["product_id"] != paymentPost["product_id"] {
			fmt.Println(`{"code":2,"msg":"商品不存在"}`)
			return
		}
		//验证金额是否正确
		if gameOrder["amount"] != paymentPost["amount"] {
			fmt.Println(`{"code":2,"msg":"金额错误"}`)
			return
		}

		//游戏发货

		//发货成功，返回：`{"code":0,"msg":"成功"}`;

		//游戏发货失败，需要再次通知，返回：`{"code":1,"msg":"数据库错误"}`;

		//成功
		fmt.Println(`{"code":0,"msg":"成功"}`)
	} else {
		//签名验证失败返回
		fmt.Println(`{"code":2, "msg":"签名验证失败"}`)
	}

	//============================================

	//在线人数实时统计
	onlineData := map[string]interface{}{
		"game_id" : 1000,
		"zone_id" : 1,
		"number": 99,
	}
	if server.GameOnline(loginKey, onlineData) {
		fmt.Println("上报在线人数成功")
	} else {
		fmt.Println("上报在线人数失败")
	}

	//============================================

	//游戏登录数据上报
	loginData := map[string]interface{}{
		"game_id" : 1000,
		"channel_id" : 1,
		"ip" : "192.168.100.1",
		"openid" : "1-123456",
		"device_udid" : "860529036567763",
		"device_type": "vivo Y51",
		"os": "android",
	}
	if server.LoginLogs(loginKey, loginData) {
		fmt.Println("上报游戏登录数据成功")
	} else {
		fmt.Println("上报游戏登录数据失败")
	}

	//============================================

	//推送消息
	pushData := map[string]interface{}{
		"game_id": 1000,
		"channel_id": 1,
		"environment": "product",
		"audience_type": "account",
		"audience_list": "1-123456",
		"message_type": "notify",
		"title": "test",
		"content": "test",
	}
	if server.Push(loginKey, pushData) {
		fmt.Println("推送消息成功")
	} else {
		fmt.Println("推送消息失败")
	}

}
