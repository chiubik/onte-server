package pkg

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

	url, err := CloudtunnelRun(ctx)
	if err != nil {
		t.Errorf("error: %s", err)
	}

	t.Log(url)

	if strings.Contains(url, "Couldn't find tunnel") {
		t.Errorf("it doesn't return url")
	}
}
