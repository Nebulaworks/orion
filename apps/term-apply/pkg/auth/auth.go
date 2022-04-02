package auth

import (
	"bufio"
	"fmt"
	"log"
	"net/http"

	gossh "golang.org/x/crypto/ssh"

	"github.com/gliderlabs/ssh"
)

func compareKeys(username string, key ssh.PublicKey) bool {
	tryfp := gossh.FingerprintSHA256(key)
	log.Printf("username: %s attempting to auth with %s", username, tryfp)

	// Should probably cache this since each auth type will invoke the service
	// reaching out to github
	r, err := http.Get(fmt.Sprintf("https://github.com/%s.keys", username))

	if err != nil {
		log.Println(err)
		return false
	}

	defer r.Body.Close()

	scanner := bufio.NewScanner(r.Body)

	for scanner.Scan() {
		linea := scanner.Text()
		line, _, _, _, err := gossh.ParseAuthorizedKey([]byte(linea))
		if err != nil {
			log.Println(err)
			return false
		}
		fp := gossh.FingerprintSHA256(line)
		if fp == tryfp {
			log.Printf("username: %s found match: %s", username, fp)
			return true
		}
		log.Printf("username: %s not a match: %s", username, fp)
	}
	if scanner.Err() != nil {
		log.Println(scanner.Err())
	}
	fmt.Printf("username: %s No match\n", username)
	return false
}

func PkHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	return compareKeys(ctx.User(), key)
}
