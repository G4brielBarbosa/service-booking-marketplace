# Feature Specification: Java Diário — Prática Deliberada + Retrieval + Revisão de Erros

**Created**: 2026-02-19  
**PRD Base**: §8.2, §5.4, §5.3, §9.1, §14, §10 (R2, R3, R6), §11 (RNF1–RNF3)

## Caso de uso *(mandatory)*

O usuário quer evoluir em Java com progresso **mensurável** e sustentável, evitando “fiz mas não aprendi”. Esta feature define um **loop diário** Telegram-first que:
- guia uma sessão curta de prática deliberada,
- checa entendimento com **retrieval** (sem consulta),
- registra **o principal aprendizado/erro**,
- e só conta como “concluído” quando o **quality gate** da tarefa for satisfeito (ver `SPEC-003`).

A feature deve funcionar em dias bons e ruins, preservando consistência com baixa fricção e mantendo tom firme, não punitivo.

## Scope & Non-goals *(recommended)*

**In scope (MVP slice)**:
- Um fluxo diário para Java com 3 componentes: prática → retrieval → registro de erro/aprendizado.
- Evidência mínima e regras de “concluído vs tentativa” alinhadas a `SPEC-003`.
- Registro de erros recorrentes (para reforço futuro).

**Non-goals (agora)**:
- Não definir editor/IDE, repositório, linguagem de testes, plataformas de exercício, nem detalhes técnicos de armazenamento.
- Não definir backlog inteligente/seleção automática ótima de tarefas (isso é `SPEC-008`/`SPEC-009`).
- Não definir revisão semanal (isso é `SPEC-007`/`SPEC-016`).

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Completar sessão de Java em dia normal com evidência mínima (Priority: P1)

O usuário quer executar uma sessão diária de Java (20–60 min total) que gere evidência real de prática e aprendizado, sem burocracia.

**Why this priority**: É parte central do MVP (PRD §14) para progresso mensurável em Java (PRD §8.2) e para “qualidade > quantidade” via gates (PRD §5.4; R2).

**Independent Test**:
- Simular um dia com tempo/energia suficientes.
- Rodar o fluxo completo e verificar que existe: (a) prática registrada, (b) retrieval registrado, (c) 1 aprendizado/erro registrado, (d) gate aplicado e resultado (aceito/rejeitado) conforme `SPEC-003`.

**Acceptance Scenarios**:

1. **Scenario**: Sessão completa conta como concluída com gate satisfeito
   - **Given** existe uma tarefa de Java no plano do dia com objetivo/ restrição clara
   - **When** o usuário executa a prática, faz o retrieval e registra aprendizado/erro com evidência mínima
   - **Then** o sistema avalia o quality gate (ver `SPEC-003`) e marca a tarefa como concluída apenas se a evidência mínima for válida e completa

2. **Scenario**: Gate falha por falta de evidência e sistema orienta o menor próximo passo
   - **Given** existe uma tarefa de Java pendente com gate definido
   - **When** o usuário tenta concluir sem fornecer a evidência mínima (ou fornece evidência incompleta)
   - **Then** o sistema não conta como concluída, explica em 1–2 frases e solicita o menor passo adicional para satisfazer o gate

3. **Scenario**: Retrieval revela baixa compreensão e aciona reforço mínimo observável
   - **Given** o usuário terminou a prática
   - **When** o usuário falha no retrieval (ex.: não consegue explicar o conceito ou erra a maioria dos itens)
   - **Then** o sistema registra “retrieval falhou” e solicita um reforço mínimo observável (ex.: reexplicar em 2–4 frases ou responder 1–2 itens adicionais), sem virar burocracia

---

### User Story 2 — Dia ruim: versão mínima viável sem “passar pano” para aprendizado (Priority: P1)

O usuário está com pouco tempo/energia e quer manter consistência sem se sentir culpado, mas sem o sistema aceitar “conclusão falsa”.

**Why this priority**: Robustez a dias ruins é princípio central (PRD RNF2) e a qualidade não pode ser sacrificada em aprendizagem (PRD §5.4).

**Independent Test**:
- Simular check-in de pouco tempo/energia.
- Verificar que o sistema oferece uma sessão mínima (Plano C) com evidência mínima e aplica gate; se o usuário não entrega evidência, não conta como concluído.

**Acceptance Scenarios**:

1. **Scenario**: Plano C de Java cabe em 5–15 minutos e gera evidência mínima
   - **Given** o usuário reporta pouco tempo e energia baixa
   - **When** solicita o mínimo viável para Java hoje
   - **Then** o sistema oferece uma versão curta (ex.: 1 micro-exercício + 1 retrieval + 1 erro/aprendizado) com critérios claros de “feito”

2. **Scenario**: Usuário faz parte mas não completa evidência — registra tentativa sem contar como concluído
   - **Given** o usuário iniciou a sessão mínima
   - **When** ele para antes de entregar a evidência mínima do gate
   - **Then** o sistema registra como “tentativa” (não concluída), preserva o próximo passo mais curto para retomar, e mantém tom não punitivo

> Regra de “tentativa vs concluído”: deve seguir a política definida em `SPEC-003` (ou, se ainda não estiver definida lá, esta SPEC deve definir um default e `SPEC-003` deve referenciar).

---

### User Story 3 — Registrar erro/aprendizado e consolidar erro recorrente (Priority: P2)

O usuário quer que o sistema ajude a transformar erros em reforço ao longo do tempo, sem exigir textos longos.

**Why this priority**: Suporta detecção de falso progresso (PRD R3) e melhora mensurável (PRD §8.2).

**Independent Test**:
- Simular 3 sessões com o mesmo erro categorizado.
- Verificar que o erro vira “recorrente” e fica disponível para reforço futuro (sem depender de backlog inteligente).

**Acceptance Scenarios**:

1. **Scenario**: Registro curto do principal erro/aprendizado após a sessão
   - **Given** o usuário concluiu a prática (com ou sem sucesso)
   - **When** o sistema pede o registro final
   - **Then** o usuário registra 1 item curto: “erro/aprendizado principal” e “correção/nota”, e o sistema salva associado ao dia

2. **Scenario**: Erro se torna recorrente após repetição
   - **Given** o mesmo erro foi registrado em múltiplas sessões recentes
   - **When** o sistema atualiza o log
   - **Then** o erro é marcado como recorrente e fica pronto para ser usado por reforços (`SPEC-003`) e priorização (`SPEC-008`)

### Edge Cases *(mandatory)*

- What happens when o usuário não consegue terminar o exercício no tempo?
  - Sistema registra “parcial/tentativa”, sugere menor próximo passo para retomar, e não conta como concluído sem gate.
- What happens when o usuário fornece evidência “vazia” (ex.: “fiz”)?
  - Sistema rejeita conclusão e pede evidência mínima concreta, com tom firme e curto.
- What happens when o usuário repete o mesmo erro por 3+ sessões?
  - Sistema marca como recorrente e sugere reforço mínimo observável na próxima oportunidade (sem virar obrigação longa).
- What happens when o usuário não pode produzir a evidência “ideal” (ambiente/privacidade)?
  - Sistema oferece alternativa equivalente **apenas se** a política de equivalência estiver definida em `SPEC-003`/`SPEC-015`; caso contrário registra bloqueio/pendência e sugere o mínimo viável.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST suportar uma sessão diária de Java composta por: (a) prática deliberada, (b) retrieval sem consulta, (c) registro de erro/aprendizado.
- **FR-002**: System MUST definir para cada sessão um objetivo/ restrição clara (em linguagem natural) e uma definição observável de “feito”.
- **FR-003**: System MUST aplicar Quality Gates para tarefas de Java e só marcar “concluído” quando a evidência mínima for válida (referência: `SPEC-003`).
- **FR-004**: System MUST, quando o gate falhar, explicar “por que falhou” e fornecer o menor próximo passo para recuperar (curto e acionável).
- **FR-005**: System MUST registrar resultado do retrieval e tratá-lo como evidência de entendimento (apto a influenciar reforço/ajuste futuro).
- **FR-006**: System MUST capturar ao menos 1 erro/aprendizado principal por sessão (curto) e associar ao dia.
- **FR-007**: System MUST identificar e marcar “erro recorrente” quando o mesmo erro reaparece ao longo do tempo, segundo limiar definido (default sugerido: ≥3 ocorrências em janela recente; a janela/limiar deve ser consistente com `SPEC-016`).

### Non-Functional Requirements

- **NFR-001**: System MUST manter interação objetiva e de baixa fricção (PRD RNF1).
- **NFR-002**: System MUST ser robusto a dias ruins, oferecendo versão mínima viável (PRD RNF2).
- **NFR-003**: System MUST manter segurança psicológica: feedback firme, não punitivo, focado em processo e ajuste (PRD RNF3).

### Key Entities *(include if feature involves data)*

- **JavaSession**: data; objetivo/ restrição; duração planejada; status (pendente/em progresso/concluída/tentativa/bloqueada).
- **JavaPracticeEvidence**: descrição curta do que foi produzido/feito; validade (válida/inválida/parcial).
- **RetrievalResult**: itens/perguntas; acertos/erros; status (ok/baixo); observações curtas.
- **LearningLogEntry**: “erro/aprendizado principal”; correção/nota; categoria.
- **RecurringError**: categoria/descrição; contagem; tendência; último visto; status (ativo/alvo/arquivado).

## Acceptance Criteria *(mandatory)*

- A tarefa diária de Java **não** conta como concluída sem satisfazer o gate definido (referência: `SPEC-003`).
- Em dia ruim, existe uma versão mínima viável (5–15 min) que preserva consistência **sem** aceitar conclusão falsa.
- Quando o gate falha, o sistema sempre retorna um motivo curto + um único próximo passo mínimo.
- Cada sessão gera ao menos 1 registro curto de erro/aprendizado; erros repetidos se tornam recorrentes e ficam consultáveis.

## Business Objectives *(mandatory)*

- Progresso mensurável em Java com evidência e checagens curtas, reduzindo “falso progresso” (PRD §8.2; §5.4; R2/R3).
- Baixa fricção e sustentabilidade diária, inclusive em dias ruins (PRD RNF1/RNF2).
- Feedback firme e não punitivo para manter consistência (PRD RNF3).

## Error Handling *(mandatory)*

- **Exercício incompleto**: registrar tentativa/parcial e sugerir o menor próximo passo; não contar como concluído sem gate.
- **Evidência insuficiente/inválida**: bloquear conclusão e pedir reenvio/ complemento mínimo.
- **Usuário some**: ao retornar, mostrar estado e o próximo passo mais curto para finalizar o gate ou registrar tentativa.
- **Overload/tempo caiu**: oferecer downgrade para versão mínima viável, mantendo critérios claros de “feito”.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Aumento sustentado de dias/semana com sessão de Java concluída com gate satisfeito.
- **SC-002**: Redução na frequência de erros recorrentes ativos ao longo de semanas.
- **SC-003**: Alta taxa de sessões com registro de aprendizado/erro (sem exigir textos longos), mantendo fricção baixa.