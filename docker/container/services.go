package container

import (
	"docky/docker/types"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
)

var baseURL = "http://v1.47"

var transport = &http.Transport{
	Dial: func(network, addr string) (net.Conn, error) {
		return net.Dial("unix", "/var/run/docker.sock")
	},
}

func getHttpClient() *http.Client {
	return &http.Client{
		Transport: transport,
	}
}

var client = getHttpClient()

func GetAll() ([]types.Container, error) {
	res, err := client.Get(baseURL + "/containers/json?all=true")

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	var containers []types.Container

	errDeco := json.NewDecoder(res.Body).Decode(&containers)

	if errDeco != nil {
		log.Fatal(err)
	}

	return containers, nil
}

func Stop(containerId string) error {
	res, err := client.Post(baseURL+"/containers/"+containerId+"/stop", "application/json", nil)

	if err != nil {
		return fmt.Errorf("Couldn't make the request %w", err)
	}

	if res.StatusCode == 304 {
		return fmt.Errorf("The container has already stoped")
	}

	if res.StatusCode == 404 {
		return fmt.Errorf("The container don't exist")
	}

	if res.StatusCode == 500 {
		return fmt.Errorf("Server error")
	}

	return nil
}

func Start(containerId string) error {
	res, err := client.Post(baseURL+"/containers/"+containerId+"/start", "application/json", nil)

	if err != nil {
		return fmt.Errorf("Couldn't make the request %w", err)
	}

	if res.StatusCode == 304 {
		return fmt.Errorf("The container has already started")
	}

	if res.StatusCode == 404 {
		return fmt.Errorf("The container don't exist")
	}

	if res.StatusCode == 500 {
		return fmt.Errorf("Server error")
	}

	return nil
}

func Restart(containerId string) error {
	res, err := client.Post(baseURL+"/containers/"+containerId+"/start", "application/json", nil)

	if err != nil {
		return fmt.Errorf("Couldn't make the request %w", err)
	}

	if res.StatusCode == 404 {
		return fmt.Errorf("The container don't exist")
	}

	if res.StatusCode == 500 {
		return fmt.Errorf("Server error")
	}

	return nil
}
