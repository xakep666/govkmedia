import QtQuick 2.0
import QtQuick.Controls 1.3
//import GoExtensions 1.0

ApplicationWindow {
    id: root

    width: 700
    height: 500
    color: "#ffffff"
    title: qsTr("GoVkMedia")
    Text {
        id: noauthtxt
        objectName: qsTr("noauthtxt")
        x: 0
        y: parent.height/2-height
        width: parent.width
        height: 27
        text: qsTr("Не авторизован")
        horizontalAlignment: Text.AlignHCenter
        font.pixelSize: 24
    }

    Button {
        id: authbtn
        objectName: qsTr("authbtn")
        x: parent.width/2-width/2
        y: parent.height/2+10
        width: 196
        height: 44
        text: qsTr("Авторизоваться")
        onClicked: {
            appEngine.showOauth()
        }
    }

    Image {
        id: avatar
        objectName: qsTr("avatar")
        x: 8
        y: 10
        width: 100
        height: 100
        visible:false
    }

    TabView{
        id: tabs
        objectName: qsTr("tabs")
        x: 210
        y: 10
        width: root.width-220
        height: root.height-20
        Tab {
            id: musictab
            title: "Музыка"
                TableView {
                    id: musictable
                    objectName: qsTr("musictable")
                    width:musictab.width
                    height:musictab.height
                    model: musiclist
                    selectionMode: SelectionMode.ContiguousSelection
                    TableViewColumn {
                        role: "artist"
                        title: "Автор"
                        width: 2*musictable.width/5
                    }
                    TableViewColumn {
                        role: "title"
                        title:"Название"
                        width: 2*musictable.width/5
                    }
                    TableViewColumn {
                        role: "duration"
                        title: "Длительность"
                        width: musictable.width/5
                    }
                    //some invisible fields to store meta-information
                    TableViewColumn {
                        role:"url"
                        visible: false
                    }
                    TableViewColumn {
                        role:"id"
                        visible:false
                    }
                    TableViewColumn {
                        role: "lyricsid"
                        visible:false
                    }
                    TableViewColumn {
                        role: "genre"
                        visible: false
                    }
                    TableViewColumn {
                        role: "album"
                        visible: false
                    }

                    ListModel {
                        id: musiclist
                        objectName: qsTr("musiclist")
                        //translate go struct to dict, direct call "append" from go not works
                        function appendStruct(m) {append(m)}
                        //wrap get method for calling from go
                        function goGet(i) {return get(i)}
                        property string selected: ""
                        function goFillSelected() {
                            selected=""
                            musictable.selection.forEach(function(index){
                                selected+=index+" "
                            })
                        }
                    }
                }
        }
        Tab {
            id: phototab
            title: "Фото"
        }
        Tab {
            id: videotab
            title: "Видео"
        }
        visible:false
    }

    Text {
        id: name
        objectName: "name"
        visible:false
        x: 8
        y: 116
        width: 199
        height: 44
        text: qsTr("")
        textFormat: Text.PlainText
        font.pixelSize: 12
    }

    Button {
        id: open
        objectName: qsTr("open")
        x: 90
        y: 166
        width: 118
        text: qsTr("Открыть")
        visible:false
        onClicked: {
            switch(tabs.currentIndex) {
            case 0:
                audioplayer.show()
            }
        }
    }

    Button {
        id: dlallbtn
        objectName: "dlallbtn"
        x: 90
        y: 224
        width: 118
        visible: false
        text: qsTr("Скачать все")
        onClicked: {
            switch(tabs.currentIndex) {
            case 0:
                appEngine.downloadAllMusic()
            }
        }
    }
    
    Button {
        id: dlselbtn
        objectName: "dlselbtn"
        x: 90
        y: 195
        width: 118
        visible: false
        text: qsTr("Скачать выбранное")
        onClicked: {
            switch(tabs.currentIndex) {
            case 0:
                appEngine.downloadSelectedMusic()
            }
        }
    }
}
