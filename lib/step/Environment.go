package step

type Environment interface {
	Set(key string, value string) error
}
