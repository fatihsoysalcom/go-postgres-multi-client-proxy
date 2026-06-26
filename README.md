# Go Postgres Multi Client Proxy

This example demonstrates a simple Go-based Multi-Client Proxy (MCP) for Postgres. It listens on a specified port, accepts multiple client connections, and forwards them to a backend Postgres database. The proxy ensures that all client and backend connections are properly closed, illustrating a basic 'zero-leak' approach to resource management.

## Language

`go`

## How to Run

1. Ensure you have Go installed (https://golang.org/doc/install).
2. Ensure a Postgres database is running and accessible (e.g., on `localhost:5432`).
3. Run the proxy: `go run main.go -port 6000 -backend localhost:5432`
4. Connect your Postgres client (e.g., `psql`) to `localhost:6000` instead of `localhost:5432`.

## Original Article

This example accompanies the Turkish article: [Go ile Sıfır Sızıntılı Postgres Çoklu İstemci Proxy (MCP) Ağ Geçidi Nasıl Oluşturulur?](https://fatihsoysal.com/blog/go-ile-sifir-sizintili-postgres-coklu-istemci-proxy-mcp-ag-gecidi-nasil-olusturulur/).

## License

MIT — see [LICENSE](LICENSE).
