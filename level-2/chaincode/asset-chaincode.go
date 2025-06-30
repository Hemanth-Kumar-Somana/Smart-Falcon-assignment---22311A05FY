package main

import (
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
    contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
type Asset struct {
    DealerID     string  `json:"DEALERID"`
    MSISDN       string  `json:"MSISDN"`
    MPIN         string  `json:"MPIN"`
    Balance      float64 `json:"BALANCE"`
    Status       string  `json:"STATUS"`
    TransAmount  float64 `json:"TRANSAMOUNT"`
    TransType    string  `json:"TRANSTYPE"`
    Remarks      string  `json:"REMARKS"`
    CreatedAt    string  `json:"CREATEDAT"`
    UpdatedAt    string  `json:"UPDATEDAT"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    assets := []Asset{
        {DealerID: "D001", MSISDN: "1234567890", MPIN: "1234", Balance: 1000.00, Status: "ACTIVE", TransAmount: 0, TransType: "INITIAL", Remarks: "Initial Balance", CreatedAt: time.Now().Format(time.RFC3339), UpdatedAt: time.Now().Format(time.RFC3339)},
        {DealerID: "D002", MSISDN: "1234567891", MPIN: "5678", Balance: 2000.00, Status: "ACTIVE", TransAmount: 0, TransType: "INITIAL", Remarks: "Initial Balance", CreatedAt: time.Now().Format(time.RFC3339), UpdatedAt: time.Now().Format(time.RFC3339)},
    }

    for _, asset := range assets {
        assetJSON, err := json.Marshal(asset)
        if err != nil {
            return err
        }

        err = ctx.GetStub().PutState(asset.DealerID, assetJSON)
        if err != nil {
            return fmt.Errorf("failed to put to world state. %v", err)
        }
    }

    return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, dealerID, msisdn, mpin string, balance float64, status, transType, remarks string) error {
    exists, err := s.AssetExists(ctx, dealerID)
    if err != nil {
        return err
    }
    if exists {
        return fmt.Errorf("the asset %s already exists", dealerID)
    }

    asset := Asset{
        DealerID:    dealerID,
        MSISDN:      msisdn,
        MPIN:        mpin,
        Balance:     balance,
        Status:      status,
        TransAmount: balance,
        TransType:   transType,
        Remarks:     remarks,
        CreatedAt:   time.Now().Format(time.RFC3339),
        UpdatedAt:   time.Now().Format(time.RFC3339),
    }
    assetJSON, err := json.Marshal(asset)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(dealerID, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, dealerID string) (*Asset, error) {
    assetJSON, err := ctx.GetStub().GetState(dealerID)
    if err != nil {
        return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if assetJSON == nil {
        return nil, fmt.Errorf("the asset %s does not exist", dealerID)
    }

    var asset Asset
    err = json.Unmarshal(assetJSON, &asset)
    if err != nil {
        return nil, err
    }

    return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, dealerID string, balance float64, status, transAmount float64, transType, remarks string) error {
    exists, err := s.AssetExists(ctx, dealerID)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("the asset %s does not exist", dealerID)
    }

    // Get existing asset
    asset, err := s.ReadAsset(ctx, dealerID)
    if err != nil {
        return err
    }

    // Update fields
    asset.Balance = balance
    asset.Status = status
    asset.TransAmount = transAmount
    asset.TransType = transType
    asset.Remarks = remarks
    asset.UpdatedAt = time.Now().Format(time.RFC3339)

    assetJSON, err := json.Marshal(asset)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(dealerID, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, dealerID string) error {
    exists, err := s.AssetExists(ctx, dealerID)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("the asset %s does not exist", dealerID)
    }

    return ctx.GetStub().DelState(dealerID)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, dealerID string) (bool, error) {
    assetJSON, err := ctx.GetStub().GetState(dealerID)
    if err != nil {
        return false, fmt.Errorf("failed to read from world state: %v", err)
    }

    return assetJSON != nil, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    var assets []*Asset
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }

        var asset Asset
        err = json.Unmarshal(queryResponse.Value, &asset)
        if err != nil {
            return nil, err
        }
        assets = append(assets, &asset)
    }

    return assets, nil
}

// GetAssetHistory returns the transaction history for an asset
func (s *SmartContract) GetAssetHistory(ctx contractapi.TransactionContextInterface, dealerID string) ([]map[string]interface{}, error) {
    resultsIterator, err := ctx.GetStub().GetHistoryForKey(dealerID)
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    var history []map[string]interface{}
    for resultsIterator.HasNext() {
        response, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }

        var asset Asset
        if len(response.Value) > 0 {
            err = json.Unmarshal(response.Value, &asset)
            if err != nil {
                return nil, err
            }
        }

        record := map[string]interface{}{
            "txId":      response.TxId,
            "timestamp": response.Timestamp,
            "isDeleted": response.IsDelete,
            "asset":     asset,
        }
        history = append(history, record)
    }

    return history, nil
}

func main() {
    assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
    if err != nil {
        log.Panicf("Error creating asset-chaincode: %v", err)
    }

    if err := assetChaincode.Start(); err != nil {
        log.Panicf("Error starting asset-chaincode: %v", err)
    }
}
