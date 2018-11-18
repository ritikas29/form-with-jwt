package utils

import (
	"encoding/json"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"net/http"
)

func Message(status bool, message string) (map[string]interface{}) {
	return map[string]interface{} {"status" : status, "message" : message}
}

func Respond(w http.ResponseWriter, data map[string] interface{})  {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func MD5Hash(email string, password string) string {
	h := md5.New()
	strHash := email + password
    h.Write([]byte(strings.ToLower(strHash)))
    return hex.EncodeToString(h.Sum(nil))
 }