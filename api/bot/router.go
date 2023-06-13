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
			OnClickDetailAdmission: telebot.HandlerFunc(Event.OnClickDetailAdmission),
			OnClickDelAdmission:    telebot.HandlerFunc(Event.OnClickDelAdmission),
			// OnClickAuthers group
			OnClickAuthers: telebot.HandlerFunc(Event.OnClickService),
			OnClickBypass:  telebot.HandlerFunc(Event.OnClickService),
			OnClickHops:    telebot.HandlerFunc(Event.OnClickService),
			OnClickIngress: telebot.HandlerFunc(Event.OnClickService),
			// OnClickServices group
			OnClickServices:      telebot.HandlerFunc(Event.OnClickServices),
			OnClickAddService:    AddServiceConversation(OnClickAddService, "/cancel"),
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
