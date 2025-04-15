package summer

import (
	"context"
	"github.com/hunterhug/golog"
	"google.golang.org/grpc/metadata"
	"path/filepath"
	"sync"
	"time"
)

var (
	// 日志和中间件单例模式，避免重复初始化
	hasInitLog     bool
	hasInitLogLock sync.Mutex

	// LogMaxAge 日志保存时间
	LogMaxAge = 30 * 24 * time.Hour
	// LogRotationTime 日志切割周期
	LogRotationTime = 24 * time.Hour
)

// 自定义的上下文日志打印
func addContextField(ctx context.Context, fields map[string]interface{}) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		// 加上自定义追踪ID
		fields["service.trace.id"] = md.Get(ContextTraceId)

		// 加上应用ID
		temp1 := md.Get(ContextAppId)
		if len(temp1) > 0 {
			fields["service.app.id"] = temp1[0]
		}

		// 加上管理员标志
		temp2 := md.Get(ContextAppAdmin)
		if len(temp2) > 0 {
			fields["service.app.admin"] = temp2[0]
		}

		// 加上用户ID
		temp3 := md.Get(ContextUserId)
		if len(temp3) > 0 {
			fields["service.app.user.id"] = temp3[0]
		}
	}
}

// 初始化日志
func initLog(appName string) error {
	hasInitLogLock.Lock()
	defer hasInitLogLock.Unlock()

	if hasInitLog {
		return nil
	}

	golog.SetLevel(golog.StringLevel(LogLevel))
	golog.AddFieldFunc(addContextField)
	golog.SetName(appName)
	if Debug {
		golog.SetIsOutputStdout(true)
	}

	if LogJson {
		golog.SetOutputJson(true)
	}

	if LogPath == "" || LogName == "" {
		hasInitLog = true
		return nil
	}

	golog.SetFileRotate(LogMaxAge, LogRotationTime)
	golog.SetOutputFile(filepath.Join(LogPath, appName), LogName)
	golog.InitLogger()
	hasInitLog = true
	return nil
}
