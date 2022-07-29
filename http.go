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

func HttpGet(h HttpTester, uri string) *http.Response {
	req, err := http.NewRequest("GET", uri, nil)
	assert.NoError(h.T(), err)

	res, err := h.HttpClient().Do(req)
	assert.NoError(h.T(), err)
	return res
}

func HttpGetAndAssert(h HttpTester, uri, wantBody string, wantCode int) {
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
//------------------------------------------------------ Struct --------------------------------------------------------
//----------------------------------------------------------------------------------------------------------------------

type HttpConfig struct {
	// Getter
	HttpClientGetter
	Tester
	TlsConfigGetter

	// Setters
	CookieSetter
	HttpClientSetter
	ProxySetter
	TestSetter
	TlsConfigSetter

	// Fields
	httpClient *http.Client
	t          *testing.T
	tlsConfig  *tls.Config
}

func NewHttpConfig(t *testing.T) HttpConfig {
	h := HttpConfig{}
	hc := &http.Client{}

	h.SetT(t)
	h.SetHttpClient(hc)
	h.SetTlsConfig(&tls.Config{})

	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	hc.Jar = jar

	return h
}

//----------------------------------------------------------------------------------------------------------------------
//---------------------------------------------------- Interfaces ------------------------------------------------------
//----------------------------------------------------------------------------------------------------------------------

type TlsTester interface {
	Tester
	TlsConfigGetter
}

type HttpTester interface {
	HttpClientGetter
	Tester
	TlsConfigGetter
}

type HttpClientGetter interface{ HttpClient() *http.Client }
type TlsConfigGetter interface{ TlsConfig() *tls.Config }

type CookieSetter interface{ SetCookies(string, []*http.Cookie) }
type HttpClientSetter interface{ SetHttpClient(*http.Client) }
type ProxySetter interface{ SetProxy(string) }
type TestSetter interface{ SetT(*testing.T) }
type TlsConfigSetter interface{ SetTlsConfig(*tls.Config) }

//----------------------------------------------------------------------------------------------------------------------
//------------------------------------------------------ Getters -------------------------------------------------------
//----------------------------------------------------------------------------------------------------------------------

func (h *HttpConfig) T() *testing.T {
	return h.t
}

func (h *HttpConfig) HttpClient() *http.Client {
	return h.httpClient
}

func (h *HttpConfig) TlsConfig() *tls.Config {
	return h.tlsConfig
}

//----------------------------------------------------------------------------------------------------------------------
//------------------------------------------------------ Setters -------------------------------------------------------
//----------------------------------------------------------------------------------------------------------------------

func (h *HttpConfig) SetCookies(uri string, cookies []*http.Cookie) {
	parsedUri, err := url.Parse(uri)
	require.NoError(h.T(), err)

	h.HttpClient().Jar.SetCookies(parsedUri, cookies)
}

func (h *HttpConfig) SetHttpClient(hc *http.Client) {
	h.httpClient = hc
}

func (h *HttpConfig) SetProxy(uri string) {
	parsedUrl, err := url.Parse(uri)
	require.NoError(h.T(), err)

	h.httpClient.Transport = &http.Transport{Proxy: http.ProxyURL(parsedUrl)}
}

func (h *HttpConfig) SetT(t *testing.T) {
	h.t = t
}

func (h *HttpConfig) SetTlsConfig(tc *tls.Config) {
	h.tlsConfig = tc
}
