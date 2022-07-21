package jwt

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type jwt struct {
	key string
}

type Body struct {
	UserId int `json:"userId"`
}

func JWT(key string) *jwt {
	return &jwt{
		key: key,
	}
}

func (k *jwt) IsValid(token string) bool {
	tSplited := strings.Split(token, ".")
	if len(tSplited) == 2 {
		body, _ := base64.StdEncoding.DecodeString(tSplited[0])
		bodyStr := k.encryptBody(body)
		if bodyStr == tSplited[1] {
			return true
		}
	}

	return false
}

func (k *jwt) GetToken(d Body) string {
	dataByte, _ := json.Marshal(d)
	signedByte := k.encryptBody(dataByte)
	dataStr := base64.StdEncoding.EncodeToString(dataByte)

	return fmt.Sprintf("%s.%s", dataStr, signedByte)
}

func (k *jwt) encryptBody(body []byte) string {
	b := string(body) + k.key
	return fmt.Sprintf("%x", sha256.Sum256([]byte(b)))
}

func (k *jwt) GetBody(token string) *Body {
	tSplited := strings.Split(token, ".")
	encryptBody, _ := base64.StdEncoding.DecodeString(tSplited[0])
	d := Body{}
	json.Unmarshal(encryptBody, &d)

	return &d
}
