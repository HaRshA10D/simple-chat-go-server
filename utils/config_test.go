package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigLoadsConfiguration(t *testing.T) {
	tmpDir1, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir1)

	err = os.Mkdir(filepath.Join(tmpDir1, "config"), 0700)
	require.NoError(t, err)

	configJSONFilePath := filepath.Join(tmpDir1, "config", "config.json")
	config := `
        {
            "ServerSettings": { "ServerPort": 9090}
        }
    `
	require.NoError(t, ioutil.WriteFile(configJSONFilePath, []byte(config), 0600))

	loadedConfig, err := LoadConfig("config", filepath.Join(tmpDir1, "config"))
	require.NoError(t, err)
	assert.Equal(t, 9090, *loadedConfig.ServerSettings.ServerPort, "Invalid Port value")
	assert.Nil(t, loadedConfig.DatabaseSettings.DBName, "Expected the missing fields in config file to be nil")

	invalidConfig := `
        {
            NotInAProperJson: { "LogLevel": "INFO"}
        }
        `
	require.NoError(t, ioutil.WriteFile(configJSONFilePath, []byte(invalidConfig), 0600))
	_, err = LoadConfig("config", filepath.Join(tmpDir1, "config"))
	assert.Error(t, err)
}
