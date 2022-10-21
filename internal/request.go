package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Request struct {
	Article string `json:"article"`
}

type Response struct {
	Data  string `json:"data"`
	Error string `json:"error"`
}

func Search(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var r Request
	err := decoder.Decode(&r)
	if err != nil {
		panic(err)
	}

	errorMessage := ""
	resp, err := RunPipeline(r.Article, "catalog", "debug")
	if err != nil {
		errorMessage = err.Error()
	}
	js := Response{
		Data:  resp,
		Error: errorMessage,
	}
	jstr, _ := json.Marshal(js)

	fmt.Fprint(w, jstr)

}
