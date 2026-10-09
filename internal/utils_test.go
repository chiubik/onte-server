package internal

import (
	"context"
	"strings"
	"testing"
)

func TestGenerateIdentifierLenght(t *testing.T) {
	identifier := GenerateIdentifier()

	if len(identifier) != 3 {
		t.Errorf("the length is not equal to 3, it's equal to %d", len(identifier))
	}
}

func TestCloudtunnelRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan string, 1)
	CloudtunnelRun(ctx, ch)

	url := <-ch
	t.Log(url)

	if strings.Contains(url, "Couldn't find tunnel") {
		t.Errorf("it doesn't return url")
	}
}
