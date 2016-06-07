package main

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/dgryski/dgoogauth"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
)

const (
	OTPConfigWindowSize = 10
	OTPConfigUTC        = true
)

// JSONResponse ajax response
type JSONResponse struct {
	Ok  bool   `json:"ok"`
	Txt string `json:"txt"`
}

// generate a new key
func generateSecret() (secret string, err error) {
	var key []byte
	key = make([]byte, sha1.Size)
	_, err = rand.Read(key)
	if err != nil {
		return
	}
	secret = base32.StdEncoding.EncodeToString(key)
	return
}

// Step 1: init OTPconfig en generate QRRCode
func handlergenQrCode(c *gin.Context) {
	var err error
	data := struct {
		user string `binding:"required"`
	}{}
	if err = c.Bind(&data); err != nil {
		return
	}

	// OTPConfig
	secret, err := generateSecret()
	if err != nil {
		return
	}

	fmt.Println(secret)

	config := dgoogauth.OTPConfig{
		Secret:     secret,
		WindowSize: OTPConfigWindowSize,
		UTC:        OTPConfigUTC,
	}

	// gen qrcode
	hasher := md5.New()
	hasher.Write([]byte(data.user))
	filename := hex.EncodeToString(hasher.Sum(nil)) + ".png"
	// gen url
	url := config.ProvisionURIWithIssuer(data.user, "otpDemo")
	log.Println(url)
	// gen & save qrcode
	err = qrcode.WriteFile(url, qrcode.Medium, 256, "./public/qrcodes/"+filename)
	// save OTPConfig in session (don't do this at /home !)
	session := sessions.Default(c)
	session.Set("optconfig", config)
	err = session.Save()
	fmt.Println(err)

	c.JSON(200, JSONResponse{true, filename})
}

// step 3 check code
func handlerCheckCode(c *gin.Context) {
	var err error
	data := struct {
		code string `binding:"required"`
	}{}
	if err = c.Bind(&data); err != nil {
		return
	}
	session := sessions.Default(c)
	configInt := session.Get("optconfig")
	if configInt == nil {
		c.JSON(400, JSONResponse{false, "no OTPConfing found on session"})
		return
	}
	config := configInt.(dgoogauth.OTPConfig)
	log.Println(configInt)
	valid, err := config.Authenticate(data.code)
	fmt.Println(err)
	if err != nil && err != dgoogauth.ErrInvalidCode {
		c.JSON(400, JSONResponse{false, err.Error()})
		return
	}
	c.JSON(200, JSONResponse{true, fmt.Sprintf("%t", valid)})
}

func main() {
	gob.Register(dgoogauth.OTPConfig{})
	r := gin.Default()
	store := sessions.NewCookieStore([]byte("bigsecret"))
	r.Use(sessions.Sessions("otp", store))
	r.Static("/public", "./public")
	// index
	r.StaticFile("/", "./public/index.html")

	// gen qrcode
	r.POST("/aj/genQRCode", handlergenQrCode)

	// check code
	r.POST("/aj/checkCode", handlerCheckCode)

	//r.GET("/", handlerIndex)
	r.Run() // listen and server on 0.0.0.0:8080
}
