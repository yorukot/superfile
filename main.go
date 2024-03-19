package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	homeListWidth   = 20
	secondListWidth = 20 // 新增的第二个列表的宽度
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list1    list.Model // 第一个列表
	list2    list.Model // 第二个列表
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list1.SetWidth(homeListWidth)
		m.list1.SetHeight(msg.Height/2 - 5)
		m.list2.SetWidth(secondListWidth) // 设置第二个列表的宽度
		m.list2.SetHeight(msg.Height/2 - 5)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list1.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list1, cmd = m.list1.Update(msg)
	m.list2, cmd = m.list2.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Not hungry? That’s cool.")
	}
	return "\n" + m.list1.View() + "\n" + m.list2.View() // 在视图中显示两个列表
}

func main() {
	// 创建第一个列表
	items1 := []list.Item{
		item("Ramen"),
		item("Tomato Soup"),
		item("Hamburgers"),
		item("Cheeseburgers"),
		item("Currywurst"),
	}
	l1 := list.New(items1, itemDelegate{}, homeListWidth, homeListWidth)
	l1.Title = "Disk:"
	l1.SetShowStatusBar(false)
	l1.SetFilteringEnabled(false)
	l1.Styles.Title = titleStyle
	l1.Styles.PaginationStyle = paginationStyle
	l1.Styles.HelpStyle = helpStyle

	// 创建第二个列表
	items2 := []list.Item{
		item("Apples"),
		item("Oranges"),
		item("Bananas"),
		item("Grapes"),
		item("Pineapples"),
	}
	l2 := list.New(items2, itemDelegate{}, homeListWidth, secondListWidth)
	l2.Title = "Pinned:"
	l2.SetShowStatusBar(false)
	l2.SetFilteringEnabled(false)
	l2.Styles.Title = titleStyle
	l2.Styles.PaginationStyle = paginationStyle
	l2.Styles.HelpStyle = helpStyle

	// 创建模型
	m := model{list1: l1, list2: l2}

	// 启动 Bubble Tea 程序
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
