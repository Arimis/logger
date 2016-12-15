// Copyright © 2016 Casa Platform
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

package logger

import (
	"github.com/casaplatform/casa"
	"github.com/casaplatform/casa/cmd/casa/environment"
	"github.com/casaplatform/mqtt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func init() {
	environment.RegisterService("logger", &Logger{})
}

// logger prints out MQTT messages
type Logger struct {
	client casa.MessageClient
	topics []string
	casa.Logger
}

func (l *Logger) UseLogger(logger casa.Logger) {
	l.Logger = logger
}

func (l *Logger) Start(config *viper.Viper) error {
	l.topics = config.GetStringSlice("Topics")
	handler := func(msg *casa.Message, err error) {
		if err != nil {
			l.Log("ERROR: " + err.Error())
			return
		}
		l.Log(msg.Topic + ": " + string(msg.Payload))
	}
	client, err := mqtt.NewClient("tcp://127.0.0.1:1883")
	if err != nil {
		return err
	}

	l.client = client
	client.Handle(handler)

	// Subscribe to specified topics
	for _, topic := range l.topics {
		err := client.Subscribe(topic)
		if err != nil {
			return errors.Wrap(err, "subscribing to '"+topic+"' failed")

		}
		l.Log("Logger subscribed to topic: " + topic)
	}

	//return errors.Wrap(errors.New("test error"), "test error generated")
	return nil
}

func (l *Logger) Stop() error {
	if l.client == nil {
		return nil
	}
	return l.client.Close()
}
