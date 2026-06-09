package health

import "context"

type Checker struct {
	redisChecker    *RedisChecker
	postgresChecker *PostgresChecker
}

func NewChecker(redisChecker *RedisChecker, postgresChecker *PostgresChecker) *Checker {
	return &Checker{
		redisChecker:    redisChecker,
		postgresChecker: postgresChecker,
	}
}

const (
	StatusOK          = "ok"
	StatusDegraded    = "degraded"
	StatusUnavailable = "unavailable"

	CheckUp   = "up"
	CheckDown = "down"
)

func (c *Checker) Check(ctx context.Context) HealthReport {
	report := HealthReport{
		Checks: make(map[string]CheckStatus),
	}
	err := c.postgresChecker.Ping(ctx)
	if err != nil {
		report.Checks["postgres"] = CheckStatus{
			Status: CheckDown,
			Error:  err.Error(),
		}
	} else {
		report.Checks["postgres"] = CheckStatus{
			Status: CheckUp,
		}
	}
	err = c.redisChecker.Ping(ctx)
	if err != nil {
		report.Checks["redis"] = CheckStatus{
			Status: CheckDown,
			Error:  err.Error(),
		}
	} else {
		report.Checks["redis"] = CheckStatus{
			Status: CheckUp,
		}
	}
	report.Status = StatusOK

	if report.Checks["redis"].Status == CheckDown {
		report.Status = StatusDegraded
	}

	if report.Checks["postgres"].Status == CheckDown {
		report.Status = StatusUnavailable
	}

	return report
}
