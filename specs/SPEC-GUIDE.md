# SPEC Guide — Como escrever SPECS (Spec Driven Development)

Este guia existe para que **outras IAs/agentes** consigam escrever cada SPEC de forma **consistente**, **testável** e focada em **O QUE** (não em “como implementar”).

## O que é uma SPEC neste projeto

- Uma SPEC define **comportamentos e resultados observáveis** do produto/feature.
- Uma SPEC **não** define stack, arquitetura, banco, frameworks, endpoints, etc.
- Uma SPEC deve ser **independentemente implementável e testável** (MVP slice).

## Template obrigatório (copie e preencha)

Use o template abaixo como base (adapte nomes/quantidade de histórias conforme necessário).

```markdown
# Feature Specification: [FEATURE NAME]

**Created**: [DATE]
**PRD Base**: [listar seções do PRD, ex.: §9.1, §10 (R1)]

## Caso de uso *(mandatory)*

[Descreva o problema a ser resolvido e os fluxos do usuário. Foque no “porquê” e “o que acontece”.]

## User Scenarios & Testing *(mandatory)*

### User Story 1 - [Brief Title] (Priority: P1)

[Descreva a jornada em linguagem natural, Telegram-first quando aplicável.]

**Why this priority**: [valor para o usuário/negócio + dependências]

**Independent Test**: [como testar esta fatia sozinha]

**Acceptance Scenarios**:

1. **Scenario**: [nome]
   - **Given** [estado inicial]
   - **When** [ação]
   - **Then** [resultado observável]

2. **Scenario**: [nome]
   - **Given** ...
   - **When** ...
   - **Then** ...

---

### User Story 2 - [Brief Title] (Priority: P2)
...

### Edge Cases *(mandatory)*

- What happens when [condição de fronteira]?
- How does system handle [erro/entrada inválida]?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST ...
- **FR-002**: System MUST ...

### Non-Functional Requirements

- **NFR-001**: System MUST ...

### Key Entities *(include if feature involves data)*

- **[Entity]**: [o que representa + atributos/chaves em nível de produto]

## Acceptance Criteria *(mandatory)*

[Consolidar as condições objetivas de “pronto” para a SPEC, derivadas dos cenários.]

## Business Objectives *(mandatory)*

[Conectar a objetivos do PRD: consistência, qualidade, adaptação, baixa fricção, segurança psicológica, privacidade.]

## Error Handling *(mandatory)*

[Comportamento esperado em situações inesperadas/entradas inválidas/ausência de dados/usuário some.]

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: [métrica mensurável, agnóstica de tecnologia]
- **SC-002**: ...
```

## Regras de escrita (para manter “Spec Driven”)

- **Sem implementação**: não mencionar frameworks, linguagens, DB, provedores, filas, serviços.
- **Requisitos testáveis**: todo FR deve ser verificável por comportamento/resultado.
- **Telegram-first**: se o fluxo é conversacional, explicitar prompts/entradas/saídas em nível de UX (sem API).
- **Qualidade e evidência**: quando o PRD exigir “quality gates”, explicitar o que conta como evidência e como avaliar.
- **Dias ruins**: sempre prever comportamento quando tempo/energia são baixos (Plano B/C e MVD quando aplicável).
- **Segurança psicológica**: feedback firme sem punição (PRD RNF3), descrevendo o tom/comportamento esperado.
- **Privacidade por padrão**: coletar o mínimo; declarar o que é guardado e por quê (PRD RNF4).
- **Dúvidas**: se algo não estiver no PRD, marcar como **[NEEDS CLARIFICATION]** em vez de inventar.

## Checklist final (o agente deve validar antes de concluir)

- A SPEC **entrega valor sozinha** e tem pelo menos 1 user story P1.
- Há **cenários Given/When/Then** suficientes para cobrir happy path + 2+ edge cases.
- FRs cobrem **entradas**, **saídas**, **estados** e **condições de sucesso/falha**.
- NFRs relevantes (RNF1–RNF4 do PRD) foram considerados quando aplicável.
- O tratamento de erros inclui: entrada ausente/ambígua, evidência inválida, usuário some, sobrecarga.
- Success Criteria (SC-xxx) são **mensuráveis** e alinhados às métricas do PRD quando existirem.

## Catálogo rápido de tratamento de erros (reutilizável)

- **Entrada ausente/ambígua**: aplicar default mínimo e pedir confirmação com 1 pergunta curta.
- **Dia ruim**: oferecer MVD/Plano C e registrar exceção sem culpa.
- **Evidência inválida**: pedir reenvio ou aceitar alternativa equivalente definida pela SPEC.
- **“Marcar como feito” sem gate**: bloquear conclusão e oferecer o caminho mais curto para gerar evidência.
- **Sobrecarga**: impor limites e sugerir pausar/adiar metas (PRD §6.2).
- **Privacidade**: explicar retenção e permitir opt-out quando aplicável.

