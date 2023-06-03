package bot

import telebot "github.com/jxo-me/gfbot"

var (
	insBotRouter = Routers{
		List: map[string]telebot.HandlerFunc{
			telebot.OnText:       Event.OnText,
			telebot.OnCallback:   Event.OnCallback,
			telebot.OnUserJoined: Event.OnUserJoined,
			// Click Callback
			OnClickAdmissions:  Event.OnClickService,
			OnClickAuthers:     Event.OnClickService,
			OnClickBypass:      Event.OnClickService,
			OnClickHops:        Event.OnClickService,
			OnClickIngress:     Event.OnClickService,
			OnClickServices:    Event.OnClickService,
			OnClickChains:      Event.OnClickService,
			OnClickHosts:       Event.OnClickService,
			OnClickResolver:    Event.OnClickService,
			OnClickLimiter:     Event.OnClickService,
			OnClickConnLimiter: Event.OnClickService,
			OnClickRateLimiter: Event.OnClickService,
			OnClickConfig:      Event.OnClickService,
			// back services
			OnBackServices: Event.OnBackServices,
			// TextCommand
			HostTextCommand: Event.OnHostTextCommand,
		},
		Btns: map[*telebot.Btn]telebot.HandlerFunc{
			//&bot.BtnBetting:       bot.Event.OnBtnBetting,
		},
	}
)

type Routers struct {
	List map[string]telebot.HandlerFunc
	Btns map[*telebot.Btn]telebot.HandlerFunc
}

func Router() *Routers {
	return &insBotRouter
}
