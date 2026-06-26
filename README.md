# go-postgres-multi-client-proxy
This example demonstrates a simple Go-based Multi-Client Proxy (MCP) for Postgres. It listens on a specified port, accepts multiple client connections, and forwards them to a backend Postgres database. The proxy ensures that all client and backend connections are properly closed, illustrating a basic 'zero-leak' approach to resource management.
