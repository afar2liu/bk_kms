package repo

import (
	"bk_kms/lib"
	"bk_kms/model/db"
)

type UserRepo struct{}

// FindByUsername 根据用户名查找用户
func (r *UserRepo) FindByUsername(username string) (*db.User, error) {
	var user db.User
	err := lib.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *UserRepo) Create(user *db.User) error {
	return lib.DB.Create(user).Error
}

// FindByID 根据ID查找用户
func (r *UserRepo) FindByID(id int) (*db.User, error) {
	var user db.User
	err := lib.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
