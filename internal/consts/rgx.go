package consts

import "regexp"

// nolint: gochecknoglobals
var (
	// RgxEmail is a regexp for validating email addresses
	RgxEmail = regexp.MustCompile(
		"^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@" +
			"[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)
