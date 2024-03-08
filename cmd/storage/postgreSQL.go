package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQL struct {
	DB *gorm.DB
}

func NewPostgreSQL() *PostgreSQL {
	return &PostgreSQL{}
}

func (s *PostgreSQL) ConnToDB() {
	var err error
	dsn := os.Getenv("DB_URL")
	s.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect PostgreSQL DB: " + err.Error())
	}
	fmt.Println("Successfully connected to PostgreSQL DB")
}

func (s *PostgreSQL) UserExist(username string, user interface{}) error {
	result := s.DB.Where("username = ?", username).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil // User with the given username does not exist
		}
		return result.Error // Some other error occurred
	}
	return fmt.Errorf("user with username: %s already exist", username) // User exist
}

func (s *PostgreSQL) EmailExist(email string, user interface{}) error {
	result := s.DB.Where("email = ?", email).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil // User with the given username does not exist
		}
		return result.Error // Some other error occurred
	}
	return fmt.Errorf("user with email: %s already exist", email) // User exist
}

func (s *PostgreSQL) SaveUser(user interface{}) error {
	tx := s.DB.Begin()
	defer tx.Rollback()
	if tx.Error != nil {
		log.Println("TX error: ", tx.Error)
		return tx.Error
	}
	if err := tx.Create(user).Error; err != nil {
		log.Println("TX error: ", err)
		return err
	}
	// Commit the transaction if everything is successful
	if err := tx.Commit().Error; err != nil {
		log.Println("TX error: ", err)
		return err
	}
	return nil
}

func (s *PostgreSQL) GetUserByUsername(username string, user interface{}) error {
	result := s.DB.Where("username = ?", username).First(user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("user with username '%s' not found", username)
	} else if result.Error != nil {
		return fmt.Errorf("can't retrieving user: %v", result.Error)
	}
	return nil
}

func (s *PostgreSQL) GetUserByEmail(email string, user interface{}) error {
	result := s.DB.Where("email = ?", email).First(user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("user with email '%s' not found", email)
	} else if result.Error != nil {
		return fmt.Errorf("can't retrieving user: %v", result.Error)
	}
	return nil
}

func (s *PostgreSQL) GetUserByID(id any, user interface{}) error {
	result := s.DB.Where("id = ?", id.(uint)).First(user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("user with id '%d' not found", id)
	} else if result.Error != nil {
		return fmt.Errorf("can't retrieving user: %v", result.Error)
	}
	return nil
}

func (s *PostgreSQL) VerifyAccount(email string, verTime sql.NullTime, user interface{}) error {
	result := s.DB.Model(user).Where("email = ?", email).Update("verified_at", verTime)
	if result.Error != nil {
		return fmt.Errorf("email verification isn't complete: %v", result.Error)
	}
	return nil
}

func (s *PostgreSQL) UpdatePswdHash(newPswdHash string, id any, user interface{}) error {
	result := s.DB.Model(user).Where("id = ?", id.(uint)).Update("pswd_hash", newPswdHash)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("user with id '%d' not found", id.(uint))
	} else if result.Error != nil {
		return fmt.Errorf("can't retrieving user: %v", result.Error)
	}
	return nil
}

func (s *PostgreSQL) UpdateVerHash(newVerHash string, t time.Time, id any, user interface{}) error {
	// result := s.DB.Model(user).Where("id = ?", id.(uint)).Updates(u{VerHash: verHash, Timeout: t})
	result := s.DB.Model(user).Where("id = ?", id.(uint)).Updates(map[string]interface{}{"ver_hash": newVerHash, "timeout": t})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("DB ERROR: user with id '%d' not found", id.(uint))
	} else if result.Error != nil {
		return fmt.Errorf("DB ERROR retrieving user: %v", result.Error)
	}
	// result := s.DB.Model(user).Where("id = ?", id.(uint)).Update("pswd_hash", newPswdHash)
	// if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	// 	return fmt.Errorf("user with id '%d' not found", id.(uint))
	// } else if result.Error != nil {
	// 	return fmt.Errorf("can't retrieving user: %v", result.Error)
	// }
	return nil
}
