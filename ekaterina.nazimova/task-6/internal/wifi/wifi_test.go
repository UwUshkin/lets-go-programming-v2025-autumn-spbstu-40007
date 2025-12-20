package wifi

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
)

func TestGetAddresses(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockWiFi := new(MockWiFiHandle)
		service := New(mockWiFi)

		addr1, _ := net.ParseMAC("00:11:22:33:44:55")
		interfaces := []*wifi.Interface{
			{Name: "wlan0", HardwareAddr: addr1},
		}

		mockWiFi.On("Interfaces").Return(interfaces, nil)

		addrs, err := service.GetAddresses()
		assert.NoError(t, err)
		assert.Len(t, addrs, 1)
		assert.Equal(t, addr1, addrs[0])
		mockWiFi.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockWiFi := new(MockWiFiHandle)
		service := New(mockWiFi)

		mockWiFi.On("Interfaces").Return(nil, errors.New("system error"))

		_, err := service.GetAddresses()
		assert.Error(t, err)
	})
}

func TestGetNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockWiFi := new(MockWiFiHandle)
		service := New(mockWiFi)

		interfaces := []*wifi.Interface{
			{Name: "wlan0"},
			{Name: "wlan1"},
		}

		mockWiFi.On("Interfaces").Return(interfaces, nil)

		names, err := service.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"wlan0", "wlan1"}, names)
	})

	t.Run("error", func(t *testing.T) {
		mockWiFi := new(MockWiFiHandle)
		service := New(mockWiFi)

		mockWiFi.On("Interfaces").Return(nil, errors.New("fail"))

		_, err := service.GetNames()
		assert.Error(t, err)
	})
}
