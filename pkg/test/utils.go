package test

import (
	"emperror.dev/errors"
	"os"
	"testing"
	"time"
)

func SkipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
		return
	}
}

func WaitUntilConditionMet(conditionToMet func() bool) error {
	timeout := 20 * time.Second

	startTime := time.Now()
	timeOutExpired := false
	meet := conditionToMet()
	for meet == false {
		if timeOutExpired {
			return errors.New("Condition not met for the test, timeout exceeded")
		}
		time.Sleep(time.Second * 2)
		meet = conditionToMet()
		timeOutExpired = time.Now().Sub(startTime) > timeout
	}

	return nil
}
