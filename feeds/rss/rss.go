// Package rss https://github.com/SlyMarbo/rss
package rss

import (
	"context"
	"sync"
	"time"

	"github.com/SlyMarbo/rss"
	"github.com/pelletier/go-toml"
	types "github.com/toffguy77/budni/internal"
	"go.uber.org/zap"
)

// MakeFeed ...
func MakeFeed(ctx context.Context, exitChan chan int) chan *rss.Item {
	logger := ctx.Value(types.ZapLogger("logger")).(*zap.SugaredLogger)
	cfg := ctx.Value(types.CfgContextKey("config")).(*toml.Tree)

	rssFeeds := cfg.Get("feeds.rss.links").([]interface{})
	logger.Infof("rss: got the following rss feeds to fetch: %v", rssFeeds)
	rssChan := make(chan *rss.Item)
	var wg sync.WaitGroup
	for _, rssFeed := range rssFeeds {
		wg.Add(1)
		go func(ctx context.Context, rssFeed string) {
			defer wg.Done()
			for {
				select {
				case <-exitChan:
					return
				default:
					item := readFead(ctx, rssFeed)
					logger.Infof("items from feed %v", rssFeed)
					rssChan <- item
				}
			}
		}(ctx, rssFeed.(string))
	}
	wg.Wait()
	return rssChan
}

func readFead(ctx context.Context, rssFeed string) *rss.Item {
	logger := ctx.Value(types.ZapLogger("logger")).(*zap.SugaredLogger)
	cfg := ctx.Value(types.CfgContextKey("config")).(*toml.Tree)
	timeout := cfg.Get("feeds.rss.timeout").(int64)

	// TODO: check if rssFeed is a http URL and RSS

	feed, err := rss.Fetch(rssFeed)
	if err != nil {

	}

	var item *rss.Item

	for {
		for _, item = range feed.Items {
			// TODO: use goroutine and WaitGroup
			news := func(item *rss.Item) *rss.Item {
				logger.Infof("one more news: %v, %v\n", item.Read, item.Title)
				return item
			}(item)
			return news
		}
		time.Sleep(time.Duration(timeout) * time.Second)
		logger.Info("rss: gone for an update")
		err = feed.Update()
		if err != nil {
			// pass
		}
	}
	return nil
}
