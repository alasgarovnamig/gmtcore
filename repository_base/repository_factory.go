package repository_base

import (
	"gorm.io/gorm"
	"sync"
)

// Singleton için factory yapısı
type RepositoryFactory struct {
	db    *gorm.DB
	cache sync.Map
}

// Yeni bir RepositoryFactory oluştur
func NewRepositoryFactory(db *gorm.DB) *RepositoryFactory {
	return &RepositoryFactory{
		db: db,
	}
}

// Generic repository'yi döndüren bağımsız fonksiyon
func GetRepository[T any](factory *RepositoryFactory) *GenericRepository[T] {
	repoType := new(T)

	// Cache içinde repository'yi kontrol et
	if repo, ok := factory.cache.Load(repoType); ok {
		return repo.(*GenericRepository[T])
	}

	// Eğer cache'de yoksa yeni repository oluştur ve cache'e ekle
	repo := NewGenericRepository[T](factory.db)
	factory.cache.Store(repoType, repo)
	return repo
}
