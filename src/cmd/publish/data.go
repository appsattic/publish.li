// --------------------------------------------------------------------------------------------------------------------
//
// This file is part of https://github.com/appsattic/publish.li
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// --------------------------------------------------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const lenLetters = len(letterBytes)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func sendJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func sendOk(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, msg string) {
	data := struct {
		Ok  bool   `json:"ok"`
		Msg string `json:"msg"`
	}{
		Ok:  false,
		Msg: msg,
	}

	sendJson(w, data)
}
