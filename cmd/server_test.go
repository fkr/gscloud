package cmd

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/gridscale/gscloud/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockServer = gsclient.Server{
	Properties: gsclient.ServerProperties{
		ObjectUUID: "xxx",
	},
}

type mockServerOp struct {
	mock.Mock
}

func (o mockServerOp) GetServerList(ctx context.Context) ([]gsclient.Server, error) {
	args := o.Called()
	return args.Get(0).([]gsclient.Server), args.Error(1)
}

func (o mockServerOp) StartServer(ctx context.Context, id string) error {
	args := o.Called(id)
	return args.Error(0)
}

func (o mockServerOp) StopServer(ctx context.Context, id string) error {
	args := o.Called(id)
	return args.Error(0)
}

func (o mockServerOp) ShutdownServer(ctx context.Context, id string) error {
	args := o.Called(id)
	return args.Error(0)
}

func (o mockServerOp) DeleteServer(ctx context.Context, id string) error {
	args := o.Called(id)
	return args.Error(0)
}

func (o mockServerOp) CreateServer(ctx context.Context, body gsclient.ServerCreateRequest) (gsclient.ServerCreateResponse, error) {
	return gsclient.ServerCreateResponse{}, nil
}

func (o mockServerOp) CreateStorage(ctx context.Context, body gsclient.StorageCreateRequest) (gsclient.CreateResponse, error) {
	return gsclient.CreateResponse{}, nil
}

func (o mockServerOp) UpdateServer(ctx context.Context, id string, body gsclient.ServerUpdateRequest) error {
	return nil
}

func (o mockServerOp) GetDeletedServers(ctx context.Context) ([]gsclient.Server, error) {
	return []gsclient.Server{}, nil
}

func (o mockServerOp) GetServer(ctx context.Context, id string) (gsclient.Server, error) {
	return gsclient.Server{}, nil
}

func (o mockServerOp) GetServerEventList(ctx context.Context, id string) ([]gsclient.Event, error) {
	return []gsclient.Event{}, nil
}

func (o mockServerOp) GetServerMetricList(ctx context.Context, id string) ([]gsclient.ServerMetric, error) {
	return []gsclient.ServerMetric{}, nil
}

func (o mockServerOp) GetServersByLocation(ctx context.Context, id string) ([]gsclient.Server, error) {
	return []gsclient.Server{}, nil
}

func (o mockServerOp) IsServerOn(ctx context.Context, id string) (bool, error) {
	return false, nil
}

func Test_ServerCommmandDelete(t *testing.T) {
	type testCase struct {
		isSuccessful   bool
		expectedFatal  bool
		expectedOutput string
	}
	testCases := []testCase{
		{
			isSuccessful:   true,
			expectedFatal:  false,
			expectedOutput: "",
		},
		{
			isSuccessful:   false,
			expectedFatal:  false,
			expectedOutput: "",
		},
	}
	rt, _ = runtime.NewTestRuntime()
	for _, tc := range testCases {
		var fatal bool
		op := mockServerOp{}
		if tc.isSuccessful {
			op.On("DeleteServer", mock.Anything).Return(nil)
		} else {
			op.On("DeleteServer", mock.Anything).Return(errors.New("test"))
			log.StandardLogger().ExitFunc = func(int) { fatal = true }
		}
		rt.SetServerOperator(op)
		r, w, _ := os.Pipe()
		os.Stdout = w
		cmd := serverRmCmd.RunE
		cmd(new(cobra.Command), []string{"rm", mockServer.Properties.ObjectUUID})
		w.Close()
		out, _ := ioutil.ReadAll(r)
		assert.Equal(t, tc.expectedFatal, fatal)
		if tc.isSuccessful {
			assert.Equal(t, tc.expectedOutput, string(out))
		}
	}
}

func Test_ServerCommmandLs(t *testing.T) {
	type testCase struct {
		isSuccessful         bool
		isJSONFormat         bool
		expectedFatal        bool
		expectedPartOfOutput string
	}
	testCases := []testCase{
		{
			isSuccessful:         true,
			isJSONFormat:         false,
			expectedFatal:        false,
			expectedPartOfOutput: mockServer.Properties.ObjectUUID,
		},
		{
			isSuccessful:         true,
			isJSONFormat:         true,
			expectedFatal:        false,
			expectedPartOfOutput: mockServer.Properties.ObjectUUID,
		},
		{
			isSuccessful:         false,
			expectedFatal:        false,
			expectedPartOfOutput: "",
		},
	}
	rt, _ = runtime.NewTestRuntime()
	for _, tc := range testCases {
		var fatal bool
		rootFlags.json = tc.isJSONFormat

		op := mockServerOp{}
		if tc.isSuccessful {
			op.On("GetServerList", mock.Anything).Return([]gsclient.Server{mockServer}, nil)
		} else {
			op.On("GetServerList", mock.Anything).Return([]gsclient.Server{}, errors.New("test"))
			log.StandardLogger().ExitFunc = func(int) { fatal = true }
		}
		rt.SetServerOperator(op)
		r, w, _ := os.Pipe()
		os.Stdout = w
		cmd := serverLsCmd.RunE
		cmd(new(cobra.Command), []string{"ls"})
		w.Close()
		out, _ := ioutil.ReadAll(r)
		assert.Equal(t, tc.expectedFatal, fatal)
		if tc.isSuccessful {
			assert.Contains(t, string(out), tc.expectedPartOfOutput)
			if tc.isJSONFormat {
				assert.Equal(t, "[{", string(out[:2]))
				assert.Equal(t, "}]", string(out[len(out)-3:len(out)-1]))
			}
		}
	}
}

func Test_ServerCommmandOn(t *testing.T) {
	type testCase struct {
		isSuccessful   bool
		expectedFatal  bool
		expectedOutput string
	}
	testCases := []testCase{
		{
			isSuccessful:   true,
			expectedFatal:  false,
			expectedOutput: "",
		},
		{
			isSuccessful:   false,
			expectedFatal:  false,
			expectedOutput: "",
		},
	}
	rt, _ = runtime.NewTestRuntime()
	for _, tc := range testCases {
		var fatal bool
		op := mockServerOp{}
		if tc.isSuccessful {
			op.On("StartServer", mock.Anything).Return(nil)
		} else {
			op.On("StartServer", mock.Anything).Return(errors.New("test"))
			log.StandardLogger().ExitFunc = func(int) { fatal = true }
		}
		rt.SetServerOperator(op)
		r, w, _ := os.Pipe()
		os.Stdout = w
		cmd := serverOnCmd.RunE
		cmd(new(cobra.Command), []string{"on", mockServer.Properties.ObjectUUID})
		w.Close()
		out, _ := ioutil.ReadAll(r)
		assert.Equal(t, tc.expectedFatal, fatal)
		if tc.isSuccessful {
			assert.Equal(t, tc.expectedOutput, string(out))
		}
	}
}

func Test_ServerCommmandOff(t *testing.T) {
	type testCase struct {
		isSuccessful    bool
		isForceShutdown bool
		expectedFatal   bool
		expectedOutput  string
	}
	testCases := []testCase{
		{
			isSuccessful:    true,
			isForceShutdown: false,
			expectedFatal:   false,
			expectedOutput:  "",
		},
		{
			isSuccessful:    true,
			isForceShutdown: true,
			expectedFatal:   false,
			expectedOutput:  "",
		},
		{
			isSuccessful:   false,
			expectedFatal:  false,
			expectedOutput: "",
		},
	}
	rt, _ = runtime.NewTestRuntime()
	for _, tc := range testCases {
		var fatal bool
		serverFlags.forceShutdown = tc.isForceShutdown

		op := mockServerOp{}
		if tc.isForceShutdown {
			if tc.isSuccessful {
				op.On("StopServer", mock.Anything).Return(nil)
			} else {
				op.On("StopServer", mock.Anything).Return(errors.New("test"))
				log.StandardLogger().ExitFunc = func(int) { fatal = true }
			}
		} else {
			if tc.isSuccessful {
				op.On("ShutdownServer", mock.Anything).Return(nil)
			} else {
				op.On("ShutdownServer", mock.Anything).Return(errors.New("test"))
				log.StandardLogger().ExitFunc = func(int) { fatal = true }
			}
		}

		rt.SetServerOperator(op)
		r, w, _ := os.Pipe()
		os.Stdout = w
		cmd := serverOffCmd.RunE
		cmd(new(cobra.Command), []string{"off", mockServer.Properties.ObjectUUID})
		w.Close()
		out, _ := ioutil.ReadAll(r)
		assert.Equal(t, tc.expectedFatal, fatal)
		if tc.isSuccessful {
			assert.Equal(t, tc.expectedOutput, string(out))
		}
		op.AssertExpectations(t)
	}
}
