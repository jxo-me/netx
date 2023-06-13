package bot

import (
	"fmt"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
)

func (h *hEvent) OnClickServices(c telebot.IContext) error {
	var (
		msg string
		//str  string
		//err  error
		user *telebot.User
	)
	if c.Message() != nil {
		user = c.Message().Sender
	}
	if c.Callback() != nil {
		user = c.Callback().Sender
	}
	msg = fmt.Sprintf("选中服务: %s %d\\.\nWhat do you want to do with the bot?", "", user.ID)
	cfg := config.Global()
	rowList := make([]telebot.Row, 0)
	btnList := make([]telebot.Btn, 0)
	selector := &telebot.ReplyMarkup{}
	for i, service := range cfg.Services {
		btnList = append(btnList, selector.Data(fmt.Sprintf("@%s", service.Name), "detailService", service.Name))
		if i%3 == 0 {
			rowList = append(rowList, selector.Row(btnList...))
			btnList = make([]telebot.Btn, 0)
		}
	}
	rowList = append(rowList, selector.Row(
		selector.Data("@添加服务", "addService", "addService"),
		selector.Data("« 返回 服务列表", "backServices", "backServices"),
	))
	//if cfg.Services != nil {
	//	str, err = ConvertJsonMsg(cfg.Services)
	//	if err != nil {
	//		return c.Reply("OnClickServices ConvertJsonMsg err:", err.Error())
	//	}
	//	msg = fmt.Sprintf(CodeMsgTpl, msg, CodeStart, str, CodeEnd)
	//}

	selector.Inline(
		rowList...,
	)
	if c.Message() != nil {
		return c.Reply(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
	}

	return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDetailService(c telebot.IContext) error {
	var (
		msg string
		str string
		err error
	)
	user := c.Callback().Sender
	serviceName := c.Callback().Data
	msg = fmt.Sprintf("选中服务: %s %d\\.\nWhat do you want to do with the bot?", "", user.ID)

	cfg := config.Global()
	var srv *config.ServiceConfig
	for _, service := range cfg.Services {
		if service.Name == serviceName {
			srv = service
		}
	}
	if srv != nil {
		str, err = ConvertJsonMsg(srv)
		if err != nil {
			return c.Reply("OnClickDetailService ConvertJsonMsg err:", err.Error())
		}
		msg = fmt.Sprintf(CodeMsgTpl, msg, CodeStart, str, CodeEnd)
	}
	selector := &telebot.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			selector.Data("@删除服务", "delService", serviceName),
			selector.Data("« 返回 服务列表", "backServices", "backServices"),
		),
	)
	return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDelService(c telebot.IContext) error {
	cmd := c.Callback().Data
	fmt.Println("OnClickDelService cmd:", cmd)
	serviceName := cmd
	svc := app.Runtime.ServiceRegistry().Get(serviceName)
	if svc == nil {
		return c.Send("object not found")
	}

	app.Runtime.ServiceRegistry().Unregister(serviceName)
	_ = svc.Close()

	_ = config.OnUpdate(func(c *config.Config) error {
		services := c.Services
		c.Services = nil
		for _, s := range services {
			if s.Name == serviceName {
				continue
			}
			c.Services = append(c.Services, s)
		}
		return nil
	})
	_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 删除成功!", cmd)})
	return h.OnClickServices(c)
}
