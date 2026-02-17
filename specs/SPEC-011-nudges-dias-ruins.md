# Feature Specification: Nudges/Lembretes “sem spam” + robustez a dias ruins (Planos B/C + MVD)

**Created**: [DATE]  
**PRD Base**: §5.3, 11 (RNF1, RNF2, RNF3)

## Caso de uso *(mandatory)*

Definir como o sistema envia lembretes e nudges que ajudam a executar sem gerar spam, e como garante que sempre exista uma opção mínima (Plano C / MVD) em dias ruins.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST permitir nudges condicionais (ex.: ausência de check-in, janela de horário, energia baixa).
- **FR-002**: System MUST oferecer opção de “silenciar”/reduzir lembretes sem quebrar o plano.
- **FR-003**: System MUST sempre ter um caminho de execução mínima (MVD) quando aplicável.

### Non-Functional Requirements

- **NFR-001**: System MUST evitar spam e manter previsibilidade da interação (PRD RNF1).
- **NFR-002**: System MUST ser robusto a dias ruins (PRD RNF2).
- **NFR-003**: System MUST manter feedback não punitivo (PRD RNF3).

## Error Handling *(mandatory)*

- Ausência prolongada: retomar com mensagem curta, sem culpa, e oferecer plano mínimo.

