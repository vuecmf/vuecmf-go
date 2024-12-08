//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package helper

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/golang/freetype"
	"github.com/vuecmf/vuecmf-go/v3/app"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
)

// 水印的位置
const (
	TopLeft int = iota
	TopRight
	BottomLeft
	BottomRight
	Center
)

type img struct {
	FileObj  *os.File //文件实例
	FileName string   //文件名
	FileExt  string   //文件扩展名
	Err      error    //异常信息
}

var imgInstance *img

// Img 获取图像处理实例
func Img() *img {
	if imgInstance == nil {
		imgInstance = &img{}
	}
	return imgInstance
}

// Load 根据文件名加载文件
//
//	参数：
//	fileName 文件名
func (im *img) Load(fileName string) *img {
	im.FileName = fileName
	im.FileExt = GetFileExt(fileName)
	f, err := os.Open(fileName)
	im.FileObj = f
	if err != nil {
		im.Err = errors.New("图像文件读取失败！" + err.Error())
	}

	return im
}

// LoadFile 直接加载文件流
//
//	参数：
//	file 文件流
//	fileName 文件名
func (im *img) LoadFile(file *os.File, fileName string) *img {
	im.FileName = fileName
	im.FileExt = GetFileExt(fileName)
	im.FileObj = file
	return im
}

// GetImage 获取image实例及图片类型
func (im *img) GetImage() (image.Image, string, error) {
	defer func(f *os.File) {
		_ = f.Close()
	}(im.FileObj)

	var imgObj image.Image
	var err error

	switch im.FileExt {
	case "gif":
		imgObj, err = gif.Decode(im.FileObj)
	case "jpeg":
		fallthrough
	case "jpg":
		imgObj, err = jpeg.Decode(im.FileObj)
	case "png":
		imgObj, err = png.Decode(im.FileObj)
	default:
		err = errors.New("未知的图像类型: " + im.FileExt)
	}

	if err != nil {
		err = errors.New("图像文件解析错误：" + err.Error())
	}

	return imgObj, im.FileExt, err
}

// SaveImage 保存图像文件
//
//	参数：
//	outImg 需要保存的image实例
//	saveFileName 保存的文件名
func (im *img) SaveImage(outImg image.Image, saveFileName string) error {
	fileExt := GetFileExt(saveFileName)
	f, _ := os.Create(saveFileName)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	buf := bufio.NewWriter(f)
	var err error

	switch fileExt {
	case "gif":
		err = gif.Encode(buf, outImg, &gif.Options{NumColors: 256})
	case "jpeg":
		fallthrough
	case "jpg":
		err = jpeg.Encode(buf, outImg, &jpeg.Options{Quality: 100})
	case "png":
		err = png.Encode(buf, outImg)
	default:
		err = errors.New("未知的图片类型: " + fileExt)
	}

	if err != nil {
		return err
	}

	err = buf.Flush()
	return err
}

// Resize 图片绽放
//
//	参数：
//		newFileName 保存新的文件名
//		width		缩放后的宽度
//		height  	缩放后的高度
//		keepRatio	是否保持等比例缩放
//		fill		填充的背景颜色 0 - 255 （R、G、B）的值共一个数值， 0 = 透明背景， 255 = 白色背景
//		centerAlign	是否以图片的中心来进行等比缩放
//		crop		是否裁切
func (im *img) Resize(newFileName string, width int, height int, keepRatio bool, fill int, centerAlign bool, crop bool) error {
	imgRes, _, err := im.GetImage()
	if err != nil {
		return err
	}

	outImg := image.NewRGBA(image.Rect(0, 0, width, height))

	if !keepRatio {
		//非等比缩放
		draw.BiLinear.Scale(outImg, outImg.Bounds(), imgRes, imgRes.Bounds(), draw.Over, nil)
	} else {
		//填充背景色
		if fill != 0 {
			fillColor := color.RGBA{R: uint8(fill), G: uint8(fill), B: uint8(fill), A: 255}
			draw.Draw(outImg, outImg.Bounds(), &image.Uniform{C: fillColor}, image.Point{}, draw.Src)
		}

		srcWidth := imgRes.Bounds().Dx()
		srcHeight := imgRes.Bounds().Dy()
		var dst image.Rectangle

		if crop == true {
			//计算需要裁切掉的宽度和高度
			cropWidth := 0
			cropHeight := 0
			var srcBounds image.Rectangle
			var tmpBounds image.Rectangle

			if srcWidth > width && srcHeight <= height {
				cropWidth = (srcWidth - width) / 2
				padding := (height - srcHeight) / 2
				srcBounds = image.Rect(0, padding, srcWidth, srcHeight+padding)
				tmpBounds = image.Rect(0, 0, srcWidth, height)
				srcHeight = height
			} else if srcWidth <= width && srcHeight > height {
				cropHeight = (srcHeight - height) / 2
				padding := (width - srcWidth) / 2
				srcBounds = image.Rect(padding, 0, srcWidth+padding, srcHeight)
				tmpBounds = image.Rect(0, 0, width, srcHeight)
				srcWidth = width
			} else if srcWidth <= width && srcHeight <= height {
				wPadding := (width - srcWidth) / 2
				hPadding := (height - srcHeight) / 2
				srcBounds = image.Rect(wPadding, hPadding, srcWidth+wPadding, srcHeight+hPadding)
				tmpBounds = image.Rect(0, 0, width, height)
				srcWidth = width
				srcHeight = height
			} else if srcWidth > width && srcHeight > height {
				wRatio := float64(srcWidth) / float64(width)
				hRatio := float64(srcHeight) / float64(height)
				if wRatio > hRatio {
					cropWidth = (srcWidth - int(float64(srcHeight)*(float64(width)/float64(height)))) / 2
				} else if wRatio < hRatio {
					cropHeight = (srcHeight - int(float64(srcWidth)*(float64(height)/float64(width)))) / 2
				}
				srcBounds = imgRes.Bounds()
				tmpBounds = srcBounds
			}

			//裁切图片
			tmp := image.NewNRGBA(tmpBounds)
			draw.Draw(tmp, srcBounds, imgRes, imgRes.Bounds().Min, draw.Over)
			imgRes = tmp.SubImage(image.Rect(cropWidth, cropHeight, srcWidth-cropWidth, srcHeight-cropHeight))

			dst = image.Rect(0, 0, width, height)

		} else {
			//计算缩放后大小
			if width*srcHeight < height*srcWidth {
				ratio := float64(width) / float64(srcWidth)
				targetHeight := int(float64(srcHeight) * ratio)
				padding := 0
				if centerAlign {
					padding = (height - targetHeight) / 2
				}
				dst = image.Rect(0, padding, width, padding+targetHeight)
			} else {
				ratio := float64(height) / float64(srcHeight)
				targetWidth := int(float64(srcWidth) * ratio)
				padding := 0
				if centerAlign {
					padding = (width - targetWidth) / 2
				}
				dst = image.Rect(padding, 0, padding+targetWidth, height)
			}
		}

		//缩放图片
		draw.ApproxBiLinear.Scale(outImg, dst.Bounds(), imgRes, imgRes.Bounds(), draw.Over, nil)
	}

	return im.SaveImage(outImg, newFileName)
}

// FontWater 给图片添加文字水印
func (im *img) FontWater(typeface []app.FontInfo) error {
	defer func(f *os.File) {
		_ = f.Close()
	}(im.FileObj)

	//新的加了水印图片覆盖原来文件
	if im.FileExt == "gif" {
		return im.gifFontWater(im.FileObj, im.FileName, typeface)
	} else {
		return im.staticFontWater(im.FileObj, im.FileName, im.FileExt, typeface)
	}
}

// gif图片水印
func (im *img) gifFontWater(imgFile *os.File, newImage string, typeface []app.FontInfo) error {
	var err error
	gifImg2, _ := gif.DecodeAll(imgFile)
	gifs := make([]*image.Paletted, 0)
	x0 := 0
	y0 := 0
	yuan := 0
	for k, gifImg := range gifImg2.Image {
		rgbImg := image.NewNRGBA(gifImg.Bounds())
		if k == 0 {
			x0 = rgbImg.Bounds().Dx()
			y0 = rgbImg.Bounds().Dy()
		}
		fmt.Printf("%v, %v\n", rgbImg.Bounds().Dx(), rgbImg.Bounds().Dy())
		if k == 0 && gifImg2.Image[k+1].Bounds().Dx() > x0 && gifImg2.Image[k+1].Bounds().Dy() > y0 {
			yuan = 1
			break
		}
		if x0 == rgbImg.Bounds().Dx() && y0 == rgbImg.Bounds().Dy() {
			for y := 0; y < rgbImg.Bounds().Dy(); y++ {
				for x := 0; x < rgbImg.Bounds().Dx(); x++ {
					rgbImg.Set(x, y, gifImg.At(x, y))
				}
			}
			rgbImg, err = im.common(rgbImg, typeface) //添加文字水印
			if err != nil {
				break
			}
			//定义一个新的图片调色板img.Bounds()：使用原图的颜色域，gifImg.Palette：使用原图的调色板
			p1 := image.NewPaletted(gifImg.Bounds(), gifImg.Palette)
			//把绘制过文字的图片添加到新的图片调色板上
			//dst：绘图的背景图
			//r：背景图的绘图区域
			//src：要绘制的图
			//sp：要绘制图src的开始点
			//op：组合方式
			draw.Draw(p1, gifImg.Bounds(), rgbImg, image.Point{}, draw.Src)
			//把添加过文字的新调色板放入调色板slice
			gifs = append(gifs, p1)
		} else {
			gifs = append(gifs, gifImg)
		}
	}
	if yuan == 1 {
		return errors.New("gif: image block is out of bounds")
	} else {
		if err != nil {
			return err
		}
		//保存到新文件中
		newFile, err := os.Create(newImage)
		if err != nil {
			return err
		}
		defer func(newFile *os.File) {
			_ = newFile.Close()
		}(newFile)

		g1 := &gif.GIF{
			Image:     gifs,
			Delay:     gifImg2.Delay,
			LoopCount: gifImg2.LoopCount,
		}
		err = gif.EncodeAll(newFile, g1)
		return err
	}
}

// png,jpeg图片水印
func (im *img) staticFontWater(imgFile *os.File, newImage, status string, typeface []app.FontInfo) (err error) {
	var staticImg image.Image
	if status == "png" {
		staticImg, _ = png.Decode(imgFile)
	} else {
		staticImg, _ = jpeg.Decode(imgFile)
	}
	//新建一个RGB图像
	rgbImg := image.NewNRGBA(staticImg.Bounds())
	for y := 0; y < rgbImg.Bounds().Dy(); y++ {
		for x := 0; x < rgbImg.Bounds().Dx(); x++ {
			//设置每个像素点颜色，颜色来自要添加水印的图像文件（相当于将图像复制一份为 rgbImg）
			rgbImg.Set(x, y, staticImg.At(x, y))
		}
	}

	rgbImg, err = im.common(rgbImg, typeface) //开始添加文字水印
	if err != nil {
		return err
	}

	//保存到新文件中
	newFile, err := os.Create(newImage)
	if err != nil {
		return err
	}
	defer func(newFile *os.File) {
		_ = newFile.Close()
	}(newFile)

	if status == "png" {
		err = png.Encode(newFile, rgbImg)
	} else {
		err = jpeg.Encode(newFile, rgbImg, &jpeg.Options{Quality: 100})
	}
	return err
}

// 添加文字水印函数
func (im *img) common(rgbImg *image.NRGBA, typeface []app.FontInfo) (*image.NRGBA, error) {
	//拷贝一个字体文件到运行目录
	fontBytes, err := os.ReadFile(app.Config().Water.WaterFont)
	if err != nil {
		return nil, errors.New("字体文件打开失败！" + err.Error())
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, errors.New("字体文件解析失败！" + err.Error())
	}
	errNum := 1
Loop:
	for _, t := range typeface {
		info := t.Message          //水印文字内容
		f := freetype.NewContext() //创建字体上下文
		f.SetDPI(72)               //设置字体分辨率（屏幕每英寸的分辨率）
		f.SetFont(font)            //设置字体
		f.SetFontSize(t.Size)      //设置字号（字体大小）
		f.SetClip(rgbImg.Bounds()) //设置文字可绘制区域的尺寸
		f.SetDst(rgbImg)           //设置最终输出的图片
		//设置字体颜色
		f.SetSrc(image.NewUniform(color.RGBA{R: t.R, G: t.G, B: t.B, A: t.A}))

		//文本左下角
		writeX := 0 //水印文字起始x位置
		writeY := 0 //水印文字起始y位置

		//f.PointToFixed 表示将字号大小转成像素值
		fontWidth := len(info) * int(f.PointToFixed(t.Size)>>7) //文字宽度
		fontHeight := int(t.Size / 72.0 * 50)                   //文字高度

		switch t.Position {
		case TopLeft:
			writeX = t.Dx
			writeY = t.Dy + fontHeight
		case TopRight:
			writeX = rgbImg.Bounds().Dx() - fontWidth - t.Dx
			writeY = t.Dy + fontHeight
		case BottomLeft:
			writeX = t.Dx
			writeY = rgbImg.Bounds().Dy() - t.Dy
		case BottomRight:
			writeX = rgbImg.Bounds().Dx() - fontWidth - t.Dx
			writeY = rgbImg.Bounds().Dy() - t.Dy
		case Center:
			writeX = (rgbImg.Bounds().Dx() - fontWidth) / 2
			writeY = (rgbImg.Bounds().Dy()-fontHeight)/2 + fontHeight
		default:
			errNum = 0
			break Loop
		}

		pt := freetype.Pt(writeX, writeY) //设置字体的位置
		_, err = f.DrawString(info, pt)
		if err != nil {
			break
		}
	}
	if errNum == 0 {
		err = errors.New("不支持的位置值（只能为0到4的整数，0=左上角， 1=右上角，2=左下角，3=右下角，4=中间）")
	}
	return rgbImg, err
}
