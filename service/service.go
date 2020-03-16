package service

import (
	"errors"
	"net/http"
	"os"

	"github.com/SeijiOmi/points-service/db"
	"github.com/SeijiOmi/points-service/entity"
	"github.com/jmcvetta/napping"
)

// Behavior 投稿サービスを提供するメソッド群
type Behavior struct{}

// User オブジェクト構造
type User struct {
	id   int
	name string
}

// GetByHelperUserID 投稿情報にユーザ情報を紐づけて取得
func (b Behavior) GetByHelperUserID(userID string) ([]entity.Post, error) {
	db := db.GetDB()
	var post []entity.Post

	if err := db.Where("helper_user_id = ?", userID).Find(&post).Error; err != nil {
		return post, err
	}

	return post, nil
}

// CreateModel 投稿情報の生成
func (b Behavior) CreateModel(inputPost entity.Post, token string) (entity.Post, error) {
	userID, err := getUserIDByToken(token)
	if err != nil {
		return inputPost, err
	}

	createPost := inputPost
	createPost.UserID = uint(userID)
	db := db.GetDB()

	if err := db.Create(&createPost).Error; err != nil {
		return createPost, err
	}

	return createPost, nil
}

// GetByID IDを元に投稿1件を取得
func (b Behavior) GetByID(id string) (entity.Post, error) {
	db := db.GetDB()
	var u entity.Post

	if err := db.Where("id = ?", id).First(&u).Error; err != nil {
		return u, err
	}

	return u, nil
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
