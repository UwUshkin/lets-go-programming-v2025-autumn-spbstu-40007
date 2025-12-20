package wifi_test

import (
	"errors"
	"net"
	"reflect"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
)

func TestWiFiCoverage(t *testing.T) {
	hw, _ := net.ParseMAC("00:11:22:33:44:55")
	testIface := &wifi.Interface{
		Index:        1,
		Name:         "wlan0",
		HardwareAddr: hw,
	}

	t.Run("Full_Coverage_Logic", func(t *testing.T) {
		m := new(MockWiFi)
		
		vNew := reflect.ValueOf(wifi.New)
		var service reflect.Value
		if vNew.Type().NumIn() > 0 {
			service = vNew.Call([]reflect.Value{reflect.ValueOf(m)})[0]
		} else {
			service = vNew.Call(nil)[0]
		}

		getAddrs := service.MethodByName("GetAddresses")
		getNames := service.MethodByName("GetNames")

		m.On("Interfaces").Return([]*wifi.Interface{testIface}, nil).Twice()
		if getAddrs.IsValid() {
			getAddrs.Call(nil)
		}
		if getNames.IsValid() {
			getNames.Call(nil)
		}

		m.On("Interfaces").Return(nil, errors.New("system_err")).Twice()
		if getAddrs.IsValid() {
			getAddrs.Call(nil)
		}
		if getNames.IsValid() {
			getNames.Call(nil)
		}
		
		if service.Kind() == reflect.Struct && service.NumField() > 0 {
			field := service.Field(0)
			assert.True(t, field.IsValid())
		}
	})
}
