package summer

import (
	"context"
	log "github.com/hunterhug/golog"
	"github.com/urfave/cli/v2"
	"os"
)

// Config 初始化配置
type Config struct {
	// 应用名
	AppName string
	// 应用使用介绍
	AppUsage string
	// 加载自定义配置
	LoadDiyFlagConfig func(map[string]cli.Flag)
	// 加载其他初始化
	LoadInitPrepare func() error
}

// Init 初始化准备
func Init(config *Config) error {
	log.SetName(config.AppName)
	log.InitLogger()

	ctx := context.Background()

	// 解析参数
	app := cli.NewApp()
	app.Name = config.AppName
	app.Usage = config.AppUsage
	flagMap := map[string]cli.Flag{}

	InitConfig(flagMap)
	config.LoadDiyFlagConfig(flagMap)

	flagList := make([]cli.Flag, 0)
	for _, v := range flagMap {
		flagList = append(flagList, v)
	}

	app.Flags = flagList
	app.Action = func(c *cli.Context) error {
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.ErrorContext(ctx, "Init: initConfig err:%s", err.Error())
		return err
	}

	log.InfoContext(ctx, "Init: initConfig success")

	// 初始化日志
	err = initLog(config.AppName)
	if err != nil {
		log.ErrorContext(ctx, "Init: initLog err:%s", err.Error())
		return err
	}

	log.InfoContext(ctx, "Init: initLog success")

	log.InfoContextWithFields(ctx, map[string]interface{}{
		"service.globalConfig": map[string]interface{}{
			"Debug":                       Debug,
			"Environment":                 Environment,
			"ListeningAddress":            ListeningAddress,
			"ListeningHTTPGateWay":        ListeningHTTPGateWay,
			"ListeningHTTPGateWayAddress": ListeningHTTPGateWayAddress,
			"ListeningHTTPProf":           ListeningHTTPProf,
			"ListeningHTTPProfAddress":    ListeningHTTPProfAddress,
			"LogLevel":                    LogLevel,
			"LogPath":                     LogPath,
			"LogName":                     LogName,
			"LogJson":                     LogJson,
		},
	}, "Init: initGlobalConfig")

	err = config.LoadInitPrepare()
	if err != nil {
		return err
	}
	return nil
}
