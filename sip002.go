package main

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func requestSipConfig(url string) (string, error) {
	resp, err := http.Get(url)
	if nil != err {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		bodyString := string(body)

		plain, err := base64.StdEncoding.DecodeString(bodyString)
		if nil != err {
			return "", err
		}
		return string(plain), nil
	}

	return "", err
}

func convertSipConfig(sipConfig string) []config {
	splitConf := strings.Split(sipConfig, "\n")
	var configs []config
	for _, c := range splitConf[:len(splitConf)-1] {
		conf, err := parseSipConfig(c)
		if nil != err {
			continue
		}
		configs = append(configs, conf)
	}
	return configs
}

func parseSipConfig(sip string) (config, error) {
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
