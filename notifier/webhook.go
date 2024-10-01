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
	"io"
	"net/http"
	"strings"

	"github.com/hantbk/vtsbackup/logger"
)

type Webhook struct {
	Base

	Service string

	method          string
	contentType     string
	buildBody       func(title, message string) ([]byte, error)
	buildWebhookURL func(url string) (string, error)
	checkResult     func(status int, responseBody []byte) error
	buildHeaders    func() map[string]string
}

type webhookPayload struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func NewWebhook(base *Base) *Webhook {
	base.viper.SetDefault("method", "POST")

	return &Webhook{
		Base:        *base,
		Service:     "Webhook",
		method:      base.viper.GetString("method"),
		contentType: "application/json",
		buildBody: func(title, message string) ([]byte, error) {
			return json.Marshal(webhookPayload{
				Title:   title,
				Message: message,
			})
		},
		buildHeaders: func() map[string]string {
			headers := make(map[string]string)
			for key, value := range base.viper.GetStringMapString("headers") {
				headers[key] = value
			}

			return headers
		},
		checkResult: func(status int, responseBody []byte) error {
			if status == 200 {
				return nil
			}

			return fmt.Errorf("status: %d, body: %s", status, string(responseBody))
		},
	}
}

func (s *Webhook) getLogger() logger.Logger {
	return logger.Tag(fmt.Sprintf("Notifier: %s", s.Service))
}

func (s *Webhook) webhookURL() (string, error) {
	url := s.viper.GetString("url")

	if s.buildWebhookURL == nil {
		return url, nil
	}

	return s.buildWebhookURL(url)
}

func (s *Webhook) notify(title string, message string) error {
	logger := s.getLogger()

	url, err := s.webhookURL()
	if err != nil {
		return err
	}

	payload, err := s.buildBody(title, message)
	if err != nil {
		return err
	}

	logger.Infof("Send notification to %s...", url)
	req, err := http.NewRequest(s.method, url, strings.NewReader(string(payload)))
	if err != nil {
		logger.Error(err)
		return err
	}

	req.Header.Set("Content-Type", s.contentType)

	if s.buildHeaders != nil {
		headers := s.buildHeaders()
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer resp.Body.Close()

	var body []byte
	if resp.Body != nil {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}

	if s.checkResult != nil {
		err = s.checkResult(resp.StatusCode, body)
		if err != nil {
			logger.Error(err)
			return nil
		}
	} else {
		logger.Infof("Response body: %s", string(body))
	}

	logger.Info("Notification sent.")

	return nil
}
