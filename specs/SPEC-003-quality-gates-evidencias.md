# Feature Specification: Quality Gates & Evidência Mínima (aprendizagem e hábitos)

**Created**: [DATE]  
**PRD Base**: §5.4, §§8.1–8.3, 10 (R2, R3), 13 (riscos)

## Caso de uso *(mandatory)*

Definir critérios objetivos para considerar tarefas como “concluídas”, exigindo evidência proporcional ao objetivo (aprendizagem vs. hábito), e reduzindo “falso progresso”.

## Instruções para o agente que vai escrever esta SPEC

- Use `specs/SPEC-GUIDE.md`.
- Esta SPEC é **cross-cutting**: deve definir regras gerais que outras SPECS referenciam.
- Descrever rubricas e evidências em nível de produto (sem pipeline técnico).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Bloquear “feito” sem evidência mínima (Priority: P1)

**Acceptance Scenarios**:

1. **Scenario**: Usuário tenta concluir tarefa de aprendizagem sem evidência
   - **Given** existe uma tarefa de inglês/java com gate definido
   - **When** usuário tenta marcar como concluída sem evidência mínima
   - **Then** o sistema não aceita a conclusão e oferece o caminho mais curto para produzir evidência

2. **Scenario**: Evidência enviada mas insuficiente/ilegível
   - **Given** tarefa requer evidência
   - **When** usuário envia evidência inválida (ex.: texto vazio/áudio inaudível)
   - **Then** o sistema solicita alternativa equivalente ou reenvio

### Edge Cases *(mandatory)*

- What happens when o usuário discorda do gate e quer “marcar mesmo assim”?
- How does system handle evidência parcial (ex.: respondeu 1 de 3 perguntas de compreensão)?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST definir “quality gates” por tipo de tarefa (aprendizagem vs hábitos).
- **FR-002**: System MUST exigir evidência mínima para tarefas de aprendizagem antes de permitir conclusão.
- **FR-003**: System MUST registrar rubricas de qualidade quando aplicável (ex.: speaking 0–2 por dimensão).
- **FR-004**: System MUST identificar padrões de “falso progresso” (ex.: repetição de erro recorrente) e acionar reforço.
- **FR-005**: System MUST permitir caminhos alternativos de evidência quando o formato principal falhar (ex.: áudio → texto curto), se definido.

### Non-Functional Requirements

- **NFR-001**: System MUST manter fricção proporcional (evidência mínima, não burocracia) (PRD RNF1).
- **NFR-002**: System MUST aplicar feedback firme sem punição (PRD RNF3).

### Key Entities *(include if feature involves data)*

- **Evidence**: tipo (texto/áudio/quiz), conteúdo/descrição, validade, timestamp.
- **RubricScore**: dimensões, pontuação, observações curtas.
- **RecurringError**: descrição, categoria, contagem, tendência.

## Acceptance Criteria *(mandatory)*

- Tarefas de aprendizagem só contam como concluídas quando a evidência mínima e o gate correspondente forem satisfeitos.

## Business Objectives *(mandatory)*

- Garantir progresso real (qualidade > quantidade) e reduzir ilusões de progresso (PRD §§5.4, 10, 13).

## Error Handling *(mandatory)*

- Evidência inválida: solicitar reenvio ou alternativa equivalente.
- Tentativa de bypass: explicar motivo e oferecer versão “mínima” do gate.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Redução de tarefas “concluídas” sem evidência (tende a zero).
- **SC-002**: Redução na recorrência de erros-alvo após ciclos de reforço.

