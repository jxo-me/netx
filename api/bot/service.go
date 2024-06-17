package bot

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/gfbot/handlers"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	parser "github.com/jxo-me/netx/x/config/parsing/service"
)

const (
	ServiceAdd         = "serviceAdd"
	ServiceUpdate      = "serviceUpdate"
	ServiceExampleJson = `
{
  "name": "service-0",
  "addr": ":8080",
  "interface": "eth0",
  "admission": "admission-0",
  "bypass": "bypass-0",
  "resolver": "resolver-0",
  "hosts": "hosts-0",
  "handler": {
    "type": "http",
    "auth": {
      "username": "gost",
      "password": "gost"
    },
    "auther": "auther-0",
    "retries": 1,
    "chain": "chain-0",
    "metadata": {
      "bar": "baz",
      "foo": "bar"
    }
  },
  "listener": {
    "type": "tcp",
    "auth": {
      "username": "user",
      "password": "pass"
    },
    "auther": "auther-0",
    "chain": "chain-0",
    "tls": {
      "certFile": "cert.pem",
      "keyFile": "key.pem",
      "caFile": "ca.pem"
    },
    "metadata": {
      "abc": "xyz",
      "def": 456
    }
  },
  "forwarder": {
    "nodes": [
      {
        "name": "target-0",
        "addr": "192.168.1.1:1234"
      },
      {
        "name": "target-1",
        "addr": "192.168.1.2:2345"
      }
    ],
    "selector": {
      "strategy": "round",
      "maxFails": 1,
      "failTimeout": 30
    }
  }
}
`
)

func (h *hEvent) OnClickServices(c telebot.Context) error {
	var (
		msg string
	)

	cfg := config.Global()
	rowList := make([]telebot.Row, 0)
	btnList := make([]telebot.Btn, 0)
	selector := &telebot.ReplyMarkup{}
	for _, service := range cfg.Services {
		btnList = append(btnList, selector.Data(fmt.Sprintf("@%s", service.Name), "detailService", service.Name))
	}
	rowList = append(rowList, selector.Split(MaxCol, btnList)...)
	rowList = append(rowList,
		selector.Row(
			selector.Data("@添加服务", "addService", "addService"),
			selector.Data("« 返回 服务列表", "backServices", "backServices"),
		))

	selector.Inline(
		rowList...,
	)
	msg = "Services List:\n"
	if c.Callback() != nil {
		return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
	}

	return c.Reply(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDetailService(c telebot.Context) error {
	var (
		msg string
		str string
		err error
	)

	serviceName := c.Callback().Data
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
			selector.Data("@更新服务", "updateService", serviceName),
			selector.Data("@删除服务", "delService", serviceName),
		),
		selector.Row(
			selector.Data("« 返回 服务列表", "backServices", "backServices"),
		),
	)
	return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDelService(c telebot.Context) error {
	cmd := c.Callback().Data
	fmt.Println("OnClickDelService cmd:", cmd)
	serviceName := cmd
	svc := app.Runtime.ServiceRegistry().Get(serviceName)
	if svc == nil {
		return c.Send(ErrNotFound)
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

func AddServiceConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startAddServiceHandler), // 入口
		map[string][]telebot.IHandler{
			ServiceAdd: {telebot.HandlerFunc(addServiceHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelServiceHandler),
			AllowReEntry: true,
		}, // options
	)
}

func UpdateServiceConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startUpdateServiceHandler), // 入口
		map[string][]telebot.IHandler{
			ServiceUpdate: {telebot.HandlerFunc(updateServiceHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelServiceHandler),
			AllowReEntry: true,
		}, // options
	)
}

func startAddServiceHandler(ctx telebot.Context) error {
	example := fmt.Sprintf(CodeTpl, CodeStart, ServiceExampleJson, CodeEnd)
	err := ctx.Send(fmt.Sprintf("请输入 服务 配置?\nExample：%s\n您可以随时键入 /cancel 来取消该操作。", example), &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(ServiceAdd)
}

func addServiceHandler(ctx telebot.Context) error {
	var (
		data config.ServiceConfig
	)

	str := ctx.Message().Text
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("addServiceHandler json.Unmarshal error:", err.Error())
	}
	if app.Runtime.ServiceRegistry().IsRegistered(data.Name) {
		return ctx.Reply(ErrDup)
	}

	svc, err := parser.ParseService(&data)
	if err != nil {
		return ctx.Reply(ErrCreate)
	}

	if err := app.Runtime.ServiceRegistry().Register(data.Name, svc); err != nil {
		_ = svc.Close()
		return ctx.Reply(ErrDup)
	}

	go svc.Serve()

	_ = config.OnUpdate(func(c *config.Config) error {
		c.Services = append(c.Services, &data)
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 添加成功!", data.Name)})
	_ = Event.OnClickServices(ctx)

	return handlers.EndConversation()
}

func startUpdateServiceHandler(ctx telebot.Context) error {
	srvName := ctx.Callback().Data
	err := ctx.Bot().Store().UpdateData(ctx, ServiceUpdate, srvName)
	if err != nil {
		return fmt.Errorf("failed UpdateData message: %w", err)
	}
	example := fmt.Sprintf(CodeTpl, CodeStart, ServiceExampleJson, CodeEnd)
	err = ctx.Send(fmt.Sprintf("请输入 服务 配置?\nExample：%s\n您可以随时键入 /cancel 来取消该操作。", example), &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(ServiceUpdate)
}

func updateServiceHandler(ctx telebot.Context) error {
	var (
		data config.ServiceConfig
	)
	state, err := ctx.Bot().Store().Get(ctx)
	if err != nil {
		return ctx.Reply("updateServiceHandler Store().Get error:", err.Error())
	}
	srvName := ""
	if sn, ok := state.Data[ServiceUpdate]; ok {
		srvName = gconv.String(sn)
	}
	fmt.Println(fmt.Sprintf("srvName :%s", srvName))
	if srvName == "" {
		return ctx.Reply(ErrInvalid)
	}
	str := ctx.Message().Text
	err = json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("updateServiceHandler json.Unmarshal error:", err.Error())
	}
	old := app.Runtime.ServiceRegistry().Get(srvName)
	if old == nil {
		return ctx.Reply(ErrInvalid)
	}
	_ = old.Close()

	data.Name = srvName

	svc, err := parser.ParseService(&data)
	if err != nil {
		return ctx.Reply(ErrCreate)
	}

	app.Runtime.ServiceRegistry().Unregister(srvName)

	if err = app.Runtime.ServiceRegistry().Register(srvName, svc); err != nil {
		_ = svc.Close()
		return ctx.Reply(ErrDup)
	}

	go svc.Serve()

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.Services {
			if c.Services[i].Name == srvName {
				c.Services[i] = &data
				break
			}
		}
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 更新成功!", data.Name)})
	_ = Event.OnClickServices(ctx)

	return handlers.EndConversation()
}

func cancelServiceHandler(ctx telebot.Context) error {
	err := ctx.Reply("当前操作已被取消。 还有什么我可以为你做的吗？", &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send cancelHandler message: %w", err)
	}
	return nil
}
