package store

import (
	"encoding/json"
	"errors"
	"goLogMyTime/common/globals"
	"goLogMyTime/types"
	"goLogMyTime/utils"
	"io/fs"
	"os"
)

const TASK_DEFINITIONS_PATH = globals.DATA_DIR + "taskDefinitions.json"
const TASK_RECORDS_PATH = globals.DATA_DIR + "taskRecords.json"
const SETTINGS_PATH = globals.DATA_DIR + "settings.json"

type FileStore struct{}

func (s *FileStore) ReadSettings() *types.Settings {
	bytes, err := os.ReadFile(SETTINGS_PATH)
	utils.Assert(err == nil, "CreateAndReadSettings > read file failed")

	var settingsRead types.Settings
	err = json.Unmarshal(bytes, &settingsRead)
	utils.Assert(err == nil, "ReadSettings > unmarshal bytes failed")

	return &settingsRead
}

func (s *FileStore) WriteSettings(settings types.Settings) {
	formatted, err := json.MarshalIndent(settings, "", " ")
	utils.Assert(err == nil, "WriteSettings > marshal failed")

	err = os.WriteFile(SETTINGS_PATH, formatted, globals.READ_WRITE_PERM)
	utils.Assert(err == nil, "WriteSettings > write failed")
}

func (s *FileStore) ReadTaskDefinitionMap() types.TaskDefinitionMap {
	bytes, err := os.ReadFile(TASK_DEFINITIONS_PATH)
	utils.Assert(err == nil, "ReadTaskDefinitionMap > read file failed")

	var taskDefinitionMap types.TaskDefinitionMap
	err = json.Unmarshal(bytes, &taskDefinitionMap)
	utils.Assert(err == nil, "ReadTaskDefinitionMap > unmarshal bytes failed")

	return taskDefinitionMap
}

func (s *FileStore) WriteTaskDefinitionMap(taskDefinitionMap types.TaskDefinitionMap) {
	formatted, err := json.MarshalIndent(taskDefinitionMap, "", " ")
	utils.Assert(err == nil, "WriteTaskDefinitionMap > marshal failed")

	err = os.WriteFile(TASK_DEFINITIONS_PATH, formatted, globals.READ_WRITE_PERM)
	utils.Assert(err == nil, "WriteTaskDefinitionMap > write failed")
}

func (s *FileStore) ReadTaskDefinitionOrThrow(taskDefinitionMap types.TaskDefinitionMap, ticket string) types.TaskDefinition {
	record, ok := taskDefinitionMap[ticket]
	utils.Assert(ok, "Task definition not found on map")

	return record
}

func (s *FileStore) WriteTimestampOfLastActivity(taskDefinitionMap types.TaskDefinitionMap, ticket string) {
	record := s.ReadTaskDefinitionOrThrow(taskDefinitionMap, ticket)
	record.TimestampOfLastActivity = utils.GetTimestamp()
	taskDefinitionMap[ticket] = record
	s.WriteTaskDefinitionMap(taskDefinitionMap)
}

func (s *FileStore) ReadTaskRecordsByDateMap() types.TaskRecordsByDateMap {
	bytes, err := os.ReadFile(TASK_RECORDS_PATH)
	utils.Assert(err == nil, "ReadTaskRecordsByDateMap > read file failed")

	var taskRecordsByDateMap types.TaskRecordsByDateMap
	err = json.Unmarshal(bytes, &taskRecordsByDateMap)
	utils.Assert(err == nil, "ReadTaskRecordsByDateMap > unmarshal bytes failed")

	return taskRecordsByDateMap
}

func (s *FileStore) WriteTaskRecordsByDateMap(taskRecordsByDateMap types.TaskRecordsByDateMap) {
	formatted, err := json.MarshalIndent(taskRecordsByDateMap, "", " ")
	utils.Assert(err == nil, "failed to marshal tasksToday")

	err = os.WriteFile(TASK_RECORDS_PATH, formatted, globals.READ_WRITE_PERM)
	utils.Assert(err == nil, "failed to write tasksToday")
}

func InitFileStore() *FileStore {
	_, err := os.ReadDir(globals.DATA_DIR)
	utils.Assert(err == nil, "you have to manually create a './data' directory inside the project")

	createTasksFileIfMissing(TASK_DEFINITIONS_PATH)
	createTasksFileIfMissing(TASK_RECORDS_PATH)
	createSettingsFileIfMissing()

	return &FileStore{}
}

func createTasksFileIfMissing(path string) {
	_, err := os.ReadFile(path)

	if errors.Is(err, fs.ErrNotExist) {
		content := []byte("{}")
		err := os.WriteFile(path, content, globals.READ_WRITE_PERM)
		utils.Assert(err == nil, "failed to create a file")
		return
	}

	utils.Assert(err == nil, "failed to read a file")
}

func createSettingsFileIfMissing() {
	_, err := os.ReadFile(SETTINGS_PATH)

	if errors.Is(err, fs.ErrNotExist) {
		content, err := json.MarshalIndent(
			types.Settings{
				ThemeIndex: 0,
			},
			"",
			" ",
		)
		utils.Assert(err == nil, "createSettingsFileIfMissing > marshal failed")

		err = os.WriteFile(SETTINGS_PATH, content, globals.READ_WRITE_PERM)
		utils.Assert(err == nil, "createSettingsFileIfMissing > failed to create a file")
	}
}
