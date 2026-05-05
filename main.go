package main

import (
	"fmt"
	"os"
	"time"

	ButtonCreateStartStop "goLogMyTime/components/ButtonCreateStartStopTask"
	ButtonStartStop "goLogMyTime/components/ButtonStartStopTask"
	"goLogMyTime/components/EmbeddedMarkdown"
	NowDoing "goLogMyTime/components/NowDoing"
	"goLogMyTime/components/TasksFavorite"
	"goLogMyTime/components/TasksRecent"
	"goLogMyTime/components/TextInputWrapper"
	"goLogMyTime/components/TimeTable"
	TodaysProgress "goLogMyTime/components/TodaysProgress"
	"goLogMyTime/components/Tutorial"
	"goLogMyTime/components/common/globals"
	"goLogMyTime/domain/taskStore"
	"goLogMyTime/domain/types"
	file "goLogMyTime/utils/file"
	"goLogMyTime/utils/layout"
	"goLogMyTime/utils/mouse"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type tickToSyncMsg time.Time

func tickToSync() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickToSyncMsg(t)
	})
}

type BaseModel struct {
	dataSource       *taskStore.TaskStore
	tutorial         *Tutorial.Model
	nowDoing         *NowDoing.Model
	todaysProgress   *TodaysProgress.Model
	tasksRecent      *TasksRecent.Model
	tasksFavorite    *TasksFavorite.Model
	timetable        *TimeTable.Model
	embeddedMarkdown *EmbeddedMarkdown.Model
}

func (m BaseModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(
		cmds,
		tickToSync(),
	)

	return tea.Batch(cmds...)
}

var prevTaskActivePtr *types.TaskRecord

func (m BaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	taskActivePtr := m.dataSource.GetTaskRecordActive()
	if taskActivePtr != prevTaskActivePtr {
		prevTaskActivePtr = taskActivePtr
		if taskActivePtr == nil {
			*m.nowDoing.TextInputDescription.StringPtr = ""
			*m.nowDoing.TextInputTask.StringPtr = ""
		} else {
			*m.nowDoing.TextInputDescription.StringPtr = taskActivePtr.Description
			*m.nowDoing.TextInputTask.StringPtr = taskActivePtr.Ticket
		}
		m.nowDoing.TextInputDescription.TextInput.SetValue(*m.nowDoing.TextInputDescription.StringPtr)
		m.nowDoing.TextInputTask.TextInput.SetValue(*m.nowDoing.TextInputTask.StringPtr)
	}

	switch msg := msg.(type) {
	case tickToSyncMsg:
		m.dataSource.Sync()
		cmds = append(cmds, tickToSync())

	case tea.WindowSizeMsg:
		sectionWidth := layout.GetSectionWidth(msg.Width)
		m.nowDoing.Width = sectionWidth
		m.tasksFavorite.Width = sectionWidth
		m.tasksRecent.Width = sectionWidth

	case tea.MouseClickMsg:
		m.nowDoing.IsFocused = mouse.IsLeftClickInBoundsOf(m.nowDoing.Zone, msg)
		m.tasksFavorite.IsFocused = mouse.IsLeftClickInBoundsOf(m.tasksFavorite.Zone, msg)
		m.tasksRecent.IsFocused = mouse.IsLeftClickInBoundsOf(m.tasksRecent.Zone, msg)
		m.timetable.IsFocused = mouse.IsLeftClickInBoundsOf(m.timetable.Zone, msg)

		if mouse.IsLeftClickInBoundsOf(globals.TextInputTaskZone, msg) {
			cmds = append(cmds, m.nowDoing.TextInputTask.TextInput.Focus())
		} else {
			m.nowDoing.TextInputTask.TextInput.Blur()
		}

		if mouse.IsLeftClickInBoundsOf(globals.TextInputDescriptionZone, msg) {
			cmds = append(cmds, m.nowDoing.TextInputDescription.TextInput.Focus())
		} else {
			m.nowDoing.TextInputDescription.TextInput.Blur()
		}

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "1", "2", "3", "4", "5":
			if m.AreTextInputsBlurred() {
				cmds = append(cmds, m.FocusSection(msg.String()))
			}
		case "h", "j", "k", "l", "down", "up", "tab", "shift+tab":
			keyString := msg.String()

			isTabKey := keyString == "tab"
			isShiftTabKey := keyString == "shift+tab"
			isHKey := keyString == "h"
			isJKey := keyString == "j"
			isKKey := keyString == "k"
			isLKey := keyString == "l"
			isDownKey := keyString == "down"
			isUpKey := keyString == "up"

			isDown := isTabKey || isJKey || isDownKey
			isUp := isShiftTabKey || isKKey || isUpKey

			switch true {

			case m.nowDoing.IsFocused:
				if isTabKey || isShiftTabKey {
					if m.nowDoing.TextInputTask.TextInput.Focused() {
						cmds = append(cmds, m.nowDoing.TextInputDescription.TextInput.Focus())
						m.nowDoing.TextInputTask.TextInput.Blur()
						break
					} else {
						cmds = append(cmds, m.nowDoing.TextInputTask.TextInput.Focus())
						m.nowDoing.TextInputDescription.TextInput.Blur()
						break
					}
				}

				if isDown && m.AreTextInputsBlurred() {
					cmds = append(cmds, m.FocusSection("2"))
					break
				}

			case m.tasksRecent.IsFocused:
				if isUp {
					isWithinSection := m.dataSource.DecrementOrResetTasksRecentCursor()
					if !isWithinSection {
						cmds = append(cmds, m.FocusSection("1"))
						break
					}
				} else if isDown {
					isWithinSection := m.dataSource.IncrementOrResetTasksRecentCursor()
					if !isWithinSection {
						cmds = append(cmds, m.FocusSection("3"))
						break
					}
				}

			case m.tasksFavorite.IsFocused:
				if isUp {
					isWithinSection := m.dataSource.DecrementOrResetTasksFavoriteCursor()
					if !isWithinSection {
						cmds = append(cmds, m.FocusSection("2"))
						break
					}
				} else if isDown {
					isWithinSection := m.dataSource.IncrementOrResetTasksFavoriteCursor()
					if !isWithinSection {
						cmds = append(cmds, m.FocusSection("4"))
						break
					}
				}

			case m.timetable.IsFocused:
				if isTabKey || isLKey {
					m.dataSource.TimeTableDaysOffsetFromToday++
				} else if isShiftTabKey || isHKey {
					m.dataSource.TimeTableDaysOffsetFromToday--
				} else if isUp {
					cmds = append(cmds, m.FocusSection("3"))
					break
				}

			}

		case "enter", "space":
			isEnterKey := msg.String() == "enter"

			switch true {

			case m.nowDoing.IsFocused:
				if isEnterKey {
					taskDefinition := m.dataSource.GetInputTaskDefinition()
					if taskDefinition != nil {
						m.dataSource.TaskToggle(taskDefinition)
					} else if m.dataSource.IsInputTaskValidAndNew() {
						m.dataSource.TaskCreate()
					}
				}

			case m.tasksRecent.IsFocused:
				taskDefinitions := m.dataSource.GetTaskDefinitionsRecent(true)
				length := len(taskDefinitions)
				if length == 0 {
					break
				}

				taskDefinition := taskDefinitions[m.dataSource.TasksRecentCursor]
				m.dataSource.TaskToggle(&taskDefinition)

			case m.tasksFavorite.IsFocused:
				taskDefinitions := m.dataSource.GetTaskDefinitionsFavorite()
				length := len(taskDefinitions)
				if length == 0 {
					break
				}

				taskDefinition := taskDefinitions[m.dataSource.TasksFavoriteCursor]
				m.dataSource.TaskToggle(&taskDefinition)

			}
		}
	}

	u1, cmd := m.nowDoing.Update(msg)
	m.nowDoing = &u1
	cmds = append(cmds, cmd)

	u2, cmd := m.tasksFavorite.Update(msg)
	m.tasksFavorite = &u2
	cmds = append(cmds, cmd)

	u3, cmd := m.tasksRecent.Update(msg)
	m.tasksRecent = &u3
	cmds = append(cmds, cmd)

	u4, cmd := m.todaysProgress.Update(msg)
	m.todaysProgress = &u4
	cmds = append(cmds, cmd)

	u5, cmd := m.timetable.Update(msg)
	m.timetable = &u5
	cmds = append(cmds, cmd)

	u6, cmd := m.embeddedMarkdown.Update(msg)
	m.embeddedMarkdown = &u6
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m BaseModel) View() tea.View {
	var view tea.View
	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion

	view.SetContent(
		zone.Scan(
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				lipgloss.JoinVertical(
					lipgloss.Right,
					m.todaysProgress.View().Content,
					m.nowDoing.View().Content,
					m.tasksRecent.View().Content,
					m.tasksFavorite.View().Content,
					m.timetable.View().Content,
					m.tutorial.View().Content,
				),
				m.embeddedMarkdown.View().Content,
			),
		),
	)

	return view
}

func (m BaseModel) FocusSection(shortcutDigit string) tea.Cmd {
	var cmd tea.Cmd

	m.nowDoing.IsFocused = shortcutDigit == m.nowDoing.Shortcut
	if m.nowDoing.IsFocused {
		return m.nowDoing.TextInputTask.TextInput.Focus()
	}

	m.tasksFavorite.IsFocused = shortcutDigit == m.tasksFavorite.Shortcut
	m.tasksRecent.IsFocused = shortcutDigit == m.tasksRecent.Shortcut
	m.timetable.IsFocused = shortcutDigit == m.timetable.Shortcut
	m.embeddedMarkdown.IsFocused = shortcutDigit == m.embeddedMarkdown.Shortcut

	if m.embeddedMarkdown.IsFocused {
		return m.embeddedMarkdown.Textarea.Focus()
	}

	return cmd
}

func (m BaseModel) AreTextInputsBlurred() bool {
	return !m.nowDoing.TextInputDescription.TextInput.Focused() &&
		!m.nowDoing.TextInputTask.TextInput.Focused() &&
		!m.embeddedMarkdown.Textarea.Focused()
}

func main() {
	zone.NewGlobal()
	file.CreateFilesIfMissing()
	dataSource := taskStore.New()
	dataSource.Sync()

	initModel := EmbeddedMarkdown.InitialModel(file.PathMarkdown, dataSource)

	app := BaseModel{
		dataSource: dataSource,
		tutorial:   &Tutorial.Model{},
		tasksRecent: &TasksRecent.Model{
			DataSource: dataSource,
			IsFocused:  false,
			Shortcut:   globals.TasksRecentSection.Shortcut,
			Zone:       globals.TasksRecentSection.Zone,

			ButtonsStartStopMap: map[string]ButtonStartStop.Model{},
		},
		tasksFavorite: &TasksFavorite.Model{
			DataSource: dataSource,
			IsFocused:  false,
			Shortcut:   globals.TasksFavoriteSection.Shortcut,
			Zone:       globals.TasksFavoriteSection.Zone,

			ButtonsStartStopMap: map[string]ButtonStartStop.Model{},
		},
		nowDoing: &NowDoing.Model{
			DataSource: dataSource,
			IsFocused:  false,
			Shortcut:   globals.NowDoingSection.Shortcut,
			Zone:       globals.NowDoingSection.Zone,

			ButtonCreateStartStop: &ButtonCreateStartStop.Model{
				DataSource: dataSource,
				IdPrefix:   "N_",
			},
			TextInputDescription: &TextInputWrapper.Model{
				Zone:      globals.TextInputDescriptionZone,
				TextInput: textinput.New(),
				StringPtr: &dataSource.InputTaskDescription,
			},
			TextInputTask: &TextInputWrapper.Model{
				Zone:      globals.TextInputTaskZone,
				TextInput: textinput.New(),
				StringPtr: &dataSource.InputTaskTicket,
			},
		},
		todaysProgress: &TodaysProgress.Model{
			DataSource: dataSource,
		},
		timetable: &TimeTable.Model{
			DataSource: dataSource,
			Shortcut:   globals.TimeTableSection.Shortcut,
			Zone:       globals.TimeTableSection.Zone,
		},
		embeddedMarkdown: &initModel,
	}

	app.nowDoing.TextInputDescription.TextInput.Placeholder = "uwa"
	app.nowDoing.TextInputDescription.TextInput.Prompt = "info: "
	app.nowDoing.TextInputDescription.TextInput.SetVirtualCursor(true)
	app.nowDoing.TextInputDescription.TextInput.CharLimit = 10
	app.nowDoing.TextInputDescription.TextInput.SetWidth(10)

	app.nowDoing.TextInputTask.TextInput.Placeholder = "bep"
	app.nowDoing.TextInputTask.TextInput.Prompt = "task: "
	app.nowDoing.TextInputTask.TextInput.SetVirtualCursor(true)
	app.nowDoing.TextInputTask.TextInput.CharLimit = 10
	app.nowDoing.TextInputTask.TextInput.SetWidth(10)

	p := tea.NewProgram(app)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
