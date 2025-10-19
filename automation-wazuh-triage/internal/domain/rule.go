package domain

import (
	"automation-wazuh-triage/internal/entity"
	"context"
)

type RuleUsecase interface {
	GetDetailRules(ctx context.Context, ruleID string) (*entity.WazuhRule, error)
	GetListRulesByFiles(ctx context.Context, filename string) ([]entity.WazuhRule, error)
}

type RuleRepository interface {
	GetDetailRules(ctx context.Context, ruleID string) (*entity.WazuhRule, error)
	GetListRulesByFiles(ctx context.Context, filename string) ([]entity.WazuhRule, error)
}
