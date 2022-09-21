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

// FontWater 给图片添加文字水印
func (w *img) FontWater(fileName string, typeface []FontInfo) error {
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
func gifFontWater(srcFile, newImage string, typeface []FontInfo) (err error) {
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
		img := image.NewNRGBA(gifImg.Bounds())
		if k == 0 {
			x0 = img.Bounds().Dx()
			y0 = img.Bounds().Dy()
		}
		fmt.Printf("%v, %v\n", img.Bounds().Dx(), img.Bounds().Dy())
		if k == 0 && gifImg2.Image[k+1].Bounds().Dx() > x0 && gifImg2.Image[k+1].Bounds().Dy() > y0 {
			yuan = 1
			break
		}
		if x0 == img.Bounds().Dx() && y0 == img.Bounds().Dy() {
			for y := 0; y < img.Bounds().Dy(); y++ {
				for x := 0; x < img.Bounds().Dx(); x++ {
					img.Set(x, y, gifImg.At(x, y))
				}
			}
			img, err2 = common(img, typeface) //添加文字水印
			if err2 != nil {
				break
			}
			//定义一个新的图片调色板img.Bounds()：使用原图的颜色域，gifImg.Palette：使用原图的调色板
			p1 := image.NewPaletted(gifImg.Bounds(), gifImg.Palette)
			//把绘制过文字的图片添加到新的图片调色板上
			draw.Draw(p1, gifImg.Bounds(), img, image.Point{}, draw.Src)
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
func staticFontWater(srcFile, newImage, status string, typeface []FontInfo) (err error) {
	//需要加水印的图片
	imgFile, err := os.Open(srcFile)

	fmt.Println("ddd")

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
	img := image.NewNRGBA(staticImg.Bounds())
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, staticImg.At(x, y))
		}
	}

	img, err = common(img, typeface) //添加文字水印
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
		err = png.Encode(newFile, img)
	} else {
		err = jpeg.Encode(newFile, img, &jpeg.Options{Quality: 100})
	}
	return err
}

//添加文字水印函数
func common(img *image.NRGBA, typeface []FontInfo) (*image.NRGBA, error) {
	//拷贝一个字体文件到运行目录
	fontBytes, err := ioutil.ReadFile(app.Config().Upload.WaterFont)
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
		f.SetClip(img.Bounds())
		f.SetDst(img)
		f.SetSrc(image.NewUniform(color.RGBA{R: t.R, G: t.G, B: t.B, A: t.A}))

		first := 0
		two := 0
		switch t.Position {
		case TopLeft:
			first = t.Dx
			two = t.Dy + int(f.PointToFixed(t.Size)>>6)
		case TopRight:
			first = img.Bounds().Dx() - len(info)*4 - t.Dx
			two = t.Dy + int(f.PointToFixed(t.Size)>>6)
		case BottomLeft:
			first = t.Dx
			two = img.Bounds().Dy() - t.Dy
		case BottomRight:
			first = img.Bounds().Dx() - len(info)*4 - t.Dx
			two = img.Bounds().Dy() - t.Dy
		case Center:
			first = (img.Bounds().Dx() - len(info)*4) / 2
			two = (img.Bounds().Dy() - t.Dy) / 2
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
	return img, err
}

// FontInfo 定义添加的文字信息
type FontInfo struct {
	Size     float64 //文字大小
	Message  string  //文字内容
	Position int     //文字存放位置
	Dx       int     //文字x轴留白距离
	Dy       int     //文字y轴留白距离
	R        uint8   //文字颜色值RGBA中的R值
	G        uint8   //文字颜色值RGBA中的G值
	B        uint8   //文字颜色值RGBA中的B值
	A        uint8   //文字颜色值RGBA中的A值
}
