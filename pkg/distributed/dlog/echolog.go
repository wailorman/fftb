package dlog

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// EchoLogger _
func EchoLogger(l logrus.FieldLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			logger := l

			method := req.Method
			logger = logger.WithField("method", method)

			status := res.Status
			logger = logger.WithField("status", status)

			path := req.URL.Path
			if path == "" {
				path = "/"
			}
			logger = logger.WithField("path", path)

			uri := req.RequestURI
			logger = logger.WithField("uri", uri)

			if requestID := req.Header.Get(echo.HeaderXRequestID); requestID != "" {
				logger = logger.WithField("request_id", requestID)
			}

			if remoteIP := c.RealIP(); remoteIP != "" {
				logger = logger.WithField("remote_ip", remoteIP)
			}

			if host := req.Host; host != "" {
				logger = logger.WithField("host", host)
			}

			if referer := req.Referer(); referer != "" {
				logger = logger.WithField("referer", referer)
			}

			if protocol := req.Proto; protocol != "" {
				logger = logger.WithField("protocol", protocol)
			}

			if userAgent := req.UserAgent(); userAgent != "" {
				logger = logger.WithField("user_agent", userAgent)
			}

			duration := stop.Sub(start).String()
			logger = logger.WithField("duration", duration)

			// reqBodyReader, _ := req.GetBody()
			// rawReqBody, _ := io.ReadAll(reqBodyReader)
			// if rawReqBody != nil {
			// 	logger.WithField("req_body", string(rawReqBody))
			// }

			logger.Debugf("%d %s %s", status, method, uri)

			return nil
		}
	}
}
