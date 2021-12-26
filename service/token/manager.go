package token

type Manager interface {
	GenerateFileToken() (string, error)
	GenerateFileID() (string, error)
	GenerateImageToken() (string, error)
}
