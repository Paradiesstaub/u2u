package main

import (
	"log"

	"github.com/Paradiesstaub/u2u/golang/iso"
	"github.com/Paradiesstaub/u2u/golang/root"
	"gopkg.in/qml.v1"
)

var controler *Controler

func main() {
	if err := qml.Run(run); err != nil {
		log.Fatal(err)
	}
}

// run main QML event loop.
func run() error {
	engine := qml.NewEngine()
	component, err := engine.LoadFile("../qml/isocreator.qml")
	if err != nil {
		return err
	}
	ctrl := NewControl(component)
	context := engine.Context()
	context.SetVar("ctrl", &ctrl)
	//controler = NewControler(iso.FakeWriter{})
	controler = NewControler(iso.NewDummyWriterLinux(root.NewExecutor()))
	ctrl.Window.Show()
	ctrl.Window.Wait()
	return nil
}

// Control is the bridge object between QML and go.
// Methods are called from QML with a lowercase starting character, e.g:
// quit() instead of Quit().
type Control struct {
	Root   qml.Object
	Window *qml.Window
}

// NewControl creates a new Control object.
func NewControl(component qml.Object) Control {
	window := component.CreateWindow(nil)
	return Control{
		Root:   window.Root(),
		Window: window,
	}
}

// CreateUsb writes the passed iso to the device.
func (ctrl *Control) CreateUsb(iso, device string) {
	controler.CreateUsb(iso, device)
}

// Quit terminates the program.
func (ctrl *Control) Quit() {
	controler.Quit()
}
