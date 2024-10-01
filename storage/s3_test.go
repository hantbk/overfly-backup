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

package storage

import (
	"testing"

	"github.com/hantbk/vtsbackup/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type serviceInfo struct {
	name, endpoint, region, storageClass string
	forcePathStyle                       bool
}

func Test_S3_open(t *testing.T) {
	viper := viper.New()
	viper.Set("bucket", "test-bucket")
	viper.Set("region", "us-east-2")

	base, err := newBase(
		config.ModelConfig{
			DumpPath: "/data/backups",
		},
		"foo/bar",
		// Creating a new base object.
		config.SubConfig{
			Type:  "s3",
			Name:  "test",
			Viper: viper,
		},
	)
	assert.NoError(t, err)

	storage := &S3{
		Base: base,
	}

	err = storage.open()
	assert.NoError(t, err)

	assert.Equal(t, "STANDARD_IA", storage.storageClass)
	assert.Equal(t, "test-bucket", storage.bucket)
	assert.Equal(t, "", storage.path)

	assert.Equal(t, 3, *storage.awsCfg.MaxRetries)
	assert.Equal(t, "us-east-2", *storage.awsCfg.Region)
	assert.Equal(t, float64(300), storage.awsCfg.HTTPClient.Timeout.Seconds())
}

func Test_providerName(t *testing.T) {
	var cases = map[string]serviceInfo{
		"s3":     {"AWS S3", "", "us-east-1", "STANDARD_IA", true},
		"minio":  {"MinIO", "", "us-east-1", "", true},
	}

	base, _ := newBase(config.ModelConfig{}, "test", config.SubConfig{})
	base.viper = viper.New()
	base.viper.SetDefault("bucket", "test-bucket")

	for service, info := range cases {
		s := &S3{Base: base, Service: service}
		s.init()

		assert.Equal(t, info.name, s.providerName(), "providerName for "+service)
		assert.Equal(t, info.endpoint, *s.defaultEndpoint(), "defaultEndpoint for "+service)
		assert.Equal(t, info.region, s.defaultRegion(), "defaultRegion for "+service)
		assert.Equal(t, info.storageClass, s.defaultStorageClass(), "defaultStorageClass for "+service)
		assert.Equal(t, info.forcePathStyle, s.forcePathStyle(), "forcePathStyle for "+service)

		assert.Equal(t, info.region, s.viper.GetString("region"))
		assert.Equal(t, info.endpoint, s.viper.GetString("endpoint"))
		assert.Equal(t, "3", s.viper.GetString("max_retries"))
		assert.Equal(t, "300", s.viper.GetString("timeout"))
	}

}
