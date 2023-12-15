package tui

import (
	"time"

	"github.com/dundee/gdu/v5/internal/common"
	"github.com/dundee/gdu/v5/pkg/path"
)

func (ui *UI) updateProgress() {
	color := "[white:red:b]"
	if ui.UseColors {
		color = "[red::b]"
	}

	progressChan := ui.Analyzer.GetProgressChan()
	doneChan := ui.Analyzer.GetDone()

	var progress common.CurrentProgress
	start := time.Now()

	for {
		select {
		case progress = <-progressChan:
		case <-doneChan:
			return
		}

		func(itemCount int, totalSize int64, currentItem string) {
			delta := time.Since(start).Round(time.Second)

			ui.app.QueueUpdateDraw(func() {
				ui.progress.SetText("[white::-] Total items: " +
					color +
					common.FormatNumber(int64(itemCount)) +
					"[white::-], size: " +
					color +
					ui.formatSize(totalSize, false, false) +
					"[white::-], elapsed time: " +
					color +
					delta.String() +
					"[white::-]\nCurrent item: [white::b]" +
					path.ShortenPath(currentItem, ui.currentItemNameMaxLen))
			})
		}(progress.ItemCount, progress.TotalSize, progress.CurrentItemName)

		time.Sleep(100 * time.Millisecond)
	}
}
