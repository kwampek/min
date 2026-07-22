package models

import (
	"fmt"
	"time"

	"github.com/mileusna/useragent"
)

type Identifier struct {
	SessionID int
	UserID    int
}

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
	SessionID  int       `db:"session_id"`
	UserID     int       `db:"user_id"`
	TokenHash  string    `db:"token_hash"`
	DeviceName string    `db:"device_name"`
	IPAddress  string    `db:"ip_address"`
	ExpiresAt  time.Time `db:"expires_at"`
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
