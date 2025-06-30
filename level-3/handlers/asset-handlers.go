package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"asset-api/fabric"
	"asset-api/models"

	"github.com/gin-gonic/gin"
)

type AssetHandler struct {
	FabricGateway *fabric.FabricGateway
}

func NewAssetHandler(fg *fabric.FabricGateway) *AssetHandler {
	return &AssetHandler{
		FabricGateway: fg,
	}
}

// InitLedger initializes the ledger with sample data
func (ah *AssetHandler) InitLedger(c *gin.Context) {
	_, err := ah.FabricGateway.Contract.SubmitTransaction("InitLedger")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ledger initialized successfully"})
}

// CreateAsset creates a new asset
func (ah *AssetHandler) CreateAsset(c *gin.Context) {
	var req models.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := ah.FabricGateway.Contract.SubmitTransaction(
		"CreateAsset",
		req.DealerID,
		req.MSISDN,
		req.MPIN,
		strconv.FormatFloat(req.Balance, 'f', 2, 64),
		req.Status,
		req.TransType,
		req.Remarks,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Asset created successfully"})
}

// GetAsset retrieves an asset by dealer ID
func (ah *AssetHandler) GetAsset(c *gin.Context) {
	dealerID := c.Param("dealerId")

	result, err := ah.FabricGateway.Contract.EvaluateTransaction("ReadAsset", dealerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var asset models.Asset
	err = json.Unmarshal(result, &asset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse asset data"})
		return
	}

	c.JSON(http.StatusOK, asset)
}

// UpdateAsset updates an existing asset
func (ah *AssetHandler) UpdateAsset(c *gin.Context) {
	dealerID := c.Param("dealerId")

	var req models.UpdateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := ah.FabricGateway.Contract.SubmitTransaction(
		"UpdateAsset",
		dealerID,
		strconv.FormatFloat(req.Balance, 'f', 2, 64),
		req.Status,
		strconv.FormatFloat(req.TransAmount, 'f', 2, 64),
		req.TransType,
		req.Remarks,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset updated successfully"})
}

// GetAllAssets retrieves all assets
func (ah *AssetHandler) GetAllAssets(c *gin.Context) {
	result, err := ah.FabricGateway.Contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var assets []models.Asset
	err = json.Unmarshal(result, &assets)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse assets data"})
		return
	}

	c.JSON(http.StatusOK, assets)
}

// DeleteAsset deletes an asset by dealer ID
func (ah *AssetHandler) DeleteAsset(c *gin.Context) {
	dealerID := c.Param("dealerId")

	_, err := ah.FabricGateway.Contract.SubmitTransaction("DeleteAsset", dealerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset deleted successfully"})
}

// GetAssetHistory retrieves the transaction history for an asset
func (ah *AssetHandler) GetAssetHistory(c *gin.Context) {
	dealerID := c.Param("dealerId")

	result, err := ah.FabricGateway.Contract.EvaluateTransaction("GetAssetHistory", dealerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var history []map[string]interface{}
	err = json.Unmarshal(result, &history)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse history data"})
		return
	}

	c.JSON(http.StatusOK, history)
}
