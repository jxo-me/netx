package bot

import (
	"fmt"
	telebot "github.com/jxo-me/gfbot"
)

var (
	Event           = hEvent{}
	HostTextCommand = "/myHosts"
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

	OnBackServices = "\fbackServices"
)

type (
	hEvent struct{}
)

func (h *hEvent) OnText(c telebot.Context) error {
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
		),
	)
	return selector
}

func (h *hEvent) OnBackServices(c telebot.Context) error {
	selector := getSelectMenus()
	return c.Edit("从下面的列表中选择一个服务:", selector)
}

func (h *hEvent) OnHostTextCommand(c telebot.Context) error {
	selector := getSelectMenus()
	return c.Send("从下面的列表中选择一个服务:", selector)
}

func (h *hEvent) OnClickService(c telebot.Context) error {
	user := c.Callback().Sender
	cmd := c.Callback().Data
	msg := fmt.Sprintf("选中服务: %s %d.\nWhat do you want to do with the bot?", cmd, user.ID)
	selector := &telebot.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			selector.Data("@update", "update", "update"),
			selector.Data("@delete", "delete", "delete"),
			selector.Data("@list", "list", "list"),
		),
		selector.Row(selector.Data("« 返回服务列表", "backServices", "backServices")),
	)
	return c.Edit(msg, selector)
}

func (h *hEvent) OnCallback(c telebot.Context) error {
	return c.Send("OnCallback")
}

func (h *hEvent) OnUserJoined(c telebot.Context) error {
	return c.Send("OnUserJoined")
}
