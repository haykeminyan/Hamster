package jwttoken

import (
	"os"
	"reflect"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/matscus/Hamster/Package/JWTToken/subset"
)

//Token - struct for token
type Token struct {
	Token string `json:"token"`
}

func (t Token) New() subset.Token {
	var token subset.Token
	token = Token{}
	return token
}

//Generate - genereta default jwt
func (t Token) Generate(role string, user string, projects []string) (tokenstring string, err error) {
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), jwt.MapClaims{
		"user":    user,
		"role":    role,
		"project": projects,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	})
	tokenstring, err = token.SignedString([]byte(os.Getenv("KEY")))
	if err != nil {
		return tokenstring, err
	}
	return tokenstring, err
}

//GenerateTemp - genereta default jwt
func (t Token) GenerateTemp(role string, user string, projects []string) (tokenstring string, err error) {
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), jwt.MapClaims{
		"user":    user,
		"role":    role,
		"project": projects,
		"exp":     time.Now().Add(time.Minute * 1).Unix(),
	})
	tokenstring, err = token.SignedString([]byte(os.Getenv("KEY")))
	if err != nil {
		return tokenstring, err
	}
	return tokenstring, err
}

//Check  - func for check validate JWT, result bool
func (t Token) Check() bool {
	token, err := jwt.Parse(t.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("KEY")), nil
	})
	if err == nil && token.Valid {
		return true
	}
	return false
}

//Parse - func to parse token and check to valid
func Parse(t string) bool {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("KEY")), nil
	})
	if err == nil && token.Valid {
		return true
	}
	return false
}

//IsAdmin - func to parse token and check to valid
func IsAdmin(t string) bool {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("KEY")), nil
	})
	if err == nil && token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		role := claims["role"]
		if role == "admin" {
			return true
		} else {
			return false
		}
	}
	return false
}

//GetUser - func to parse token and check to valid
func GetUser(t string) string {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("KEY")), nil
	})
	if err == nil && token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		role := claims["user"]
		return role.(string)
	}
	return ""
}

//GetUserProjects - func to parse token and check to valid
func GetUserProjects(t string) []string {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("KEY")), nil
	})
	if err == nil && token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		role := claims["project"]
		value := reflect.ValueOf(role)
		len := value.Len()
		res := make([]string, 0, len)
		for i := 0; i < value.Len(); i++ {
			temp := value.Index(i).Interface()
			res = append(res, temp.(string))
		}
		return res
	}
	return nil
}
