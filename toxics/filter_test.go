package toxics_test

import (
	"github.com/Shopify/toxiproxy/meta"
	"github.com/Shopify/toxiproxy/stream"
	"github.com/Shopify/toxiproxy/toxics"
	"testing"
)

func TestFilterToxicIsNoOpForUnfilteredAddresses(t *testing.T) {
	hello := "hello world"
	data := []byte(hello)
	filter := &toxics.FilterToxic{
		FilteredAddresses: []string{},
	}

	input := make(chan *stream.StreamChunk)
	output := make(chan *stream.StreamChunk)
	connectionMeta := meta.ConnectionMeta{
		DownstreamAddress: "127.0.0.1",
	}
	stub := toxics.NewToxicStub(input, output, &connectionMeta)

	done := make(chan bool)
	go func() {
		filter.Pipe(stub)
		done <- true
	}()
	defer func() {
		close(input)
		for {
			select {
			case <-done:
				return
			case <-output:
			}
		}
	}()

	input <- &stream.StreamChunk{Data: data}

	buf := make([]byte, 0, len(data))

	c, ok := <-output
	if !ok {
		t.Errorf("Output channel closed unexpectedly!")
		return
	}

	buf = append(buf, c.Data...)
	result := string(buf)

	if result != hello {
		t.Errorf("Expected %s but got %s\n", hello, result)
	}
}

func TestFilterToxicIsDroppingConnectionsForFilteredAddresses(t *testing.T) {
	hello := "hello world"
	data := []byte(hello)
	filter := &toxics.FilterToxic{
		FilteredAddresses: []string{"127.0.0.1"},
	}

	input := make(chan *stream.StreamChunk)
	output := make(chan *stream.StreamChunk)
	connectionMeta := meta.ConnectionMeta{
		DownstreamAddress: "127.0.0.1",
	}
	stub := toxics.NewToxicStub(input, output, &connectionMeta)

	done := make(chan bool)
	go func() {
		filter.Pipe(stub)
		done <- true
	}()
	defer func() {
		close(input)
		for {
			select {
			case <-done:
				return
			case <-output:
			}
		}
	}()

	input <- &stream.StreamChunk{Data: data}

	_, ok := <-output
	if ok {
		t.Errorf("Output channel was not closed!")
	}
}