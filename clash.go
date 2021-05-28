package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"strings"
)

type clash struct {
	Proxies []clashProxy `yaml:"proxies"`
}

type clashProxy struct {
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
	Server     string `yaml:"server"`
	Port       int    `yaml:"port"`
	Cipher     string `yaml:"cipher"`
	Password   string `yaml:"password"`
	Plugin     string `yaml:"plugin"`
	PluginOpts struct {
		Mode string `yaml:"mode"`
		Host string `yaml:"host"`
	} `yaml:"plugin-opts"`
}

func requestClashConfig(url string) (clash, error) {
	resp, err := http.Get(url)
	if nil != err {
		return clash{}, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return clash{}, err
		}
		config, err := unmarshalClashConfig(body)
		return config, err
	}

	return clash{}, errors.New("request not OK")
}

func unmarshalClashConfig(clashConfig []byte) (clash, error) {
	c := clash{}
	err := yaml.Unmarshal(clashConfig, &c)
	if nil != err {
		return c, err
	}
	return c, nil
}

func convertClashConfig(proxies []clashProxy) []config {
	var configs []config
	for _, proxy := range proxies {
		if proxy.Type != "ss" {
			continue
		}
		c := parseClashConfig(proxy)
		configs = append(configs, c)
	}
	return configs
}

func parseClashConfig(proxy clashProxy) config {
	c := config{}

	c.Server = proxy.Server
	c.ServerPort = proxy.Port
	c.Remarks = proxy.Name
	c.Method = proxy.Cipher
	c.Password = proxy.Password
	if strings.Contains(proxy.Plugin, "obfs") {
		c.Plugin = "obfs-local"

		c.PluginOpts = fmt.Sprintf("obfs=%s;obfs-host=%s", proxy.PluginOpts.Mode, proxy.PluginOpts.Host)
	}
	c.UdpFallback = c.profile
	return c
}
