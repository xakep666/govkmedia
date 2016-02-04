import QtQuick 2.0
import QtQuick.Controls 1.3
import QtMultimedia 5.0
import QtQuick.Controls.Styles 1.3

ApplicationWindow{
    id: playerwindow
    width: 600
    height: 600
    maximumWidth: width
    maximumHeight: height
    minimumWidth: width
    minimumHeight: height
    title: "VkAudio"
    TextArea {
        id: lyrics
        objectName: qsTr("lyrics")
        x: 7
        y: 10
        width: playerwindow.width-17
        height: playerwindow.height-107
    }

    ProgressBar {
        id: bufferprogress
        x: 63
        y: playerwindow.width-91
        width: playerwindow.width-126
        height: 23
        style: ProgressBarStyle {
            background: Rectangle {
                radius: 2
                color: "lightgray"
                border.color: "gray"
                border.width: 1
            }
            progress: Rectangle{
                color: "orange"
                border.color: "steelblue"
            }
        }
    }

    ProgressBar {
        id: playprogress
        x: 63
        y: playerwindow.width-91
        width: playerwindow.width-126
        height: 23
        style: ProgressBarStyle {
            background: Rectangle {
                radius: 2
                color: "transparent"
                border.color: "gray"
                border.width: 1
            }
            progress: Rectangle{
                color: "lightsteelblue"
                border.color: "steelblue"
            }
        }
    }


    Text {
        id: durationText
        x: 7
        y: 509
        width: 50
        height: 22
        text: qsTr("")
        font.pixelSize: 12
    }

    Text {
        id: etaText
        x: playerwindow.width-57
        y: playerwindow.height-91
        width: 47
        height: 23
        text: qsTr("")
        font.pixelSize: 12
    }

    ToolButton {
        id: prev
        x: 7
        y: playerwindow.height-52
        width: 42
        height: 39
        iconSource: "qrc:///icons/prev.png"
        style: ButtonStyle {
            background:
                Rectangle {
                    border.width: control.activeFocus ? 2 : 1
                    border.color: "#888"
                    radius: 4
                    gradient: Gradient {
                        GradientStop { position: 0 ; color: control.pressed ? "#ccc" : "#eee" }
                        GradientStop { position: 1 ; color: control.pressed ? "#aaa" : "#ccc" }
                    }
            }
        }
        onClicked: playerengine.prev()
    }

    ToolButton {
        id: playpause
        x: 55
        y: playerwindow.height-52
        state:"playing"
        iconSource: "qrc:///icons/pause.png"
        width: 42
        height: 39
        style: ButtonStyle {
                    background:
                        Rectangle {
                            border.width: control.activeFocus ? 2 : 1
                            border.color: "#888"
                            radius: 4
                            gradient: Gradient {
                                GradientStop { position: 0 ; color: control.pressed ? "#ccc" : "#eee" }
                                GradientStop { position: 1 ; color: control.pressed ? "#aaa" : "#ccc" }
                            }
                    }
        }
        states:[
            State{
                name:"playing"
                PropertyChanges {
                    target: playpause
                    iconSource: "qrc:///icons/pause.png"
                }
            },
            State{
                name: "paused"
                PropertyChanges {
                    target: playpause
                    iconSource: "qrc:///icons/play.png"
                }
            }
        ]
        onClicked: {
            if(state==="playing") {
                state="paused"
                mplayer.pause()
            } else {
                state="playing"
                mplayer.play()
            }
        }
    }

    ToolButton {
        id: stop
        x: 103
        y: playerwindow.height-52
        width: 42
        height: 39
        iconSource: "qrc:///icons/stop.png"
        style: ButtonStyle {
            background:
                Rectangle {
                    border.width: control.activeFocus ? 2 : 1
                    border.color: "#888"
                    radius: 4
                    gradient: Gradient {
                        GradientStop { position: 0 ; color: control.pressed ? "#ccc" : "#eee" }
                        GradientStop { position: 1 ; color: control.pressed ? "#aaa" : "#ccc" }
                    }
            }
        }
        onClicked: {
            mplayer.stop()
            playpause.state="paused"
        }
    }


    ToolButton {
        id: next
        x: 151
        y: playerwindow.height-52
        width: 42
        height: 39
        iconSource: "qrc:///icons/next.png"
        style: ButtonStyle {
            background:
                Rectangle {
                    border.width: control.activeFocus ? 2 : 1
                    border.color: "#888"
                    radius: 4
                    gradient: Gradient {
                        GradientStop { position: 0 ; color: control.pressed ? "#ccc" : "#eee" }
                        GradientStop { position: 1 ; color: control.pressed ? "#aaa" : "#ccc" }
                    }
            }
        }
        onClicked: playerengine.next()
    }
    MediaPlayer {
        id: mplayer
        objectName: qsTr("mplayer")
        autoLoad: true
        function resetProgress() {
            etaText.text=totimetext(Math.floor(vkduration))
            durationText.text=totimetext(0)
            bufferprogress.value=0
            playprogress.value=0
            playslider.maximumValue=Math.floor(vkduration)
            playslider.value=0
        }
        onPlaying: {
            playtimer.start()
            playpause.state="playing"
            updateTitle(vkartist,vktitle)
        }
        onPaused: {
            playpause.state="paused"
        }
        onStopped: {
            playpause.state="paused"
            resetProgress()
        }
    }

    ToolButton {
        id: muteunmute
        x: playerwindow.width-212
        y: playerwindow.height-52
        width: 42
        height: 39
        style: ButtonStyle {
            background:
                Rectangle {
                    border.width: control.activeFocus ? 2 : 1
                    border.color: "#888"
                    radius: 4
                    gradient: Gradient {
                        GradientStop { position: 0 ; color: control.pressed ? "#ccc" : "#eee" }
                        GradientStop { position: 1 ; color: control.pressed ? "#aaa" : "#ccc" }
                    }
            }
        }
        states:[
            State{
                name: "unmuted"
                when: mplayer.muted==false
                PropertyChanges {
                    target: muteunmute
                    iconSource: "qrc:///icons/unmuted.png"
                }
            },
            State {
                name: "muted"
                when: mplayer.muted==true
                PropertyChanges{
                    target:muteunmute
                    iconSource: "qrc:///icons/muted.png"
                }
            }
        ]
        onClicked: {
            if(mplayer.muted) {
                mplayer.muted=false
                state="unmuted"
            } else {
                mplayer.muted=true
                state="muted"
            }
        }
    }
    function updateTitle(artist,title) {
        if (artist!=="" && title!=="")
            playerwindow.title=artist+" - "+title+" - VkAudio"
        else
            playerwindow.title="VkAudio"
    }
    function totimetext(d) {
        d = Number(d)
        var h = Math.floor(d / 3600)
        if (h==NaN) h=0
        var m = Math.floor(d % 3600 / 60)
        if (m==NaN) m=0
        var s = Math.floor(d % 3600 % 60)
        if (m==NaN) m=0
        return ((h > 0 ? h + ":" + (m < 10 ? "0" : "") : "") + m + ":" + (s < 10 ? "0" : "") + s)
    }

    Slider {
        id: volslider
        x: playerwindow.width-163
        y: playerwindow.height-43
        width: 153
        height: 22
        updateValueWhileDragging: true
        minimumValue: 0
        maximumValue: 1
        stepSize: 0.01
        value:1
        onValueChanged: {
            mplayer.volume=volslider.value
        }
    }
    Slider {
        id: playslider
        x: 63
        y: playerwindow.width-91
        width: playerwindow.width-126
        height: 23
        minimumValue: 0
        stepSize: 1
        updateValueWhileDragging: false
        property bool doseek: true
        style: SliderStyle{
            groove: Rectangle {
                visible:false
            }
            handle: Rectangle {
                implicitHeight: playprogress.height+8
                implicitWidth: 10
                color: "lightgray"
                border.color: "gray"
                border.width: 2
                radius:12
            }
        }
        onValueChanged: {
            //needed for updates by timer
            if (!doseek) return
            //if we can not seek turn slider back
            if(mplayer.seekable) {
                mplayer.seek(playslider.value*1000)
                playprogress.value=mplayer.position/mplayer.duration
            } else {
                playslider.value=0
            }
        }
    }

    Timer {
        id: playtimer
        repeat: true
        interval: 1000
        onTriggered: {
            bufferprogress.value=mplayer.bufferProgress
            playprogress.value=mplayer.position/mplayer.duration
            if (mplayer.duration-mplayer.position<1000) playerengine.next() //MediaPlayer doesn`t have signal like onEnded
            durationText.text=totimetext(mplayer.position/1000)
            etaText.text=totimetext((mplayer.duration-mplayer.position)/1000)
            playslider.doseek=false
            playslider.value=mplayer.position/1000
            playslider.doseek=true
        }
    }
    onClosing: {
        mplayer.stop()
    }
}

