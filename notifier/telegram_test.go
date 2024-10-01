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

func Test_Telegram(t *testing.T) {
	base := &Base{
		viper: viper.New(),
	}

	s := NewTelegram(base)
	s.viper.Set("token", "123213:this-is-my-token")
	s.viper.Set("chat_id", "@backuptest")

	assert.Equal(t, "Telegram", s.Service)
	assert.Equal(t, "POST", s.method)
	assert.Equal(t, "application/json", s.contentType)

	body, err := s.buildBody("This is title", "This is body")
	assert.NoError(t, err)
	assert.Equal(t, `{"chat_id":"@backuptest","text":"This is title\n\nThis is body"}`, string(body))

	url, err := s.buildWebhookURL("")
	assert.NoError(t, err)
	assert.Equal(t, "https://api.telegram.org/bot123213:this-is-my-token/sendMessage", url)

	respBody := `{"ok":true,"result":{"message_id":1,"from":{"id":123213,"is_bot":true,"first_name":"Backup","username":"BackupBot"},"chat":{"id":-100123213,"title":"Backup Test","type":"supergroup"},"date":1610000000,"text":"This is title\n\nThis is body"}}`
	err = s.checkResult(200, []byte(respBody))
	assert.NoError(t, err)

	respBody = `{"ok":false,"error_code":403,"description":"Forbidden: bot was blocked by the user"}`
	err = s.checkResult(403, []byte(respBody))
	assert.EqualError(t, err, "status: 403, body: "+respBody)
}
