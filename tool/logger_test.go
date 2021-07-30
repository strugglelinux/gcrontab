package tool

import (
	"os"
	"sync"
	"testing"
)

func TestLogger_SetVerbose(t *testing.T) {
	type fields struct {
		verbose  bool
		mux      sync.Mutex
		dir      string
		file     *os.File
		filePath string
		cache    Cache
	}
	type args struct {
		isVerbose bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				verbose:  tt.fields.verbose,
				mux:      tt.fields.mux,
				dir:      tt.fields.dir,
				file:     tt.fields.file,
				filePath: tt.fields.filePath,
				cache:    tt.fields.cache,
			}
			l.SetVerbose(tt.args.isVerbose)
		})
	}
}
