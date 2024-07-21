package log

import (
	"fmt"
	"time"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// EchoLogger returns echo middleware which will log incoming requests.
func EchoLogger(l logger.Logger, opts ...Option) echo.MiddlewareFunc {
	cfg := config{}
	for _, opt := range opts {
		opt.apply(&cfg)
	}

	if cfg.Skipper == nil {
		cfg.Skipper = middleware.DefaultSkipper
	}

	requestMiddleware := middleware.RequestLoggerWithConfig(
		middleware.RequestLoggerConfig{
			Skipper:          cfg.Skipper,
			LogRequestID:     true,
			LogRemoteIP:      true,
			LogHost:          true,
			LogMethod:        true,
			LogURI:           true,
			LogUserAgent:     true,
			LogStatus:        true,
			LogError:         true,
			LogLatency:       true,
			LogContentLength: true,
			LogResponseSize:  true,

			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				l.Infow(
					fmt.Sprintf(
						"[Request Middleware] REQUEST: uri: %v, status: %v\n",
						v.URI,
						v.Status,
					),
					logger.Fields{
						"uri":           v.URI,
						"status":        v.Status,
						"id":            v.RequestID,
						"remote_ip":     v.RemoteIP,
						"host":          v.Host,
						"method":        v.Method,
						"user_agent":    v.UserAgent,
						"error":         v.Error,
						"latency":       v.Latency.Nanoseconds(),
						"latency_human": v.Latency.String(),
						"bytes_in":      v.ContentLength,
						"bytes_out":     v.ResponseSize,
					},
				)

				return nil
			},
		},
	)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cfg.Skipper(c) {
				return requestMiddleware(next)(c)
			}

			start := time.Now()

			err := requestMiddleware(next)(c)
			if err != nil {
				// handle echo error in this middleware and raise echo errorhandler func and our custom error handler
				// when we call c.Error more than once, `c.Response().Committed` becomes true and response doesn't write to client again in our error handler
				// Error will update response status with occurred error object status code
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			fields := map[string]interface{}{
				"remote_ip":  c.RealIP(),
				"latency":    time.Since(start).String(),
				"host":       req.Host,
				"request":    fmt.Sprintf("%s %s", req.Method, req.RequestURI),
				"status":     res.Status,
				"size":       res.Size,
				"user_agent": req.UserAgent(),
			}

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}
			fields["request_id"] = id

			n := res.Status
			switch {
			case n >= 500:
				l.Errorw(
					"EchoServer logger middleware: Server error",
					fields,
				)
			case n >= 400:
				l.Errorw(
					"EchoServer logger middleware: Client error",
					fields,
				)
			case n >= 300:
				l.Errorw(
					"EchoServer logger middleware: Redirection",
					fields,
				)
			default:
				l.Infow("EchoServer logger middleware: Success", fields)
			}

			return nil
		}
	}
}
