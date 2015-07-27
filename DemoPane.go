package main

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"strings"
	"time"

	"fmt"
	"github.com/ninjasphere/gestic-tools/go-gestic-sdk"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/config"
	"github.com/ninjasphere/sphere-go-led-controller/fonts/O4b03b"
	"github.com/ninjasphere/sphere-go-led-controller/util"
)

var tapInterval = config.MustDuration("uber.tapInterval")
var updateOnTap = config.MustBool("uber.updateOnTap")
var introDuration = config.MustDuration("uber.introDuration")
var visibleTimeout = config.MustDuration("uber.visibilityTimeout") // Time between frames rendered before we reset the ui.
var updateInterval = config.MustDuration("uber.updateInterval")

var imageLogo = util.LoadImage(util.ResolveImagePath("logo.png"))
var border = util.LoadImage(util.ResolveImagePath("imageSurge.png"))

var stateImages map[string]util.Image
var stateImageNames []string

// loadImages saves the PNG and GIF files in the images directory into the stateImages map
func loadImages() {
	files, err := ioutil.ReadDir("./images")

	if err != nil {
		panic("Couldn't load images: " + err.Error())
	}

	stateImages = make(map[string]util.Image)

	for _, f := range files {

		if strings.HasSuffix(f.Name(), ".gif") || strings.HasSuffix(f.Name(), ".png") {
			name := strings.TrimSuffix(strings.TrimSuffix(f.Name(), ".png"), ".gif")

			log.Infof("Found state image: " + name)
			// TODO - check; I don't think we need the "/"
			stateImages[name] = util.LoadImage(util.ResolveImagePath("/" + f.Name()))
			// also save names of images used as keys in images map
			stateImageNames = append(stateImageNames, name)
		}
	}
}

// DemoPane stores the data we want to access.
// The struct doesn't need any particular fields but must implement the remote.pane interface functions
type DemoPane struct {
	lastTap       time.Time
	lastDoubleTap time.Time

	displayingIntro bool
	introTimeout    *time.Timer

	visible        bool
	visibleTimeout *time.Timer

	updateTimer *time.Timer

	keepAwake        bool
	keepAwakeTimeout *time.Timer

	// TODO: can data be stored in app, not pane? - Just pass app to NewDemoPane, I think
	myText string
	test   bool
	number int
	f      float64
}

// NewDemoPane creates a DemoPane with the data and timers initialised
// It doesn't need to do much more than create a struct if you want
func NewDemoPane(conn *ninja.Connection) *DemoPane {

	pane := &DemoPane{
		lastTap: time.Now(),
		number:  0,
		f:       0.0,
	}

	pane.test = false

	// TODO: figure these timers out - how do they work?

	pane.visibleTimeout = time.AfterFunc(0, func() {
		// TODO: keepAwake is unused, I think
		pane.keepAwake = false
		pane.visible = false
	})

	pane.introTimeout = time.AfterFunc(0, func() {
		pane.displayingIntro = false
	})

	pane.updateTimer = time.AfterFunc(0, func() {
		log.Infof("updateTimer...")
		if !pane.visible {
			return
		}

		err := pane.UpdateData(false)
		if err != nil {
			log.Errorf("Failed to get uber data: %s", err)
			pane.updateTimer.Reset(time.Second * 5)
		}
	})

	pane.keepAwakeTimeout = time.AfterFunc(0, func() {
		pane.keepAwake = false
	})

	return pane
}

func (p *DemoPane) UpdateData(once bool) error {
	if !once && p.visible {
		p.updateTimer.Reset(updateInterval)
	}
	if p.test {
		//		p.myText = time.ANSIC
		p.myText = "Test"
	} else {
		p.myText = ":)"
	}

	//	p.number++
	p.f += 2.0

	return nil
}

// Gesture is called by the system when the LED matrix receives any kind of gesture
func (p *DemoPane) Gesture(gesture *gestic.GestureMessage) {
	log.Infof("%v gesture received", gesture.Gesture.Gesture.String())

	if gesture.Tap.Active() && time.Since(p.lastTap) > tapInterval {
		p.lastTap = time.Now()

		log.Infof("Tap!")

		p.number++
		p.number %= len(stateImageNames)

		p.test = true

		if updateOnTap {
			go p.UpdateData(true)
		}

		//		img := image.NewRGBA(image.Rect(0, 0, 16, 16))
		//
		//		drawText := func(text string, col color.RGBA, top int) {
		//			width := O4b03b.Font.DrawString(img, 0, 8, text, color.Black)
		//			start := int(16 - width - 1)
		//
		//			O4b03b.Font.DrawString(img, start, top, text, col)
		//		}
		//
		//		drawText("N/A", color.RGBA{253, 151, 32, 255}, 2)
		//			} else {
		//				drawText(fmt.Sprintf("%dm", 3), color.RGBA{253, 151, 32, 255}, 2)
		//			}

		//			drawText(fmt.Sprintf("%.1fx", 2.1), color.RGBA{69, 175, 249, 255}, 9)

	}

	if gesture.DoubleTap.Active() && time.Since(p.lastDoubleTap) > tapInterval {
		p.lastDoubleTap = time.Now()

		log.Infof("Double Tap!")

		p.test = false
		//		p.number = 0
		go p.UpdateData(true)
	}
}

// KeepAwake is needed as it's part of the remote.pane interface
func (p *DemoPane) KeepAwake() bool {
	return true
}

// IsEnabled is needed as it's part of the remote.pane interface
func (p *DemoPane) IsEnabled() bool {
	return true
}

// Render is called by the system repeatedly when the pane is visible
// It should return the RGBA image to be rendered on the LED matrix
func (p *DemoPane) Render() (*image.RGBA, error) {
	//	log.Infof("Rendering DemoPane (visible = %v)", p.visible)
	p.visibleTimeout.Reset(visibleTimeout)

	if !p.visible {
		p.visible = true
		p.displayingIntro = true

		p.introTimeout.Reset(introDuration)

		go p.UpdateData(false)
	}

	if p.displayingIntro {
		return imageLogo.GetNextFrame(), nil
	}

	// img here is an empty 16*16 RGBA image for the Draw function to draw into
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))

	// set one of the images loaded at the start to be displayed
	// (p.number is just an index to change so we can see different images)
	stateImg, ok := stateImages[stateImageNames[p.number]]
	if !ok {
		panic("Unknown state/image")
	}
	// Draw (built-in Go function) draws the frame from stateImg into the img 'image' starting at 4th parameter, "Over" the top
	draw.Draw(img, img.Bounds(), stateImg.GetNextFrame(), image.Point{0, 0}, draw.Over)

	// TODO - can we make draw text non-local function definition?
	// draw Text
//	drawText := func(text string, col color.RGBA, top int, offsetY int) {
//		width := O4b03b.Font.DrawString(img, 0, 8, text, color.Black)
//		start := int(16 - width + offsetY)
//
//		O4b03b.Font.DrawString(img, start, top, text, col)
//	}

	// TODO - 2 states with taps to switch - image/text...
	drawText := func(text string, col color.RGBA, top int) {
		width := O4b03b.Font.DrawString(img, 0, 8, text, color.Black)
		start := int(16 - width - 1)

		O4b03b.Font.DrawString(img, start, top, text, col)
	}

	// TODO - something else here :)
	if !p.test {
		drawText("N/A", color.RGBA{253, 151, 32, 255}, 2)
	} else {

		drawText(fmt.Sprintf("%dm", p.number), color.RGBA{253, 151, 32, 255}, 2)
		drawText(fmt.Sprintf("%.1f", p.f), color.RGBA{253, 151, 32, 255}, 9)
		p.f += 0.5
	}

	drawText(fmt.Sprintf("%s", p.myText), color.RGBA{69, 175, 249, 255}, 9)

	draw.Draw(img, img.Bounds(), border.GetNextFrame(), image.Point{0, 0}, draw.Over)

	// return the image we've created by drawing to it
	return img, nil
}
