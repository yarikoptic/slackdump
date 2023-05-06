package convert

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/rusq/dlog"
	"github.com/rusq/fsadapter"
	"github.com/rusq/slackdump/v2/internal/chunk"
	"github.com/rusq/slackdump/v2/logger"
	"github.com/slack-go/slack"
)

const (
	testSrcDir = "../../tmp/slackdump_20230506_120330" // TODO: fix manual nature of this/obfuscate
)

var testLogger = dlog.New(os.Stderr, "unit ", log.Lshortfile|log.LstdFlags, true)

func TestChunkToExport_Validate(t *testing.T) {
	srcDir, err := chunk.OpenDir(testSrcDir)
	if err != nil {
		t.Fatal(err)
	}
	var testTrgDir = t.TempDir()

	type fields struct {
		Src          *chunk.Directory
		Trg          fsadapter.FS
		UploadDir    string
		IncludeFiles bool
		SrcFileLoc   func(*slack.Channel, *slack.File) string
		TrgFileLoc   func(*slack.Channel, *slack.File) string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"empty", fields{}, true},
		{"no source", fields{Trg: fsadapter.NewDirectory(testTrgDir)}, true},
		{"no target", fields{Src: srcDir}, true},
		{
			"valid, no files",
			fields{
				Src:          srcDir,
				Trg:          fsadapter.NewDirectory(testTrgDir),
				IncludeFiles: false,
			},
			false,
		},
		{
			"valid, include files, but no location functions",
			fields{
				Src:          srcDir,
				Trg:          fsadapter.NewDirectory(testTrgDir),
				IncludeFiles: true,
			},
			true,
		},
		{
			"valid, include files, with location functions",
			fields{
				Src:          srcDir,
				Trg:          fsadapter.NewDirectory(testTrgDir),
				IncludeFiles: true,
				SrcFileLoc: func(*slack.Channel, *slack.File) string {
					return ""
				},
				TrgFileLoc: func(*slack.Channel, *slack.File) string {
					return ""
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ChunkToExport{
				src:          tt.fields.Src,
				trg:          tt.fields.Trg,
				includeFiles: tt.fields.IncludeFiles,
				srcFileLoc:   tt.fields.SrcFileLoc,
				trgFileLoc:   tt.fields.TrgFileLoc,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ChunkToExport.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChunkToExport_Convert(t *testing.T) {
	cd, err := chunk.OpenDir(testSrcDir)
	if err != nil {
		t.Fatal(err)
	}
	testTrgDir, err := os.MkdirTemp("", "slackdump")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(testTrgDir)
	// var testTrgDir = t.TempDir()
	fsa, err := fsadapter.NewZipFile(filepath.Join(testTrgDir, "slackdump.zip"))
	if err != nil {
		t.Fatal(err)
	}
	defer fsa.Close()

	c := NewChunkToExport(cd, fsa, WithIncludeFiles())

	ctx := logger.NewContext(context.Background(), testLogger)
	if err := c.Convert(ctx); err != nil {
		t.Fatal(err)
	}
}
