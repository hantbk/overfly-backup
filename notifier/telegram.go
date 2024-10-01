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
	"encoding/json"
	"fmt"

	"github.com/hantbk/vtsbackup/helper"
)

type telegramPayload struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

const DEFAULT_TELEGRAM_ENDPOINT = "api.telegram.org"

func NewTelegram(base *Base) *Webhook {
	return &Webhook{
		Base:        *base,
		Service:     "Telegram",
		method:      "POST",
		contentType: "application/json",
		buildWebhookURL: func(url string) (string, error) {
			token := base.viper.GetString("token")
			endpoint := DEFAULT_TELEGRAM_ENDPOINT
			if base.viper.IsSet("endpoint") {
				endpoint = base.viper.GetString("endpoint")
			}

			endpoint = helper.FormatEndpoint(endpoint)

			return fmt.Sprintf("%s/bot%s/sendMessage", endpoint, token), nil
		},
		buildBody: func(title, message string) ([]byte, error) {
			chat_id := base.viper.GetString("chat_id")

			payload := telegramPayload{
				ChatID: chat_id,
				Text:   fmt.Sprintf("%s\n\n%s", title, message),
			}

			return json.Marshal(payload)
		},
		checkResult: func(status int, body []byte) error {
			if status != 200 {
				return fmt.Errorf("status: %d, body: %s", status, string(body))
			}

			return nil
		},
	}
}
