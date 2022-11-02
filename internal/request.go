package internal

import (
	"net/http"
)

type ChainElem struct {
	Command string   `json:"command"`
	Params  []string `json:"params"`
}

type Request struct {
	Chain []ChainElem `json:"chain"`
}

type Response struct {
	Data   string `json:"data"`
	Custom string `json:"custom"`
	Error  string `json:"error"`
}

func Search(w http.ResponseWriter, req *http.Request) {
	//decoder := json.NewDecoder(req.Body)
	//var r Request
	//err := decoder.Decode(&r)
	//if err != nil {
	//	panic(err)
	//}
	//
	//errorMessage := ""
	//
	//dp, err := RunPipeline(r, "debug")
	//if err != nil {
	//	errorMessage = err.Error()
	//}
	//js := Response{
	//	Data:   data,
	//	Custom: custom,
	//	Error:  errorMessage,
	//}
	//jstr, _ := json.Marshal(js)
	//
	//fmt.Fprint(w, string(jstr))

}
