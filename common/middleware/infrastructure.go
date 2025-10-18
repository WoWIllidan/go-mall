package middleware

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/WoWBytePaladin/go-mall/common/logger"
	"github.com/WoWBytePaladin/go-mall/common/util"
	"github.com/gin-gonic/gin"
)

// infrastructure 中存放项目运行需要的基础中间价

func StartTrace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.Request.Header.Get("traceid")
		pSpanId := c.Request.Header.Get("spanid")
		spanId := util.GenerateSpanID(c.Request.RemoteAddr)
		if traceId == "" { // 如果traceId 为空，证明是链路的发端，把它设置成此次的spanId，发端的spanId是root spanId
			traceId = spanId // trace 标识整个请求的链路, span则标识链路中的不同服务
		}
		c.Set("traceid", traceId)
		c.Set("spanid", spanId)
		c.Set("pspanid", pSpanId)
		c.Next()
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// 包装一下 gin.ResponseWriter，通过这种方式拦截写响应
// 让gin写响应的时候先写到 bodyLogWriter 再写gin.ResponseWriter ，
// 这样利用中间件里输出访问日志时就能拿到响应了
// https://stackoverflow.com/questions/38501325/how-to-log-response-body-in-gin
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LogAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		//保存body
		reqBody, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(reqBody))
		start := time.Now()
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		accessLog(c, "access_start", time.Since(start), reqBody, nil)
		defer func() {
			accessLog(c, "access_end", time.Since(start), reqBody, blw.body.String())
		}()
		c.Next()

		return
	}
}

func accessLog(c *gin.Context, accessType string, dur time.Duration, body []byte, dataOut interface{}) {
	req := c.Request
	bodyStr := string(body)
	query := req.URL.RawQuery
	path := req.URL.Path
	// TODO: 实现Token认证后再把访问日志里也加上token记录
	// token := c.Request.Header.Get("token")
	logger.New(c).Info("AccessLog",
		"type", accessType,
		"ip", c.ClientIP(),
		//"token", token,
		"method", req.Method,
		"path", path,
		"query", query,
		"body", bodyStr,
		"output", dataOut,
		"time(ms)", int64(dur/time.Millisecond))
}

// GinPanicRecovery 自定义gin recover输出
func GinPanicRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.New(c).Error("http request broken pipe", "path", c.Request.URL.Path, "error", err, "request", string(httpRequest))
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				logger.New(c).Error("http_request_panic", "path", c.Request.URL.Path, "error", err, "request", string(httpRequest), "stack", string(debug.Stack()))

				c.AbortWithError(http.StatusInternalServerError, err.(error))
			}
		}()
		c.Next()
	}
}
