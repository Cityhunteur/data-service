package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	v1 "github.com/cityhunteur/data-service/api/v1"
	"github.com/cityhunteur/data-service/internal/database"
)

type DataHandler struct {
	db *database.DataDB
}

// NewDataHandler creates API handlers for the data API.
func NewDataHandler(_ context.Context, db *database.DataDB) (*DataHandler, error) {
	if db == nil {
		return nil, fmt.Errorf("db must be non-nil")
	}
	return &DataHandler{
		db: db,
	}, nil
}

func (dh *DataHandler) GetData(c *gin.Context) {
	title := c.DefaultQuery("title", "")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Query param 'title' missing."})
		return
	}
	d, err := dh.db.GetData(c, title)
	if errors.Is(err, database.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "No data record found with specified title"})
		return
	}
	c.JSON(http.StatusOK, &v1.DataResponse{
		ID:        d.ID,
		Title:     d.Title,
		Timestamp: v1.Time3339(d.Timestamp),
	})
}

func (dh *DataHandler) PostData(c *gin.Context) {
	var d v1.Data
	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if d.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request field 'title' missing."})
		return
	}
	record, err := dh.db.InsertData(c, d.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, &v1.DataResponse{
		ID:        record.ID,
		Title:     record.Title,
		Timestamp: v1.Time3339(record.Timestamp),
	})
}
