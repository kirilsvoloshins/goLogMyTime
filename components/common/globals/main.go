package globals

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
