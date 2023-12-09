package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// ServeHTTP 简单http接口
func ServeHTTP() {
	http.HandleFunc("/poster", postHandle)
	http.ListenAndServe(":"+C.ListenPort, nil)
}

type httpReq struct {
	AccessToken string `form:"access_token" json:"access_token"`
	Scene       string `form:"scene" json:"scene"`
	Page        string `form:"page" json:"page"`
	Width       string `form:"width" json:"width"` //二维码的宽度，单位 px，最小 280px，最大 1280px
	EnvVersion  string `form:"env_version" json:"env_version"` //要打开的小程序版本。正式版为 "release"，体验版为 "trial"，开发版为 "develop"。默认是正式版。
	ImageURL    string `form:"image_url" json:"image_url"` //主图http链接或文件链接，多个,号隔开
	QrcodeURL   string `form:"qrcode_url" json:"qrcode_url"` //可选 二维码http链接或文件链接
	OutputDIR   string `form:"output_dir" json:"output_dir"`     //保存文件目录，可选
	OutputFileName   string `form:"output_filename" json:"output_filename"`     //保存文件名称，可选，多个,号隔开
	Bottom		string `form:"bottom" json:"bottom"` //二维码距离背景图底部距离
	Right		string `form:"right" json:"right"` //二维码距离背景图右侧距离
	QrcodeWidth string `form:"qrcode_width" json:"qrcode_width"` //海报上二维码宽度
	ImageType 	string `form:"image_type" json:"image_type"` //可选 默认 jpg 生成海报的图片类型 png jpg
}

type httpResp struct {
	Poster string `json:"poster"` //海报图文件名或链接
	QRCode string `json:"qrcode"` //二维码文件名或链接
}

func postHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Println("请使用POST方法")
		writeResponse(w, map[string]interface{}{"error": 1, "message": "请使用POST方法"})
		return
	}
	var rq httpReq
	switch r.Header.Get("Content-Type") {
	case "application/json":
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("请求参数异常", err)
			writeResponse(w, map[string]interface{}{"error": 1, "request": rq, "message": "请求参数异常：缺少参数"})
			return
		}
		err = json.Unmarshal(data, &rq)
		if err != nil {
			log.Println("请求参数异常", err)
			writeResponse(w, map[string]interface{}{"error": 1, "request": rq, "message": "请求参数异常:不是有效的JSON格式"})
			return
		}
	default:
		rq.AccessToken = r.PostFormValue("access_token")
		rq.Scene = r.PostFormValue("scene")
		rq.Page = r.PostFormValue("page")
		rq.Width = r.PostFormValue("width")
		rq.EnvVersion = r.PostFormValue("env_version")
		rq.ImageURL = r.PostFormValue("image_url")
		rq.OutputDIR = r.PostFormValue("output_dir")
		rq.OutputFileName = r.PostFormValue("output_filename")
		rq.Bottom = r.PostFormValue("bottom")
		rq.Right = r.PostFormValue("right")
		rq.QrcodeWidth = r.PostFormValue("qrcode_width")
		rq.ImageType = r.PostFormValue("image_type")
		rq.QrcodeURL = r.PostFormValue("qrcode_url")
	}

	if  rq.Scene == "" || rq.Page == "" {
		log.Println("POST参数无效")
		writeResponse(w, map[string]interface{}{"error": 1, "request": rq, "message": "POST参数无效"})
		return
	}

	// 获取小程序码保存到本地
	_width, err := strconv.Atoi(rq.Width)
	if err != nil {
		log.Println("width 必须是有效的数字")
		writeResponse(w, map[string]interface{}{"error": 1, "request": rq, "message": "width 必须是有效的数字"})
		return
	}

	_bottom, err := strconv.Atoi(rq.Bottom)
	if err != nil {
		log.Println("bottom 必须是有效的数字")
		writeResponse(w, map[string]interface{}{"error": 1, "request": rq, "message": "bottom 必须是有效的数字"})
		return
	}

	_right, err := strconv.Atoi(rq.Right)
	if err != nil {
		log.Println("right 必须是有效的数字")
		writeResponse(w, map[string]interface{}{"error": 1, "request": rq, "message": "right 必须是有效的数字"})
		return
	}

	_qrcodeWidth, err := strconv.Atoi(rq.QrcodeWidth)
	if err != nil {
		log.Println("qrcode_width 必须是有效的数字")
		writeResponse(w, map[string]interface{}{"error": 1, "request": rq, "message": "qrcode_width 必须是有效的数字"})
		return
	}

	if(rq.OutputDIR == "") {
		rq.OutputDIR = C.OutputDIR
	}
	var qrcodeName string
	if rq.QrcodeURL == "" {
		qrcodeImageName, err := RequestQRCode(QRCodeReq{
			Scene: rq.Scene,
			Page:  rq.Page,
			Width: _width,
			EnvVersion: rq.EnvVersion,
		}, "", rq.AccessToken, rq.OutputDIR)
	
		if err != nil {
			log.Println("获取小程序码失败", err)
			writeResponse(w, map[string]interface{}{"error": 1, "request": rq, "message": "获取小程序码失败"})
			return
		}
		qrcodeName = qrcodeImageName
		// 使用生成的二维码来生成海报
		rq.QrcodeURL = rq.OutputDIR + qrcodeImageName
	} else {
		qrcodeName = rq.QrcodeURL
	}

	if rq.ImageType == "" {
		rq.ImageType = "jpg"
	}

	posterName, err := DrawPoster(Style{
		ImageURL:       rq.ImageURL,
		QRCodeURL:      rq.QrcodeURL,
		OutputFileName: rq.OutputFileName,
		OutputDIR:      rq.OutputDIR,
		ImageType: 		rq.ImageType,
		QrcodeWidth: 	_qrcodeWidth,
		Bottom:			_bottom,
		Right:			_right,
	})
	if err != nil {
		log.Println("海报生成失败", err)
		writeResponse(w, map[string]interface{}{"error": 1, "request": rq, "message": "海报生成失败"})
		return
	}
	writeResponse(w, map[string]interface{}{"error": 0, "result": map[string]string{"poster": posterName, "qrcode": qrcodeName}})
}

func writeResponse(w http.ResponseWriter, d map[string]interface{}) {
	data, _ := json.MarshalIndent(d, "", "  ")
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}
