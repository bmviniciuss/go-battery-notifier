package reader

import (
	"bytes"
	"github.com/bmviniciuss/go-battery-notifier/domain"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type UbuntuReader struct{}

func NewUbuntuReader() *UbuntuReader {
	return &UbuntuReader{}
}

func (u *UbuntuReader) Read() (*domain.BatteryStatus, error) {
	raw, err := u.getRawDataFromSystem()
	if err != nil {
		return &domain.BatteryStatus{}, err
	}

	b, err := u.parseRawDataToDomain(raw)
	if err != nil {
		return &domain.BatteryStatus{}, err
	}
	return b, nil
}

func (u *UbuntuReader) getRawDataFromSystem() (string, error) {
	var out bytes.Buffer
	c := exec.Command("acpi", "-b")
	c.Stdout = &out
	err := c.Run()
	if err != nil {
		return "", nil
	}
	return out.String(), nil
}

func (u *UbuntuReader) parseRawDataToDomain(rawData string) (*domain.BatteryStatus, error) {
	re := regexp.MustCompile(`Battery 0: (?P<status>\w+), (?P<level>\d+)%`)
	c := re.FindStringSubmatch(rawData)

	stateString := c[1]
	level, err := strconv.Atoi(c[2])

	if err != nil {
		return &domain.BatteryStatus{}, err
	}

	batteryState := u.getBatteryStateFromSystemStr(stateString)
	uintLevel := uint8(level)

	return domain.NewBatteryStatus(uintLevel, batteryState), nil
}

func (u *UbuntuReader) getBatteryStateFromSystemStr(state string) domain.BatteryState {
	cleanStr := strings.TrimSpace(strings.ToLower(state))
	if cleanStr == "charging" {
		return domain.Charging
	} else if cleanStr == "discharging" {
		return domain.Discharging
	} else {
		return domain.Unknown
	}
}
