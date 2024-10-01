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

package notifier

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_Webhook(t *testing.T) {
	base := &Base{
		viper: viper.New(),
	}

	base.viper.Set("headers", map[string]string{
		"Authorization": "Bearer this-is-token",
	})

	s := NewWebhook(base)
	assert.Equal(t, "Webhook", s.Service)
	assert.Equal(t, "POST", s.method)
	assert.Equal(t, "application/json", s.contentType)

	base.viper.Set("method", "PUT")
	s = NewWebhook(base)
	assert.Equal(t, "PUT", s.method)

	body, err := s.buildBody("This is title", "This is body")
	assert.NoError(t, err)
	assert.Equal(t, `{"title":"This is title","message":"This is body"}`, string(body))

	headers := s.buildHeaders()
	assert.Equal(t, "Bearer this-is-token", headers["Authorization"])

	err = s.checkResult(200, []byte(`{"status":"ok"}`))
	assert.NoError(t, err)

	respBody := `{"status":"error"}`
	err = s.checkResult(403, []byte(respBody))
	assert.EqualError(t, err, "status: 403, body: "+respBody)
}
