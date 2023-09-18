package health

import "context"

func NewUnhealthyHealthService() UnhealthyHealthService {
	return UnhealthyHealthService{}
}

type UnhealthyHealthService struct{}

func (service UnhealthyHealthService) CheckHealth(context.Context) Check {
	return Check{
		"postgres": Status{Status: StatusDown},
		"redis":    Status{Status: StatusDown},
	}
}
