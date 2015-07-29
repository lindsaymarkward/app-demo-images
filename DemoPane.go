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
	"github.com/ninjasphere/sphere-go-led-controller/fonts/O4b03b"
	"github.com/ninjasphere/sphere-go-led-controller/util"
)

var tapInterval = time.Millisecond * 500
var introDuration = time.Millisecond * 1500

// load a particular image - for a 'logo' in this case
var imageLogo = util.LoadImage(util.ResolveImagePath("logo.png"))
var border = util.LoadImage(util.ResolveImagePath("border-green.gif"))

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
			stateImages[name] = util.LoadImage(util.ResolveImagePath(f.Name()))
			// also save names of images used as keys in images map
			stateImageNames = append(stateImageNames, name)
		}
	}
}

// DemoPane stores the data we want to access.
// The struct doesn't need any particular fields but must implement the remote.pane interface functions
type DemoPane struct {
	lastTap         time.Time
	lastDoubleTap   time.Time
	lastTapLocation gestic.Location

	displayingIntro bool
	introTimeout    *time.Timer
	visible bool

	isImageMode bool
	imageIndex  int
	app *App
}

// NewDemoPane creates a DemoPane with the data and timers initialised
// It doesn't need to do much more than create a struct if you want
// the app is passed in so that the pane can access the data and methods in it
func NewDemoPane(a *App) *DemoPane {

	pane := &DemoPane{
		lastTap: time.Now(),
		isImageMode: true,
		imageIndex:  0,
		app: a,
	}

	// AfterFunc(0, ...) creates a timer with no duration so it doesn't fire until Reset is called
	// timers only fire once, unless they are reset again
	pane.introTimeout = time.AfterFunc(0, func() {
		log.Infof("introTimeout func firing...")
		pane.displayingIntro = false
	})

	return pane
}

// Gesture is called by the system when the LED matrix receives any kind of gesture
// it only seems to receive tap gestures ("GestureNone" type)
func (p *DemoPane) Gesture(gesture *gestic.GestureMessage) {
	log.Infof("gesture received - %v, %v", gesture.Touch, gesture.Position)

	// check the second last touch location since the most recent one before a tap is usually blank it seems
	lastLocation := p.lastTapLocation
	p.lastTapLocation = gesture.Touch

	if gesture.Tap.Active() && time.Since(p.lastTap) > tapInterval {
		p.lastTap = time.Now()

		log.Infof("Tap! %v", lastLocation)

		// change between images - right or left
		if lastLocation.East {
			p.imageIndex++
			p.imageIndex %= len(stateImageNames)
		} else {
			p.imageIndex--
			if p.imageIndex < 0 {
				p.imageIndex = len(stateImageNames) - 1
			}
		}
		log.Infof("Showing image: %v", stateImageNames[p.imageIndex])
	}

	if gesture.DoubleTap.Active() && time.Since(p.lastDoubleTap) > tapInterval {
		p.lastDoubleTap = time.Now()

		log.Infof("Double Tap!")

		// change between image and text displaying (in Render)
		p.isImageMode = !p.isImageMode
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

	if !p.visible {
		p.visible = true
		p.displayingIntro = true
		p.introTimeout.Reset(introDuration)
	}

	// simply return the logo image
	if p.displayingIntro {
		return imageLogo.GetNextFrame(), nil
	}

	// create an empty 16*16 RGBA image for the Draw function to draw into (to be returned)
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))

	// display either images or some text
	if p.isImageMode {
		// set one of the images loaded at the start to be displayed
		// (p.imageIndex is just an index to change so we can see different images)
		stateImg := stateImages[stateImageNames[p.imageIndex]]
		// Draw (built-in Go function) draws the frame from stateImg into the img 'image' starting at 4th parameter, "Over" the top
		draw.Draw(img, img.Bounds(), stateImg.GetNextFrame(), image.Point{0, 0}, draw.Over)

	} else {
		// let's make a dynamic colour :)
//		red := uint8(float64(p.imageIndex) / float64(len(stateImageNames)) * 255)
		red := uint8(float64(time.Now().Second()) / float64(60) * 255)
		// draw the index up the top
		drawText(fmt.Sprintf("%2d", p.imageIndex), color.RGBA{red, 250, 250, 255}, 2, img)
		// draw the text from app down the bottom
		drawText(p.app.textData, color.RGBA{253, 151, 32, 255}, 9, img)
		// add a border to the text (you can combine multiple images/text - just keep drawing into img
		if p.imageIndex == 0 {
			draw.Draw(img, img.Bounds(), border.GetNextFrame(), image.Point{0, 0}, draw.Over)
		}
	}

	// return the image we've created to be rendered to the matrix
	return img, nil
}

// drawText is a helper function to draw a string of text into an image
func drawText(text string, col color.RGBA, top int, img *image.RGBA) {
	width := O4b03b.Font.DrawString(img, 0, 8, text, color.Black)
	start := int(16 - width - 1)

	O4b03b.Font.DrawString(img, start, top, text, col)
}