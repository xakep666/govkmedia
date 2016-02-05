package downloaders

import (
    "io"
    "os"
    "fmt"
    "path/filepath"
    "errors"
    "time"
    "net/http"
    "strconv"
)

type Downloadable interface {
    AllocateFile(string,string,string) error
    ActualPath() string
    FileDescriptor() *os.File
    Download() error
    BytesDownloaded() uint64
    BytesTotal() uint64
    Progress() float64
    MomentDownloadSpeed() float64 //bytes per second
    EstimatedTime() uint64 //in seconds, gui reformates to XXhXXmXXs
    CancelDownload()
}

type ControllableReadWriter struct {
    io.Reader
    io.Writer
    transferred uint64
    lastCalled *time.Time
    lastTransferred uint64
    momentSpeed float64
    controlChan <-chan workerState
    currentState workerState
}

type DownloadableFile struct {
    Url string
    fpath string
    fdesc *os.File
    isDownloading bool
    size uint64
    source *ControllableReadWriter
}

type workerState int
const (
    running = iota
    paused
    stopped
)

//function which allows pause,resume,stop IO operation, error: on stop - EOF
func (m *ControllableReadWriter) ioLocker() error {
    for {
        newstate,wasmsg:=<-m.controlChan
        //if we recieved nothing, do not change state
        if wasmsg { m.currentState=newstate }
        switch m.currentState {
            //on running, we unlock operation and perform it
            case running: return nil
            //on stopped, return EOF
            case stopped: return io.EOF
            //on paused, hold lock until reciving stopped/running
            case paused: time.Sleep(time.Millisecond) //prevent overloaing
            default: time.Sleep(time.Millisecond)
        }
    }
}

func (m *ControllableReadWriter) updateValues(transferred uint64) {
    //calculate instantaneous speed (deltaSize/deltaTime)
    m.momentSpeed=float64(m.lastTransferred-uint64(transferred))/float64(time.Now().Sub(*m.lastCalled).Seconds())
    now:=time.Now()
    m.lastCalled=&now
    m.lastTransferred=transferred
    m.transferred+=uint64(transferred)
}

func (m *ControllableReadWriter) Read(buf []byte) (int,error) {
    err:=m.ioLocker()
    if err!=nil { return 0,err }
    if m.lastCalled==nil {
        now:=time.Now()
        m.lastCalled=&now
    }
    read,err:=m.Reader.Read(buf)
    m.updateValues(uint64(read))
    return read,err
}

func (m *ControllableReadWriter) Write(buf []byte) (int,error) {
    err:=m.ioLocker()
    if err!=nil { return 0,err }
    if m.lastCalled==nil {
        now:=time.Now()
        m.lastCalled=&now
    }
    read,err:=m.Writer.Write(buf)
    m.updateValues(uint64(read))
    return read,err
}

func (d *DownloadableFile) AllocateFile(dir,name,ext string) error{
    if d.fdesc!=nil { return nil }
    //if exist, try increment number in brackets
    fpath:=filepath.Join(dir,fmt.Sprintf("%s.%s",name,ext))
    fdesc,err:=os.Create(fpath)
    if !os.IsExist(err) {
        return err
    }
    n:=1
    for os.IsExist(err) {
        fpath:=filepath.Join(dir,fmt.Sprintf("%s (%d).%s",name,n,ext))
        fdesc,err=os.Create(fpath)
        n++
    }
    if !os.IsExist(err) {
        return err
    }
    d.fpath=fpath
    d.fdesc=fdesc
    return nil
}

func (d *DownloadableFile) ActualPath() string {
    return d.fpath
}

func (d *DownloadableFile) FileDescriptor() *os.File {
    return d.fdesc
}

//worker gorutine
func (d *DownloadableFile) Download(echan chan error,state <- chan workerState) {
    if d.fdesc==nil {
        echan<-errors.New("Cannot start downloading,file not created")
        return
    }
    //one worker per one file
    if d.isDownloading {
        return
    } else {
        d.isDownloading=true
    }
    resp,err:=http.Get(d.Url)
    if err!=nil {
        echan<-err
        return
    }
    defer resp.Body.Close()
    d.size,_=strconv.ParseUint(resp.Header.Get("Content-Length"),10,64)
    //wrap to enable progress counting
    d.source=&ControllableReadWriter{Reader: resp.Body, controlChan: state}
    _,err=io.Copy(d.fdesc,d.source)
    if err==io.EOF {
        d.removeFile()
    } else {
        if err!=nil {
            echan<-err
            return
        }
    }
}

//needed to remove partitally downloaded file
func (d *DownloadableFile) removeFile() {
    os.Remove(d.fpath)
    d.fpath=""
    d.fdesc=nil
}

func (d *DownloadableFile) BytesDownloaded() uint64 {
    return d.source.transferred
}

func (d *DownloadableFile) BytesTotal() uint64 {
    return d.size
}

func (d *DownloadableFile) Progress() float64 {
    return float64(d.source.transferred)/float64(d.size)
}

func (d *DownloadableFile) MomentDownloadSpeed() float64 {
    return d.source.momentSpeed
}

func (d *DownloadableFile) EstimatedTime() uint64 {
    return (d.size-d.source.transferred)/uint64(d.source.momentSpeed)
}
