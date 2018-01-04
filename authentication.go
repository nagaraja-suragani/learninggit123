package doctorservice

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore" 
)

const (
	privKeyPath = "keys/app.rsa"
	pubKeyPath  = "keys/app.rsa.pub"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)
var verifyBytes, signBytes []byte

func initKeys() {
	var err error

	signBytes, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatalf("Error reading private key: %v", err)
	}
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatalf("Error parsing private key: %v", err)
	}
	verifyBytes, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatalf("Error reading public key: %v", err)
	}
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatalf("Error parsing public key: %v", err)
	}
}

// Response is strcut
type Response struct {
	Data string `json:"data"`
}

// Token is also strcut
type Token struct {
	Token string `json:"token"`
}

// ProtectedHandler is handler
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {

	response := Response{"Gained access to protected resource"}
	JSONResponse(response, w)

}

//RegisterHandler is using to insert user into database
func RegisterHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		panic(err)
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	log.Printf("%#v Getting hashedPass ", hashedPass)
	if err != nil {
		panic(err)
	}

	keys := make([]*datastore.Key, 1)
	ctx := appengine.NewContext(r)
	ctx, err = appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		panic(err)
	}
	user.Password = string(hashedPass)
	log.Printf("%#v Getting user.Password ", user.Password)
	user.Username = strings.ToLower(user.Username)
	log.Printf("%#v Getting kk ", user.Username)
	keys[0] = datastore.NewKey(ctx, "User", user.Username, 0, nil)

	_, err = datastore.Put(ctx, keys[0], user)
	if err != nil {
		panic(err)
	}

	//user.ID = k.IntID()
	json.NewEncoder(w).Encode(user)
}

//LoginHandler is using to check wherther user existed or not
//If User Existed In the datastore send tocken back
func LoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		panic(err)
	}

	ctx := appengine.NewContext(r)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		panic(err)
	}

	// get one
	userdatastore := &User{}
	user.Username = strings.ToLower(user.Username)
	key := datastore.NewKey(ctx, "User", user.Username, 0, nil)
	err = datastore.Get(ctx, key, userdatastore)

	//if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.FormValue("password"))) != nil {
	if err != nil || bcrypt.CompareHashAndPassword([]byte(userdatastore.Password), []byte(user.Password)) != nil {
		// there is an err, there is a NO user
		//fmt.Fprint(w, "false")
		w.WriteHeader(http.StatusForbidden)
		fmt.Println("Error logging in")
		fmt.Fprint(w, "Invalid credentials")
		return
	}

	//create a rsa 256 signer
	signer := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "admin",
		//"exp": time.Now().Add(time.Minute * 20).Unix(),
		"CustomUserInfo": struct {
			Name string
			Role string
		}{user.Username, "Member"}})

	//set claims

	tokenString, err := signer.SignedString(signKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		log.Printf("Error signing token: %v\n", err)
	}
	//create a token instance using the token string
	response := Token{tokenString}
	cookie := &http.Cookie{
		Name:  "session",
		Value: tokenString,
		Path:  "/",
		//		UNCOMMENT WHEN DEPLOYED:
		//		Secure: true,
		//		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	JSONResponse(response, w)

}

func Forgotpassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		panic(err)
	}

	ctx := appengine.NewContext(r)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		panic(err)
	}
	userdata := &User{}
	user.Username = strings.ToLower(user.Username)
	key := datastore.NewKey(ctx, "User", user.Username, 0, nil)
	err = datastore.Get(ctx, key, userdata)

	log.Printf("%#v Getting user.Password ", user.Username)

	if err != nil {
		// there is an err, there is a NO user
		fmt.Fprint(w, " there is a NO user for these username")
		return
	} else {

		hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPass)
		userdata.Password = user.Password
		_, err = datastore.Put(ctx, key, userdata)
		if err != nil {
			panic(err)
		}

	}

}

//ValidateTokenMiddleware is a AUTH TOKEN VALIDATION
func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	//validate token
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if err == nil {

		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorised access to this resource")
	}

}

//JSONResponse is a HELPER FUNCTIONS
func JSONResponse(response interface{}, w http.ResponseWriter) {

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func startServer() {
	/*http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/forgot", Forgotpassword)*/

	//PROTECTED ENDPOINTS
	http.Handle("/resource/", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(ProtectedHandler)),
	))

	log.Println("Now listening...")

}

/*
func init() {
	initKeys()
	startServer()
}*/
