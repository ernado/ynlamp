package ynlamp

import (
	"math"
)

func temperatureToRGB(temp float64) (int, int, int) {
	temp = temp / 100.0
	var red, green, blue float64

	if temp <= 66 {
		red = 255
	} else {
		red = temp - 60
		red = 329.698727446 * math.Pow(red, -0.1332047592)
		if red < 0 {
			red = 0
		}
		if red > 255 {
			red = 255
		}
	}

	if temp <= 66 {
		green = temp
		green = 99.4708025861*math.Log(green) - 161.1195681661
		if green < 0 {
			green = 0
		}
		if green > 255 {
			green = 255
		}
	} else {
		green = temp - 60
		green = 288.1221695283 * math.Pow(green, -0.0755148492)
		if green < 0 {
			green = 0
		}
		if green > 255 {
			green = 255
		}
	}

	if temp >= 66 {
		blue = 255
	} else {
		if temp <= 19 {
			blue = 0
		} else {
			blue = temp - 10
			blue = 138.5177312231*math.Log(blue) - 305.0447927307
			if blue < 0 {
				blue = 0
			}
			if blue > 255 {
				blue = 255
			}
		}
	}

	return int(math.Round(red)), int(math.Round(green)), int(math.Round(blue))
}
