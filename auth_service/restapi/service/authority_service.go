package service

import (
	"context"

	"github.com/jae2274/auth-service/auth_service/common/domain"
	"github.com/jae2274/auth-service/auth_service/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func GetAllAuthorities(ctx context.Context, exec boil.ContextExecutor) ([]*models.Authority, error) {
	return models.Authorities(models.AuthorityWhere.AuthorityCode.NEQ(domain.AuthorityAdmin), qm.OrderBy(models.AuthorityColumns.AuthorityID)).All(ctx, exec)
}
