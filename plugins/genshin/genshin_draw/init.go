package genshin_draw

import (
	"math/rand"
	"time"

	"github.com/RicheyJang/PaimengBot/manager"
	"github.com/RicheyJang/PaimengBot/utils"
	"github.com/RicheyJang/PaimengBot/utils/consts"

	zero "github.com/wdvxdr1123/ZeroBot"
)

const GenshinDrawPoolDir = consts.GenshinDataDir + "/pool"
const GenshinPoolPicDir = GenshinDrawPoolDir + "/pic"

var info = manager.PluginInfo{
	Name: "模拟原神抽卡",
	Usage: `用法：
	原神[卡池名]一发：来一发！
	原神[卡池名]十连：来个十连！
	原神当前卡池：查看当前原神卡池列表
	注：[卡池名]请以"原神当前卡池"中的为准哦
示例：
	原神十连：十连抽常驻池
	原神武器十连：十连抽UP武器池（若存在）
	原神角色十连：十连抽UP角色池（若仅存在1个UP角色池）
	原神角色1十连：十连抽第1个UP角色池（若存在两个或以上UP角色池）`,
	SuperUsage: `更新指令：
	原神抽卡更新：强制更新图片素材以及卡池信息
	另外每天2点10分会自动更新
config-plugin配置项：
	genshin_draw.skip.normal4: 角色、武器池中所不包含的4星角色列表`,
	Classify: "原神相关",
}
var proxy *manager.PluginProxy

func init() {
	proxy = manager.RegisterPlugin(info)
	if proxy == nil {
		return
	}
	proxy.OnCommands([]string{"原神抽卡更新", "原神卡池更新"}, zero.SuperUserPermission).SetBlock(true).SetPriority(3).Handle(updateAllForce)
	proxy.OnCommands([]string{"原神当前卡池", "原神当前祈愿"}).SetBlock(true).SetPriority(3).Handle(showNowPool)
	proxy.OnRegex(`^原神(.*)(10|十)[发连]`).SetBlock(true).SetPriority(4).Handle(drawTenCard)
	proxy.OnRegex(`^原神(.*)[1一][发连]`).SetBlock(true).SetPriority(5).Handle(drawOneCard)
	_, _ = proxy.AddScheduleDailyFunc(2, 10, func() { _ = updateAll() })
	proxy.AddConfig("skip.normal4", []string{"丽莎", "安柏", "凯亚"})
	if !utils.DirExists(GenshinDrawPoolDir) || !utils.DirExists(GenshinPoolPicDir) {
		go updateAll()
	}
	rand.Seed(time.Now().Unix())
}

const (
	PoolNormal = iota
	PoolCharacter
	PoolWeapon
)

var poolTypeMap = map[int]string{
	PoolNormal:    "常驻",
	PoolCharacter: "角色",
	PoolWeapon:    "武器",
}

func getPrefixByType(tp int) (prefix string) {
	return poolTypeMap[tp]
}
