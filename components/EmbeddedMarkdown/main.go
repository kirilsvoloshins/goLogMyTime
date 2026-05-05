package EmbeddedMarkdown

import (
	"goLogMyTime/components/common/globals"
	commonStyles "goLogMyTime/components/common/styles"
	"goLogMyTime/domain/taskStore"
	"goLogMyTime/utils/layout"
	"goLogMyTime/utils/mouse"
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
	height    int
	width     int
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
	ta.ShowLineNumbers = false

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
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height - 4
		m.width = layout.GetEditorWidth(msg.Width)

	case tea.MouseClickMsg:
		// TODO: handle in parent component
		if mouse.IsLeftClickInBoundsOf(m.Zone, msg) {
			cmd = m.Textarea.Focus()
			m.IsFocused = true
			return m, cmd
		} else {
			m.IsFocused = false
			m.Textarea.Blur()
			return m, nil
		}

	case tea.KeyPressMsg:
		// TODO: handle in parent component
		if m.Textarea.Focused() {
			switch msg.String() {
			case "esc":
				m.IsFocused = false
				m.Textarea.Blur()
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

	// if m.Textarea.Focused() {
	// 	if _, ok := msg.(tea.KeyPressMsg); ok {
	// 		cmds = append(cmds, saveFileCmd(m.Filepath, m.Textarea.Value()))
	// 	}
	// }

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	var title = lipgloss.Sprintf("Notes [%s]", m.Shortcut)

	m.Textarea.SetHeight(m.height)
	m.Textarea.SetWidth(m.width)

	editorView := m.Textarea.View()
	statusBar := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(m.status)

	baseStyle := commonStyles.EmbeddedSectionContainerBlurred
	if m.IsFocused {
		baseStyle = commonStyles.EmbeddedSectionContainerFocused
	}

	return tea.NewView(zone.Mark(
		m.Zone,
		baseStyle.Width(m.width).UnsetAlign().Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					commonStyles.SectionTitle.Render(title),
					" ",
					statusBar,
					"\n",
				),
				editorView,
			),
		),
	))
}
