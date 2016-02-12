package main

import (
    "strconv"
    "time"
    "dialogboxes"
    "audioplayer"
    "requestwrapper"
)

func (ae *AppEngine) loadAvatar() {
    params:=map[string]string{
        "user_ids":strconv.Itoa(ae.UserId),
        "fields":"photo_100",
    }
    mwroot:=ae.MainWindow.Root()
    avatar:=mwroot.ObjectByName("avatar")
    resp,err:=ae.MakeRequest("users.get",params)
    if err!=nil {
        dialogboxes.ShowErrorDialog("Не удалось загрузить аватар: "+err.Error())
        return
    }
    /*in response we get json like this
      response: [ {<user1data>},{<user2data>},... ]
      here we extract first item of array and converting to map[string]...
      to get info about user
    */
    resp=resp["response"].([]interface{})[0].(map[string]interface{})
    avatar.Set("source",resp["photo_100"])
}

func (ae *AppEngine) loadName() {
    params:=map[string]string{
        "user_ids":strconv.Itoa(ae.UserId),
        "fields":"",
    }
    mwroot:=ae.MainWindow.Root()
    namefield:=mwroot.ObjectByName("name")
    resp,err:=ae.MakeRequest("users.get",params)
    if err!=nil {
        dialogboxes.ShowErrorDialog("Не удалось загрузить имя: "+err.Error())
        return
    }
    resp=resp["response"].([]interface{})[0].(map[string]interface{})
    namestr:=resp["first_name"].(string)
    namestr+=" "
    namestr+=resp["last_name"].(string)
    namefield.Set("text",namestr)
}

type MusicItem struct {
    Artist string
    Title string
    Duration string
    Id float64
    Url string
    LyricsId float64
    Genre string
    Album string
}

func (ae *AppEngine) loadAudios(uid int) {
    params:=map[string]string{
        "owner_id":strconv.Itoa(uid),
        "count":"6000",
    }
    mwroot:=ae.MainWindow.Root()
    model:=mwroot.ObjectByName("musiclist")
    resp,err:=ae.MakeRequest("audio.get",params)
    if err!=nil {
        dialogboxes.ShowErrorDialog("Не удалось загрузить аудиозаписи: "+err.Error())
        return
    }
    model.Call("clear")
    content:=resp["response"].(map[string]interface{})
    items:=content["items"].([]interface{})
    tmplist:=[]audioplayer.PtrSong{}
    //get albums
    resp,err=ae.MakeRequest("audio.getAlbums",map[string]string{"owner_id":strconv.Itoa(uid)})
    if err!=nil {
        dialogboxes.ShowErrorDialog("Не удалось загрузить аудиозаписи: "+err.Error())
        return
    }
    //generate map for albums id->name
    albummap:=map[float64]string{}
    gotalbums:=resp["response"].(map[string]interface{})["items"].([]map[string]interface{})
    for _,v:=range gotalbums {
        id:=v["id"].(float64)
        title:=v["title"].(string)
        albummap[id]=title
    }
    for _,v:=range items {
        mp:=v.(map[string]interface{})
        item:=MusicItem{}
        item.Artist=mp["artist"].(string)
        item.Title=mp["title"].(string)
        item.Id=mp["id"].(float64)
        item.Url=mp["url"].(string)
        lyricsid,present:=mp["lyrics_id"].(float64)
        if !present {lyricsid=0}
        item.LyricsId=lyricsid
        duration:=mp["duration"].(float64)
        duration_obj,_:=time.ParseDuration(strconv.FormatFloat(duration,'g',-1,64)+"s")
        item.Duration=duration_obj.String()
        genreId,present:=mp["genre_id"].(float64)
        if present {
            item.Genre=requestwrapper.VkAudioGenres[genreId]
        } else {
            item.Genre=""
        }
        albumid,present:=mp["album_id"].(float64)
        if present {
            item.Album=albummap[albumid]
        } else {
            item.Album=""
        }
        model.Call("appendStruct",item)
        tmplist=append(tmplist,audioplayer.PtrSong{Artist: item.Artist,
                                                   Title: item.Title,
                                                   Url: item.Url,
                                                   Duration: duration,
                                                   LyricsId: item.LyricsId})
    }
    player:=audioplayer.Engine{}
    ae.QMLEngine.Context().SetVar("audioplayer",&player)
    player.Initialize(ae.Token,tmplist)
}
