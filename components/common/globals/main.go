package globals

import tea "charm.land/bubbletea/v2"

type InteractiveSection struct {
	Shortcut string
	Zone     string
}

var NowDoingSection = InteractiveSection{
	Shortcut: "1",
	Zone:     "now-doing",
}

var TasksRecentSection = InteractiveSection{
	Shortcut: "2",
	Zone:     "tasks-recent",
}

var TasksFavoriteSection = InteractiveSection{
	Shortcut: "3",
	Zone:     "tasks-favorite",
}

var TimeTableSection = InteractiveSection{
	Shortcut: "4",
	Zone:     "time-table",
}

var EmbeddedFileSection = InteractiveSection{
	Shortcut: "5",
	Zone:     "embeddedFile",
}

var TextInputDescriptionZone = "description-input"
var TextInputTaskZone = "task-input"
var CreateStartStopZone = "create-start-stop"
var PrevDateZone = "dayBefore"
var NextDateZone = "dayAfter"

const (
	TASK_INPUT_EMPTY = iota
	TASK_INPUT_VALID_NEW
	TASK_INPUT_VALID_INACTIVE
	TASK_INPUT_VALID_ACTIVE
)

var shortcuts = []string{
	"1", "2", "3", "4", "5",
}

type FocusSectionMsg string
type BlurSectionMsg string

func FocusSection(section string) tea.Cmd {
	var cmds []tea.Cmd

	for _, shortcut := range shortcuts {
		if section == shortcut {
			cmds = append(cmds, FocusSectionCmd(shortcut))
		} else {
			cmds = append(cmds, BlurSectionCmd(shortcut))
		}
	}

	return tea.Batch(cmds...)
}

func FocusSectionCmd(section string) tea.Cmd {
	return func() tea.Msg {
		return FocusSectionMsg(section)
	}
}

func BlurSectionCmd(section string) tea.Cmd {
	return func() tea.Msg {
		return BlurSectionMsg(section)
	}
}
