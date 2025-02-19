package mongo

import (
	"testing"

	"github.com/shellhub-io/shellhub/api/cache"
	"github.com/shellhub-io/shellhub/api/pkg/dbtest"
	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/pkg/api/paginator"
	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestPublicKeyCreate(t *testing.T) {
	data := initData()

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())

	err := mongostore.PublicKeyCreate(data.Context, &data.PublicKey)
	assert.NoError(t, err)
}

func TestPublicKeysList(t *testing.T) {
	data := initData()

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())

	err := mongostore.PublicKeyCreate(data.Context, &data.PublicKey)
	assert.NoError(t, err)

	var publicKeyList []models.PublicKey
	publicKeyList = append(publicKeyList, data.PublicKey)

	returnedPublicKeyList, count, err := mongostore.PublicKeyList(data.Context, paginator.Query{Page: -1, PerPage: -1})
	assert.Equal(t, publicKeyList, returnedPublicKeyList)
	assert.Equal(t, 1, count)
	assert.NoError(t, err)
}

func TestPublicKeyGet(t *testing.T) {
	data := initData()

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())

	err := mongostore.PublicKeyCreate(data.Context, &data.PublicKey)
	assert.NoError(t, err)

	k, err := mongostore.PublicKeyGet(data.Context, data.PublicKey.Fingerprint, data.PublicKey.TenantID)
	assert.NoError(t, err)
	assert.NotEmpty(t, k)
}

func TestPublicKeyUpdate(t *testing.T) {
	data := initData()

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())

	err := mongostore.PublicKeyCreate(data.Context, &data.PublicKey)
	assert.NoError(t, err)

	err = mongostore.PublicKeyCreate(data.Context, &models.PublicKey{
		Data: []byte("teste"), Fingerprint: "fingerprint", TenantID: "tenant2", PublicKeyFields: models.PublicKeyFields{Name: "teste2", Filter: models.PublicKeyFilter{Hostname: ".*"}},
	})
	assert.NoError(t, err)

	update := &models.PublicKeyUpdate{
		PublicKeyFields: models.PublicKeyFields{Name: "teste2", Filter: models.PublicKeyFilter{Hostname: ".*"}},
	}

	k, err := mongostore.PublicKeyUpdate(data.Context, data.PublicKey.Fingerprint, "tenant2", update)
	assert.NoError(t, err)
	assert.Equal(t, k, &models.PublicKey{
		Data: []byte("teste"), Fingerprint: "fingerprint", TenantID: "tenant2", PublicKeyFields: models.PublicKeyFields{Name: "teste2", Filter: models.PublicKeyFilter{Hostname: ".*"}},
	})

	_, err = mongostore.PublicKeyUpdate(data.Context, "fingerprint2", "tenant1", update)
	assert.EqualError(t, err, store.ErrNoDocuments.Error())
}

func TestPublicKeyDelete(t *testing.T) {
	data := initData()

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())

	err := mongostore.PublicKeyCreate(data.Context, &data.PublicKey)
	assert.NoError(t, err)

	err = mongostore.PublicKeyDelete(data.Context, data.PublicKey.Fingerprint, data.PublicKey.TenantID)
	assert.NoError(t, err)
}
