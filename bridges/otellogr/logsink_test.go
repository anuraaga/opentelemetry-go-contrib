// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
package otellogr

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
	"go.opentelemetry.io/otel/log/global"
)

type mockLoggerProvider struct {
	embedded.LoggerProvider
}

func (mockLoggerProvider) Logger(name string, options ...log.LoggerOption) log.Logger {
	return nil
}

func TestNewConfig(t *testing.T) {
	customLoggerProvider := mockLoggerProvider{}

	for _, tt := range []struct {
		name    string
		options []Option

		wantConfig config
	}{
		{
			name: "with no options",

			wantConfig: config{
				provider: global.GetLoggerProvider(),
			},
		},
		{
			name: "with a custom instrumentation scope",
			options: []Option{
				WithVersion("42.0"),
			},

			wantConfig: config{
				version:  "42.0",
				provider: global.GetLoggerProvider(),
			},
		},
		{
			name: "with a custom logger provider",
			options: []Option{
				WithLoggerProvider(customLoggerProvider),
			},

			wantConfig: config{
				provider: customLoggerProvider,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantConfig, newConfig(tt.options))
		})
	}
}

func TestNewLogSink(t *testing.T) {
	const name = "test_logsink"
	provider := global.GetLoggerProvider()

	for _, tt := range []struct {
		name       string
		options    []Option
		wantLogger log.Logger
	}{
		{
			name:       "with default options",
			wantLogger: provider.Logger(name),
		},
		{
			name: "with version and schema URL",
			options: []Option{
				WithVersion("1.0"),
				WithSchemaURL("https://example.com"),
			},
			wantLogger: provider.Logger(name,
				log.WithInstrumentationVersion("1.0"),
				log.WithSchemaURL("https://example.com"),
			),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			hook := NewLogSink(name, tt.options...)
			assert.NotNil(t, hook)
			assert.Equal(t, tt.wantLogger, hook.logger)
		})
	}
}