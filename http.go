package tUtils

import (
	"crypto/tls"
	httpHelper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//--------------------------------------------------- Functions --------------------------------------------------------
//----------------------------------------------------------------------------------------------------------------------

// AssertHttpRes checks if statusCode == wantCode && that gotBody contains wantBody
func AssertHttpRes(t Tester, res *http.Response, wantBody string, wantCode int) {
	b, err := ioutil.ReadAll(res.Body)
	gotBody := string(b)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), wantCode, res.StatusCode)
	assert.Contains(t.T(), gotBody, wantBody)
}

func HttpGet(h HttpConfig, uri string) *http.Response {
	req, err := http.NewRequest("GET", uri, nil)
	assert.NoError(h.T(), err)

	res, err := h.HttpClient().Do(req)
	assert.NoError(h.T(), err)
	return res
}

func HttpGetAndAssert(h HttpConfig, uri, wantBody string, wantCode int) {
	res := HttpGet(h, uri)
	AssertHttpRes(h, res, wantBody, wantCode)
}

// HttpGetWithRetry Runs a http GET requests with retries & Checks if returned statusCode is equal to the expected
// statusCode
//	url 		string
//  wantCode 	int				// want: http status code
//	retries		int				// number of retries
//	sleep		time.Duration 	// time between retries
func HttpGetWithRetry(t TlsTester, url string, wantCode, retries int, sleep time.Duration) {
	httpHelper.HttpGetWithRetryWithCustomValidation(
		t.T(),
		url,
		t.TlsConfig(),
		retries,
		sleep,
		func(gotCode int, _ string) bool {
			return gotCode == wantCode
		},
	)
}

//----------------------------------------------------------------------------------------------------------------------
//------------------------------------------------ HttpConfigBuilder ---------------------------------------------------
//----------------------------------------------------------------------------------------------------------------------

type httpConfigBuilder struct {
	HttpConfigBuilder
	httpConfig httpConfig
}

func NewHttpConfigBuilder() HttpConfigBuilder {
	return &httpConfigBuilder{httpConfig: httpConfig{}}
}

func (b *httpConfigBuilder) Build() HttpConfig {
	return &b.httpConfig
}

func (b *httpConfigBuilder) SetCookies(uri string, cookies []*http.Cookie) HttpConfigBuilder {
	parsedUri, err := url.Parse(uri)
	require.NoError(b.httpConfig.T(), err)
	b.httpConfig.HttpClient().Jar.SetCookies(parsedUri, cookies)
	return b
}

func (b *httpConfigBuilder) SetHttpClient(hc *http.Client) HttpConfigBuilder {
	b.httpConfig.httpClient = hc
	return b
}

func (b *httpConfigBuilder) SetProxy(uri string) HttpConfigBuilder {
	parsedUrl, err := url.Parse(uri)
	require.NoError(b.httpConfig.T(), err)
	b.httpConfig.httpClient.Transport = &http.Transport{Proxy: http.ProxyURL(parsedUrl)}
	return b
}

func (b *httpConfigBuilder) SetT(t *testing.T) HttpConfigBuilder {
	b.httpConfig.t = t
	return b
}

func (b *httpConfigBuilder) SetTlsConfig(tc *tls.Config) HttpConfigBuilder {
	b.httpConfig.tlsConfig = tc
	return b
}

//----------------------------------------------------------------------------------------------------------------------
//----------------------------------------------------- HttpConfig -----------------------------------------------------
//----------------------------------------------------------------------------------------------------------------------

type httpConfig struct {
	HttpConfig
	httpClient *http.Client
	t          *testing.T
	tlsConfig  *tls.Config
}

func NewHttpConfig(t *testing.T) HttpConfig {
	b := NewHttpConfigBuilder()
	hc := &http.Client{}

	b.
		SetT(t).
		SetHttpClient(hc).
		SetTlsConfig(&tls.Config{})

	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	hc.Jar = jar

	return b.Build()
}

func (h *httpConfig) T() *testing.T {
	return h.t
}

func (h *httpConfig) HttpClient() *http.Client {
	return h.httpClient
}

func (h *httpConfig) TlsConfig() *tls.Config {
	return h.tlsConfig
}
