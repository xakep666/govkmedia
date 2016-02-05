package downloaders

type ListDownloadable interface{
    TotalProgress() float64
    Start()
    Pause()
    Resume()
    Cancel()
}
