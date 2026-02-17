# Feature Specification: Métricas & Registros — consistência, rubricas, tendências, erros recorrentes

**Created**: [DATE]  
**PRD Base**: §9.2, §§8.1–8.5 (métricas), 10 (R5)

## Caso de uso *(mandatory)*

Definir quais métricas e registros mínimos o sistema mantém para suportar planejamento adaptativo, quality gates e revisão semanal (consistência, rubricas, sono/energia, erros recorrentes), com baixa fricção.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST registrar consistência por meta (dias/semana) e permitir consulta rápida do progresso do dia/semana.
- **FR-002**: System MUST registrar rubricas de qualidade (quando aplicável) e tendências por semana.
- **FR-003**: System MUST registrar e consolidar erros recorrentes por domínio (inglês/java) e tendência de redução.
- **FR-004**: System MUST registrar sono/energia com tendência (regularidade, qualidade percebida, energia manhã).

### Non-Functional Requirements

- **NFR-001**: System MUST manter captura de dados com mínima digitação (PRD RNF1).
- **NFR-002**: System MUST aplicar privacidade por padrão (PRD RNF4).

## Error Handling *(mandatory)*

- Dados faltantes: lidar com lacunas sem penalizar e propor o mínimo para recuperar baseline.

