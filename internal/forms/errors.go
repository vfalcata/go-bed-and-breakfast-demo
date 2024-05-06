package forms

type errors map[string][]string

// Adds an error message for a form field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// errors are stored here fifo, so if multiple errors proc, the lowest index one triggers first
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
