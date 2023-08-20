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
    downloadCtxsChan := make(chan *DownloadCtx)
    var jobs []DownloadCtx

    startedGoRoutines := 0
    pw := progress.NewWriter()
    pw.SetStyle(progress.StyleDefault)
    pw.SetTrackerLength(10)
    pw.SetNumTrackersExpected(len(*links))
    pw.SetUpdateFrequency(time.Millisecond * 100)
    pw.Style().Colors = progress.StyleColorsExample
    pw.SetTrackerPosition(progress.PositionLeft)
    pw.Style().Visibility.ETA = false
    pw.Style().Visibility.ETAOverall = false
    pw.Style().Visibility.Speed = false
    pw.Style().Visibility.Percentage = false
    pw.Style().Visibility.Value = false
    pw.Style().Visibility.TrackerOverall = true
    pw.Style().Options.TimeInProgressPrecision = time.Second
    pw.Style().Options.TimeDonePrecision = time.Second

    go pw.Render()

    for _, val := range *links {
        ctx := createDownloadCtx(val)

        jobs = append(jobs, *ctx)

        pw.AppendTracker(ctx.tracker)

        go startJob(ctx, downloadCtxsChan)
        startedGoRoutines++
    }

    for i := 0; i < startedGoRoutines; i++ {
        ctx := <- downloadCtxsChan
        if ctx.err != nil {
            ctx.tracker.MarkAsErrored()
        } else {
            ctx.tracker.MarkAsDone()
        } 
    }

    time.Sleep(time.Millisecond * 100)
    pw.Stop()
    for pw.IsRenderInProgress() {}

}

func createDownloadCtx(link string) *DownloadCtx {
    var ctx DownloadCtx
    ctx.tracker = &progress.Tracker{}
    ctx.tracker.SetValue(0)
    ctx.tracker.Message = link
    ctx.tracker.Units = progress.UnitsDefault
    ctx.link = link
    ctx.status = JOB_STATUS_NEW
    return &ctx
}

func startJob(ctx *DownloadCtx, downloadCtxs chan *DownloadCtx) {
    
    getTitle(ctx)

    if ctx.err != nil {
        downloadCtxs <- ctx
        return
    }
    ctx.tracker.UpdateMessage(ctx.title)

    ytdlpDownload(ctx, downloadCtxs)

    downloadCtxs <- ctx
}

func getTitle(ctx *DownloadCtx) {

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
    ctx *DownloadCtx, downloadCtxs chan *DownloadCtx) {

    cmd := exec.Command(YT_DLP, ctx.link)
    _, err := cmd.Output()
    if err != nil {
        ctx.err = &err
        ctx.status = JOB_STATUS_ERR
        downloadCtxs <- ctx
        return
    }
    ctx.status = JOB_STATUS_COMPLETED

}
