package downloaders

import (
    "github.com/nareix/curl"
    id3 "github.com/mikkyang/id3-go"
    v2 "github.com/mikkyang/id3-go/v2"
    "govkmedia/requestwrapper"
    "strconv"
)

type DownloadableMusic struct {
    DownloadableFile
    Artist string
    Title string
    Album string
    Genre string
    LyricsId float64
    AccessToken string
}

//reimplement download to insert info into id3 tags
func (dm *DownloadableMusic) Do() (resp curl.Response,err error) {
    resp,err=dm.DownloadableFile.Do()
    if err!=nil {
        return
    }
    mp3file,err:=id3.Open(dm.ActualPath())
    defer mp3file.Close()
    //add id3 info
    mp3file.SetTitle(dm.Title)
    mp3file.SetArtist(dm.Artist)
    if dm.Album!="" {mp3file.SetAlbum(dm.Album)}
    if dm.Genre!="" {mp3file.SetGenre(dm.Genre)}
    //download and insert lyrics
    if dm.LyricsId!=0 {
        rac:=requestwrapper.RequestAccesser{Token: dm.AccessToken}
        parms:=map[string]string{"lyrics_id":strconv.FormatFloat(dm.LyricsId,'f',-1,64)}
        vkresp,err:=rac.MakeRequest("audio.getLyrics",parms)
        if err!=nil {
            return resp,err
        }
        lyrics:=vkresp["response"].(map[string]interface{})["text"].(string)
        frt:=v2.V23FrameTypeMap["USLT"]
        lyricsFrame:=v2.NewTextFrame(frt,lyrics)
        mp3file.AddFrames(lyricsFrame)
    }
    return
}
