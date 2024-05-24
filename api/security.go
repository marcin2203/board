package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	_ "encoding/json"
	"fmt"

	"github.com/cristalhq/jwt/v5"
)

type Claims struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

func encryptPasswordSHA256(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString([]byte(h[:]))
}

func decriptJWT(token string) Claims {
	var myClaims Claims
	jwtToken, _ := jwt.ParseNoVerify([]byte(token))
	json.Unmarshal(jwtToken.Claims(), &myClaims)
	return myClaims
}
func verifyJWT(token string) (Claims, error) {
	var myClaims Claims
	verifier, err := jwt.NewVerifierHS(jwt.HS256, []byte("secret"))
	if err != nil {
		fmt.Println(err)
	}
	err = jwt.ParseClaims([]byte(token), verifier, &myClaims)
	return myClaims, err
}
func getJWTFrom(email string, nickname string) *jwt.Token {
	key := []byte("secret")
	signer, _ := jwt.NewSignerHS(jwt.HS256, key)
	builder := jwt.NewBuilder(signer)

	claims := &Claims{
		Email:    email,
		Nickname: nickname,
	}

	token, _ := builder.Build(claims)
	return token
}

// func main() {
// 	// CREATE JWT

// 	key := []byte("secret")
// 	signer, _ := jwt.NewSignerHS(jwt.HS256, key)
// 	builder := jwt.NewBuilder(signer)

// 	claims := &Claims{
// 		Email:    "12@a.com",
// 		Nickname: "12",
// 	}

// 	token, _ := builder.Build(claims)
// 	fmt.Println(token)

// 	// VERIFY
// 	myjwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6IjEyMzRAZ21haWwuY29tIiwibmlja25hbWUiOiIxMjM0In0.fvHtYcqXvxVhx3VJiJfCLFX9p_yGnG5EfEG_p2k0syc"
// 	VerifyJWT, err := verifyJWT(myjwt)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println("VerifyJWT: ", VerifyJWT)
// 	DecriptJWT := decriptJWT(myjwt)
// 	VerifyJWT, err = verifyJWT(token.String())

// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println("DecriptJWT: ", DecriptJWT)
// 	fmt.Println("2nd VerifyJWT: ", err, DecriptJWT)

// 	test := getJWTFrom("adam", "rolo")
// 	fmt.Println(test.String())
// 	VerifyJWT, err = verifyJWT(test.String())
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println("Verify GetJWTFrom(): ", VerifyJWT)
// 	//GET CLAIMS
// 	var retrivedClaims Claims
// 	json.Unmarshal(token.Claims(), &retrivedClaims)
// 	fmt.Println(retrivedClaims)

// 	mypass := encryptPasswordSHA256("1234")

// 	sql := "03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4"
// 	fmt.Println(mypass, sql)
// 	expr := (sql == mypass)

// 	fmt.Println(expr)
// }
