package helper

import (
	"errors"
	"fmt"
	"github.com/golang/freetype"
	"github.com/vuecmf/vuecmf-go/app"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
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
}

var imgInstance *img

func Img() *img {
	if imgInstance == nil {
		imgInstance = &img{}
	}
	return imgInstance
}

//
func (w *img) Resize(fileName string) error {
	imgFile, err := os.Open(fileName)
	if err != nil {
		return errors.New("打开文件失败！" + err.Error())
	}
	im, _, err := image.Decode(imgFile)
	if err != nil {
		return errors.New("解析图片失败！" + err.Error())
	}

	rgbImg := im.(*image.YCbCr)
	subImg := rgbImg.SubImage(image.Rect(0,0, 300,300)).(*image.YCbCr)

	f, err := os.Create("uploads/test.jpg")     //创建文件
	defer f.Close()                   //关闭文件
	jpeg.Encode(f, subImg, nil)       //写入文件
	return err
}

// FontWater 给图片添加文字水印
func (w *img) FontWater(fileName string, typeface []app.FontInfo) error {
	imgFile, err := os.Open(fileName)
	if err != nil {
		return errors.New("打开文件失败！" + err.Error())
	}

	defer func(imgFile *os.File) {
		_ = imgFile.Close()
	}(imgFile)

	_, str, err := image.DecodeConfig(imgFile)
	if err != nil {
		return err
	}

	//新的加了水印图片覆盖原来文件
	if str == "gif" {
		err = gifFontWater(fileName, fileName, typeface)
	} else {
		err = staticFontWater(fileName, fileName, str, typeface)
	}
	return err
}

//gif图片水印
func gifFontWater(srcFile, newImage string, typeface []app.FontInfo) (err error) {
	imgFile, err := os.Open(srcFile)

	if err != nil {
		return errors.New("打开文件失败！" + err.Error())
	}

	defer func(imgFile *os.File) {
		_ = imgFile.Close()
	}(imgFile)

	var err2 error
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
			rgbImg, err2 = common(rgbImg, typeface) //添加文字水印
			if err2 != nil {
				break
			}
			//定义一个新的图片调色板img.Bounds()：使用原图的颜色域，gifImg.Palette：使用原图的调色板
			p1 := image.NewPaletted(gifImg.Bounds(), gifImg.Palette)
			//把绘制过文字的图片添加到新的图片调色板上
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
		if err2 != nil {
			return err2
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

//png,jpeg图片水印
func staticFontWater(srcFile, newImage, status string, typeface []app.FontInfo) (err error) {
	//需要加水印的图片
	imgFile, err := os.Open(srcFile)

	if err != nil {
		return errors.New("打开文件失败！" + err.Error() + srcFile)
	}

	defer func(imgFile *os.File) {
		_ = imgFile.Close()
	}(imgFile)

	var staticImg image.Image
	if status == "png" {
		staticImg, _ = png.Decode(imgFile)
	} else {
		staticImg, _ = jpeg.Decode(imgFile)
	}
	rgbImg := image.NewNRGBA(staticImg.Bounds())
	for y := 0; y < rgbImg.Bounds().Dy(); y++ {
		for x := 0; x < rgbImg.Bounds().Dx(); x++ {
			rgbImg.Set(x, y, staticImg.At(x, y))
		}
	}

	rgbImg, err = common(rgbImg, typeface) //添加文字水印
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

//添加文字水印函数
func common(rgbImg *image.NRGBA, typeface []app.FontInfo) (*image.NRGBA, error) {
	//拷贝一个字体文件到运行目录
	fontBytes, err := ioutil.ReadFile(app.Conf().Water.WaterFont)
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
		info := t.Message
		f := freetype.NewContext()
		f.SetDPI(108)
		f.SetFont(font)
		f.SetFontSize(t.Size)
		f.SetClip(rgbImg.Bounds())
		f.SetDst(rgbImg)
		f.SetSrc(image.NewUniform(color.RGBA{R: t.R, G: t.G, B: t.B, A: t.A}))

		first := 0
		two := 0
		switch t.Position {
		case TopLeft:
			first = t.Dx
			two = t.Dy + int(f.PointToFixed(t.Size)>>6)
		case TopRight:
			first = rgbImg.Bounds().Dx() - len(info)*4 - t.Dx
			two = t.Dy + int(f.PointToFixed(t.Size)>>6)
		case BottomLeft:
			first = t.Dx
			two = rgbImg.Bounds().Dy() - t.Dy
		case BottomRight:
			first = rgbImg.Bounds().Dx() - len(info)*4 - t.Dx
			two = rgbImg.Bounds().Dy() - t.Dy
		case Center:
			first = (rgbImg.Bounds().Dx() - len(info)*4) / 2
			two = (rgbImg.Bounds().Dy() - t.Dy) / 2
		default:
			errNum = 0
			break Loop
		}

		pt := freetype.Pt(first, two)
		_, err = f.DrawString(info, pt)
		if err != nil {
			break
		}
	}
	if errNum == 0 {
		err = errors.New("坐标值不对")
	}
	return rgbImg, err
}

