package application

import (
	"errors"
	"fmt"
	"github.com/bmviniciuss/go-battery-notifier/domain"
	"github.com/bmviniciuss/go-battery-notifier/notification"
	"github.com/bmviniciuss/go-battery-notifier/reader"
	"log"
	"time"
)

type Application struct {
	Reader   reader.BatteryReader
	Notifier notification.Notifier
	title    string
	UnplugAt uint8
	PlugAt   uint8
	Interval int // in seconds
	Debug    bool
}

func NewApplication(reader reader.BatteryReader, notifier notification.Notifier, unplugAt uint8, plugAt uint8, interval int, debug bool) (*Application, error) {
	err := validateApplicationParams(unplugAt, plugAt, interval)
	if err != nil {
		return &Application{}, err
	}
	return &Application{
		Reader:   reader,
		Notifier: notifier,
		UnplugAt: unplugAt,
		PlugAt:   plugAt,
		Interval: interval,
		Debug:    debug,
		title:    "Go Battery Notifier",
	}, nil
}

func validateApplicationParams(unplugAt uint8, plugAt uint8, interval int) error {
	if interval <= 0 {
		return errors.New("interval must be a positive value")
	} else if unplugAt < 0 || unplugAt > 100 {
		return errors.New("unplug must be in the range [0-100]")
	} else if plugAt < 0 || plugAt > 100 {
		return errors.New("plug value must be in the range [0-100]")
	}
	return nil
}

func (a *Application) Run() error {
	a.debugPrintNow()
	var err error
	err = a.getReading()
	duration := time.Second * time.Duration(a.Interval)
	for range time.Tick(duration) {
		a.debugPrintNow()
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			err = a.getReading()
		}()
	}
	return nil
}

func (a *Application) getReading() error {
	status, err := a.getBatteryStatus()
	if err != nil {
		return err
	}
	a.interpretBatteryStatus(status)
	return nil
}

func (a *Application) getBatteryStatus() (*domain.BatteryStatus, error) {
	status, err := a.Reader.Read()
	if err != nil {
		return &domain.BatteryStatus{}, err
	}
	return status, nil
}

func (a *Application) interpretBatteryStatus(s *domain.BatteryStatus) {
	if a.Debug {
		fmt.Printf("DEBUG: %v\n", s)
	}
	if s.State == domain.Charging && s.Level >= a.UnplugAt {
		unplugMessage := fmt.Sprintf("Your battery is already at %d%%! Please unplug now.", a.UnplugAt)
		if a.Debug {
			fmt.Printf("Attention! %s\n", unplugMessage)
		}
		a.Notifier.Notify(a.title, unplugMessage)
	} else if s.State == domain.Discharging && s.Level <= a.PlugAt {
		plugMessage := fmt.Sprintf("Your battery is already at %d%%! Please plug now.", a.PlugAt)
		if a.Debug {
			fmt.Printf("Attention! You battery is already at %d%%. Please plug now.\n", a.PlugAt)
		}
		a.Notifier.Notify(a.title, plugMessage)
	}
}

func (a *Application) debugPrintNow() {
	if a.Debug {
		fmt.Println(time.Now())
	}
}
