package server

import (
	_ "io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultHandler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(defaultHandler))
	defer ts.Close()
	
	res, err := http.Get(ts.URL)
	require.NoError(t, err)
	assert.Equal(t, 404, res.StatusCode)
	require.NoError(t, err)
	assert.Equal(t, JsonContentType, res.Header.Get(ContentTypeHeader))
	err = ReadServerError(res)
	require.Error(t, err)
	require.IsType(t, ServerError{}, err)
	servErr, _ := err.(ServerError)
	assert.Equal(t, 404, servErr.StatusCode)
}

func TestCoreResourcesHandler_BadFormat(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(coreResourcesHandler))
	defer ts.Close()
	
	res, err := http.Get(ts.URL + serveCoreResourcesRootPath + "/foo_bar")
	require.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)
	require.NoError(t, err)
	assert.Equal(t, JsonContentType, res.Header.Get(ContentTypeHeader))
	err = ReadServerError(res)
	require.Error(t, err)
	require.IsType(t, ServerError{}, err)
	servErr, _ := err.(ServerError)
	assert.Equal(t, 400, servErr.StatusCode)

	res, err = http.Get(ts.URL + serveCoreResourcesRootPath + "/foo/bar")
	require.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)
	require.NoError(t, err)
	assert.Equal(t, JsonContentType, res.Header.Get(ContentTypeHeader))
	err = ReadServerError(res)
	require.Error(t, err)
	require.IsType(t, ServerError{}, err)
	servErr, _ = err.(ServerError)
	assert.Equal(t, 400, servErr.StatusCode)

	res, err = http.Get(ts.URL + serveCoreResourcesRootPath + "/foo/-bars")
	require.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)
	require.NoError(t, err)
	assert.Equal(t, JsonContentType, res.Header.Get(ContentTypeHeader))
	err = ReadServerError(res)
	require.Error(t, err)
	require.IsType(t, ServerError{}, err)
	servErr, _ = err.(ServerError)
	assert.Equal(t, 400, servErr.StatusCode)

	res, err = http.Get(ts.URL + serveCoreResourcesRootPath + "/foo/bars/-baz")
	require.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)
	require.NoError(t, err)
	assert.Equal(t, JsonContentType, res.Header.Get(ContentTypeHeader))
	err = ReadServerError(res)
	require.Error(t, err)
	require.IsType(t, ServerError{}, err)
	servErr, _ = err.(ServerError)
	assert.Equal(t, 400, servErr.StatusCode)
}