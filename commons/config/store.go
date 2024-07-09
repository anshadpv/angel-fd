package config

import "github.com/angel-one/fd-core/constants"

const (
	DS = "xxxxx"
)

type SecretStore struct {
}

func initStore() map[string]string {
	store := map[string]string{
		constants.DATABASE_USERNAME:   client.getStringSecretD(constants.DBUsername, DS),
		constants.DATABASE_PASSWORD:   client.getStringSecretD(constants.DBPassword, DS),
		constants.DATABASE_URL:        client.getStringSecretD(constants.DBHost, DS),
		constants.DATABASE_NAME:       client.getStringSecretD(constants.DBName, DS),
		constants.JWTSymmetricKey:     client.getStringSecretD(constants.JWTSymmetricKey, DS),
		constants.ProfileServiceToken: client.getStringSecretD(constants.ProfileServiceToken, DS),
		constants.UpswingClientId:     client.getStringSecretD(constants.UpswingClientId, DS),
		constants.UpswingGrantType:    client.getStringSecretD(constants.UpswingGrantType, DS),
		constants.UpswingClientSecret: client.getStringSecretD(constants.UpswingClientSecret, DS),
		constants.UpswingScope:        client.getStringSecretD(constants.UpswingScope, DS),
	}
	return store
}
