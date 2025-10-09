package rules

import (
	"fmt"
	"github.com/nalonluap/ci-turbo/internal/app/service"
	"github.com/nalonluap/ci-turbo/internal/domain"
)

var _ service.Rule = (*ParallelismRule)(nil)

type ParallelismRule struct{}

func NewParallelismRule() *ParallelismRule {
	return &ParallelismRule{}
}

func (r *ParallelismRule) Execute(p *domain.Pipeline) []domain.Finding {
	// TODO: Реализовать логику поиска задач, которые можно распараллелить.
	fmt.Println("INFO: Executing ParallelismRule...")
	return nil // Placeholder
}
