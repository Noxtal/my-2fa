package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"
)

var db = make(map[string][2]string)

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

	ticker := time.NewTicker(time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				db = make(map[string][2]string)
				log.Printf("Current database was wiped!")
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	log.Printf("Listening on 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	users, ok := r.URL.Query()["user"]
	if !ok || len(users[0]) < 1 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "'user' parameter is missing!")
		return
	}

	user := users[0]

	passwords, ok := r.URL.Query()["password"]
	if !ok || len(passwords[0]) < 1 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "'password' parameter is missing!")
		return
	}

	password := passwords[0]
	hash := sha256.Sum256([]byte(password))

	switch r.Method {
	case "GET":
		db[user] = [2]string{string(hash[:]), Generate()}
		fmt.Fprintf(w, db[user][1])
		log.Printf("Registered user '%s' with code '%s'", user, db[user][1])
		return
	case "POST":
		codes, ok := r.URL.Query()["code"]
		if !ok || len(codes[0]) < 1 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "'code' parameter is missing!")
			return
		}
		code := codes[0]

		if c, ok := db[user]; ok && c[0] == string(hash[:]) && c[1] == code {
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
