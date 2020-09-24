package config

import (
	"context"
	"io/ioutil"

	"github.com/pelletier/go-toml"
	types "github.com/toffguy77/budni/internal"
	"go.uber.org/zap"
)

// GetCfg loads the configuration from file
func GetCfg(ctx context.Context, path string) (*toml.Tree, error) {
	logger := ctx.Value(types.ZapLogger("logger")).(*zap.SugaredLogger)

	if path == "" {
		path = "config/config.toml"
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Infof("config: %v", err)
		return nil, err
	}

	cfg, err := toml.Load(string(content))
	if err != nil {
		logger.Infof("config: %v", err)
		return nil, err
	}
	logger.Infof("passed: %v", cfg)

	return cfg, nil
}
