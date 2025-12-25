package part

import (
	"context"
	"log"
	"sync"

	def "github.com/delyke/go_workspace_example/inventory/internal/repository"
	repoModel "github.com/delyke/go_workspace_example/inventory/internal/repository/model"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	mu    sync.RWMutex
	parts map[string]*repoModel.Part
}

func NewRepository() *repository {
	r := &repository{
		parts: make(map[string]*repoModel.Part),
	}
	err := r.Init(context.Background())
	if err != nil {
		log.Printf("init parts err: %v", err)
	}
	return r
}
