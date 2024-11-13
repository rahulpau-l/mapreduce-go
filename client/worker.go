package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Worker struct {
	Address string
	NodeID  int
}

func InitializeWorker(address string, nodeID int) {
	worker := Worker{Address: address, NodeID: nodeID}
	worker.ListenForIncomingRequests()
}

func (w *Worker) Map(key, value string) {
	// TODO
}

func greetings(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

type info struct {
	Data string `json:"data"`
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}

func (w *Worker) ListenForIncomingRequests() {
	http.HandleFunc("/", greetings)
	http.HandleFunc("/mr", handlePost)
	log.Fatal(http.ListenAndServe(w.Address, nil))
}
