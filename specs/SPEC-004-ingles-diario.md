# Feature Specification: Inglês Diário — Input + Output + Retrieval (rubrica + erros recorrentes)

**Created**: 2026-02-17  
**PRD Base**: §8.1, §§5.3–5.4, §9.1, §10 (R2, R3), §11 (RNF1–RNF4), §14

## Caso de uso *(mandatory)*

O usuário quer progredir em inglês com ênfase em **speaking** e **comprehensible input**, mas tende a cair em “falso progresso” (ex.: consumir conteúdo sem entender, ou “treinar speaking” sem evidência). Esta feature define um **loop diário** (Telegram-first como UX conversacional) que:

- Mantém **consistência** em dias bons e ruins (PRD §9.1; RNF2).
- Exige **evidência mínima proporcional** para contar como “feito” (PRD §5.4; R2) e reduzir “falso progresso” (PRD R3).
- Registra **qualidade** de speaking por rubrica e acompanha **erros recorrentes** como alvos de reforço (PRD §8.1).
- Mantém **baixa fricção**, com instruções objetivas e decisões mínimas (PRD §5.3; R6; RNF1).
- Dá feedback com **segurança psicológica** (firme, não punitivo) (PRD RNF3).
- Opera com **privacidade por padrão**, coletando o mínimo e explicando o que é guardado e por quê (PRD RNF4).

O loop diário é composto por três blocos, com durações de referência do PRD:

- **Input compreensível** (10–30 min) + checagem de compreensão (PRD §8.1).
- **Output guiado** (5–15 min) com evidência de produção (PRD §8.1).
- **Retrieval** (3–7 min) sem consulta para consolidar vocabulário/padrões (PRD §8.1).

## Instruções para o agente que vai escrever esta SPEC

- Use `specs/SPEC-GUIDE.md`.
- Alinhar “quality gates/evidência mínima” com a lógica definida em `SPEC-003` (sem depender de implementação).
- Não prescrever conteúdos específicos (canais/apps); foque no comportamento e critérios.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Completar inglês diário com evidência mínima (Priority: P1)

O usuário quer concluir um treino diário de inglês que **conte como progresso real**: input com checagem de compreensão, output com evidência e rubrica, e retrieval curto “sem olhar”, com caminho claro mesmo em dia ruim.

**Why this priority**: É o “MVP slice” que materializa o valor do produto para inglês (PRD §14) e implementa o núcleo de qualidade/evidência (PRD §5.4; R2, R3), com robustez (RNF2) e baixa fricção (RNF1).

**Independent Test**: Simular uma sessão diária completa (Planos A/B/C) e verificar se o sistema: (a) coleta evidências mínimas, (b) aplica gates, (c) registra rubrica e (d) registra/atualiza erros recorrentes, sem depender de revisão semanal.

**Acceptance Scenarios**:

1. **Scenario**: Plano A — completar o loop com evidência mínima
   - **Given** usuário informa tempo ≥ 30 min e energia adequada (PRD §9.1)
   - **When** inicia o inglês diário
   - **Then** o sistema orienta e coleta evidências para input + output + retrieval, e só considera “feito” quando os gates mínimos forem cumpridos (PRD §5.4; R2)

2. **Scenario**: Input — checagem de compreensão com 3 perguntas
   - **Given** usuário escolheu um conteúdo de input compreensível (PRD §8.1)
   - **When** termina o bloco de input
   - **Then** responde a **3 perguntas** de compreensão sobre o conteúdo (PRD §8.1) e o sistema registra o resultado (acertos/erros)

3. **Scenario**: Output — speaking com rubrica e autoavaliação rápida
   - **Given** usuário tem 5–15 min para output (PRD §8.1)
   - **When** produz um output (ex.: monólogo curto 1–3 min) e envia evidência de que produziu
   - **Then** o sistema registra a rubrica de autoavaliação com as dimensões do PRD (clareza/fluidez/correção aceitável/vocabulário-variedade; 0–2 por dimensão) (PRD §8.1)

4. **Scenario**: Retrieval — recall sem consulta
   - **Given** usuário completou input e/ou output no dia
   - **When** inicia o bloco de retrieval
   - **Then** o usuário faz recall de **5–10 itens** (expressões/padrões) **sem consultar** (PRD §8.1) e o sistema registra se foi concluído e a percepção de dificuldade (ex.: fácil/médio/difícil)

5. **Scenario**: Eros recorrentes — registrar e escolher 1 “erro do dia”
   - **Given** usuário terminou o output e percebeu 1–3 erros que se repetem
   - **When** registra “erros recorrentes” do dia
   - **Then** o sistema adiciona/atualiza a lista de erros recorrentes e marca pelo menos 1 como foco imediato (PRD §8.1; R3)

6. **Scenario**: Plano C — dia ruim com versão mínima que ainda gera evidência
   - **Given** usuário informa pouco tempo e energia baixa (PRD §9.1; RNF2)
   - **When** inicia o inglês diário
   - **Then** o sistema oferece um Plano C (MVD) com passos mínimos e evidência proporcional (ex.: output muito curto + rubrica rápida; ou input curto + checagem mínima), sem culpar o usuário (PRD RNF3) e mantendo a interação curta (PRD RNF1)

---

### User Story 2 - Selecionar “alvo da semana” e reforçar erros recorrentes (Priority: P2)

O usuário quer transformar erros recorrentes em **ações claras de reforço**, para reduzir repetição do mesmo padrão ao longo das semanas.

**Why this priority**: Conecta diretamente ao objetivo de reduzir “falso progresso” e reagir a padrões de erro (PRD R3; §8.1), mantendo o sistema orientado por evidência e melhoria contínua (PRD §5.4).

**Independent Test**: Com uma lista de erros recorrentes existente, executar o fluxo de seleção de “alvo da semana” e verificar que o sistema (a) define 1 alvo, (b) pede evidência de prática relacionada e (c) registra tendência do erro ao longo da semana.

**Acceptance Scenarios**:

1. **Scenario**: Escolher 1 alvo semanal baseado em recorrência
   - **Given** existe uma lista de erros recorrentes com contagem/tendência (PRD §8.1)
   - **When** o usuário inicia a semana ou solicita “definir alvo”
   - **Then** o sistema ajuda a escolher **1 erro-alvo** para a semana e explica por que (ex.: “erro repetiu X vezes”) (PRD §9.2)

2. **Scenario**: Reforço do alvo — evidência mínima
   - **Given** existe um erro-alvo da semana definido
   - **When** o usuário conclui o output diário
   - **Then** o sistema solicita um micro-reforço ligado ao alvo (ex.: 2–3 frases/mini-treino) e registra que houve tentativa (PRD R3; §8.1)

---

### User Story 3 - Ajustar dificuldade do input pelo desempenho de compreensão (Priority: P3)

O usuário quer manter input “no nível certo”: compreensível o suficiente para gerar aprendizado, mas desafiador o bastante para progredir.

**Why this priority**: Sustenta o princípio de “evidência > intuição” e melhora a qualidade do input ao longo do tempo (PRD §3; §8.1).

**Independent Test**: Rodar 3 sessões com resultados de compreensão (alta vs baixa) e verificar que o sistema recomenda ajuste de dificuldade/forma do input, sem prescrever fontes específicas.

**Acceptance Scenarios**:

1. **Scenario**: Compreensão repetidamente baixa aciona ajuste de abordagem
   - **Given** o usuário errou repetidamente a checagem de compreensão em dias recentes (PRD §8.1)
   - **When** inicia o próximo bloco de input
   - **Then** o sistema recomenda reduzir dificuldade (ex.: tema mais familiar, duração menor, mais contexto) e reforça que isso é ajuste, não falha (PRD RNF3)

### Edge Cases *(mandatory)*

- What happens when o usuário **não pode enviar áudio** por privacidade/ambiente/barulho? **[NEEDS CLARIFICATION]**: a evidência alternativa aceitável conta como “output de speaking” ou deve ser tratada como “output alternativo” (não equivalente)?
- How does system handle **input com compreensão consistentemente baixa** (ex.: < 2/3 acertos por vários dias)? Deve ajustar dificuldade, sugerir duração menor, ou pausar input e focar em vocabulário base? (PRD §8.1) **[NEEDS CLARIFICATION]**: política de “quando ajustar” (quantos dias/qual limiar).
- What happens when o usuário tenta **marcar como feito** sem cumprir gate de evidência (PRD §5.4; `SPEC-003`)?
- What happens when o usuário envia evidência **ilegível/inaudível** ou incompleta?
- What happens when o usuário tem **0 tempo** e só quer “não quebrar sequência”? (RNF2) — qual é o MVD mínimo aceitável para contar como “feito”? **[NEEDS CLARIFICATION]**.
- How does system handle quando o usuário não lembra quais foram seus “steps” do dia e pede para revisar o que fez (PRD §2: “consultar os steps do dia atual”)?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST oferecer um fluxo diário de inglês com blocos **Input**, **Output** e **Retrieval** (PRD §8.1; §14).
- **FR-002**: System MUST suportar **Planos A/B/C** para o inglês diário com base em tempo/energia informados (PRD §9.1; R1; RNF2).
- **FR-003**: System MUST definir e aplicar **quality gates** para considerar a sessão “concluída”, alinhados a `SPEC-003` (PRD §5.4; R2).
- **FR-004**: System MUST, no bloco de Input, exigir e registrar uma **checagem de compreensão** com **3 perguntas** e suas respostas (PRD §8.1).
- **FR-005**: System MUST, no bloco de Output, exigir e registrar **evidência de produção** e uma **rubrica de autoavaliação** com dimensões e escala do PRD (PRD §8.1).
- **FR-006**: System MUST, no bloco de Retrieval, guiar o usuário a fazer recall de **5–10 itens** sem consulta e registrar conclusão e percepção de dificuldade (PRD §8.1).
- **FR-007**: System MUST permitir registrar e manter uma lista de **erros recorrentes** e associá-los a contagem/tendência (PRD §8.1; R3).
- **FR-008**: System MUST permitir selecionar e registrar um **“alvo da semana”** (derivado de erros recorrentes) e vincular evidências de reforço a esse alvo (PRD §9.2; §8.1).
- **FR-009**: System MUST impedir “marcar como feito” quando os gates mínimos não foram cumpridos e oferecer o caminho mais curto para gerar evidência (PRD §5.4; R2; `SPEC-003`).
- **FR-010**: System MUST permitir ao usuário consultar um **resumo do que foi feito no dia** (steps/blocos concluídos e evidências registradas) (PRD §2; §9.1).
- **FR-011**: System MUST registrar o mínimo necessário para medir progresso e, quando aplicável, explicar “o que é guardado e por quê” (PRD RNF4). **[NEEDS CLARIFICATION]**: quais itens (ex.: áudios) são retidos e por quanto tempo.

### Non-Functional Requirements

- **NFR-001**: System MUST manter interação curta, previsível e com baixa carga cognitiva (PRD RNF1; §5.3).
- **NFR-002**: System MUST ser robusto a dias ruins e sempre oferecer um Plano C (MVD) que mantenha consistência sem burnout (PRD RNF2; §4).
- **NFR-003**: System MUST fornecer feedback firme e orientado a aprendizado, sem humilhação/punição (PRD RNF3; §3).
- **NFR-004**: System MUST operar com privacidade por padrão: coletar o mínimo e dar clareza do uso/retensão dos dados (PRD RNF4). **[NEEDS CLARIFICATION]**: política de retenção/remoção de evidências (especialmente áudio).

### Key Entities *(include if feature involves data)*

- **EnglishSession**: data; plano (A/B/C); blocos concluídos (input/output/retrieval); duração por bloco; status de “concluída”; resumo do dia (PRD §9.1).
- **ComprehensionCheck**: referência do conteúdo (descrição curta); 3 perguntas; respostas; acertos; observações de dificuldade (PRD §8.1).
- **OutputEvidence**: tipo (speaking/shadowing/conversa/output alternativo); descrição; validade; observações (PRD §8.1; §5.4). **[NEEDS CLARIFICATION]**: o que conta como “output alternativo” quando não há áudio.
- **SpeakingRubric**: scores por dimensão (clareza/fluidez/correção aceitável/vocabulário-variedade); escala 0–2; notas curtas (PRD §8.1).
- **RetrievalAttempt**: itens-alvo (5–10); conclusão; percepção de dificuldade; observações (PRD §8.1).
- **RecurringError**: descrição; categoria (opcional); exemplos; contagem; tendência; status (ativo/arquivado) (PRD §8.1; R3).
- **WeeklyTarget**: erro-alvo da semana; justificativa; evidências de reforço associadas; resultado da semana (PRD §9.2).

## Acceptance Criteria *(mandatory)*

- O sistema oferece um inglês diário com Planos A/B/C baseado em tempo/energia informados (PRD §9.1; RNF2).
- Uma sessão diária de inglês **só conta como concluída** quando os **quality gates** mínimos forem cumpridos (PRD §5.4; R2; `SPEC-003`).
- Input diário inclui checagem de compreensão com **3 perguntas** e registro de resultado (PRD §8.1).
- Output diário registra evidência de produção e rubrica (dimensões e escala do PRD) quando houver speaking (PRD §8.1).
- Retrieval diário é guiado como recall “sem olhar” de **5–10 itens** e é registrado (PRD §8.1).
- O sistema mantém lista de erros recorrentes e permite eleger 1 alvo semanal com evidências de reforço (PRD §8.1; §9.2; R3).
- Em dia ruim, existe um Plano C (MVD) com evidência proporcional, sem tom punitivo (PRD RNF2; RNF3).
- O usuário consegue consultar “o que foi feito hoje” (steps) e o estado do inglês diário no dia atual (PRD §2).

## Business Objectives *(mandatory)*

- **Consistência com qualidade**: aumentar dias/semana com treino real (não apenas “consumo”) (PRD §8.1; §9.1).
- **Reduzir falso progresso**: exigir evidência mínima e reagir a padrões de erro recorrente (PRD §5.4; R2; R3).
- **Baixa fricção**: tornar o loop executável com decisões mínimas e instruções objetivas (PRD §5.3; RNF1).
- **Segurança psicológica**: manter o usuário em movimento com feedback firme e não punitivo (PRD RNF3).
- **Privacidade por padrão**: coletar apenas o necessário e manter clareza sobre retenção/uso (PRD RNF4).

## Error Handling *(mandatory)*

- **Entrada ausente/ambígua** (tempo/energia não informados): aplicar defaults mínimos e pedir 1 confirmação curta antes de sugerir Plano A/B/C (alinhado ao catálogo do `SPEC-GUIDE`) (PRD RNF1).
- **Dia ruim / sobrecarga**: oferecer Plano C (MVD) imediatamente, reforçar que “consistência mínima” é sucesso do dia, e evitar linguagem de culpa (PRD RNF2; RNF3).
- **Evidência ausente**: não permitir conclusão; oferecer o caminho mais curto para gerar evidência mínima (PRD §5.4; R2; `SPEC-003`).
- **Evidência inválida/ilegível/inaudível**: solicitar reenvio ou alternativa definida pela SPEC; registrar tentativa sem punição (PRD RNF3). **[NEEDS CLARIFICATION]**: alternativa aceitável quando áudio não é possível.
- **Checagem de compreensão com baixo desempenho repetido**: orientar ajuste de dificuldade/duração e registrar a decisão como adaptação (não falha) (PRD §8.1; RNF3). **[NEEDS CLARIFICATION]**: limiares e janela para disparar ajuste.
- **Tentativa de bypass (“marca como feito mesmo assim”)**: explicar o porquê do gate e oferecer uma versão mínima do gate para o dia (PRD §5.4; RNF3; `SPEC-003`).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Atingir e manter uma meta de consistência semanal de inglês (ex.: dias/semana com sessão concluída por gates) (PRD §8.1; §9.1).
- **SC-002**: Aumentar o volume de prática real: minutos semanais de input registrados + número de outputs com evidência (PRD §8.1).
- **SC-003**: Melhorar a qualidade do speaking ao longo de semanas: tendência positiva na rubrica (média por dimensão) e/ou aumento de “clareza” e “fluidez” (PRD §8.1).
- **SC-004**: Reduzir recorrência do(s) erro(s)-alvo: queda na contagem de erros recorrentes do alvo semanal ao longo de ciclos (PRD §8.1; §9.2; R3).
- **SC-005**: Manter robustez a dias ruins: taxa de dias de baixa energia em que o usuário ainda conclui um Plano C (MVD) com evidência proporcional (PRD RNF2).

