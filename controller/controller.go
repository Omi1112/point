package controller

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/SeijiOmi/points-service/entity"
	"github.com/SeijiOmi/points-service/service"
	"github.com/gin-gonic/gin"
)

// Create action: POST /points
func Create(c *gin.Context) {
	var inputPoint entity.Point
	if err := bindJSON(c, &inputPoint); err != nil {
		return
	}
	type tokenStru struct {
		Token string `json:"token"`
	}
	var token tokenStru
	if err := bindJSON(c, &token); err != nil {
		return
	}

	var b service.Behavior
	createdPoint, err := b.CreateModel(inputPoint, token.Token)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusCreated, createdPoint)
	}
}

// Show action: GET /points/:id
func Show(c *gin.Context) {
	id := c.Params.ByName("id")
	var b service.Behavior
	p, err := b.GetByUserID(id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusOK, p)
	}
}

// Sum action: GET /sum/:id
func Sum(c *gin.Context) {
	id := c.Params.ByName("id")
	var b service.Behavior
	p, err := b.GetSumNumberByUserID(id)

	response := struct {
		Total int
	}{
		p,
	}

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func bindJSON(c *gin.Context, data interface{}) error {
	buf := make([]byte, 2048)
	n, _ := c.Request.Body.Read(buf)
	b := string(buf[0:n])
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(b)))
	if err := c.BindJSON(data); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println("bind JSON err")
		fmt.Println(err)
		return err
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(b)))
	return nil
}
