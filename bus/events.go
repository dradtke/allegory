package bus

const (
	_ = EventId(^uint(0) - iota)

	// Handler signature: func(cmd string)
	ConsoleCommandEvent
)
