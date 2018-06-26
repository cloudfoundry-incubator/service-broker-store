package credhub_shims

import (
	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
)

//go:generate counterfeiter -o ./credhub_fakes/credhub_auth_fake.go . CredhubAuth
type CredhubAuth interface {
	UaaClientCredentials(clientId, clientSecret string) auth.Builder
}

type CredhubAuthShim struct {
}

func (c *CredhubAuthShim) UaaClientCredentials(clientId, clientSecret string) auth.Builder {
	return auth.UaaClientCredentials(clientId, clientSecret)
}

//go:generate counterfeiter -o ./credhub_fakes/credhub_fake.go . Credhub
type Credhub interface {
	SetJSON(name string, value values.JSON, overwrite credhub.Mode) (credentials.JSON, error)
	GetLatestJSON(name string) (credentials.JSON, error)
	Delete(name string) error
}

type CredhubShim struct {
	delegate *credhub.CredHub
}

func NewCredhubShim(
	url string,
	caCert string,
	clientID string,
	clientSecret string,
	uaaCACert string,
	authShim CredhubAuth,
) (Credhub, error) {
	var delegate *credhub.CredHub
	var err error
	auth := credhub.Auth(authShim.UaaClientCredentials(clientID, clientSecret))

	caCerts := []string{}
	if caCert != "" {
		caCerts = append(caCerts, caCert)
	}

	if uaaCACert != "" {
		caCerts = append(caCerts, uaaCACert)
	}

	if len(caCerts) > 0 {
		delegate, err = credhub.New(url, credhub.CaCerts(caCerts...), auth)
	} else {
		delegate, err = credhub.New(url, auth)
	}

	if err != nil {
		return nil, err
	}

	return &CredhubShim{
		delegate: delegate,
	}, nil
}

func (ch *CredhubShim) SetJSON(name string, value values.JSON, overwrite credhub.Mode) (credentials.JSON, error) {
	return ch.delegate.SetJSON(name, value, overwrite)
}

func (ch *CredhubShim) GetLatestJSON(name string) (credentials.JSON, error) {
	return ch.delegate.GetLatestJSON(name)
}

func (ch *CredhubShim) Delete(name string) error {
	return ch.delegate.Delete(name)
}
