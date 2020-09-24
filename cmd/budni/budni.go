package main

import (
	"context"
	"log"

	"github.com/toffguy77/budni/config"
	types "github.com/toffguy77/budni/internal"
	bot "github.com/toffguy77/budni/telegram"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	ctx := context.WithValue(context.Background(), types.ZapLogger("logger"), sugar)

	cfg, err := config.GetCfg(ctx, "config/config.toml")
	if err != nil {
		log.Fatalf("budni: %v", err)
	}

	c := types.CfgContextKey("config")
	ctx = context.WithValue(ctx, c, cfg)

	//fmt.Println(cfg.GetArray("feeds.twitter.hashtags"))

	bot.Start(ctx)
}
