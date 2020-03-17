package service

import (
	"errors"
	"net/http"
	"os"

	"github.com/SeijiOmi/points-service/db"
	"github.com/SeijiOmi/points-service/entity"
	"github.com/jmcvetta/napping"
)

// Behavior ポイントサービスを提供するメソッド群
type Behavior struct{}

// GetByUserID ユーザIDを元にポイント情報一覧を取得する。
func (b Behavior) GetByUserID(userID string) ([]entity.Point, error) {
	db := db.GetDB()
	var point []entity.Point

	if err := db.Where("user_id = ?", userID).Find(&point).Error; err != nil {
		return point, err
	}

	return point, nil
}

// CreateModel ポイント情報の登録
func (b Behavior) CreateModel(inputPoint entity.Point, token string) (entity.Point, error) {
	userID, err := getUserIDByToken(token)
	if err != nil {
		return inputPoint, err
	}

	createPoint := inputPoint
	createPoint.UserID = uint(userID)
	db := db.GetDB()

	if err := db.Create(&createPoint).Error; err != nil {
		return createPoint, err
	}

	return createPoint, nil
}

// GetSumNumberByUserID ユーザIDを元にポイント所持数を取得する。
func (b Behavior) GetSumNumberByUserID(id string) (int, error) {
	db := db.GetDB()
	var total int

	err := db.Table("points").
		Select("sum(number) as total").
		Where("user_id = ?", id).
		Group("user_id").Row().
		Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func getUserIDByToken(token string) (int, error) {
	response := struct {
		ID int
	}{}
	error := struct {
		Error string
	}{}

	baseURL := os.Getenv("USER_URL")
	resp, err := napping.Get(baseURL+"/auth/"+token, nil, &response, &error)

	if err != nil {
		return 0, err
	}

	if resp.Status() == http.StatusBadRequest {
		return 0, errors.New("token invalid")
	}

	return response.ID, nil
}
