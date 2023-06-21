package repos

type GenericRepo[T any] interface {
	Create(T) T
	GetList() []T
	GetOne(string) (T, error)
	Update(string, T) (T, error)
	DeleteOne(string) (bool, error)
}