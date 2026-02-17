# Feature Specification: Revisão Semanal — Painel + 3 decisões + alvos da semana

**Created**: [DATE]  
**PRD Base**: §§5.2, 5.5, 9.2, 10 (R5), §14

## Caso de uso *(mandatory)*

Definir como o sistema conduz uma revisão semanal curta (10–20 min) que sintetiza métricas e evidencia gargalos, levando o usuário a 3 decisões claras (manter/ajustar/pausar) e à seleção de “alvos da semana”.

## Instruções para o agente que vai escrever esta SPEC

- Use `specs/SPEC-GUIDE.md`.
- Não desenhar dashboard técnico; descreva apenas as informações mínimas que devem ser apresentadas e decisões resultantes.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Completar revisão semanal e sair com decisões acionáveis (Priority: P1)

**Acceptance Scenarios**:

1. **Scenario**: Revisão semanal com painel mínimo
   - **Given** existe histórico da semana (check-ins, planos, rubricas, sono/energia)
   - **When** usuário inicia revisão semanal
   - **Then** o sistema apresenta consistência por meta, qualidade (rubricas) e gargalos principais

2. **Scenario**: Decisões e alvos definidos
   - **Given** o painel foi apresentado
   - **When** usuário escolhe manter/ajustar/pausar
   - **Then** o sistema registra decisões e define alvos da semana (ex.: 1 em inglês, 1 em java, 1 em sono/saúde)

### Edge Cases *(mandatory)*

- What happens when há poucos dados na semana (usuário ausente)?
- How does system handle sinais de overload (muitas metas intensivas ativas)?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST compilar um painel semanal mínimo: consistência por meta, qualidade/rubricas, gargalos, sono/energia (tendência).
- **FR-002**: System MUST guiar o usuário a 3 decisões: manter, ajustar ou pausar.
- **FR-003**: System MUST permitir selecionar “alvos da semana” e registrar o compromisso.
- **FR-004**: System MUST sugerir próximos experimentos/intervenções com base nos gargalos observados.

### Non-Functional Requirements

- **NFR-001**: System MUST manter revisão concluível em 10–20 min (PRD §9.2).
- **NFR-002**: System MUST manter tom não punitivo (PRD RNF3).

### Key Entities *(include if feature involves data)*

- **WeeklyReview**: período, painel, decisões, alvos, recomendações.
- **GoalStatus**: consistência, qualidade, gargalos.

## Acceptance Criteria *(mandatory)*

- Revisão semanal resulta em: painel + 3 decisões registradas + alvos da semana definidos.

## Business Objectives *(mandatory)*

- Garantir adaptação contínua e foco em evidência (PRD §§5.5, 9.2).

## Error Handling *(mandatory)*

- Poucos dados: explicitar limitação e definir a menor ação para coletar baseline na próxima semana.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Percentual de semanas com revisão completa.
- **SC-002**: Melhora de métricas alvo (consistência/qualidade/sono) após 2–4 ciclos.

