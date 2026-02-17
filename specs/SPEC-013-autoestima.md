# Feature Specification: Autoestima — registros curtos + revisão de padrões + micro-exposições

**Created**: [DATE]  
**PRD Base**: §8.5, 11 (RNF3)

## Caso de uso *(mandatory)*

Definir como o sistema reduz autocrítica paralisante e aumenta autoeficácia com registros curtos (gatilho → pensamento → resposta alternativa), revisão semanal de padrões e micro-exposições (“ações de coragem”).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST permitir registro curto quando houver autocrítica forte (gatilho, pensamento, resposta alternativa).
- **FR-002**: System MUST conduzir revisão semanal de padrões e escolher 1 experimento para a semana seguinte.
- **FR-003**: System MUST registrar “ações de coragem” por semana e refletir tendência.

### Non-Functional Requirements

- **NFR-001**: System MUST manter segurança psicológica e tom não punitivo (PRD RNF3).

## Error Handling *(mandatory)*

- Usuário em alta emoção: oferecer intervenção mínima e postergar reflexão longa.

