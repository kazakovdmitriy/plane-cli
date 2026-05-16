package api

import "fmt"

func EndpointWorkItems() string                { return "work-items/" }
func EndpointWorkItem(id string) string         { return fmt.Sprintf("work-items/%s/", id) }
func EndpointWorkItemComments(id string) string { return fmt.Sprintf("work-items/%s/comments/", id) }
func EndpointCycles() string                    { return "cycles/" }
func EndpointCycle(id string) string            { return fmt.Sprintf("cycles/%s/", id) }
func EndpointModules() string                   { return "modules/" }
func EndpointModule(id string) string           { return fmt.Sprintf("modules/%s/", id) }
func EndpointStates() string                    { return "states/" }
func EndpointState(id string) string            { return fmt.Sprintf("states/%s/", id) }
func EndpointLabels() string                    { return "issue-labels/" }
func EndpointMembers() string                   { return "members/" }
func EndpointPages() string                     { return "pages/" }
func EndpointPage(id string) string             { return fmt.Sprintf("pages/%s/", id) }
