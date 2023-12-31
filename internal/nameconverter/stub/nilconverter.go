package stub

type stub struct {
}

func New() *stub {
	return &stub{}
}

func (s *stub) Convert(text string) (string, error) {
	return text, nil
}
