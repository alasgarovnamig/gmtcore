package repository_base

import (
	"fmt"

	"github.com/alasgarovnamig/gmtcore/search"
	"gorm.io/gorm"
)

// Singleton repository'ler için genel bir yapı
type GenericRepository[T any] struct {
	db *gorm.DB
}

// Yeni bir GenericRepository oluştur
func NewGenericRepository[T any](db *gorm.DB) *GenericRepository[T] {
	return &GenericRepository[T]{db: db}
}

// GetByID ile Preload kullanarak veri çekme
func (r *GenericRepository[T]) GetByIDWithPreload(id uint, preloads ...string) (T, error) {
	var entity T
	query := r.db

	// Preload edilecek ilişkileri ekle
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	err := query.First(&entity, id).Error
	if err != nil {
		return entity, fmt.Errorf("could not find entity: %w", err)
	}
	return entity, nil
}

// Create fonksiyonu
func (r *GenericRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

// Update fonksiyonu
func (r *GenericRepository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

// Delete fonksiyonu
func (r *GenericRepository[T]) Delete(id uint) error {
	return r.db.Delete(new(T), id).Error
}

// Search fonksiyonu (Preload ile)
// Dinamik Search fonksiyonu
func (r *GenericRepository[T]) SearchWithCriteria(criteriaList []search.Criteria, preloads ...string) ([]T, error) {
	var entities []T
	query := r.db

	// Preload edilecek ilişkileri ekle
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	// SearchCriteria'ya göre sorguyu oluştur
	for _, criteria := range criteriaList {
		switch criteria.Operation {
		case search.Equal:
			query = query.Where(fmt.Sprintf("%s = ?", criteria.Key), criteria.Value)
		case search.NotEqual:
			query = query.Where(fmt.Sprintf("%s != ?", criteria.Key), criteria.Value)
		case search.GreaterThan:
			query = query.Where(fmt.Sprintf("%s > ?", criteria.Key), criteria.Value)
		case search.LessThan:
			query = query.Where(fmt.Sprintf("%s < ?", criteria.Key), criteria.Value)
		case search.GreaterThanEqual:
			query = query.Where(fmt.Sprintf("%s >= ?", criteria.Key), criteria.Value)
		case search.LessThanEqual:
			query = query.Where(fmt.Sprintf("%s <= ?", criteria.Key), criteria.Value)
		case search.In:
			query = query.Where(fmt.Sprintf("%s IN (?)", criteria.Key), criteria.Value)
		case search.NotIn:
			query = query.Where(fmt.Sprintf("%s NOT IN (?)", criteria.Key), criteria.Value)
		case search.Match:
			query = query.Where(fmt.Sprintf("%s LIKE ?", criteria.Key), "%"+criteria.Value.(string)+"%")
		case search.MatchStart:
			query = query.Where(fmt.Sprintf("%s LIKE ?", criteria.Key), criteria.Value.(string)+"%")
		case search.MatchEnd:
			query = query.Where(fmt.Sprintf("%s LIKE ?", criteria.Key), "%"+criteria.Value.(string))

		// Join Child: Bir child tabloya join yapar
		case search.JoinChild:
			joinKey, ok := criteria.Value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid child key")
			}
			query = query.Joins(fmt.Sprintf("JOIN %s ON %s.%s = ?", criteria.Key, criteria.Key, joinKey))

		// Join Grand Child: Child ve grandchild tablolarını join yapar
		case search.JoinGrandChild:
			keys, ok := criteria.Value.([]string)
			if !ok || len(keys) < 2 {
				return nil, fmt.Errorf("invalid child or grandchild keys")
			}
			childKey, grandChildKey := keys[0], keys[1]
			query = query.Joins(fmt.Sprintf("JOIN %s ON %s.%s = ?", criteria.Key, criteria.Key, childKey)).
				Joins(fmt.Sprintf("JOIN %s ON %s.%s = ?", childKey, childKey, grandChildKey))

		// AnyOf: Bir alanın koleksiyon içindeki herhangi bir değerle eşleşip eşleşmediğini kontrol eder
		case search.AnyOf:
			values, ok := criteria.Value.([]interface{})
			if !ok {
				return nil, fmt.Errorf("any of requires a slice of values")
			}
			orConditions := ""
			for _, v := range values {
				orConditions += fmt.Sprintf("%s = '%v' OR ", criteria.Key, v)
			}
			orConditions = orConditions[:len(orConditions)-4] // Son OR kaldır
			query = query.Where(orConditions)

		// IsMember: Bir alanın koleksiyon içinde olup olmadığını kontrol eder
		case search.IsMember:
			values, ok := criteria.Value.([]interface{})
			if !ok {
				return nil, fmt.Errorf("is member requires a slice of values")
			}
			for _, v := range values {
				query = query.Where(fmt.Sprintf("? = ANY(%s)", criteria.Key), v)
			}
		}
	}

	// Sorguyu çalıştır ve sonuçları getir
	err := query.Find(&entities).Error
	if err != nil {
		return nil, fmt.Errorf("could not search entities: %w", err)
	}

	return entities, nil
}
