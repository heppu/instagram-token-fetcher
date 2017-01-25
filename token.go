package igtoken

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	authUrl  string = "https://api.instagram.com/oauth/authorize/"
	loginUrl string = "https://www.instagram.com"
	scopeSep string = "+"

	BASIC          Scope = "basic"
	PUBLIC_CONTENT Scope = "public_content"
	FOLLOWER_LIST  Scope = "follower_list"
	COMMENTS       Scope = "comments"
	RELATIONSHIPS  Scope = "relationships"
	LIKES          Scope = "likes"
)

type Scope string

type TokenClient struct {
	clientId     string
	redirectUrl  string
	userNick     string
	userPassword string
	client       *http.Client
}

func NewClient(clientId, redirectUrl, userNick, userPassword string) (tc *TokenClient) {
	jar, _ := cookiejar.New(nil)
	return &TokenClient{
		clientId:     clientId,
		redirectUrl:  redirectUrl,
		userNick:     userNick,
		userPassword: userPassword,
		client: &http.Client{
			Jar: jar,
		},
	}
}

func (t *TokenClient) GetToken(scopes ...Scope) (accessToken string, err error) {
	// Create server to receive token
	var u *url.URL
	var port string
	var closeServer func() bool
	u, err = url.Parse(t.redirectUrl)
	if _, port, err = net.SplitHostPort(u.Host); err != nil {
		return
	}
	if closeServer, err = startServer(port, u.Path); err != nil {
		return
	}
	defer closeServer()

	// Get login page
	var r *http.Response
	var req *http.Request
	if req, err = http.NewRequest("GET", authUrl, new(bytes.Buffer)); err != nil {
		return
	}

	values := req.URL.Query()
	values.Add("client_id", t.clientId)
	values.Add("redirect_uri", t.redirectUrl)
	values.Add("response_type", "token")

	if len(scopes) > 0 {
		values.Add("scope", joinScopes(scopes))
	}

	req.URL.RawQuery = values.Encode()
	r, err = t.client.Do(req)

	// Do login
	var csrfToken, action string
	if csrfToken, action, err = parseForm(r.Body); err != nil {
		return
	}
	if r, err = t.login(csrfToken, action); err != nil {
		return
	}

	// If we have already authorized this app we will get token straight away
	f := r.Request.URL.Fragment
	if len(f) > 13 && f[0:13] == "access_token=" {
		accessToken = f[13:]
		return
	}

	// We did not get token from redirect url.
	// This could mean couple of things:
	//     - We got error message in json format
	//     - App was not authorized so we got auth form
	//     - We could some other error message

	var b []byte
	if b, err = ioutil.ReadAll(r.Body); err != nil {
		return
	}

	// We got empty body but no access token so something is wrong
	if len(b) == 0 {
		err = fmt.Errorf("Empty body with status code %d", r.StatusCode)
	}

	// Looks like we got error message in json format
	if b[0] == '{' {
		err = fmt.Errorf("%s", b)
		return
	}

	// Authorize app for account
	if csrfToken, action, err = parseForm(ioutil.NopCloser(bytes.NewReader(b))); err != nil {
		return
	}
	if r, err = t.auth(csrfToken, action); err != nil {
		return
	}

	f = r.Request.URL.Fragment
	if len(f) > 13 && f[0:13] == "access_token=" {
		accessToken = f[13:]
		return
	}

	alert := parseAlert(ioutil.NopCloser(bytes.NewReader(b)))
	if alert != "" {
		err = fmt.Errorf("%s", alert)
		return
	}

	err = fmt.Errorf("%s", b)
	return
}

func (t *TokenClient) login(csrftToken, action string) (r *http.Response, err error) {
	formData := url.Values{}
	formData.Add("csrfmiddlewaretoken", csrftToken)
	formData.Add("username", t.userNick)
	formData.Add("password", t.userPassword)

	var req *http.Request
	if req, err = http.NewRequest("POST", loginUrl+action, bytes.NewBufferString(formData.Encode())); err != nil {
		return
	}

	req.Header.Add("Referer", loginUrl+action)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if r, err = t.client.Do(req); err != nil {
		return
	}

	return
}

func (t *TokenClient) auth(csrftToken, action string) (r *http.Response, err error) {
	formData := url.Values{}
	formData.Add("csrfmiddlewaretoken", csrftToken)
	formData.Add("allow", "Authorize")

	var req *http.Request
	if req, err = http.NewRequest("POST", loginUrl+action, bytes.NewBufferString(formData.Encode())); err != nil {
		return
	}

	req.Header.Add("Referer", loginUrl+action)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	r, err = t.client.Do(req)
	return
}
