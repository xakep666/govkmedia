import QtQuick 2.0
import QtQuick.Controls 1.3
import QtQuick.Dialogs 1.2
ApplicationWindow {
    width: 600
    height: 100
    minimumHeight: height
    maximumHeight: height
    minimumWidth: width
    maximumWidth: width
    id: selectwindow

    signal gotpath
    property bool cancelled: false
    function closeAndSendPath() {
        gotpath()
        selectwindow.close()
    }
    function setupForFolder() {
        selectwindow.title=qsTr("Выберете папку")
        chosentxt.text=qsTr("Выбрана папка:")
        filedialog.selectFolder=true
    }

    Button {
        id: okbtn
        x: 424
        y: 70
        text: qsTr("OK")
        onClicked: closeAndSendPath()
    }

    Button {
        id: selectbtn
        x: 10
        y: 70
        text: qsTr("Выбрать")
        isDefault: true
        onClicked: filedialog.open()
    }

    Button {
        id: cancelbtn
        x: 510
        y: 70
        text: qsTr("Отмена")
        onClicked: {
            selectwindow.cancelled=true
            closeAndSendPath()
        }
    }

    Image {
        id: image1
        x: 10
        y: 12
        width: 45
        height: 47
        source: "qrc:///icons/info.png"
    }

    Text {
        id: chosentxt
        x: 64
        y: 12
        width: 526
        height: 15
        horizontalAlignment: Text.AlignLeft
        font.pixelSize: 12
    }

    Text {
        id: pathtxt
        x: 64
        y: 33
        width: 526
        height: 26
        wrapMode: Text.WordWrap
        font.pixelSize: 12
    }
    FileDialog {
        id: filedialog
        objectName: "filedialog"
        onAccepted: {
            pathtxt.text=""
            if(fileUrl!=="") {
                pathtxt.text=fileUrl.toString().replace("file://","")
            } else {
                for (var url in fileUrls) {
                    url=url.toString().replace("file://","")
                    pathtxt.text+=url.split("/\\").pop()+""
                }
            }
        }
    }
}

