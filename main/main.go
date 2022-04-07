package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	driver "github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
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

	conn, err := arangohttp.NewConnection(arangohttp.ConnectionConfig{
		Endpoints: []string{"http://10.125.29.146:8540"},
	})

	if err != nil {
		panic("Http connection failed")
	}

	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication("root", "ohve7uthePhi"),
	})

	if err != nil {
		panic("Client open failed")
	}

	db, err := client.Database(nil, "githist_iop")
	if err != nil {
		panic("DB open failed")
	}

	var cursor driver.Cursor

	querystring := "FOR commit_obj in 1..10000000 ANY 'issue/HIVE-4019' ANY issueAlias, INBOUND autoAlias, INBOUND commitIssue OPTIONS {uniqueVertices: \"global\", bfs:true} FILTER IS_SAME_COLLECTION(commit, commit_obj) let branches = UNIQUE(FOR v in 1..10000000 ANY commit_obj INBOUND parentCommit, INBOUND branchHead OPTIONS {uniqueVertices: \"global\", bfs:true} FILTER IS_SAME_COLLECTION(branch, v) RETURN {name:v.branchName, remote:v.remoteURI}) RETURN { hash:commit_obj.hash, message:commit_obj.message, componentName:commit_obj.componentName, tickets:commit_obj.tickets, commitTime: commit_obj.commitTime, isRevert: commit_obj.isRevert, branches: branches }"

	watch := stopwatch.Start()
	cursor, err = db.Query(nil, querystring, nil)
	watch.Stop()
	fmt.Printf(inRed("Total: %v\n"), watch.String())

	if err != nil {
		fmt.Println("Failed query")
		fmt.Println(err)
	}

	defer cursor.Close()

	panic("DONE")

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
