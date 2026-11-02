# Changelog

Todas as alterações relevantes deste projeto serão documentadas neste arquivo.

O projeto segue o padrão [Keep a Changelog](https://keepachangelog.com/pt-BR/1.1.0/)
e utiliza [Semantic Versioning](https://semver.org/lang/pt-BR/).

## [Unreleased]

## [1.0.0] - 2026-11-02

### Adicionado

- Estrutura inicial do microsserviço de multiplicação em Go.
- Endpoint REST `POST /multiplicar`.
- Suporte à multiplicação de números inteiros e decimais.
- Validação dos operandos obrigatórios.
- Tratamento de erros com respostas JSON padronizadas.
- Testes automatizados utilizando `testing` e `httptest`.
- Dockerfile multi-stage para conteinerização.
- Documentação de instalação, execução e testes.
- Contrato inicial da API.

### Corrigido

- Validação do corpo JSON para rejeitar valores `null`, arrays e strings.
- Rejeição de campos desconhecidos nas requisições.
- Validação de múltiplos objetos JSON no corpo da requisição.
- Ampliação da cobertura de testes automatizados para 14 cenários.

### Validação

- 14 cenários de testes automatizados aprovados.
- Análise estática executada com `go vet ./...`.
- Endpoint REST validado localmente.
- Imagem Docker reconstruída e validada após as correções.
