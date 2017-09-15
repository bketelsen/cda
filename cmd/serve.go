// +build !darwin
// +build !windows

// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
)

var (
	dbPath string
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the CDA URL Shortening service",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func serve() {

	if os.Getenv("BASE_URL") == "" {
		baseURL = "cda.ms"
	}
	if os.Getenv("DB_PATH") == "" {
		var err error
		dbPath, err = os.Getwd()
		if err != nil {
			log.Fatal("Unable to get executable path", err)
		}
	}
	fmt.Println(path.Join(dbPath, "db.sqlite"))
	db := sqlite{Path: path.Join(dbPath, "db.sqlite")}
	db.Init()

	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	r := mux.NewRouter()
	r.HandleFunc("/save",
		func(response http.ResponseWriter, request *http.Request) {
			encodeHandler(response, request, db, baseURL)
		}).Methods("POST")
	r.HandleFunc("/{id}", func(response http.ResponseWriter, request *http.Request) {
		decodeHandler(response, request, db)
	})
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	h := cors.Default().Handler(r)
	log.Println("Starting server on port :1337")
	log.Fatal(http.ListenAndServe(":1337", handlers.LoggingHandler(os.Stdout, h)))
}

const base string = "0123456789bcdfghjkmnpqrstvwxyzBCDFGHJKLMNPQRSTVWXYZ"

func encode(n int64) string {
	var s string
	var num = float64(n)

	for num > 0 {
		s = string(base[int(num)%len(base)]) + s
		num = math.Floor(num / float64(len(base)))
	}

	return s
}

func decode(s string) int {
	var num = 0
	for _, element := range strings.Split(s, "") {
		num = num*len(base) + strings.Index(base, element)
	}

	return num
}

func decodeHandler(response http.ResponseWriter, request *http.Request, db Database) {
	id := decode(mux.Vars(request)["id"])
	url, err := db.Get(id)
	if err != nil {
		http.Error(response, `{"error": "No such URL"}`, http.StatusNotFound)
		return
	}
	http.Redirect(response, request, url, 301)
}

func encodeHandler(response http.ResponseWriter, request *http.Request, db Database, baseURL string) {
	decoder := json.NewDecoder(request.Body)
	var data struct {
		URL string `json:"url"`
	}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(response, `{"error": "Unable to parse json"}`, http.StatusBadRequest)
		return
	}

	if !govalidator.IsURL(data.URL) {
		http.Error(response, `{"error": "Not a valid URL"}`, http.StatusBadRequest)
		return
	}

	id, err := db.Save(data.URL)
	if err != nil {
		log.Println(err)
		return
	}

	resp := map[string]string{"url": baseURL + encode(id), "id": encode(id), "error": ""}
	jsonData, _ := json.Marshal(resp)
	response.Write(jsonData)

}
