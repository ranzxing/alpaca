package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRoundTripper(t *testing.T) {
	handler := func(w http.ResponseWriter, req *http.Request) { fmt.Fprintln(w, "It works!") }
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	rt, err := newRoundTripper("tcp", server.Listener.Addr().String())
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
	written, err := io.Copy(ioutil.Discard, resp.Body)
	require.Nil(t, err)
	assert.Equal(t, resp.ContentLength, written)

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
