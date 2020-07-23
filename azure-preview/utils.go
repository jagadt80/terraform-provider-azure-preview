package azurepreview

import (
	"fmt"
	"strings"
)

func parseSubscriptionID(input string) (string, error) {
	parts := strings.Split(input, "/")
	if len(parts) != 3 {
		return "", fmt.Errorf("error parsing Subscription ID: unexpected format: %q", input)
	}

	return parts[2], nil
}
