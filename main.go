package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"mapreducego/client"
)

// take in a file
// split it into parts
// send it to other nodes - how? http? rpc? sockets?
// get back data

var (
	nodeFlag    = flag.Bool("node", false, "is this instance a node or not?")
	addressFlag = flag.String("address", "", "The workers need addresses to send requests to")
	idFlag      = flag.Int("id", 0, "The workers have IDs")
)

type Nodes struct {
	NumberOfNodes int      `json:"nodes"`
	Locations     []string `json:"node_locations"`
}

func readConfig(filename string) Nodes {
	configFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := io.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var nodes Nodes
	if err := json.Unmarshal(bytes, &nodes); err != nil {
		log.Fatal(err)
	}

	return nodes
}

type file struct {
	contents []byte
	node     Nodes
}

func newFile(contents []byte, node Nodes) file {
	return file{contents, node}
}

func readFile(filename string, node Nodes) file {
	contents, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	return newFile(contents, node)
}

func (f *file) split_content() {
	s := bytes.Split(f.contents, []byte("\n"))
	numberOfLines := len(s) - 1
	fmt.Println(numberOfLines)

	work := numberOfLines / f.node.NumberOfNodes
	fmt.Println(work)

	for i := 0; i < numberOfLines; i++ {
		jsonData := []byte(fmt.Sprintf(`{"data": "%s"}`, s[i]))
		fmt.Println(string(jsonData))
		f.send_to_nodes("http://127.0.0.1:9000", jsonData)
	}
}

func (f *file) send_to_nodes(address string, jsonData []byte) {
	req, err := http.NewRequest(http.MethodPost, address+"/mr", bytes.NewReader(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	_, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	if !*nodeFlag {
		n := readConfig("config.json")
		f := readFile("words.txt", n)
		f.split_content()

	} else {
		fmt.Println("booting up node...")
		client.InitializeWorker(*addressFlag, *idFlag)
	}
}
