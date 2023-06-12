package bot

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/netx/x/config"
	"strings"
)

var (
	Event              = hEvent{}
	StartTextCommand   = "/start"
	NodeTextCommand    = "/myHosts"
	ParsingTextCommand = "/parsing"
	GostTextCommand    = "/gost"
	// Click group
	OnClickAdmissions  = "\fAdmissions"
	OnClickAuthers     = "\fAuthers"
	OnClickBypass      = "\fBypass"
	OnClickHops        = "\fHops"
	OnClickIngress     = "\fIngress"
	OnClickServices    = "\fServices"
	OnClickChains      = "\fChains"
	OnClickHosts       = "\fHosts"
	OnClickResolver    = "\fResolver"
	OnClickLimiter     = "\fLimiter"
	OnClickConnLimiter = "\fConnLimiter"
	OnClickRateLimiter = "\fRateLimiter"
	OnClickConfig      = "\fConfig"
	OnClickSaveConfig  = "\fsaveConfig"
	OnClickNode        = "\fNode"
	OnClickAddNode     = "\fAddNode"

	OnBackServices = "\fbackServices"
	OnBackHosts    = "\fbackHosts"
)

type (
	hEvent struct{}
)

func (h *hEvent) OnText(c telebot.IContext) error {
	return c.Send("OnText")
}

func getSelectMenus() *telebot.ReplyMarkup {
	selector := &telebot.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			selector.Data("@Admissions", "Admissions", "Admissions"),
			selector.Data("@Authers", "Authers", "Authers"),
			selector.Data("@Bypass", "Bypass", "Bypass"),
		),
		selector.Row(
			selector.Data("@Hops", "Hops", "Hops"),
			selector.Data("@Ingress", "Ingress", "Ingress"),
			selector.Data("@Services", "Services", "Services"),
		),
		selector.Row(
			selector.Data("@Chains", "Chains", "Chains"),
			selector.Data("@Hosts", "Hosts", "Hosts"),
			selector.Data("@Resolver", "Resolver", "Resolver"),
		),
		selector.Row(
			selector.Data("@Limiter", "Limiter", "Limiter"),
			selector.Data("@ConnLimiter", "ConnLimiter", "ConnLimiter"),
			selector.Data("@RateLimiter", "RateLimiter", "RateLimiter"),
		),
		selector.Row(
			selector.Data("@Config", "Config", "Config"),
			selector.Data("« 返回 节点列表", "backHosts", "backHosts"),
		),
	)
	return selector
}

func getSelectHosts() *telebot.ReplyMarkup {
	selector := &telebot.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			selector.Data("@Node", "Node", "Node"),
		),
		selector.Row(
			selector.Data("@AddNode", "AddNode", "AddNode"),
		),
	)
	return selector
}

func (h *hEvent) OnBackServices(c telebot.IContext) error {
	selector := getSelectMenus()
	return c.Edit("从下面的列表中选择一个服务:", selector)
}

func (h *hEvent) OnBackHosts(c telebot.IContext) error {
	selector := getSelectHosts()
	return c.Edit("从下面的列表中选择一个节点:", selector)
}

func (h *hEvent) OnClickNode(c telebot.IContext) error {
	selector := getSelectHosts()
	return c.Send("从下面的列表中选择一个服务:", selector)
}

func (h *hEvent) OnClickService(c telebot.IContext) error {
	cmd := c.Callback().Data
	switch strings.ToLower(cmd) {
	case "config":
		return h.OnClickConfig(c)
	}
	user := c.Callback().Sender
	msg := fmt.Sprintf("选中服务: %s %d\\.\nWhat do you want to do with the bot?", cmd, user.ID)
	selector := &telebot.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			//selector.Data("@update", "update", "update"),
			selector.Data("@add", "add", "add"),
			selector.Data("@update", "update", "update"),
		),
		selector.Row(selector.Data("« 返回 服务列表", "backServices", "backServices")),
	)
	return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnCallback(c telebot.IContext) error {
	cmd := c.Callback().Data
	return c.Send(fmt.Sprintf("OnCallback:%s", cmd))
}

func (h *hEvent) OnUserJoined(c telebot.IContext) error {
	return c.Send("OnUserJoined")
}

func (h *hEvent) OnParsingCommand(c telebot.IContext) error {
	var (
		services stringList
		nodes    stringList
	)

	payload := c.Message().Text
	str := strings.Split(payload, " ")
	cmd := flag.NewFlagSet(gconv.String(str[:1]), flag.ContinueOnError)
	cmd.Var(&services, "L", "service list")
	cmd.Var(&nodes, "F", "chain node list")
	err := cmd.Parse(str[1:])
	if err != nil {
		return c.Reply("OnParsingCommand err:", err.Error())
	}
	cfg, err := buildConfigFromCmd(services, nodes)
	if err != nil {
		return c.Reply("OnParsingCommand err:", err.Error())
	}
	var buf bytes.Buffer
	bio := bufio.NewWriter(&buf)
	err = cfg.Write(bio, "json")
	if err != nil {
		return c.Reply("OnParsingCommand err:", err.Error())
	}
	err = bio.Flush()
	if err != nil {
		return c.Reply("OnParsingCommand err:", err.Error())
	}
	start := "```"
	end := "```"
	tpl := `
%s
%s
%s
`
	msg := fmt.Sprintf(tpl, start, buf.String(), end)
	return c.Reply(msg, &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnStartCommand(c telebot.IContext) error {
	payload := c.Message().Text
	user := c.Message().Sender
	str := strings.Split(payload, " ")
	token := ""
	if len(str) >= 2 {
		token = strings.Join(str[1:], "")
	}
	return c.Send(fmt.Sprintf("欢迎 %s 加入 参数:%s", user.Username, token))
}

func (h *hEvent) OnGostCommand(c telebot.IContext) error {
	var (
		services stringList
		nodes    stringList
	)

	payload := c.Message().Text
	str := strings.Split(payload, " ")
	cmd := flag.NewFlagSet(gconv.String(str[:1]), flag.ContinueOnError)
	cmd.Var(&services, "L", "service list")
	cmd.Var(&nodes, "F", "chain node list")
	err := cmd.Parse(str[1:])
	if err != nil {
		return c.Reply("OnParsingCommand err:", err.Error())
	}
	cfg, err := buildConfigFromCmd(services, nodes)
	if err != nil {
		return c.Reply("OnParsingCommand err:", err.Error())
	}
	config.Set(cfg)
	for _, svc := range buildService(cfg) {
		svc := svc
		go func() {
			svc.Serve()
		}()
	}
	cfgNew := config.Global()
	var buf bytes.Buffer
	bio := bufio.NewWriter(&buf)
	err = cfgNew.Write(bio, "json")
	if err != nil {
		return c.Reply("OnParsingCommand err:", err.Error())
	}
	err = bio.Flush()
	if err != nil {
		return c.Reply("OnParsingCommand err:", err.Error())
	}
	start := "```"
	end := "```"
	tpl := `
%s
%s
%s
`
	msg := fmt.Sprintf(tpl, start, buf.String(), end)
	return c.Reply(msg, &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
}
