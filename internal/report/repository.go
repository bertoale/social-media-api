package report

import "gorm.io/gorm"

type Repository interface {
	Create(report *Report) error
	FindByID(id uint) (*Report, error)
	FindAll() ([]*Report, error)
	Update(report *Report) error
}

type repository struct {
	db *gorm.DB
}

// FindAll implements Repository.
func (r *repository) FindAll() ([]*Report, error) {
	var reports []*Report

	err := r.db.
		Preload("User").
		Preload("Post").
		Find(&reports).Error

	if err != nil {
		return nil, err
	}

	return reports, nil
}

// Create implements Repository.
func (r *repository) Create(report *Report) error {
	return r.db.Create(report).Error
}

// FindByID implements Repository.
func (r *repository) FindByID(id uint) (*Report, error) {
	var report Report

	err := r.db.
		Preload("User").
		Preload("Post").
		First(&report, id).Error

	if err != nil {
		return nil, err
	}

	return &report, nil
}

// Update implements Repository.
func (r *repository) Update(report *Report) error {
	return r.db.Save(report).Error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}
