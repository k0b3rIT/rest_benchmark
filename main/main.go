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

type ArangoTest struct {
	Query  string
	Params []map[string]string
}

func main() {

	baseUrlPtr := flag.String("host", "", "host te execute the query against")
	configPtr := flag.String("config", "", "config file path")
	benchmarkTypePtr := flag.String("type", "", "benchmark type")

	flag.Parse()

	baseUrl := *baseUrlPtr
	configFilePath := *configPtr
	benchmarkType := *benchmarkTypePtr

	if baseUrl == "" {
		panic("you have to define the host")
	}

	if configFilePath == "" {
		panic("you have to define the config file")
	}

	if benchmarkType != "api" && benchmarkType != "arango" {
		panic("invalid type")
	}

	if benchmarkType == "api" {
		apiBenchmark(baseUrl, configFilePath)
	}

	if benchmarkType == "arango" {
		arangoBenchmark(baseUrl, configFilePath)
	}

}

func arangoBenchmark(baseUrl string, configFilePath string) {
	connector, err := NewArangoDbConnector(baseUrl, "root", "ohve7uthePhi", "githist_iop")
	if err != nil {
		panic(err)
	}

	configStr, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Print(err)
	}

	var arangoTests []ArangoTest

	err = json.Unmarshal([]byte(configStr), &arangoTests)

	if err != nil {
		panic(fmt.Sprint("Unable to parse config file ", err))
	}

	for _, arangoTest := range arangoTests {
		executeArangoTest(connector, arangoTest.Query, arangoTest.Params)
	}
}

func executeArangoTest(connector *ArangoConnector, querystring string, params []map[string]string) {

	for _, param := range params {
		queryWithParams := substituteParams(querystring, param, "@", "")
		watch := stopwatch.Start()
		fmt.Printf("%s", queryWithParams)

		_, err := connector.ExecuteQuery(queryWithParams)
		if err != nil {
			panic(err)
		}

		watch.Stop()

		dur := watch.Milliseconds() * 1000 * 1000
		fmt.Printf("\t%v\n", dur.Milliseconds())
	}
}

func apiBenchmark(baseUrl string, configFilePath string) {
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
		apiWithParams := substituteParams(api, param, "{", "}")

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

func substituteParams(api string, params map[string]string, paramPrefix string, paramPostfix string) string {
	afterReplace := api
	for k, v := range params {
		afterReplace = strings.Replace(afterReplace, paramPrefix+k+paramPostfix, v, -1)
	}
	return afterReplace
}

func inRed(str string) string {
	return Red + str + Reset
}
