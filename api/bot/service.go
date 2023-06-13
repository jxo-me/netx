package bot

import (
	"bufio"
	"bytes"
	"fmt"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/netx/api/handler"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
)

func (h *hEvent) OnClickServices(c telebot.IContext) error {
	user := c.Callback().Sender
	cmd := c.Callback().Data
	msg := fmt.Sprintf("选中服务: %s %d\\.\nWhat do you want to do with the bot?", cmd, user.ID)
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
		//selector.Data("@update", "update", "update"),
		selector.Data("@addService", "addService", "addService"),
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

func (h *hEvent) OnClickDetailService(c telebot.IContext) error {
	user := c.Callback().Sender
	cmd := c.Callback().Data
	serviceName := cmd
	msg := fmt.Sprintf("选中服务: %s %d\\.\nWhat do you want to do with the bot?", "", user.ID)
	start := "```"
	end := "```"
	tpl := `
%s
%s
%s
%s
`
	cfg := config.Global()
	var srv *config.ServiceConfig
	for _, service := range cfg.Services {
		if service.Name == serviceName {
			srv = service
		}
	}
	var buf bytes.Buffer
	bio := bufio.NewWriter(&buf)

	err := handler.Write(bio, srv, "json")
	if err != nil {
		return c.Reply("OnClickService cfg.Write err:", err.Error())
	}
	err = bio.Flush()
	if err != nil {
		return c.Reply("OnClickService bio.Flush err:", err.Error())
	}
	msg = fmt.Sprintf(tpl, msg, start, buf.String(), end)

	selector := &telebot.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			//selector.Data("@update", "update", "update"),
			selector.Data("@delService", "delService", serviceName),
			selector.Data("« 返回 服务列表", "backServices", "backServices"),
		),
	)
	return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDelService(c telebot.IContext) error {
	//user := c.Callback().Sender
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

	return c.Send(fmt.Sprintf("%s 删除成功!", cmd))
}
