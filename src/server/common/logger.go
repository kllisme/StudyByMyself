package common

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/kataras/iris.v5"
)

func SetupLogger() {

	isDevelopment := viper.GetBool("isDevelopment")

	appName := viper.GetString("name")

	if isDevelopment {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetOutput(os.Stderr)
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logFilePath := viper.GetString("server.log.path")
		logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
		if err != nil {
			logrus.Fatalf("open file error :%s \n", logFilePath)
			TeardownLogger()
		}
		logrus.SetLevel(logrus.WarnLevel)
		logrus.SetOutput(logFile)
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	Logger = logrus.WithFields(logrus.Fields{
		"system": appName,
	})
}

func TeardownLogger() {
	logrus.Printf("logger file stream closed.")
	logFile.Close()
}

var (
	Logger  *logrus.Entry
	logFile *os.File
	Log     = func(ctx *iris.Context, result *Result) {
		var processTime int64
		startAt := ctx.Get("startAt")
		if startAt != nil {
			startAt := startAt.(int64)
			endAt := time.Now().UnixNano() / 1000000
			processTime = endAt - startAt
		} else {
			processTime = -1
		}

		body := string(ctx.Request.Body()[:])

		alarmID := "0"

		handle := strings.Split(ctx.GetHandlerName(), "/")
		_interface := handle[1] + ":" + handle[len(handle)-2] + ":" + handle[len(handle)-1]
		_status := -1
		var err error
		if len(result.Status) >= 2 {
			_status, err = strconv.Atoi(result.Status[len(result.Status)-2 : len(result.Status)])
		} else {
			_status, err = strconv.Atoi(result.Status)
		}
		if _status != 0 && err == nil {
			_interface = "error:" + _interface
			alarmID = "1"
		} else {
			result.Data = struct{}{}
		}

		userId, _ := ctx.Session().GetInt(viper.GetString("server.session.user.id"))
		appName := viper.GetString("name")
		Logger := logrus.WithFields(logrus.Fields{
			"@source":    ctx.LocalIP().String(),
			"@timestamp": time.Now().Format(time.RFC3339),
			"@fields": map[string]interface{}{
				"userId":      userId,
				"fromtype":    appName,
				"host":        ctx.HostString(),
				"interface":   _interface,
				"method":      ctx.MethodString(),
				"ip":          ctx.RemoteAddr(),
				"query":       ctx.URLParams(),
				"param":       ctx.ParamsSentence(),
				"body":        body,
				"alarmID":     alarmID,
				"path":        ctx.PathString(),
				"processTime": processTime,
				"result":      result,
				"msg":         result.Msg,
				"data":        result.Data,
				"status":      result.Status,
				"exception":   result.Exception,
				"appCode":     appName + ":" + result.Code,
				"system":      appName,
				"totype":      appName,
			},
		})
		Logger.Warningln(result.Status)
	}
)
