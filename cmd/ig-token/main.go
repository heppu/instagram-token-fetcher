package main

import (
	"fmt"
	"os"

	"github.com/heppu/instagram-token-fetcher"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	clientId    = kingpin.Flag("client-id", "application client ID").Short('i').Required().String()
	redirectUrl = kingpin.Flag("url", "redirect url").Short('r').Required().String()
	igUser      = kingpin.Flag("user", "instagram username").Short('u').Required().String()
	igPassword  = kingpin.Flag("password", "password for user").Short('p').Required().String()
	scope       = kingpin.Flag("scope", "scope for access token").Short('s').Strings()
)

func main() {
	kingpin.Parse()
	tokenClient := igtoken.NewClient(*clientId, *redirectUrl, *igUser, *igPassword)

	scopes := make([]igtoken.Scope, 0)
	for _, s := range *scope {
		scopes = append(scopes, igtoken.Scope(s))
	}

	token, err := tokenClient.GetToken(scopes...)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	fmt.Print(token)
}
