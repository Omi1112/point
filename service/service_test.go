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
	os.Setenv("USER_URL", "http://post-mock-user:3000")
	db.Init()
	initPostTable()
}

// テスト実施後共通処理
func teardown() {
	os.Setenv("USER_URL", tmpBaseUserURL)
	initPostTable()
	db.Close()
}

/*
	ここからが個別のテスト実装
*/

func TestGetByHelperUserID(t *testing.T) {
	initPostTable()
	createDefaultPoint(1)
	createDefaultPoint(1)

	var b Behavior
	points, err := b.GetByHelperUserID("1")
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(points))
}

func createDefaultPoint(userID uint) entity.Point {
	db := db.GetDB()
	point := pointDefault
	point.UserID = userID

	db.Create(&point)
	return point
}

func initPostTable() {
	db := db.GetDB()
	var u entity.Point
	db.Delete(&u)
}
