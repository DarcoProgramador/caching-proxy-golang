package utils

import (
	"log/slog"
	"net"
	"net/url"
	"strings"
)

func IsUrl(str string) bool {
	url, err := url.ParseRequestURI(str)
	if err != nil {
		slog.Error("url-parse-error", slog.String("Error", err.Error()))
		return false
	}

	address := net.ParseIP(url.Host)
	slog.Info("url-info", "host", address)

	if address == nil {
		slog.Info("url-info", "host", url.Host)

		return strings.Contains(url.Host, ".")
	}

	return true
}
