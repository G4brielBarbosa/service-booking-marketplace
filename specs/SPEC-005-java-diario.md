# Feature Specification: Java Diário — Prática Deliberada + Retrieval + Revisão de Erros

**Created**: [DATE]  
**PRD Base**: §8.2, §14, 10 (R2, R3)

## Caso de uso *(mandatory)*

Definir um loop diário de estudo/prática em Java com evidência por produção (código) e checagem de entendimento (retrieval), incluindo revisão e catalogação de erros para reforço ao longo do tempo.

## Instruções para o agente que vai escrever esta SPEC

- Use `specs/SPEC-GUIDE.md`.
- Referencie `SPEC-003` para gates/evidências e descreva como se aplica aqui.
- Não especificar IDE, repositórios, plataformas; descreva somente critérios e resultados.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Completar prática deliberada com “definição de feito” (Priority: P1)

**Acceptance Scenarios**:

1. **Scenario**: Exercício concluído com definição de pronto
   - **Given** usuário escolhe uma tarefa de prática (kata/exercício)
   - **When** finaliza a tarefa e fornece evidência mínima
   - **Then** a tarefa só é considerada concluída quando atender “definição de feito” e registrar aprendizado/erro principal

2. **Scenario**: Retrieval curto para checar entendimento
   - **Given** após a prática
   - **When** usuário responde mini-quiz/explica conceito sem consultar
   - **Then** o sistema registra desempenho e sugere reforço se necessário

### Edge Cases *(mandatory)*

- What happens when usuário não consegue terminar o exercício no tempo?
- How does system handle usuário repetindo o mesmo erro por 3+ sessões?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST suportar prática deliberada com restrição/objetivo claro e definição de “feito”.
- **FR-002**: System MUST exigir um registro curto de raciocínio (2–5 linhas) e do principal aprendizado/erro.
- **FR-003**: System MUST suportar retrieval (quiz/explicação) sem consulta e registrar resultado.
- **FR-004**: System MUST registrar erros recorrentes por categoria e permitir reforço direcionado.
- **FR-005**: System MUST aplicar quality gates para concluir tarefas de Java (ver `SPEC-003`).

### Non-Functional Requirements

- **NFR-001**: System MUST manter interação objetiva e de baixa fricção (PRD RNF1).

### Key Entities *(include if feature involves data)*

- **JavaPracticeSession**: tarefa, duração, evidência, definição de feito.
- **RetrievalCheck**: tópico, perguntas/itens, desempenho.
- **ErrorLogEntry**: categoria, descrição, correção, recorrência.

## Acceptance Criteria *(mandatory)*

- Tarefa de Java só conta como concluída quando o gate definido for satisfeito e houver registro do aprendizado/erro.

## Business Objectives *(mandatory)*

- Progresso mensurável por código e avaliações curtas, evitando “fiz mas não aprendi” (PRD §8.2).

## Error Handling *(mandatory)*

- Exercício incompleto: registrar progresso, definir próximo passo mínimo e não contar como “feito” sem gate.
- Evidência insuficiente: solicitar complemento mínimo.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Aumento do tempo semanal de prática deliberada sustentado.
- **SC-002**: Redução de erros recorrentes por categoria ao longo de semanas.

