媒体合作广告点击调用接口说明
=========================

 本广告点击统计对接分为两种方式:

* Server to server
* 客户端302跳转

## Server to server 接口说明

### 点击接口地址

`http://ad.stat.mztgame.com/ad/s2s-hits`
 
### 接口参数
 
| 参数名 | 说明 | 必填 |
|------|------|------|
| ad_id | 对接广告ID(巨人提供) | 是 |
| game_id | 巨人游戏应用 id(巨人提供) | 是|
| idfa | IOS IDFA | 否 |
| imei | 设备IMEI(Android)  | 否 |
| mac | IOS 或Android Mac | 否 |
| ua | 用户若是通过网页点击的， 为浏览器的 User-Agent 信息 | 否 |
| ip | 用户IP | 否 |
| callback |  广告商追踪使用url，启动时会回调 |  否 |
 
### 对接地址范例

```
http://ad.stat.mztgame.com/ad/s2s-hits?ad_id=2000001&game_id=5012&idfa=(idfa)&ip=1.1.1.1&ua=(ua)
```

### 响应说明

HTTP 响应状态码为 204，无响应内容

## 客户端302跳转接口说明

### 点击接口地址

`http://ad.stat.mztgame.com/ad/client-hits`

### 接口参数
 
| 参数名 | 说明 | 必填 |
|------|------|------|
| ad_id | 对接广告ID(巨人提供) | 是 |
| game_id | 巨人游戏应用 id(巨人提供) | 是|
| idfa | IOS IDFA | 否 |
| imei | 设备IMEI(Android)  | 否 |
| mac | IOS 或Android Mac | 否 |
| default_url | 默认跳转URL|  是 |
| ios_url | 程序识别 IOS UA 跳转URL|  否 |
| android_url | 程序识别 Android UA跳转URL |  否 |
| callback |  广告商追踪使用url，启动时会回调 |  否 |

> 注：以上参数 `default_url`, `ios_url`, `android_url` 限制域名：`*.[m]ztgame.com`, `itunes.apple.com`,`*.[m]ztgame.com.cn`
 
### 对接地址范例

* 客户端跳转方式

```
http://ad.stat.mztgame.com/ad/client-hits?ad_id=2000001&game_id=5012&idfa=(idfa)&default_url=https://itunes.apple.com&ios_url=&android_url=
```
 
### 响应说明

HTTP 响应状态码为 302，跳转地址是 `default_url`, `ios_url`, `android_url` 参数传入值。
