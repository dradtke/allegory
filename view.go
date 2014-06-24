package gopher

type View interface {
    HandleEvent(msg interface{}) bool
}
