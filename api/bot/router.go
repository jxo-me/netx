package bot

import (
	telebot "github.com/jxo-me/gfbot"
)

var (
	insBotRouter = Routers{
		List: map[string]telebot.IHandler{
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
			OnClickIngress:       telebot.HandlerFunc(Event.OnClickIngresses),
			OnClickAddIngress:    AddIngressConversation(OnClickAddIngress, "/cancel"),
			OnClickUpdateIngress: UpdateIngressConversation(OnClickUpdateIngress, "/cancel"),
			OnClickDetailIngress: telebot.HandlerFunc(Event.OnClickDetailIngress),
			OnClickDelIngress:    telebot.HandlerFunc(Event.OnClickDelIngress),
			// OnClickServices group
			OnClickServices:      telebot.HandlerFunc(Event.OnClickServices),
			OnClickAddService:    AddServiceConversation(OnClickAddService, "/cancel"),
			OnClickUpdateService: UpdateServiceConversation(OnClickUpdateService, "/cancel"),
			OnClickDetailService: telebot.HandlerFunc(Event.OnClickDetailService),
			OnClickDelService:    telebot.HandlerFunc(Event.OnClickDelService),
			// OnClickChains group
			OnClickChains:      telebot.HandlerFunc(Event.OnClickChains),
			OnClickAddChain:    AddChainConversation(OnClickAddChain, "/cancel"),
			OnClickUpdateChain: UpdateChainConversation(OnClickUpdateChain, "/cancel"),
			OnClickDetailChain: telebot.HandlerFunc(Event.OnClickDetailChain),
			OnClickDelChain:    telebot.HandlerFunc(Event.OnClickDelChain),
			// OnClickHosts group
			OnClickHosts:       telebot.HandlerFunc(Event.OnClickHosts),
			OnClickAddHosts:    AddHostsConversation(OnClickAddHosts, "/cancel"),
			OnClickUpdateHosts: UpdateHostsConversation(OnClickUpdateHosts, "/cancel"),
			OnClickDetailHosts: telebot.HandlerFunc(Event.OnClickDetailHosts),
			OnClickDelHosts:    telebot.HandlerFunc(Event.OnClickDelHosts),
			// OnClickResolver group
			OnClickResolvers:      telebot.HandlerFunc(Event.OnClickResolvers),
			OnClickAddResolver:    AddResolverConversation(OnClickAddResolver, "/cancel"),
			OnClickUpdateResolver: UpdateResolverConversation(OnClickUpdateResolver, "/cancel"),
			OnClickDetailResolver: telebot.HandlerFunc(Event.OnClickDetailResolver),
			OnClickDelResolver:    telebot.HandlerFunc(Event.OnClickDelResolver),
			// OnClickLimiter group
			OnClickLimiters:      telebot.HandlerFunc(Event.OnClickLimiters),
			OnClickAddLimiter:    AddLimiterConversation(OnClickAddLimiter, "/cancel"),
			OnClickUpdateLimiter: UpdateLimiterConversation(OnClickUpdateLimiter, "/cancel"),
			OnClickDetailLimiter: telebot.HandlerFunc(Event.OnClickDetailLimiter),
			OnClickDelLimiter:    telebot.HandlerFunc(Event.OnClickDelLimiter),
			// OnClickConnLimiter group
			OnClickConnLimiters:      telebot.HandlerFunc(Event.OnClickConnLimiters),
			OnClickAddConnLimiter:    AddConnLimiterConversation(OnClickAddConnLimiter, "/cancel"),
			OnClickUpdateConnLimiter: UpdateConnLimiterConversation(OnClickUpdateConnLimiter, "/cancel"),
			OnClickDetailConnLimiter: telebot.HandlerFunc(Event.OnClickDetailConnLimiter),
			OnClickDelConnLimiter:    telebot.HandlerFunc(Event.OnClickDelConnLimiter),
			// OnClickRateLimiter group
			OnClickRateLimiters:      telebot.HandlerFunc(Event.OnClickRateLimiters),
			OnClickAddRateLimiter:    AddRateLimiterConversation(OnClickAddRateLimiter, "/cancel"),
			OnClickUpdateRateLimiter: UpdateRateLimiterConversation(OnClickUpdateRateLimiter, "/cancel"),
			OnClickDetailRateLimiter: telebot.HandlerFunc(Event.OnClickDetailRateLimiter),
			OnClickDelRateLimiter:    telebot.HandlerFunc(Event.OnClickDelRateLimiter),
			// OnClickConfig group
			OnClickConfig:     telebot.HandlerFunc(Event.OnClickConfig),
			OnClickSaveConfig: telebot.HandlerFunc(Event.OnClickSaveConfig),

			OnClickNode:    telebot.HandlerFunc(Event.OnBackServices),
			OnClickAddNode: AddNodeConversation(OnClickAddNode, "/cancel"),
			// back services
			OnBackServices: telebot.HandlerFunc(Event.OnBackServices),
			OnBackHosts:    telebot.HandlerFunc(Event.OnClickNode),

			telebot.OnText:       telebot.HandlerFunc(Event.OnText),
			telebot.OnCallback:   telebot.HandlerFunc(Event.OnCallback),
			telebot.OnUserJoined: telebot.HandlerFunc(Event.OnUserJoined),

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
