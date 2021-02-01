package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/", parse)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func parse(w http.ResponseWriter, req *http.Request) {
	queryStr := req.URL.Query()
	urlStr := queryStr.Get("url")
	if "" == urlStr {
		http.Error(w, "Missing url", http.StatusBadRequest)
		return
	}
	sip002, err := requestSIP002(urlStr)
	if nil != err {
		http.Error(w, "Bad SIP002", http.StatusBadRequest)
		return
	}
	splitConf := strings.Split(sip002, "\n")
	var configs []config
	for _, c := range splitConf[:len(splitConf)-1] {
		conf, err := parseSIP002(c)
		if nil != err {
			http.Error(w, "Bad SIP002", http.StatusBadRequest)
			return
		}
		configs = append(configs, conf)
	}
	data, _ := json.Marshal(configs)
	_, _ = w.Write(data)
}

func requestSIP002(url string) (string, error) {
	resp, err := http.Get(url)
	if nil != err {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		bodyString := string(bodyBytes)

		plain, err := base64.StdEncoding.DecodeString(bodyString)
		if nil != err {
			return "", err
		}
		return string(plain), nil
	}

	return "", err
}

func parseSIP002(sip string) (config, error) {
	u, err := url.Parse(sip)
	if nil != err {
		return config{}, err
	}

	c := config{}

	host := u.Host
	serverAndPort := strings.Split(host, ":")
	c.Server = serverAndPort[0]
	c.ServerPort, _ = strconv.Atoi(serverAndPort[1])
	c.Remarks = u.Fragment

	userInfo, err := base64.RawURLEncoding.DecodeString(u.User.String())
	if nil != err {
		return config{}, err
	}
	methodAndPwd := strings.Split(string(userInfo), ":")
	c.Method = methodAndPwd[0]
	c.Password = methodAndPwd[1]

	query := u.Query()
	if pluginQuery := query["plugin"]; len(pluginQuery) > 0 {
		pluginConfig := pluginQuery[0]
		index := strings.Index(pluginConfig, ";")
		c.Plugin = pluginConfig[:index]
		c.PluginOpts = pluginConfig[index+1:]
	}
	return c, nil
}

type config struct {
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
	Method     string `json:"method"`
	Password   string `json:"password"`
	Plugin     string `json:"plugin,omitempty"`
	PluginOpts string `json:"plugin_opts,omitempty"`
	Remarks    string `json:"remarks,omitempty"`
}
