package bot

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/gfbot/handlers"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	parser "github.com/jxo-me/netx/x/config/parsing/ingress"
)

const (
	IngressAdd         = "ingressAdd"
	IngressUpdate      = "ingressUpdate"
	IngressExampleJson = `
{
  "name": "ingress-0",
  "rules": [
    {
      "hostname": "example.com",
      "endpoint": "4d21094e-b74c-4916-86c1-d9fa36ea677b"
    },
    {
      "hostname": "example.org",
      "endpoint": "ac74d9dd-3125-442a-a7c1-f9e49e05faca"
    }
  ]
}
`
)

func (h *hEvent) OnClickIngresses(c telebot.Context) error {
	var (
		msg string
	)

	cfg := config.Global()
	rowList := make([]telebot.Row, 0)
	btnList := make([]telebot.Btn, 0)
	selector := &telebot.ReplyMarkup{}
	for _, service := range cfg.Ingresses {
		btnList = append(btnList, selector.Data(fmt.Sprintf("@%s", service.Name), "detailIngress", service.Name))
	}
	rowList = append(rowList, selector.Split(MaxCol, btnList)...)
	rowList = append(rowList, selector.Row(
		selector.Data("@添加Ingress", "addIngress", "addIngress"),
		selector.Data("« 返回 服务列表", "backServices", "backServices"),
	))

	selector.Inline(
		rowList...,
	)
	msg = "Ingresss List:\n"
	if c.Callback() != nil {
		return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
	}

	return c.Reply(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDetailIngress(c telebot.Context) error {
	var (
		msg string
		str string
		err error
	)

	serviceName := c.Callback().Data
	cfg := config.Global()
	var srv *config.IngressConfig
	for _, service := range cfg.Ingresses {
		if service.Name == serviceName {
			srv = service
		}
	}
	if srv != nil {
		str, err = ConvertJsonMsg(srv)
		if err != nil {
			return c.Reply("OnClickDetailIngress ConvertJsonMsg err:", err.Error())
		}
		msg = fmt.Sprintf(CodeMsgTpl, msg, CodeStart, str, CodeEnd)
	}
	selector := &telebot.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			selector.Data("@更新Ingress", "updateIngress", serviceName),
			selector.Data("@删除Ingress", "delIngress", serviceName),
		),
		selector.Row(
			selector.Data("« 返回 服务列表", "backServices", "backServices"),
		),
	)
	return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDelIngress(c telebot.Context) error {
	serviceName := c.Callback().Data
	svc := app.Runtime.IngressRegistry().Get(serviceName)
	if svc == nil {
		return c.Send(ErrNotFound)
	}

	app.Runtime.IngressRegistry().Unregister(serviceName)

	_ = config.OnUpdate(func(c *config.Config) error {
		ingresses := c.Ingresses
		c.Ingresses = nil
		for _, s := range ingresses {
			if s.Name == serviceName {
				continue
			}
			c.Ingresses = append(c.Ingresses, s)
		}
		return nil
	})
	_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 删除成功!", serviceName)})
	return h.OnClickIngresses(c)
}

func AddIngressConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startAddIngressHandler), // 入口
		map[string][]telebot.IHandler{
			IngressAdd: {telebot.HandlerFunc(addIngressHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelIngressHandler),
			AllowReEntry: true,
		}, // options
	)
}

func UpdateIngressConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startUpdateIngressHandler), // 入口
		map[string][]telebot.IHandler{
			IngressUpdate: {telebot.HandlerFunc(updateIngressHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelIngressHandler),
			AllowReEntry: true,
		}, // options
	)
}

func startAddIngressHandler(ctx telebot.Context) error {
	example := fmt.Sprintf(CodeTpl, CodeStart, IngressExampleJson, CodeEnd)
	err := ctx.Send(fmt.Sprintf("请输入 Ingress 配置?\nExample：%s\n您可以随时键入 /cancel 来取消该操作。", example), &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(IngressAdd)
}

func addIngressHandler(ctx telebot.Context) error {
	var (
		data config.IngressConfig
	)

	str := ctx.Message().Text
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("addIngressHandler json.Unmarshal error:", err.Error())
	}
	v := parser.ParseIngress(&data)
	if err = app.Runtime.IngressRegistry().Register(data.Name, v); err != nil {
		return ctx.Reply(ErrDup)
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		c.Ingresses = append(c.Ingresses, &data)
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 添加成功!", data.Name)})
	_ = Event.OnClickIngresses(ctx)

	return handlers.EndConversation()
}

func startUpdateIngressHandler(ctx telebot.Context) error {
	srvName := ctx.Callback().Data
	err := ctx.Bot().Store().UpdateData(ctx, IngressUpdate, srvName)
	if err != nil {
		return fmt.Errorf("failed UpdateData message: %w", err)
	}
	example := fmt.Sprintf(CodeTpl, CodeStart, IngressExampleJson, CodeEnd)
	err = ctx.Send(fmt.Sprintf("请输入 Ingress 配置?\nExample：%s\n您可以随时键入 /cancel 来取消该操作。", example), &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(IngressUpdate)
}

func updateIngressHandler(ctx telebot.Context) error {
	var (
		data config.IngressConfig
	)
	state, err := ctx.Bot().Store().Get(ctx)
	if err != nil {
		return ctx.Reply("updateIngressHandler Store().Get error:", err.Error())
	}
	srvName := ""
	if sn, ok := state.Data[IngressUpdate]; ok {
		srvName = gconv.String(sn)
	}
	fmt.Println(fmt.Sprintf("srvName :%s", srvName))
	if srvName == "" {
		return ctx.Reply(ErrInvalid)
	}
	str := ctx.Message().Text
	err = json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("updateIngressHandler json.Unmarshal error:", err.Error())
	}

	if !app.Runtime.IngressRegistry().IsRegistered(srvName) {
		return ctx.Reply(ErrNotFound)
	}

	data.Name = srvName

	v := parser.ParseIngress(&data)

	app.Runtime.IngressRegistry().Unregister(srvName)

	if err = app.Runtime.IngressRegistry().Register(srvName, v); err != nil {
		return ctx.Reply(ErrDup)
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.Ingresses {
			if c.Ingresses[i].Name == srvName {
				c.Ingresses[i] = &data
				break
			}
		}
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 更新成功!", data.Name)})
	_ = Event.OnClickIngresses(ctx)

	return handlers.EndConversation()
}

func cancelIngressHandler(ctx telebot.Context) error {
	err := ctx.Reply("添加 Ingress 已被取消。 还有什么我可以为你做的吗？", &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send cancelHandler message: %w", err)
	}
	//return handlers.EndConversation()
	return nil
}
