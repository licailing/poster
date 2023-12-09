package common

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"strings"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
	"github.com/anthonynsimon/bild/transform"
)

// Style 样式
type Style struct {
	ImageURL       string     //主图http链接或文件链接
	QRCodeURL      string     //二维码http链接或文件链接
	OutputFileName string     //保存文件名，可选
	OutputDIR      string     //保存文件目录，可选
	ImageType	   string  	  //生成海报的图片类型 png jpg
	Bottom	   	   int  	  //二维码距离背景图底部距离
	Right	   	   int  	  //二维码距离背景图右侧距离
	QrcodeWidth	   int  	  //海报上二维码宽度
}

func DrawPoster(s Style)(fileName string, err error) {
	arr := strings.Split(s.ImageURL, ",")
	nameArr := strings.Split(s.OutputFileName, ",")
	newNameArr := make([]string, len(arr), (cap(arr))*2)
	copy(newNameArr, nameArr)
	var images []string
	for index,item := range arr {
		fmt.Println(index, item)
		if item == "" {
			continue
		}
		s.OutputFileName = newNameArr[index]
		s.ImageURL = item
	 	name, err := DrawPosterImage(s)
		if err != nil {
			continue
		}
		images = append(images, name)
	}
	return strings.Join(images, ","), nil;
}

// DrawPoster 绘制海报
func DrawPosterImage(s Style) (fileName string, err error) {
	if s.ImageURL == "" || s.QRCodeURL == "" {
		return "", errors.New("样式参数缺失")
	}
	if s.OutputFileName == "" {
		s.OutputFileName = GetRandomName(16)
		s.OutputFileName = s.OutputFileName + "." + s.ImageType
	}

	// 获取网络图片并解码
	picRd, err := getResourceReader(s.ImageURL)
	if err != nil {
		log.Println("未找到图片资源")
		return "", err
	}
	pic, _ , err := image.Decode(picRd)
	if err != nil {
		log.Println("图片加载失败", err)
		return "", err
	}

	bounds := pic.Bounds();
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, pic, image.Point{}, draw.Src)

	// 获取二维码并绘制
	qr, err := getResourceReader(s.QRCodeURL)
	if err != nil {
		log.Println("二维码资源获取失败", err)
		return "", err
	}
	qrcode, _ , err := image.Decode(qr)
	qrcodeResized := transform.Resize(qrcode, s.QrcodeWidth, s.QrcodeWidth, transform.Linear)
	
	draw.Draw(rgba,
		image.Rectangle{
			image.Point{int(bounds.Dx()-s.Right-s.QrcodeWidth), int(bounds.Dy()-s.Bottom-s.QrcodeWidth)},
			rgba.Bounds().Max,
		},
		qrcodeResized,
		image.Point{0, 0},
		draw.Src)

	// 保存
sv:
	img := rgba.SubImage(rgba.Bounds())
	f, err := os.OpenFile(s.OutputDIR+s.OutputFileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Println("未找到保存目录", err)
		log.Println("即将创建保存目录")
		os.Mkdir(s.OutputDIR, 0666)
		goto sv
	}
	if(s.ImageType == "png") {
		err = png.Encode(f, img)
	} else {
		err = jpeg.Encode(f, img, nil)
	}
	if err != nil {
		fmt.Println("保存图片失败", err)
		return "", err
	}
	defer f.Close()
	return s.OutputFileName, nil
}

func getResourceReader(src string) (r *bytes.Reader, err error) {
	if src[0:4] == "http" {
		resp, err := http.Get(src)
		if err != nil {
			return r, err
		}
		defer resp.Body.Close()
		fileBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return r, err
		}
		r = bytes.NewReader(fileBytes)
	} else {
		fileBytes, err := ioutil.ReadFile(src)
		if err != nil {
			return nil, err
		}
		r = bytes.NewReader(fileBytes)
	}
	return r, nil
}

// GetRandomName 获取随机的文件名
func GetRandomName(length int) (name string) {
	dic := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	maxL := len([]byte(dic))
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		name = name + string(dic[rand.Intn(maxL)])
	}
	return name
}
