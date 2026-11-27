package store

import (
	"goLogMyTime/types"
)

type StorageStore struct {
	File *FileStore
}

func InitStorageStore() *StorageStore {
	return &StorageStore{
		File: InitFileStore(),
	}
}

func (s *StorageStore) ReadSettings() *types.Settings {
	return s.File.ReadSettings()
}

func (s *StorageStore) WriteSettings(settings types.Settings) {
	s.File.WriteSettings(settings)
}

func (s *StorageStore) ReadTaskDefinitionMap() types.TaskDefinitionMap {
	return s.File.ReadTaskDefinitionMap()
}

func (s *StorageStore) WriteTaskDefinitionMap(taskDefinitionMap types.TaskDefinitionMap) {
	s.File.WriteTaskDefinitionMap(taskDefinitionMap)
}

func (s *StorageStore) ReadTaskDefinitionOrThrow(taskDefinitionMap types.TaskDefinitionMap, ticket string) types.TaskDefinition {
	return s.File.ReadTaskDefinitionOrThrow(taskDefinitionMap, ticket)
}

func (s *StorageStore) WriteTimestampOfLastActivity(taskDefinitionMap types.TaskDefinitionMap, ticket string) {
	s.File.WriteTimestampOfLastActivity(taskDefinitionMap, ticket)
}

func (s *StorageStore) ReadTaskRecordsByDateMap() types.TaskRecordsByDateMap {
	return s.File.ReadTaskRecordsByDateMap()
}

func (s *StorageStore) WriteTaskRecordsByDateMap(taskRecordsByDateMap types.TaskRecordsByDateMap) {
	s.File.WriteTaskRecordsByDateMap(taskRecordsByDateMap)
}
