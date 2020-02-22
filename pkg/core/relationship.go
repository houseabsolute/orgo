package core

type relationship struct {
	fromTable   string
	toTable     string
	fromColumns []string
	toColumns   []string
}

type HasOne struct {
	relationship
}

type HasMany struct {
	relationship
}

type ManyToMany struct {
	from HasMany
	to   HasMany
}
