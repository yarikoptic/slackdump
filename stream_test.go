package slackdump

import (
	"context"
	"os"
	"path/filepath"
	"runtime/trace"
	"testing"
	"time"

	"github.com/rusq/chttp"
	"github.com/rusq/slackdump/v2/internal/cache"
	"github.com/rusq/slackdump/v2/internal/chunk"
	"github.com/rusq/slackdump/v2/internal/chunk/proctest"
	"github.com/slack-go/slack"
)

var expandedLimits = Limits{
	Workers:         10,
	DownloadRetries: 10,
	Tier2: TierLimits{
		Boost:   100,
		Burst:   100,
		Retries: 20,
	},
	Tier3: TierLimits{
		Boost:   100,
		Burst:   100,
		Retries: 20,
	},
	Request: RequestLimit{
		Conversations: 200,
		Channels:      200,
		Replies:       1000,
	},
}

const testConversation = "C01SPFM1KNY"

func TestChannelStream(t *testing.T) {
	t.Skip("skipping test")
	ucd, err := os.UserCacheDir()
	if err != nil {
		t.Fatal(err)
	}
	m, err := cache.NewManager(filepath.Join(ucd, "slackdump"))
	if err != nil {
		t.Fatal(err)
	}
	wsp, err := m.Current()
	if err != nil {
		t.Fatal(err)
	}
	prov, err := m.Auth(context.Background(), wsp, nil)
	if err != nil {
		t.Fatal(err)
	}

	sd := slack.New(prov.SlackToken(), slack.OptionHTTPClient(chttp.Must(chttp.New("https://slack.com", prov.Cookies()))))

	f, err := os.Create("record.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rec := chunk.NewRecorder(f)
	defer rec.Close()

	cs := newChannelStream(sd, &DefOptions.Limits, time.Time{}, time.Time{})
	if err := cs.Stream(context.Background(), testConversation, rec); err != nil {
		t.Fatal(err)
	}
}

func TestRecorderStream(t *testing.T) {
	t.Skip()
	ctx, task := trace.NewTask(context.Background(), "TestRecorderStream")
	defer task.End()
	f, err := os.Open("record.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rgnNewSrv := trace.StartRegion(ctx, "NewServer")
	srv := proctest.NewServer(f)
	rgnNewSrv.End()
	defer srv.Close()
	sd := slack.New("test", slack.OptionAPIURL(srv.URL+"/api/"))

	w, err := os.Create("replay_record.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rec := chunk.NewRecorder(w)
	defer rec.Close()

	rgnStream := trace.StartRegion(ctx, "Stream")
	cs := newChannelStream(sd, &expandedLimits, time.Time{}, time.Time{})
	if err := cs.Stream(ctx, testConversation, rec); err != nil {
		t.Fatal(err)
	}
	rgnStream.End()
}

func TestReplay(t *testing.T) {
	t.Skip("skipping test")
	f, err := os.Open("record.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	srv := proctest.NewServer(f)
	defer srv.Close()
	sd := slack.New("test", slack.OptionAPIURL(srv.URL+"/api/"))

	reachedEnd := false
	for i := 0; i < 100_000; i++ {
		resp, err := sd.GetConversationHistory(&slack.GetConversationHistoryParameters{ChannelID: testConversation})
		if err != nil {
			t.Fatalf("error on iteration %d: %s", i, err)
		}
		if !resp.HasMore {
			reachedEnd = true
			t.Log("no more messages")
			break
		}
	}
	if !reachedEnd {
		t.Fatal("didn't reach end of stream")
	}
}
