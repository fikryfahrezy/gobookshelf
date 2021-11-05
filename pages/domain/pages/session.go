package pages

type Session interface {
	Create(v string) string
	Get(k string) string
	Delete(k string)
}
