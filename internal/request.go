package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type ChainElem struct {
	Command string   `json:"command"`
	Params  []string `json:"params"`
}

type Request struct {
	Chain []ChainElem `json:"chain"`
}

type Response struct {
	Data   []map[string]string `json:"data"`
	Custom map[string]string   `json:"custom"`
	Error  string              `json:"error"`
}

func GetImage(w http.ResponseWriter, req *http.Request) {
	b, err := os.ReadFile(strings.ReplaceAll(req.URL.String(), "/", ""))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(b)

}

func Search(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	var r Request
	err := decoder.Decode(&r)
	if err != nil {
		panic(err)
	}

	errorMessage := ""

	br, err := MakeBrowser("debug")
	if err != nil {
		log.Fatal(err)
	}
	defer br.Close()

	err = br.RunPipeline(r)
	if err != nil {
		log.Fatal(err)
	}

	js := Response{
		Data:   br.Data,
		Custom: br.CustomData,
		Error:  errorMessage,
	}
	jstr, _ := json.Marshal(js)

	fmt.Fprint(w, string(jstr))

}
