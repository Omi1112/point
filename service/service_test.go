package service

import (
	"net/http"
	"os"
	"testing"

	"github.com/SeijiOmi/points-service/db"
	"github.com/SeijiOmi/points-service/entity"
	"github.com/stretchr/testify/assert"
)

/*
	テストの前準備
*/

var client = new(http.Client)
var pointDefault = entity.Point{Number: 100, Comment: "testComment"}
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
}

// テスト実施後共通処理
func teardown() {
	os.Setenv("USER_URL", tmpBaseUserURL)
	initPointTable()
	db.Close()
}

/*
	ここからが個別のテスト実装
*/

func TestGetByUserID(t *testing.T) {
	initPointTable()
	createDefaultPoint(1)
	createDefaultPoint(2)

	var b Behavior
	points, err := b.GetByUserID("1")
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(points))
}

func TestCreateModel(t *testing.T) {
	initPointTable()
	inputPoint := pointDefault

	var b Behavior
	point, err := b.CreateModel(inputPoint, "testToken")
	assert.Equal(t, nil, err)
	assert.NotEqual(t, 0, point.ID)
	assert.NotEqual(t, 0, point.UserID)
	assert.Equal(t, pointDefault.Number, point.Number)
	assert.Equal(t, pointDefault.Comment, point.Comment)
}

func TestGetSumNumberByUserID(t *testing.T) {
	initPointTable()
	createDefaultPoint(1)
	createDefaultPoint(1)
	createDefaultPoint(1)
	createDefaultPoint(2)

	db := db.GetDB()
	point := pointDefault
	point.UserID = 1
	point.Number = -100
	db.Create(&point)

	var b Behavior
	total, err := b.GetSumNumberByUserID("1")
	assert.Equal(t, nil, err)
	assert.Equal(t, 200, total)
}

func TestGetSumNumberByUserIDUnknownUserID(t *testing.T) {
	initPointTable()
	createDefaultPoint(1)
	createDefaultPoint(2)

	var b Behavior
	total, err := b.GetSumNumberByUserID("3")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, 0, total)
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
