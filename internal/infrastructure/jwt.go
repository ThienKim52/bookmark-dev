package infrastructure

import (
	"github.com/ThienKim52/golang-dev/pkg/jwtutils"
)

func CreateJWTProvider() (jwtutils.JWTGenerator, jwtutils.JWTValidator) {
	jwtGenerator, err := jwtutils.NewJWTGenerator("./private.pem")
	if err != nil {
		panic(err)
	}
	jwtValidator, err := jwtutils.NewJWTValidator("./public.pem")
	if err != nil {
		panic(err)
	}
	return jwtGenerator, jwtValidator
}