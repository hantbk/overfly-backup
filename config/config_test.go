// Copyright © 2024 Ha Nguyen <captainnemot1k60@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testConfigFile = "../vtsbackup_test.yml"
)

func init() {
	os.Setenv("S3_ACCESS_KEY_ID", "xxxxxxxxxxxxxxxxxxxx")
	os.Setenv("S3_SECRET_ACCESS_KEY", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	if err := Init(testConfigFile); err != nil {
		panic(err.Error())
	}
}

func TestModelsLength(t *testing.T) {
	assert.Equal(t, Exist, true)
	assert.Equal(t, len(Models), 4)
}

func TestModel(t *testing.T) {
	model := GetModelConfigByName("test")

	assert.Equal(t, model.Name, "test")
	assert.Equal(t, model.Description, "test backup")

	// compress_with
	assert.Equal(t, model.CompressWith.Type, "tgz")
	assert.NotNil(t, model.CompressWith.Viper)

	// storages
	assert.Equal(t, model.DefaultStorage, "minio")
	assert.Equal(t, model.Storages["minio"].Type, "minio")
	assert.Equal(t, model.Storages["minio"].Viper.GetString("bucket"), "vtsbackup-test")
	assert.Equal(t, model.Storages["minio"].Viper.GetString("endpoint"), "http://127.0.0.1:9000")
	assert.Equal(t, model.Storages["minio"].Viper.GetString("access_key_id"), "test-user")
	assert.Equal(t, model.Storages["minio"].Viper.GetString("secret_access_key"), "test-user-secret")

	// notifiers
	assert.Equal(t, model.Notifiers["telegram"].Type, "telegram")
	assert.Equal(t, model.Notifiers["telegram"].Viper.GetString("chat_id"), "@vtsbackuptest")
	assert.Equal(t, model.Notifiers["telegram"].Viper.GetString("token"), "your-token-here")

	// archive
	includes := model.Archive.GetStringSlice("includes")
	assert.Len(t, includes, 1)
	assert.Contains(t, includes, "/Users/hant/Documents/")
	// schedule
	schedule := model.Schedule
	assert.Equal(t, true, schedule.Enabled)
	assert.Equal(t, "0 0 * * *", schedule.Cron)
}

func Test_otherModels(t *testing.T) {
	model := GetModelConfigByName("normal_files")

	// default_storage
	assert.Equal(t, model.DefaultStorage, "scp")

	// schedule
	schedule := model.Schedule
	assert.Equal(t, true, schedule.Enabled)
	assert.Equal(t, "", schedule.Cron)
	assert.Equal(t, "1day", schedule.Every)
	assert.Equal(t, "0:30", schedule.At)

	model = GetModelConfigByName("test_model")
	assert.Equal(t, false, model.Schedule.Enabled)
}

func Test_ScheduleConfig_String(t *testing.T) {
	schedule := ScheduleConfig{
		Enabled: true,
		Every:   "1day",
		At:      "0:30",
	}
	assert.Equal(t, schedule.String(), "every 1day at 0:30")

	schedule = ScheduleConfig{
		Enabled: true,
		Every:   "1day",
	}
	assert.Equal(t, schedule.String(), "every 1day")

	schedule = ScheduleConfig{
		Enabled: true,
		Cron:    "5 4 * * sun",
	}

	assert.Equal(t, schedule.String(), "cron 5 4 * * sun")

	schedule = ScheduleConfig{
		Enabled: false,
	}
	assert.Equal(t, schedule.String(), "disabled")
}

func TestExpandEnv(t *testing.T) {
	model := GetModelConfigByName("expand_env")

	assert.Equal(t, model.Storages["s3"].Type, "s3")
	assert.Equal(t, model.Storages["s3"].Viper.GetString("access_key_id"), "xxxxxxxxxxxxxxxxxxxx")
	assert.Equal(t, model.Storages["s3"].Viper.GetString("secret_access_key"), "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
}

func TestWebConfig(t *testing.T) {
	assert.Equal(t, Web.Host, "0.0.0.0")
	assert.Equal(t, Web.Port, "1201")
	assert.Equal(t, Web.Username, "admin")
	assert.Equal(t, Web.Password, "admin")
}

func TestInitWithNotExistsConfigFile(t *testing.T) {
	err := Init("config/path/not-exist.yml")
	assert.NotNil(t, err)
}

func TestWatchConfigToReload(t *testing.T) {
	err := Init(testConfigFile)
	assert.Nil(t, err)

	lastUpdatedAt := UpdatedAt.UnixNano()
	time.Sleep(1 * time.Millisecond)

	// Touch `testConfigFile` to trigger file changes event
	err = updateFile(testConfigFile)
	assert.Nil(t, err)

	// Wait for reload
	time.Sleep(10 * time.Millisecond)

	// check config reload updated_at
	assert.NotEqual(t, lastUpdatedAt, UpdatedAt.UnixNano())
}

func updateFile(path string) error {
	// Open file and write it again without any changes
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
