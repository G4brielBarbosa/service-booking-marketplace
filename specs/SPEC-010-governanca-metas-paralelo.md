# Feature Specification: Personalização progressiva & governança de metas em paralelo (limites de overload)

**Created**: [DATE]  
**PRD Base**: §6.2, 10 (R7), 13 (riscos)

## Caso de uso *(mandatory)*

Definir como o sistema suporta metas em paralelo com limites explícitos (ex.: no máximo 2 metas intensivas por ciclo), iniciando simples e aumentando sofisticação conforme aprende padrões do usuário.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST impor limite de metas intensivas ativas por ciclo (base: PRD §6.2).
- **FR-002**: System MUST permitir pausar/retomar metas com registro do motivo.
- **FR-003**: System MUST aumentar personalização progressivamente com base em dados observados (PRD R7).

### Non-Functional Requirements

- **NFR-001**: System MUST reduzir overload e frustração (PRD §13).

## Error Handling *(mandatory)*

- Conflitos de prioridade: solicitar escolha objetiva (1 prioridade absoluta) e propor ajustes.

