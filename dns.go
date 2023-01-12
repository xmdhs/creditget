//go:build androidgodns

package main

import (
	"context"
	"net"
)

const bootstrapDNS = "223.5.5.5:53"

func init() {
	var dialer net.Dialer
	net.DefaultResolver = &net.Resolver{
		PreferGo: false,
		Dial: func(context context.Context, _, _ string) (net.Conn, error) {
			conn, err := dialer.DialContext(context, "udp", bootstrapDNS)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
}
