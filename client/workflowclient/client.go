package workflowclient

import (
	"crypto/tls"
	"fmt"
	"go.temporal.io/sdk/client"
	"log"
	"net"
)

const (
	temporalHostPort     = "demo-cloud-ops.ps13i.tmprl.cloud:7233"
	temporalNamespace    = "demo-cloud-ops.ps13i"
	controlPlaneHostPort = "saas-api.tmprl-test.cloud:443"
	//apiKeyValue          = "tmprl_myCSpZXG5EyGYOhCCriejuCf814zFRNz_Np9b0JC6knytsMz3261G1VT85x4vIJ5fBPdhuKK4SuTbK0XruUBgobydixE7icGB"
	// demo3 (in the ps13i account)
	//apiKeyValue = "tmprl_46L87rLmDTmve2mqycurdcJE9DEDNF7j_VBZGpFXcL3WD4UWVtQ6wxz9nRmRuKNMs6BSFBtlkNF1BgSxvZAGL3Wi5lXQGWwj2"
	// demo4 (in the temporal-dev account)
	apiKeyValue = "tmprl_aWsGvnW3p4kwGbocxW83PWaOF0MbsMLQ_thyBko0S3ukKi9k9ZsTNZOgcYcVHGlNEZSuBD6hyoLvoFr6kRC1ZC0GPnjXFb1aL"

	certFilePath = "/Users/liang/demo/demo-cloud-ops.ps13i.pem"
	keyFilePath  = "/Users/liang/demo/demo-cloud-ops.ps13i.key"
)

func MustGetClient() client.Client {
	tlsConfig, err := getTLSConfig()
	if err != nil {
		log.Fatalln("failed to create TLS config", err)
	}

	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort:          temporalHostPort,
		Namespace:         temporalNamespace,
		ConnectionOptions: client.ConnectionOptions{TLS: tlsConfig},
		Logger:            NewDefaultLogger(),
	})

	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	return c
}

func getTLSConfig() (*tls.Config, error) {
	serverName, _, parseErr := net.SplitHostPort(temporalHostPort)
	if parseErr != nil {
		return nil, fmt.Errorf("failed to split hostport %s: %w", temporalHostPort, parseErr)
	}
	var cert tls.Certificate
	var err error
	cert, err = tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS from files: %w", err)
	}
	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ServerName:         serverName,
		InsecureSkipVerify: false,
	}, nil
}
