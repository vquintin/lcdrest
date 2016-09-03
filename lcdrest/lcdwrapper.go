package lcdrest

import 	"github.com/vquintin/i2c"

const lineSize = 20
const screenSize = 4 * lineSize

type lcdWrapper struct {
    lcd *i2c.Lcd
}

func (lw *lcdWrapper) Write(buf []byte) (int, error) {
    lcd.Clear()
    lcd.Home()
    buffer := buf[:screenSize]
    buffer = reorderLines(buffer)
    return lcd.Write(buffer)
}

func reorderLines(buffer []byte) []byte {
    size := bufferSize(len(buffer))
    result := make([]byte, size)
    for i,_ := range result {
        result[i] = ' '
    }
    for i, v := range buffer {
        result[position(i)] = buffer[i]
    }
}

func bufferSize(n int) int {
    if n < lineSize {
        return n
    } else if n < 2 * lineSize {
        n + lineSize
    } else if n < 3 * lineSize {
        return 3 * lineSize
    }
    return n
}

func position(i int) int {
    if i < lineSize {
        return n
    } else if i < 2 * lineSize {
        return n + lineSize
    } else if n < 3 * lineSize {
        return n - lineSize
    }
    return n
}