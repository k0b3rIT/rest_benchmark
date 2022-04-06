package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bradhe/stopwatch"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

type ApiTest struct {
	Api    string
	Params []map[string]string
}

func main() {

	baseUrlPtr := flag.String("host", "", "host te execute the query against")
	configPtr := flag.String("config", "", "config file path")

	flag.Parse()

	baseUrl := *baseUrlPtr
	configFilePath := *configPtr

	if baseUrl == "" {
		panic("you have to define the host")
	}

	if configFilePath == "" {
		panic("you have to define the config file")
	}

	configStr, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Print(err)
	}

	var apiTests []ApiTest

	err = json.Unmarshal([]byte(configStr), &apiTests)

	if err != nil {
		panic(fmt.Sprint("Unable to parse config file, %v", err))
	}

	fmt.Printf("Start banchmark on %s\n", baseUrl)

	totalWatch := stopwatch.Start()

	for _, apiTest := range apiTests {
		executeApiTest(baseUrl, apiTest.Api, apiTest.Params)
	}

	totalWatch.Stop()
	fmt.Printf(inRed("Total: %v\n"), totalWatch.String())
}

func executeApiTest(baseUrl string, api string, params []map[string]string) {

	for _, param := range params {
		apiWithParams := substituteParams(api, param)

		watch := stopwatch.Start()
		requestUrl := fmt.Sprintf("%s%s", baseUrl, apiWithParams)
		fmt.Printf("%s", apiWithParams)
		_, err := http.Get(requestUrl)
		if err != nil {
			panic(fmt.Sprint("Failed to send request to [%s]", requestUrl))
		}

		watch.Stop()
		dur := watch.Milliseconds() * 1000 * 1000
		fmt.Printf("\t%v\n", dur.Milliseconds())
	}

}

func substituteParams(api string, params map[string]string) string {
	afterReplace := api
	for k, v := range params {
		afterReplace = strings.Replace(afterReplace, "{"+k+"}", v, 1)
	}
	return afterReplace
}

func inRed(str string) string {
	return Red + str + Reset
}
