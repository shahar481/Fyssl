package base

type Target interface {
	ProcessTarget(buffer *[]byte) (*[]byte, error)
	Close()
}