package root

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Executor executes commands as root user.
type Executor interface {
	// Exec runs the passed command and displays the
	// message msg in the root pop-up window.
	Exec(cmd, msg string) string
}

// NewExecutor is a factory method to return an Executor
// implementation (Linux only right now).
func NewExecutor() Executor {
	desktop := os.Getenv("XDG_CURRENT_DESKTOP") // KDE, Unity, GNOME, LXDE
	switch desktop {
	case "KDE":
		return new(kdeExecutor)
	case "Unity", "GNOME":
		return new(gnomeExecutor)
	default:
		log.Fatalf("Plattform %s is not supported.", desktop)
		return nil
	}
}

type kdeExecutor struct{}

// Exec runs the passed command cmd using the sh shell and kdesudo.
func (r kdeExecutor) Exec(cmd, message string) string {
	checkIfInstalled("kdesudo")
	b, err := exec.Command("/bin/sh", "-c",
		fmt.Sprintf("kdesudo -d --comment '%s' %s", message, cmd)).Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

type gnomeExecutor struct{}

// Exec runs the passed command cmd using the sh shell and gksudo.
func (r gnomeExecutor) Exec(cmd, message string) string {
	checkIfInstalled("gksudo")
	b, err := exec.Command("/bin/sh", "-c",
		fmt.Sprintf("gksudo --message '%s' %s", message, cmd)).Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func checkIfInstalled(app string) {
	cmd := exec.Command("which", app)
	if err := cmd.Run(); err != nil {
		log.Fatalf("Application '%v' is not installed!\n", app)
	}
}
