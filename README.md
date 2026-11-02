# Microsserviço de Multiplicação — Time 3

Microsserviço REST desenvolvido em Go para a Calculadora Distribuída, como parte da atividade de Gerenciamento de Configuração e Mudanças de Software (GCMS).

## Tecnologias

- Go
- Biblioteca padrão `net/http`
- Docker com multi-stage build
- Testes automatizados com `testing` e `httptest`

## Endpoint

**POST /multiplicar**

Realiza a multiplicação de dois operandos numéricos.

### Requisição

```json
{
  "a": 10,
  "b": 5
}
```

### Resposta

HTTP `200 OK`

```json
{
  "resultado": 50
}
```

## Executando localmente

É necessário ter Go instalado em versão compatível com o `go.mod`.

```bash
go run ./cmd/server
```

O serviço estará disponível em `http://localhost:8003`.

## Executando com Docker

```bash
docker build -t time3-multiplicacao:dev .
docker run --rm -p 8003:8003 time3-multiplicacao:dev
```

## Testando o endpoint

```bash
curl -X POST http://localhost:8003/multiplicar \
  -H "Content-Type: application/json" \
  -d '{"a":10,"b":5}'
```

## Testes automatizados

```bash
go test -count=1 -v ./...
go vet ./...
```

Os testes incluem operações com inteiros, decimais, negativos, zero e validação de requisições inválidas.

## Contrato de API

O serviço segue o contrato compartilhado entre as equipes, documentado em `docs/API_CONTRACT.md`.

Alterações incompatíveis no contrato devem ser discutidas e registradas por meio de Issues e Pull Requests.

## Versionamento

O projeto adota Semantic Versioning (SemVer).

A versão `v1.0.0` será publicada após aprovação do contrato, revisão do Pull Request e validação da integração.

## Integração

Este microsserviço integra a Calculadora Distribuída, composta por quatro serviços independentes de adição, subtração, multiplicação e divisão, orquestrados posteriormente pelo Docker Compose.
