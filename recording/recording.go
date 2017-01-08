package recording

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var (
	ErrMismatchWrite = errors.New("recording: did not write the same number of bytes that were read")
)

// Recording ...
type Recording struct {
	ctx      context.Context
	url      string
	fname    string
	cancel   context.CancelFunc
	started  time.Time
	restarts int

	Debug bool
	Err   error
}

// New creates a new Recording of the given URL to the given filename for output.
func New(url, fname string) (*Recording, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Hour)

	r := &Recording{
		ctx:     ctx,
		url:     url,
		fname:   fname,
		cancel:  cancel,
		started: time.Now(),
	}

	return r, nil
}

func (r *Recording) Cancel() {
	r.cancel()
}

func (r *Recording) Done() <-chan struct{} {
	return r.ctx.Done()
}

// OutputFilename gets the output filename originally passed into New.
func (r *Recording) OutputFilename() string {
	return r.fname
}

// StartTime gets start time
func (r *Recording) StartTime() time.Time {
	return r.started
}

// Start blockingly starts the recording and returns the error if one is encountered while streaming.
// This should be stopped in another goroutine.
func (r *Recording) Start() error {
	sr, err := exec.LookPath("streamripper")
	if err != nil {
		return err
	}

	dir, err := ioutil.TempDir("", strconv.Itoa(rand.Int()))
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	cmd := exec.Command(sr, r.url, "-d", ".", "-a", r.fname)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	log.Printf("%s: %v", cmd.Path, cmd.Args)

	err = cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		defer r.Cancel()
		err := cmd.Wait()
		if err != nil {
			log.Println(err)
		}
	}()

	defer r.cancel()

	for {
		time.Sleep(250 * time.Millisecond)

		select {
		case <-r.ctx.Done():
			return cmd.Process.Signal(os.Interrupt)
		default:
		}
	}
}
