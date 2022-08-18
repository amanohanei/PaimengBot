package search_image

import (
	"github.com/RicheyJang/PaimengBot/manager"
	"github.com/RicheyJang/PaimengBot/utils"
	"github.com/RicheyJang/PaimengBot/utils/consts"

	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var info = manager.PluginInfo{
	Name: "搜图",
	Usage: `用法：
	搜图 [图片]: 调用SauceNAO的API搜索图片`,
	Classify: "实用工具",
}
var proxy *manager.PluginProxy

func init() {
	proxy = manager.RegisterPlugin(info)
	if proxy == nil {
		return
	}
	proxy.OnCommands([]string{"搜图", "什么本子", "本子"}, zero.OnlyToMe).SetBlock(true).SecondPriority().Handle(searchImageHandler)
	proxy.AddAPIConfig(consts.APIOfSauceNAOAPIKey, "saucenao.com")
	proxy.AddConfig("api_key", "")
}

func searchImageHandler(ctx *zero.Ctx) {
	urls := utils.GetImageURLs(ctx.Event)
	if len(urls) == 0 { // 没有发图，等待他的下一条消息
		ctx.SendChain(message.At(ctx.Event.UserID), message.Text("图呢？"))
		urls = utils.GetImageURLs(utils.WaitNextMessage(ctx))
		if len(urls) == 0 { // 依旧没有发图
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("那就算啦"))
			return
		}
	}
	// 上锁，防止重复调用
	if proxy.LockUser(0) {
		ctx.Send("有请求正在处理中哦")
		return
	}
	defer proxy.UnlockUser(0)
	// 只查询第一张图
	msg, err := SearchImage(urls[0], utils.IsMessagePrimary(ctx))
	if err != nil {
		log.Warnf("SearchImage err: user=%v,url=%v,err=%v", ctx.Event.UserID, urls[0], err)
	}
	msg = append(message.Message{message.At(ctx.Event.UserID)}, msg...)
	ctx.SendChain(msg...)
}

// SearchAnime 搜番，参数：url为图片链接，返回整理后需要发出的消息体
func SearchImage(url string, showAdult bool) (message.Message, error) {
	return searchImageBySauceNAO(url, true)
}
