# Contrato de API — Calculadora Distribuída

**Versão do contrato:** `v1.0.0`  
**Status:** Contrato inicial aprovado pelos quatro times; extensões de erros aguardando alinhamento  
**Disciplina:** Gerenciamento de Configuração e Mudança de Software (GCMS)

## 1. Objetivo

Estabelecer o contrato de comunicação REST/JSON compartilhado pelos quatro microsserviços independentes da Calculadora Distribuída.

O contrato define os endpoints, as estruturas de requisição e resposta, as regras de validação, o tratamento de erros e os procedimentos de gerenciamento de mudanças.

Cada equipe possui autonomia para escolher a linguagem de programação, o framework e a estrutura interna de implementação, desde que respeite o contrato estabelecido.

## 2. Endpoints

| Time | Operação | Método HTTP | Endpoint |
| --- | --- | --- | --- |
| Time 1 | Adição | `POST` | `/somar` |
| Time 2 | Subtração | `POST` | `/subtrair` |
| Time 3 | Multiplicação | `POST` | `/multiplicar` |
| Time 4 | Divisão | `POST` | `/dividir` |

Cada endpoint pertence ao microsserviço de responsabilidade da respectiva equipe.

Os endereços e as portas dos serviços serão definidos durante a integração por meio do Docker Compose, sem alterar os endpoints estabelecidos neste documento.

## 3. Contrato de requisição

Todas as operações recebem dois operandos numéricos por meio de um objeto JSON.

**Método HTTP:** `POST`  
**Content-Type:** `application/json`

### 3.1 Estrutura da requisição

```json
{
  "a": 10,
  "b": 5
}
```

### 3.2 Definição dos campos

| Campo | Tipo | Obrigatório | Descrição |
| --- | --- | --- | --- |
| `a` | number | Sim | Primeiro operando |
| `b` | number | Sim | Segundo operando |

### 3.3 Regras de validação

- Os campos `a` e `b` são obrigatórios.
- Ambos devem conter valores numéricos JSON.
- São aceitos números inteiros, decimais, positivos, negativos e zero.
- Strings numéricas, como `"10"`, não são aceitas.
- Valores `null` não são aceitos como operandos.
- Campos obrigatórios ausentes não devem receber valores padrão.
- O corpo deve conter exatamente um objeto JSON válido.
- Arrays e valores JSON primitivos não são aceitos como corpo da requisição.
- Campos desconhecidos devem ser rejeitados.
- O cabeçalho `Content-Type` deve identificar o formato `application/json`.

Requisições que não atendam a essas regras devem retornar os erros definidos na seção 5.

## 4. Contrato de resposta

### 4.1 Resposta de sucesso

**Status HTTP:** `200 OK`  
**Content-Type:** `application/json`

```json
{
  "resultado": 50
}
```

### 4.2 Definição dos campos

| Campo | Tipo | Obrigatório | Descrição |
| --- | --- | --- | --- |
| `resultado` | number | Sim | Resultado da operação matemática |

### 4.3 Regras de saída

- O campo `resultado` deve estar presente em todas as respostas de sucesso.
- O valor retornado deve ser um número JSON finito.
- A resposta deve utilizar a mesma estrutura em todos os microsserviços.
- O resultado não deve sofrer arredondamento deliberado.
- Valores como `NaN`, `Infinity` e `-Infinity` não podem ser retornados como resultados JSON.

### 4.4 Exemplos das operações

| Endpoint | Requisição | Resposta (`200 OK`) |
| --- | --- | --- |
| `/somar` | `{"a":10,"b":5}` | `{"resultado":15}` |
| `/subtrair` | `{"a":10,"b":5}` | `{"resultado":5}` |
| `/multiplicar` | `{"a":10,"b":5}` | `{"resultado":50}` |
| `/dividir` | `{"a":10,"b":5}` | `{"resultado":2}` |

## 5. Tratamento de erros

Todas as respostas de erro devem utilizar a seguinte estrutura JSON:

```json
{
  "erro": {
    "codigo": "INVALID_INPUT",
    "mensagem": "Os operandos a e b são obrigatórios."
  }
}
```

### 5.1 Estrutura da resposta de erro

| Campo | Tipo | Obrigatório | Descrição |
| --- | --- | --- | --- |
| `erro` | object | Sim | Objeto com os detalhes do erro |
| `erro.codigo` | string | Sim | Identificador padronizado do erro |
| `erro.mensagem` | string | Sim | Descrição legível do problema |

### 5.2 Códigos de erro

| Situação | Status HTTP | `erro.codigo` |
| --- | --- | --- |
| Campo obrigatório ausente | `400 Bad Request` | `INVALID_INPUT` |
| Operando não numérico | `400 Bad Request` | `INVALID_INPUT` |
| Operando com valor `null` | `400 Bad Request` | `INVALID_INPUT` |
| JSON sintaticamente inválido | `400 Bad Request` | `INVALID_JSON` |
| Corpo que não seja objeto JSON | `400 Bad Request` | `INVALID_JSON` |
| Múltiplos valores JSON no corpo | `400 Bad Request` | `INVALID_JSON` |
| Campo desconhecido | `400 Bad Request` | `INVALID_JSON` |
| Método HTTP não permitido | `405 Method Not Allowed` | `METHOD_NOT_ALLOWED` |
| Tipo de conteúdo não suportado | `415 Unsupported Media Type` | `UNSUPPORTED_MEDIA_TYPE` |
| Divisão por zero | `422 Unprocessable Content` | `DIVISION_BY_ZERO` |
| Resultado numérico não finito | `422 Unprocessable Content` | `INVALID_RESULT` |
| Erro inesperado do servidor | `500 Internal Server Error` | `INTERNAL_ERROR` |

**Nota:** `METHOD_NOT_ALLOWED`, `UNSUPPORTED_MEDIA_TYPE` e `INVALID_RESULT` são extensões propostas com base na implementação do Time 3 e precisam ser homologadas como parte do contrato compartilhado.

### 5.3 Exemplos de erros

**Exemplo 1 — Operando obrigatório ausente**

Requisição:

```json
{
  "a": 10
}
```

Resposta — `400 Bad Request`:

```json
{
  "erro": {
    "codigo": "INVALID_INPUT",
    "mensagem": "Os operandos a e b são obrigatórios."
  }
}
```

**Exemplo 2 — Operando não numérico**

Requisição:

```json
{
  "a": "dez",
  "b": 5
}
```

Resposta — `400 Bad Request`:

```json
{
  "erro": {
    "codigo": "INVALID_INPUT",
    "mensagem": "Os operandos devem ser números."
  }
}
```

**Exemplo 3 — JSON malformado**

Resposta — `400 Bad Request`:

```json
{
  "erro": {
    "codigo": "INVALID_JSON",
    "mensagem": "JSON inválido."
  }
}
```

**Exemplo 4 — Divisão por zero**

Endpoint: `POST /dividir`

Requisição:

```json
{
  "a": 10,
  "b": 0
}
```

Resposta — `422 Unprocessable Content`:

```json
{
  "erro": {
    "codigo": "DIVISION_BY_ZERO",
    "mensagem": "Não é permitido dividir por zero."
  }
}
```

**Exemplo 5 — Resultado fora do intervalo permitido**

Resposta — `422 Unprocessable Content`:

```json
{
  "erro": {
    "codigo": "INVALID_RESULT",
    "mensagem": "Resultado fora do intervalo permitido."
  }
}
```

### 5.4 Regras gerais de tratamento de erros

- Todas as respostas de erro devem utilizar `Content-Type: application/json`.
- Os códigos de erro devem seguir os identificadores definidos neste contrato.
- As mensagens devem ser claras e consistentes.
- Erros de validação do cliente não devem retornar HTTP `500`.
- Erros inesperados do servidor devem retornar HTTP `500`, sem exposição de informações internas sensíveis.
- O status HTTP e o código retornado devem corresponder à situação identificada.

## 6. Versionamento e gerenciamento de mudanças

O contrato inicial foi aprovado pelas quatro equipes, estabelecendo a versão `v1.0.0` como referência de integração.

A publicação das versões de software deve seguir o padrão Semantic Versioning (SemVer).

### 6.1 Política de versionamento

| Tipo | Exemplo | Descrição |
| --- | --- | --- |
| MAJOR | `v2.0.0` | Mudanças incompatíveis com versões anteriores |
| MINOR | `v1.1.0` | Adição de funcionalidades compatíveis |
| PATCH | `v1.0.1` | Correções que preservam a compatibilidade |

### 6.2 Mudanças incompatíveis

São exemplos de alterações que exigem análise de breaking change:

- Renomear endpoints existentes.
- Alterar métodos HTTP estabelecidos.
- Renomear campos obrigatórios.
- Alterar tipos de dados esperados.
- Modificar a estrutura de respostas.
- Remover códigos de erro definidos no contrato.
- Tornar obrigatórios campos anteriormente opcionais.

### 6.3 Processo formal de mudanças

Toda solicitação de alteração contratual deve seguir o fluxo:

1. Criar uma Issue descrevendo a alteração proposta.
2. Documentar sua justificativa e os impactos previstos.
3. Discutir a mudança com as equipes envolvidas.
4. Definir a classificação da alteração conforme o SemVer.
5. Implementar a alteração em uma branch específica.
6. Abrir um Pull Request vinculado à Issue.
7. Realizar revisão e validação da mudança.
8. Atualizar o contrato e a documentação correspondente.
9. Registrar a alteração no `CHANGELOG.md`.
10. Publicar a versão correspondente após aprovação.

Mudanças incompatíveis não devem ser incorporadas unilateralmente por uma equipe.

## 7. Integração e infraestrutura

### 7.1 Microsserviços independentes

Cada equipe é responsável por:

- Implementar o microsserviço correspondente.
- Manter seu próprio repositório Git.
- Disponibilizar um Dockerfile na raiz.
- Implementar o endpoint definido neste contrato.
- Garantir o formato padronizado de requisição e resposta.
- Documentar a execução no `README.md`.
- Registrar alterações no `CHANGELOG.md`.
- Manter rastreabilidade por Issues, Pull Requests e tags.

### 7.2 Docker

Cada microsserviço deve ser capaz de executar independentemente em um container Docker.

A imagem deve conter os recursos necessários à execução, sem depender de ferramentas de desenvolvimento instaladas no ambiente do avaliador.

### 7.3 Docker Compose

A integração entre os quatro microsserviços será realizada por meio de um arquivo central `docker-compose.yml`.

O ambiente integrado deverá permitir a inicialização por meio do comando:

```bash
docker compose up
```

A orquestração deve preservar o funcionamento independente de cada microsserviço e os endpoints definidos no contrato.

### 7.4 Responsabilidades por equipe

| Time | Responsabilidade |
| --- | --- |
| Time 1 | Implementar e disponibilizar o serviço de adição |
| Time 2 | Implementar e disponibilizar o serviço de subtração |
| Time 3 | Implementar e disponibilizar o serviço de multiplicação |
| Time 4 | Implementar e disponibilizar o serviço de divisão |
| Integração | Configurar e validar a execução conjunta dos serviços |

## 8. Governança e controle de configuração

O desenvolvimento deve manter rastreabilidade das alterações realizadas nos repositórios.

### 8.1 Branches

As alterações devem ser desenvolvidas em branches específicas, evitando modificações diretas na branch principal.

### 8.2 Issues

As Issues devem documentar funcionalidades, correções, solicitações de alteração e problemas identificados.

### 8.3 Pull Requests

As alterações devem ser integradas por Pull Requests, permitindo revisão e registro das modificações realizadas.

### 8.4 Tags e releases

As entregas estáveis devem possuir tags de versão conforme o SemVer.

A versão `v1.0.0` representa a primeira versão estável, após sua implementação e validação.

A aprovação do contrato não implica automaticamente a publicação das versões `v1.0.0` dos microsserviços.

## 9. Aprovação inicial do contrato

As quatro equipes aprovaram o contrato inicial de comunicação REST/JSON.

- [x] Time 1 — Adição
- [x] Time 2 — Subtração
- [x] Time 3 — Multiplicação
- [x] Time 4 — Divisão

**Estado do contrato inicial:** Aprovado.

**Extensões adicionais de tratamento de erros:** Pendentes de homologação entre as equipes.

## 10. Considerações finais

Este documento estabelece a referência de comunicação entre os quatro microsserviços da Calculadora Distribuída.

Todas as equipes devem respeitar os endpoints, os formatos JSON e as regras estabelecidas no contrato aprovado.

Qualquer alteração posterior deve seguir o processo formal de gerenciamento de mudanças, mantendo a documentação atualizada e a rastreabilidade no Git.

As extensões identificadas neste documento devem ser alinhadas entre as equipes antes de sua incorporação definitiva ao contrato compartilhado.
