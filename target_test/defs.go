package target_test

import "time"

// database model
type Todo struct {
	Id               int
	Id_target        int
	Name             string
	Period           int
	Last_test_time   time.Time
	Next_test_time   time.Time
	Next_test_server int
	Type             string
	Test             string
	Validation       string
}

type Target_test interface {
	New() Target_test
	Test() string
}
