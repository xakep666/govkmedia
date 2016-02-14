package audioplayer

import (
    "govkmedia/dialogboxes"
    "govkmedia/requestwrapper"
    "strconv"
    "time"
)

func (e *Engine) StartPlay() {
    mwroot:=e.mainWindow.Root()
    playerobj:=mwroot.ObjectByName("mplayer")
    curplaying:=e.playlist[e.currentPlaying]
    playerobj.Set("source",curplaying.Url)
    time.Sleep(time.Second) //need to avoid playback artefacts
    context:=e.qmlEngine.Context()
    context.SetVar("vkartist",curplaying.Artist)
    context.SetVar("vktitle",curplaying.Title)
    context.SetVar("vkduration",curplaying.Duration)
    e.LoadLyrics()
    playerobj.Call("resetProgress")
    playerobj.Call("play")
}

func (e *Engine) LoadLyrics() {
    curplaying:=e.playlist[e.currentPlaying]
    if curplaying.LyricsId==0 { return }
    ra:=requestwrapper.RequestAccesser{Token: e.accessToken}
    parms:=map[string]string{"lyrics_id":strconv.FormatFloat(curplaying.LyricsId,'f',-1,64)}
    resp,err:=ra.MakeRequest("audio.getLyrics",parms)
    if err!=nil {
        dialogboxes.ShowErrorDialog(err.Error())
        return
    }
    lyricsfield:=e.mainWindow.Root().ObjectByName("lyrics")
    lyrics:=resp["response"].(map[string]interface{})["text"].(string)
    lyricsfield.Set("text",lyrics)
}

func (e *Engine) Next() {
    if e.currentPlaying+1==len(e.playlist) { return }
    e.currentPlaying++
    e.StartPlay()
}

func (e *Engine) Prev() {
    if e.currentPlaying==0 { return }
    e.currentPlaying--
    e.StartPlay()
}
