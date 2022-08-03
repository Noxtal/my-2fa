package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
)

var db = make(map[string]string)

const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const length = int64(len(charset))

func Generate() string {
	b := make([]byte, 6)
	for i := 0; i < 6; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(length))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

func main() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	users, ok := r.URL.Query()["user"]
	if !ok || len(users[0]) < 1 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "'user' parameter is missing!")
		return
	}

	user := users[0]

	switch r.Method {
	case "GET":
		db[user] = Generate()
		fmt.Fprintf(w, db[user])
		log.Printf("Registered user '%s' with code '%s'", user, db[user])
		return
	case "POST":
		codes, ok := r.URL.Query()["code"]
		if !ok || len(codes[0]) < 1 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "'code' parameter is missing!")
			return
		}
		code := codes[0]

		if c, ok := db[user]; ok && c == code {
			delete(db, user)
			fmt.Fprintf(w, "Access Granted (200)")
			log.Printf("Access granted to user '%s' with code '%s'", user, code)
			return
		} else {
			w.WriteHeader(http.StatusForbidden)
			log.Printf("Access denied to user '%s' with code '%s'", user, code)
			fmt.Fprintf(w, "Unauthorized (401)")
			return
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error (500)")
		return
	}
}
