package services

import (
	"Asane/internal/yandere"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/Comdex/imgo"
)

func makeMessageYandereRandomR18Illust(params []string) string {
	post := yandere.Client.GetRandomExplicitPost()
	imageDir := ""
	if imageDir = os.Getenv("CQHTTP_IMAGES_DIR"); imageDir == "" {
		log.Fatal("缺少环境变量：Yandere下载文件夹")
	}
	imageName := path.Base(post.JpegURL)
	imagePath := path.Join(imageDir, imageName)

	err := downloadFile(post.JpegURL, imagePath)
	if err != nil {
		log.Fatal(err)
	}
	processIllust(imagePath, post.JpegHeight, post.JpegWidth)

	return fmt.Sprintf("https://yande.re/post/show/%d [CQ:image,file=%s]", post.ID, imageName)
}

func processIllust(file string, height int, width int) {
	log.Info("开始图片反和谐处理")
	defer log.Info("反和谐处理完毕")

	var raw = [][][]uint8{}
	var err error
	if height > width && height > 2000 {
		width = int(2000.0 / float64(height) * float64(width))
		height = 2000
		raw, err = imgo.ResizeForMatrix(file, width, height)
	} else if width > height && width > 2000 {
		height = int(2000.0 / float64(width) * float64(height))
		width = 2000
		raw, err = imgo.ResizeForMatrix(file, width, height)
	} else {
		raw = imgo.MustRead(file)
	}

	if err != nil {
		log.Fatal(err)
	}

	triple := imgo.NewRGBAMatrix(height*2, width)

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			triple[i][j][0] = 255 - raw[i][j][0]
			triple[i][j][1] = 255 - raw[i][j][1]
			triple[i][j][2] = 255 - raw[i][j][2]
			triple[i][j][3] = 255
		}
	}

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			triple[i+height][j] = raw[i][j]
		}
	}

	err = imgo.SaveAsJPEG(file, triple, 80)
	if err != nil {
		log.Fatal(err)
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
