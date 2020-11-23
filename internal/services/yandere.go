package services

import (
	"Asane/internal/api/yandere"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Comdex/imgo"
)

const imageResolutionLimit = 3000000
const noise = 300

func yandereSerchTags(params []string) string {
	if len(params) == 0 {
		return "请输入需要搜索的tag\n例：\nasane tag loli"
	}

	result, err := yandere.Client.SearchTags(params[0])

	if err != nil {
		return err.Error()
	}

	str := []string{}
	for _, respObj := range result {
		str = append(str, respObj.Name)
	}

	return fmt.Sprintf("搜索结果：\n%s", strings.Join(str, "\n"))
}

func yandereRandomR18Illust(params []string) string {
	if len(params) > 4 {
		params = params[:3]
	}
	post, err := yandere.Client.RandomExplicitPost(params)
	if err != nil {
		return err.Error()
	}

	imageDir := ""
	if imageDir = os.Getenv("CQHTTP_IMAGES_DIR"); imageDir == "" {
		log.Fatal("缺少环境变量：Yandere下载文件夹")
	}
	imageName := path.Base(post.JpegURL)
	imagePath := path.Join(imageDir, imageName)

	err = downloadFile(post.JpegURL, imagePath)
	if err != nil {
		log.Error(err)
	}
	processIllust(imagePath, post.JpegHeight, post.JpegWidth)

	return fmt.Sprintf("https://yande.re/post/show/%d [CQ:image,file=%s]", post.ID, imageName)
}

func processIllust(file string, height int, width int) {
	log.Info("开始图片反和谐处理")
	defer log.Info("反和谐处理完毕")

	var raw = [][][]uint8{}
	var err error
	if resolution := height * width; resolution > imageResolutionLimit {
		log.Tracef("图片需要压缩（当前分辨率：%d=%d*%d）", resolution, height, width)
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
	newWidth := width + noise
	newImage := imgo.NewRGBAMatrix(newHeight, newWidth)

	log.Trace("添加噪点")
	halfNoise := noise / 2
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < newHeight; i++ {
		for j := 0; j < newWidth; j++ {
			if i >= halfNoise && i < halfNoise+height && j >= halfNoise && j < halfNoise+width {
				newImage[i][j] = raw[i-halfNoise][j-halfNoise]
			} else {
				newImage[i][j][0] = uint8(rand.Intn(255))
				newImage[i][j][1] = uint8(rand.Intn(255))
				newImage[i][j][2] = uint8(rand.Intn(255))
				newImage[i][j][3] = uint8(255)
			}
		}
	}
	log.Trace("添加噪点完毕")

	err = imgo.SaveAsJPEG(file, newImage, 80)
	if err != nil {
		log.Error(err)
	}
}

func downloadFile(url string, filepath string) error {
	log.Info("开始下载图片")
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
