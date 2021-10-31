package gauth

import "github.com/joeshaw/envdecode"

type Config struct {
	CredAddress string `env:"GAUTH_URL_ADDRESS,default=0.0.0.0:3000"`
	WebURL      string `env:"GAUTH_WEB_URL,default=http://localhost:3000"`

	EmailVerifierToken string `env:"GAUTH_EMAIL_VERIFIER_TOKEN"`
	PwnedPasswordsURL  string `env:"GAUTH_PWNED_PASSWORDS_URL"`
}

func NewConfig() (Config, error) {
	cfg := Config{}
	if err := envdecode.Decode(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}