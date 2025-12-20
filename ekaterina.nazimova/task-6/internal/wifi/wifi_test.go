package wifi_test



import (

    "errors"

    "testing"



    "github.com/UwUshkin/task-6/internal/wifi"

    "github.com/stretchr/testify/require"

)



var errSystem = errors.New("system error")



func TestWiFiService(t *testing.T) {

    t.Parallel()



    t.Run("GetAddresses_success", func(t *testing.T) {

        t.Parallel()

        mockHandle := new(MockWiFi)

        service := wifi.New(mockHandle)



        mockHandle.On("Interfaces").Return(nil, nil)



        _, err := service.GetAddresses()

        require.NoError(t, err)

    })



    t.Run("GetAddresses_error", func(t *testing.T) {

        t.Parallel()

        mockHandle := new(MockWiFi)

        service := wifi.New(mockHandle)



        mockHandle.On("Interfaces").Return(nil, errSystem)



        _, err := service.GetAddresses()

        require.Error(t, err)

    })



    t.Run("GetNames_success", func(t *testing.T) {

        t.Parallel()

        mockHandle := new(MockWiFi)

        service := wifi.New(mockHandle)



        mockHandle.On("Interfaces").Return(nil, nil)



        _, err := service.GetNames()

        require.NoError(t, err)

    })

}

