package taskStore

import (
	"goLogMyTime/components/common/globals"
	"goLogMyTime/domain/types"
	"goLogMyTime/utils/assert"
	"goLogMyTime/utils/conversions"
	"goLogMyTime/utils/date"
	"goLogMyTime/utils/file"
	"sort"
	"time"
)

type TaskStore struct {
	TodaysDate                   time.Time
	TodaysDateString             string
	TaskDefinitionMap            types.TaskDefinitionMap
	TaskRecordsByDateMap         types.TaskRecordsByDateMap
	TaskInputState               int
	InputTaskTicket              string
	InputTaskDescription         string
	TasksRecentCursor            int
	TasksFavoriteCursor          int
	TimeTableDaysOffsetFromToday int
}

func New() *TaskStore {
	taskStore := &TaskStore{}
	taskStore.Sync()
	return taskStore
}

func (s *TaskStore) Sync() {
	s.TodaysDate = time.Now()
	s.TaskDefinitionMap = file.ReadTaskDefinitionMap()
	s.TaskRecordsByDateMap = file.ReadTaskRecordsByDateMap()
	todaysDateString := date.DateToString(s.TodaysDate)
	if todaysDateString != s.TodaysDateString {
		s.TasksStopAtDate(s.TodaysDateString)
		s.TodaysDateString = todaysDateString
	}
}

func (s *TaskStore) IsInputTaskExisting() bool {
	ticket := s.InputTaskTicket

	if len(ticket) == 0 {
		return false
	}

	_, existing := s.TaskDefinitionMap[ticket]
	return existing
}

func (s *TaskStore) IsInputTaskValidAndNew() bool {
	ticket := s.InputTaskTicket

	if len(ticket) == 0 {
		return false
	}

	_, existing := s.TaskDefinitionMap[ticket]
	return !existing
}

func (s *TaskStore) GetTaskDefinitionsFavorite() types.TaskDefinitionsFavorite {
	return conversions.MapPointerToArrayFavorite(&s.TaskDefinitionMap)
}

func (s *TaskStore) GetTaskDefinitionsRecent(excludeCurrent bool) types.TaskDefinitionsRecent {
	taskDefinitionsRecent := conversions.MapPointerToArrayRecent(&s.TaskDefinitionMap)
	if !excludeCurrent {
		return taskDefinitionsRecent
	}

	taskRecordActive := s.GetTaskRecordActive()
	if taskRecordActive == nil {
		return taskDefinitionsRecent
	}

	var taskDefinitions types.TaskDefinitionsRecent
	taskDefinitionActive := s.GetTaskDefinitionFromTaskRecord(taskRecordActive)
	for _, taskDefinition := range taskDefinitionsRecent {
		if taskDefinition.Ticket != taskDefinitionActive.Ticket {
			taskDefinitions = append(taskDefinitions, taskDefinition)
		}
	}
	return taskDefinitions
}

// TODO: check if does not crash on empty definition file :)
func (s *TaskStore) GetTaskDefinitionMostRecent() *types.TaskDefinition {
	taskDefinitionsRecent := s.GetTaskDefinitionsRecent(false)
	if len(taskDefinitionsRecent) > 0 {
		return &taskDefinitionsRecent[0]
	}

	return nil
}

func (s *TaskStore) GetTaskDefinitionFromTaskRecord(
	taskRecord *types.TaskRecord,
) *types.TaskDefinition {
	v, ok := s.TaskDefinitionMap[taskRecord.Ticket]
	assert.Assert(ok, "Failed to resolve ticket definition")

	return &v
}

// TODO: what is returned if not ok?
func (s *TaskStore) GetTasksRecordsToday() *types.TaskRecords {
	v := s.TaskRecordsByDateMap[s.TodaysDateString]
	return &v
}

func (s *TaskStore) GetTaskRecordActive() *types.TaskRecord {
	return conversions.GetTaskRecordActive(s.GetTasksRecordsToday())
}

func (s *TaskStore) TaskCreate() {
	ticket := s.InputTaskTicket

	assert.Assert(s.IsInputTaskValidAndNew(), "New task must have valid and unique 'ticket'")

	s.TaskDefinitionMap[ticket] = types.TaskDefinition{
		Ticket:                  ticket,
		Info:                    s.InputTaskDescription,
		IsFavorite:              false,
		TimestampOfLastActivity: 0,
	}
	file.WriteTaskDefinitionMap(s.TaskDefinitionMap)
	s.Sync()

	taskDefinition, ok := s.TaskDefinitionMap[ticket]
	assert.Assert(ok, "Failed to read definition of a new task")

	s.TaskStart(&taskDefinition)
}

func (s *TaskStore) IsTaskStarted(taskDefinition *types.TaskDefinition) bool {
	taskRecordActive := s.GetTaskRecordActive()
	return taskRecordActive != nil && taskDefinition.Ticket == taskRecordActive.Ticket
}

func (s *TaskStore) TaskToggle(taskDefinition *types.TaskDefinition) {
	if s.IsTaskStarted(taskDefinition) {
		s.TasksStopAtDate(s.TodaysDateString)
	} else {
		s.TaskStart(taskDefinition)
	}
}

func (s *TaskStore) TaskStart(taskDefinition *types.TaskDefinition) {
	s.TasksStopAtDate(s.TodaysDateString)

	var taskRecordsAtDate, ok = s.TaskRecordsByDateMap[s.TodaysDateString]
	if !ok {
		// TODO: find a better way :)
		s.TaskRecordsByDateMap[s.TodaysDateString] = []types.TaskRecord{}
		taskRecordsAtDate, ok = s.TaskRecordsByDateMap[s.TodaysDateString]
		assert.Assert(ok, "come on, should have worked now")
	}

	s.TaskRecordsByDateMap[s.TodaysDateString] = append(
		taskRecordsAtDate,
		types.TaskRecord{
			Ticket:           taskDefinition.Ticket,
			Description:      taskDefinition.Info,
			TimestampOfStart: date.GetTimestamp(),
			IsActive:         true,
			TimestampOfEnd:   0,
		},
	)

	file.WriteTaskRecordsByDateMap(s.TaskRecordsByDateMap)
	file.WriteTimestampOfLastActivity(s.TaskDefinitionMap, *taskDefinition)

	s.InputTaskTicket = taskDefinition.Ticket
	s.InputTaskDescription = taskDefinition.Info
	s.TaskInputState = globals.TASK_INPUT_VALID_ACTIVE
	s.Sync()
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
			taskRecordsAtDate[i].TimestampOfEnd = date.GetTimestamp()
			taskDefinition, ok := s.TaskDefinitionMap[taskRecordsAtDate[i].Ticket]
			assert.Assert(ok, "failed to map task record to its definition")

			file.WriteTimestampOfLastActivity(s.TaskDefinitionMap, taskDefinition)
		}
	}

	file.WriteTaskRecordsByDateMap(s.TaskRecordsByDateMap)

	s.InputTaskTicket = ""
	s.InputTaskDescription = ""
	s.TaskInputState = globals.TASK_INPUT_EMPTY
	s.Sync()
}

func (s *TaskStore) IncrementOrResetTasksRecentCursor() bool {
	length := len(s.GetTaskDefinitionsRecent(true))
	if length == 0 || s.TasksRecentCursor == length-1 {
		return false
	}

	s.TasksRecentCursor++
	return true
}

func (s *TaskStore) DecrementOrResetTasksRecentCursor() bool {
	length := len(s.GetTaskDefinitionsRecent(true))
	if length == 0 || s.TasksRecentCursor == 0 {
		return false
	}

	s.TasksRecentCursor--
	return true
}

func (s *TaskStore) IncrementOrResetTasksFavoriteCursor() bool {
	length := len(s.GetTaskDefinitionsFavorite())
	if length == 0 || s.TasksFavoriteCursor == length-1 {
		return false
	}

	s.TasksFavoriteCursor++
	return true
}

func (s *TaskStore) DecrementOrResetTasksFavoriteCursor() bool {
	length := len(s.GetTaskDefinitionsFavorite())
	if length == 0 || s.TasksFavoriteCursor == 0 {
		return false
	}

	s.TasksFavoriteCursor--
	return true
}

func (s *TaskStore) GetDateAndTimeTable() (string, *types.TimeTable) {
	date1 := s.TodaysDate.AddDate(0, 0, s.TimeTableDaysOffsetFromToday)
	dateString := date.DateToString(date1)
	timeTable := s.GetTimeTable(dateString)
	return dateString, timeTable
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
			end = date.GetTimestamp()
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

func (s *TaskStore) GetTodaysWorkTimeMin() int {
	dateBeingChecked := date.DateToString(s.TodaysDate)
	timetable := s.GetTimeTable(dateBeingChecked)
	if timetable == nil {
		return 0
	}

	var workTimeMin int
	for _, entry := range *timetable {
		workTimeMin += entry.Minutes
	}

	return workTimeMin
}

func (s *TaskStore) GetInputTaskDefinition() *types.TaskDefinition {
	activeTaskRecord := s.GetTaskRecordActive()
	if activeTaskRecord != nil {
		taskDefinition, readOk := s.TaskDefinitionMap[activeTaskRecord.Ticket]
		assert.Assert(readOk, "missing definition of an active task")
		return &taskDefinition
	}

	inputTaskDefinition, inputOk := s.TaskDefinitionMap[s.InputTaskTicket]
	if inputOk {
		return &inputTaskDefinition
	}

	return nil
}

func (s *TaskStore) GetTaskInputState() int {
	taskDefinition := s.GetInputTaskDefinition()
	if taskDefinition != nil {
		if s.IsTaskStarted(taskDefinition) {
			return globals.TASK_INPUT_VALID_ACTIVE
		}
		return globals.TASK_INPUT_VALID_INACTIVE
	}

	if s.IsInputTaskValidAndNew() {
		return globals.TASK_INPUT_VALID_NEW
	}

	return globals.TASK_INPUT_EMPTY
}

var WORK_DAY_IN_MINUTES = 8 * 60

func (s *TaskStore) GetTodaysProgreesPercent() float64 {
	var res float64
	todaysWorkTimeMin := s.GetTodaysWorkTimeMin()
	res = float64((100 * todaysWorkTimeMin) / WORK_DAY_IN_MINUTES)
	return res
}
