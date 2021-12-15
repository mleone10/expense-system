package api

type Config struct {
	CognitoClientId     string
	CognitoClientSecret string
	ClientHostname      string
	ClientScheme        string
}

func (c Config) getClientHostname() string {
	return c.ClientHostname
}

func (c Config) getClientScheme() string {
	return c.ClientScheme
}

func (c Config) getCognitoClientId() string {
	return c.CognitoClientId
}

func (c Config) getCognitoClientSecret() string {
	return c.CognitoClientSecret
}
