package bus

const (
    _ uint = ^uint(0) - iota

    // Handler signature: func(cmd string)
    ConsoleCommandEvent
)
