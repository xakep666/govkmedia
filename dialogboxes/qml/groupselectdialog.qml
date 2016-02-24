import QtQuick 2.0
import QtQuick.Controls 1.3
import GoExtensions 1.0
ApplicationWindow {
    id: mainwindow
    property int selectedgid:-1
    signal selectedgidchanged
    width:640
    height:480
    Component {
        id: groupDelegate
        Rectangle{
                id: group
                objectName: "grouprow"
                x: 0
                y: 0
                width: mainwindow.width-20
                height: 116

                Image {
                    x: 8
                    y: 8
                    width: 100
                    height: 100
                    source: avatar
                }

                Text {
                    x: 114
                    y: 8
                    width: group.width-125
                    height: 32
                    styleColor: "#0e01ec"
                    font.bold: true
                    font.pixelSize: 20
                    text: name
                }

                MouseArea {
                    x:114
                    y:8
                    width: group.width-125
                    height: 32
                    onClicked: {
                        mainwindow.selectedgid=gid
                        mainwindow.close()
                    }
                }

                Text {
                    id: subscriberstxt
                    visible: group.ListView.view.subscribersvisible
                    objectName: "subrscriberstxt"
                    x: 114
                    y: 52
                    height: 13
                    color: "#dacfe3"
                    text: qsTr("Подписчиков:")
                    font.pixelSize: 12
                }

                Text {
                    visible: group.ListView.view.subscribersvisible
                    id: subscribersnum
                    objectName: "subscribers"
                    x: 200
                    y: 52
                    width: group.width-210
                    height: 13
                    font.pixelSize: 12
                    color: "#dacfe3"
                    text: subscribers
                }
            }
    }
    ListView {
        x:10
        y:50
        objectName: "grouplistview"
        width: mainwindow.width-20
        height: mainwindow.height-90
        id: grouplistview
        model: grouplist
        delegate: groupDelegate
        property bool subscribersvisible: true
    }
    ListModel {
        id: grouplist
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

    TextField {
        id: searchfield
        x: 10
        y: 22
        width: mainwindow.width-200
        height: 22
        placeholderText: qsTr("Поиск")
    }

    Button {
        id: searchbtn
        x: mainwindow.width-180
        y: 22
        text: qsTr("Найти")
        onClicked: engine.searchGroups(searchfield.text)
    }

    Button {
        id: mygrpbtn
        x: mainwindow.width-90
        y: 22
        text: qsTr("Мои группы")
        onClicked: engine.loadUserGroups(engine.myUid)
    }
    onClosing: selectedgidchanged()
}
