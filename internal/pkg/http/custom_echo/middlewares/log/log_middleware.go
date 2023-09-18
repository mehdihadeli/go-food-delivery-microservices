package log

import (
	"fmt"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"

	"github.com/labstack/echo/v4"
)

func EchoLogger(logger logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				// handle echo error in this middleware and raise echo errorhandler func and our custom error handler (problem details handler)
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
				logger.Errorw("Echo logger middleware: Server error", fields)
			case n >= 400:
				logger.Errorw("Echo logger middleware: Client error", fields)
			case n >= 300:
				logger.Errorw("Echo logger middleware: Redirection", fields)
			default:
				logger.Infow("Echo logger middleware: Success", fields)
			}

			return nil
		}
	}
}
