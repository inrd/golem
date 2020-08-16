package main

var history [][]byte

const historySize = 8

func PushToHistory(data []byte) {
	if len(history) >= historySize {
		history = history[1:]
	}
	history = append(history, data)
}

func Undo() []byte {
	if len(history) < 1 {
		return nil
	}
	pop := history[len(history)-1]
	history = history[:len(history)-1]
	return pop
}
