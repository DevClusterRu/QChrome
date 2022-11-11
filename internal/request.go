package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type ChainElem struct {
	Command string   `json:"command"`
	Params  []string `json:"params"`
}

type Request struct {
	Token string      `json:"token"`
	Chain []ChainElem `json:"chain"`
}

type Response struct {
	Token  string              `json:"token"`
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

	var br *Instance
	var ok bool

	if r.Token == "" {
		br, err = MakeBrowser("debug")
		if err != nil {
			log.Fatal(err)
		}
		Browsers[br.Token] = br
		go func(t string) {
			//Браузер живет 120 сек и умирает
			time.Sleep(5 * time.Minute)
			br.Close()
			delete(Browsers, t)
			log.Println("Browser is deleted")
		}(r.Token)
	} else {
		fmt.Println(Browsers)

		if br, ok = Browsers[r.Token]; !ok {
			fmt.Fprint(w, "Browser time is out")
			return
		}
	}

	err = br.RunPipeline(r)
	if err != nil {
		log.Fatal(err)
	}

	js := Response{
		Data:   br.Data,
		Custom: br.CustomData,
		Error:  errorMessage,
		Token:  br.Token,
	}
	jstr, _ := json.Marshal(js)

	fmt.Fprint(w, string(jstr))

}
