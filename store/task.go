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
	TaskDefinitionsFavorite types.TaskDefinitions
	TaskDefinitionsRecent   types.TaskDefinitions
}

func (s *TaskStore) DoesTaskExist(ticket types.Ticket) (*types.TaskDefinition, bool) {
	if len(ticket) == 0 {
		return nil, false
	}

	taskDef, exists := s.GetTaskDefinition(ticket)
	return &taskDef, exists
}

func (s *TaskStore) GetTaskDefinition(ticket types.Ticket) (types.TaskDefinition, bool) {
	taskDef, exists := s.TaskDefinitionMap[ticket]
	return taskDef, exists
}

func (s *TaskStore) GetTaskDefinitionFromRecord(taskRecord *types.TaskRecord) *types.TaskDefinition {
	if taskRecord == nil {
		return nil
	}

	taskDef, _ := s.GetTaskDefinition(taskRecord.Ticket)
	return &taskDef
}

func (s *TaskStore) IsTaskValidAndNew(ticket types.Ticket, title string) bool {
	if len(ticket) == 0 || len(title) == 0 {
		return false
	}

	_, exists := s.DoesTaskExist(ticket)
	return !exists
}

func (s *TaskStore) IsTaskStarted(ticket string) bool {
	if s.TaskRecordActive == nil {
		return false
	}

	return ticket == s.TaskRecordActive.Ticket
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
	taskDef, exists := s.GetTaskDefinition(ticket)
	if exists {
		return s.TaskToggle(taskDef.Ticket)
	}

	if s.IsTaskValidAndNew(ticket, title) {
		s.TaskCreateAndStart(ticket, title)
		return globals.TASK_GOT_STARTED
	}

	return globals.TASK_NOTHING_HAPPENED
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

func (s *TaskStore) IsTaskFavorite(ticket types.Ticket) bool {
	taskDef, exists := s.GetTaskDefinition(ticket)
	return exists && taskDef.IsFavorite
}

func (s *TaskStore) TaskFavoriteToggle(ticket types.Ticket) {
	taskDef, exists := s.GetTaskDefinition(ticket)
	utils.Assert(exists, "does not exist")
	taskDef.IsFavorite = !taskDef.IsFavorite
	s.TaskDefinitionMap[ticket] = taskDef
	utils.WriteTaskDefinitionMap(s.TaskDefinitionMap)
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

		if taskRecord.TimestampOfEnd == 0 {
			end = utils.GetTimestamp()
		}

		timetableMap[taskRecord.Ticket] += end - taskRecord.TimestampOfStart
	}

	var timetable types.TimeTable

	// TODO: better way of converting map to array
	for key, valueMs := range timetableMap {
		timetable = append(timetable, types.TimeEntry{
			Ticket:  key,
			Minutes: int(valueMs / 1000 / 60),
		})
	}

	sort.Slice(timetable, func(i, j int) bool {
		return timetable[i].Minutes > timetable[j].Minutes
	})

	return &timetable
}

func (s *TaskStore) GetTodaysProgress(todaysDate time.Time) float32 {
	timetable := s.GetTimeTable(utils.DateToString(todaysDate))
	if timetable == nil {
		return 0
	}

	var workTimeMin int
	for _, entry := range *timetable {
		workTimeMin += entry.Minutes
	}

	return float32(workTimeMin / globals.WORK_DAY_IN_MINUTES)
}

func (s *TaskStore) GetDateAndTimeTable(todaysDate time.Time, offsetFromToday int) (string, *types.TimeTable) {
	date1 := todaysDate.AddDate(0, 0, offsetFromToday)
	dateString := utils.DateToString(date1)
	timeTable := s.GetTimeTable(dateString)
	return dateString, timeTable
}

// Sync time only below
func (s *TaskStore) GetTaskDefinitionsFavorite() types.TaskDefinitions {
	return utils.MapPointerToArrayFavorite(&s.TaskDefinitionMap)
}

func (s *TaskStore) GetTaskDefinitionsRecent() types.TaskDefinitions {
	taskDefsRecent := utils.MapPointerToArrayRecent(&s.TaskDefinitionMap)
	if s.TaskRecordActive == nil {
		return taskDefsRecent
	}

	var taskDefs types.TaskDefinitions
	taskDefActive, exists := s.GetTaskDefinition(s.TaskRecordActive.Ticket)
	for _, taskDef := range taskDefsRecent {
		if !exists || taskDef.Ticket != taskDefActive.Ticket {
			taskDefs = append(taskDefs, taskDef)
		}
	}
	return taskDefs
}

func (s *TaskStore) GetTaskRecordActive() *types.TaskRecord {
	tasksToday := s.TaskRecordsByDateMap[utils.TodaysString()]
	return utils.GetTaskRecordActive(&tasksToday)
}
