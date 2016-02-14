package audioplayer

import (
    "gopkg.in/qml.v1"
    "log"
    "govkmedia/dialogboxes"
)

//go:generate genqrc qml icons

type PtrSong struct {
    Artist string
    Title string
    Url string
    LyricsId float64
    Duration float64
}

type Engine struct {
    qmlEngine *qml.Engine
    mainWindow *qml.Window
    playlist []PtrSong
    currentPlaying int
    accessToken string
    initialized bool
}

func (e *Engine) Initialize(token string, playlist []PtrSong) {
    e.qmlEngine=qml.NewEngine()
    component, err := e.qmlEngine.LoadFile("qrc:///qml/player.qml")
    if err != nil {
        log.Panicln(err)
    }
    context := e.qmlEngine.Context()
    e.mainWindow = component.CreateWindow(nil)
    context.SetVar("playerengine",e)
    e.accessToken=token
    e.playlist=playlist
    e.initialized=true
}

func (e *Engine) Show() {
    if !e.initialized {
         dialogboxes.ShowErrorDialog("Движок плеера не инициализирован")
         return
    }
    e.mainWindow.Show()
    e.currentPlaying=0
    e.StartPlay()
}
