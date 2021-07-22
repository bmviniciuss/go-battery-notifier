package reader

import "github.com/bmviniciuss/go-battery-notifier/domain"

type BatteryReader interface {
	Read() (*domain.BatteryStatus, error)
}
