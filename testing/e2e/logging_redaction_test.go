package e2e

import (
	"encoding/json"
	"testing"

	"github.com/sensu/sensu-go/types"
	"github.com/sensu/sensu-go/types/dynamic"
	"github.com/stretchr/testify/assert"
)

func TestLoggingRedaction(t *testing.T) {
	t.Parallel()

	// Initializes sensuctl
	sensuctl, cleanup := newSensuCtl(t)
	defer cleanup()

	// Start the agent
	agentConfig := agentConfig{
		ID:               "TestLoggingRedaction",
		CustomAttributes: `{"ec2_access_key": "P@ssw0rd!","secret": "P@ssw0rd!"}`,
		Redact:           []string{"ec2_access_key"},
	}
	agent, cleanup := newAgent(agentConfig, sensuctl, t)
	defer cleanup()

	// Allow time agent connection to be established, etcd to start,
	// keepalive to be sent, etc.
	var output []byte
	var err error
	if err := backoff.Retry(func(retry int) (bool, error) {
		if output, err = sensuctl.run("entity", "info", agent.ID,
			"--organization", sensuctl.Organization,
			"--environment", sensuctl.Environment); err != nil {
			// The command returned an error, let's retry
			return false, nil
		}

		// At this point the attempt was successful
		return true, nil
	}); err != nil {
		t.Errorf("no keepalive received: %s", string(output))
	}

	entity := &types.Entity{}
	assert.NoError(t, json.Unmarshal(output, entity))

	// Make sure the ec2_access_key attribute is redacted, which indicates it was
	// received as such in keepalives
	i, _ := entity.Get("ec2_access_key")
	assert.Equal(t, dynamic.Redacted, i)

	// Make sure the secret attribute is not redacted, because it was not
	// specified in the redact flag
	i, _ = entity.Get("secret")
	assert.NotEqual(t, dynamic.Redacted, i)
}
