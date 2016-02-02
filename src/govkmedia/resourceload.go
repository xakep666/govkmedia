package main

import (
    "strconv"
    "time"
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
        ae.showErrorDialog("Не удалось загрузить аватар: "+err.Error())
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
        ae.showErrorDialog("Не удалось загрузить имя: "+err.Error())
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
        ae.showErrorDialog("Не удалось загрузить аудиозаписи: "+err.Error())
        return
    }
    model.Call("clear")
    content:=resp["response"].(map[string]interface{})
    items:=content["items"].([]interface{})
    for _,v:=range items {
        mp:=v.(map[string]interface{})
        artist:=mp["artist"].(string)
        title:=mp["title"].(string)
        id:=mp["id"].(float64)
        url:=mp["url"].(string)
        lyricsid,present:=mp["lyrics_id"].(float64)
        if !present {lyricsid=0}
        duration:=mp["duration"].(float64)
        duration_obj,_:=time.ParseDuration(strconv.FormatFloat(duration,'g',-1,64)+"s")
        durationstr:=duration_obj.String()
        item:=MusicItem{Artist: artist, Title: title, Duration: durationstr,Id: id,LyricsId: lyricsid,Url:url}
        model.Call("appendStruct",item)
    }
}
