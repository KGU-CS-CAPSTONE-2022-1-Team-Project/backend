package owner

type ObjectValidator interface {
	validate() error
}

func Validate(v ObjectValidator) error {
	return v.validate()
}
