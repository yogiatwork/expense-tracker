package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yogiatwork/expense-tracker/model"
	"github.com/yogiatwork/expense-tracker/service"
)

func GetSheetMetadataHandler(c *gin.Context) {
	err := service.SheetService.GetSheetMetadata()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sheet metadata retrieved successfully"})
}

func ExpenseEntryHandler(c *gin.Context) {
	exp := model.ExpenseEntry{}
	if err := c.BindJSON(&exp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := service.SheetService.AppendSheet(exp.ToTransaction())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data appended successfully"})
}

func IncomeEntryHandler(c *gin.Context) {
	earn := model.EarnedEntry{}
	if err := c.BindJSON(&earn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := service.SheetService.AppendSheet(earn.ToTransaction())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data appended successfully"})
}
