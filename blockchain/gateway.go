package blockchain

import (
	"crypto/x509"
	"fmt"
	"gin-fabric-connector/common"
	"os"
	"path"
	"sync"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	once sync.Once
	conn *grpc.ClientConn
)

type GatewayManager struct{}

func (g *GatewayManager) Init() (err error) {
	once.Do(func() {
		err = createGRPCConnection()
		if err != nil {
			fmt.Println("Failed to initialize gRPC connection:", err.Error())
			return
		}
	})
	return
}

func (g *GatewayManager) GetGateway() (gateway *client.Gateway, err error) {
	id, err := getIdentity()
	if err != nil {
		return
	}

	sign, err := getSign()
	if err != nil {
		return
	}

	gateway, err = client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(conn),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		fmt.Println("Failed to create Fabric gateway:", err.Error())
	}
	return
}

func createGRPCConnection() (err error) {
	config := common.GetConfig()

	certificate, err := getCertificate(config.PeerTLSCert)
	if err != nil {
		fmt.Println("Failed to load peer TLS certificate:", err.Error())
		return
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, config.PeerId)

	conn, err = grpc.Dial(config.PeerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		fmt.Println("Failed to create gRPC connection:", err.Error())
	}
	return
}

func getIdentity() (*identity.X509Identity, error) {
	config := common.GetConfig()

	cert, err := getCertificate(config.ClientCert)
	if err != nil {
		fmt.Println("Failed to load client certificate:", err.Error())
		return nil, err
	}

	id, err := identity.NewX509Identity(config.MSPID, cert)
	if err != nil {
		fmt.Println("Failed to create client identity:", err.Error())
		return nil, err
	}
	return id, nil
}

func getSign() (sign identity.Sign, err error) {
	config := common.GetConfig()

	files, err := os.ReadDir(config.ClientKey)
	if err != nil {
		fmt.Println("Failed to read key folder:", err.Error())
		return nil, err
	}

	privateKeyPEM, err := os.ReadFile(path.Join(config.ClientKey, files[0].Name()))
	if err != nil {
		fmt.Println("Failed to read key file:", err.Error())
		return nil, err
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		fmt.Println("Failed to create private key from PEM:", err.Error())
		return nil, err
	}

	sign, err = identity.NewPrivateKeySign(privateKey)
	if err != nil {
		fmt.Println("Failed to create sign:", err.Error())
		return nil, err
	}
	return
}

func getCertificate(filename string) (cert *x509.Certificate, err error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Failed to read certificate file:", err.Error())
		return nil, err
	}
	cert, err = identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		fmt.Println("Failed to create certificate from PEM:", err.Error())
	}
	return
}
