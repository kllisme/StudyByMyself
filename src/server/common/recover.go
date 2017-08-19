package common

import (
	"fmt"
	"runtime"
	"strconv"

	"gopkg.in/kataras/iris.v5"
)

func getRequestLogs(ctx *iris.Context) string {
	var status, ip, method, path string
	status = strconv.Itoa(ctx.Response.StatusCode())
	path = ctx.PathString()
	method = ctx.MethodString()
	ip = ctx.RemoteAddr()
	return fmt.Sprintf("%v %s %s %s", status, path, method, ip)
}

func NewRecover() iris.HandlerFunc {
	return func(ctx *iris.Context) {
		defer func() {
			if err := recover(); err != nil {
				if ctx.IsStopped() {
					return
				}

				var stacktrace string
				for i := 1; ; i++ {
					_, f, l, got := runtime.Caller(i)
					if !got {
						break

					}

					stacktrace += fmt.Sprintf("%s:%d\n", f, l)
				}

				// when stack finishes
				logMessage := fmt.Sprintf("Recovered from a route's Handler('%s')\n", ctx.GetHandlerName())
				logMessage += fmt.Sprintf("At Request: %s\n", getRequestLogs(ctx))
				logMessage += fmt.Sprintf("Trace: %s\n", err)
				logMessage += fmt.Sprintf("\n%s\n", stacktrace)
				ctx.Log(logMessage)

				ctx.StopExecution()
				ctx.EmitError(iris.StatusInternalServerError)

				result := &Result{}
				_result := result.New("000002", nil)
				fmt.Println("isError:", _result.IsError)
				_result.Exception = logMessage
				_Log(ctx, _result)
			}
		}()

		ctx.Next()
	}
}
