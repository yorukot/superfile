package processbar

type ProcessChannelFullError struct {
}

func (p *ProcessChannelFullError) Error() string {
	return "process channel is full"
}

type NoProcessFoundError struct {
	id string
}

func (p *NoProcessFoundError) Error() string {
	return "no process with id : " + p.id
}

type ProcessAlreadyExistsError struct {
	id string
}

func (p *ProcessAlreadyExistsError) Error() string {
	return "process already exists with id : " + p.id
}
