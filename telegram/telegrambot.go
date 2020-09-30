package telegrambot

import (
	"context"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/pelletier/go-toml"
	"github.com/toffguy77/budni/feeds/rss"
	types "github.com/toffguy77/budni/internal"
	"go.uber.org/zap"
)

// Start ...
func Start(ctx context.Context) {
	logger := ctx.Value(types.ZapLogger("logger")).(*zap.SugaredLogger)
	cfg := ctx.Value(types.CfgContextKey("config")).(*toml.Tree)
	token := cfg.Get("bot.token").(string)

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Errorf("telegrambot: ", err)
	}

	logger.Infof("telegrambot: authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	// TODO: check for /start command at first >> remember update.Message.Chat.ID

	ProcessUpdates(ctx, bot, updates)
}

// ProcessUpdates ...
func ProcessUpdates(ctx context.Context, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	logger := ctx.Value(types.ZapLogger("logger")).(*zap.SugaredLogger)
	for update := range updates {
		chRssForExit := make(chan int)
		chRss := rss.MakeFeed(ctx, chRssForExit)

		for {
			select {
			case x := <-chRss:
				// TODO: check if message was already sent before
				logger.Infof("news: %v", x)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, x.Title+"\n\n"+x.Link)
				bot.Send(msg)
			}
		}
		close(chRssForExit)

		// логируем от кого какое сообщение пришло
		logger.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)
	}
}
