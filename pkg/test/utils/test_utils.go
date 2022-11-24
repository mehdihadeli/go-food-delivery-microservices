package testUtils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
)

func SkipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
		return
	}
}

func WaitUntilConditionMet(conditionToMet func() bool, timeout ...time.Duration) error {
	timeOutTime := 20 * time.Second
	if len(timeout) >= 0 && timeout != nil {
		timeOutTime = timeout[0]
	}

	startTime := time.Now()
	timeOutExpired := false
	meet := conditionToMet()
	for meet == false {
		if timeOutExpired {
			return errors.New("Condition not met for the test, timeout exceeded")
		}
		time.Sleep(time.Second * 2)
		meet = conditionToMet()
		timeOutExpired = time.Now().Sub(startTime) > timeOutTime
	}

	return nil
}

func HttpRecorder(t *testing.T, e *echo.Echo, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}
