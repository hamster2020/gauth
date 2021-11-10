package gauthclient

func (gc GauthClient) Authenticate(email, password string) error {
	gc.authClient.SetEmail(email)
	gc.authClient.SetPassword(password)
	return gc.authClient.Authenticate()
}

func (gc GauthClient) Logout() error {
	return gc.authClient.Logout()
}
