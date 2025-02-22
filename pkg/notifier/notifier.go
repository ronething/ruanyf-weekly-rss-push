package notifier

type Message struct {
	Title string
	URL   string
}

type Notifier interface {
	Notify(msg Message) error
}
