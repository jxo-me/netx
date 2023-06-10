package bot

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/netx/x/config"
	"strings"
)

var (
	Event              = hEvent{}
	NodeTextCommand    = "/myHosts"
	ParsingTextCommand = "/parsing"
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
	user := c.Callback().Sender
	cmd := c.Callback().Data
	msg := fmt.Sprintf("选中服务: %s %d\\.\nWhat do you want to do with the bot?", cmd, user.ID)
	cfg := config.Global()
	var buf bytes.Buffer
	bio := bufio.NewWriter(&buf)
	err := cfg.Write(bio, "json")
	if err != nil {
		return err
	}
	err = bio.Flush()
	if err != nil {
		return err
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
			selector.Data("@add", "add", "add"),
			selector.Data("@update", "update", "update"),
		),
		selector.Row(selector.Data("« 返回 服务列表", "backServices", "backServices")),
	)
	return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnCallback(c telebot.IContext) error {
	return c.Send("OnCallback")
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
	//flag.Parse()
	cmd := flag.NewFlagSet("/parsing", flag.ContinueOnError)
	cmd.Var(&services, "L", "service list")
	cmd.Var(&nodes, "F", "chain node list")
	err := cmd.Parse(strings.Split(payload, " "))
	if err != nil {
		return c.Send("OnParsingCommand err:", err.Error())
	}
	cfg, err := buildConfigFromCmd(services, nodes)
	if err != nil {
		return c.Send("OnParsingCommand err:", err.Error())
	}
	var buf bytes.Buffer
	bio := bufio.NewWriter(&buf)
	err = cfg.Write(bio, "json")
	if err != nil {
		return err
	}
	err = bio.Flush()
	if err != nil {
		return err
	}
	//return c.Send("OnParsingCommand")
	start := "```"
	end := "```"
	tpl := `
%s
%s
%s
`
	msg := fmt.Sprintf(tpl, start, buf.String(), end)
	return c.Edit(msg, &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
}
