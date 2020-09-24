package telegrambot

import (
	"context"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/pelletier/go-toml"
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
	bot.Self.FirstName = cfg.Get("bot.firstname").(string)
	bot.Self.LastName = cfg.Get("bot.lastname").(string)

	logger.Infof("telegrambot: authorized on account %s", bot.Self.UserName)

	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates, err := bot.GetUpdatesChan(u)

	// в канал updates прилетают структуры типа Update
	// вычитываем их и обрабатываем
	for update := range updates {
		// универсальный ответ на любое сообщение
		reply := "Не знаю что сказать"
		if update.Message == nil {
			continue
		}

		// логируем от кого какое сообщение пришло
		logger.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// свитч на обработку комманд
		// комманда - сообщение, начинающееся с "/"
		switch update.Message.Command() {
		case "start":
			reply = "Привет. Я телеграм-бот"
		case "hello":
			reply = "world"
		}

		// создаем ответное сообщение
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		// отправляем
		bot.Send(msg)
	}
}
