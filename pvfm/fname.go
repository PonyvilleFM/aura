package pvfm

import (
	"fmt"
	"time"

	"git.xeserv.us/PVFM/magi-v2/pvl"
)

func GenFilename() (string, error) {
	cal, err := pvl.Get()
	if err != nil {
		return "", nil
	}

	now := cal.Result[0]

	localTime := time.Now()
	thentime := time.Unix(now.StartTime, 0)
	if thentime.Unix() < localTime.Unix() {
		// return fmt.Sprintf("%s - %s.mp3", now.Title, localTime.Format(time.RFC822)), nil
	}

	return fmt.Sprintf("%s - %s.mp3", now.Title, localTime.Format(time.RFC822)), nil

	// return "", errors.New("pvfm: no DJ is live, cannot make filename")
}
