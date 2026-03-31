package auth

import (
	"testing"

	internalconfig "github.com/router-for-me/CLIProxyAPI/v6/internal/config"
)

func TestResolveOAuthUpstreamModel_SuffixPreservation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		aliases map[string][]internalconfig.OAuthModelAlias
		channel string
		input   string
		want    string
	}{
		{
			name: "numeric suffix preserved",
			aliases: map[string][]internalconfig.OAuthModelAlias{
				"gemini-cli": {{Name: "gemini-2.5-pro-exp-03-25", Alias: "gemini-2.5-pro"}},
			},
			channel: "gemini-cli",
			input:   "gemini-2.5-pro(8192)",
			want:    "gemini-2.5-pro-exp-03-25(8192)",
		},
		{
			name: "level suffix preserved",
			aliases: map[string][]internalconfig.OAuthModelAlias{
				"claude": {{Name: "claude-sonnet-4-5-20250514", Alias: "claude-sonnet-4-5"}},
			},
			channel: "claude",
			input:   "claude-sonnet-4-5(high)",
			want:    "claude-sonnet-4-5-20250514(high)",
		},
		{
			name: "no suffix unchanged",
			aliases: map[string][]internalconfig.OAuthModelAlias{
				"gemini-cli": {{Name: "gemini-2.5-pro-exp-03-25", Alias: "gemini-2.5-pro"}},
			},
			channel: "gemini-cli",
			input:   "gemini-2.5-pro",
			want:    "gemini-2.5-pro-exp-03-25",
		},
		{
			name: "config suffix takes priority",
			aliases: map[string][]internalconfig.OAuthModelAlias{
				"claude": {{Name: "claude-sonnet-4-5-20250514(low)", Alias: "claude-sonnet-4-5"}},
			},
			channel: "claude",
			input:   "claude-sonnet-4-5(high)",
			want:    "claude-sonnet-4-5-20250514(low)",
		},
		{
			name: "auto suffix preserved",
			aliases: map[string][]internalconfig.OAuthModelAlias{
				"gemini-cli": {{Name: "gemini-2.5-pro-exp-03-25", Alias: "gemini-2.5-pro"}},
			},
			channel: "gemini-cli",
			input:   "gemini-2.5-pro(auto)",
			want:    "gemini-2.5-pro-exp-03-25(auto)",
		},
		{
			name: "none suffix preserved",
			aliases: map[string][]internalconfig.OAuthModelAlias{
				"gemini-cli": {{Name: "gemini-2.5-pro-exp-03-25", Alias: "gemini-2.5-pro"}},
			},
			channel: "gemini-cli",
			input:   "gemini-2.5-pro(none)",
			want:    "gemini-2.5-pro-exp-03-25(none)",
		},
		{
			name: "kimi suffix preserved",
			aliases: map[string][]internalconfig.OAuthModelAlias{
				"kimi": {{Name: "kimi-k2.5", Alias: "k2.5"}},
			},
			channel: "kimi",
			input:   "k2.5(high)",
			want:    "kimi-k2.5(high)",
		},
		{
			name: "case insensitive alias lookup with suffix",
			aliases: map[string][]internalconfig.OAuthModelAlias{
				"gemini-cli": {{Name: "gemini-2.5-pro-exp-03-25", Alias: "Gemini-2.5-Pro"}},
			},
			channel: "gemini-cli",
			input:   "gemini-2.5-pro(high)",
			want:    "gemini-2.5-pro-exp-03-25(high)",
		},
		{
			name: "no alias returns empty",
			aliases: map[string][]internalconfig.OAuthModelAlias{
				"gemini-cli": {{Name: "gemini-2.5-pro-exp-03-25", Alias: "gemini-2.5-pro"}},
			},
			channel: "gemini-cli",
			input:   "unknown-model(high)",
			want:    "",
		},
		{
			name: "wrong channel returns empty",
			aliases: map[string][]internalconfig.OAuthModelAlias{
				"gemini-cli": {{Name: "gemini-2.5-pro-exp-03-25", Alias: "gemini-2.5-pro"}},
			},
			channel: "claude",
			input:   "gemini-2.5-pro(high)",
			want:    "",
		},
		{
			name: "empty suffix filtered out",
			aliases: map[string][]internalconfig.OAuthModelAlias{
				"gemini-cli": {{Name: "gemini-2.5-pro-exp-03-25", Alias: "gemini-2.5-pro"}},
			},
			channel: "gemini-cli",
			input:   "gemini-2.5-pro()",
			want:    "gemini-2.5-pro-exp-03-25",
		},
		{
			name: "incomplete suffix treated as no suffix",
			aliases: map[string][]internalconfig.OAuthModelAlias{
				"gemini-cli": {{Name: "gemini-2.5-pro-exp-03-25", Alias: "gemini-2.5-pro(high"}},
			},
			channel: "gemini-cli",
			input:   "gemini-2.5-pro(high",
			want:    "gemini-2.5-pro-exp-03-25",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mgr := NewManager(nil, nil, nil)
			mgr.SetConfig(&internalconfig.Config{})
			mgr.SetOAuthModelAlias(tt.aliases)

			auth := createAuthForChannel(tt.channel)
			got := mgr.resolveOAuthUpstreamModel(auth, tt.input)
			if got != tt.want {
				t.Errorf("resolveOAuthUpstreamModel(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func createAuthForChannel(channel string) *Auth {
	switch channel {
	case "gemini-cli":
		return &Auth{Provider: "gemini-cli"}
	case "claude":
		return &Auth{Provider: "claude", Attributes: map[string]string{"auth_kind": "oauth"}}
	case "vertex":
		return &Auth{Provider: "vertex", Attributes: map[string]string{"auth_kind": "oauth"}}
	case "codex":
		return &Auth{Provider: "codex", Attributes: map[string]string{"auth_kind": "oauth"}}
	case "aistudio":
		return &Auth{Provider: "aistudio"}
	case "antigravity":
		return &Auth{Provider: "antigravity"}
	case "qwen":
		return &Auth{Provider: "qwen"}
	case "iflow":
		return &Auth{Provider: "iflow"}
	case "kimi":
		return &Auth{Provider: "kimi"}
	default:
		return &Auth{Provider: channel}
	}
}

func TestOAuthModelAliasChannel_Kimi(t *testing.T) {
	t.Parallel()

	if got := OAuthModelAliasChannel("kimi", "oauth"); got != "kimi" {
		t.Fatalf("OAuthModelAliasChannel() = %q, want %q", got, "kimi")
	}
}

func TestApplyOAuthModelAlias_SuffixPreservation(t *testing.T) {
	t.Parallel()

	aliases := map[string][]internalconfig.OAuthModelAlias{
		"gemini-cli": {{Name: "gemini-2.5-pro-exp-03-25", Alias: "gemini-2.5-pro"}},
	}

	mgr := NewManager(nil, nil, nil)
	mgr.SetConfig(&internalconfig.Config{})
	mgr.SetOAuthModelAlias(aliases)

	auth := &Auth{ID: "test-auth-id", Provider: "gemini-cli"}

	resolvedModel := mgr.applyOAuthModelAlias(auth, "gemini-2.5-pro(8192)")
	if resolvedModel != "gemini-2.5-pro-exp-03-25(8192)" {
		t.Errorf("applyOAuthModelAlias() model = %q, want %q", resolvedModel, "gemini-2.5-pro-exp-03-25(8192)")
	}
}

func TestResolveOAuthUpstreamModelPool_MultipleModels(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		aliases map[string][]internalconfig.OAuthModelAlias
		channel string
		input   string
		want    []string
	}{
		{
			name: "two models with same alias returns pool",
			aliases: map[string][]internalconfig.OAuthModelAlias{
				"gemini-cli": {
					{Name: "gemini-2.5-pro-exp-03-25", Alias: "my-model"},
					{Name: "gemini-2.5-flash-8b", Alias: "my-model"},
				},
			},
			channel: "gemini-cli",
			input:   "my-model",
			want:    []string{"gemini-2.5-pro-exp-03-25", "gemini-2.5-flash-8b"},
		},
		{
			name: "single model returns nil (use resolveOAuthUpstreamModel)",
			aliases: map[string][]internalconfig.OAuthModelAlias{
				"gemini-cli": {
					{Name: "gemini-2.5-pro-exp-03-25", Alias: "my-model"},
				},
			},
			channel: "gemini-cli",
			input:   "my-model",
			want:    nil,
		},
		{
			name: "no alias returns nil",
			aliases: map[string][]internalconfig.OAuthModelAlias{
				"gemini-cli": {
					{Name: "gemini-2.5-pro-exp-03-25", Alias: "my-model"},
				},
			},
			channel: "gemini-cli",
			input:   "unknown-model",
			want:    nil,
		},
		{
			name: "suffix preserved in pool",
			aliases: map[string][]internalconfig.OAuthModelAlias{
				"claude": {
					{Name: "claude-sonnet-4-5-20250514", Alias: "cs4"},
					{Name: "claude-haiku-4-5-20250514", Alias: "cs4"},
				},
			},
			channel: "claude",
			input:   "cs4(high)",
			want:    []string{"claude-sonnet-4-5-20250514(high)", "claude-haiku-4-5-20250514(high)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mgr := NewManager(nil, nil, nil)
			mgr.SetConfig(&internalconfig.Config{})
			mgr.SetOAuthModelAlias(tt.aliases)

			auth := createAuthForChannel(tt.channel)
			got := mgr.resolveOAuthUpstreamModelPool(auth, tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("resolveOAuthUpstreamModelPool(%q) returned %v, want %v", tt.input, got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("resolveOAuthUpstreamModelPool(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestOAuthModelPoolRotation(t *testing.T) {
	t.Parallel()

	// Test that model pool rotates correctly using the rotation logic
	aliases := map[string][]internalconfig.OAuthModelAlias{
		"gemini-cli": {
			{Name: "gemini-2.5-pro", Alias: "my-model"},
			{Name: "gemini-2.5-flash", Alias: "my-model"},
			{Name: "gemini-2.5-flash-lite", Alias: "my-model"},
		},
	}

	mgr := NewManager(nil, nil, nil)
	mgr.SetConfig(&internalconfig.Config{})
	mgr.SetOAuthModelAlias(aliases)

	// Create auth with proper attributes for channel resolution
	auth := &Auth{
		ID:         "test-auth",
		Provider:   "gemini-cli",
		Attributes: map[string]string{},
	}

	// Debug: verify pool resolution
	pool := mgr.resolveOAuthUpstreamModelPool(auth, "my-model")
	if len(pool) != 3 {
		t.Fatalf("resolveOAuthUpstreamModelPool returned %d models, want 3: %v", len(pool), pool)
	}

	// Simulate 3 consecutive requests to test rotation
	// The pool rotation should work like OpenAICompat pool
	calls := make([]string, 0, 9)
	for i := 0; i < 9; i++ {
		models, pooled := mgr.preparedExecutionModels(auth, "my-model")
		if !pooled {
			t.Fatalf("expected pooled models, got %v", models)
		}
		if len(models) != 3 {
			t.Fatalf("expected 3 models in pool, got %d", len(models))
		}
		// First model is the one that will be tried
		calls = append(calls, models[0])
	}

	// Each model should be called 3 times (round-robin)
	expected := []string{"gemini-2.5-pro", "gemini-2.5-flash", "gemini-2.5-flash-lite",
		"gemini-2.5-pro", "gemini-2.5-flash", "gemini-2.5-flash-lite",
		"gemini-2.5-pro", "gemini-2.5-flash", "gemini-2.5-flash-lite"}

	if len(calls) != len(expected) {
		t.Fatalf("calls = %v, want %v", calls, expected)
	}
	for i := range calls {
		if calls[i] != expected[i] {
			t.Errorf("call %d = %q, want %q", i, calls[i], expected[i])
		}
	}
}
