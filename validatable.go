package command

//CommandValidator command validator interface
type CommandValidator interface {
	Validate() error
}
