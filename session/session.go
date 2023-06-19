package session

type Session interface {
	SetSessionID(ID string)
	GetSessionID() string
	Set(key string, val any)
	Get(key string, val any) bool
	Del(key string)
	All() map[string]any
	Load()
	Save()
	Flush()
}
