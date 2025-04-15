package summer

import (
	"github.com/urfave/cli/v2"
)

// 全局配置
var (
	// Debug 调试模式
	Debug = false

	// ListeningAddress TCP监听地址
	ListeningAddress = ":10000"

	// ListeningHTTPGateWayAddress HTTP监听地址
	ListeningHTTPGateWayAddress = ":10001"
	ListeningHTTPGateWay        = false

	// ListeningHTTPProfAddress HTTP监听地址
	ListeningHTTPProfAddress = ":10002"
	ListeningHTTPProf        = false

	// LogLevel LogMaxAge 日志配置
	LogLevel = "debug"   // 日志级别
	LogJson  = false     // 日志输出格式
	LogPath  = "./log"   // 日志路径
	LogName  = "log.log" // 日志文件名

	// Environment 哪个环境
	Environment = "Development"
)

// InitConfig 初始化配置
func InitConfig(flagMap map[string]cli.Flag) {
	flagMap["debug"] = &cli.BoolFlag{
		Name:        "debug",
		Value:       Debug,
		Usage:       "debug mode",
		Destination: &Debug,
		EnvVars:     []string{"debug"},
	}

	flagMap["address"] = &cli.StringFlag{
		Name:        "address",
		Value:       ListeningAddress,
		Usage:       "gRPC Address listened by this program, format IP:PORT",
		Destination: &ListeningAddress,
		EnvVars:     []string{"address"},
	}

	flagMap["open_http_gateway"] = &cli.BoolFlag{
		Name:        "open_http_gateway",
		Value:       ListeningHTTPGateWay,
		Usage:       "open http gateway Address listened",
		Destination: &ListeningHTTPGateWay,
		EnvVars:     []string{"open_http_gateway"},
	}

	flagMap["http_gateway_address"] = &cli.StringFlag{
		Name:        "http_gateway_address",
		Value:       ListeningHTTPGateWayAddress,
		Usage:       "http gateway Address listened by this program, format IP:PORT",
		Destination: &ListeningHTTPGateWayAddress,
		EnvVars:     []string{"http_gateway_address"},
	}

	flagMap["open_http_pprof"] = &cli.BoolFlag{
		Name:        "open_http_pprof",
		Value:       ListeningHTTPProf,
		Usage:       "open http PProf Address listened",
		Destination: &ListeningHTTPProf,
		EnvVars:     []string{"open_http_pprof"},
	}

	flagMap["http_pprof_address"] = &cli.StringFlag{
		Name:        "http_pprof_address",
		Value:       ListeningHTTPProfAddress,
		Usage:       "http PProf Address listened by this program, format IP:PORT",
		Destination: &ListeningHTTPProfAddress,
		EnvVars:     []string{"http_pprof_address"},
	}

	flagMap["log_level"] = &cli.StringFlag{
		Name:        "log_level",
		Value:       LogLevel,
		Usage:       "Valid values: debug, info, warn, error",
		Destination: &LogLevel,
		EnvVars:     []string{"log_level"},
	}

	flagMap["log_path"] = &cli.StringFlag{
		Name:        "log_path",
		Value:       LogPath,
		Usage:       "Log file path, if the file path does not exist, it will be created",
		Destination: &LogPath,
		EnvVars:     []string{"log_path"},
	}

	flagMap["log_json"] = &cli.BoolFlag{
		Name:        "log_json",
		Value:       LogJson,
		Usage:       "Log to json",
		Destination: &LogJson,
		EnvVars:     []string{"log_json"},
	}

	flagMap["log_name"] = &cli.StringFlag{
		Name:        "log_name",
		Value:       LogName,
		Usage:       "Log file name",
		Destination: &LogName,
		EnvVars:     []string{"log_name"},
	}

	flagMap["environment"] = &cli.StringFlag{
		Name:        "environment",
		Value:       Environment,
		Usage:       "e.g. test or production",
		Destination: &Environment,
		EnvVars:     []string{"environment"},
	}
}
