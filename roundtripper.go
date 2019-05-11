package main

import (
	"bufio"
	"net"
	"net/http"
)

type roundTripper struct {
	conn   net.Conn
	reader *bufio.Reader
}

func newRoundTripper(network, address string) (*roundTripper, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	rd := bufio.NewReader(conn)
	return &roundTripper{conn, rd}, nil
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := req.Write(rt.conn); err != nil {
		return nil, err
	}
	return http.ReadResponse(rt.reader, req)
}

func (rt *roundTripper) Close() {
	rt.conn.Close()
}
