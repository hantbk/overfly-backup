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

package helper

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/dustin/go-humanize"
	"github.com/hako/durafmt"
	"github.com/hantbk/vtsbackup/logger"
)

const (
	progressbarTemplate = `{{string . "time"}} {{string . "prefix"}}{{bar . "[" "=" "=" "-" "]"}} {{percent .}} ({{speed .}})`
)

type ProgressBar struct {
	bar        *pb.ProgressBar
	FileLength int64
	Reader     io.Reader
	logger     logger.Logger
	startTime  time.Time
}

func NewProgressBar(myLogger logger.Logger, reader *os.File) ProgressBar {
	info, _ := reader.Stat()
	fileLength := info.Size()

	bar := pb.ProgressBarTemplate(progressbarTemplate).Start64(fileLength)
	bar.SetWidth(100)
	bar.Set("time", time.Now().Format(logger.TimeFormat))
	bar.Set("prefix", myLogger.Prefix())

	multiReader := bar.NewProxyReader(reader)

	progressBar := ProgressBar{bar, fileLength, multiReader, myLogger, time.Now()}
	progressBar.start()

	return progressBar
}

func (p ProgressBar) start() {
	logger := p.logger

	logger.Infof("-> Uploading (%s)...", humanize.Bytes(uint64(p.FileLength)))
}

func (p ProgressBar) Errorf(format string, err ...any) error {
	p.bar.Finish()

	return fmt.Errorf(format, err...)
}

func (p ProgressBar) Done(url string) {
	logger := p.logger

	p.bar.Finish()
	t := time.Now()
	elapsed := t.Sub(p.startTime)

	logger.Info(fmt.Sprintf("Uploaded: %s (Duration %v)", url, durafmt.Parse(elapsed).LimitFirstN(2).String()))
}
