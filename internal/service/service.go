package service

type URLS struct {
	URLWord string
	URLPDF  string
}

type Gendocer interface {
	Gendoc([]byte, int) (*URLS, error)
	File(path string) ([]byte, error)
}

type Service struct {
	GendocService
}

func NewService() *Service {
	return &Service{
		GendocService: *NewGengocService(),
	}
}
