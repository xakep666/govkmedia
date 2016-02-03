package audioplayer

import (
    _ "gopkg.in/qml.v1"
)

func (e *Engine) Play() {
    mwroot:=e.MainWindow.Root()
    playerobj:=mwroot.ObjectByName("mplayer")
    playerobj.Set("source",e.Playlist[e.CurrentPlaying].Url)
    playerobj.Call("play")
}

func (e *Engine) Pause() {
    mwroot:=e.MainWindow.Root()
    playerobj:=mwroot.ObjectByName("mplayer")
    playerobj.Call("pause")
}
