package memory

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"payment-service/internal/domain/category"
	"payment-service/pkg/store"
)

type CategoryRepository struct {
	db map[string]category.Entity
	sync.RWMutex
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{
		db: make(map[string]category.Entity),
	}
}

func (r *CategoryRepository) Select(ctx context.Context) (dest []category.Entity, err error) {
	r.RLock()
	defer r.RUnlock()

	dest = make([]category.Entity, 0, len(r.db))
	for _, data := range r.db {
		dest = append(dest, data)
	}

	return
}

func (r *CategoryRepository) SelectByParentID(ctx context.Context, parentID string) (dest []category.Entity, err error) {
	r.RLock()
	defer r.RUnlock()

	dest = make([]category.Entity, 0, len(r.db))
	for _, data := range r.db {
		dest = append(dest, data)
	}

	return
}

func (r *CategoryRepository) Create(ctx context.Context, data category.Entity) (dest string, err error) {
	r.Lock()
	defer r.Unlock()

	id := r.generateID()
	data.ID = id
	r.db[id] = data

	return id, nil
}

func (r *CategoryRepository) Get(ctx context.Context, id string) (dest category.Entity, err error) {
	r.RLock()
	defer r.RUnlock()

	dest, ok := r.db[id]
	if !ok {
		err = store.ErrorNotFound
		return
	}

	return
}

func (r *CategoryRepository) Update(ctx context.Context, id string, data category.Entity) (err error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.db[id]; !ok {
		return store.ErrorNotFound
	}
	r.db[id] = data

	return
}

func (r *CategoryRepository) Delete(ctx context.Context, id string) (err error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.db[id]; !ok {
		return store.ErrorNotFound
	}
	delete(r.db, id)

	return
}

func (r *CategoryRepository) generateID() string {
	return uuid.New().String()
}
