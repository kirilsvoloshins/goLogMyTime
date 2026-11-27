package store

import (
	"goLogMyTime/common/colors"
	"goLogMyTime/common/globals"
	"goLogMyTime/common/styles"
	"goLogMyTime/types"
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

	themeIndexNext := s.ThemeIndex + step
	if themeIndexNext >= 0 && themeIndexNext < len(colors.Palettes) {
		s.ThemeIndex = themeIndexNext
		s.Style = styles.GetStyle(s.ThemeIndex)
	}
}

func (s *UiStore) IncrementCursor(cursor *int, numTasks int) bool {
	if numTasks == 0 || *cursor == numTasks-1 {
		return false
	}

	*cursor++
	return true
}

func (s *UiStore) DecrementCursor(cursor *int, numTasks int) bool {
	if numTasks == 0 || *cursor == 0 {
		return false
	}

	*cursor--
	return true
}

func (s *UiStore) ChangeDaysOffset(increment bool) {
	if increment {
		s.TimeTableDaysOffsetFromToday++
	} else {
		s.TimeTableDaysOffsetFromToday--
	}
}
