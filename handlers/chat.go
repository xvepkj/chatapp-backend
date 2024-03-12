package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"github.com/xvepkj/chatapp-backend/db"
	"github.com/xvepkj/chatapp-backend/models"
	"gorm.io/gorm"
)

func AddMessage(c *gin.Context, dbConn *gorm.DB) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.AddMessage(dbConn, &message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "message added successfully", "user": message})
}

func AddMessageWebSocket(dbConn *gorm.DB, message *models.Message) error {

	if err := db.AddMessage(dbConn, message); err != nil {
		return err
	}

	return nil
}

func GetMessagesBetween(c *gin.Context, dbConn *gorm.DB) {
	senderID := c.Param("senderID")
	receiverId := c.Param("receiverID")

	messages, err := db.GetMessagesBetween(dbConn, senderID, receiverId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "messages not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// ExportMessagesToExcel exports messages to an Excel file and sends it in the response
func ExportMessagesToExcel(c *gin.Context, dbConn *gorm.DB) {
	senderID := c.Param("senderID")
	receiverId := c.Param("receiverID")

	messages, err := db.GetMessagesBetween(dbConn, senderID, receiverId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "messages not found"})
		return
	}

	// Generate Excel file
	file, err := messagesToExcel(messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate Excel file"})
		return
	}

	// Set headers for Excel file download
	c.Header("Content-Disposition", "attachment; filename=messages.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// Write Excel file content to response
	err = file.Write(c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write Excel file to response"})
		return
	}
}

func messagesToExcel(messages []models.Message) (*excelize.File, error) {
	// Create a new Excel file
	file := excelize.NewFile()

	// Create a new sheet in the Excel file
	sheet := "Messages"
	file.NewSheet(sheet)

	// Write headers to the sheet
	headers := []string{"ID", "SenderID", "RecipientID", "Content", "Timestamp"}
	for i, header := range headers {
		col := string(rune('A'+i)) + "1"
		file.SetCellValue(sheet, col, header)
	}

	// Write messages to the sheet
	for i, msg := range messages {
		row := i + 2
		file.SetCellValue(sheet, "A"+strconv.Itoa(row), msg.ID)
		file.SetCellValue(sheet, "B"+strconv.Itoa(row), msg.SenderID)
		file.SetCellValue(sheet, "C"+strconv.Itoa(row), msg.ReceipientID)
		file.SetCellValue(sheet, "D"+strconv.Itoa(row), msg.Content)
		file.SetCellValue(sheet, "E"+strconv.Itoa(row), msg.Timestamp.String())
	}

	return file, nil
}
