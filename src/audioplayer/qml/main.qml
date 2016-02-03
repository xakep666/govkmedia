import QtQuick 2.0
import QtQuick.Controls 1.3
import QtMultimedia 5.0
import QtQuick.Controls.Styles 1.3
import "tooltipcreator.js" as TooltipCreator
ApplicationWindow{
    id: playerwindow
    width: 600
    height: 600
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

    Text {
        id: duration
        x: 7
        y: 509
        width: 50
        height: 22
        text: qsTr("")
        font.pixelSize: 12
    }

    Text {
        id: eta
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
        iconSource: "../icons/prev.png"
        style: playerbtn
        onClicked: playerengine.prev()
    }

    ToolButton {
        id: playpause
        x: 55
        y: playerwindow.height-52
        state:"playing"
        iconSource: "../icons/pause.png"
        width: 42
        height: 39
        style: playerbtn
        states:[
            State{
                name:"playing"
                PropertyChanges {
                    target: playpause
                    iconSource: "../icons/pause.png"
                }
            },
            State{
                name: "paused"
                PropertyChanges {
                    target: playpause
                    iconSource: "../icons/play.png"
                }
            }
        ]
    }

    ToolButton {
        id: stop
        x: 103
        y: playerwindow.height-52
        width: 42
        height: 39
        iconSource: "../icons/stop.png"
        style: playerbtn
    }


    ToolButton {
        id: next
        x: 151
        y: playerwindow.height-52
        width: 42
        height: 39
        iconSource: "../icons/next.png"
        style: playerbtn
        onClicked: playerengine.next()
    }
    MediaPlayer {
        id: mplayer
        objectName: qsTr("mplayer")
        var secplayed=0
        //will call it from go
        function startPlay() {
            updateTitle(metaData.author,metaData.title)
            eta.text=totimetext(duration/1000)
            duration.text=totimetext(0)
            mplayer.play()
            playtimer.start()
            playslider.maximumValue=duration/1000
        }
        onPlaying: {
            playpause.state="playing"
        }
        onPaused: {
            playpause.state="paused"
            playtimer.stop()
        }
        onStopped: {
            playprogress.value=0
            eta.text=totimetext(duration/1000)
            duration.text=totimetext(0)
            playtimer.stop()
            updateTitle("","")
        }

    }

    ToolButton {
        id: muteunmute
        x: playerwindow.width-212
        y: playerwindow.height-52
        width: 42
        height: 39
        style: playerbtn
        states:[
            State{
                name: "unmuted"
                when: mplayer.muted==false
                PropertyChanges {
                    target: muteunmute
                    iconSource: "../icon/unmuted.png"
                }
            },
            State {
                name: "muted"
                when: mplayer.muted==true
                PropertyChanges{
                    target:muteunmute
                    iconSource: "../icon/muted.png"
                }
            }
        ]

    }
    function updateTitle(artist,title) {
        if (artist!=="" && title!=="")
            title=artist+" - "+title+" - VkAudio"
        else
            title="VkAudio"
    }
    function totimetext(secs) {
        d = Number(d)
        var h = Math.floor(d / 3600)
        var m = Math.floor(d % 3600 / 60)
        var s = Math.floor(d % 3600 % 60)
        return ((h > 0 ? h + ":" + (m < 10 ? "0" : "") : "") + m + ":" + (s < 10 ? "0" : "") + s)
    }

    ButtonStyle {
        id: playerbtn
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

    Slider {
        id: volslider
        x: 437
        y: playerwindow.height-43
        width: 153
        height: 22
        updateValueWhileDragging: true
        minimumValue: 0
        maximumValue: 1
        stepSize: 0.01
        onValueChanged: {
            var ttx = volslider.x+volslider.width/2
            var tty = volslider.y-10
            TooltipCreator.create(toString(volslider.value*100)+"%",volslider,{x:ttx,y:tty}).show()
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
        style: SliderStyle{
            groove: Rectangle {
                visible:false
            }
            handle: Rectangle {
                implicitHeight: playprogress.height+8
                implicitWidth: 10
                color: control.pressed ? "white" : "lightgray"
                border.color: "gray"
                border.width: 2
                radius:12
            }
        }
        onValueChanged: {
            //if we can not seek turn slider back
            if(mplayer.seekable) {
                mplayer.secplayed=playslider.value
                mplayer.seek(playslider.value*1000)
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
            mplayer.secplayed++
            playprogress.value+=mplayer.position/mplayer.duration
            duration.text=totimetext(mplayer.position/1000)
            eta.text=totimetext((mplayer.duration-mplayer.position)/1000)
        }
    }
}

