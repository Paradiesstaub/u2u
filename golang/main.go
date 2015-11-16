package main

import (
	"log"
	"teststuff/isocreator/golang/root"

	"github.com/Paradiesstaub/u2u/golang/iso"
	"gopkg.in/qml.v1"
)

var controler *Controler

func main() {
	if err := qml.Run(run); err != nil {
		log.Fatal(err)
	}
}

// run main QML event loop
func run() error {
	engine := qml.NewEngine()
	component, err := engine.LoadFile("../qml/isocreator.qml")
	if err != nil {
		return err
	}
	b := NewBridge(component)
	engine.Context().SetVar("b", &b)
	controler = NewControler(
		iso.NewDummyWriterLinux(root.NewExecutor()),
		//iso.FakeWriter{},
		b,
	)
	b.Window.Show()
	b.Window.Wait()
	return nil
}

// Bridge object between QML and go.
// Methods are called from QML with a lowercase starting character,
// e.g: quit() instead of Quit().
type Bridge struct {
	Root   qml.Object
	Window *qml.Window
}

// NewBridge creates a new Control object.
func NewBridge(component qml.Object) Bridge {
	window := component.CreateWindow(nil)
	return Bridge{
		Root:   window.Root(),
		Window: window,
	}
}

// CreateUsb writes the passed iso to the device.
func (b *Bridge) CreateUsb(iso, device string) {
	controler.CreateUsb(iso, device)
}

// CheckShowRunButton checks if the run button should be displayed in QML.
func (b *Bridge) CheckShowRunButton(iso string) bool {
	return controler.CheckShowRunButton(iso)
}

// Quit terminates the program.
func (b *Bridge) Quit() {
	controler.Quit()
}
