package domain

type BatteryStatus struct {
	Level uint8
	State BatteryState
}

func NewBatteryStatus(level uint8, state BatteryState) *BatteryStatus {
	return &BatteryStatus{Level: level, State: state}
}

type BatteryState int

const (
	Unknown     = 0
	Charging    = 1
	Discharging = 2
)

func (b BatteryState) String() string {
	switch b {
	case Charging:
		return "Charging"
	case Discharging:
		return "Discharging"
	case Unknown:
	default:
		return "Unknown"
	}
	return "none"
}
