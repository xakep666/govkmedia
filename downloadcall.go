package main

import (
	"fmt"
    "gopkg.in/qml.v1"
	"govkmedia/dialogboxes"
	"govkmedia/downloaders"
)

func (ae *AppEngine) DownloadAllMusic() {
	mwroot := ae.MainWindow.Root()
	go func() {
		dstpath := dialogboxes.SelectFolderDialog()
		if dstpath == "" {
			return
		}
        table:=mwroot.ObjectByName("musiclist")
        count:=table.Int("count")
        var dl []downloaders.Downloadable
		for i:=0;i<count;i++ {
            qmlItem:=table.Call("goGet",i).(*qml.Common)
            fmt.Println(qmlItem)
            item:=MusicItem{Artist: qmlItem.String("artist"),
                Album: qmlItem.String("album"),
                Genre: qmlItem.String("genre"),
                Title: qmlItem.String("title"),
                LyricsId: qmlItem.Float64("lyricsId"),
                Url: qmlItem.String("url"),
            }
			dlItem := downloaders.DownloadableMusic{Artist: item.Artist,
				Album:       item.Album,
				Genre:       item.Genre,
				Title:       item.Title,
				LyricsId:    item.LyricsId,
				AccessToken: ae.RequestAccesser.Token,
			}
			dlItem.SetSource(item.Url)
			dlItem.AllocateFile(dstpath, fmt.Sprintf("%s - %s", item.Artist, item.Title), "mp3")
			dl = append(dl, &dlItem)
		}
		_, err := downloaders.Initialize(dl, 3)
		if err != nil {
			dialogboxes.ShowErrorDialog(err.Error())
		}
	}()
}
