// Package ntpclient requests NTP server to return the exact current time.
package ntpclient

import (
	"fmt"
	"time"

	"github.com/beevik/ntp"
)

func GetCurrentTime() (time.Time, error) {
	ntpTime, err := ntp.Time("pool.ntp.org")
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to communicate with NTP server: %w", err)
	}

	return ntpTime, nil
}
