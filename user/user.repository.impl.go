package users

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type userRepositoryImpl struct {
	DB *gorm.DB
}

// NewUserRepository for entry point of user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{DB: db}
}

func (repository *userRepositoryImpl) Insert(user Users) (Users, error) {
	repository.DB.Create(&user)

	return user, nil
}

func (repository *userRepositoryImpl) Read(id uuid.UUID) (Users, error) {
	var user Users

	whereArgs := map[string]interface{}{"id": id}

	result := repository.DB.Table("users").Take(&user, whereArgs)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, errors.New("user not found")
	}

	return user, nil
}

func (repository *userRepositoryImpl) Update(id uuid.UUID, user Users) (Users, error) {
	user.ID = id

	result := repository.DB.Model(&user).Update(user).First(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (repository *userRepositoryImpl) Delete(id uuid.UUID) bool {
	var user Users

	user.ID = id

	result := repository.DB.Delete(&user)
	fmt.Println(result)
	if result.Error != nil {
		return false
	}

	return true
}

func (repository *userRepositoryImpl) ValidateUsernameOrEmail(username string, email string) bool {
	var user Users

	result := repository.DB.Where("username = ?", username).Or("email = ?", email).Find(&user)

	return result.Error == nil
}
