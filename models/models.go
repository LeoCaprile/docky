package models

import (
	"docky/docker/container"
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

type UpdateTableRowsMsg struct{}
type RemoveErrMsg struct{}

type message struct {
	error error
	info  string
}

type model struct {
	table table.Model
	message
}

var (
	notificationStyles = lipgloss.NewStyle().Background(lipgloss.Color("#6495ED")).Foreground(lipgloss.Color("#FFFFFF")).Render
	errorStyles        = lipgloss.
				NewStyle().
				Background(lipgloss.TerminalColor(lipgloss.Color("#FF0000"))).
				Foreground(lipgloss.Color("#FFFFFF")).Render
)

func updateTable() tea.Msg {
	time.Sleep(time.Duration(500) * time.Millisecond)
	return UpdateTableRowsMsg{}
}

func removeAllMsg() tea.Msg {
	time.Sleep(time.Duration(3) * time.Second)
	return RemoveErrMsg{}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(updateTable)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.table.MoveUp(1)
		case "down", "j":
			m.table.MoveDown(1)
		case "s":
			row := m.table.SelectedRow()
			errStop := container.Stop(row[0])
			if errStop != nil {
				m.message.error = errStop
			}
			return m, removeAllMsg
		case "u":
			row := m.table.SelectedRow()
			errStart := container.Start(row[0])
			if errStart != nil {
				m.message.error = errStart
			}
			return m, removeAllMsg
		case "r":
			row := m.table.SelectedRow()
			errRestart := container.Restart(row[0])
			if errRestart != nil {
				m.message.error = errRestart
				return m, removeAllMsg
			}
			m.message.info = "Container" + lipgloss.NewStyle().Bold(true).Render(row[0]) + "restarted..."
			return m, removeAllMsg
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case UpdateTableRowsMsg:
		t := getContainerTableRows()
		m.table.SetRows(t)
		return m, updateTable
	case RemoveErrMsg:
		m.message = message{}
	}

	return m, nil
}

func (m model) View() string {
	var msg string

	if m.message.error != nil {
		msg = errorStyles(fmt.Sprint(m.message.error))
	}
	if m.message.info != "" {
		msg = notificationStyles(m.message.info)
	}

	table := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Render(m.table.View() + "\n" + m.table.HelpView() + "\n" + msg)

	return table

}

func getContainerTableRows() []table.Row {
	containers, err := container.GetAll()

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
	containers, err := container.GetAll()

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

func GetInitialModel() model {
	t := UpdateContainers()
	return model{table: t, message: message{}}
}
