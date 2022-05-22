package search_image

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/RicheyJang/PaimengBot/utils"
	"github.com/RicheyJang/PaimengBot/utils/client"
	"github.com/RicheyJang/PaimengBot/utils/consts"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func searchImageBySauceNAO(url string, showAdult bool) (msg message.Message, err error) {
	api := proxy.GetAPIConfig(consts.APIOfSauceNAOAPIKey)
	if len(api) == 0 {
		return message.Message{message.Text("失败了...")}, fmt.Errorf("api of SauceNAO is empty")
	}
	// 整理API URL
	if !strings.HasPrefix(api, "http://") && !strings.HasPrefix(api, "https://") {
		api = "https://" + api
	}
	if !strings.HasSuffix(api, "/") {
		api += "/"
	}

	api = fmt.Sprintf("%ssearch.php?db=999&output_type=2&url=%s&api_key=%s", api, url, proxy.GetConfigString("api_key"))
	log.Infof(api)
	// 调用
	c := client.NewHttpClient(nil)
	rsp, err := c.GetGJson(api)
	log.Infof(rsp.String())
	if err != nil {
		return message.Message{message.Text("出错了...")}, err
	}
	results := rsp.Get("results").Array()
	// log.Infof("删除超时的<%v>插件[%v]子限流器", pl.Key, key)
	// log.Infof()
	if len(results) == 0 { // 没有结果
		return message.Message{message.Text(fmt.Sprintf("%v也不知道", utils.GetBotNickname()))},
			fmt.Errorf("result is empty, error=%v", rsp.Get("error"))
	}
	result := results[0]
	// 解析
	// title := formatMoeResultTitle(result)
	title := result.Get("data").Get("title").String()
	// imgMsg := message.Image(result.Get("image").String()) // 直接以URL格式发送
	imgMsg := message.Image(result.Get("header").Get("thumbnail").String())

	text := fmt.Sprintf("相似度：%v\n作者：%s\n", formatMoeResultSimilarity(result.Get("header").Get("similarity")), result.Get("data").Get("author_name").String())
	creator := fmt.Sprint(result.Get("data").Get("creator").String() + "\n")
	pixiv := fmt.Sprintf("pixiv_id：%s\n", result.Get("data").Get("pixiv_id").String())
	source := fmt.Sprintf("源：%s\n", result.Get("data").Get("source").String())
	return message.Message{message.Text(title), imgMsg, message.Text(text), message.Text(creator), message.Text(pixiv), message.Text(source)}, nil
}

func formatMoeResultSimilarity(similarity gjson.Result) string {
	if !similarity.Exists() {
		return "未知"
	}
	org := similarity.Float()
	str := strconv.FormatFloat(org, 'f', 2, 64) + "%"
	if org <= 90 {
		str += "(较低)"
	}
	return str
}

func getMoeResultEpisode(episode gjson.Result) string {
	if !episode.Exists() {
		return "?"
	}
	i := episode.Int()
	return strconv.FormatInt(i, 10)
}

func formatMoeResultTime(from gjson.Result) string {
	if !from.Exists() {
		return "未知时间"
	}
	org := math.Floor(from.Float())
	return fmt.Sprintf("%02d:%02d", int(org)/60, int(org)%60)
}

func formatMoeResultTitle(result gjson.Result) string {
	title := result.Get("anilist").Get("title")
	if !title.Exists() {
		return result.Get("filename").String()
	}
	res := title.Get("native").String()
	if title.Get("chinese").Type != gjson.Null {
		res += "\n" + title.Get("chinese").String()
	} else if title.Get("english").Type != gjson.Null {
		res += "\n" + title.Get("english").String()
	} else if title.Get("romaji").Type != gjson.Null {
		res += "\n" + title.Get("romaji").String()
	}
	return res
}
