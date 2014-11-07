package db

type Db map[string]interface{}

func (db Db) Get(key string, reply *interface{}) error {
	*reply = db[key]
	return nil
}

func NewDb() Db {
	return make(Db)
}
