package middleware

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"
	"webhook-basket/model"

	"github.com/gin-gonic/gin"
)

func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer

		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		c.Request.Body = ioutil.NopCloser(&buf)

		// log.Println(string(body))
		// log.Println(c.Request.Header)

		header := time.Now().Format("2006-01-02 15:04:05") + " - "
		header += c.Request.RemoteAddr + " - "
		header += c.Request.Method + " - " + c.Request.URL.Path + " - "
		header += c.Request.UserAgent() + " - "
		body = bytes.Replace(body, []byte("\r\n"), []byte(""), -1)
		body = bytes.Replace(body, []byte("\n"), []byte(""), -1)
		data := string(header) + string(body) + "\n"

		if _, err := model.FileRequests.WriteString(data); err != nil {
			log.Println(err)
			return
		}

		c.Next()
	}
}

func LogFormatter(param gin.LogFormatterParams) string {
	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
		param.ClientIP,
		param.TimeStamp.Format(time.RFC1123),
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}
