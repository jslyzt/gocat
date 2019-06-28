package cat

import (
	"time"

	"github.com/jslyzt/gocat/ccat"
)

type catInstance struct {
}

func Instance() *catInstance {
	return &catInstance{}
}

func (t *catInstance) flush(m ccat.Messager) {
	ccat.Send(m)
}

func (t *catInstance) NewTransaction(mtype, name string) *ccat.Transaction {
	return ccat.NewTransaction(mtype, name, t.flush)
}

func (t *catInstance) NewCompletedTransactionWithDuration(mtype, name string, durationInNano int64) {
	var trans = t.NewTransaction(mtype, name)
	trans.SetDuration(durationInNano)
	if durationInNano > 0 && durationInNano < 60*time.Second.Nanoseconds() {
		trans.SetTimestamp(time.Now().UnixNano() - durationInNano)
	}
	trans.SetStatus(ccat.SUCCESS)
	trans.Complete()
}

func (t *catInstance) NewEvent(mtype, name string) *ccat.Event {
	return &ccat.Event{
		Message: *ccat.NewMessage(mtype, name, t.flush),
	}
}

func (t *catInstance) NewHeartbeat(mtype, name string) *ccat.Heartbeat {
	return &ccat.Heartbeat{
		Message: *ccat.NewMessage(mtype, name, t.flush),
	}
}

func (t *catInstance) LogEvent(mtype, name string, args ...string) {
	var e = t.NewEvent(mtype, name)
	if len(args) > 0 {
		e.SetStatus(args[0])
	}
	if len(args) > 1 {
		e.AddData(args[1])
	}
	e.Complete()
}

func (t *catInstance) LogError(err error, args ...string) {
	var category = "error"
	if len(args) > 0 {
		category = args[0]
	}
	var e = t.NewEvent("Exception", category)
	var buf = newStacktrace(2, err)
	e.SetStatus(ccat.FAIL)
	e.AddData(buf.String())
	e.Complete()
}

func (t *catInstance) LogMetricForCount(mname string, args ...int) {
	if len(args) == 0 {
		ccat.LogMetricForCount(mname, 1)
	} else {
		ccat.LogMetricForCount(mname, args[0])
	}
}

func (t *catInstance) LogMetricForDuration(mname string, duration int64) {
	ccat.LogMetricForDuration(mname, duration)
}
