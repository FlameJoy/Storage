package models

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"testezhik/cmd/storage"
	"time"
	"unicode"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/golang-jwt/jwt/v5"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func init() {
	verifier = verifier.EnableDomainSuggest()
	if serveLoc == "remote" {
		verifier.EnableSMTPCheck()
	}
	dispEmailsDomains := mustDispEmailDomains()
	verifier = verifier.AddDisposableDomains(dispEmailsDomains)
}

var (
	minEntropyBits, _ = strconv.Atoi(os.Getenv("minEntropyBits"))
	verifier          = emailverifier.NewVerifier()
	serveLoc          = os.Getenv("server_location")
)

type User struct {
	gorm.Model
	MongoID    string       `bson:"_id,omitempty"`
	Username   string       `bson:"username"`
	Email      string       `bson:"email"`
	PswdHash   string       `bson:"pswdHash"`
	VerifiedAt sql.NullTime `bson:"verifiedAt"`
	VerHash    string       `bson:"verHash"`
	Timeout    time.Time    `bson:"timeout"`
}

func (u *User) ValidateUsername() error {
	for _, char := range u.Username {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return fmt.Errorf("username must contain only letters and numbers")
		}
	}
	if len(u.Username) > 4 && len(u.Username) < 25 {
		return nil
	}
	return fmt.Errorf("username must be greater than 4 and less than 25 characters")
}

func mustDispEmailDomains() (dispEmailDomains []string) {
	file, err := os.Open("../../disposable_email_blocklist.txt")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dispEmailDomains = append(dispEmailDomains, scanner.Text())
	}
	return dispEmailDomains
}

func (u *User) ValidateEmail() error {
	result, err := verifier.Verify(u.Email)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("email verification failed: %v", err)
	}
	if !result.Syntax.Valid {
		return errors.New("email syntax is invalid")
	}
	if result.Disposable {
		return errors.New("sorry we don't accept disposable email addresses")
	}
	if result.Suggestion != "" {
		return errors.New("email address is not reachable, looking for: " + result.Suggestion + " instead of " + result.Email + "?")
	}
	if result.Reachable == "no" {
		return errors.New("email address is unreachable")
	}
	if !result.HasMxRecords {
		return errors.New("domain entered not properly setup to recieve emails, MX record not found")
	}
	return nil
}

func (u *User) ValidatePswd(pswd1, pswd2 string) error {
	// Compare pswds
	for i := 0; i < len(pswd1); i++ {
		if pswd1[i] != pswd2[i] {
			return errors.New("passwords aren't comparable")
		}
	}
	// Validation
	err := passwordvalidator.Validate(pswd1, float64(minEntropyBits))
	if err != nil {
		return err
	}
	// Check length
	if len(pswd1) >= 8 && len(pswd1) <= 24 {
		return nil
	}
	return errors.New("password must be greater than 7 and less than 25 characters")
}

func (u *User) UserExist() error {
	err := storage.DB.UserExist(u.Username, &User{})
	return err // User exist
}

func (u *User) EmailExist() error {
	err := storage.DB.EmailExist(u.Email, &User{})
	return err // User exist
}

func (u *User) New(pswd string) error {
	// Create hash from pswd
	hash, err := bcrypt.GenerateFromPassword([]byte(pswd), bcrypt.DefaultCost)
	u.PswdHash = string(hash)
	if err != nil {
		log.Println("Hash generation error: ", err)
		return err
	}
	// Init rand source
	source := rand.NewSource(time.Now().Unix())
	rng := rand.New(source)
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	randBytes := make([]byte, 64)
	// Create a random slice of chars to create emailVerPswd
	for i := 0; i < 64; i++ {
		randBytes[i] = chars[rng.Intn(len(chars)-1)]
	}
	emailVerPswd := string(randBytes)
	b, err := bcrypt.GenerateFromPassword([]byte(emailVerPswd), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return err
	}
	u.VerHash = string(b)
	// Create user timeout after 24 hours
	u.Timeout = time.Now().Local().AddDate(0, 0, 1)
	if err != nil {
		fmt.Println("can't parse time:", err)
		return err
	}
	// Save user
	if err = storage.DB.SaveUser(u); err != nil {
		return err
	}
	// Send ver email
	if err = u.SendVerLink(emailVerPswd); err != nil {
		log.Println("Error: ", err)
		return err
	}
	return nil
}

func (u *User) SendVerLink(emailVerPswd string) error {
	// Send email
	var domName string
	var serveLoc = os.Getenv("server_location")
	var useHTTPS = os.Getenv("useHTTPS")
	var domainName = os.Getenv("domain_name")
	if serveLoc == "local" {
		domName = "http://localhost:8080"
	} else {
		if useHTTPS == "true" {
			domName = "https://" + domainName
		} else {
			domName = "http://" + domainName
		}
	}
	subject := "Email Verification"
	HTMLbody := `<html><h1>CLick link to verify email</h1><a href="` + domName + `/verify-email/` + u.Username + `/` + emailVerPswd + `">Click here to verify email</a></html>`
	err := u.SendEmail(subject, HTMLbody)
	if err != nil {
		fmt.Println("Can't send verification email")
		return err
	}
	return nil
}

func (u *User) SendEmail(subject, HTMLbody string) error {
	to := []string{u.Email}
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	var fromEmail = os.Getenv("from_Email_Addr")
	var smtpPswd = os.Getenv("SMTP_pswd")
	var entityName = os.Getenv("entity_name")
	auth := smtp.PlainAuth("", fromEmail, smtpPswd, host)
	msg := []byte(
		"From: " + entityName + ": <" + fromEmail + ">\r\n" +
			"To: " + u.Email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-version: 1.0\r\n" +
			"Content-type: text/html; charset=\"utf8\";\r\n" +
			"\r\n" +
			HTMLbody)
	err := smtp.SendMail(address, auth, fromEmail, to, msg)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetUserByUsername() error {
	err := storage.DB.GetUserByUsername(u.Username, u)
	return err
}

func (u *User) GetUserByID() error {
	var id any
	switch storage.DB.(type) {
	case *storage.PostgreSQL:
		id = u.ID
	case *storage.MondoDB:
		id = u.MongoID
	default:
		return fmt.Errorf("unknown database struct type: %v", storage.DB)
	}
	err := storage.DB.GetUserByID(id, u)
	return err
}

func (u *User) GetUserByEmail() error {
	err := storage.DB.GetUserByEmail(u.Email, u)
	return err
}

func (u *User) VerifyAccount() error {
	u.VerifiedAt.Time = time.Now().Local()
	u.VerifiedAt.Valid = true
	err := storage.DB.VerifyAccount(u.Email, u.VerifiedAt, u)
	return err
}

func (u *User) NewToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Local().Add(time.Hour * 24).Unix()
	// Additional data
	claims["userID"] = strconv.Itoa(int(u.ID))
	claims["userMongoID"] = u.MongoID
	secretKey := os.Getenv("secret")
	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (u *User) ChangePswd(newPswd string) error {
	b, err := bcrypt.GenerateFromPassword([]byte(newPswd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newPswdHash := string(b)
	u.PswdHash = newPswdHash
	var id any
	switch storage.DB.(type) {
	case *storage.PostgreSQL:
		id = u.ID
	case *storage.MondoDB:
		id = u.MongoID
	default:
		return fmt.Errorf("unknown database struct type: %v", storage.DB)
	}
	err = storage.DB.UpdatePswdHash(newPswdHash, id, u)
	return err
}

func (u *User) NewEmailVerPswd() error {
	source := rand.NewSource(time.Now().Unix())
	rng := rand.New(source)
	alphaNumRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	randRunes := make([]rune, 64)
	for i := 0; i < len(randRunes); i++ {
		randRunes[i] = alphaNumRunes[rng.Intn(len(alphaNumRunes)-1)]
	}
	emailVerPswd := string(randRunes)
	b, err := bcrypt.GenerateFromPassword([]byte(emailVerPswd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	verHash := string(b)
	// Insert new data
	t := time.Now()
	t = t.Add(2 * time.Hour)
	var id any
	switch storage.DB.(type) {
	case *storage.PostgreSQL:
		id = u.ID
	case *storage.MondoDB:
		id = u.MongoID
	default:
		return fmt.Errorf("unknown database struct type: %v", storage.DB)
	}
	storage.DB.UpdateVerHash(verHash, t, id, u)
	// Send email
	var domName string
	var serveLoc = os.Getenv("server_location")
	var useHTTPS = os.Getenv("useHTTPS")
	var domainName = os.Getenv("domain_name")
	if serveLoc == "local" {
		domName = "http://localhost:8080"
	} else {
		if useHTTPS == "true" {
			domName = "https://" + domainName
		} else {
			domName = "http://" + domainName
		}
	}
	subject := "Account recovery"
	HTMLbody := `<html><h1>CLick link to change password</h1><a href="` + domName + `/account-recovery/` + u.Username + `/` + emailVerPswd + `">Click here to change password</a></html>`
	err = u.SendEmail(subject, HTMLbody)
	if err != nil {
		return err
	}
	return nil
}
