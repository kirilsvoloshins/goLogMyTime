package globals

import (
	"goLogMyTime/types"
	"os"

	tea "charm.land/bubbletea/v2"
)

const READ_WRITE_PERM os.FileMode = 0644
const DATA_DIR = "./data/"
const MARKDOWN_PATH = DATA_DIR + "main.md"

var WORK_DAY_IN_MINUTES float32 = 8 * 60

const (
	ZoneUnfocused       = "none"
	ZoneTextInputTitle  = "nowdoing-input-title"
	ZoneTextInputTask   = "nowdoing-input-task"
	ZoneCreateStartStop = "nowdoing-create-start-stop"
	ZonePrevDate        = "timetable-dayBefore"
	ZoneNextDate        = "timetable-dayAfter"
)

type InteractiveSection struct {
	Shortcut types.SectionShortcut
	Zone     types.Zone
}

type SectionsType struct {
	NowDoing      InteractiveSection
	TasksRecent   InteractiveSection
	TasksFavorite InteractiveSection
	TimeTable     InteractiveSection
	EmbeddedFile  InteractiveSection
}

var InteractiveSections = SectionsType{
	NowDoing: InteractiveSection{
		Shortcut: "1",
		Zone:     "now-doing",
	},
	TasksRecent: InteractiveSection{
		Shortcut: "2",
		Zone:     "tasks-recent",
	},
	TasksFavorite: InteractiveSection{
		Shortcut: "3",
		Zone:     "tasks-favorite",
	},
	TimeTable: InteractiveSection{
		Shortcut: "4",
		Zone:     "time-table",
	},
	EmbeddedFile: InteractiveSection{
		Shortcut: "5",
		Zone:     "embeddedFile",
	},
}

func allSections() []InteractiveSection {
	return []InteractiveSection{
		InteractiveSections.NowDoing,
		InteractiveSections.TasksRecent,
		InteractiveSections.TasksFavorite,
		InteractiveSections.TimeTable,
		InteractiveSections.EmbeddedFile,
	}
}

var ALL_SECTION_ZONES = func() []types.Zone {
	var zones []types.Zone
	for _, s := range allSections() {
		zones = append(zones, s.Zone)
	}
	return zones
}()

var ALL_ZONES = func() []types.Zone {
	return append(
		[]types.Zone{
			ZoneTextInputTitle,
			ZoneTextInputTask,
		},
		ALL_SECTION_ZONES...,
	)
}()

func ShortcutToZone(shortcut types.SectionShortcut) types.Zone {
	for _, s := range allSections() {
		if s.Shortcut == shortcut {
			return s.Zone
		}
	}
	return ZoneUnfocused
}

func FocusZone(zoneToFocus types.Zone) tea.Cmd {
	var cmdsFocus []tea.Cmd
	var cmdsBlur []tea.Cmd

	if zoneToFocus == ZoneUnfocused {
		cmdsFocus = append(cmdsFocus, focusZoneCmd(ZoneUnfocused))

		for _, sectionZone := range ALL_ZONES {
			cmdsBlur = append(cmdsBlur, blurZoneCmd(sectionZone))
		}
	} else {
		shouldKeepNowDoingFocused := zoneToFocus == ZoneTextInputTask || zoneToFocus == ZoneTextInputTitle || zoneToFocus == ZoneCreateStartStop

		for _, sectionZone := range ALL_ZONES {
			if sectionZone == zoneToFocus {
				cmdsFocus = append(cmdsFocus, focusZoneCmd(sectionZone))
			} else if !(sectionZone == InteractiveSections.NowDoing.Zone && shouldKeepNowDoingFocused) {
				cmdsBlur = append(cmdsBlur, blurZoneCmd(sectionZone))
			}
		}
	}

	return tea.Batch(append(cmdsFocus, cmdsBlur...)...)
}

func IsZoneOfAnyTextInput(zone types.Zone) bool {
	return zone == ZoneTextInputTitle ||
		zone == ZoneTextInputTask ||
		zone == InteractiveSections.EmbeddedFile.Zone
}

func focusZoneCmd(zone types.Zone) tea.Cmd {
	return func() tea.Msg {
		return types.FocusZoneMsg(zone)
	}
}

func blurZoneCmd(zone types.Zone) tea.Cmd {
	return func() tea.Msg {
		return types.BlurZoneMsg(zone)
	}
}

const (
	TASK_GOT_STARTED      types.TaskStateChange = "TASK_GOT_STARTED"
	TASK_GOT_STOPPED      types.TaskStateChange = "TASK_GOT_STOPPED"
	TASK_NOTHING_HAPPENED types.TaskStateChange = "TASK_NOTHING_HAPPENED"
)
