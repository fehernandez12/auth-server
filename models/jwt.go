package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type JSONSerializable interface {
	ToJSON() ([]byte, error)
}

// Jwt is a struct that represents a JSON Web Token.
// A JSON Web Token (JWT) is a compact, URL-safe means of representing claims to be transferred between two parties.
// The claims in a JWT are encoded as a JSON object that is used as the payload of a JSON Web Signature (JWS) structure or as the plaintext of a JSON Web Encryption (JWE) structure, enabling the claims to be digitally signed or integrity protected with a Message Authentication Code (MAC) and/or encrypted.
// Source: https://tools.ietf.org/html/rfc7519
//
// A JWT is composed of three parts: a header, a payload, and a signature.
// The header contains information about how the JWT signature should be computed.
// The payload contains the claims of the token.
// The signature is used to verify the message wasn't changed along the way, and, in the case of tokens signed with a private key, it can also verify that the sender of the JWT is who it says it is.
// Source: https://jwt.io/introduction/
type Jwt struct {
	Header    *Header  `json:"header"`
	Payload   *Payload `json:"payload"`
	message   string
	Signature string `json:"signature"`
}

// Header is a struct that represents the header of a JSON Web Token.
// The header contains information about how the JWT signature should be computed.
// The header contains the following fields:
// alg (algorithm): Algorithm used to sign the JWT
// typ (type): Type of the token
type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// Payload is a struct that represents the payload of a JSON Web Token.
// The payload contains the claims of the token.
// The claims are statements about an entity (typically, the user) and additional data.
// The payload contains the following claims:
// iss (issuer): Issuer of the JWT
// sub (subject): Subject of the JWT (the user)
// aud (audience): Recipient for which the JWT is intended
// exp (expiration time): Time after which the JWT expires
// iat (issued at time): Time at which the JWT was issued
// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)
type Payload struct {
	Iss   string   `json:"iss"`
	Sub   string   `json:"sub"`
	Aud   string   `json:"aud"`
	Exp   int64    `json:"exp"`
	Iat   int64    `json:"iat"`
	Jti   string   `json:"jti"`
	Scope []string `json:"scope"`
}

// NewJwt is a function that creates a new Jwt.
func NewJwt(payload *Payload, typ string) (*Jwt, error) {
	var jwt Jwt
	alg := os.Getenv("AUTH_SERVER_JWT_ALG")
	jwt.Header = newHeader(alg, typ)
	jwt.Payload = payload
	jwt.makeMessage()
	err := jwt.makeSignature(alg)
	if err != nil {
		return nil, err
	}
	return &jwt, nil
}

// ToJSON is a method that converts a Jwt to JSON.
func (j *Jwt) ToJSON() ([]byte, error) {
	return json.Marshal(j)
}

func (j *Jwt) Token() (string, error) {
	return j.message + "." + j.Signature, nil
}

func (j *Jwt) makeMessage() error {
	header, err := j.Header.ToJSON()
	if err != nil {
		return err
	}
	headerBase64 := base64.RawURLEncoding.EncodeToString(header)
	payload, err := j.Payload.ToJSON()
	if err != nil {
		return err
	}
	payloadBase64 := base64.RawURLEncoding.EncodeToString(payload)
	j.message = headerBase64 + "." + payloadBase64
	return nil
}

func (j *Jwt) makeSignature(alg string) error {
	switch alg {
	case "HS256":
		j.makeHS256Signature()
	default:
		return errors.New("invalid algorithm")
	}
	return nil
}

func (j *Jwt) makeHS256Signature() {
	secret := os.Getenv("AUTH_SERVER_JWT_SECRET")
	log.Printf("secret: %s", secret)
	hasher := hmac.New(sha256.New, []byte(secret))
	hasher.Write([]byte(j.message))
	signature := hasher.Sum(nil)
	j.Signature = base64.RawURLEncoding.EncodeToString(signature)
}

func (j *Jwt) CheckSignature(signature string) (bool, error) {
	switch j.Header.Alg {
	case "HS256":
		return j.checkHS256Signature(signature)
	default:
		return false, errors.New("invalid algorithm")
	}
}

func (j *Jwt) checkHS256Signature(signature string) (bool, error) {
	if j.Signature != signature {
		return false, nil
	}
	return true, nil
}

// NewHeader is a function that creates a new Header.
func newHeader(alg string, typ string) *Header {
	return &Header{
		Alg: alg,
		Typ: typ,
	}
}

// ToJSON is a method that converts a Header to JSON.
func (h *Header) ToJSON() ([]byte, error) {
	return json.Marshal(h)
}

// NewPayload is a function that creates a new Payload.
func NewPayload(sub string, aud string, exp int, scope string) *Payload {
	jti := uuid.New().String()
	return &Payload{
		Iss:   os.Getenv("AUTH_SERVER_JWT_ISS"),
		Sub:   sub,
		Aud:   aud,
		Exp:   time.Now().Add(time.Hour * time.Duration(exp)).Unix(),
		Iat:   time.Now().Unix(),
		Jti:   jti,
		Scope: strings.Split(scope, " "),
	}
}

// ToJSON is a method that converts a Payload to JSON.
func (p *Payload) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}
