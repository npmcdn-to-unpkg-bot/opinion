package publisher

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"net/http"
	"time"

	"encoding/json"
	"log"

	"github.com/boltdb/bolt"


	"github.com/kataras/iris"
)

const (
	IrisContextField = "Session"
	XSRFCookieName   = "XSRF-TOKEN"
	TokenHeaderField = "X-XSRF-TOKEN"
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

func ReadSession(ctx *iris.Context) (Session, error) {
	v := ctx.Get(IrisContextField)
	if v == nil {
		return Session{}, errors.New("Wrong session in cookie")
	}

	s, ok := v.(Session)
	if !ok {
		return Session{}, errors.New("Wrong session in cookie")
	}

	return s, nil
}

func AngularAuth(db *bolt.DB) iris.HandlerFunc {
	return func(c *iris.Context) {
		err := Auther(c, db)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.JSON(NewFailResponse(err))
			return
		}
	}
}

func Auther(c *iris.Context, db *bolt.DB) error {
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

	c.Set(IrisContextField, sess)
	c.Next()
	return nil
}

func AngularSignIn(coll *bolt.DB, findUser FindUser, cPass ConvertPassword, expireTime time.Duration) iris.HandlerFunc {
	return func(c *iris.Context) {
		err := Signer(c, coll, findUser, cPass, expireTime)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.JSON(NewFailResponse(err))
		}
	}
}

func Signer(c *iris.Context, db *bolt.DB, findUser FindUser, convertPassword ConvertPassword, expireTime time.Duration) error {

	type auth struct {
		Email    string
		Password string
	}

	var a auth
	err := c.ReadJSON(&a)
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
		c.EmitError(500)
		c.Write(err.Error())
		return err
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

	http.SetCookie(c.ResponseWriter, &cookie)

	c.JSON(resp)

	return nil
}
