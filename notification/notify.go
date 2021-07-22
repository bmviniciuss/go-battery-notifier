package notification

import "github.com/gen2brain/beeep"

type Notifier interface {
	Notify(title string, message string) error
}

type BeeepNotifier struct {
	Icon string
}

func NewBeeepNotifier() *BeeepNotifier {
	return &BeeepNotifier{Icon: ""}
}

func (bn *BeeepNotifier) Notify(title string, message string) error {
	return beeep.Notify(title, message, "")
}
