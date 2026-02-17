# Feature Specification: Onboarding e Diagnóstico Leve (1–2 semanas)

**Created**: [DATE]  
**PRD Base**: §5.1, §14

## Caso de uso *(mandatory)*

Definir como o usuário configura o assistente pela primeira vez e como o sistema estabelece uma linha de base (sono/inglês/java/autoestima/energia/tempo) para calibrar planos e medir progresso.

## Instruções para o agente que vai escrever esta SPEC

- Use `specs/SPEC-GUIDE.md` como padrão.
- Foque em **O QUE** deve ser coletado/produzido ao final do onboarding e como isso habilita o restante do produto.
- Não invente instrumentos clínicos; se precisar de algo não especificado no PRD, marque como **[NEEDS CLARIFICATION]**.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Completar onboarding mínimo e receber MVD (Priority: P1)

[Descrever a jornada.]

**Why this priority**: [valor]

**Independent Test**: [como testar só o onboarding]

**Acceptance Scenarios**:

1. **Scenario**: Onboarding mínimo completo
   - **Given** usuário novo sem dados prévios
   - **When** responde às perguntas de metas, restrições e preferências mínimas
   - **Then** o sistema registra uma baseline mínima e define um MVD

2. **Scenario**: Onboarding parcial (tempo curto)
   - **Given** usuário com pouco tempo
   - **When** responde somente ao mínimo obrigatório
   - **Then** o sistema conclui onboarding “mínimo” e agenda/coleta o restante em etapas

### Edge Cases *(mandatory)*

- What happens when o usuário não sabe definir metas anuais com clareza?
- How does system handle respostas conflitantes (ex.: “pouco tempo” e “muitas metas intensivas”)?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST coletar metas anuais suportadas e permitir que o usuário selecione quais estarão ativas no ciclo atual.
- **FR-002**: System MUST coletar restrições operacionais mínimas (tempo, horários, limitações) para calibrar planos.
- **FR-003**: System MUST capturar uma linha de base mínima para: sono, inglês, java, autoestima e contexto (energia/estresse).
- **FR-004**: System MUST definir um “mínimo viável diário (MVD)” que funcione em dia ruim.
- **FR-005**: System MUST produzir um resumo do onboarding que o usuário consiga revisar depois.
- **FR-006**: System MUST permitir retomada do onboarding caso seja interrompido.

### Non-Functional Requirements

- **NFR-001**: System MUST minimizar digitação e manter interação curta (PRD RNF1).
- **NFR-002**: System MUST coletar apenas o mínimo necessário e explicar por quê (PRD RNF4).

### Key Entities *(include if feature involves data)*

- **OnboardingSession**: status, respostas mínimas, pendências, timestamps.
- **BaselineSnapshot**: baseline por domínio (sono/inglês/java/autoestima/contexto).
- **MinimumViableDaily (MVD)**: lista curta de ações sustentáveis.

## Acceptance Criteria *(mandatory)*

- Ao final do onboarding mínimo, existe baseline mínima registrada + MVD definido + resumo revisável.

## Business Objectives *(mandatory)*

- Reduzir fricção inicial e criar base para planejamento adaptativo e métricas (PRD §§5.1–5.2).

## Error Handling *(mandatory)*

- Se o usuário não responder: manter sessão e permitir retomar sem perda.
- Se houver dados insuficientes: completar onboarding mínimo e marcar lacunas como **[NEEDS CLARIFICATION]** para coleta posterior.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Usuário completa onboarding mínimo em \u2264 10 minutos totais (somando interações).
- **SC-002**: Usuário consegue revisar o resumo do onboarding a qualquer momento.

