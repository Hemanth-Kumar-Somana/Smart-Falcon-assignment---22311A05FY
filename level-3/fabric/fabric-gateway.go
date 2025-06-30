package fabric

import (
    "crypto/x509"
    "fmt"
    "io/ioutil"
    "path"
    "time"

    "github.com/hyperledger/fabric-gateway/pkg/client"
    "github.com/hyperledger/fabric-gateway/pkg/identity"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
)

type FabricGateway struct {
    Gateway  *client.Gateway
    Contract *client.Contract
}

func NewFabricGateway() (*FabricGateway, error) {
    // Path to crypto materials
    cryptoPath := "../level-1-setup/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com"
    
    // Path to user private key directory
    keyPath := path.Join(cryptoPath, "users", "User1@org1.example.com", "msp", "keystore")
    
    // Path to user certificate
    certPath := path.Join(cryptoPath, "users", "User1@org1.example.com", "msp", "signcerts", "cert.pem")
    
    // Path to peer tls certificate
    tlsCertPath := path.Join(cryptoPath, "peers", "peer0.org1.example.com", "tls", "ca.crt")

    // Read the user certificate
    cert, err := ioutil.ReadFile(certPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read certificate file: %w", err)
    }

    // Read the user private key
    keyFiles, err := ioutil.ReadDir(keyPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read private key directory: %w", err)
    }
    if len(keyFiles) == 0 {
        return nil, fmt.Errorf("no private key files found")
    }
    
    privateKeyPath := path.Join(keyPath, keyFiles[0].Name())
    privateKey, err := ioutil.ReadFile(privateKeyPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read private key file: %w", err)
    }

    // Create identity
    id, err := identity.NewX509Identity("Org1MSP", cert, privateKey)
    if err != nil {
        return nil, fmt.Errorf("failed to create identity: %w", err)
    }

    // Read TLS certificate
    tlsCert, err := ioutil.ReadFile(tlsCertPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read TLS certificate: %w", err)
    }

    // Create certificate pool and add the TLS certificate
    certPool := x509.NewCertPool()
    certPool.AppendCertsFromPEM(tlsCert)

    // Create gRPC connection
    conn, err := grpc.Dial("localhost:7051", grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(certPool, "peer0.org1.example.com")))
    if err != nil {
        return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
    }

    // Create gateway
    gateway, err := client.Connect(
        id,
        client.WithSign(identity.NewPrivateKeySign(privateKey)),
        client.WithHash(identity.SHA256),
        client.WithClientConnection(conn),
        client.WithEvaluateTimeout(5*time.Second),
        client.WithEndorseTimeout(15*time.Second),
        client.WithSubmitTimeout(5*time.Second),
        client.WithCommitStatusTimeout(1*time.Minute),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to connect to gateway: %w", err)
    }

    // Get network
    network := gateway.GetNetwork("mychannel")
    
    // Get contract
    contract := network.GetContract("asset")

    return &FabricGateway{
        Gateway:  gateway,
        Contract: contract,
    }, nil
}

func (fg *FabricGateway) Close() {
    fg.Gateway.Close()
}
