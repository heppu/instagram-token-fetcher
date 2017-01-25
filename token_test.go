package igtoken

import "testing"

// Set these to test client
const (
	clientId    = ""
	redirectUrl = ""
	igUser      = ""
	igPassword  = ""
)

func TestScope(t *testing.T) {
	expected := "basic+comments+likes"
	scopes := []Scope{BASIC, COMMENTS, LIKES}
	output := joinScopes(scopes)
	if expected != output {
		t.Errorf("'%s' does not equal '%s'\n", expected, output)
	}
}

func TestGetToken(t *testing.T) {
	if clientId == "" || redirectUrl == "" || igUser == "" || igPassword == "" {
		t.Skip("Can not test client until required parameters are set. Check token_test.go file.")
	}
	tokenClient := NewClient(clientId, redirectUrl, igUser, igPassword)
	if _, err := tokenClient.GetToken(PUBLIC_CONTENT); err != nil {
		t.Error(err)
	}
}
