## app-demo-images
Ninja Sphere Application DEMO of adding a pane to the LED Matrix for displaying images and text and responding to tap gestures

**Stuff you can do:**

 - Tap once on the right/east side of the spheramid to advance to the next image
 - Tap once on the left/west side for the previous image
 - Double tap to switch between image and text displaying
 - East/West taps still work in text mode
 - Put whatever 16*16 GIF or PNG images you like in the `images` folder

Run with something like:

`DEBUG=* ./app-demo-images --mqtt.host=ninjasphere.local --mqtt.port=1883 --serial=yourSerial# --led.host=ninjasphere.local`

