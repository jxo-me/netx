package bot

import (
	"fmt"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/netx/x/config"
)

func (h *hEvent) OnClickNodes(c telebot.IContext) error {
	var (
		msg string
		str string
		err error
	)
	user := c.Callback().Sender
	cmd := c.Callback().Data
	msg = fmt.Sprintf("选中服务: %s %d\\.\nWhat do you want to do with the bot?", cmd, user.ID)
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
		selector.Data("@saveConfig", "saveConfig", "saveConfig"),
		selector.Data("« 返回 服务列表", "backServices", "backServices"),
	))
	if cfg.Services != nil {
		str, err = ConvertJsonMsg(cfg.Services)
		if err != nil {
			return c.Reply("OnClickNodes ConvertJsonMsg err:", err.Error())
		}
		msg = fmt.Sprintf(CodeMsgTpl, msg, CodeStart, str, CodeEnd)
	}
	selector.Inline(
		rowList...,
	)

	return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}
