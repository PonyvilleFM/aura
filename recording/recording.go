package recording

import (
	"bufio"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	ErrMismatchWrite = errors.New("recording: did not write the same number of bytes that were read")
)

// Recording ...
type Recording struct {
	ctx     context.Context
	url     string
	fname   string
	fout    *os.File
	cancel  context.CancelFunc
	started time.Time

	Debug bool
	Err   error
}

// New creates a new Recording of the given URL to the given filename for output.
func New(url, fname string) (*Recording, error) {
	fout, err := os.Create(fname)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Hour)

	r := &Recording{
		ctx:     ctx,
		url:     url,
		fname:   fname,
		fout:    fout,
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
	resp, err := http.Get(r.url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	defer r.fout.Close()
	defer r.cancel()

	reader := bufio.NewReader(resp.Body)

	c := time.NewTicker(5 * time.Second)
	defer c.Stop()

	buf := make([]byte, 65536)

	for {
		time.Sleep(250 * time.Millisecond)

		select {
		case <-r.ctx.Done():
			r.fout.Sync()
			return nil
		case <-c.C:
			if r.Debug {
				log.Println("Syncing file")
			}
			err := r.fout.Sync()
			if err != nil {
				r.Err = err
				return err
			}
		default:
			nr, err := reader.Read(buf)
			if err != nil {
				r.Err = err
				return err
			}

			if r.Debug {
				log.Printf("%d bytes read", nr)
			}

			buf = buf[:nr]

			_, err = r.fout.Write(buf)
			if err != nil {
				r.Err = err
				return err
			}
		}
	}
}
