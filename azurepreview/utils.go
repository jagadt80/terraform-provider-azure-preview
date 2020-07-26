package azurepreview

import (
	"fmt"
	"strings"
)

func expandStringSlice(input []interface{}) *[]string {
	result := make([]string, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, item.(string))
		} else {
			result = append(result, "")
		}
	}
	return &result
}

func flattenStringSlice(input *[]string) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, item)
		}
	}
	return result
}

func parseSubscriptionID(input string) (string, error) {
	parts := strings.Split(input, "/")
	if len(parts) != 3 {
		return "", fmt.Errorf("error parsing Subscription ID: unexpected format: %q", input)
	}

	return parts[2], nil
}

type budgetResource struct {
	Scope      string
	BudgetName string
}

func parseBudgetID(input string) (*budgetResource, error) {
	parts := strings.Split(input, "/providers/Microsoft.Consumption/budgets")
	if len(parts) != 2 {
		return nil, fmt.Errorf("error parsing Budget resource ID: unexpected format: %q", input)
	}

	return &budgetResource{
		Scope:      parts[0],
		BudgetName: parts[1],
	}, nil
}
