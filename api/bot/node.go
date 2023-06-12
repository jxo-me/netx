package bot

import (
	"bufio"
	"bytes"
	"fmt"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/netx/api/handler"
	"github.com/jxo-me/netx/x/config"
)

func (h *hEvent) OnClickNodes(c telebot.IContext) error {
	user := c.Callback().Sender
	cmd := c.Callback().Data
	msg := fmt.Sprintf("选中服务: %s %d\\.\nWhat do you want to do with the bot?", cmd, user.ID)
	cfg := config.Global()
	rowList := make([]telebot.Row, 0)
	btnList := make([]telebot.Btn, 0)
	selector := &telebot.ReplyMarkup{}
	for i, service := range cfg.Services {
		btnList = append(btnList, selector.Data(fmt.Sprintf("@%s", service.Name), fmt.Sprintf("%s", service.Name), service.Name))
		if i%3 == 0 {
			rowList = append(rowList, selector.Row(btnList...))
			btnList = make([]telebot.Btn, 0)
		}
	}
	rowList = append(rowList, selector.Row(
		//selector.Data("@update", "update", "update"),
		selector.Data("@saveConfig", "saveConfig", "saveConfig"),
		selector.Data("« 返回 服务列表", "backServices", "backServices"),
	))
	var buf bytes.Buffer
	bio := bufio.NewWriter(&buf)

	err := handler.Write(bio, cfg.Services, "json")
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
	selector.Inline(
		rowList...,
	)
	return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}
