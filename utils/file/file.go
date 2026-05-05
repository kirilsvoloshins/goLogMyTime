package file

import (
	"encoding/json"
	"goLogMyTime/domain/types"
	"goLogMyTime/utils/assert"
	"goLogMyTime/utils/date"
	"os"
)

var PathMarkdown = "./data/main.md"
var pathTaskDefinitions = "data/taskDefinitions.json"
var pathTaskRecords = "data/taskRecords.json"
var READ_WRITE_PERM os.FileMode = 0644

func isFileMissing(path string) bool {
	_, err := os.ReadFile(path)
	return err != nil
}

func createFileIfMissing(path string) {
	if isFileMissing(path) {
		content := []byte("{}")
		err := os.WriteFile(path, content, READ_WRITE_PERM)
		assert.Assert(err == nil, "failed to create a file")
	}
}

func CreateFilesIfMissing() {
	createFileIfMissing(pathTaskDefinitions)
	createFileIfMissing(pathTaskRecords)
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
