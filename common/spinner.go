package common

import (
    "github.com/briandowns/spinner"
    "time"
)

var s = spinner.New(spinner.CharSets[0], 100*time.Millisecond)

func SpinnerStart(prefix string) {
    s.Prefix = prefix
    s.Start()
}

func SpinnerStop() {
    s.Stop()
}
