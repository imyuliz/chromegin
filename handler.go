package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"log"
	"math"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, err error) bool {
	if err != nil {
		log.Println(err)
		//logrus.WithError(err).Error("gin context http handler error")
		c.JSON(200, ResJob{
			Code: 400,
			Msg:  err.Error(),
			Uri:  "",
		})
		return true
	}
	return false
}
func sha256String(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func takeShot(arg *ReqJob) (res *ResJob, err error) {
	//set time-out
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout

	if arg.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(arg.Timeout)*time.Second)
		defer cancel()
	}

	// capture screenshot of an element
	var buf []byte
	if err := chromedp.Run(ctx, makeActions(arg, &buf)); err != nil {
		return nil, err
	}
	uri := sha256String(buf) + ".png"
	fp := filepath.Join(staticDir, uri)
	if err := ioutil.WriteFile(fp, buf, 0o644); err != nil {
		return nil, err
	}

	dataString := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf)

	return &ResJob{
		Code: 200,
		Msg:  "OK",
		Uri:  uri,
		Url:  "",
		B64:  dataString,
	}, err

}

// makeActions takes a screenshot of a specific element.
func makeActions(arg *ReqJob, res *[]byte) chromedp.Tasks {
	ts := chromedp.Tasks{
		chromedp.Navigate(arg.Url),
	}
	// if arg.PxWidth > 0 && arg.PxHeight > 0 {
	// 	ts = append(ts, chromedp.EmulateViewport(arg.PxWidth, arg.PxHeight))
	// }
	if arg.PxWidth > 0 && arg.PxHeight > 0 {
		ts = append(ts, chromedp.EmulateViewport(arg.PxWidth, arg.PxHeight, func(sdmop *emulation.SetDeviceMetricsOverrideParams, steep *emulation.SetTouchEmulationEnabledParams) {
			sdmop.DeviceScaleFactor = 3
		}))
	}
	if arg.Wait > 0 {
		wtFn := func(ctx context.Context) error {
			time.Sleep(time.Duration(arg.Wait) * time.Second)
			return nil
		}
		ts = append(ts, chromedp.ActionFunc(wtFn))
	}
	if arg.Sel != "" {
		ts = append(ts, chromedp.WaitVisible(arg.Sel), chromedp.Screenshot(arg.Sel, res, chromedp.NodeVisible, chromedp.ByID))
	} else {
		if arg.Quality < 1 {
			arg.Quality = 80
		}

		fullScreenFn := func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(arg.Quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}

		ts = append(ts, chromedp.ActionFunc(fullScreenFn))

	}

	return ts
}
