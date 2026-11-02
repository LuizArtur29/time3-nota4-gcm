# Contrato de API — Calculadora Distribuída

**Versão proposta:** `v1.0.0`  
**Status:** Proposta inicial — pendente de aprovação dos quatro times  
**Disciplina:** Gerenciamento de Configuração e Mudança de Software (GCMS)

## 1. Objetivo

Estabelecer o contrato REST/JSON compartilhado pelos quatro microsserviços independentes da Calculadora Distribuída, antes do início da implementação.

O contrato define endpoints, estruturas de requisição e resposta, regras de validação, tratamento de erros e procedimento para mudanças. A linguagem de programação, o framework e a implementação interna ficam a critério de cada equipe.

## 2. Endpoints

| Time | Operação | Método | Endpoint |
| --- | --- | --- | --- |
| 1 | Adição | `POST` | `/somar` |
| 2 | Subtração | `POST` | `/subtrair` |
| 3 | Multiplicação | `POST` | `/multiplicar` |
| 4 | Divisão | `POST` | `/dividir` |

Cada endpoint pertence ao microsserviço do respectivo time. O endereço e a porta de cada serviço serão definidos na integração via Docker Compose e não fazem parte deste contrato inicial.

## 3. Requisição

Todas as operações recebem o mesmo corpo JSON:

```json
{
  "a": 10,
  "b": 5
}
```

### Regras de entrada

- `a`: campo obrigatório, do tipo número JSON.
- `b`: campo obrigatório, do tipo número JSON.
- São aceitos valores inteiros, decimais e negativos.
- Strings numéricas, como `"10"`, não são aceitas.
- Os serviços não devem atribuir valores padrão a campos ausentes.
- O corpo deve ser um objeto JSON válido.
- A requisição deve utilizar `Content-Type: application/json`.

## 4. Resposta de sucesso

**Status HTTP:** `200 OK`  
**Content-Type:** `application/json`

```json
{
  "resultado": 15
}
```

O campo `resultado` é obrigatório e contém um número JSON finito. O formato é idêntico nos quatro serviços; apenas o valor calculado muda. Não deve haver arredondamento deliberado do resultado.

### Exemplos

| Endpoint | Requisição | Resposta (`200 OK`) |
| --- | --- | --- |
| `/somar` | `{"a":10,"b":5}` | `{"resultado":15}` |
| `/subtrair` | `{"a":10,"b":5}` | `{"resultado":5}` |
| `/multiplicar` | `{"a":10,"b":5}` | `{"resultado":50}` |
| `/dividir` | `{"a":10,"b":5}` | `{"resultado":2}` |

## 5. Respostas de erro

Todas as respostas de erro devem seguir a estrutura:

```json
{
  "erro": {
    "codigo": "INVALID_INPUT",
    "mensagem": "Os campos 'a' e 'b' devem ser números."
  }
}
```

| Situação | Status HTTP | `erro.codigo` |
| --- | --- | --- |
| Campo obrigatório ausente | `400 Bad Request` | `INVALID_INPUT` |
| Operando não numérico | `400 Bad Request` | `INVALID_INPUT` |
| JSON sintaticamente inválido | `400 Bad Request` | `INVALID_JSON` |
| Divisão por zero | `422 Unprocessable Content` | `DIVISION_BY_ZERO` |
| Erro inesperado no servidor | `500 Internal Server Error` | `INTERNAL_ERROR` |

### Exemplo: divisão por zero

**Requisição:** `POST /dividir`

```json
{
  "a": 10,
  "b": 0
}
```

**Resposta:** `422 Unprocessable Content`

```json
{
  "erro": {
    "codigo": "DIVISION_BY_ZERO",
    "mensagem": "Não é permitido dividir por zero."
  }
}
```

Os códigos de erro e o formato da resposta são fixos. As mensagens devem ser claras e consistentes. Resultados que não possam ser representados como números JSON finitos devem retornar erro; a classificação específica desse caso será acordada antes de sua implementação, caso seja necessária.

## 6. Versionamento e gestão de mudanças

O contrato começará na versão `v1.0.0` **após aprovação das quatro equipes**. Até lá, esta é apenas uma proposta e não deve ser tratada como versão estável publicada.

As mudanças devem seguir o versionamento semântico (SemVer):

- **MAJOR:** mudanças incompatíveis no contrato, como renomeação de campos ou alteração da estrutura de resposta.
- **MINOR:** adições compatíveis de funcionalidades.
- **PATCH:** correções que preservem a compatibilidade.

Qualquer proposta de alteração deve ser registrada por meio de **Issue**, discutida entre as equipes e implementada por **Pull Request**, com revisão, histórico em `CHANGELOG.md` e tag de versão quando apropriado.

## 7. Integração e entregáveis

- Cada equipe implementará seu microsserviço em repositório próprio, com um `Dockerfile` na raiz.
- A integração geral ocorrerá por um `docker-compose.yml` central, que deverá permitir executar o sistema com `docker compose up`.
- Cada repositório deve documentar as instruções de uso em `README.md` e o histórico de alterações em `CHANGELOG.md`.
- Uma cópia aprovada deste documento deve estar acessível a todas as equipes, evitando divergência entre implementações.

## 8. Aprovação inicial

- [ ] Time 1 — Adição
- [ ] Time 2 — Subtração
- [ ] Time 3 — Multiplicação
- [ ] Time 4 — Divisão

**Observação:** Este documento contém propostas de padronização técnicas definidas para o projeto. O enunciado da atividade exige um contrato JSON previamente acordado, quatro microsserviços REST, conteinerização, Docker Compose e rastreabilidade de mudanças, mas não determina os nomes de campos, os métodos HTTP nem os códigos de erro específicos adotados aqui.
