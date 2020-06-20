package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/go-vgo/robotgo"
)

var (
	pause      = time.Millisecond * 100
	debug      = false
	zero       = false
	left       = false
	right      = false
	mid        = false
	scrolled   = false
	keys       = [4]*bool{&zero, &left, &right, &mid}
	keyPressed = uint8(0)
)

func main() {
	flag.BoolVar(&debug, "d", false, "print debug")
	path := flag.String("p", "/dev/hidraw0", "path to device file")
	flag.DurationVar(&pause, "t", time.Millisecond * 100, "Milliseconds to wait for middle button emulation." +
		" If two buttons pressed simultaneously before timeout, middle button will be pressed." +
		" If one button is pressed more than timeout, usual button will be pressed")
	flag.Parse()

	deb(`"BB" - before click processing; "AA" - after click processing; 1 is left button; 2 is right button; 3 is middle button.
raw: {h='horizont rel' v='vertical rel' k='key_pressed'};
prog: {pk='previous key pressed; kp='key programm is holding down' s='was scrolling'};
pressed: { 1_button 2_button 3_button }
`)

	f, err := os.Open(*path)
	if err != nil {
		log.Fatalln("Error opening f!!!", err)
		return
	}
	defer f.Close()

	// declare chunk size
	const maxSz = 4

	// create buffer
	b := make([]byte, maxSz)

	k, pk := uint8(0), uint8(0) // k='key pressed', pk='previous key pressed'

	for {
		// read content to buffer
		readTotal, err := f.Read(b)
		if err != nil {
			log.Fatalln("Reading error", err)
		}

		if readTotal != 4 || b[0] != 24 {
			continue
		}

		k = b[1]

		if k > 9 {
			continue
		}

		h, v := int(int8(b[2])), int(int8(b[3]))
		go m(h, v)

		deb("BB raw: {h=%d v=%d k=%d}; prog: {pk=%d; kp=%d; scr: %t;}; pressed: { %t %t %t }", h, v, k, pk, keyPressed, scrolled, *keys[1], *keys[2], *keys[3])
		// if mouse click do not change
		if k != pk {
			switch k {
			case 0:
				u()
			case 1, 2:
				if keyPressed != 3 {
					keyPressed = k

					go func() {
						time.Sleep(pause)
						if !mid {
							d(k)
						}
					}()
				}
			case 3:
				keyPressed = k
				d(k)
			}
		}

		deb("AA raw: {h=%d v=%d k=%d}; prog: {pk=%d; kp=%d; scr: %t;}; pressed: { %t %t %t }", h, v, k, pk, keyPressed, scrolled, *keys[1], *keys[2], *keys[3])
		pk = k
	}

}

func u() {
	b := keyPressed
	defer func() {
		*keys[b] = false
		keyPressed = 0
	}()

	if b == 0 {
		return
	} else if b == 3 && scrolled { // if we was scrolling, we don't press middle button
		scrolled = false
		return
	}

	if *keys[b] == true && b != 3 { // if it's middle button we press it only on "up", so it's always is in up state here
		deb("up %d", b)
		robotgo.MouseToggle("up", bt(b))
	} else {
		deb("down %d", b)
		robotgo.MouseToggle("down", bt(b))
		deb("up %d", b)
		robotgo.MouseToggle("up", bt(b))
	}

	*keys[b] = false
	keyPressed = 0
}

func d(b uint8) {
	if *keys[b] == true || keyPressed == 0 {
		return
	}

	if b != 3 { // if it's middle button we press it on "up" action
		deb("down %d", b)
		robotgo.MouseToggle("down", bt(b))
		*keys[b] = true
		keyPressed = b
	}
}

func m(h, v int) {
	if h == 0 && v == 0 {
		return
	}

	if keyPressed != 3 {
		robotgo.MoveRelative(e(h), e(v))
	} else {
		robotgo.Scroll(scr(h), scr(v), 0)
		scrolled = true
	}
}

func scr(x int) int {
	if x == 0 {
		return 0
	}

	f := float64(x)

	return int( -f / math.Abs(f) * math.Sqrt(math.Abs(f)))
}

func e(x int) int {
	f := float64(x)
	return int(f * math.Log(math.Abs(f)*math.Exp(math.Abs(f))))
}

func deb(format string, a ...interface{}) {
	if !debug {
		return
	}

	fmt.Printf(format + "\n", a...)
}

func bt(x uint8) string {
	switch x {
	case 1:
		return "left"
	case 2:
		return "right"
	case 3:
		return "center"
	}

	return "left"
}
