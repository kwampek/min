package models

import (
	"fmt"
	"time"

	"github.com/mileusna/useragent"
)

type Device struct {
	DeviceName string
	IPAddress  string
}

func DeviceFromRaw(userAgent, IP string) Device {
	return Device{
		DeviceName: ParseUserAgent(userAgent),
		IPAddress:  IP,
	}
}

type Session struct {
	UserID     int
	TokenHash  string
	DeviceName string
	IPAddress  string
	ExpiresAt  time.Time
}

func ParseUserAgent(userAgent string) string {
	info := useragent.Parse(userAgent)

	switch {
	case info.Name != "" && info.OS != "":
		return fmt.Sprintf("%s on %s", info.Name, info.OS)
	case info.Name != "":
		return info.Name
	case info.OS != "":
		return info.OS
	default:
		return "Unknown device"
	}
}
