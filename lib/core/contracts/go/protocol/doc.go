// Package protocol is the wire contract between the Microruntimes connected
// agent and the server it supervises runtimes for: an MCP endpoint at
// <console-url>/api/mcp, called with JSON-RPC 2.0 tools/call requests.
//
// Every operation is driven by authenticated console intent, never by
// repository prose. A field that was not measured is null, never zero.
package protocol
