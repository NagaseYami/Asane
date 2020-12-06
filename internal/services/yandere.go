package services

import (
	"Asane/internal/api/yandere"
	"fmt"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

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

func yandereRandomExplicitIllust(params []string) string {
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

	err = DownloadFile(post.JpegURL, imagePath)
	if err != nil {
		log.Error(err)
	}
	ProcessIllust(imagePath, post.JpegHeight, post.JpegWidth)

	return fmt.Sprintf("https://yande.re/post/show/%d [CQ:image,file=%s]", post.ID, imageName)
}

func yandereRandomSafeIllust(params []string) string {
	if len(params) > 4 {
		params = params[:3]
	}
	post, err := yandere.Client.RandomSafePost(params)
	if err != nil {
		return err.Error()
	}

	imageDir := ""
	if imageDir = os.Getenv("CQHTTP_IMAGES_DIR"); imageDir == "" {
		log.Fatal("缺少环境变量：Yandere下载文件夹")
	}
	imageName := path.Base(post.JpegURL)
	imagePath := path.Join(imageDir, imageName)

	err = DownloadFile(post.JpegURL, imagePath)
	if err != nil {
		log.Error(err)
	}

	return fmt.Sprintf("https://yande.re/post/show/%d [CQ:image,file=%s]", post.ID, imageName)
}
