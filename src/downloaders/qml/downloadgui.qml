import QtQuick 2.0
import QtQuick.Controls 1.3
import QtQuick.Controls.Styles 1.3
//import GoExtensions 1.0
ApplicationWindow {
    id: dlwindow
    title: qsTr("Загрузка")
    width: 400
    height: 90
    property int origwidth: 400
    property int origheight: 90
    maximumWidth: width
    minimumWidth: width
    maximumHeight: height
    minimumHeight: height
    onClosing: { engine.destruct() }
    ProgressBar {
        id: totalprogressbar
        objectName: "totalprogressbar"
        x: 10
        y: 39
        width: 326
        height: 23
    }

    ToolButton {
        id: startpausebtn
        x: 342
        y: 39
        width: 23
        height: 23
        state: "running"
        property bool isRunning: true
        onClicked: {
            if (isRunning) {
                engine.pause()
                isRunning=false
            } else {
                engine.resume()
                isRunning=true
            }
        }

        states: [
            State{
                name: "running"
                when: startpausebtn.isRunning
                PropertyChanges {
                    target: startpausebtn
                    iconSource: "qrc:///icons/pause.png"
                }
            },
            State{
                name: "paused"
                when: !startpausebtn.isRunning
                PropertyChanges {
                    target: startpausebtn
                    iconSource: "qrc:///icons/play.png"
                }
            }
        ]
    }

    ToolButton {
        id: cancelbtn
        x: 371
        y: 39
        width: 24
        height: 23
        iconSource: "qrc:///icons/stop.png"
        onClicked: {
            engine.cancel()
            startpausebtn.isRunning=false
        }
    }

    Text {
        id: text1
        x: 10
        y: 18
        text: qsTr("Загрузка: ")
        font.pixelSize: 12
    }

    Text {
        id: percents
        objectName: "percents"
        x: 81
        y: 18
        width: 28
        height: 15
        font.pixelSize: 12
        function setPercents(p) {
            text= (isNaN(p) ? 0 : Math.ceil(p))+"%"
        }
    }

    Text {
        id: collapseexpandtext
        x: 10
        y: 68
        width: 90
        height: 15
        property bool isCollapsed: true
        state:"collapsed"
        horizontalAlignment: Text.AlignRight
        font.underline: true
        verticalAlignment: Text.AlignBottom
        font.pixelSize: 12
        states: [
            State {
                name:"collapsed"
                when:collapseexpandtext.isCollapsed
                PropertyChanges {
                    target: collapseexpandtext
                    text: qsTr("Развернуть")
                }
                PropertyChanges {
                    target: iconrotate
                    angle:0
                }
                PropertyChanges {
                    target: dlwindow
                    height: dlwindow.origheight
                    maximumHeight: dlwindow.origheight
                    minimumHeight: dlwindow.origheight
                }
            },
            State{
                name:"expanded"
                when:!collapseexpandtext.isCollapsed
                PropertyChanges {
                    target: collapseexpandtext
                    text: qsTr("Свернуть")
                }
                PropertyChanges {
                    target: iconrotate
                    angle: 90
                }
                PropertyChanges {
                    target: dlwindow
                    height: dlwindow.origheight+filetable.height+10
                    maximumHeight: dlwindow.origheight+filetable.height+10
                    minimumHeight: dlwindow.origheight+filetable.height+10
                }
            }
        ]
        transitions: [
            Transition {
                from: "collapsed"
                to: "expanded"
                RotationAnimation {
                    duration: 500
                    direction: RotationAnimation.Clockwise
                }
            },
            Transition {
                from: "expanded"
                to: "collapsed"
                RotationAnimation {
                    duration: 500
                    direction: RotationAnimation.Counterclockwise
                }
            }
        ]

        Image {
            id: collapseexpandicon
            x: 0
            y: 0
            width: collapseexpandtext.height
            height: collapseexpandtext.height
            source: "qrc:///icons/play.png"
            transform: Rotation {
                id: iconrotate
                axis {x:0;y:0;z:1}
                origin {
                    x:collapseexpandicon.width/2
                    y:collapseexpandicon.height/2
                }
                angle:0
            }
        }
        MouseArea {
            id: mousereact
            x: 0
            y: 0
            width: collapseexpandtext.width
            height: collapseexpandtext.height
            onClicked: collapseexpandtext.isCollapsed=!collapseexpandtext.isCollapsed
        }
        TableView {
            id: filetable
            objectName: "filetable"
            visible: !collapseexpandtext.isCollapsed
            x:0
            y:collapseexpandtext.height+5
            width:dlwindow.width-20
            height:250
            model: dllist
            TableViewColumn {
                role: "fname"
                title: qsTr("Имя файла")
                width:filetable.width*3/5
            }
            TableViewColumn {
                role: "dlspeed"
                title: qsTr("Скорость")
                width:filetable.width*1/5
            }
            TableViewColumn {
                role: "dlprogress"
                id: dlprogressColumn
                title: qsTr("Прогресс")
                width:filetable.width*1/5
                delegate: ProgressBar {
                        id: dlprogressBar
                        width: dlprogressColumn.width
                        value: filetable.model.get(styleData.row).dlprogress
                    }
            }
            ListModel {
                id: dllist
                objectName: "dllist"
                function appendStruct(m) {
                    append(m)
                }
                function back() {
                    return get(count-1)
                }
            }
        }
    }
}

