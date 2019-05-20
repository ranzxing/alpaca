package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
)

type connector struct {
	conn   net.Conn
	reader *bufio.Reader
}

func newConnector(network, address string) (*connector, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	rd := bufio.NewReader(conn)
	return &connector{conn, rd}, nil
}

func (rt *connector) RoundTrip(req *http.Request) (*http.Response, error) {
	if rt.conn == nil {
		return nil, errors.New("connection closed, can't send request")
	}
	if err := req.Write(rt.conn); err != nil {
		return nil, err
	}
	resp, err := http.ReadResponse(rt.reader, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var body bytes.Buffer
	// TODO: Not sure if we should trust content length here. Do we even need the response body,
	// or can we just discard it?
	if _, err = io.CopyN(&body, resp.Body, resp.ContentLength); err != nil {
		return nil, err
	}
	resp.Body = ioutil.NopCloser(&body)
	return resp, err
}

func (rt *connector) hijack() net.Conn {
	rt.conn = nil
	return rt.conn
}

func (rt *connector) Close() {
	if rt.conn != nil {
		rt.conn.Close()
		rt.conn = nil
	}
}
