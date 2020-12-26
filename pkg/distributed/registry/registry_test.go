package registry

import (
	"testing"

	"github.com/google/uuid"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/files"

	"github.com/stretchr/testify/assert"
)

func Test__NewRegistry(t *testing.T) {
	assert := assert.New(t)

	dbPath, err := files.NewTempFile("ru.wailorman.fftb_test", "sqlite.db")

	assert.Nil(err, "Creating database tmp file")
	defer dbPath.Remove()

	_, err = NewSqliteRegistry(dbPath.FullPath())

	assert.Nil(err, "Creating registry")
}

func Test__PersistTask(t *testing.T) {
	assert := assert.New(t)

	dbPath, err := files.NewTempFile("ru.wailorman.fftb_test", "sqlite.db")

	assert.Nil(err, "Creating database tmp file")
	defer dbPath.Remove()

	registry, err := NewSqliteRegistry(dbPath.FullPath())

	assert.Nil(err, "Creating registry")

	convertTask := &models.ConvertTask{
		Identity:             uuid.New().String(),
		OrderIdentity:        "",
		Type:                 models.ConvertV1Type,
		StorageClaimIdentity: "",

		Muxer:            "",
		VideoCodec:       "",
		HWAccel:          "",
		VideoBitRate:     "",
		VideoQuality:     0,
		Preset:           "",
		Scale:            "",
		KeyframeInterval: 0,
	}

	err = registry.PersistTask(convertTask)

	assert.Nil(err, "Persisting task")

	// TODO: validate FindByID response
}

func Test__PersistOrder(t *testing.T) {
	assert := assert.New(t)

	dbPath, err := files.NewTempFile("ru.wailorman.fftb_test", "sqlite.db")

	assert.Nil(err, "Creating database tmp file")
	defer dbPath.Remove()

	registry, err := NewSqliteRegistry(dbPath.FullPath())

	assert.Nil(err, "Creating registry")

	convertOrder := &models.ConvertOrder{
		Identity: uuid.New().String(),
		Type:     models.ConvertV1Type,

		Muxer:            "",
		VideoCodec:       "",
		HWAccel:          "",
		VideoBitRate:     "",
		VideoQuality:     0,
		Preset:           "",
		Scale:            "",
		KeyframeInterval: 0,
	}

	err = registry.PersistOrder(convertOrder)

	assert.Nil(err, "Persisting task")
}
