package conversions

import (
	"goLogMyTime/domain/types"
	"sort"
	"time"
)

func mapPointerToArray(taskDefinitionMap *types.TaskDefinitionMap) []types.TaskDefinition {
	out := []types.TaskDefinition{}

	// TODO: better way of converting map to array
	for _, value := range *taskDefinitionMap {
		out = append(out, value)
	}

	return out
}

func MapPointerToArrayByModified(taskDefinitionMap *types.TaskDefinitionMap) []types.TaskDefinition {
	tasks := mapPointerToArray(taskDefinitionMap)

	sort.Slice(tasks, func(i, j int) bool {
		aUnix := time.Unix(tasks[i].TimestampOfLastActivity, 0)
		bUnix := time.Unix(tasks[j].TimestampOfLastActivity, 0)
		return aUnix.After(bUnix)
	})

	return tasks
}

func MapPointerToArrayRecent(taskDefinitionMap *types.TaskDefinitionMap) []types.TaskDefinition {
	tasksSortedByModified := MapPointerToArrayByModified(taskDefinitionMap)

	var tasksRecentCapped types.TaskDefinitionsRecent
	if len(tasksSortedByModified) >= 5 {
		tasksRecentCapped = make(types.TaskDefinitionsRecent, 5)
	} else {
		tasksRecentCapped = make(types.TaskDefinitionsRecent, len(tasksSortedByModified))
	}

	copy(tasksRecentCapped, tasksSortedByModified)
	return tasksRecentCapped
}

// TODO: sort them once and later reuse here :)
func MapPointerToArrayFavorite(taskDefinitionMap *types.TaskDefinitionMap) []types.TaskDefinition {
	tasksFavorite := []types.TaskDefinition{}
	tasksSortedByModified := MapPointerToArrayByModified(taskDefinitionMap)

	for _, task := range tasksSortedByModified {
		if task.IsFavorite {
			tasksFavorite = append(tasksFavorite, task)
		}
	}

	return tasksFavorite
}

func GetTaskRecordActive(taskRecordsTodayPtr *types.TaskRecords) *types.TaskRecord {
	if taskRecordsTodayPtr == nil {
		return nil
	}

	var records = *taskRecordsTodayPtr
	for i := range records {
		if records[i].IsActive {
			return &records[i]
		}
	}
	return nil
}
