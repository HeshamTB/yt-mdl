package main

import (
	"os/exec"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
)


const YT_DLP string = "yt-dlp"
const YT_DLP_FLAG_GET_TITLE string = "--get-title"

const JOB_STATUS_NEW = 1
const JOB_STATUS_COMPLETED = 2
const JOB_STATUS_ERR = 3

type DownloadCtx struct {
    link string
    title string
    progress uint8
    status uint8
    tracker *progress.Tracker
    err *error
}


func Download(links *[]string) {
    // Start and manage goroutines here
    // Now get titles only
    downloadCtxsChan := make(chan *DownloadCtx)
    getTitleErrChan := make(chan *DownloadCtx)
    startedGoRoutines := 0
    pw := progress.NewWriter()
    pw.SetTrackerLength(10)
    pw.SetNumTrackersExpected(len(*links))
    pw.SetStyle(progress.StyleDefault)
    pw.SetUpdateFrequency(time.Millisecond * 100)
    pw.Style().Visibility.ETA = false
	pw.Style().Colors = progress.StyleColorsExample
    pw.Style().Visibility.Speed = false
    pw.Style().Visibility.Percentage = false
    pw.Style().Visibility.Value = false
    pw.Style().Visibility.TrackerOverall = true

    go pw.Render()

    for _, val := range *links {
        ctx := createDownloadCtx(val)

        pw.AppendTracker(ctx.tracker)

        go startJob(&ctx, downloadCtxsChan, getTitleErrChan)
        startedGoRoutines++
    }

    for startedGoRoutines > 0 {
        select {
        case ctx := <- downloadCtxsChan:
            ctx.tracker.UpdateMessage(ctx.title)
            startedGoRoutines--
            
        case ctx := <- getTitleErrChan:
            ctx.tracker.MarkAsErrored()
            startedGoRoutines--
        }
    }

 
    pw.Stop()

    for pw.IsRenderInProgress() {}
}

func createDownloadCtx(link string) DownloadCtx {
    var ctx DownloadCtx
    ctx.tracker = &progress.Tracker{}
    ctx.tracker.SetValue(0)
    ctx.tracker.Message = link
    ctx.tracker.Units = progress.UnitsDefault
    ctx.link = link
    ctx.status = JOB_STATUS_NEW
    return ctx
}

func startJob(
    ctx *DownloadCtx, downloadCtxs chan *DownloadCtx, errs chan *DownloadCtx) {
    
    getTitle(ctx)

    if ctx.err != nil {
        errs <- ctx
        return
    }
    ctx.tracker.UpdateMessage(ctx.title)

    ytdlpDownload(ctx, downloadCtxs, errs)

    if ctx.err != nil {
        errs <- ctx
        return
    }

    downloadCtxs <- ctx
}
func getTitle(
    ctx *DownloadCtx) {

    cmd := exec.Command(YT_DLP, YT_DLP_FLAG_GET_TITLE, ctx.link)
    stdout, err := cmd.Output()
    if err != nil {
        ctx.err = &err
        ctx.status = JOB_STATUS_ERR
        return
    }
    title := string(stdout)
    title = strings.TrimSpace(title)
    ctx.title = title

}

func ytdlpDownload(
    ctx *DownloadCtx, downloadCtxs chan *DownloadCtx, errs chan *DownloadCtx) {

    cmd := exec.Command(YT_DLP, ctx.link)
    stdout, err := cmd.Output()
    if err != nil {
        ctx.err = &err
        ctx.status = JOB_STATUS_ERR
        errs <- ctx
        return
    }
    title := string(stdout)
    title = strings.TrimSpace(title)
    ctx.title = title

    downloadCtxs <- ctx
    
}
