package EmbeddedMarkdown

import (
	"goLogMyTime/components/common/globals"
	"goLogMyTime/domain/taskStore"
	"goLogMyTime/utils/file"
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
	DataSource *taskStore.TaskStore
	Zone       string
	Shortcut   string
	IsFocused  bool
	Filepath   string
	Textarea   *textarea.Model
	status     string
	Width      int
}

func saveFileCmd(path string, content string) tea.Cmd {
	return func() tea.Msg {
		err := os.WriteFile(path, []byte(content), file.READ_WRITE_PERM)
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

	return Model{
		DataSource: dataSource,
		Shortcut:   globals.EmbeddedFileSection.Shortcut,
		Zone:       globals.EmbeddedFileSection.Zone,
		Filepath:   filepath,
		Textarea:   &ta,
		status:     "Loaded: " + filepath,
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

	case saveFileMsg:
		if msg.err != nil {
			m.status = "Error while saving: " + msg.err.Error()
		} else {
			m.status = "Saved " + m.Filepath
		}
	}

	var taCmd tea.Cmd
	ta, taCmd := m.Textarea.Update(msg)
	m.Textarea = &ta
	cmds = append(cmds, taCmd)

	if key, ok := msg.(tea.KeyMsg); ok && m.Textarea.Focused() {
		if key.String() == "esc" {
			m.blur()
		} else {
			cmds = append(cmds, saveFileCmd(m.Filepath, m.Textarea.Value()))
		}
	}

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

func (m *Model) View() tea.View {
	title := lipgloss.Sprintf("[%s] Notes", m.Shortcut)
	style := m.DataSource.Style
	styleSection := style.GetEmbeddedSectionStyle(m.IsFocused)

	styles := m.Textarea.Styles()
	styles.Focused.CursorLine = m.DataSource.Style.TextHighContrast
	styles.Focused.Text = m.DataSource.Style.TextHighContrast
	styles.Blurred.CursorLine = m.DataSource.Style.TextHighContrast
	styles.Blurred.Text = m.DataSource.Style.TextHighContrast
	m.Textarea.SetStyles(styles)

	return tea.NewView(zone.Mark(
		m.Zone,
		styleSection.Width(m.Width).Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					style.SectionTitle.Render(title),
					" ",
					style.EmbeddedStatusBar.Render(m.status),
					"\n",
				),
				m.Textarea.View(),
			),
		),
	))
}
