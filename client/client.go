package client

type client interface {
	MakePIMRequest(url string, role string) (string, error)
}
