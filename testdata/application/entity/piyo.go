package entity

type Piyo struct {
	Tel string
}

func NewPiyo(tel string) *Piyo {
	return &Piyo{
		Tel: tel,
	}
}
