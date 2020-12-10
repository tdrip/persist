package SQLL

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	sli "github.com/tdrip/logger/pkg/interfaces"
	per "gitlab.venafi.com/dave.brancato/mariachi/backend/persist/pkg/interfaces"
)

type SQLLiteDatastore struct {
	per.IPersistantStorage

	database *sql.DB

	Filename string

	Log sli.ISimpleLogger

	StorageHandlers map[string]per.IStorageHandler
}

func CreateSQLLiteDatastore(log sli.ISimpleLogger, filename string) *SQLLiteDatastore {
	sqlds := SQLLiteDatastore{}
	sqlds.Filename = filename
	sqlds.Log = log
	StorageHandlers := make(map[string]per.IStorageHandler)
	sqlds.StorageHandlers = StorageHandlers
	return &sqlds
}

func CreateOpenSQLLiteDatastore(log sli.ISimpleLogger, filename string) *SQLLiteDatastore {
	sqlds := CreateSQLLiteDatastore(log, filename)
	sqlds.Open()
	return sqlds
}

func (sqlds *SQLLiteDatastore) Open() {
	// further checks to be added here like checking filepath is correct etc
	sqlds.database, _ = sql.Open("sqlite3", sqlds.Filename)
}

// Storage Handlers
// in this implementation that is Tables

func (sqlds *SQLLiteDatastore) GetStorageHandler(name string) (per.IStorageHandler, bool) {
	res, ok := sqlds.StorageHandlers[name]
	return res, ok
}

func (sqlds *SQLLiteDatastore) SetStorageHander(name string, handler per.IStorageHandler) {
	sqlds.Log.LogDebugf("SetStorageHander", "Setting %s, %v", name, handler
	sqlds.StorageHandlers[name] = handler
	sqlds.Log.LogDebugf("SetStorageHander", "Handlers = %v", name, sqlds.StorageHandlers)
}

func (sqlds *SQLLiteDatastore) RemoveStorageHandler(name string) bool {
	_, ok := sqlds.StorageHandlers[name]
	if ok {
		delete(sqlds.StorageHandlers, name)
	}
	return ok
}

func (sqlds *SQLLiteDatastore) GetAllStorageHandlers() map[string]per.IStorageHandler {
	return sqlds.StorageHandlers
}

// Get/Set the logging for the interface
func (sqlds *SQLLiteDatastore) GetLog() sli.ISimpleLogger {
	return sqlds.Log
}

func (sqlds *SQLLiteDatastore) SetLog(logger sli.ISimpleLogger) {
	sqlds.Log = logger
}

func (sqlds *SQLLiteDatastore) GetDatabase() *sql.DB {
	return sqlds.database
}

func (sqlds *SQLLiteDatastore) CreateStructures() per.IQueryResult {

	success := NewEmptySucceedSQLLiteQueryResult()
	for key, element := range sqlds.StorageHandlers {
		sqlds.Log.LogDebugf("CreateStructures", "Handling %s, %v", key, element)
		res := element.CreateStructures()
		if !res.QuerySucceeded() {
			success = NewEmptyFailedSQLLiteQueryResult()
		}
	}
	return success
}
