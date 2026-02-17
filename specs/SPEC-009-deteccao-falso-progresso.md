# Feature Specification: Detecção de falhas reais (“falso progresso”) e reforço automático

**Created**: [DATE]  
**PRD Base**: 10 (R3), 13 (riscos)

## Caso de uso *(mandatory)*

Definir como o sistema detecta quando o usuário “fez” mas não houve aprendizagem/progresso real (ex.: repetição de erros, baixa qualidade em rubrica) e como reage com reforço, ajuste de plano e feedback.

## Instruções para o agente que vai escrever esta SPEC

- Use `specs/SPEC-GUIDE.md`.
- Referencie `SPEC-003` (gates) e `SPEC-016` (métricas).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST identificar padrões de recorrência de erro/baixa qualidade ao longo do tempo.
- **FR-002**: System MUST acionar um mecanismo de reforço (ex.: tarefa direcionada) quando um padrão atingir limiar definido.
- **FR-003**: System MUST ajustar recomendações do plano com base na evidência (Evidência > Intuição).

### Non-Functional Requirements

- **NFR-001**: System MUST manter feedback não punitivo (PRD RNF3).

## Error Handling *(mandatory)*

- Evidência insuficiente: não concluir aprendizado e solicitar mínima evidência adicional.

