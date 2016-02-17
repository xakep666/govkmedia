import QtQuick 2.0
import QtQuick.Controls 1.3
ApplicationWindow {
    id: mainwindow
    Component {
        id: groupDelegate
        Row{
                id: group
                x: 0
                y: 0
                width: mainwindow.width-20
                height: 116

                Image {
                    id: groupavatar
                    x: 8
                    y: 8
                    width: 100
                    height: 100
                    source: grouplistview.model.get(styleData.row).avatar
                }

                Text {
                    id: groupname
                    x: 114
                    y: 8
                    width: group.width-125
                    height: 32
                    styleColor: "#0e01ec"
                    font.bold: true
                    font.pixelSize: 20
                    text: grouplistview.model.get(styleData.row).name
                }

                MouseArea {
                    x:114
                    y:8
                    width: group.width-125
                    height: 32
                    onClicked: grouplist.selectedgid=grouplistview.model.get(styleData.row).gid
                }

                Text {
                    id: text1
                    x: 114
                    y: 52
                    height: 13
                    color: "#dacfe3"
                    text: qsTr("Подписчиков:")
                    font.pixelSize: 12
                }

                Text {
                    id: subscribers
                    x: 200
                    y: 52
                    width: group.width-210
                    height: 13
                    font.pixelSize: 12
                    color: "#dacfe3"
                    text: grouplistview.model.get(styleData.row).name
                }
            }
    }
    ListView {
        x:10
        y:10
        width: mainwindow.width-20
        height: mainwindow.height-50
        id: grouplistview
        model: grouplist
        delegate: groupDelegate
    }
    ListModel {
        id: grouplist
        property int selectedgid:-1
        objectName: "grouplist"
        function appendStruct(m) { append(m) }
    }

    Button {
        id: cancelbtn
        y: mainwindow.height-30
        x: mainwindow.width-90
        text: qsTr("Отмена")
        onClicked: mainwindow.close()
    }

    Text {
        id: text2
        x: 10
        y: mainwindow.height-30
        width: 229
        height: 21
        text: qsTr("Выберете группу, кликнув по названию")
        verticalAlignment: Text.AlignVCenter
        horizontalAlignment: Text.AlignHCenter
        font.pixelSize: 12
    }
}
