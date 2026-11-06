package main

import (
	"fmt"
	"os"
	"time"

	"goLogMyTime/common/globals"
	"goLogMyTime/common/styles"
	ButtonCreateStartStop "goLogMyTime/components/ButtonCreateStartStopTask"
	ButtonStartStop "goLogMyTime/components/ButtonStartStopTask"
	"goLogMyTime/components/EmbeddedMarkdown"
	"goLogMyTime/components/HintsContextAware"
	"goLogMyTime/components/InputWrapper"
	"goLogMyTime/components/NowDoing"
	"goLogMyTime/components/TasksFavorite"
	"goLogMyTime/components/TasksRecent"
	"goLogMyTime/components/TimeTable"
	"goLogMyTime/components/TodaysProgress"
	"goLogMyTime/store"
	"goLogMyTime/types"
	"goLogMyTime/utils"

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
	store             *store.RootStore
	hintsContextAware *HintsContextAware.Model
	nowDoing          *NowDoing.Model
	todaysProgress    *TodaysProgress.Model
	tasksRecent       *TasksRecent.Model
	tasksFavorite     *TasksFavorite.Model
	timetable         *TimeTable.Model
	embeddedMarkdown  *EmbeddedMarkdown.Model
}

func (m BaseModel) Init() tea.Cmd {
	return tickToSync()
}

// var prevTaskActivePtr *types.TaskRecord

func (m *BaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	// m.store.TaskInputState = m.store.Task.GetTaskInputState()

	// taskActivePtr := m.store.Task.GetTaskRecordActive()
	// if taskActivePtr != prevTaskActivePtr {
	// 	prevTaskActivePtr = taskActivePtr

	// 	if taskActivePtr != nil {
	// 		m.store.InputDescription = taskActivePtr.Description
	// 		m.store.InputTicket = taskActivePtr.Ticket
	// 	}

	// 	desc := m.store.InputDescription
	// 	task := m.store.InputTicket

	// 	*m.nowDoing.InputDescription.StringPtr = desc
	// 	m.nowDoing.InputDescription.TextInput.SetValue(desc)

	// 	*m.nowDoing.InputTask.StringPtr = tas
	// 	m.nowDoing.InputTask.TextInput.SetValue(task)
	// }

	switch msg := msg.(type) {
	case types.FocusZoneMsg:
		zoneToFocus := types.Zone(msg)
		m.store.Ui.FocusedZone = zoneToFocus

	case tickToSyncMsg:
		m.store.Sync()
		cmds = append(cmds, tickToSync())

	case tea.WindowSizeMsg:
		sectionWidth := styles.Sizes.SectionWidth
		m.nowDoing.Width = sectionWidth
		m.tasksFavorite.Width = sectionWidth
		m.tasksRecent.Width = sectionWidth
		m.todaysProgress.Width = sectionWidth
		m.timetable.Width = sectionWidth
		m.embeddedMarkdown.Width = msg.Width - sectionWidth
		m.hintsContextAware.Width = msg.Width

	case tea.MouseClickMsg:
		found := false
		for _, zone := range globals.ALL_ZONES {
			if utils.IsLeftClickInBoundsOf(zone, msg) {
				cmds = append(cmds, globals.FocusZone(zone))
				found = true
				break
			}
		}

		if !found {
			cmds = append(cmds, globals.FocusZone(globals.ZoneUnfocused))
		}

	case tea.KeyPressMsg:
		keyString := msg.String()

		if keyString == "ctrl+c" {
			return m, tea.Quit
		}

		if m.IsAnyTextInputFocused() {
			break
		}

		switch keyString {
		case "/":
			cmds = append(cmds, globals.FocusZone(globals.ZoneTextInputTask))

		case "i":
			cmds = append(cmds, globals.FocusZone(globals.InteractiveSections.EmbeddedFile.Zone))

		case "t", "T":
			next := keyString == "T"
			m.store.Ui.ChangeThemeIndex(next)

		case "1", "2", "3", "4", "5":
			cmds = append(cmds, globals.FocusZone(globals.ShortcutToZone(types.SectionShortcut(keyString))))

		case "j", "k", "l", "down", "up", "tab", "shift+tab":
			if utils.IsRightMovementKeyMsg(msg) {
				cmds = append(cmds, globals.FocusZone(globals.InteractiveSections.EmbeddedFile.Zone))
				break
			}

			isDown := utils.IsDownMovementKeyMsg(msg)
			isUp := utils.IsUpMovementKeyMsg(msg)

			switch true {

			case m.nowDoing.IsFocused:
				if isDown {
					cmds = append(cmds, globals.FocusZone(getNextSectionZone(m.nowDoing)))
				}

			case m.tasksRecent.IsFocused:
				numTasks := len(m.store.Task.TaskDefinitionsRecent)

				if isUp && !m.store.Ui.DecrementCursor(&m.store.Ui.TasksRecentCursor, numTasks) {
					cmds = append(cmds, globals.FocusZone(getPrevSectionZone(m.tasksRecent)))
				}

				if isDown && !m.store.Ui.IncrementCursor(&m.store.Ui.TasksRecentCursor, numTasks) {
					cmds = append(cmds, globals.FocusZone(getNextSectionZone(m.tasksRecent)))
				}

			case m.tasksFavorite.IsFocused:
				numTasks := len(m.store.Task.TaskDefinitionsFavorite)

				if isUp && !m.store.Ui.DecrementCursor(&m.store.Ui.TasksFavoriteCursor, numTasks) {
					cmds = append(cmds, globals.FocusZone(getPrevSectionZone(m.tasksFavorite)))
				}

				if isDown && !m.store.Ui.IncrementCursor(&m.store.Ui.TasksFavoriteCursor, numTasks) {
					cmds = append(cmds, globals.FocusZone(getNextSectionZone(m.tasksFavorite)))
				}

			case m.timetable.IsFocused:
				if isUp {
					cmds = append(cmds, globals.FocusZone(getPrevSectionZone(m.timetable)))
				}
			}
		}

	}

	m.hintsContextAware, cmd = m.hintsContextAware.Update(msg)
	cmds = append(cmds, cmd)

	m.nowDoing, cmd = m.nowDoing.Update(msg)
	cmds = append(cmds, cmd)

	m.tasksFavorite, cmd = m.tasksFavorite.Update(msg)
	cmds = append(cmds, cmd)

	m.tasksRecent, cmd = m.tasksRecent.Update(msg)
	cmds = append(cmds, cmd)

	m.todaysProgress, cmd = m.todaysProgress.Update(msg)
	cmds = append(cmds, cmd)

	m.timetable, cmd = m.timetable.Update(msg)
	cmds = append(cmds, cmd)

	m.embeddedMarkdown, cmd = m.embeddedMarkdown.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m BaseModel) View() tea.View {
	var view tea.View
	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion
	view.BackgroundColor = m.store.Ui.Style.BgColor

	view.SetContent(
		zone.Scan(
			lipgloss.JoinVertical(
				lipgloss.Top,
				m.hintsContextAware.View().Content,
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					lipgloss.JoinVertical(
						lipgloss.Right,
						m.todaysProgress.View().Content,
						m.nowDoing.View().Content,
						m.tasksRecent.View().Content,
						m.tasksFavorite.View().Content,
						m.timetable.View().Content,
					),
					m.embeddedMarkdown.View().Content,
				),
			),
		),
	)

	return view
}

func (m BaseModel) IsAnyTextInputFocused() bool {
	return m.nowDoing.TextInputTitle.TextInput.Focused() ||
		m.nowDoing.TextInputTask.TextInput.Focused() ||
		m.embeddedMarkdown.Textarea.Focused()
}

func main() {
	utils.ReadDataDir()
	utils.CreateTasksFileIfMissing(utils.TASK_DEFINITIONS_PATH)
	utils.CreateTasksFileIfMissing(utils.TASK_RECORDS_PATH)
	utils.CreateSettingsIfMissing()
	settings := utils.ReadSettings()

	rootStore := store.New(settings)

	zone.NewGlobal()

	embeddedMarkdownModel := EmbeddedMarkdown.InitialModel(utils.MARKDOWN_PATH, rootStore)

	app := BaseModel{
		store:             rootStore,
		hintsContextAware: &HintsContextAware.Model{Store: rootStore},
		tasksRecent: &TasksRecent.Model{
			Store:               rootStore,
			Zone:                globals.InteractiveSections.TasksRecent.Zone,
			Shortcut:            globals.InteractiveSections.TasksRecent.Shortcut,
			ButtonsStartStopMap: map[string]ButtonStartStop.Model{},
		},
		tasksFavorite: &TasksFavorite.Model{
			Store:               rootStore,
			Zone:                globals.InteractiveSections.TasksFavorite.Zone,
			Shortcut:            globals.InteractiveSections.TasksFavorite.Shortcut,
			ButtonsStartStopMap: map[string]ButtonStartStop.Model{},
		},
		nowDoing: &NowDoing.Model{
			Store:                 rootStore,
			Zone:                  globals.InteractiveSections.NowDoing.Zone,
			Shortcut:              globals.InteractiveSections.NowDoing.Shortcut,
			ButtonCreateStartStop: &ButtonCreateStartStop.Model{Store: rootStore},
			TextInputTitle: &InputWrapper.Model{
				Store:     rootStore,
				Zone:      globals.ZoneTextInputTitle,
				TextInput: textinput.New(),
				StringPtr: &rootStore.Ui.InputTitle,
			},
			TextInputTask: &InputWrapper.Model{
				Store:     rootStore,
				Zone:      globals.ZoneTextInputTask,
				TextInput: textinput.New(),
				StringPtr: &rootStore.Ui.InputTicket,
			},
		},
		todaysProgress: &TodaysProgress.Model{Store: rootStore},
		timetable: &TimeTable.Model{
			Store:    rootStore,
			Zone:     globals.InteractiveSections.TimeTable.Zone,
			Shortcut: globals.InteractiveSections.TimeTable.Shortcut,
		},
		embeddedMarkdown: &embeddedMarkdownModel,
	}

	app.nowDoing.TextInputTitle.TextInput.Placeholder = "Something great"
	app.nowDoing.TextInputTitle.TextInput.Prompt = "Title: "
	app.nowDoing.TextInputTitle.TextInput.SetVirtualCursor(true)
	app.nowDoing.TextInputTitle.TextInput.CharLimit = 15
	app.nowDoing.TextInputTitle.TextInput.SetWidth(15)

	app.nowDoing.TextInputTask.TextInput.Placeholder = "DEV-1234"
	app.nowDoing.TextInputTask.TextInput.Prompt = "Task: "
	app.nowDoing.TextInputTask.TextInput.SetVirtualCursor(true)
	app.nowDoing.TextInputTask.TextInput.CharLimit = 15
	app.nowDoing.TextInputTask.TextInput.SetWidth(15)

	p := tea.NewProgram(&app)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func getPrevSectionZone(modelPtr any) types.Zone {
	switch modelPtr.(type) {
	case *TasksRecent.Model:
		return globals.InteractiveSections.NowDoing.Zone
	case *TasksFavorite.Model:
		return globals.InteractiveSections.TasksRecent.Zone
	case *TimeTable.Model:
		return globals.InteractiveSections.TasksFavorite.Zone
	}

	return globals.ZoneUnfocused
}

func getNextSectionZone(modelPtr any) types.Zone {
	switch modelPtr.(type) {
	case *NowDoing.Model:
		return globals.InteractiveSections.TasksRecent.Zone
	case *TasksRecent.Model:
		return globals.InteractiveSections.TasksFavorite.Zone
	case *TasksFavorite.Model:
		return globals.InteractiveSections.TimeTable.Zone
	}

	return globals.ZoneUnfocused
}
