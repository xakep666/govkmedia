import QtQuick 2.0
import QtQuick.Controls 1.3
import QtQuick.Window 2.0
Window {
    width: 500
    height:100
    title: qsTr("Ошибка")
    id:errorwindow
    Image {
        id: errorimg
        x: 8
        y: 8
        width: 75
        height: 84
        source: "qrc:///icons/error.png"
    }

    Text {
        id: errortext
        objectName: qsTr("errortext")
        x: 89
        y: 8
        width: 403
        height: 53
        text: qsTr("")
        font.pixelSize: 12
    }

    Button {
        id: okbtn
        x: 412
        y: 69
        text: qsTr("OK")
        onClicked: errorwindow.close()
    }

}

