package handlers

import (
	"eda/database"
	"eda/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func ListFacts(c *gin.Context) {
	var facts []models.Fact
	database.DB.Db.Find(&facts)

	c.JSON(http.StatusOK, facts)
}

func CreateFact(c *gin.Context) {
	fact := new(models.Fact)
	data, _ := io.ReadAll(c.Request.Body)
	if e := json.Unmarshal(data, &fact); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": e.Error()})
		return
	}
	database.DB.Db.Create(&fact)
	c.JSON(http.StatusOK, fact)
}
