Level 1: Setup Instructions
Prerequisites

Install Docker and Docker Compose
Install Go (version 1.19 or later)
Install Git
Install curl and jq

Setup Steps

Clone Fabric Samples Repository

bashmkdir hyperledger-fabric-assignment
cd hyperledger-fabric-assignment
mkdir level-1-setup
cd level-1-setup

curl -sSL https://bit.ly/2ysbOFE | bash -s
cd fabric-samples/test-network

Start the Test Network

bash./network.sh down
./network.sh up createChannel -ca

Verify Network is Running

bashdocker ps
# You should see containers for orderer, peer0.org1, peer0.org2, and CAs