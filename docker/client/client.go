package client

import (
	"fmt"
	"net"
	"net/http"

	"github.com/charmbracelet/lipgloss"
)

var transport = &http.Transport{
	Dial: func(network, addr string) (net.Conn, error) {
		return net.Dial("unix", "/var/run/docker.sock")
	},
}

func GetDockerHttpClient() *http.Client {
	return &http.Client{
		Transport: transport,
	}
}

var (
	client                = GetDockerHttpClient()
	titleErrorStyle       = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#FF00000")).Foreground(lipgloss.Color("#FFFFFFF")).MarginLeft(4).Render
	errorStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00000")).MarginLeft(8).Render
	errorDescriptionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00000")).MarginLeft(4).Render
)

func TestConnection() {
	_, err := client.Get("http://v1.47/containers/json?all=true")

	if err != nil {
		fmt.Println(titleErrorStyle("DOCKY CONNECTION ERROR:"))
		fmt.Println(errorDescriptionStyle("Docky can't connect to the docker engine unix socket"))
		fmt.Println(errorDescriptionStyle("Here are some posible solutions:\n"))
		fmt.Println(errorStyle("· If you are linux user you probably don't have permissions as a non root user"))
		fmt.Println(errorStyle("try to add your user to the docker group"))
		fmt.Println(errorStyle("see: https://docs.docker.com/engine/install/linux-postinstall/ \n"))
		fmt.Println(errorStyle("· If you are macos/windows user you probably don't have docker desktop running \n"))
	}

}
