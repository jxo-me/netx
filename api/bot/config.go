package bot

import (
	"bufio"
	"bytes"
	"fmt"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/netx/x/config"
	"os"
	"time"
)

func (h *hEvent) OnClickConfig(c telebot.IContext) error {
	user := c.Callback().Sender
	cmd := c.Callback().Data
	msg := fmt.Sprintf("选中服务: %s %d\\.\nWhat do you want to do with the bot?", cmd, user.ID)
	cfg := config.Global()
	var buf bytes.Buffer
	bio := bufio.NewWriter(&buf)
	err := cfg.Write(bio, "json")
	if err != nil {
		return c.Reply("OnClickService cfg.Write err:", err.Error())
	}
	err = bio.Flush()
	if err != nil {
		return c.Reply("OnClickService bio.Flush err:", err.Error())
	}
	start := "```"
	end := "```"
	tpl := `
%s
%s
%s
%s
`
	msg = fmt.Sprintf(tpl, msg, start, buf.String(), end)
	selector := &telebot.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			//selector.Data("@update", "update", "update"),
			selector.Data("@saveConfig", "saveConfig", "saveConfig"),
			selector.Data("« 返回 服务列表", "backServices", "backServices"),
		),
	)
	return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickSaveConfig(c telebot.IContext) error {
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
