package handlers

import "github.com/daneofmanythings/calcuroller/internal/config"

var Repo *Repository

type Repository struct {
	App *config.Config
}

func NewRepo(config *config.Config) *Repository {
	return &Repository{
		App: config,
	}
}

func LinkRepository(r *Repository) {
	Repo = r
}
