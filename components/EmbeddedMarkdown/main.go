package EmbeddedMarkdown

import (
	"goLogMyTime/components/common/globals"
	commonStyles "goLogMyTime/components/common/styles"
	"goLogMyTime/domain/taskStore"
	"os"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type saveFileMsg struct {
	err error
}

type Model struct {
	Zone      string
	Shortcut  string
	IsFocused bool
	Filepath  string
	Textarea  textarea.Model
	status    string
	Width     int
}

func saveFileCmd(path string, content string) tea.Cmd {
	return func() tea.Msg {
		err := os.WriteFile(path, []byte(content), 0644)
		return saveFileMsg{err: err}
	}
}

func InitialModel(filepath string, dataSource *taskStore.TaskStore) Model {
	content, err := os.ReadFile(filepath)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	ta := textarea.New()
	ta.SetValue(string(content))
	ta.ShowLineNumbers = true
	ta.Prompt = ""

	styles := ta.Styles()
	styles.Focused.CursorLine = lipgloss.NewStyle()
	styles.Focused.Text = lipgloss.NewStyle()
	ta.SetStyles(styles)

	return Model{
		Shortcut: globals.EmbeddedFileSection.Shortcut,
		Zone:     globals.EmbeddedFileSection.Zone,
		Filepath: filepath,
		Textarea: ta,
		status:   "Loaded: " + filepath,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Textarea.SetHeight(msg.Height - 4)
		m.Textarea.SetWidth(m.Width - 2)

	case globals.FocusSectionMsg:
		if string(msg) == m.Shortcut {
			cmds = append(cmds, m.focus())
		}

	case globals.BlurSectionMsg:
		if string(msg) == m.Shortcut {
			m.blur()
		}

	case tea.KeyPressMsg:
		if m.Textarea.Focused() {
			switch msg.String() {
			case "esc":
				m.blur()
				return m, nil

			default:
				cmds = append(cmds, saveFileCmd(m.Filepath, m.Textarea.Value()))
			}
		}

	case saveFileMsg:
		if msg.err != nil {
			m.status = "Error while saving: " + msg.err.Error()
		} else {
			m.status = "Saved " + m.Filepath
		}
	}

	var taCmd tea.Cmd
	m.Textarea, taCmd = m.Textarea.Update(msg)
	cmds = append(cmds, taCmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) focus() tea.Cmd {
	m.IsFocused = true
	return m.Textarea.Focus()
}

func (m *Model) blur() {
	m.IsFocused = false
	m.Textarea.Blur()
}

func (m Model) View() tea.View {
	title := lipgloss.Sprintf("[%s] Notes", m.Shortcut)

	return tea.NewView(zone.Mark(
		m.Zone,
		m.getStyle().Padding(0, 0, 0, 2).Width(m.Width).UnsetAlign().Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					commonStyles.SectionTitle.Render(title),
					" ",
					commonStyles.EmbeddedStatusBar.Render(m.status),
					"\n",
				),
				m.Textarea.View(),
			),
		),
	))
}

func (m *Model) getStyle() lipgloss.Style {
	if m.IsFocused {
		return commonStyles.EmbeddedSectionContainerFocused
	}

	return commonStyles.EmbeddedSectionContainerBlurred
}
