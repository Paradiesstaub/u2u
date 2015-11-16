import QtQuick 2.2
import QtQuick.Window 2.0
import QtQuick.Controls 1.1

// widget doc: http://doc.qt.io/qt-5/qtquickcontrols-index.html

Window {
    width: 640
    height: 480

    property int margin: 16

    /*
    MouseArea {
        anchors.fill: parent
        onClicked: ctrl.quit()
    }
    */

    Text {
        id: info
        x: margin
        y: margin
        text: "Create USB Start-Up Disk."
    }
    TextField {
        id: isoPath
        width: 400
        x: margin
        anchors.top: info.bottom
        anchors.topMargin: margin
        placeholderText: "path to iso file"
        onTextChanged: runButtonVisibilityCheck()
    }

    TextField {
        id: devicePath
        width: 240
        x: margin
        anchors.top: isoPath.bottom
        anchors.topMargin: margin
        placeholderText: "path to the device, e.g: /dev/sdc"
        onTextChanged: runButtonVisibilityCheck()
    }

    function runButtonVisibilityCheck() {
        btn.visible = isoPath.text.length > 0 && devicePath.text.length > 0
    }

    Button {
        id: btn
        visible: false
        text: "Run"
        x: margin
        anchors.top: devicePath.bottom
        anchors.topMargin: margin
        onClicked: {
            ctrl.createUsb(isoPath.text, devicePath.text)
        }
    }
}
