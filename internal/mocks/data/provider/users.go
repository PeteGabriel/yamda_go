package provider

import "yamda_go/internal/models"

type UserProviderMock struct {
	InsertMock     func(*models.User) (*models.User, error)
	GetByEmailMock func(string) (*models.User, error)
	UpdateMock     func(user *models.User) error
	DeleteMock     func(id int64) error
}

func (u UserProviderMock) Insert(user *models.User) (*models.User, error) {
	return u.InsertMock(user)
}

func (u UserProviderMock) GetByEmail(email string) (*models.User, error) {
	return u.GetByEmailMock(email)
}

func (u UserProviderMock) Update(user *models.User) error {
	return u.UpdateMock(user)
}

func (u UserProviderMock) Delete(id int64) error {
	return u.DeleteMock(id)
}
