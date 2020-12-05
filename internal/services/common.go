package services

import (
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Comdex/imgo"
)

const imageResolutionLimit = 3000000
const qqGroupImageWidthMax = 1280
const noise = 300

//ProcessIllust 图片压缩，反和谐等处理
func ProcessIllust(file string, height int, width int) {
	log.Info("开始图片反和谐处理")
	defer log.Info("反和谐处理完毕")

	var raw = [][][]uint8{}
	var err error
	if width > qqGroupImageWidthMax {
		log.Tracef("图片宽度超过QQ压图阈值，需要压缩（当前分辨率：%d=%d*%d）", height*width, height, width)
		height = int(float64(qqGroupImageWidthMax) / float64(width) * float64(height))
		width = qqGroupImageWidthMax
		raw, err = imgo.ResizeForMatrix(file, width, height)
		log.Tracef("压缩完毕（压缩后分辨率：%d=%d*%d）", height*width, height, width)
	} else if resolution := height * width; resolution > imageResolutionLimit {
		log.Tracef("图片分辨率太大，需要压缩（当前分辨率：%d=%d*%d）", resolution, height, width)
		scale := math.Sqrt(float64(imageResolutionLimit) / float64(resolution))

		height = int(float64(height) * scale)
		width = int(float64(width) * scale)

		raw, err = imgo.ResizeForMatrix(file, width, height)
		log.Tracef("压缩完毕（压缩后分辨率：%d=%d*%d）", height*width, height, width)
	} else {
		raw = imgo.MustRead(file)
	}

	if err != nil {
		log.Error(err)
	}

	newHeight := height + noise
	newImage := imgo.NewRGBAMatrix(newHeight, width)

	log.Trace("复制像素")
	halfNoise := noise / 2

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			newImage[i+halfNoise][j] = raw[i][j]
		}
	}
	log.Trace("复制完毕")

	log.Trace("添加噪点")
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < halfNoise; i++ {
		for j := 0; j < width; j++ {
			newImage[i][j][0] = uint8(rand.Intn(255))
			newImage[i][j][1] = uint8(rand.Intn(255))
			newImage[i][j][2] = uint8(rand.Intn(255))
			newImage[i][j][3] = uint8(255)
		}
	}
	for i := 0; i < halfNoise; i++ {
		for j := 0; j < width; j++ {
			newImage[i+halfNoise+height][j][0] = uint8(rand.Intn(255))
			newImage[i+halfNoise+height][j][1] = uint8(rand.Intn(255))
			newImage[i+halfNoise+height][j][2] = uint8(rand.Intn(255))
			newImage[i+halfNoise+height][j][3] = uint8(255)
		}
	}
	log.Trace("添加噪点完毕")

	err = imgo.SaveAsJPEG(file, newImage, 80)
	if err != nil {
		log.Error(err)
	}
}

//DownloadFile 下载文件
func DownloadFile(url string, filepath string) error {
	log.Info("开始下载文件")
	defer log.Info("下载完毕")
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
