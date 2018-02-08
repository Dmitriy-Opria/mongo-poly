package jwt

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"golang.org/x/crypto/scrypt"
	"strings"
	"time"
)

type (
	ID = uint32
)

var (
	jwtSecret = []byte{}
	sign      = HS256
)

type JwtSession struct {
	Perm    uint32 `json:"perm,omitempty"`
	Group   uint32 `json:"group,omitempty"`
	UID     uint32 `json:"uid,omitempty"`
	Expired int64  `json:"exp,omitempty"`
	Email   string `json:"email,omitempty"`
	IP      string `json:"ip,omitempty"`
}

func HS256(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

func HS384(message, key []byte) []byte {
	mac := hmac.New(sha512.New384, key)
	mac.Write(message)
	return mac.Sum(nil)
}

func SetJWTSecret(secret []byte) {
	if len(secret) < 32 {
		panic("len(jwtSecret) < 32")
	}
	var jwtSalt = []byte(`V1fiCcYjH;0t}h4(Vpo7"bn1$^fw.0`)
	data, err := scrypt.Key(secret, jwtSalt, 16384, 8, 1, 64)
	if err != nil {
		panic(err)
	}
	jwtSecret = data
}

func (s *JwtSession) JWTString() string {
	jsonBytes := bytes.Buffer{}
	json.NewEncoder(&jsonBytes).Encode(s)
	if jsonBytes.Len() > 0 {
		jsonBytes.Truncate(jsonBytes.Len() - 1)
	}
	base64Bytes := bytes.Buffer{}
	encoderBase64 := base64.NewEncoder(base64.RawURLEncoding, &base64Bytes)
	encoderBase64.Write(jsonBytes.Bytes())
	encoderBase64.Close()
	data := sign(base64Bytes.Bytes(), jwtSecret)
	base64Bytes.WriteString(".")
	encoderBase64 = base64.NewEncoder(base64.RawURLEncoding, &base64Bytes)
	encoderBase64.Write(data)
	encoderBase64.Close()
	return base64Bytes.String()
}

func (s *JwtSession) String() string {
	jsonBytes := bytes.Buffer{}
	json.NewEncoder(&jsonBytes).Encode(s)
	if jsonBytes.Len() > 0 {
		jsonBytes.Truncate(jsonBytes.Len() - 1)
	}
	return jsonBytes.String()
}

func Decode(str string) (*JwtSession, bool) {
	if index := strings.Index(str, "."); index >= 0 {
		jsonBase64, hmacBase64 := str[:index], str[index+1:]
		if signBytes, err := base64.RawURLEncoding.DecodeString(hmacBase64); err == nil {
			if testBytes := sign([]byte(jsonBase64), jwtSecret); bytes.Equal(signBytes, testBytes) {
				base64Reader := base64.NewDecoder(base64.RawURLEncoding, strings.NewReader(jsonBase64))
				session := new(JwtSession)
				if err = json.NewDecoder(base64Reader).Decode(session); err == nil {
					return session, session.Expired > time.Now().Unix()
				}
			}
		}
	}
	return nil, false
}
