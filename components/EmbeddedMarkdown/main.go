package EmbeddedMarkdown

import (
	"goLogMyTime/common/globals"
	"goLogMyTime/store"
	"goLogMyTime/types"
	"goLogMyTime/utils"
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
	Store     *store.RootStore
	Zone      types.Zone
	Shortcut  types.SectionShortcut
	IsFocused bool
	Filepath  string
	Textarea  *textarea.Model
	status    string
	Width     int
}

func saveFileCmd(path string, content string) tea.Cmd {
	return func() tea.Msg {
		err := os.WriteFile(path, []byte(content), utils.READ_WRITE_PERM)
		return saveFileMsg{err: err}
	}
}

func InitialModel(filepath string, store *store.RootStore) Model {
	content, err := os.ReadFile(filepath)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	ta := textarea.New()
	ta.SetValue(string(content))
	ta.ShowLineNumbers = true
	ta.Prompt = ""

	return Model{
		Store:    store,
		Shortcut: globals.InteractiveSections.EmbeddedFile.Shortcut,
		Zone:     globals.InteractiveSections.EmbeddedFile.Zone,
		Filepath: filepath,
		Textarea: &ta,
		status:   "Loaded: " + filepath,
	}
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Textarea.SetHeight(msg.Height - 4)
		m.Textarea.SetWidth(m.Width - 2)

	case types.FocusZoneMsg:
		if types.Zone(msg) == m.Zone {
			m.IsFocused = true
			cmds = append(cmds, m.Textarea.Focus())
		}

	case types.BlurZoneMsg:
		if types.Zone(msg) == m.Zone {
			m.IsFocused = false
			m.Textarea.Blur()
		}

	case saveFileMsg:
		if msg.err != nil {
			m.status = "Error while saving: " + msg.err.Error()
		} else {
			m.status = "Saved " + m.Filepath
		}
	}

	ta, taCmd := m.Textarea.Update(msg)
	m.Textarea = &ta
	cmds = append(cmds, taCmd)

	if key, ok := msg.(tea.KeyMsg); ok && m.Textarea.Focused() {
		if key.String() == "esc" {
			cmds = append(cmds, globals.FocusZone(globals.InteractiveSections.NowDoing.Zone))
		} else {
			cmds = append(cmds, saveFileCmd(m.Filepath, m.Textarea.Value()))
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() tea.View {
	title := lipgloss.Sprintf("[%s] Notes", m.Shortcut)
	style := m.Store.Ui.Style
	styleSection := style.GetEmbeddedSectionStyle(m.IsFocused)

	styles := m.Textarea.Styles()
	styles.Focused.CursorLine = m.Store.Ui.Style.TextHighContrast
	styles.Focused.Text = m.Store.Ui.Style.TextHighContrast
	styles.Blurred.CursorLine = m.Store.Ui.Style.TextHighContrast
	styles.Blurred.Text = m.Store.Ui.Style.TextHighContrast
	m.Textarea.SetStyles(styles)

	return tea.NewView(zone.Mark(
		string(m.Zone),
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
