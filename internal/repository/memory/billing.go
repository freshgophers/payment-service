package memory

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"payment-service/internal/domain/billing"
	"payment-service/pkg/store"
)

type BillingRepository struct {
	db map[string]billing.Entity
	sync.RWMutex
}

func NewBillingRepository() *BillingRepository {
	return &BillingRepository{
		db: make(map[string]billing.Entity),
	}
}

func (r *BillingRepository) Select(ctx context.Context) (dest []billing.Entity, err error) {
	r.RLock()
	defer r.RUnlock()

	dest = make([]billing.Entity, 0, len(r.db))
	for _, data := range r.db {
		dest = append(dest, data)
	}

	return
}

func (r *BillingRepository) SelectByParentID(ctx context.Context, parentID string) (dest []billing.Entity, err error) {
	r.RLock()
	defer r.RUnlock()

	dest = make([]billing.Entity, 0, len(r.db))
	for _, data := range r.db {
		dest = append(dest, data)
	}

	return
}

func (r *BillingRepository) Create(ctx context.Context, data billing.Entity) (dest string, err error) {
	r.Lock()
	defer r.Unlock()

	id := r.generateID()
	data.ID = id
	r.db[id] = data

	return id, nil
}

func (r *BillingRepository) Get(ctx context.Context, id string) (dest billing.Entity, err error) {
	r.RLock()
	defer r.RUnlock()

	dest, ok := r.db[id]
	if !ok {
		err = store.ErrorNotFound
		return
	}

	return
}

func (r *BillingRepository) Update(ctx context.Context, id string, data billing.Entity) (err error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.db[id]; !ok {
		return store.ErrorNotFound
	}
	r.db[id] = data

	return
}

func (r *BillingRepository) Delete(ctx context.Context, id string) (err error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.db[id]; !ok {
		return store.ErrorNotFound
	}
	delete(r.db, id)

	return
}

func (r *BillingRepository) generateID() string {
	return uuid.New().String()
}
