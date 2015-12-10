package root

// Executor executes commands as root user.
type Executor interface {
	// Exec runs the passed command and displays the
	// message msg in the root pop-up window.
	Exec(cmd, msg string) string
}
