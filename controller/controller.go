package controller

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"

	"github.com/SeijiOmi/points-service/entity"
	"github.com/SeijiOmi/points-service/service"
)

// Create action: POST /points
func Create(c *gin.Context) {
	var inputPost entity.Post
	if err := bindJSON(c, &inputPost); err != nil {
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
	createdPost, err := b.CreateModel(inputPost, token.Token)

	if err != nil {
		c.AbortWithStatus(400)
		fmt.Println(err)
	} else {
		c.JSON(201, createdPost)
	}
}

// Show action: GET /points/:id
func Show(c *gin.Context) {
	id := c.Params.ByName("id")
	var b service.Behavior
	p, err := b.GetByID(id)

	if err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, p)
	}
}

// Sum action: GET /sum/:id
func Sum(c *gin.Context) {
	id := c.Params.ByName("id")
	var b service.Behavior
	p, err := b.GetByID(id)

	if err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, p)
	}
}

func bindJSON(c *gin.Context, data interface{}) error {
	buf := make([]byte, 2048)
	n, _ := c.Request.Body.Read(buf)
	b := string(buf[0:n])
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(b)))
	if err := c.BindJSON(data); err != nil {
		c.AbortWithStatus(400)
		fmt.Println("bind JSON err")
		fmt.Println(err)
		return err
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(b)))
	return nil
}
