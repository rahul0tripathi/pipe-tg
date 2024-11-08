package tg

import (
	"errors"
	"time"
)

const (
	_maxCodePollingTimeout = time.Minute * 2
)

var (
	ErrAuthExp = errors.New("auth expired, reauthenticate")
)

type TGSessionConfig struct {
	Version int `json:"Version"`
	Data    struct {
		Config struct {
			BlockedMode  bool `json:"BlockedMode"`
			ForceTryIpv6 bool `json:"ForceTryIpv6"`
			Date         int  `json:"Date"`
			Expires      int  `json:"Expires"`
			TestMode     bool `json:"TestMode"`
			ThisDC       int  `json:"ThisDC"`
			DCOptions    []struct {
				Flags             int         `json:"Flags"`
				Ipv6              bool        `json:"Ipv6"`
				MediaOnly         bool        `json:"MediaOnly"`
				TCPObfuscatedOnly bool        `json:"TCPObfuscatedOnly"`
				CDN               bool        `json:"CDN"`
				Static            bool        `json:"Static"`
				ThisPortOnly      bool        `json:"ThisPortOnly"`
				ID                int         `json:"ID"`
				IPAddress         string      `json:"IPAddress"`
				Port              int         `json:"Port"`
				Secret            interface{} `json:"Secret"`
			} `json:"DCOptions"`
			DCTxtDomainName string `json:"DCTxtDomainName"`
			TmpSessions     int    `json:"TmpSessions"`
			WebfileDCID     int    `json:"WebfileDCID"`
		} `json:"Config"`
		DC        int    `json:"DC"`
		Addr      string `json:"Addr"`
		AuthKey   string `json:"AuthKey"`
		AuthKeyID string `json:"AuthKeyID"`
		Salt      int64  `json:"Salt"`
	} `json:"Data"`
}
