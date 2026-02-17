# Feature Specification: Vida saudável — planejamento semanal + métricas + limites por dor/energia

**Created**: [DATE]  
**PRD Base**: §8.4

## Caso de uso *(mandatory)*

Definir como o sistema apoia hábitos de saúde com foco em sustentabilidade: planejamento semanal (força/cardio), ajuste por energia/dor/agenda e métricas simples.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST apoiar um plano semanal com sessões de força e atividade moderada ajustáveis.
- **FR-002**: System MUST registrar métricas semanais (minutos totais, consistência) e sinais de excesso (dor/fadiga).
- **FR-003**: System MUST adaptar recomendações quando houver dor/fadiga reportada.

### Non-Functional Requirements

- **NFR-001**: System MUST evitar abordagem “tudo ou nada”.

## Error Handling *(mandatory)*

- Dor/limitação: reduzir carga e sugerir alternativa segura em nível de produto (sem prescrição médica).

