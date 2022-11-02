package rest

import (
	"encoding/json"
	"fmt"
	"github.com/78planet/nomadcoin/blockchain"
	"github.com/78planet/nomadcoin/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type url string

func (u url) MarshalText() (text []byte, err error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL         url    `json:"url,omitempty"`
	Method      string `json:"method,omitempty"`
	Description string `json:"description,omitempty"`
	Payload     string `json:"payload,omitempty"`
}

type addBlockBody struct {
	Message string
}

type errorResponse struct {
	ErrorMessage string `json:"error_message,omitempty"`
}

func (u *urlDescription) String() {
	fmt.Printf("url: %v, Method: %s, Description: %s, Payload: %s", u.URL, u.Method, u.Description, u.Payload)
}

var port string = ":4000"

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add a block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{height}"),
			Method:      "GET",
			Description: "See a block",
		},
	}
	rw.Header().Add("Content-Type", "application-json")
	utils.HandleErr(json.NewEncoder(rw).Encode(data))
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rw.Header().Add("Content-Type", "application-json")
		err := json.NewEncoder(rw).Encode(blockchain.GetBlockchain().AllBlocks())
		utils.HandleErr(err)
	case "POST":
		var addBlockBody addBlockBody
		err := json.NewDecoder(r.Body).Decode(&addBlockBody)
		utils.HandleErr(err)
		blockchain.GetBlockchain().AddBlock(addBlockBody.Message)
		rw.WriteHeader(http.StatusCreated)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application-json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["height"])
	utils.HandleErr(err)
	block, err := blockchain.GetBlockchain().GetBlock(id)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		utils.HandleErr(encoder.Encode(errorResponse{fmt.Sprint(err)}))
	}
	utils.HandleErr(encoder.Encode(block))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application-json")
		next.ServeHTTP(rw, r)
	})
}

func Start(aPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", aPort)
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/block/{height:[0-9]+}", block).Methods("GET")
	fmt.Printf("Listening on http://localhost%s \n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
