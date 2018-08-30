package sources

// images sources

type SourceWrapper interface {
	GetOne() []string
	GetAll() []string
}
