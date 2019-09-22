package mailgun

import (
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var apiKeyArg = flag.String("apikey", "", "Mailgun API key")
var domainArg = flag.String("domain", "", "Domain for mailgun api")

func TestSendEmail(t *testing.T) {
	if *apiKeyArg == "" || *domainArg == "" {
		t.Skip("mailgun api key and/or domain are missing")
	}

	emailer := New(*apiKeyArg, *domainArg, true)

	now := time.Now().Unix()
	from := fmt.Sprintf("test%d@%s", now, *domainArg)
	to := fmt.Sprintf("fake-recipient-%d@zood.xyz", now)
	subj := fmt.Sprintf("Subject %d", now)
	txtMsg := fmt.Sprintf("Text body: %d https://zood.xyz", now)

	err := emailer.SendEmail(from, to, subj, txtMsg, nil)
	require.NoError(t, err)
}
