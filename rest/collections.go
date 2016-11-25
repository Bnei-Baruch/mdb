package rest

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

type CreateCollection struct {
	Type        string `json:"type" binding:"required"`
	UID         string `json:"uid" binding:"max=8"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Language    string `json:"language" binding:"max=2"`
}

func CollectionsCreateHandler(c *gin.Context) {
	var cc CreateCollection
	if c.BindJSON(&cc) == nil {
		cl := new(models.Collection)
		cl.TypeID = 2
		cl.UID = utils.GenerateUID(8)
		cl.Name.Text = cc.Name
		cl.Name.Language = cc.Language
		cl.Description.Text = cc.Description
		cl.Description.Language = cc.Language
		c.JSON(http.StatusCreated, cl)
	}

}

