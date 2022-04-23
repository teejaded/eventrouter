/*
Copyright 2017 The Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sinks

import (
	"bytes"
	"net"
	"sync"
	"time"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
)

// SyslogSink implements the EventSinkInterface
type SyslogSink struct {
	address string
	conn    net.Conn
	mu      sync.Mutex
}

// NewSyslogSink will create a new SyslogSink with default options, returned as an EventSinkInterface
func NewSyslogSink(address string) (EventSinkInterface, error) {
	s := SyslogSink{
		address: address,
	}

	err := s.connect()
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *SyslogSink) connect() error {
	glog.Infof("Connecting to %s", s.address)
	var err error
	s.conn, err = net.DialTimeout("tcp", s.address, 10*time.Second)
	if err != nil {
		glog.Errorf("Connection error: %v", err)
		return err
	}
	return nil
}

// UpdateEvents implements EventSinkInterface.UpdateEvents
func (s *SyslogSink) UpdateEvents(eNew *v1.Event, eOld *v1.Event) {
	buf := bytes.NewBuffer(make([]byte, 0, 4096))

	eData := NewEventData(eNew, eOld)
	eData.WriteRFC5424(buf)

	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.conn.Write(buf.Bytes())
	if err != nil {
		glog.Errorf("Error transmitting log, reconnecting: %v", err)

		s.conn.Close()
		err := s.connect()
		if err != nil {
			panic(err.Error())
		}

		_, err = s.conn.Write(buf.Bytes())
		if err != nil {
			panic(err.Error())
		}
	}
}
