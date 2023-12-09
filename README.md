# POSTER 小程序分享海报

## GO建纯后端绘制的小程序海报图片

生成带二维码的分享海报图片，常见的解决方式是在前端页面进行canvas绘制并获取快照图片，后端生成的案例寥寥无几。  
`poster`就是后端生成海报图片的一种简单实现，其适用于微信小程序海报分享，提示用户长按保存到系统相册后可以转发给朋友、分享到朋友圈。用户点击朋友圈图片，长按识别二维码后就可以跳转到小程序中。

### 海报案例

可一次生成多个海报image_url(多个,号隔开),如果需要生成背景透明的海报image_url背景使用png图片，image_url设为png

![尺寸设计](./design.jpg)
qrcode_width:170 海报上二维码宽度 为正方形
bottom:42 二维码距离背景图底部距离
right:42 二维码距离背景图右侧距离

![海报](https://s1.ax1x.com/2020/04/16/JkS2pn.jpg)

## 开始

### 运行

```bash
# 首次运行会提示填写配置，填写完后再次运行
go run main.go
# 调试ok后在后台运行
nohup go run main.go &
```

## 配置文件

```json
{
    "mapp": {
        "app_id": "小程序app_id，可选，除非命令行测试运行，否则请用自己的业务系统维护的access_token传值替代",
        "app_secret": "小程序app_secret，可选，除非命令行测试运行，否则请用自己的业务系统维护的access_token传值替代"
    },
    "output_dir": "海报保存目录，可选，默认：./output/",
    "listen_port": "http服务监听端口，可选，默认：2020",
    "fontfile_path": "字体路径，可选，默认: ./resources/font.ttc"
}
```

### HTTP调用

请求链接 http://127.0.0.1:2020/poster

```http
POST / HTTP/1.1
Host: 127.0.0.1:2020
Content-Type: application/x-www-form-urlencoded

access_token=&
scene=code=1&
page=pages/index/index&
width=280&
image_url=./sourceImages/test/poster1.jpg,./sourceImages/test/poster2.jpg&
env_version=trial&
output_dir=./output/test/&
bottom=42&
right=42&
qrcode_width=170&
image_type=jpg
```
**参数**  

`access_token` 是微信请求微信接口的access_token，需要在自己的业务系统中保持唯一，避免单独使用appid 和 appsecret生成，生成小程序码需要使用到这个参数。  

`scene`、`page`、`env_version`是小程序获取二维码的参数。

`width`是二维码的宽度，单位 px，最小 280px，最大 1280px。

`env_version`是要打开的小程序版本。正式版为 "release"，体验版为 "trial"，开发版为 "develop"。默认是正式版。

`image_url`是主图http链接或文件链接，多个,号隔开。

`output_dir`是保存文件目录，可选

`output_filename`是保存文件名称，可选，多个,号隔开

`bottom`是二维码距离背景图底部距离

`right`是二维码距离背景图右侧距离

`qrcode_width`是海报上二维码宽度

`image_type`是可选 默认 jpg 生成海报的图片类型 png jpg

**正常返回**  

```json
{
  "error": 0,
  "result": {
    "poster": "P50LAvUwq20yQCbC.jpg,hFQejimQ6YMMo7W6.jpg",
    "qrcode": "d7VWwZQRprO56Zun.jpg"
  }
}
```
输出目录中会生成海报图片和小程序码图片，将目录暴露到web目录中即可访问到（比如 `ln -s /your/output/dir /var/www/html/poster`）。

HTTP请求也支持JSON格式  

多张图片以及生成jpg海报

请求链接 http://127.0.0.1:2020/poster
```http
POST / HTTP/1.1
Host: 127.0.0.1:2020
Content-Type: application/json

{
    "access_token": "",
    "scene": "code=1",
    "page": "pages/index/index",
    "width": "280",
    "env_version": "trial",
    "image_url": "./sourceImages/test/poster1.jpg,./sourceImages/test/poster2.jpg",
    "output_dir": "./output/test/",
    "bottom": "42",
    "right": "42",
    "qrcode_width": "170",
    "image_type": "jpg"
}
```
返回
```json
{
  "error": 0,
  "result": {
    "poster": "cOGyIpfTuBlto18t.jpg,Bq6P9fNwfafY0NNZ.jpg",
    "qrcode": "d7VWwZQRprO56Zun.jpg"
  }
}
```
生成png海报
```http
POST / HTTP/1.1
Host: 127.0.0.1:2020
Content-Type: application/json

{
    "access_token": "",
    "scene": "code=1",
    "page": "pages/index/index",
    "width": "280",
    "env_version": "trial",
    "image_url": "./sourceImages/test/poster1.png",
    "output_filename": "test.png",
    "output_dir": "./output/test/",
    "bottom": "42",
    "right": "42",
    "qrcode_width": "170",
    "image_type": "png"
}
```
返回
```json
{
  "error": 0,
  "result": {
    "poster": "test.png",
    "qrcode": "sIaXCxws1IADl2uI.jpg"
  }
}
```
## 总结

这次更新取消强制填写appid和appsecret参数，改为由用户自己实现access_token的获取并传值，解决小程序生成二维码时，access_token容易被自己的业务系统挤掉导致失败的问题。  
移除了复杂的调用方式，使用最简单的http表单数据格式和json格式发送请求，简化到开箱即用。  
