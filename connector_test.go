package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestConnector(t *testing.T) {
	handler := func(w http.ResponseWriter, req *http.Request) { fmt.Fprintln(w, "It works!") }
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	rt, err := newConnector("tcp", server.Listener.Addr().String())
	require.Nil(t, err)
	defer rt.Close()

	req, err := http.NewRequest(http.MethodConnect, server.URL, nil)
	require.Nil(t, err)
	u, err := url.Parse(server.URL)
	require.Nil(t, err)
	req.Host = u.Host
	resp, err := rt.RoundTrip(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequest(http.MethodGet, server.URL, nil)
	require.Nil(t, err)
	resp, err = rt.RoundTrip(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	require.Nil(t, err)
	assert.Equal(t, "It works!\n", string(body))
}

func TestConnectorWithInterleavedReads(t *testing.T) {
	handler := func(w http.ResponseWriter, req *http.Request) { fmt.Fprintln(w, "It works!") }
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	rt, err := newConnector("tcp", server.Listener.Addr().String())
	require.Nil(t, err)
	defer rt.Close()
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.Nil(t, err)

	resp1, err := rt.RoundTrip(req)
	require.Nil(t, err)
	defer resp1.Body.Close()

	resp2, err := rt.RoundTrip(req)
	require.Nil(t, err)
	defer resp2.Body.Close()

	assert.Equal(t, http.StatusOK, resp1.StatusCode)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)
}
