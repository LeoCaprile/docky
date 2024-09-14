package main

import (
	"docky/models"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var transport = &http.Transport{
	Dial: func(network, addr string) (net.Conn, error) {
		return net.Dial("unix", "/var/run/docker.sock")
	},
}

func GetHttpClient() *http.Client {
	return &http.Client{
		Transport: transport,
	}
}

type UpdateMsg = struct{}

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.table.MoveUp(1)
		case "down", "j":
			m.table.MoveDown(1)
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case UpdateMsg:
		{
			t := GetTableRows()
			m.table.SetRows(t)
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Render(m.table.View() + "\n" + m.table.HelpView() + "\n")
}

func getDockerContainers() ([]docker.Container, error) {
	client := GetHttpClient()

	res, err := client.Get("http://v1.47/containers/json")

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	var containers []docker.Container

	errDeco := json.NewDecoder(res.Body).Decode(&containers)

	if errDeco != nil {
		log.Fatal(err)
	}

	return containers, nil

}

func GetTableRows() []table.Row {
	containers, err := getDockerContainers()
	if err != nil {
		log.Fatal(err)
	}
	var rows []table.Row

	for _, container := range containers {
		naturalTime := humanize.Time(time.Unix(int64(container.Created), 0))
		toInsertRow := table.Row{container.ID, container.Image, container.Command, naturalTime, container.Status, container.Names[0]}
		rows = append(rows, toInsertRow)
	}

	return rows
}

func UpdateContainers() table.Model {
	containers, err := getDockerContainers()

	if err != nil {
		log.Fatal(err)
	}

	columns := []table.Column{
		{Title: "Id", Width: 20},
		{Title: "Image", Width: 20},
		{Title: "Command", Width: 20},
		{Title: "Created", Width: 20},
		{Title: "Status", Width: 20},
		{Title: "Name", Width: 20},
	}

	var rows []table.Row

	for _, container := range containers {
		naturalTime := humanize.Time(time.Unix(int64(container.Created), 0))
		toInsertRow := table.Row{container.ID, container.Image, container.Command, naturalTime, container.Status, container.Names[0]}
		rows = append(rows, toInsertRow)
	}

	t := table.New(table.WithColumns(columns), table.WithRows(rows), table.WithFocused(true), table.WithHeight(7))

	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("12")).
		Bold(false)
	t.SetStyles(s)

	return t
}

func getInitialModel() model {
	t := UpdateContainers()
	return model{table: t}
}

func updateTable(program *tea.Program) {
	for {
		time.Sleep(1000)
		program.Send(UpdateMsg{})
	}
}

func main() {

	program := tea.NewProgram(getInitialModel())

	go updateTable(program)

	if _, err := program.Run(); err != nil {
		fmt.Println("Error running programgetInitialModel", err)
		os.Exit(1)
	}

}
