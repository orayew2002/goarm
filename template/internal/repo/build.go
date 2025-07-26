package repo

type Repo struct {
	db DB
}

func NewRepo(db DB) *Repo {
	return &Repo{
		db: db,
	}
}
