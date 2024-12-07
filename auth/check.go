package auth

import "github.com/valyentdev/cli/config"

func IsLoggedIn() bool {
	_, err := config.RetrieveAPIKey()
	if err != nil {
		return false
	}

	// TODO: Check key against Valyent's console API.

	return true
}
