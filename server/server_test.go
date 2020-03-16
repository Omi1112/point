package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/SeijiOmi/points-service/db"
	"github.com/SeijiOmi/points-service/entity"
	"github.com/jmcvetta/napping"
	"github.com/stretchr/testify/assert"
)

/*
	テストの前準備
*/

var client = new(http.Client)
var testServer *httptest.Server
var pointDefault = entity.Point{Number: 100}
var tmpBaseUserURL string

// テストを統括するテスト時には、これが実行されるイメージでいる。
func TestMain(m *testing.M) {
	// テスト実施前の共通処理（自作関数）
	setup()
	ret := m.Run()
	// テスト実施後の共通処理（自作関数）
	teardown()
	os.Exit(ret)
}

// テスト実施前共通処理
func setup() {
	tmpBaseUserURL = os.Getenv("USER_URL")
	os.Setenv("USER_URL", "http://point-mock-user:3000")
	db.Init()
	initPointTable()
	router := router()
	testServer = httptest.NewServer(router)
}

// テスト実施後共通処理
func teardown() {
	testServer.Close()
	initPointTable()
	db.Close()
	os.Setenv("USER_URL", tmpBaseUserURL)
}

/*
	ここからが個別のテスト実装
*/

func TestPointsPost(t *testing.T) {
	inputPoint := struct {
		Number int    `json:"number"`
		Token  string `json:"token"`
	}{
		100,
		"testToken",
	}
	input, _ := json.Marshal(inputPoint)
	resp, _ := http.Post(testServer.URL+"/points", "application/json", bytes.NewBuffer(input))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestPointsGetByUserID(t *testing.T) {
	response := []entity.Point{}
	error := struct {
		Error string
	}{}

	initPointTable()
	createDefaultPoint(1)
	createDefaultPoint(1)

	resp, err := napping.Get(testServer.URL+"/points/1", nil, &response, &error)
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, resp.Status())
	assert.Equal(t, 2, len(response))
}

func TestSumGetByUserID(t *testing.T) {
	response := struct {
		Total int
	}{}
	error := struct {
		Error string
	}{}

	initPointTable()
	createDefaultPoint(1)
	createDefaultPoint(1)

	resp, err := napping.Get(testServer.URL+"/sum/1", nil, &response, &error)
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, resp.Status())
	assert.Equal(t, 200, response.Total)
}

func createDefaultPoint(userID uint) entity.Point {
	db := db.GetDB()
	point := pointDefault
	point.UserID = userID

	db.Create(&point)
	return point
}

func initPointTable() {
	db := db.GetDB()
	var u entity.Point
	db.Delete(&u)
}
