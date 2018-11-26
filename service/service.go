package service

import (
	log "github.com/cihub/seelog"
	"github.com/daiguadaidai/blingbling/config"
	"github.com/daiguadaidai/blingbling/reviewserver"
	"sync"
)

func Run(_httpServerConfig *config.HttpServerConfig, _config *config.ReviewConfig) {
	defer log.Flush()
	logger, _ := log.LoggerFromConfigAsBytes([]byte(config.LogDefautConfig()))
	log.ReplaceLogger(logger)

	config.SetReviewConfig(_config) // 设置全局的sql review配置文件

	var wg sync.WaitGroup

	wg.Add(1)
	go reviewserver.StartReviewServer(_httpServerConfig, &wg)

	wg.Wait()
}
