package service

import (
	"github.com/daiguadaidai/blingbling/config"
	"sync"
	"github.com/daiguadaidai/blingbling/reviewserver"
	log "github.com/cihub/seelog"
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


