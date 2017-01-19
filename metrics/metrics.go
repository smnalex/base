package metrics

type Counter interface {
	Name() string
	With(string) Counter
	Add(delta uint64)
}
