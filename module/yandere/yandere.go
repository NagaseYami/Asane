package yandere

import (
	"fmt"
	"path"
	"strings"

	"github.com/NagaseYami/asane/system"
	"github.com/NagaseYami/asane/util"

	log "github.com/sirupsen/logrus"
)

func YandereSerchTags(params []string) string {
	if len(params) == 0 {
		return "请输入需要搜索的tag\n例：\nasane tag loli"
	}

	result, err := httpClient.SearchTags(params[0])

	if err != nil {
		return err.Error()
	}

	str := []string{}
	for _, respObj := range result {
		str = append(str, respObj.Name)
	}

	return fmt.Sprintf("搜索结果：\n%s", strings.Join(str, "\n"))
}

func YandereRandomExplicitIllust(params []string) string {
	if len(params) > 4 {
		params = params[:3]
	}
	post, err := httpClient.RandomExplicitPost(params)
	if err != nil {
		return err.Error()
	}

	imageName := path.Base(post.JpegURL)
	imagePath := path.Join(system.ImageDirectoryPath, imageName)

	err = util.DownloadFile(post.JpegURL, imagePath)
	if err != nil {
		log.Error(err)
	}
	util.ProcessIllust(imagePath, post.JpegHeight, post.JpegWidth)

	return fmt.Sprintf("[CQ:image,file=%s]\nhttps://yande.re/post/show/%d", imagePath, post.ID)
}

func YandereRandomSafeIllust(params []string) string {
	if len(params) > 4 {
		params = params[:3]
	}
	post, err := httpClient.RandomSafePost(params)
	if err != nil {
		return err.Error()
	}

	imageName := path.Base(post.JpegURL)
	imagePath := path.Join(system.ImageDirectoryPath, imageName)

	err = util.DownloadFile(post.JpegURL, imagePath)
	if err != nil {
		log.Error(err)
	}

	return fmt.Sprintf("[CQ:image,file=%s]\nhttps://yande.re/post/show/%d", imagePath, post.ID)
}
