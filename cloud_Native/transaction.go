package cloud_Native

import (
	"bufio"
	"fmt"
	"os"
)

type Event struct {
	Sequence  uint64
	Key       string
	Value     string
	EventType EventType
}

type TransactionLogger struct {
	events       chan<- Event // Write-only channel for sending events
	errors       <-chan error // Read-only channel for receiving errors
	lastSequence uint64       // The last used event sequence number
	file         *os.File     // The location of the transaction log
}

func NewFileTransactionLogger(fileName string) (*TransactionLogger, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}
	return &TransactionLogger{file: file}, nil
}

func (l *TransactionLogger) Run() {
	events := make(chan Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		for e := range events {
			l.lastSequence++
			_, err := fmt.Fprintf(
				l.file,
				"%d\t%d\t%s\t%s\n",
				l.lastSequence, e.EventType, e.Key, e.Value)
			if err != nil {
				if err != nil {
					errors <- err
					return
				}
			}
		}
	}()
}

func (l *TransactionLogger) ReadEvents() (<-chan Event, <-chan error) {

	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	go func() {
		var e Event
		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			_, err := fmt.Sscanf(
				line, "%d\t%d\t%s\t%s",
				&e.Sequence, &e.EventType, &e.Key, &e.Value)

			if err != nil {
				outError <- fmt.Errorf("input parse error: %w", err)
				return
			}

			if l.lastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}

			outEvent <- e
		}
		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}

	}()
	return outEvent, outError
}

func (l *TransactionLogger) WritePut(key, value string) {
	l.events <- Event{
		Key:       key,
		Value:     value,
		EventType: EventPut,
	}
}

func (l *TransactionLogger) WriteDelete(key string) {
	l.events <- Event{
		Key:       key,
		EventType: EventDelete,
	}
}

func (l *TransactionLogger) Err() <-chan error {
	return l.errors
}
