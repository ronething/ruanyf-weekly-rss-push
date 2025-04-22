package notifier

type Message struct {
	Title           string
	URL             string
	Recommendations string
}

type Notifier interface {
	Notify(msg Message) error
}
