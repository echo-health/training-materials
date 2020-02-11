package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	env, err := readVolume()
	if err != nil {
		fmt.Printf("Failed reading volume")
	}
	fmt.Fprintf(w, "Hello, %s!", env)
}

type Env struct {
	environment string
}

func readVolume() (string, error) {
	output, err := ioutil.ReadFile("/etc/config/envfile")
	if err != nil {
		return "", err
	}
	env := Env{}
	err = json.Unmarshal(output, &env)
	if err != nil {
		return "", err
	}
	return env.environment, nil
}
