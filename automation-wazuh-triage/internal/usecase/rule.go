package usecase

import (
	"automation-wazuh-triage/internal/domain"
	"automation-wazuh-triage/internal/entity"
	"context"
)

type ruleUsecase struct {
	ruleRepo domain.RuleRepository
}

func NewRuleUsecase(ruleRepo domain.RuleRepository) domain.RuleUsecase {
	return &ruleUsecase{
		ruleRepo: ruleRepo,
	}
}

func (u *ruleUsecase) GetDetailRules(ctx context.Context, ruleID string) (*entity.WazuhRule, error) {
	return u.ruleRepo.GetDetailRules(ctx, ruleID)
}

func (u *ruleUsecase) GetListRulesByFiles(ctx context.Context, filename string) ([]entity.WazuhRule, error) {
	return u.ruleRepo.GetListRulesByFiles(ctx, filename)
}
