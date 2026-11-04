package utils

import (
	"encoding/json"
	"errors"
	"goLogMyTime/types"
	"io/fs"
	"os"
)

const DATA_DIR = "./data/"
const MARKDOWN_PATH = DATA_DIR + "main.md"
const TASK_DEFINITIONS_PATH = DATA_DIR + "taskDefinitions.json"
const TASK_RECORDS_PATH = DATA_DIR + "taskRecords.json"
const SETTINGS_PATH = DATA_DIR + "settings.json"
const READ_WRITE_PERM os.FileMode = 0644

func ReadDataDir() {
	_, err := os.ReadDir(DATA_DIR)
	Assert(err == nil, "you have to manually create a './data' directory inside the project")
}

func CreateTasksFileIfMissing(path string) {
	_, err := os.ReadFile(path)

	if errors.Is(err, fs.ErrNotExist) {
		content := []byte("{}")
		err := os.WriteFile(path, content, READ_WRITE_PERM)
		Assert(err == nil, "failed to create a file")
		return
	}

	Assert(err == nil, "failed to read a file")
}

func CreateSettingsIfMissing() {
	_, err := os.ReadFile(SETTINGS_PATH)

	if errors.Is(err, fs.ErrNotExist) {
		content, err := json.MarshalIndent(
			types.Settings{
				ThemeIndex: 0,
			},
			"",
			" ",
		)
		Assert(err == nil, "createSettingsFileIfMissing > marshal failed")

		err = os.WriteFile(SETTINGS_PATH, content, READ_WRITE_PERM)
		Assert(err == nil, "createSettingsFileIfMissing > failed to create a file")
	}
}

func ReadSettings() *types.Settings {
	bytes, err := os.ReadFile(SETTINGS_PATH)
	Assert(err == nil, "CreateAndReadSettings > read file failed")

	var settingsRead types.Settings
	err = json.Unmarshal(bytes, &settingsRead)
	Assert(err == nil, "ReadSettings > unmarshal bytes failed")

	return &settingsRead
}

func WriteSettings(settings types.Settings) {
	formatted, err := json.MarshalIndent(settings, "", " ")
	Assert(err == nil, "WriteSettings > marshal failed")

	err = os.WriteFile(SETTINGS_PATH, formatted, READ_WRITE_PERM)
	Assert(err == nil, "WriteSettings > write failed")
}

func ReadTaskDefinitionMap() types.TaskDefinitionMap {
	bytes, err := os.ReadFile(TASK_DEFINITIONS_PATH)
	Assert(err == nil, "ReadTaskDefinitionMap > read file failed")

	var taskDefinitionMap types.TaskDefinitionMap
	err = json.Unmarshal(bytes, &taskDefinitionMap)
	Assert(err == nil, "ReadTaskDefinitionMap > unmarshal bytes failed")

	return taskDefinitionMap
}

func WriteTaskDefinitionMap(taskDefinitionMap types.TaskDefinitionMap) {
	formatted, err := json.MarshalIndent(taskDefinitionMap, "", " ")
	Assert(err == nil, "WriteTaskDefinitionMap > marshal failed")

	err = os.WriteFile(TASK_DEFINITIONS_PATH, formatted, READ_WRITE_PERM)
	Assert(err == nil, "WriteTaskDefinitionMap > write failed")
}

func ReadTaskDefinitionOrThrow(taskDefinitionMap types.TaskDefinitionMap, ticket string) types.TaskDefinition {
	record, ok := taskDefinitionMap[ticket]
	Assert(ok, "Task definition not found on map")

	return record
}

func WriteTimestampOfLastActivity(taskDefinitionMap types.TaskDefinitionMap, ticket string) {
	record := ReadTaskDefinitionOrThrow(taskDefinitionMap, ticket)
	record.TimestampOfLastActivity = GetTimestamp()
	taskDefinitionMap[ticket] = record
	WriteTaskDefinitionMap(taskDefinitionMap)
}

func ReadTaskRecordsByDateMap() types.TaskRecordsByDateMap {
	bytes, err := os.ReadFile(TASK_RECORDS_PATH)
	Assert(err == nil, "ReadTaskRecordsByDateMap > read file failed")

	var taskRecordsByDateMap types.TaskRecordsByDateMap
	err = json.Unmarshal(bytes, &taskRecordsByDateMap)
	Assert(err == nil, "ReadTaskRecordsByDateMap > unmarshal bytes failed")

	return taskRecordsByDateMap
}

func WriteTaskRecordsByDateMap(taskRecordsByDateMap types.TaskRecordsByDateMap) {
	formatted, err := json.MarshalIndent(taskRecordsByDateMap, "", " ")
	Assert(err == nil, "failed to marshal tasksToday")

	err = os.WriteFile(TASK_RECORDS_PATH, formatted, READ_WRITE_PERM)
	Assert(err == nil, "failed to write tasksToday")
}
