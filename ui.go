package main

import (
	"fmt"
	"io"
	"os"

	listview "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	viewport viewport.Model
	list     listview.Model
}

type navItem string

func (n navItem) FilterValue() string { return string(n) }

type navItemDelegate struct{}

func (navItemDelegate) Height() int                                 { return 1 }
func (navItemDelegate) Spacing() int                                { return 0 }
func (navItemDelegate) Update(_ tea.Msg, _ *listview.Model) tea.Cmd { return nil }
func (navItemDelegate) Render(w io.Writer, m listview.Model, index int, listItem listview.Item) {
	i, ok := listItem.(navItem)
	if !ok {
		return
	}

	fmt.Fprintf(w, "%s", i)
}

func text() string {
	res := ""
	for i := 0; i < 100; i++ {
		res += fmt.Sprintf("%d\n———————\n", i)
	}
	return res
}

func NewModel() *model {
	m := &model{
		list:     buildList(),
		viewport: viewport.New(10, 10),
	}

	m.viewport.SetContent(text())

	return m
}

func buildList() listview.Model {
	var rows []listview.Item
	for i := 0; i < 10; i++ {
		rows = append(rows, navItem(fmt.Sprintf("Row %d", i)))
	}

	navigation := listview.New(rows, navItemDelegate{}, 20, 100)

	return navigation
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "esc" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		navWidth := 20

		m.viewport.Width = msg.Width - navWidth
		m.viewport.Height = msg.Height
		m.list.SetHeight(msg.Height - 2)
	}

	// m.navigation, cmd = m.navigation.Update(msg)
	// cmds = append(cmds, cmd)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	panelStyle := lipgloss.NewStyle().Margin(1, 0)
	sidebar := panelStyle.Render("bug") + "\n" + m.list.View()
	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, m.viewport.View())
}

func main() {
	p := tea.NewProgram(
		NewModel(),
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
