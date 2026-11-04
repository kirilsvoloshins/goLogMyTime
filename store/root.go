package store

import (
	"goLogMyTime/common/globals"
	"goLogMyTime/common/styles"
	"goLogMyTime/types"
	"goLogMyTime/utils"
	"time"
)

type RootStore struct {
	Task *TaskStore
	Ui   *UiStore
}

func New(settings *types.Settings) *RootStore {
	rootStore := &RootStore{}
	rootStore.Task = &TaskStore{RootStore: rootStore}
	rootStore.Ui = &UiStore{ThemeIndex: settings.ThemeIndex, FocusedZone: globals.ZoneUnfocused}

	rootStore.Sync()
	return rootStore
}

func (s *RootStore) Sync() {
	s.Ui.TodaysDate = time.Now()
	s.Task.TaskDefinitionMap = utils.ReadTaskDefinitionMap()
	s.Task.TaskRecordsByDateMap = utils.ReadTaskRecordsByDateMap()
	s.Task.TaskRecordActive = s.Task.GetTaskRecordActive()
	s.Task.TaskDefinitionActive = s.Task.GetTaskDefinitionFromRecord(s.Task.TaskRecordActive)
	s.Task.TaskDefinitionsFavorite = s.Task.GetTaskDefinitionsFavorite()
	s.Task.TaskDefinitionsRecent = s.Task.GetTaskDefinitionsRecent()
	todaysDateString := utils.DateToString(s.Ui.TodaysDate)
	if todaysDateString != s.Ui.TodaysDateString {
		// FIXME: crashes the app at midnight
		// s.TasksStopAtDate(s.TodaysDateString)
		s.Ui.TodaysDateString = todaysDateString
	}
	s.Ui.Style = styles.GetStyle(s.Ui.ThemeIndex)
}
