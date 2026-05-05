package types

type TaskDefinition struct {
	Ticket                  string `json:"ticket"`
	Info                    string `json:"info"`
	IsFavorite              bool   `json:"isFavorite"`
	TimestampOfLastActivity int64  `json:"timestampOfLastActivity"`
}

type TaskDefinitionMap map[string]TaskDefinition
type TaskDefinitionsRecent []TaskDefinition
type TaskDefinitionsFavorite []TaskDefinition

type TaskRecord struct {
	Ticket           string `json:"ticket"`
	Description      string `json:"description"`
	TimestampOfStart int64  `json:"timestampOfStart"`
	IsActive         bool   `json:"isActive"`
	TimestampOfEnd   int64  `json:"timestampOfEnd"`
}

type TaskRecordsByDateMap map[string]TaskRecords
type TaskRecords []TaskRecord

type TimeTableMap map[string]int64 // task: minutes
type TimeEntry struct {
	Name    string
	Minutes int
}
type TimeTable []TimeEntry
