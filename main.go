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
	return tickToSync()
}

var prevTaskActivePtr *types.TaskRecord

func (m *BaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.dataSource.TaskInputState = m.dataSource.GetTaskInputState()

	taskActivePtr := m.dataSource.GetTaskRecordActive()
	if taskActivePtr != prevTaskActivePtr {
		prevTaskActivePtr = taskActivePtr

		if taskActivePtr != nil {
			m.dataSource.InputTaskDescription = taskActivePtr.Description
			m.dataSource.InputTaskTicket = taskActivePtr.Ticket
		}

		desc := m.dataSource.InputTaskDescription
		task := m.dataSource.InputTaskTicket

		*m.nowDoing.TextInputDescription.StringPtr = desc
		m.nowDoing.TextInputDescription.TextInput.SetValue(desc)

		*m.nowDoing.TextInputTask.StringPtr = task
		m.nowDoing.TextInputTask.TextInput.SetValue(task)
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
		m.tutorial.Width = sectionWidth
		m.todaysProgress.Width = sectionWidth
		m.timetable.Width = sectionWidth
		m.embeddedMarkdown.Width = layout.GetEditorWidth(msg.Width)

	case tea.MouseClickMsg:
		switch true {

		case mouse.IsLeftClickInBoundsOf(m.nowDoing.Zone, msg):
			cmds = append(cmds, globals.FocusSection(m.nowDoing.Shortcut))

		case mouse.IsLeftClickInBoundsOf(m.tasksFavorite.Zone, msg):
			cmds = append(cmds, globals.FocusSection(m.tasksFavorite.Shortcut))

		case mouse.IsLeftClickInBoundsOf(m.tasksRecent.Zone, msg):
			cmds = append(cmds, globals.FocusSection(m.tasksRecent.Shortcut))

		case mouse.IsLeftClickInBoundsOf(m.timetable.Zone, msg):
			cmds = append(cmds, globals.FocusSection(m.timetable.Shortcut))

		case mouse.IsLeftClickInBoundsOf(m.embeddedMarkdown.Zone, msg):
			cmds = append(cmds, globals.FocusSection(m.embeddedMarkdown.Shortcut))

		}

	case tea.KeyPressMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "1", "2", "3", "4", "5":
			if m.AreTextInputsBlurred() {
				return m, globals.FocusSection(msg.String())
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
				if isDown && m.AreTextInputsBlurred() {
					cmds = append(cmds, globals.FocusSection(m.tasksRecent.Shortcut))
					return m, tea.Batch(cmds...)
				}

			case m.tasksRecent.IsFocused:
				if isUp {
					isWithinSection := m.dataSource.DecrementOrResetTasksRecentCursor()
					if !isWithinSection {
						cmds = append(cmds, globals.FocusSection(m.nowDoing.Shortcut))
						return m, tea.Batch(cmds...)
					}
				} else if isDown {
					isWithinSection := m.dataSource.IncrementOrResetTasksRecentCursor()
					if !isWithinSection {
						cmds = append(cmds, globals.FocusSection(m.tasksFavorite.Shortcut))
						return m, tea.Batch(cmds...)
					}
				}

			case m.tasksFavorite.IsFocused:
				if isUp {
					isWithinSection := m.dataSource.DecrementOrResetTasksFavoriteCursor()
					if !isWithinSection {
						cmds = append(cmds, globals.FocusSection(m.tasksRecent.Shortcut))
						return m, tea.Batch(cmds...)
					}
				} else if isDown {
					isWithinSection := m.dataSource.IncrementOrResetTasksFavoriteCursor()
					if !isWithinSection {
						cmds = append(cmds, globals.FocusSection(m.timetable.Shortcut))
						return m, tea.Batch(cmds...)
					}
				}

			case m.timetable.IsFocused:
				if isTabKey || isLKey {
					m.dataSource.TimeTableDaysOffsetFromToday++
					return m, tea.Batch(cmds...)
				} else if isShiftTabKey || isHKey {
					m.dataSource.TimeTableDaysOffsetFromToday--
					return m, tea.Batch(cmds...)
				} else if isUp {
					cmds = append(cmds, globals.FocusSection(m.tasksFavorite.Shortcut))
					return m, tea.Batch(cmds...)
				}
			}

		case "enter", "space":
			isEnterKey := msg.String() == "enter"

			switch true {

			case m.nowDoing.IsFocused:
				if isEnterKey {
					taskDefinition := m.dataSource.GetInputTaskDefinition()
					if taskDefinition != nil {
						wasActive := m.dataSource.IsTaskStarted(taskDefinition)
						m.dataSource.TaskToggle(taskDefinition)
						if wasActive {
							cmds = append(cmds, globals.FocusSection(m.nowDoing.Shortcut))
						} else {
							m.nowDoing.BlurOnlyInputs()
						}
						return m, tea.Batch(cmds...)
					} else if m.dataSource.IsInputTaskValidAndNew() {
						m.dataSource.TaskCreate()
						m.nowDoing.BlurOnlyInputs()
						return m, tea.Batch(cmds...)
					}
				}

			case m.tasksRecent.IsFocused:
				taskDefinitions := m.dataSource.GetTaskDefinitionsRecent(true)
				length := len(taskDefinitions)
				if length == 0 {
					return m, tea.Batch(cmds...)
				}

				taskDefinition := taskDefinitions[m.dataSource.TasksRecentCursor]
				m.dataSource.TaskToggle(&taskDefinition)
				return m, tea.Batch(cmds...)

			case m.tasksFavorite.IsFocused:
				taskDefinitions := m.dataSource.GetTaskDefinitionsFavorite()
				length := len(taskDefinitions)
				if length == 0 {
					return m, tea.Batch(cmds...)
				}

				taskDefinition := taskDefinitions[m.dataSource.TasksFavoriteCursor]
				m.dataSource.TaskToggle(&taskDefinition)
				return m, tea.Batch(cmds...)
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

func (m BaseModel) AreTextInputsBlurred() bool {
	return !m.nowDoing.TextInputDescription.TextInput.Focused() &&
		!m.nowDoing.TextInputTask.TextInput.Focused() &&
		!m.embeddedMarkdown.Textarea.Focused()
}

func (m BaseModel) getTextInputValues() (string, string) {
	taskActivePtr := m.dataSource.GetTaskRecordActive()

	if taskActivePtr == nil {
		return "", ""
	}

	return taskActivePtr.Description, taskActivePtr.Ticket
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
				DataSource: dataSource,
				Zone:       globals.TextInputDescriptionZone,
				TextInput:  textinput.New(),
				StringPtr:  &dataSource.InputTaskDescription,
			},
			TextInputTask: &TextInputWrapper.Model{
				DataSource: dataSource,
				Zone:       globals.TextInputTaskZone,
				TextInput:  textinput.New(),
				StringPtr:  &dataSource.InputTaskTicket,
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

	p := tea.NewProgram(&app)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
