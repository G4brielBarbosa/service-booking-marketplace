# Feature Specification: Backlog Inteligente e Priorização baseada em lacunas observadas

**Created**: [DATE]  
**PRD Base**: §5.2, 10 (R4)

## Caso de uso *(mandatory)*

Definir como o sistema mantém e prioriza um backlog de próximos passos baseado em lacunas observadas (erros recorrentes, baixa consistência, baixa qualidade), e não apenas em vontade momentânea.

## Instruções para o agente que vai escrever esta SPEC

- Use `specs/SPEC-GUIDE.md`.
- Relacione com `SPEC-009` (falhas reais) e `SPEC-016` (métricas/registros).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Receber sugestões acionáveis de backlog (Priority: P1)

**Acceptance Scenarios**:

1. **Scenario**: Lacuna detectada gera item de backlog
   - **Given** usuário errou o mesmo padrão 3 vezes (inglês/java) ou falhou em consistência
   - **When** o sistema atualiza backlog
   - **Then** um item de reforço é criado e priorizado

### Edge Cases *(mandatory)*

- What happens when existem itens demais (overload)?
- How does system handle lacunas em múltiplas metas intensivas ao mesmo tempo?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST manter uma lista de backlog com itens derivados de lacunas observadas.
- **FR-002**: System MUST priorizar backlog com base em impacto esperado e contexto (tempo/energia).
- **FR-003**: System MUST limitar quantidade de itens ativos para evitar overload.

### Non-Functional Requirements

- **NFR-001**: System MUST manter simplicidade e baixa fricção (PRD RNF1).

### Key Entities *(include if feature involves data)*

- **BacklogItem**: origem (lacuna), prioridade, status, data de criação, evidência associada.

## Error Handling *(mandatory)*

- Dados insuficientes: marcar recomendações como **[NEEDS CLARIFICATION]** e coletar mínimo no próximo check-in.

