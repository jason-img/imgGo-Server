package handle

import (
	"bytes"
	"errors"
	"github.com/chai2010/webp"
	"log"
	"os"
	//gifToWebp "github.com/sizeofint/gif-to-webp"
	"golang.org/x/image/bmp"
	"image"
	"image/jpeg"
	"image/png"
)

const (
	ImageWebp = iota
	ImageJpeg
	ImagePng
	ImageBmp
	ImageGif
	ImageTif
)

type Image struct {
	FilePath  string
	Data      []byte
	ImageType int
	Ext       string
	Width     int
	Height    int
}

// Open /**
func (i *Image) Open(filePath string, extName string) (err error) {
	i.FilePath = filePath
	i.Data, _ = os.ReadFile(filePath)
	i.Ext = extName

	switch i.Ext {
	case "jpg", "jpeg":
		i.ImageType = ImageJpeg
	case "png":
		i.ImageType = ImagePng
	case "bmp":
		i.ImageType = ImageBmp
	case "gif":
		i.ImageType = ImageGif
	case "tif", "tiff":
		i.ImageType = ImageTif
	case "webp":
		i.ImageType = ImageWebp
	}

	reader := bytes.NewReader(i.Data)
	img, _, err := image.Decode(reader)
	if err != nil {
		return err
	}
	b := img.Bounds()
	i.Width = b.Max.X
	i.Height = b.Max.Y
	return nil
}

// ToWebP 转换为WEBP
func (i *Image) ToWebP(quality float32) (out []byte, err error) {
	var img image.Image
	reader := bytes.NewReader(i.Data)
	lossLess := false //是否无损压缩
	Exact := false    //透明部分消失
	switch i.ImageType {
	case ImageJpeg:
		img, _ = jpeg.Decode(reader)
		break
	case ImagePng:
		img, _ = png.Decode(reader)
		lossLess = true
		Exact = true
		break
	case ImageBmp:
		img, _ = bmp.Decode(reader)
		break
	//case ImageGif:
	//	return i.gitToWebP(i.Data, quality)
	case ImageWebp:
		return i.Data, nil
	default:
		log.Printf("暂不支持将%s转换为webp\n", i.Ext)
	}
	if img == nil {
		msg := "image file " + i.FilePath + " is corrupted or not supported"
		err = errors.New(msg)
		return nil, err
	}
	var buf bytes.Buffer
	if err = webp.Encode(&buf, img, &webp.Options{Lossless: lossLess, Exact: Exact, Quality: quality}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// gitToWebP git 转webP
//func (i *Image) gitToWebP(gifBin []byte, quality float32) (webPBin []byte, err error) {
//	converter := gifToWebp.NewConverter()
//	converter.LoopCompatibility = false
//	//0 有损压缩  1无损压缩
//	converter.WebPConfig.SetLossless(0)
//	//压缩速度  0-6  0最快 6质量最好
//	converter.WebPConfig.SetMethod(0)
//	converter.WebPConfig.SetQuality(quality)
//	//搞不懂什么意思,例子是这样用的
//	converter.WebPAnimEncoderOptions.SetKmin(9)
//	converter.WebPAnimEncoderOptions.SetKmax(17)
//
//	return converter.Convert(gifBin)
//}
