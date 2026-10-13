package NowDoing

import (
	ButtonCreateStartStop "goLogMyTime/components/ButtonCreateStartStopTask"
	"goLogMyTime/components/TextInputWrapper"
	"goLogMyTime/components/common/globals"
	"goLogMyTime/domain/taskStore"
	"goLogMyTime/domain/types"
	"goLogMyTime/utils/mouse"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type Model struct {
	DataSource            *taskStore.TaskStore
	IsFocused             bool
	Shortcut              string
	Zone                  string
	TextInputTask         *TextInputWrapper.Model
	TextInputDescription  *TextInputWrapper.Model
	ButtonCreateStartStop *ButtonCreateStartStop.Model
	Width                 int
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case globals.FocusSectionMsg:
		if string(msg) == m.Shortcut {
			isAnyFieldFocused := m.TextInputDescription.TextInput.Focused() ||
				m.TextInputTask.TextInput.Focused()
			cmds = append(cmds, m.focus(isAnyFieldFocused))
		}

	case globals.BlurSectionMsg:
		if string(msg) == m.Shortcut {
			m.Blur()
		}

	case tea.MouseClickMsg:
		if m.DataSource.TaskInputState == globals.TASK_INPUT_VALID_ACTIVE {
			break
		}

		if mouse.IsLeftClickInBoundsOf(m.TextInputTask.Zone, msg) {
			m.TextInputDescription.Blur()
			cmds = append(cmds, m.TextInputTask.Focus())
		} else if mouse.IsLeftClickInBoundsOf(m.TextInputDescription.Zone, msg) {
			m.TextInputTask.Blur()
			cmds = append(cmds, m.TextInputDescription.Focus())
		}

	case tea.KeyPressMsg:
		if m.DataSource.TaskInputState == globals.TASK_INPUT_VALID_ACTIVE {
			break
		}

		keyString := msg.String()
		isTabKey := keyString == "tab"
		isShiftTabKey := keyString == "shift+tab"
		if !isTabKey && !isShiftTabKey {
			break
		}

		if m.TextInputDescription.TextInput.Focused() {
			m.TextInputDescription.TextInput.Blur()
			cmds = append(cmds, m.TextInputTask.Focus())
		} else if m.TextInputTask.TextInput.Focused() {
			m.TextInputTask.TextInput.Blur()
			cmds = append(cmds, m.TextInputDescription.Focus())
		}
	}

	u1, cmd := m.TextInputTask.Update(msg)
	m.TextInputTask = &u1
	cmds = append(cmds, cmd)

	u2, cmd := m.TextInputDescription.Update(msg)
	m.TextInputDescription = &u2
	cmds = append(cmds, cmd)

	u3, cmd := m.ButtonCreateStartStop.Update(msg)
	m.ButtonCreateStartStop = &u3
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) focus(isAnyFieldFocused bool) tea.Cmd {
	m.IsFocused = true

	if m.DataSource.TaskInputState != globals.TASK_INPUT_VALID_ACTIVE &&
		!isAnyFieldFocused {
		m.TextInputDescription.TextInput.Blur()
		return m.TextInputTask.TextInput.Focus()
	}

	return nil
}

func (m *Model) Blur() {
	m.IsFocused = false
	m.BlurOnlyInputs()
}

func (m *Model) BlurOnlyInputs() {
	m.TextInputDescription.TextInput.Blur()
	m.TextInputTask.TextInput.Blur()
}

func (m Model) View() tea.View {
	style := m.DataSource.Style

	return tea.NewView(zone.Mark(
		m.Zone,
		style.GetSectionStyle(m.IsFocused).Width(m.Width).Render(
			lipgloss.JoinVertical(
				lipgloss.Right,
				style.SectionTitle.Render(m.getTitle()),
				lipgloss.JoinHorizontal(
					lipgloss.Bottom,
					m.TextInputDescription.View().Content,
					m.TextInputTask.View().Content,
					m.ButtonCreateStartStop.View().Content,
				),
			),
		),
	))
}

func (m Model) getTitle() string {
	taskRecordActive := m.DataSource.GetTaskRecordActive()

	if taskRecordActive == nil {
		return lipgloss.Sprintf("Start doing [%s]", m.Shortcut)
	}

	return lipgloss.Sprintf("Now doing [%s]", m.Shortcut)
}

var DESCRIPTION_MAX_LENGTH = 6 // "review"
var DESCRIPTION_MAX_INDEX = DESCRIPTION_MAX_LENGTH - 1

func getStringNowDoing(taskRecordActive *types.TaskRecord) string {
	var sb strings.Builder
	if len(taskRecordActive.Description) > 0 {
		sb.WriteString("\"")
		if len(taskRecordActive.Description) > DESCRIPTION_MAX_LENGTH {
			runes := []rune(taskRecordActive.Description)
			sb.WriteString(string(runes[0:DESCRIPTION_MAX_INDEX]) + "...")
		} else {
			sb.WriteString(taskRecordActive.Description)
		}
		sb.WriteString("\" ")
	}

	return sb.String() + taskRecordActive.Ticket + " "
}
