package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
)

type vaccinePlanTemplateRepositoryImpl struct {
	db *gorm.DB
}

// NewVaccinePlanTemplateRepository 创建疫苗计划模板仓储实现
func NewVaccinePlanTemplateRepository(db *gorm.DB) repository.VaccinePlanTemplateRepository {
	return &vaccinePlanTemplateRepositoryImpl{db: db}
}

func (r *vaccinePlanTemplateRepositoryImpl) FindAll(ctx context.Context) ([]*entity.VaccinePlanTemplate, error) {
	var templates []*entity.VaccinePlanTemplate
	err := r.db.WithContext(ctx).
		Order("sort_order ASC, age_in_months ASC, dose_number ASC").
		Find(&templates).Error
	return templates, err
}

func (r *vaccinePlanTemplateRepositoryImpl) FindByID(ctx context.Context, templateID int64) (*entity.VaccinePlanTemplate, error) {
	var template entity.VaccinePlanTemplate
	err := r.db.WithContext(ctx).Where("id = ?", templateID).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *vaccinePlanTemplateRepositoryImpl) Create(ctx context.Context, template *entity.VaccinePlanTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

func (r *vaccinePlanTemplateRepositoryImpl) BatchCreate(ctx context.Context, templates []*entity.VaccinePlanTemplate) error {
	if len(templates) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(templates, 100).Error
}
