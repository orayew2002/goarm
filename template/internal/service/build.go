package service

type Service struct {
	repo RepoInterface
}

func NewService(repo RepoInterface) *Service {
	return &Service{
		repo: repo,
	}
}
