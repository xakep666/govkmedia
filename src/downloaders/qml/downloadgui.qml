import QtQuick 2.0
import QtQuick.Controls 1.3
import QtQuick.Controls.Styles 1.3
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
    ProgressBar {
        id: totalprogressbar
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
        states: [
            State{
                name: "running"
                PropertyChanges {
                    target: startpausebtn
                    iconSource: "qrc:///icons/pause.png"
                }
            },
            State{
                name: "paused"
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
    }

    Text {
        id: collpseexpandtext
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
                PropertyChanges {
                    target: collpseexpandtext
                    isCollapsed: false
                    text: qsTr("Развернуть")
                }
                PropertyChanges {
                    target: iconrotate
                    angle:0
                }
                PropertyChanges {
                    target: dlwindow
                    width: dlwindow.origwidth
                }
            },
            State{
                name:"expanded"
                PropertyChanges {
                    target: collpseexpandtext
                    isCollapsed: true
                    text: qsTr("Свернуть")
                }
                PropertyChanges {
                    target: iconrotate
                    angle: -90
                }
                PropertyChanges {
                    target: dlwindow
                    width: dlwindow.origwidth+filetable.width+10
                }
            }
        ]
        transitions: [
            Transition {
                from: "collapsed"
                to: "expanded"
                RotationAnimation {
                    duration: 1000
                    target: collpseexpandicon
                    direction: RotationAnimation.Clockwise
                }
            },
            Transition {
                from: "expanded"
                to: "collapsed"
                RotationAnimation {
                    duration: 1000
                    target: collpseexpandicon
                    direction: RotationAnimation.Counterclockwise
                }
            }
        ]
        Image {
            id: collpseexpandicon
            x: 0
            y: 0
            width: collpseexpandtext.height
            height: collpseexpandtext.height
            source: "qrc:///icons/play.png"
            transform: Rotation {
                id: iconrotate
                axis {x:0;y:0;z:1}
                origin {
                    x:collpseexpandicon.width/2
                    y:collpseexpandicon.height/2
                }
                angle:0
            }
        }
        MouseArea {
            id: mousereact
            x: 0
            y: 0
            width: collpseexpandtext.width
            height: collpseexpandtext.height
        }
        TableView {
            id: filetable
            visible: !collpseexpandtext.isCollapsed
            x:0
            y:collpseexpandtext.height+5
            width:dlwindow.width-20
            height:250
        }
    }
}

