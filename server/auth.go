package main

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"net/http"
	"time"

	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"log"
)

const (
	GinContextField = "Session"
	XSRFCookieName = "XSRF-TOKEN"
	TokenHeaderField = "X-XSRF-TOKEN"
	NameRequestField = "Email"
	PassRequestField = "Password"
)

var (
	SignInErr = errors.New("Sign in error")
)

type (
	SuccessResponse struct {
		Status string
		Data   interface{}
	}

	FailResponse struct {
		Status string
		Err    string
	}

	UserIDData struct {
		ID string
	}

	Session struct {
		Token   string    `bson:"Token"`
		UserID  string    `bson:"UserID"`
		Expires time.Time `bson:"Expires"`
	}

	User interface {
		ID() string
		PASSWORD() string
	}

	FindUser        func(string) (User, error)
	ConvertPassword func(string) string
)

func NewSuccessResponse(data interface{}) SuccessResponse {
	return SuccessResponse{
		Status: "success",
		Data:   data,
	}
}

func NewFailResponse(err interface{}) FailResponse {
	return FailResponse{
		Status: "fail",
		Err:    fmt.Sprintf("%v", err),
	}
}

func NewSessionToken() (string, error) {
	buf := make([]byte, 2)

	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	c := sha256.New()
	hash := fmt.Sprintf("%x", c.Sum(buf))

	return hash, nil
}

func NewSha512Password(pass string) string {
	hash := sha512.New()
	tmp := hash.Sum([]byte(pass))
	passHash := fmt.Sprintf("%x", tmp)
	return passHash
}

func ReadSession(ctx *gin.Context) (Session, error) {
	v, ok := ctx.Get(GinContextField)
	if !ok {
		return Session{}, errors.New("Wrong session in cookie")
	}

	s, ok := v.(Session)
	if !ok {
		return Session{}, errors.New("Wrong session in cookie")
	}

	return s, nil
}

// Middleware Decorator:
// Handles Angularjs Default Authentication
// Sendet man über den angular http Serviecs ein Request und erhält
// daraufhin ein Response mit einem Cookie welcher ein XSRF-Token Feld
// enthält wird der hinterlegte Token für zukünftige Request verwendet.
// Der Token wird als HTTP-Header-Feld X-XSRF-Token versand. Diesen
// Eigenschaft kann man für die Benutzer Authentifikation verwenden.
//
// Die Middleware fügt ein Feld Session zum gin Context.
//
// Die Middleware erwartet ein Session Collection mit den selben
// Feldern wie der Session Typ
//
// Example:
// app := gin.New()
//
// func protectedHandler(c *gin.Context) {
//      // Access only for succesfully authenticated user
// }
//
// s,_ := db.Dial("mongodb://127.0.0.1:27017")
// db := s.DB("DBName")
// auth := AngularAuth(*mgo.Database, "SessionCollName")
// app.GET(auth, portectedHandler)
//
func AngularAuth(db *bolt.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := Auther(c, db)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				NewFailResponse(err))
			c.Abort()
		}
	}
}

func Auther(c *gin.Context, db *bolt.DB) error {
	token := c.Request.Header.Get(TokenHeaderField)
	if token == "" {
		cookie, err := c.Request.Cookie(XSRFCookieName)
		if err != nil {
			return errors.New("Cookie not found")
		}
		token = cookie.Value
		if token == "" {
			return errors.New("Header not found")
		}
	}

	var sess Session
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(SessionsBucket)
		buf := b.Get([]byte(token))

		return json.Unmarshal(buf, &sess)
	})
	if err != nil {
		log.Println(err)
		return nil
	}

	if &sess == nil {
		return errors.New("Session not found")

	}

	if sess.Expires.Before(time.Now()) {
		return errors.New("Session expired")
	}

	c.Set(GinContextField, sess)
	c.Next()
	return nil
}

func AngularSignIn(coll *bolt.DB, findUser FindUser, cPass ConvertPassword, expireTime time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := Signer(c, coll, findUser, cPass, expireTime)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				NewFailResponse(err))
		}
	}
}

func Signer(c *gin.Context, db *bolt.DB, findUser FindUser, convertPassword ConvertPassword, expireTime time.Duration) error {

	type auth struct {
		Email    string
		Password string
	}

	var a auth
	err := c.BindJSON(&a)
	if err != nil {
		return err
	}

	passHash := convertPassword(a.Password)

	log.Println(passHash, a)

	user, err := findUser(a.Email)
	if err != nil {
		return err
	}

	if user.PASSWORD() != passHash {
		return SignInErr
	}

	resp := NewSuccessResponse(user.(*Publisher))

	sessionToken, err := NewSessionToken()
	if err != nil {
		return err
	}

	expire := time.Now().Add(expireTime)

	session := Session{
		UserID:  user.ID(),
		Token:   sessionToken,
		Expires: expire,
	}

	buf, err := json.Marshal(session)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(SessionsBucket)
		err := b.Put([]byte(sessionToken), buf)
		return err
	})
	if err != nil {
		c.Error(err)
	}

	cookie := http.Cookie{
		Name:    XSRFCookieName,
		Value:   sessionToken,
		Expires: expire,
		// Setze Path auf / ansonsten kann angularjs
		// diese Cookie nicht finden und in späteren
		// Request nicht mitsenden.
		Path: "/",
	}

	http.SetCookie(c.Writer, &cookie)

	c.JSON(http.StatusOK, resp)

	return nil
}
