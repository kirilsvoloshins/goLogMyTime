package store

import (
	"goLogMyTime/common/globals"
	"goLogMyTime/types"
	"goLogMyTime/utils"
	"sort"
	"time"
)

type TaskStore struct {
	RootStore               *RootStore
	TaskDefinitionMap       types.TaskDefinitionMap
	TaskRecordsByDateMap    types.TaskRecordsByDateMap
	TaskRecordActive        *types.TaskRecord
	TaskDefinitionActive    *types.TaskDefinition
	TaskDefinitionsFavorite types.TaskDefinitionsFavorite
	TaskDefinitionsRecent   types.TaskDefinitionsRecent
}

func (s *TaskStore) IsTaskValidAndNew(ticket types.Ticket, title string) bool {
	if len(ticket) == 0 || len(title) == 0 {
		return false
	}

	_, exists := s.DoesTaskExist(ticket)
	return !exists
}

func (s *TaskStore) DoesTaskExist(ticket types.Ticket) (*types.TaskDefinition, bool) {
	if len(ticket) == 0 {
		return nil, false
	}

	taskDefinition, existing := s.TaskDefinitionMap[ticket]
	return &taskDefinition, existing
}

func (s *TaskStore) IsTaskStarted(ticket string) bool {
	if s.TaskRecordActive == nil {
		return false
	}

	return ticket == s.TaskRecordActive.Ticket
}

func (s *TaskStore) GetTaskDefinition(ticket types.Ticket) *types.TaskDefinition {
	taskDef, exists := s.TaskDefinitionMap[ticket]
	if exists {
		return &taskDef
	}

	return nil
}

func (s *TaskStore) GetTaskDefinitionFromRecord(taskRecord *types.TaskRecord) *types.TaskDefinition {
	if taskRecord == nil {
		return nil
	}

	return s.GetTaskDefinition(taskRecord.Ticket)
}

func (s *TaskStore) TaskCreateAndStart(ticket string, title string) {
	utils.Assert(s.IsTaskValidAndNew(ticket, title), "New task must have valid and unique 'ticket'")

	s.TaskDefinitionMap[ticket] = types.TaskDefinition{
		Ticket:                  ticket,
		Info:                    title,
		IsFavorite:              false,
		TimestampOfLastActivity: 0,
	}
	utils.WriteTaskDefinitionMap(s.TaskDefinitionMap)

	s.TaskStart(ticket)
}

func (s *TaskStore) TaskToggle(ticket string) types.TaskStateChange {
	if s.IsTaskStarted(ticket) {
		s.TasksStopAtDate(utils.TodaysString())
		return globals.TASK_GOT_STOPPED
	}

	s.TaskStart(ticket)
	return globals.TASK_GOT_STARTED
}

func (s *TaskStore) TaskToggleIfValid(ticket string, title string) types.TaskStateChange {
	taskDef := s.GetTaskDefinition(ticket)
	if taskDef == nil {
		if s.IsTaskValidAndNew(ticket, title) {
			s.TaskCreateAndStart(ticket, title)
			return globals.TASK_GOT_STARTED
		}

		return globals.TASK_NOTHING_HAPPENED
	}

	return s.TaskToggle(taskDef.Ticket)
}

func (s *TaskStore) TaskStart(ticket string) {
	todaysString := utils.TodaysString()
	s.TasksStopAtDate(todaysString)

	var taskRecordsAtDate, ok = s.TaskRecordsByDateMap[todaysString]
	if !ok {
		// TODO: find a better way :)
		s.TaskRecordsByDateMap[todaysString] = []types.TaskRecord{}
		taskRecordsAtDate, ok = s.TaskRecordsByDateMap[todaysString]
		utils.Assert(ok, "come on, should have worked now")
	}

	taskDefinition, ok := s.TaskDefinitionMap[ticket]
	utils.Assert(ok, "missing task def")

	s.TaskRecordsByDateMap[todaysString] = append(
		taskRecordsAtDate,
		types.TaskRecord{
			Ticket:           taskDefinition.Ticket,
			Description:      taskDefinition.Info,
			TimestampOfStart: utils.GetTimestamp(),
			IsActive:         true,
			TimestampOfEnd:   0,
		},
	)

	utils.WriteTaskRecordsByDateMap(s.TaskRecordsByDateMap)
	utils.WriteTimestampOfLastActivity(s.TaskDefinitionMap, ticket)

	// s.InputTicket = taskDefinition.Ticket
	// s.InputDescription = taskDefinition.Info
	s.RootStore.Sync()
}

func (s *TaskStore) TasksStopAtDate(dateString string) {
	taskRecordsAtDate, ok := s.TaskRecordsByDateMap[dateString]
	if !ok {
		return
	}

	// TODO: make a util for stopping a record...
	for i := range taskRecordsAtDate {
		if taskRecordsAtDate[i].IsActive {
			taskRecordsAtDate[i].IsActive = false
			taskRecordsAtDate[i].TimestampOfEnd = utils.GetTimestamp()
			taskDefinition, ok := s.TaskDefinitionMap[taskRecordsAtDate[i].Ticket]
			utils.Assert(ok, "failed to map task record to its definition")

			utils.WriteTimestampOfLastActivity(s.TaskDefinitionMap, taskDefinition.Ticket)
		}
	}

	utils.WriteTaskRecordsByDateMap(s.TaskRecordsByDateMap)

	// s.InputTicket = ""
	// s.InputDescription = ""
	s.RootStore.Sync()
}

func (s *TaskStore) GetTimeTable(dateString string) *types.TimeTable {
	allRecordsAtDate, dateHasRecords := s.TaskRecordsByDateMap[dateString]
	if !dateHasRecords || len(allRecordsAtDate) == 0 {
		return nil
	}

	timetableMap := types.TimeTableMap{}

	for _, taskRecord := range allRecordsAtDate {
		end := taskRecord.TimestampOfEnd

		// task in progress
		if taskRecord.TimestampOfEnd == 0 {
			end = utils.GetTimestamp()
		}

		durationMs := end - taskRecord.TimestampOfStart

		_, alreadyInMap := timetableMap[taskRecord.Ticket]
		if alreadyInMap {
			timetableMap[taskRecord.Ticket] += durationMs
		} else {
			timetableMap[taskRecord.Ticket] = durationMs
		}
	}

	var timetable types.TimeTable

	// TODO: better way of converting map to array
	for key, valueMs := range timetableMap {
		valueMinutes := int(valueMs / 1000 / 60)
		timetable = append(timetable, types.TimeEntry{
			Name:    key,
			Minutes: valueMinutes,
		})
	}

	sort.Slice(timetable, func(i, j int) bool {
		return timetable[i].Minutes > timetable[j].Minutes
	})

	return &timetable
}

func (s *TaskStore) GetTodaysWorkTimeMin(todaysDate time.Time) int {
	timetable := s.GetTimeTable(utils.DateToString(todaysDate))
	if timetable == nil {
		return 0
	}

	var workTimeMin int
	for _, entry := range *timetable {
		workTimeMin += entry.Minutes
	}

	return workTimeMin
}

func (s *TaskStore) GetDateAndTimeTable(todaysDate time.Time, offsetFromToday int) (string, *types.TimeTable) {
	date1 := todaysDate.AddDate(0, 0, offsetFromToday)
	dateString := utils.DateToString(date1)
	timeTable := s.GetTimeTable(dateString)
	return dateString, timeTable
}

// Sync time only below
func (s *TaskStore) GetTaskDefinitionsFavorite() types.TaskDefinitionsFavorite {
	return utils.MapPointerToArrayFavorite(&s.TaskDefinitionMap)
}

func (s *TaskStore) GetTaskDefinitionsRecent() types.TaskDefinitionsRecent {
	taskDefsRecent := utils.MapPointerToArrayRecent(&s.TaskDefinitionMap)
	if s.TaskRecordActive == nil {
		return taskDefsRecent
	}

	var taskDefs types.TaskDefinitionsRecent
	taskDefActive := s.GetTaskDefinition(s.TaskRecordActive.Ticket)
	for _, taskDef := range taskDefsRecent {
		if taskDef.Ticket != taskDefActive.Ticket {
			taskDefs = append(taskDefs, taskDef)
		}
	}
	return taskDefs
}

func (s *TaskStore) GetTaskRecordActive() *types.TaskRecord {
	tasksToday := s.TaskRecordsByDateMap[utils.TodaysString()]
	return utils.GetTaskRecordActive(&tasksToday)
}
