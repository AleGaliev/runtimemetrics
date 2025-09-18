package retry

import (
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"net"
	"time"
)

type Retry struct {
	RetryCount int
	Interval   []time.Duration
}

func CreateRetry() Retry {
	return Retry{
		RetryCount: 3,
		Interval:   []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second},
	}
}

func CheckConectionProblem(err error) bool {
	if err != nil {
		return false
	}

	var conectionError net.Error
	if errors.As(err, &conectionError) {
		return conectionError.Timeout()
	}

	var postgresErr *pgconn.PgError
	if errors.As(err, &postgresErr) {
		return pgerrcode.IsConnectionException(postgresErr.Code)
	}
	return false
}

func (retry *Retry) RetryConnection(dbFunction func() error) error {
	var err error
	for i := 0; i <= retry.RetryCount; i++ {
		err = dbFunction()
		if err == nil {
			return nil
		}

		if !CheckConectionProblem(err) {
			return fmt.Errorf("connection problem: %v", err)
		}
		time.Sleep(retry.Interval[i])
	}
	if err != nil {
		return fmt.Errorf("max retry connect: %v", err)
	}
	return nil
}
