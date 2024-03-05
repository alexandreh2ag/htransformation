package rename_test

import (
	"context"
	"net/http"
	"regexp"
	"testing"

	"github.com/alexandreh2ag/htransformation/pkg/handler/rename"
	"github.com/alexandreh2ag/htransformation/pkg/tests/assert"
	"github.com/alexandreh2ag/htransformation/pkg/tests/require"
	"github.com/alexandreh2ag/htransformation/pkg/types"
)

func TestRenameHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		rule           types.Rule
		requestHeaders map[string]string
		wantHeaders    map[string]string
		wantHost       string
	}{
		{
			name: "no transformation",
			rule: types.Rule{
				Header: "not-existing",
			},
			requestHeaders: map[string]string{
				"Foo": "Bar",
			},
			wantHeaders: map[string]string{
				"Foo": "Bar",
			},
		},
		{
			name: "one transformation",
			rule: types.Rule{
				Header: "Test",
				Value:  "X-Testing",
			},
			requestHeaders: map[string]string{
				"Foo":  "Bar",
				"Test": "Success",
			},
			wantHeaders: map[string]string{
				"Foo":       "Bar",
				"X-Testing": "Success",
			},
		},
		{
			name: "override host request",
			rule: types.Rule{
				Header: "X-Host",
				Value:  "Host",
			},
			requestHeaders: map[string]string{
				"Foo":    "Bar",
				"X-Host": "example.com",
			},
			wantHeaders: map[string]string{
				"Foo":  "Bar",
				"Host": "example.com",
			},
			wantHost: "example.com",
		},
		{
			name: "Deletion",
			rule: types.Rule{
				Header: "Test",
			},
			requestHeaders: map[string]string{
				"Foo":  "Bar",
				"Test": "Success",
			},
			wantHeaders: map[string]string{
				"Foo":  "Bar",
				"Test": "",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
			require.NoError(t, err)

			for hName, hVal := range test.requestHeaders {
				req.Header.Add(hName, hVal)
			}

			test.rule.Regexp = regexp.MustCompile(test.rule.Header)

			rename.Handle(nil, req, test.rule)

			for hName, hVal := range test.wantHeaders {
				assert.Equal(t, hVal, req.Header.Get(hName))
			}
			if test.wantHost != "" {
				assert.Equal(t, test.wantHost, req.Host)
			}
		})
	}
}

func TestValidation(t *testing.T) {
	testCases := []struct {
		name    string
		rule    types.Rule
		wantErr bool
	}{
		{
			name:    "no rules",
			wantErr: true,
		},
		{
			name: "missing header value",
			rule: types.Rule{
				Header: ".",
				Type:   types.Rename,
			},
			wantErr: true,
		},
		{
			name: "invalid regexp",
			rule: types.Rule{
				Header: "(",
				Type:   types.Rename,
			},
			wantErr: true,
		},
		{
			name: "valid rule",
			rule: types.Rule{
				Header: "not-empty",
				Value:  "not-empty",
				Type:   types.Rename,
			},
			wantErr: false,
		},
	}

	for _, test := range testCases {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := rename.Validate(test.rule)
			t.Log(err)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
