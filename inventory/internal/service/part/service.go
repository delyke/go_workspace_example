package part

import (
	"github.com/delyke/go_workspace_example/inventory/internal/repository"
	def "github.com/delyke/go_workspace_example/inventory/internal/service"
)

var _ def.PartService = (*service)(nil)

type service struct {
	partRepository repository.PartRepository
}

func NewService(partRepository repository.PartRepository) *service {
	return &service{
		partRepository: partRepository,
	}
}
