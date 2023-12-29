package common

import (
	"os"
	"path"
	"sync"

	"github.com/spf13/viper"
)

var BASEDIR = getEnv("FABRIC_CONNECTOR_BASEDIR", os.ExpandEnv("$HOME/Projects/fabric-samples/test-network"))

var (
	once   sync.Once
	config *Config
)

type Config struct {
	Auth0Audience string
	Auth0Domain   string
	BaseDir       string
	MSPID         string
	ClientKey     string
	ClientCert    string
	PeerId        string
	PeerTLSCert   string
	PeerEndpoint  string
}

// GetConfig creates a Config singleton instance from environment variables using Viper
func GetConfig() Config {
	once.Do(func() {
		viper.AutomaticEnv()
		viper.SetEnvPrefix("FABRIC_CONNECTOR")

		config = &Config{
			Auth0Audience: getEnv("AUTH0_AUDIENCE", ""),
			Auth0Domain:   getEnv("AUTH0_DOMAIN", ""),
			MSPID:         getEnv("MSPID", "Org1MSP"),
			ClientKey:     getEnv("CLIENT_KEY", path.Join(BASEDIR, "/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/")),
			ClientCert:    getEnv("CLIENT_CERT", path.Join(BASEDIR, "/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem")),
			PeerId:        getEnv("PEER_ID", "peer0.org1.example.com"),
			PeerTLSCert:   getEnv("PEER_TLS_CERT", path.Join(BASEDIR, "/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt")),
			PeerEndpoint:  getEnv("PEER_ENDPOINT", "localhost:7051"),
		}
	})
	return *config
}

// getEnv returns the environment variable that matches key, or a default value
func getEnv(key, defaultValue string) string {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	}
	return value
}
