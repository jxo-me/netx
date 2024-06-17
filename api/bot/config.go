package bot

import (
	"fmt"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/netx/x/config"
	"os"
	"time"
)

func (h *hEvent) OnClickConfig(c telebot.Context) error {
	var (
		msg string
		str string
		err error
	)
	cfg := config.Global()
	if cfg != nil {
		str, err = ConvertJsonMsg(cfg)
		if err != nil {
			return c.Reply("OnClickDetailService ConvertJsonMsg err:", err.Error())
		}
		msg = fmt.Sprintf(CodeTpl, CodeStart, str, CodeEnd)
	}

	selector := &telebot.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			selector.Data("@saveConfig", "saveConfig", "saveConfig"),
			selector.Data("« 返回 服务列表", "backServices", "backServices"),
		),
	)
	return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickSaveConfig(c telebot.Context) error {
	t := time.Now().Format(time.DateOnly)
	file := fmt.Sprintf("./gost_%s.json", t)
	f, err := os.Create(file)
	if err != nil {
		return c.Send("OnClickSaveConfig os.Create err:", err.Error())
	}
	defer f.Close()

	if err := config.Global().Write(f, "json"); err != nil {
		return c.Send("OnClickSaveConfig config.Global().Write err:", err.Error())
	}
	msg := fmt.Sprintf("文件名:%s\n保存成功！", file)
	selector := &telebot.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			selector.Data("« 返回 服务列表", "backServices", "backServices"),
		),
	)
	return c.Send(msg, &telebot.SendOptions{ReplyMarkup: selector})
}
