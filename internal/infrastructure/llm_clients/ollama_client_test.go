package llm_clients_test

import (
	"context"
	"fmt"
	"llm-agent-go/internal/infrastructure/llm_clients"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func newTestClient(serverURL string) llm_clients.OllamaClient {
	logger := zerolog.Nop()
	return llm_clients.NewOllamaClient(serverURL, "test-model", logger)
}

func TestGenerate_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"response":"hello"}`)
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)

	res, err := client.Generate(context.Background(), "hi")
	if err != nil {
		t.Fatal(err)
	}

	if res != "hello" {
		t.Fatalf("expected hello got %s", res)
	}
}

func TestGenerate_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)

	_, err := client.Generate(context.Background(), "fail")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestStream_Chunks(t *testing.T) {
	body := `
{"response":"hel","done":false}
{"response":"lo","done":false}
{"done":true}
`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, strings.TrimSpace(body))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)

	ch, err := client.Stream(context.Background(), "hi")
	if err != nil {
		t.Fatal(err)
	}

	var out string
	for s := range ch {
		out += s
	}

	if out != "hello" {
		t.Fatalf("expected hello got %s", out)
	}
}

func TestStream_ContextCancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 2)
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := client.Stream(ctx, "hi")
	if err == nil {
		t.Fatal("expected context error")
	}
}

func TestHealth_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)

	if err := client.Health(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestHealth_500(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)

	err := client.Health(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}
