package db

type Db map[string]interface{}

func (db Db) Get(key string) interface{} {
	return db[key]
}

func (db Db) Set(key string, value interface{}) interface{} {
	db[key] = value
	return value
}

func NewDb() *Db {
	return new(Db)
}
