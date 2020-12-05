package services

import (
	"Asane/internal/api/sankaku"
	"fmt"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

func sankakuSerchTags(params []string) string {
	if len(params) == 0 {
		return "请输入需要搜索的tag\n例：\nasane tag loli"
	}

	result, err := sankaku.Client.SearchTags(params[0])

	if err != nil {
		return err.Error()
	}

	str := []string{}
	for _, respObj := range result {
		str = append(str, respObj.Name)
	}

	return fmt.Sprintf("搜索结果：\n%s", strings.Join(str, "\n"))
}

func sankakuRandomR18Illust(params []string) string {
	if len(params) > 4 {
		params = params[:3]
	}
	post, err := sankaku.Client.RandomExplicitPost(params)
	if err != nil {
		return err.Error()
	}

	imageDir := ""
	if imageDir = os.Getenv("CQHTTP_IMAGES_DIR"); imageDir == "" {
		log.Fatal("缺少环境变量：图片下载文件夹")
	}
	imageName := strings.Split(path.Base(post.SampleURL),"?")[0]
	imagePath := path.Join(imageDir, imageName)

	err = DownloadFile(post.SampleURL, imagePath)
	if err != nil {
		log.Error(err)
	}
	ProcessIllust(imagePath, post.SampleHeight, post.SampleWidth)

	return fmt.Sprintf("https://chan.sankakucomplex.com/post/show/%d [CQ:image,file=%s]", post.ID, imageName)
}
