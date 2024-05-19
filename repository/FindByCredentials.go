package repository

import (
	"github.com/AdluAghnia/nyoba-fiber/connection"
	"github.com/AdluAghnia/nyoba-fiber/models"
	"golang.org/x/crypto/bcrypt"
)

func FindByCredentials(username, password string) (*models.User, error) {
	user := models.User{}
	conn, err := connection.InitiliazedDB()
	if err != nil {
		return nil, err
	}

	err = conn.QueryRow("SELECT id, Name, Passwd FROM User WHERE Name = ?", username).
		Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	return &user, nil
}
