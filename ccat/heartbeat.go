package ccat

type Heartbeat struct {
	Message
}

func (e *Heartbeat) Complete() {
	e.Message.flush(e)
}
