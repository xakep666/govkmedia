package audioplayer

import (
    "gopkg.in/qml.v1"
    "log"
    "dialogboxes"
)

//go:generate genqrc qml icons

type PtrSong struct {
    Artist string
    Title string
    Url string
    LyricsId float64
}

type Engine struct {
    QMLEngine *qml.Engine
    MainWindow *qml.Window
    Playlist []PtrSong
    CurrentPlaying int
    AccessToken string
    initialized bool
}

func (e *Engine) Initialize(token string, playlist []PtrSong) {
    e.QMLEngine=qml.NewEngine()
    component, err := e.QMLEngine.LoadFile("qrc:///qml/player.qml")
    if err != nil {
        log.Panicln(err)
    }
    context := e.QMLEngine.Context()
    e.MainWindow = component.CreateWindow(nil)
    context.SetVar("playerengine",e)
    e.Playlist=playlist
    e.initialized=true
}

func (e *Engine) Show() {
    if !e.initialized {
         dialogboxes.ShowErrorDialog("Движок плеера не инициализирован")
         return
    }
    e.MainWindow.Show()
    e.CurrentPlaying=0
    e.MainWindow.Root().Call("resetProgress")
    e.Play()
}
