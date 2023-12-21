package bot

import (
	"fmt"
	telebot "github.com/jxo-me/gfbot"
	"strings"
)

var (
	Event              = hEvent{}
	StartTextCommand   = "/start"
	NodeTextCommand    = "/myHosts"
	ParsingTextCommand = "/parsing"
	GostTextCommand    = "/gost"
	WebAppTextCommand  = "/webApp"
	GameTextCommand    = "/game"
	ContactTextCommand = "/contact"
	VenueTextCommand   = "/venue"
	MenuTextCommand    = "/menu"
	// Click group
	OnClickAdmissions        = "\fAdmissions"
	OnClickAddAdmission      = "\faddAdmission"
	OnClickUpdateAdmission   = "\fupdateAdmission"
	OnClickDetailAdmission   = "\fdetailAdmission"
	OnClickDelAdmission      = "\fdelAdmission"
	OnClickAuthers           = "\fAuthers"
	OnClickAddAuther         = "\faddAuther"
	OnClickUpdateAuther      = "\fupdateAuther"
	OnClickDetailAuther      = "\fdetailAuther"
	OnClickDelAuther         = "\fdelAuther"
	OnClickBypass            = "\fBypass"
	OnClickAddBypass         = "\faddBypass"
	OnClickUpdateBypass      = "\fupdateBypass"
	OnClickDetailBypass      = "\fdetailBypass"
	OnClickDelBypass         = "\fdelBypass"
	OnClickHops              = "\fHops"
	OnClickAddHop            = "\faddHop"
	OnClickUpdateHop         = "\fupdateHop"
	OnClickDetailHop         = "\fdetailHop"
	OnClickDelHop            = "\fdelHop"
	OnClickIngress           = "\fIngress"
	OnClickAddIngress        = "\faddIngress"
	OnClickUpdateIngress     = "\fupdateIngress"
	OnClickDetailIngress     = "\fdetailIngress"
	OnClickDelIngress        = "\fdelIngress"
	OnClickServices          = "\fServices"
	OnClickAddService        = "\faddService"
	OnClickUpdateService     = "\fupdateService"
	OnClickDetailService     = "\fdetailService"
	OnClickDelService        = "\fdelService"
	OnClickChains            = "\fChains"
	OnClickAddChain          = "\faddChain"
	OnClickDetailChain       = "\fdetailChain"
	OnClickUpdateChain       = "\fupdateChain"
	OnClickDelChain          = "\fdelChain"
	OnClickHosts             = "\fHosts"
	OnClickAddHosts          = "\faddHosts"
	OnClickUpdateHosts       = "\fupdateHosts"
	OnClickDetailHosts       = "\fdetailHosts"
	OnClickDelHosts          = "\fdelHosts"
	OnClickResolvers         = "\fResolvers"
	OnClickAddResolver       = "\faddResolver"
	OnClickUpdateResolver    = "\fupdateResolver"
	OnClickDetailResolver    = "\fdetailResolver"
	OnClickDelResolver       = "\fdelResolver"
	OnClickRouters           = "\fRouters"
	OnClickAddRouter         = "\faddRouter"
	OnClickUpdateRouter      = "\fupdateRouter"
	OnClickDetailRouter      = "\fdetailRouter"
	OnClickDelRouter         = "\fdelRouter"
	OnClickLimiters          = "\fLimiters"
	OnClickAddLimiter        = "\faddLimiter"
	OnClickUpdateLimiter     = "\fupdateLimiter"
	OnClickDetailLimiter     = "\fdetailLimiter"
	OnClickDelLimiter        = "\fdelLimiter"
	OnClickConnLimiters      = "\fCLimiters"
	OnClickAddConnLimiter    = "\faddConnLimiter"
	OnClickUpdateConnLimiter = "\fupdateConnLimiter"
	OnClickDetailConnLimiter = "\fdetailConnLimiter"
	OnClickDelConnLimiter    = "\fdelConnLimiter"
	OnClickRateLimiters      = "\fRLimiters"
	OnClickAddRateLimiter    = "\faddRateLimiter"
	OnClickUpdateRateLimiter = "\fupdateRateLimiter"
	OnClickDetailRateLimiter = "\fdetailRateLimiter"
	OnClickDelRateLimiter    = "\fdelRateLimiter"
	OnClickConfig            = "\fConfig"
	OnClickSaveConfig        = "\fsaveConfig"
	OnClickNode              = "\fNode"
	OnClickAddNode           = "\fAddNode"

	OnBackServices = "\fbackServices"
	OnBackHosts    = "\fbackHosts"
)

type (
	hEvent struct{}
)

func (h *hEvent) OnText(c telebot.IContext) error {
	btn := &telebot.MenuButton{
		Type:   telebot.MenuButtonWebApp,
		Text:   "WebApp",
		WebApp: &telebot.WebApp{URL: "https://dev.us.jxo.me"},
	}
	_ = c.Bot().SetMenuButton(c.Sender(), btn)
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
			selector.Data("@Resolvers", "Resolvers", "Resolvers"),
		),
		selector.Row(
			selector.Data("@Limiters", "Limiters", "Limiters"),
			selector.Data("@CLimiters", "CLimiters", "CLimiters"),
			selector.Data("@RLimiters", "RLimiters", "RLimiters"),
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
	return c.Edit("从下面的列表中选择一个服务:", &telebot.SendOptions{Protected: true, ReplyMarkup: selector})
}

func (h *hEvent) OnBackHosts(c telebot.IContext) error {
	selector := getSelectHosts()
	return c.Edit("从下面的列表中选择一个节点:", &telebot.SendOptions{Protected: true, ReplyMarkup: selector})
}

func (h *hEvent) OnClickNode(c telebot.IContext) error {
	selector := getSelectHosts()
	if c.Callback() != nil {
		return c.Edit("从下面的列表中选择一个服务:", &telebot.SendOptions{Protected: true, ReplyMarkup: selector})
	}
	return c.Send("从下面的列表中选择一个服务:", &telebot.SendOptions{Protected: true, ReplyMarkup: selector})
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

func (h *hEvent) OnWebAppCommand(c telebot.IContext) error {
	selector := &telebot.ReplyMarkup{}
	selector.Inline(selector.Row(selector.WebApp("Press me", &telebot.WebApp{URL: "https://dev.us.jxo.me"})))
	return c.Reply(
		fmt.Sprintf("Hello, I'm @%s.\nYou can use me to run a (very) simple telegram webapp demo!",
			c.Message().Sender.Username), &telebot.SendOptions{
			ReplyMarkup: selector,
		},
	)
}

func (h *hEvent) OnGameTextCommand(c telebot.IContext) error {
	//cmd := c.Callback().Data
	_ = c.Send("OnGameTextCommand")
	btnPlayGame := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{{{Text: "开始游戏", CallbackGame: &telebot.CallbackGame{}}}},
	}
	game := &telebot.Game{Name: "HGsport"}

	return c.Send(game, &telebot.SendOptions{ReplyMarkup: btnPlayGame})
}

func (h *hEvent) OnContactTextCommand(c telebot.IContext) error {
	//phone := c.Message().Text
	_ = c.Send("OnContactTextCommand")

	contact := &telebot.Contact{PhoneNumber: "+19106353888", FirstName: "Mickey"}

	return c.Send(contact, &telebot.SendOptions{})
}

func (h *hEvent) OnVenueTextCommand(c telebot.IContext) error {
	//phone := c.Message().Text
	_ = c.Send("OnVenueTextCommand")

	venue := &telebot.Venue{
		Location: telebot.Location{Lat: 1.3109116, Lng: 103.7536497},
		Title:    "Fischerbell",
		Address:  "79 Pandan Loop, Singapore 128282",
	}
	return c.Send(venue, &telebot.SendOptions{})
}

func (h *hEvent) OnMyChatMember(c telebot.IContext) error {
	// 邀请bot加入channel
	chatMember := c.ChatMember()
	oldRole := chatMember.OldChatMember.Role
	newRole := chatMember.NewChatMember.Role
	if oldRole == "left" && newRole == "administrator" {
		fmt.Println("加入Channel", chatMember.Chat.Title)
	} else if oldRole == "administrator" && newRole == "left" {
		fmt.Println("移除Channel", chatMember.Chat.Title)
	}
	if newRole == telebot.Administrator {
		return c.Send(fmt.Sprintf("邀请bot 加入Channel:%s", chatMember.Chat.Title))
	}
	return nil
	//fmt.Println("邀请bot 加入Channel:", chatMember.Chat.Title)
	//return c.Send(fmt.Sprintf("邀请bot 加入Channel:%s", chatMember.Chat.Title))
}
