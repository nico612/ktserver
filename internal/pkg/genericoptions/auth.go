package genericoptions

import "time"

// Auth is a generic type
type AuthOptions struct {
	Issuer         string        `json:"issuer,omitempty" mapstructure:"issuer,omitempty"`
	ExpireTime     time.Duration `json:"expire-time,omitempty" mapstructure:"expire-time,omitempty"`
	RefreshTime    time.Duration `json:"refresh-time,omitempty" mapstructure:"refresh-time,omitempty"`
	PrivateKeyFile string        `json:"private-key-file,omitempty" mapstructure:"private-key-file,omitempty"`
	PublicKeyFile  string        `json:"public-key-file,omitempty" mapstructure:"public-key-file,omitempty"`
}

// NewAuthOptions create a `zero` value instance.
func NewAuthOptions() *AuthOptions {
	return &AuthOptions{
		Issuer:         "ktserver",
		PrivateKeyFile: "./certs/auth.key",
		PublicKeyFile:  "./certs/auth.pub",
		ExpireTime:     7 * 24 * time.Hour,
		RefreshTime:    30 * 24 * time.Hour,
	}
}
