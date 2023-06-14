package bot

import (
	telebot "github.com/jxo-me/gfbot"
)

var (
	insBotRouter = Routers{
		List: map[string]telebot.IHandler{
			telebot.OnText:       telebot.HandlerFunc(Event.OnText),
			telebot.OnCallback:   telebot.HandlerFunc(Event.OnCallback),
			telebot.OnUserJoined: telebot.HandlerFunc(Event.OnUserJoined),
			// Click Callback
			OnClickAdmissions:      telebot.HandlerFunc(Event.OnClickAdmissions),
			OnClickAddAdmission:    AddAdmissionConversation(OnClickAddAdmission, "/cancel"),
			OnClickUpdateAdmission: UpdateAdmissionConversation(OnClickUpdateAdmission, "/cancel"),
			OnClickDetailAdmission: telebot.HandlerFunc(Event.OnClickDetailAdmission),
			OnClickDelAdmission:    telebot.HandlerFunc(Event.OnClickDelAdmission),
			// OnClickAuthers group
			OnClickAuthers:      telebot.HandlerFunc(Event.OnClickAuthers),
			OnClickAddAuther:    AddAutherConversation(OnClickAddAuther, "/cancel"),
			OnClickUpdateAuther: UpdateAutherConversation(OnClickUpdateAuther, "/cancel"),
			OnClickDetailAuther: telebot.HandlerFunc(Event.OnClickDetailAuther),
			OnClickDelAuther:    telebot.HandlerFunc(Event.OnClickDelAuther),
			// OnClickBypass group
			OnClickBypass:       telebot.HandlerFunc(Event.OnClickBypasses),
			OnClickAddBypass:    AddBypassConversation(OnClickAddBypass, "/cancel"),
			OnClickUpdateBypass: UpdateBypassConversation(OnClickUpdateBypass, "/cancel"),
			OnClickDetailBypass: telebot.HandlerFunc(Event.OnClickDetailBypass),
			OnClickDelBypass:    telebot.HandlerFunc(Event.OnClickDelBypass),
			// OnClickHops group
			OnClickHops:      telebot.HandlerFunc(Event.OnClickHops),
			OnClickAddHop:    AddHopConversation(OnClickAddHop, "/cancel"),
			OnClickUpdateHop: UpdateHopConversation(OnClickUpdateHop, "/cancel"),
			OnClickDetailHop: telebot.HandlerFunc(Event.OnClickDetailHop),
			OnClickDelHop:    telebot.HandlerFunc(Event.OnClickDelHop),
			// OnClickIngress group
			OnClickIngress: telebot.HandlerFunc(Event.OnClickService),
			// OnClickServices group
			OnClickServices:      telebot.HandlerFunc(Event.OnClickServices),
			OnClickAddService:    AddServiceConversation(OnClickAddService, "/cancel"),
			OnClickUpdateService: UpdateServiceConversation(OnClickUpdateService, "/cancel"),
			OnClickDetailService: telebot.HandlerFunc(Event.OnClickDetailService),
			OnClickDelService:    telebot.HandlerFunc(Event.OnClickDelService),

			OnClickChains:   telebot.HandlerFunc(Event.OnClickService),
			OnClickHosts:    telebot.HandlerFunc(Event.OnClickService),
			OnClickResolver: telebot.HandlerFunc(Event.OnClickService),
			// OnClickLimiter group
			OnClickLimiter:     telebot.HandlerFunc(Event.OnClickService),
			OnClickConnLimiter: telebot.HandlerFunc(Event.OnClickService),
			OnClickRateLimiter: telebot.HandlerFunc(Event.OnClickService),
			// OnClickConfig group
			OnClickConfig:     telebot.HandlerFunc(Event.OnClickConfig),
			OnClickSaveConfig: telebot.HandlerFunc(Event.OnClickSaveConfig),

			OnClickNode:    telebot.HandlerFunc(Event.OnBackServices),
			OnClickAddNode: AddNodeConversation(OnClickAddNode, "/cancel"),
			// back services
			OnBackServices: telebot.HandlerFunc(Event.OnBackServices),
			OnBackHosts:    telebot.HandlerFunc(Event.OnClickNode),
			// TextCommand
			StartTextCommand:   telebot.HandlerFunc(Event.OnStartCommand),
			NodeTextCommand:    telebot.HandlerFunc(Event.OnClickNode),
			ParsingTextCommand: telebot.HandlerFunc(Event.OnParsingCommand),
			GostTextCommand:    telebot.HandlerFunc(Event.OnGostCommand),
		},
		Btns: map[*telebot.Btn]telebot.IHandler{
			//&bot.BtnBetting:       bot.Event.OnBtnBetting,
		},
	}
)

type Routers struct {
	List map[string]telebot.IHandler
	Btns map[*telebot.Btn]telebot.IHandler
	Test map[string]telebot.IHandler
}

func Router() *Routers {
	return &insBotRouter
}
