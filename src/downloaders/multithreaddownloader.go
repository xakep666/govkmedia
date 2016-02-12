package downloaders

import (
    "os"
    "fmt"
    "path/filepath"
    "time"
    "github.com/nareix/curl"
)

type Downloadable interface {
    SetSource(string)
    AllocateFile(string,string,string)
    ActualPath() string
    Do() (curl.Response,error)
    Pause()
    Start()
    Stop()
    Progress(curl.MonitorProgressCb,time.Duration) *curl.Request
}

type DownloadableFile struct {
    fpath string
    *curl.Request
    ctrl *curl.Control
    isStopped bool
}

func (d *DownloadableFile) SetSource(url string) {
    d.Request=curl.New(url)
}

func (d *DownloadableFile) AllocateFile(dir,name,ext string) {
    isExist:=true
    n:=0
    var allocpath string
    for isExist {
        if n==0 {
            if ext!="" { allocpath=filepath.Join(dir,fmt.Sprintf("%s.%s",name,ext))
            } else { allocpath=filepath.Join(dir,name) }
        } else {
            if ext!="" { allocpath=filepath.Join(dir,fmt.Sprintf("%s (%d).%s",name,n,ext))
            } else { allocpath=filepath.Join(dir,fmt.Sprintf("%s (%d)",name,n)) }
        }
        finfo,err:=os.Stat(allocpath)
        if err!=nil {
            isExist=false
        } else {
            isExist=finfo.Mode().IsRegular()
        }
        n++
    }
    d.fpath=allocpath
    d.SaveToFile(allocpath)
    d.ctrl=d.ControlDownload()
}

func (d *DownloadableFile) ActualPath() string { return d.fpath }
func (d *DownloadableFile) Start() { if !d.isStopped { d.ctrl.Resume() } }
func (d *DownloadableFile) Stop() {
    d.ctrl.Stop()
    os.Remove(d.fpath) //remove incomplete downloads
    d.isStopped=true
}
func (d *DownloadableFile) Pause() { d.ctrl.Pause() }
