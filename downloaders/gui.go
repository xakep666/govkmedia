package downloaders

import (
    "github.com/jeffail/tunny"
    "gopkg.in/qml.v1"
    _"dialogboxes"
    "errors"
    "log"
    "time"
    "path/filepath"
    "github.com/nareix/curl"
)

//go:generate genqrc qml icons

type DisplayableItem struct {
    Fname string
    Dlspeed string
    Dlprogress float64
}

type DownloadEngine struct {
    engine *qml.Engine
    mainwindow *qml.Window
    threadPool *tunny.WorkPool
    srcs []Downloadable
    totaldl []int64
    completedl []int64
}

func Initialize(srcs []Downloadable,threads int) (d *DownloadEngine,err error) {
    if srcs==nil || threads<=0 { return nil,errors.New("Неверный аргумент при инициализации загрузчика")}
    d=new(DownloadEngine)
    d.registerQMLTypes()
    d.srcs=srcs
    d.engine=qml.NewEngine()
    component,err:=d.engine.LoadFile("qrc:///qml/downloadgui.qml")
    if err!=nil {
        log.Panicln(err)
    }
    d.engine.Context().SetVar("engine",d)
    d.mainwindow=component.CreateWindow(nil)
    d.mainwindow.Show()
    model:=d.mainwindow.Root().ObjectByName("filetable").ObjectByName("dllist")
    d.threadPool,err=tunny.CreatePool(threads,func(object interface{}) interface{} {
        dlo:=object.(Downloadable)
        d.completedl=append(d.completedl,0)
        d.totaldl=append(d.totaldl,0)
        model.Call("appendStruct",&DisplayableItem{Fname: filepath.Base(dlo.ActualPath()),Dlspeed:"0B/s",Dlprogress:0.0})
        item:=model.Call("back").(qml.Object)
        index:=model.Int("count")-1
        dlo.Progress(func (p curl.ProgressStatus) {
            if p.Size!=0 { d.completedl[index]=p.Size }
            if p.ContentLength!=0 { d.totaldl[index]=p.ContentLength }
            d.updateTotalProgress()
            //when download completes, percents sets to 0
            if p.Percent!=0 { item.Set("dlprogress",p.Percent) }
            item.Set("dlspeed",curl.PrettySpeedString(p.Speed))
        },500*time.Millisecond)
        _,err:=object.(Downloadable).Do()
        if err!=nil { /*dialogboxes.ShowErrorDialog(err.Error())*/ log.Println(err.Error()) }
        return nil
    }).Open()
    if err!=nil {
        log.Panicln(err)
    }
    for _,v:=range srcs {
        d.threadPool.SendWorkAsync(v,func(interface{},error) {})
    }
    return
}

func (d *DownloadEngine) Destruct() {
    d.Cancel()
    //d.threadPool.Close()
}

func (d *DownloadEngine) AppendDownload(dl Downloadable) {
    d.srcs=append(d.srcs,dl)
    d.threadPool.SendWork(dl)
}

func (d *DownloadEngine) Resume() {
    for i:=0;i<len(d.srcs);i++ { d.srcs[i].Start() }
}

func (d *DownloadEngine) Pause() {
    for i:=0;i<len(d.srcs);i++ { d.srcs[i].Pause() }
}

func (d *DownloadEngine) Cancel() {
    for i:=0;i<len(d.srcs);i++ { d.srcs[i].Stop() }
}

//percentage
func (d *DownloadEngine) TotalProgress() float64 {
    var totaldl,completedl float64
    for i:=0;i<len(d.totaldl);i++ {
        totaldl+=float64(d.totaldl[i])
        completedl+=float64(d.completedl[i])
    }
    return completedl/totaldl
}

func (d *DownloadEngine) updateTotalProgress() {
    mwroot:=d.mainwindow.Root()
    pr:=d.TotalProgress()
    mwroot.ObjectByName("totalprogressbar").Set("value",pr)
    mwroot.ObjectByName("percents").Call("setPercents",pr*100)
}

func (d *DownloadEngine) registerQMLTypes() {
    types:=[]qml.TypeSpec{
        {
            Init: func(m *DisplayableItem, obj qml.Object) {},
        },
    }
    qml.RegisterTypes("GoExtensions",1,0,types)
}
