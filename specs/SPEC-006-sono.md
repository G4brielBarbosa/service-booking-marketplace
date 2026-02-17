# Feature Specification: Sono — Diário + Rotina Pré-sono + 1 Intervenção simples/semana

**Created**: [DATE]  
**PRD Base**: §6.1, §8.3, §14, 11 (RNF1, RNF2, RNF3, RNF4)

## Caso de uso *(mandatory)*

Definir como o sistema ajuda o usuário a melhorar sono com baixa fricção: diário mínimo, rotina pré-sono realista e intervenção semanal simples, monitorando tendências (regularidade, qualidade, energia).

## Instruções para o agente que vai escrever esta SPEC

- Use `specs/SPEC-GUIDE.md`.
- Não transformar em protocolo clínico; manter no escopo de produto do PRD.
- Descrever como “cumpriu sono” é contabilizado (gate) e como lidar com dias ruins.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Registrar diário de sono e obter feedback (Priority: P1)

**Acceptance Scenarios**:

1. **Scenario**: Diário preenchido ao acordar
   - **Given** novo dia
   - **When** usuário preenche diário de sono (30–60s)
   - **Then** o sistema registra métricas mínimas e dá feedback curto

2. **Scenario**: Rotina pré-sono mínima em dia ruim
   - **Given** usuário reporta baixa energia/estresse alto
   - **When** solicita plano mínimo
   - **Then** o sistema oferece rotina pré-sono pequena e viável

### Edge Cases *(mandatory)*

- What happens when o usuário esquece de preencher o diário por 2+ dias?
- How does system handle dias com horários muito irregulares?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST coletar diário de sono mínimo (horário dormir/levantar, qualidade percebida, despertares quando aplicável).
- **FR-002**: System MUST suportar definição de horário alvo e rotinas pré-sono pequenas e realistas.
- **FR-003**: System MUST definir o que conta como “cumpriu sono” (gate) de forma não punitiva.
- **FR-004**: System MUST oferecer 1 intervenção simples por semana e registrar aderência.
- **FR-005**: System MUST fornecer tendência semanal (regularidade, qualidade, energia).

### Non-Functional Requirements

- **NFR-001**: System MUST manter fricção mínima (PRD RNF1).
- **NFR-002**: System MUST ser robusto a dias ruins (PRD RNF2).
- **NFR-003**: System MUST manter segurança psicológica no feedback (PRD RNF3).
- **NFR-004**: System MUST aplicar privacidade por padrão no registro de dados (PRD RNF4).

### Key Entities *(include if feature involves data)*

- **SleepDiaryEntry**: dormir/levantar, qualidade, energia, observações curtas.
- **SleepRoutinePlan**: passos mínimos, horário alvo.
- **WeeklyIntervention**: tipo/descrição, objetivo, aderência.

## Acceptance Criteria *(mandatory)*

- Usuário consegue registrar diário e receber feedback; em dias ruins existe plano mínimo sem culpa.

## Business Objectives *(mandatory)*

- Tratar sono como fundação de execução e aprendizagem (PRD §6.1).

## Error Handling *(mandatory)*

- Ausência de dados: oferecer retomar com defaults e registrar lacuna sem penalizar.
- Dados incoerentes: pedir confirmação curta e seguir com plano mínimo.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Aumento de regularidade (redução de variância dos horários).
- **SC-002**: Melhora em qualidade percebida e/ou energia pela manhã ao longo de semanas.

