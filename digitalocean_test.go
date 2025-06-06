package digitalocean

import (
	"context"
	"strings"
	"testing"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/digitalocean"
)

func TestCaddyModule(t *testing.T) {
	p := Provider{}
	info := p.CaddyModule()

	if "dns.providers.digitalocean" != string(info.ID) {
		t.Errorf("Module ID should be correct, got: %s, want: %s", string(info.ID), "dns.providers.digitalocean")
	}

	instance := info.New()
	_, ok := instance.(*Provider)

	if !ok {
		t.Error("New should return a *Provider instance")
	}
}

func newTestDispenser(input string) *caddyfile.Dispenser {
	return caddyfile.NewTestDispenser(input)
}

func TestProvision(t *testing.T) {
	p := &Provider{Provider: &digitalocean.Provider{APIToken: "{env.DO_API_TOKEN}"}}

	// Set environment variable for testing placeholder replacement
	t.Setenv("DO_API_TOKEN", "actual_token")

	ctx, cancel := caddy.NewContext(caddy.Context{Context: context.TODO()})
	defer cancel()

	err := p.Provision(ctx)
	if err != nil {
		t.Errorf("Provision should not return an error, got: %v", err)
	}
	if "actual_token" != p.Provider.APIToken {
		t.Errorf("API token should be replaced, got: %s, want: %s", p.Provider.APIToken, "actual_token")
	}

	// Test with no placeholder
	pNoPlaceholder := &Provider{Provider: &digitalocean.Provider{APIToken: "direct_token"}}
	err = pNoPlaceholder.Provision(ctx)
	if err != nil {
		t.Errorf("Provision should not return an error with direct token, got: %v", err)
	}
	if "direct_token" != pNoPlaceholder.Provider.APIToken {
		t.Errorf("API token should remain unchanged, got: %s, want: %s", pNoPlaceholder.Provider.APIToken, "direct_token")
	}
}

func TestUnmarshalCaddyfile(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedErr      string
		expectedAPIToken string
	}{
		{
			name:             "token as argument",
			input:            "digitalocean mytoken",
			expectedAPIToken: "mytoken",
		},
		{
			name:             "token in block",
			input:            "digitalocean {\n api_token mytoken \n}",
			expectedAPIToken: "mytoken",
		},
		{
			name:        "token as argument and in block",
			input:       "digitalocean mytoken {\n api_token anothertoken \n}",
			expectedErr: "API token already set",
		},
		{
			name:        "missing token",
			input:       "digitalocean",
			expectedErr: "missing API token",
		},
		{
			name:        "missing token in block",
			input:       "digitalocean {\n}",
			expectedErr: "missing API token",
		},
		{
			name:        "api_token already set",
			input:       "digitalocean {\n api_token mytoken \n api_token anothertoken \n}",
			expectedErr: "API token already set",
		},
		{
			name:        "unrecognized subdirective",
			input:       "digitalocean {\n unknown_directive value \n}",
			expectedErr: "unrecognized subdirective 'unknown_directive'",
		},
		{
			name:        "token as arg with extra arg",
			input:       "digitalocean mytoken extraarg",
			expectedErr: "wrong argument count or unexpected line ending after 'extraarg'",
		},
		{
			name:        "api_token in block with extra arg",
			input:       "digitalocean {\n api_token mytoken extraarg \n}",
			expectedErr: "unrecognized subdirective 'extraarg'",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := &Provider{Provider: new(digitalocean.Provider)}
			d := newTestDispenser(tc.input)
			err := p.UnmarshalCaddyfile(d)

			if tc.expectedErr != "" {
				if err == nil {
					t.Errorf("Expected error: %s, but got no error", tc.expectedErr)
				} else if !strings.Contains(err.Error(), tc.expectedErr) {
					t.Errorf("Expected error to contain: %s, got: %v", tc.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
				if tc.expectedAPIToken != p.Provider.APIToken {
					t.Errorf("Expected API token: %s, got: %s", tc.expectedAPIToken, p.Provider.APIToken)
				}
			}
		})
	}
}
