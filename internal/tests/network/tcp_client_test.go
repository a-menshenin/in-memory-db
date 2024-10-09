package network

import (
	"errors"
	"fmt"
	"testing"
	"time"
	"umemory/internal/network"
	"umemory/internal/network/mock"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

type TCPClientTestCase struct {
	name string
	request string
	expected []byte
	prepare func()
	expectedErr string
}

func TestTCPClient(t *testing.T) {
	ctrl := gomock.NewController(t) 
	defer ctrl.Finish()

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	mockConn := mock.NewMockConn(ctrl)

	duration, err := time.ParseDuration("1m")
	if err != nil {
		fmt.Println("time.ParseDuration error")

		t.Fail()
	}

	maxMsgSize := 1024
	connDeadline := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	cfg := network.TCPClientConfig{
		MaxMessageSize: &maxMsgSize,
		ConnectionDeadline: &connDeadline,
		IdleTimeout: &duration,
	}
	client, err := network.NewTCPClient(cfg, mockConn, zap.NewNop())
	if err != nil {
		t.Errorf("network.NewTCPClient error: %s", err.Error())
	}

	var testCases = []TCPClientTestCase{
		{
			name: "Send setIdle timeout error",
			request: "case1",
			expected: []byte(""),
			prepare: func() {
				cfg := network.TCPClientConfig{
					MaxMessageSize: &maxMsgSize,
					ConnectionDeadline: &connDeadline,
				}
				client, err = network.NewTCPClient(cfg, mockConn, zap.NewNop())
				if err != nil {
					t.Errorf("network.NewTCPClient error: %s", err.Error())
				}

				mockConn.EXPECT().SetDeadline(connDeadline).Return(errors.New("err"))
			},
			expectedErr: "Client internal error",
		},
		{
			name: "Conn write error",
			request: "case2",
			expected: []byte(""),
			prepare: func() {
				mockConn.EXPECT().SetDeadline(connDeadline).Return(nil)
				mockConn.EXPECT().Write([]byte("case2")).Return(0, errors.New("err"))
			},
			expectedErr: "Client send data error",
		},
		{
			name: "Conn read error",
			request: "case3",
			expected: []byte(""),
			prepare: func() {
				mockConn.EXPECT().SetDeadline(connDeadline).Return(nil)
				mockConn.EXPECT().Write([]byte("case3")).Return(1, nil)
				mockConn.EXPECT().Read(make([]byte, maxMsgSize)).Return(0, errors.New("err"))
			},
			expectedErr: "Client read data error",
		},
		{
			name: "Conn read error: count >= c.maxMessageSize",
			request: "case4",
			expected: []byte(""),
			prepare: func() {
				mockConn.EXPECT().SetDeadline(connDeadline).Return(nil)
				mockConn.EXPECT().Write([]byte("case4")).Return(1, nil)
				mockConn.EXPECT().Read(make([]byte, maxMsgSize)).Return(maxMsgSize + 1, nil)
			},
			expectedErr: "Small buffer size",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			res, err := client.Send([]byte(tt.request))
			if string(res) != string(tt.expected) {
				t.Errorf("expected res: %v \nactual res: %v", tt.expected, res)
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("expected err: %v \nactual err: %v", tt.expectedErr, err.Error())
			}
		})
	}
}
