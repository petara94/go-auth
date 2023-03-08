package repo

type RepoError string

func (r RepoError) Error() string {
	return string(r)
}

const (
	ErrNotFound RepoError = "not found"
)
