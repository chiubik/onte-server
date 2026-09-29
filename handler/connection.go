package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"x/flap/pkg"
)

type User struct {
	Identifier string `json:"identifier"`
}

func IdentifierHandler(w http.ResponseWriter, r *http.Request) {
	user := User{}

	user.Identifier = pkg.GenerateIdentifier()
	body, err := json.Marshal(&user)
	if err != nil {
		fmt.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
