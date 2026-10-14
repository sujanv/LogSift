package cluster

import "strings"

func ExtractVariables(template, message string) []string {
	tmplTokens := strings.Fields(template)
	msgTokens := strings.Fields(message)

	if len(tmplTokens) != len(msgTokens) {
		return nil
	}

	var vars []string
	for i := 0; i < len(tmplTokens); i++ {
		if tmplTokens[i] == "<*>" || tmplTokens[i] == "<*NUM*>" || tmplTokens[i] == "<*HEX*>" || tmplTokens[i] == "<*UUID*>" {
			vars = append(vars, msgTokens[i])
		}
	}
	return vars
}
