package models

import (
	"fmt"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"gopkg.in/mgo.v2"
	u "github.com/form-with-jwt/utils"
	"os"
	"golang.org/x/crypto/bcrypt"
	"time"

)

/*
JWT claims struct
*/
type Token struct {
	Id uint64
	jwt.StandardClaims
}

//a struct to rep user account
type Account struct {
	ID     string   `json:"id" bson:"id"`
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Token string `json:"token" bson:"token"`
	ProfileImage string `json:"profileImage" bson:"profileImage"`
}

//Validate incoming user details...
func (account *Account) Validate() (map[string] interface{}, bool) {

	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Invalid Email address"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password is required with atleast 6 characters"), false
	}

	//Email must be unique
	temp := &Account{}

	//check for errors and duplicate emails
	db, err := mgo.Dial("localhost")
    if err != nil {
		//log.Fatal("cannot dial mongo", err)
		fmt.Println("cannot dial mongo", err)
		return nil, false
    }
    defer db.Close()
	err = db.DB(os.Getenv("db_name")).C("users1").Find(bson.M{"email": account.Email}).Select(bson.M{"username": "", "email": "", "token": ""}).One(&temp)
	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (account *Account) Create() (map[string] interface{}) {

	db, err := mgo.Dial("localhost")
    if err != nil {
		//log.Fatal("cannot dial mongo", err)
		fmt.Println("cannot dial mongo", err)
		return nil
    }
    defer db.Close()

	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	
	//idd,_ := strconv.ParseUint(account.ID, 10, 64);
	//tk := &Token{Id: idd}
	//token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	//fmt.Println(tk)
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	claims := make(jwt.MapClaims)

	claims["foo"] = "bar"
	claims["email"] = account.Email
	claims["password"] = account.Password
	claims["id"] = account.ID
    claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
    // Sign and get the complete encoded token as a string
    tokenString, err := token.SignedString("XKDIEP1323s")
    //return tokenString, err
	//tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "" //delete password
	account.ID = u.MD5Hash(account.Email, account.Password)

	if err := db.DB(os.Getenv("db_name")).C("users1").Insert(&account); err != nil {//why
		fmt.Println("Signup unsuccessful.")
        return nil
    } else {
		fmt.Println("Signup successful.")
    }

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

func Login(email, password string) (map[string]interface{}) {

	accountPayload := &Account{}

	db, err := mgo.Dial("localhost")
    if err != nil {
		fmt.Println("cannot dial mongo", err)
		return nil
    }
	defer db.Close()
	
    err = db.DB(os.Getenv("db_name")).C("users1").Find(bson.M{"email": email}).Select(bson.M{"username": "", "email": "", "token": "", "id":""}).One(&accountPayload)
    if err != nil {
		resp:= u.Message(false, "No Such User")
		resp["account"] = nil
		return resp
	}
	

	err = bcrypt.CompareHashAndPassword([]byte(accountPayload.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	accountPayload.Password = ""
	//Worked! Logged In
	idd,_ := strconv.ParseUint(accountPayload.ID, 0, 64)// converting unit64 to base of string as 0 or till 36 only
	tk := &Token{Id: idd}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	accountPayload.Token = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In")// used for claims 
	resp["account"] = accountPayload
	fmt.Println("login successful.")
	return resp
}

func GetUser(email, password string) *Account {
	db, err := mgo.Dial("localhost")
    if err != nil {
		//log.Fatal("cannot dial mongo", err)
		fmt.Println("cannot dial mongo", err)
		return nil
    }
    defer db.Close()
	acc := &Account{}
    err = db.DB("db_name").C("users1").Find(bson.M{"email": email,"password":password}).Select(bson.M{"username": "", "email": "", "token": "", "profileImage": ""}).One(&acc)
    if err != nil {
        panic(err)
	}
	return acc
	
}