package downloaders

import (
    "github.com/jeffail/tunny"
    "gopkg.in/qml.v1"
    "dialogboxes"
    "errors"
)

type DownloadEngine struct {
    controlChannels []chan WorkerState
    engine *qml.Engine
    mainwindow *qml.Window
    threadPool *tunny.WorkPool
    state WorkerState
    srcs []Downloadable
}

func Initialize(srcs []Downloadable,threads int) (d *DownloadEngine,err error) {
    if srcs==nil || threads<=0 { return nil,errors.New("Invalid argument on downloader init")}
    d.srcs=srcs
    d.state=Running
    //gui init code here
    d.threadPool=tunny.CreatePool(threads,func(object interface{}) interface{} {
        cchan:=make(chan WorkerState)
        d.controlChannels=append(d.controlChannels,cchan)
        err:=object.(Downloadable).Download(cchan)
        if err!=nil { dialogboxes.ShowErrorDialog(err.Error()) }
        return nil
    })
    for _,v:=range srcs {
        d.threadPool.SendWork(v)
    }
    return
}

func (d *DownloadEngine) Destruct() {
    d.threadPool.Close()
}

func (d *DownloadEngine) AppendDownload(dl Downloadable) {
    d.srcs=append(d.srcs,dl)
    d.threadPool.SendWork(dl)
}

func (d *DownloadEngine) Resume() {
    if d.state!=Paused { return }
    for _,v:=range d.controlChannels { v<-Running }
    d.state=Running
}

func (d *DownloadEngine) Pause() {
    if d.state!=Running { return }
    for _,v:=range d.controlChannels { v<-Paused }
    d.state=Paused
}

//gui deactivates buttons
func (d *DownloadEngine) Cancel() {
    for _,v:=range d.controlChannels { v<-Stopped }
    d.state=Stopped
}

//percentage
func (d *DownloadEngine) TotalProgress() float64 {
    var totaldl,completedl float64
    for _,v:=range d.srcs {
        totaldl+=float64(v.BytesTotal())
        completedl+=float64(v.BytesDownloaded())
    }
    return completedl/totaldl
}
