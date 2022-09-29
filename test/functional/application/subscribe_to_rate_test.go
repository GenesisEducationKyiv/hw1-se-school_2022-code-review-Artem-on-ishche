package application

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

func TestSubscribe_WithCorrectInput_ReturnsNoError(t *testing.T) {
	address, _ := models.NewEmailAddress("my_address@domain.extension")
	subscribeServiceImpl := getSubscribeService([]models.EmailAddress{})

	err := subscribeServiceImpl.Subscribe(address, &pair)

	assert.Nil(t, err)
}

func TestSubscribe_WithAlreadySubscribedAddress_ReturnsError(t *testing.T) {
	addressString := "my_address@domain.extension"
	address, _ := models.NewEmailAddress(addressString)
	subscribeServiceImpl := getSubscribeService([]models.EmailAddress{*address})

	err := subscribeServiceImpl.Subscribe(address, &pair)

	assert.NotNil(t, err)
	assert.Equal(t, services.ErrEmailAddressAlreadyExists(addressString), err)
}

func TestSubscribe_SuccessivelyAddingSameAddress_ReturnsError(t *testing.T) {
	addressString := "my_address@domain.extension"
	address, _ := models.NewEmailAddress(addressString)
	subscribeServiceImpl := getSubscribeService([]models.EmailAddress{})

	err := subscribeServiceImpl.Subscribe(address, &pair)
	assert.Nil(t, err)

	err = subscribeServiceImpl.Subscribe(address, &pair)
	assert.NotNil(t, err)
	assert.Equal(t, services.ErrEmailAddressAlreadyExists(addressString), err)
}

func TestSubscribe_WithMultipleNewAddresses_ReturnsNoError(t *testing.T) {
	addressStrings := []string{
		"address@what.com",
		"artem.mykytyshyn@gmail.com",
		"someone@some.mail",
	}
	subscribeServiceImpl := getSubscribeService([]models.EmailAddress{})

	for _, addressString := range addressStrings {
		address, _ := models.NewEmailAddress(addressString)

		err := subscribeServiceImpl.Subscribe(address, &pair)

		assert.Nil(t, err)
	}
}
