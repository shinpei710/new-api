package model

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/stretchr/testify/require"
)

func TestUpdateOptionPersistsAndUpdatesOptionMap(t *testing.T) {
	truncateTables(t)

	originalRetryTimes := common.RetryTimes
	t.Cleanup(func() {
		common.RetryTimes = originalRetryTimes
	})

	common.OptionMapRWMutex.Lock()
	originalOptionMap := common.OptionMap
	common.OptionMap = map[string]string{}
	common.OptionMapRWMutex.Unlock()
	t.Cleanup(func() {
		common.OptionMapRWMutex.Lock()
		common.OptionMap = originalOptionMap
		common.OptionMapRWMutex.Unlock()
	})

	require.NoError(t, UpdateOption("RetryTimes", "3"))

	var saved Option
	require.NoError(t, DB.Where("key = ?", "RetryTimes").First(&saved).Error)
	require.Equal(t, "3", saved.Value)

	common.OptionMapRWMutex.RLock()
	require.Equal(t, "3", common.OptionMap["RetryTimes"])
	common.OptionMapRWMutex.RUnlock()
	require.Equal(t, 3, common.RetryTimes)
}
