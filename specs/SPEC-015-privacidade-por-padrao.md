# Feature Specification: Privacidade por padrão — dados mínimos e transparência

**Created**: [DATE]  
**PRD Base**: 11 (RNF4), 10 (R6)

## Caso de uso *(mandatory)*

Definir como o sistema aplica privacidade por padrão: coleta mínima, clareza do que é guardado e por quê, e controle do usuário sobre retenção/opt-out, sem prejudicar as funções essenciais.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST declarar claramente quais dados são armazenados para cada fluxo (check-in, evidência, métricas) e por qual objetivo.
- **FR-002**: System MUST permitir ao usuário revisar seus dados registrados (resumos) e solicitar remoção/limitação conforme definido.
- **FR-003**: System MUST minimizar coleta e evitar dados sensíveis não necessários.

### Non-Functional Requirements

- **NFR-001**: System MUST coletar o mínimo necessário e ser transparente (PRD RNF4).

## Error Handling *(mandatory)*

- Pedido de armazenamento de conteúdo sensível: alertar e oferecer alternativa de registro mínimo.

