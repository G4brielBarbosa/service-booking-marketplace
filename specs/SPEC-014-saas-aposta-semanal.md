# Feature Specification: SaaS — “aposta semanal” (bloco profundo) + microblocos

**Created**: [DATE]  
**PRD Base**: §8.6, §6.2

## Caso de uso *(mandatory)*

Definir como o sistema garante progresso consistente no SaaS sem canibalizar fundamentos: 1 bloco profundo semanal + 1–2 microblocos, sempre com definição de pronto e resultado esperado.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST suportar planejamento e execução de 1 bloco profundo semanal com resultado esperado definido.
- **FR-002**: System MUST suportar 1–2 microblocos de manutenção sem sobrecarregar o plano diário.
- **FR-003**: System MUST exigir definição de “pronto” para contabilizar o bloco como concluído.

## Error Handling *(mandatory)*

- Overload: reduzir escopo do bloco e manter apenas microblocos mínimos.

