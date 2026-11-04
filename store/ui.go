package store

import (
	"goLogMyTime/common/colors"
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
