package main

import (
	"testing"
)

func TestClashRequest(t *testing.T) {
	url := "https://renzhesub.com/link/juuF3eC9v7oyxZqk?list=clash"
	sipConfig, err := requestClashConfig(url)
	if nil != err {
		t.Error(err)
	}
	t.Log(sipConfig)
	TestParseClashConfig(t)
}

func TestParseClashConfig(t *testing.T) {
	url := "https://renzhesub.com/link/juuF3eC9v7oyxZqk?list=clash"
	config, err := requestClashConfig(url)
	if nil != err {
		t.Error(err)
	}
	t.Log(config)
	configs := convertClashConfig(config.Proxies)
	t.Log(configs)
}
