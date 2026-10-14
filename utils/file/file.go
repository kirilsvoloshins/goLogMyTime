package file

import (
	"encoding/json"
	"errors"
	"goLogMyTime/domain/types"
	"goLogMyTime/utils/assert"
	"goLogMyTime/utils/date"
	"io/fs"
	"os"
)

const DATA_DIR = "./data/"
const PathMarkdown = "./data/main.md"
const pathTaskDefinitions = "data/taskDefinitions.json"
const pathTaskRecords = "data/taskRecords.json"
const READ_WRITE_PERM os.FileMode = 0644

func IsDataDirMissing() bool {
	_, err := os.ReadDir(DATA_DIR)
	return errors.Is(err, fs.ErrNotExist)
}

func createFileIfMissing(path string) bool {
	_, err := os.ReadFile(path)
	if err == nil {
		return true
	}

	if errors.Is(err, fs.ErrNotExist) {
		content := []byte("{}")
		err := os.WriteFile(path, content, READ_WRITE_PERM)
		assert.Assert(err == nil, "failed to create a file")
		return true
	}

	return false
}

func CreateFilesIfMissing() bool {
	ok1 := createFileIfMissing(pathTaskDefinitions)
	ok2 := createFileIfMissing(pathTaskRecords)
	return ok1 && ok2
}

func ReadTaskDefinitionMap() types.TaskDefinitionMap {
	bytes, err := os.ReadFile(pathTaskDefinitions)
	assert.Assert(err == nil, "ReadTaskDefinitionMap > read file failed")

	var taskDefinitionMap types.TaskDefinitionMap
	err = json.Unmarshal(bytes, &taskDefinitionMap)
	assert.Assert(err == nil, "ReadTaskDefinitionMap > unmarshal bytes failed")

	return taskDefinitionMap
}

func WriteTaskDefinitionMap(taskDefinitionMap types.TaskDefinitionMap) {
	formatted, err := json.MarshalIndent(taskDefinitionMap, "", " ")
	assert.Assert(err == nil, "WriteTaskDefinitionMap > marshal failed")

	err = os.WriteFile(pathTaskDefinitions, formatted, READ_WRITE_PERM)
	assert.Assert(err == nil, "WriteTaskDefinitionMap > write failed")
}

func WriteTimestampOfLastActivity(taskDefinitionMap types.TaskDefinitionMap, taskDefinition types.TaskDefinition) {
	record, ok := taskDefinitionMap[taskDefinition.Ticket]
	assert.Assert(ok, "WriteTimestampOfLastActivity > failed to find on map")

	record.TimestampOfLastActivity = date.GetTimestamp()
	taskDefinitionMap[taskDefinition.Ticket] = record
	WriteTaskDefinitionMap(taskDefinitionMap)
}

func ReadTaskRecordsByDateMap() types.TaskRecordsByDateMap {
	bytes, err := os.ReadFile(pathTaskRecords)
	assert.Assert(err == nil, "ReadTaskRecordsByDateMap > read file failed")

	var taskRecordsByDateMap types.TaskRecordsByDateMap
	err = json.Unmarshal(bytes, &taskRecordsByDateMap)
	assert.Assert(err == nil, "ReadTaskRecordsByDateMap > unmarshal bytes failed")

	return taskRecordsByDateMap
}

func WriteTaskRecordsByDateMap(taskRecordsByDateMap types.TaskRecordsByDateMap) {
	formatted, err := json.MarshalIndent(taskRecordsByDateMap, "", " ")
	assert.Assert(err == nil, "failed to marshal tasksToday")

	err = os.WriteFile(pathTaskRecords, formatted, READ_WRITE_PERM)
	assert.Assert(err == nil, "failed to write tasksToday")
}
