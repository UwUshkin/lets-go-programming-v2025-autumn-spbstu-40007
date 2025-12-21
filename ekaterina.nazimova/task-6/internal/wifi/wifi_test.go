package wifi_test

import (
	"errors"
	"net"
	"testing"

	mdlayher "github.com/mdlayher/wifi"
	"github.com/UwUshkin/task-6/internal/wifi"
	"github.com/stretchr/testify/require"
)

var errWifiSys = errors.New("system error")

func TestWiFi(t *testing.T) {
	t.Parallel()
	hw, _ := net.ParseMAC("00:11:22:33:44:55")
	testIface := &mdlayher.Interface{
		Index:        1,
		Name:         "wlan0",
		HardwareAddr: hw,
	}

	t.Run("FullCoverage", func(t *testing.T) {
		t.Parallel()
		mockWiFi := new(MockWiFiHandle)
		
		service := wifi.New(mockWiFi)

		mockWiFi.On("Interfaces").Return([]*mdlayher.Interface{testIface}, nil).Once()
		_, err := service.GetAddresses()
		require.NoError(t, err)

		mockWiFi.On("Interfaces").Return(nil, errWifiSys).Once()
		_, err = service.GetAddresses()
		require.Error(t, err)

		mockWiFi.On("Interfaces").Return([]*mdlayher.Interface{testIface}, nil).Once()
		_, err = service.GetNames()
		require.NoError(t, err)

		mockWiFi.On("Interfaces").Return(nil, errWifiSys).Once()
		_, err = service.GetNames()
		require.Error(t, err)
	})
}
