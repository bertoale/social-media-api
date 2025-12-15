package blog

import "gorm.io/gorm"

type Repository interface {
	Create(blog *Blog) error
	FindByID(id uint) (*Blog, error)
	Update(blog *Blog) error
	Delete(id uint) error
	FindAll() ([]*Blog, error)
}

type repository struct {
	db *gorm.DB
}

// Create implements Repository.
func (r *repository) Create(blog *Blog) error {
	return r.db.Create(blog).Error
}

// Delete implements Repository.
func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Blog{}, id).Error
}

// FindAll implements Repository.
func (r *repository) FindAll() ([]*Blog, error) {
	var blogs []*Blog
	if err := r.db.Preload("Author").Find(&blogs).Error; err != nil {
		return nil, err
	}
	return blogs, nil
}

// FindByID implements Repository.
func (r *repository) FindByID(id uint) (*Blog, error) {
	var blog Blog
	if err := r.db.Preload("Author").First(&blog, id).Error; err != nil {
		return nil, err
	}
	return &blog, nil
}

// Update implements Repository.
func (r *repository) Update(blog *Blog) error {
	return r.db.Save(blog).Error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
