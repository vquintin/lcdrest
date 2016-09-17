package lcdrest

import "github.com/vquintin/i2c"

const lineSize = 20
const screenSize = 4 * lineSize

type lcdWrapper struct {
	lcd *i2c.Lcd
}

func (lw *lcdWrapper) Write(buf []byte) (int, error) {
	lw.lcd.Clear()
	lw.lcd.Home()
	buffer := buf
	if len(buf) > screenSize {
		buffer = buf[:screenSize]
	}
	buffer = reorderLines(buffer)
	return lw.lcd.Write(buffer)
}

func reorderLines(buffer []byte) []byte {
	size := bufferSize(len(buffer))
	result := make([]byte, size)
	for i := range result {
		result[i] = ' '
	}
	for i, v := range buffer {
		result[position(i)] = v
	}
	return result
}

func bufferSize(n int) int {
	if n < lineSize {
		return n
	} else if n < 2*lineSize {
		return n + lineSize
	} else if n < 3*lineSize {
		return 3 * lineSize
	}
	return n
}

func position(i int) int {
	if i < lineSize {
		return i
	} else if i < 2*lineSize {
		return i + lineSize
	} else if i < 3*lineSize {
		return i - lineSize
	}
	return i
}
