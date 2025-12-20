package wifi 

import (
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, fmt.Errorf("interfaces error: %w", args.Error(1))
	}
	
	ifaces, ok := args.Get(0).([]*wifi.Interface)
	if !ok {
		return nil, fmt.Errorf("type assertion failed")
	}
	
	return ifaces, args.Error(1)
}
