package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/UwUshkin/task-6/internal/wifi"
	netwifi "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errWifi = errors.New("system error")

func TestGetAddresses(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockWiFi := new(MockWiFiHandle)
		service := wifi.New(mockWiFi)

		addr1, _ := net.ParseMAC("00:11:22:33:44:55")
		interfaces := []*netwifi.Interface{
			{Name: "wlan0", HardwareAddr: addr1},
		}

		mockWiFi.On("Interfaces").Return(interfaces, nil)

		addrs, err := service.GetAddresses()
		require.NoError(t, err)
		assert.Len(t, addrs, 1)
		assert.Equal(t, addr1, addrs[0])
		mockWiFi.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		mockWiFi := new(MockWiFiHandle)
		service := wifi.New(mockWiFi)

		mockWiFi.On("Interfaces").Return(nil, errWifi)

		_, err := service.GetAddresses()
		require.Error(t, err)
	})
}

func TestGetNames(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockWiFi := new(MockWiFiHandle)
		service := wifi.New(mockWiFi)

		interfaces := []*netwifi.Interface{
			{Name: "wlan0"},
			{Name: "wlan1"},
		}

		mockWiFi.On("Interfaces").Return(interfaces, nil)

		names, err := service.GetNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"wlan0", "wlan1"}, names)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		mockWiFi := new(MockWiFiHandle)
		service := wifi.New(mockWiFi)

		mockWiFi.On("Interfaces").Return(nil, errWifi)

		_, err := service.GetNames()
		require.Error(t, err)
	})
}
