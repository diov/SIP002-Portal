package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	port := "8081"
	if p := os.Getenv("PORT"); "" != p {
		port = p
	}

	http.HandleFunc("/", parse)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func parse(w http.ResponseWriter, req *http.Request) {
	queryValues := req.URL.Query()
	urlStr := queryValues.Get("url")
	if "" == urlStr {
		http.Error(w, "Missing url", http.StatusBadRequest)
		return
	}

	switch queryValues.Get("type") {
	case "sip002":
		sip002, err := requestSipConfig(urlStr)
		if nil != err {
			http.Error(w, "Bad SIP002", http.StatusBadRequest)
			return
		}
		configs := convertSipConfig(sip002)
		data, _ := json.Marshal(configs)
		_, _ = w.Write(data)

	case "clash":
		clashConfig, err := requestClashConfig(urlStr)
		if nil != err {
			http.Error(w, "Bad Clash", http.StatusBadRequest)
			return
		}
		configs := convertClashConfig(clashConfig.Proxies)
		data, err := json.Marshal(configs)
		if nil != err {
			http.Error(w, err.Error(), http.StatusExpectationFailed)
		}
		_, _ = w.Write(data)
	}
}

type config struct {
	profile
	Route       string  `json:"route"`
	UdpFallback profile `json:"udp_fallback"`
}

type profile struct {
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
	Method     string `json:"method"`
	Password   string `json:"password"`
	Plugin     string `json:"plugin,omitempty"`
	PluginOpts string `json:"plugin_opts,omitempty"`
	Remarks    string `json:"remarks,omitempty"` // 名字
}
