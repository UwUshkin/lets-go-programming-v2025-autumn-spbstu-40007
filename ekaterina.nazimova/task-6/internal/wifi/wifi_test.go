package wifi 

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errWifiSys = errors.New("system error")

func TestGetAddresses(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockWiFi := new(MockWiFiHandle)
		service := New(mockWiFi)

		addr1, _ := net.ParseMAC("00:11:22:33:44:55")
		ifaces := []*wifi.Interface{
			{Name: "wlan0", HardwareAddr: addr1},
		}

		mockWiFi.On("Interfaces").Return(ifaces, nil)

		addrs, err := service.GetAddresses()
		require.NoError(t, err)
		assert.Len(t, addrs, 1)
		assert.Equal(t, addr1, addrs[0])
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		mockWiFi := new(MockWiFiHandle)
		service := New(mockWiFi)
		mockWiFi.On("Interfaces").Return(nil, errWifiSys)

		_, err := service.GetAddresses()
		require.Error(t, err)
	})
}

func TestGetNames(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockWiFi := new(MockWiFiHandle)
		service := New(mockWiFi)

		ifaces := []*wifi.Interface{
			{Name: "wlan0"},
		}

		mockWiFi.On("Interfaces").Return(ifaces, nil)

		names, err := service.GetNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"wlan0"}, names)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		mockWiFi := new(MockWiFiHandle)
		service := New(mockWiFi)
		mockWiFi.On("Interfaces").Return(nil, errWifiSys)

		_, err := service.GetNames()
		require.Error(t, err)
	})
}
