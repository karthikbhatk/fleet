// Package execuser is used to run applications from a high privilege user (root on Unix,
// SYSTEM service on Windows) as the current login user.
package execuser

type eopts struct {
	env        [][2]string
	stderrPath string //nolint:unused
}

// Option allows configuring the application.
type Option func(*eopts)

// WithEnv sets environment variables for the application.
func WithEnv(name, value string) Option {
	return func(a *eopts) {
		a.env = append(a.env, [2]string{name, value})
	}
}

// Run runs an application as the current login user.
// It assumes the caller is running with high privileges (root on Unix, SYSTEM on Windows).
//
// It returns after starting the child process.
func Run(path string, opts ...Option) error {
	var o eopts
	for _, fn := range opts {
		fn(&o)
	}
	return run(path, o)
}
