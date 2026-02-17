# Feature Specification: Rotina Diária (Telegram-first) — Check-in + Plano A/B/C + Execução Guiada

**Created**: [DATE]  
**PRD Base**: §§5.2, 5.3, 9.1, 10 (R1, R6), 11 (RNF1, RNF2, RNF3), §14

## Caso de uso *(mandatory)*

Definir como o usuário faz um check-in diário rápido e como o sistema retorna um plano executável (A/B/C) com tarefas priorizadas, instruções objetivas e baixa fricção, adequado para dias bons e ruins.

## Instruções para o agente que vai escrever esta SPEC

- Use `specs/SPEC-GUIDE.md`.
- Trate “Telegram-first” como **UX conversacional**: entradas (tempo/energia/estresse) e saídas (plano) em linguagem natural.
- Não detalhar implementações (ex.: bot commands); foque nos comportamentos.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Receber plano diário adaptado por tempo/energia (Priority: P1)

[Descrever jornada.]

**Why this priority**: [valor]

**Independent Test**: [simular check-in e verificar plano retornado]

**Acceptance Scenarios**:

1. **Scenario**: Check-in completo e plano A recomendado
   - **Given** usuário informa tempo \u2265 30 min e energia alta
   - **When** envia check-in diário
   - **Then** recebe plano com 1 prioridade absoluta, 1–2 complementares e 1 fundação

2. **Scenario**: Energia baixa e ativação do Plano C
   - **Given** usuário informa energia muito baixa e pouco tempo
   - **When** envia check-in diário
   - **Then** recebe MVD/Plano C que mantém consistência sem burnout

### Edge Cases *(mandatory)*

- What happens when o usuário não informa tempo ou energia?
- How does system handle usuário pedindo “muitas tarefas” em um dia ruim?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST coletar tempo disponível e energia/estresse (mínimo) para escolher Plano A/B/C.
- **FR-002**: System MUST retornar um plano com: 1 prioridade absoluta, 1–2 complementares e 1 tarefa de fundação (quando aplicável).
- **FR-003**: System MUST fornecer instruções objetivas: o que fazer, por quanto tempo e qual critério de qualidade (se aplicável).
- **FR-004**: System MUST suportar planos alternativos quando o usuário perder o horário ou reduzir tempo disponível.
- **FR-005**: System MUST manter a interação diária curta e previsível (PRD RNF1).

### Non-Functional Requirements

- **NFR-001**: System MUST ser robusto a dias ruins, sempre oferecendo um plano mínimo (PRD RNF2).
- **NFR-002**: System MUST manter feedback honesto sem humilhação (PRD RNF3).

### Key Entities *(include if feature involves data)*

- **DailyCheckIn**: tempo, energia, estresse/humor, timestamp.
- **DailyPlan**: plano (A/B/C), lista de tarefas, prioridades, critérios.

## Acceptance Criteria *(mandatory)*

- Com um check-in mínimo, o usuário sempre recebe um plano executável (A/B/C) com prioridades claras.

## Business Objectives *(mandatory)*

- Reduzir fricção e aumentar consistência diária (PRD §§5.2–5.3, §9.1).

## Error Handling *(mandatory)*

- Se entradas forem omitidas: aplicar defaults mínimos e pedir 1 confirmação curta.
- Se usuário sumir: permitir retomada e oferecer plano atualizado pelo tempo restante.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Usuário conclui check-in e recebe plano em \u2264 2 minutos de interação.
- **SC-002**: Em dias de baixa energia, o sistema ainda mantém consistência (taxa de dias com MVD concluído).

