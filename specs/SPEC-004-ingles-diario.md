# Feature Specification: Inglês Diário — Input + Output + Retrieval (rubrica + erros recorrentes)

**Created**: 2026-02-19  
**PRD Base**: §8.1, §5.3, §5.4, §5.2, §9.1, §14, §10 (R2, R3, R6), §11 (RNF1–RNF4)

## Caso de uso *(mandatory)*

O usuário quer melhorar inglês com foco em **speaking** e **comprehensible input**, com progresso mensurável e evitando “falso progresso”. Esta feature define um loop diário Telegram-first que:
- guia **input compreensível** (10–30 min) com checagem simples de compreensão,
- guia **output** (speaking) curto e avaliado por **rubrica**,
- inclui **retrieval** (3–7 min) para consolidar itens sem consulta,
- registra erros recorrentes e aciona reforços quando necessário,
- e aplica **quality gates**: tarefa só conta como concluída com evidência mínima adequada (`SPEC-003`).

O loop deve funcionar em dias bons e ruins (Plano A/B/C), mantendo baixa fricção, segurança psicológica e privacidade por padrão (`SPEC-015`).

## Scope & Non-goals *(recommended)*

**In scope (MVP slice)**:
- Um fluxo diário com 3 blocos (input, output, retrieval) e versões mínimas para dia ruim.
- Rubrica de output (speaking) e registro mínimo de evidência.
- Registro de erros recorrentes e alvos (para reforço / revisão semanal).
- Regras de conclusão (gate) alinhadas a `SPEC-003`.

**Non-goals (agora)**:
- Não escolher plataforma específica de conteúdo, professor, app, curso.
- Não definir STT, scoring automático, modelo de pronúncia; a rubrica pode ser autoavaliada inicialmente.
- Não desenhar pipeline técnico de áudio/armazenamento.

## Definições *(recommended)*

- **Input compreensível**: conteúdo levemente acima do nível, com compreensão global.
- **Output (speaking)**: gravação curta (ou alternativa equivalente se definida) para praticar produção.
- **Retrieval**: recall ativo de itens (palavras/expressões/padrões) sem consultar.
- **Rubrica de speaking (default PRD)**: clareza (0–2), fluidez (0–2), correção aceitável (0–2), vocabulário/variedade (0–2). Total 0–8.
- **Erro recorrente**: erro/padrão que aparece com frequência suficiente para virar alvo (limiar consistente com `SPEC-016`).

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Completar loop diário de inglês em dia normal com evidência mínima (Priority: P1)

O usuário quer fazer inglês todos os dias com um loop que realmente melhora speaking e compreensão, com prova mínima de que aconteceu.

**Why this priority**: É parte do MVP (PRD §14) e concretiza qualidade > quantidade (PRD §5.4; R2).

**Independent Test**:
- Simular um dia “bom” com tempo/energia.
- Verificar que: input tem checagem, output tem rubrica, retrieval é registrado, gate aplicado e conclusão registrada.

**Acceptance Scenarios**:

1. **Scenario**: Input concluído com checagem de compreensão
   - **Given** existe uma tarefa de inglês do dia com bloco de input
   - **When** o usuário completa o input
   - **Then** o sistema solicita 3 perguntas curtas de compreensão (ou equivalente) e registra respostas mínimas
   - **And** só conta input como concluído se a checagem mínima for satisfeita (ver gate em `SPEC-003`)

2. **Scenario**: Output (speaking) com rubrica preenchida
   - **Given** existe um bloco de output no dia
   - **When** o usuário envia speaking (ex.: gravação curta) e preenche rubrica
   - **Then** o sistema registra rubrica e marca o bloco como concluído apenas se evidência mínima estiver válida

3. **Scenario**: Retrieval curto registrado
   - **Given** o usuário completou (ou está completando) o loop do dia
   - **When** o sistema solicita retrieval de 5–10 itens
   - **Then** o usuário responde sem consulta e o sistema registra desempenho (ok / baixo) e itens-alvo

---

### User Story 2 — Dia ruim: manter consistência com versão mínima sem “passar pano” (Priority: P1)

O usuário tem pouco tempo/energia e quer um mínimo viável que preserve consistência e progresso real.

**Why this priority**: RNF2 (dias ruins) é requisito central do produto; sem isso o hábito morre.

**Independent Test**:
- Simular dia com pouco tempo/energia.
- Validar que existe versão mínima para inglês, com gate proporcional e sem aceitar conclusão vazia.

**Acceptance Scenarios**:

1. **Scenario**: Plano C de inglês em 5–15 min
   - **Given** o usuário reporta pouco tempo/energia
   - **When** pede “mínimo viável de inglês”
   - **Then** o sistema oferece uma versão curta (ex.: input 5–10 min + 1 pergunta de compreensão + 1 output de 30–60s OU alternativa mínima definida + retrieval de 3 itens)

2. **Scenario**: Usuário não pode fazer áudio (ambiente/privacidade)
   - **Given** o usuário está sem privacidade para gravar áudio
   - **When** o output do dia exige speaking
   - **Then** o sistema oferece alternativa equivalente **se** definida pela política de evidência (`SPEC-003` + `SPEC-015`)
   - **Else** registra bloqueio/pendência de forma não punitiva e oferece mínimo viável do dia que ainda seja observável

---

### User Story 3 — Registrar erros recorrentes e definir “alvo da semana” (Priority: P2)

O usuário quer que o sistema identifique padrões de erro e transforme em foco prático.

**Why this priority**: Implementa R3 (falso progresso) e alimenta revisão semanal (PRD §9.2).

**Independent Test**:
- Simular 5 dias com o mesmo erro aparecendo.
- Validar que o erro vira recorrente e fica elegível para “alvo da semana”.

**Acceptance Scenarios**:

1. **Scenario**: Erro recorrente registrado após speaking
   - **Given** o usuário fez output e percebe um erro
   - **When** registra “erro principal do dia”
   - **Then** o sistema adiciona ao log de erros e incrementa recorrência quando aplicável

2. **Scenario**: Alvo da semana sugerido
   - **Given** existem erros recorrentes ativos na semana
   - **When** chega a revisão semanal (`SPEC-007`) ou o usuário pede foco
   - **Then** o sistema sugere 1 alvo de inglês (erro/padrão) com um reforço simples para a semana

---

### User Story 4 — Transparência e privacidade ao lidar com evidências (Priority: P2)

O usuário quer sentir segurança ao enviar evidências (principalmente áudio).

**Why this priority**: RNF4 é requisito explícito e evidências podem ser sensíveis.

**Independent Test**:
- No momento de pedir speaking, perguntar “o que você guarda?” e ativar opt-out.
- Verificar explicação e modo mínimo.

**Acceptance Scenarios**:

1. **Scenario**: Explicação curta do que é guardado
   - **Given** o sistema pede evidência de speaking
   - **When** o usuário pergunta “o que você guarda?”
   - **Then** o sistema explica em 1–2 frases o mínimo guardado e oferece controle/opt-out (ver `SPEC-015`)

2. **Scenario**: Opt-out de conteúdo sensível não quebra o loop
   - **Given** o usuário ativou não guardar conteúdo sensível
   - **When** executa o loop diário
   - **Then** o sistema opera em modo mínimo e registra apenas o necessário (ex.: rubrica/agregados), conforme política

## Edge Cases *(mandatory)*

- What happens when o usuário “faz input” mas não consegue responder as perguntas?
  - Sistema não conta como concluído; sugere reduzir dificuldade do conteúdo e oferece 1 passo mínimo de compreensão.
- What happens when o usuário tenta marcar como feito sem rubrica/evidência?
  - Gate falha e o sistema pede o menor próximo passo (ver `SPEC-003`).
- What happens when o usuário repete o mesmo erro por 3+ dias?
  - Sistema marca como recorrente e aciona reforço curto (`SPEC-009`) / backlog (`SPEC-008`).
- What happens when o usuário quer fazer muito e está em overload?
  - Sistema reduz para mínimo viável e reforça foco em qualidade.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST suportar um bloco diário de **input** com checagem mínima de compreensão (default: 3 perguntas curtas).
- **FR-002**: System MUST suportar um bloco diário de **output (speaking)** com evidência mínima e **rubrica** (default PRD 0–8) registrada.
- **FR-003**: System MUST suportar **retrieval** diário (5–10 itens; versão mínima em dia ruim) e registrar desempenho de forma simples.
- **FR-004**: System MUST aplicar Quality Gates para inglês e só contar como concluído quando evidência mínima for válida (referência: `SPEC-003`).
- **FR-005**: System MUST oferecer versão mínima viável (Plano C) para dias ruins, mantendo progresso observável (PRD RNF2).
- **FR-006**: System MUST permitir registrar 1 erro/aprendizado principal do dia e consolidar erros recorrentes (alinhado a `SPEC-016`).
- **FR-007**: System MUST alimentar foco semanal (alvo de inglês) para revisão semanal (`SPEC-007`) e backlog (`SPEC-008`).
- **FR-008**: System MUST aplicar privacidade por padrão para evidências sensíveis e permitir opt-out/modo mínimo (referência: `SPEC-015`).

### Non-Functional Requirements

- **NFR-001**: System MUST manter baixa fricção e interação curta (PRD RNF1).
- **NFR-002**: System MUST ser robusto a dias ruins (PRD RNF2).
- **NFR-003**: System MUST manter segurança psicológica (PRD RNF3): feedback firme, não punitivo.
- **NFR-004**: System MUST aplicar privacidade por padrão (PRD RNF4): minimização + transparência + controle.

### Key Entities *(include if feature involves data)*

- **EnglishInputSession**: data; duração_aprox; conteúdo_descrito (alto nível); checagem_de_compreensão (respostas); status.
- **SpeakingEvidence**: data; status (válida/inválida); sensível (sim/não); política (guardar/não guardar); referência mínima.
- **SpeakingRubric**: data; clareza/fluidez/correção/vocabulário (0–2); total; status (completa/parcial).
- **EnglishRetrieval**: data; itens; desempenho (ok/baixo); itens_alvo.
- **EnglishErrorLog**: erro/padrão; contagem; último visto; status (ativo/recorrente/alvo).

## Acceptance Criteria *(mandatory)*

- O loop diário de inglês pode ser concluído com evidência mínima: input com checagem, output com rubrica, retrieval registrado.
- Em dia ruim existe Plano C executável em 5–15 min, sem aceitar conclusão vazia.
- Tarefas não contam como concluídas quando gate falha; o sistema sempre dá o menor próximo passo.
- Erros recorrentes podem ser registrados e ficam elegíveis para alvo da semana/reforço.
- Privacidade por padrão funciona: transparência + opt-out/modo mínimo.

## Business Objectives *(mandatory)*

- Melhorar speaking e compreensão com práticas eficazes (PRD §8.1) e evitar falso progresso via gates (PRD §5.4; R2/R3).
- Sustentar consistência diária com baixo atrito e versões mínimas (RNF1/RNF2).
- Manter confiança com privacidade por padrão (RNF4).

## Error Handling *(mandatory)*

- **Checagem de compreensão falhou**: reduzir dificuldade e sugerir próxima peça de input mais adequada; não contar como concluído.
- **Evidência inválida**: pedir reenvio ou alternativa equivalente quando definida.
- **Usuário some**: permitir retomar pelo próximo passo mínimo.
- **Ambiente sem privacidade**: operar em modo mínimo e registrar bloqueios explicitamente, sem culpa.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Aumento de consistência semanal no inglês (dias/semana) com gates satisfeitos.
- **SC-002**: Tendência de melhora na rubrica de speaking ao longo de semanas.
- **SC-003**: Redução de erros recorrentes ativos após alvos/reforços.
- **SC-004**: Em dias ruins, aumento da taxa de conclusão do mínimo viável sem abandono do hábito.