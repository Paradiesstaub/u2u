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
        //onClicked: Qt.quit()
        onClicked: b.quit()
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

    ComboBox {
        id: dropdown
        width: 240
        x: margin
        anchors.top: isoPath.bottom
        anchors.topMargin: margin
        onCurrentIndexChanged: {
            // ignore initial event when widget is setup
            if (currentText !== "") {
                runButtonVisibilityCheck()
            }
        }
    }

    // setDropdownItems gets a string passed as argument containing
    // a string-array (items are separated by ':'). The array is
    // than assigned to the ComboBox as model.
    function setDropdownItems(arr) {
        dropdown.model = arr.split(":")
    }

    Button {
        id: btn
        visible: false
        text: "Run"
        x: margin
        anchors.top: dropdown.bottom
        anchors.topMargin: margin
        onClicked: {
            b.createUsb(isoPath.text, dropdown.currentText)
        }
    }

    function runButtonVisibilityCheck() {
        btn.visible = b.checkShowRunButton(isoPath.text)
    }
}
