package store

import (
	"goLogMyTime/common/colors"
	"goLogMyTime/common/globals"
	"goLogMyTime/common/styles"
	"goLogMyTime/types"
	"goLogMyTime/utils"
	"time"
)

type UiStore struct {
	ThemeIndex                   int
	Style                        styles.Style
	TodaysDate                   time.Time
	TodaysDateString             string
	TasksRecentCursor            int
	TasksFavoriteCursor          int
	TimeTableDaysOffsetFromToday int
	InputTitle                   string
	InputTicket                  string
	FocusedZone                  types.Zone
}

func (s *UiStore) GetZone(next bool) types.Zone {
	if next {
		switch s.FocusedZone {
		case globals.InteractiveSections.NowDoing.Zone:
			return globals.InteractiveSections.TasksRecent.Zone
		case globals.InteractiveSections.TasksRecent.Zone:
			return globals.InteractiveSections.TasksFavorite.Zone
		case globals.InteractiveSections.TasksFavorite.Zone:
			return globals.InteractiveSections.TimeTable.Zone
		case globals.InteractiveSections.TimeTable.Zone:
			return globals.InteractiveSections.TimeTable.Zone
		}
	}

	switch s.FocusedZone {
	case globals.InteractiveSections.NowDoing.Zone:
		return globals.InteractiveSections.NowDoing.Zone
	case globals.InteractiveSections.TasksRecent.Zone:
		return globals.InteractiveSections.NowDoing.Zone
	case globals.InteractiveSections.TasksFavorite.Zone:
		return globals.InteractiveSections.TasksRecent.Zone
	case globals.InteractiveSections.TimeTable.Zone:
		return globals.InteractiveSections.TasksFavorite.Zone
	}

	return globals.InteractiveSections.NowDoing.Zone
}

func (s *UiStore) ChangeThemeIndex(next bool) {
	step := -1

	if next {
		step = 1
	}

	cursor := utils.ClampCursor(s.ThemeIndex+step, len(colors.Palettes))
	s.ThemeIndex = cursor
	s.Style = styles.GetStyle(s.ThemeIndex)
}

func (s *UiStore) IncrementCursor(pCursor *int, numTasks int) bool {
	if numTasks == 0 {
		return false
	}

	next := utils.ClampCursor(*pCursor+1, numTasks)
	if next <= *pCursor {
		return false
	}

	*pCursor++
	return true
}

func (s *UiStore) DecrementCursor(pCursor *int, numTasks int) bool {
	if numTasks == 0 {
		return false
	}

	next := utils.ClampCursor(*pCursor-1, numTasks)
	if next >= *pCursor {
		return false
	}

	*pCursor--
	return true
}

func (s *UiStore) ChangeDaysOffset(increment bool) {
	if increment {
		s.TimeTableDaysOffsetFromToday++
	} else {
		s.TimeTableDaysOffsetFromToday--
	}
}
