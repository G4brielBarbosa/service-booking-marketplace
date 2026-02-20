# Feature Specification: Nudges/Lembretes “sem spam” + Robustez a dias ruins (degradação A→B→C→MVD)

**Created**: 2026-02-19  
**PRD Base**: §5.1, §5.2, §5.3, §6.2, §9.1, §10 (R6), §11 (RNF1–RNF4), §13, §14  
**Related**: `SPEC-002` (rotina diária/planos), `SPEC-003` (gates/evidência), `SPEC-015` (privacidade), `SPEC-016` (métricas)

## Caso de uso *(mandatory)*

O usuário inevitavelmente tem dias ruins (pouco tempo/energia, estresse, ambiente ruim) e pode sumir/atrasar respostas. O sistema precisa manter o usuário “em movimento” com **nudges úteis**, sem virar spam, e com **degradação graciosa**:

- se não responde ao check-in → oferecer alternativa leve automaticamente;
- se não executa tarefas → propor versão mais simples (sem culpa) e preservar consistência mínima;
- se há sequência de falhas/sumiço → reduzir expectativa, oferecer MVD e eventualmente pausar proatividade.

Esta SPEC define:
- quando enviar lembretes (e quando **não** enviar),
- limites de frequência (“budget de nudges”),
- escalonamento A→B→C→MVD,
- comportamento em ausência prolongada,
- e requisitos de tom/segurança psicológica/privacidade.

> Importante: esta SPEC não define como planos são gerados nem conteúdo das tarefas (isso é `SPEC-002`). Aqui é **proatividade, anti-spam e degradação**.

## Scope & Non-goals *(recommended)*

**In scope (MVP slice)**:
- Política anti-spam com limites claros (por tarefa e por dia).
- Lembretes consolidados (uma mensagem para várias pendências).
- Timeouts padrão para “não respondeu” e “não começou”.
- Degradação automática para Plano B/C e MVD.
- Pausa automática de nudges após ausência prolongada.

**Non-goals (agora)**:
- Não otimizar com ML; regras simples e observáveis bastam.
- Não implementar sistema de notificações avançado fora do Telegram.

## Definições *(recommended)*

- **Proativo**: mensagem enviada sem o usuário iniciar a conversa.
- **Nudge budget (dia)**: máximo de mensagens proativas por dia para evitar spam.
- **Quiet hours**: janela de sono/descanso em que o sistema não envia nudges (configurável).
- **Degradação**: reduzir exigência mantendo identidade do hábito (A→B→C→MVD).
- **MVD**: mínimo viável diário do ciclo/metas ativas (definido no onboarding e ajustável), tipicamente 5–15 min total.

## Políticas padrão (defaults) *(recommended)*

> O usuário pode ajustar depois, mas o MVP precisa de defaults.

### Anti-spam (budgets)
- **Budget diário**: no máximo **3** mensagens proativas/dia.
- **Por tarefa**: no máximo **2** lembretes proativos/dia por tarefa.
- **Intervalo mínimo**: pelo menos **3 horas** entre lembretes proativos (mesmo que sejam tarefas diferentes), exceto 1 “check-in do dia”.
- **Consolidação**: se houver 2+ pendências, enviar **1** mensagem consolidada em vez de várias.

### Quiet hours
- Respeitar quiet hours do usuário.
- **Default** (se ainda não configurado): **22:00–07:00** no fuso do usuário (até ser ajustado no onboarding).

### Timeouts
- **Check-in sem resposta**: se não houver resposta em **90 min**, oferecer automaticamente um **Plano B leve** (“assumi dia corrido”).
- Se continuar sem resposta até **6 horas** após o check-in, oferecer **Plano C/MVD** (uma única mensagem, se ainda houver budget do dia).
- **Tarefa sem início**: se uma tarefa foi sugerida e não há sinal de início em **4 horas**, enviar 1 lembrete curto (respeitando budget/quiet hours).
- **Tarefa ainda pendente**: se após o 1º lembrete passar **6 horas** sem avanço, oferecer “versão menor” ou MVD (sem insistir além disso no mesmo dia).

### Ausência prolongada
- Se não houver qualquer interação do usuário por **7 dias**, o sistema envia **1** mensagem final de “check de bem-estar + como retomar” e depois **pausa nudges proativos** até o usuário retornar.

### Privacidade por padrão em nudges
- Nudges não devem incluir conteúdo sensível (ex.: transcrições/trechos emocionais) por padrão; apenas títulos neutros (“Inglês: loop mínimo”, “Sono: diário”).
- Opt-outs e retenção seguem `SPEC-015`.

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Check-in não respondido → oferta automática de Plano B/C (Priority: P1)

**Why this priority**: Robustez a dias ruins (RNF2) e baixa carga cognitiva (RNF1).

**Independent Test**: Simular check-in enviado, sem resposta; validar oferta automática de B e depois C/MVD, respeitando quiet hours e budget.

**Acceptance Scenarios**:

1. **Scenario**: Timeout de check-in oferece Plano B
   - **Given** o sistema enviou o check-in do dia
   - **When** passam 90 minutos sem resposta
   - **Then** o sistema envia uma mensagem curta oferecendo um Plano B leve (com 1 prioridade + 1 fundação mínima), com opção “ok” para aceitar ou enviar check-in para personalizar

2. **Scenario**: Ausência continua → oferta única de Plano C/MVD
   - **Given** o usuário não respondeu ao Plano B
   - **When** passam 6 horas desde o check-in e ainda há budget do dia
   - **Then** o sistema oferece Plano C/MVD em 1 mensagem (“vamos manter o essencial hoje?”) sem tom punitivo

3. **Scenario**: Usuário responde depois → priorizar check-in real
   - **Given** o sistema já ofereceu Plano B/C
   - **When** o usuário envia check-in real
   - **Then** o sistema ajusta para o contexto real e não “penaliza” o atraso

---

### User Story 2 — Tarefa não executada → alternativa progressivamente mais simples (Priority: P1)

**Why this priority**: Evita que uma tarefa trave o dia inteiro; mantém identidade do hábito (RNF2).

**Independent Test**: Simular tarefa pendente sem início; validar lembrete + oferta de simplificação + MVD, sem exceder budgets.

**Acceptance Scenarios**:

1. **Scenario**: Lembrete curto após ausência de início
   - **Given** existe uma tarefa planejada hoje
   - **When** passam 4 horas sem qualquer sinal de início
   - **Then** o sistema envia 1 lembrete curto: “Quer começar agora ou prefere reduzir?”

2. **Scenario**: Sem avanço após lembrete → oferecer versão menor
   - **Given** o lembrete foi enviado
   - **When** passam 6 horas sem avanço e ainda há budget do dia
   - **Then** o sistema oferece uma versão menor (2 opções no máximo) ou MVD se o contexto indicar dia ruim

3. **Scenario**: Múltiplas tarefas pendentes → consolidação + foco no mínimo
   - **Given** 2+ tarefas pendentes no dia
   - **When** o sistema decide lembrar
   - **Then** envia 1 mensagem consolidada propondo “mínimo essencial” (MVD) e perguntando por 1 escolha simples (“topa o mínimo hoje?”)

---

### User Story 3 — Lembretes sem spam (budgets + quiet hours + não interromper execução) (Priority: P1)

**Why this priority**: Spam destrói o produto; lembretes precisam ser raros e úteis (RNF1/RNF3).

**Independent Test**: Simular dia com múltiplas pendências, com quiet hours, e com tarefa em progresso.

**Acceptance Scenarios**:

1. **Scenario**: Não enviar durante quiet hours
   - **Given** é horário dentro de quiet hours
   - **When** um lembrete estaria elegível
   - **Then** o sistema não envia e reavalia após o fim da janela

2. **Scenario**: Não enviar se tarefa está em progresso
   - **Given** a tarefa está marcada como “em progresso” ou há evidência parcial
   - **When** um lembrete seria elegível
   - **Then** o sistema não envia

3. **Scenario**: Budget diário impede insistência
   - **Given** o sistema já enviou 3 nudges proativos no dia
   - **When** novas pendências surgem
   - **Then** o sistema não envia mais proativos naquele dia e deixa para a próxima interação do usuário

---

### User Story 4 — Sequência de dias ruins → MVD e redução de expectativa (Priority: P1)

**Why this priority**: Previne “uma semana perdida” e protege segurança psicológica (RNF2/RNF3).

**Independent Test**: Simular 3 dias seguidos com baixa energia / baixa execução; validar redução de escopo e oferta de MVD sem culpa.

**Acceptance Scenarios**:

1. **Scenario**: 3 dias seguidos de baixa execução → foco em MVD
   - **Given** houve 3 dias recentes com execução muito baixa (ou apenas MVD)
   - **When** chega um novo dia
   - **Then** o sistema sugere começar direto pelo MVD e evita planos ambiciosos até a tendência melhorar

2. **Scenario**: Execução parcial é celebrada e não vira bronca
   - **Given** o usuário fez parte do dia (ex.: só input, sem speaking)
   - **When** o sistema sumariza o dia
   - **Then** reconhece a execução parcial, oferece um próximo passo mínimo opcional e não registra como “falha completa”

---

### User Story 5 — Ausência prolongada → pausa proativa e retomada simples (Priority: P2)

**Why this priority**: Respeito e anti-spam; o usuário não deve ser perseguido.

**Independent Test**: Simular 7 dias sem mensagens do usuário; validar 1 última mensagem e depois silêncio até retorno.

**Acceptance Scenarios**:

1. **Scenario**: Pausar nudges após 7 dias sem interação
   - **Given** o usuário está ausente há 7 dias
   - **When** o sistema alcança o limiar de ausência prolongada
   - **Then** envia 1 mensagem final curta (“quando voltar, eu te ajudo a retomar com o mínimo”) e pausa proatividade até o usuário responder

## Edge Cases *(mandatory)*

- **Energia/tempo inconsistentes** (ex.: energia 2, tempo 60): priorizar energia e oferecer mínimo; opcionalmente pedir 1 confirmação curta.
- **Usuário diz “fiz” sem evidência em tarefas com gate**: tratar via `SPEC-003` (não concluir), e aqui apenas oferecer o menor passo para evidenciar (sem insistir repetidamente no dia).
- **Usuário pede “não me lembre”**: reduzir intensidade (modo silencioso) e registrar preferência; nudges passam a ser apenas check-in (ou nenhum), conforme escolha do usuário (alinhado a RNF1/RNF3 e `SPEC-015`).

## Requirements *(mandatory)*

### Functional Requirements
- **FR-001**: System MUST ter budgets anti-spam (diário e por tarefa) e consolidar lembretes quando houver múltiplas pendências.
- **FR-002**: System MUST respeitar quiet hours (configurável) e não enviar nudges durante essa janela.
- **FR-003**: System MUST detectar “check-in sem resposta” e oferecer automaticamente Plano B após 90 min; e Plano C/MVD após 6h (respeitando budget).
- **FR-004**: System MUST detectar “tarefa sem início” e enviar 1 lembrete após 4h; e oferecer simplificação/MVD após mais 6h sem avanço (respeitando budgets).
- **FR-005**: System MUST oferecer degradação progressiva A→B→C→MVD sem exigir decisões complexas (no máximo 1 pergunta curta por vez).
- **FR-006**: System MUST reduzir expectativa e sugerir MVD quando houver sequência recente de dias ruins (default: 3 dias).
- **FR-007**: System MUST pausar nudges proativos após ausência prolongada (default: 7 dias) até retorno do usuário.
- **FR-008**: System MUST aplicar privacidade por padrão nos nudges (mensagens neutras; sem conteúdo sensível por padrão), alinhado a `SPEC-015`.

### Non-Functional Requirements
- **NFR-001**: System MUST manter simplicidade: mensagens curtas, escolhas pequenas, sem jargão (PRD RNF1).
- **NFR-002**: System MUST ser robusto a dias ruins: sempre existir um mínimo viável (PRD RNF2).
- **NFR-003**: System MUST manter segurança psicológica: tom não punitivo, foco em processo e pequenas vitórias (PRD RNF3).
- **NFR-004**: System MUST evitar spam como requisito de qualidade (implícito em RNF1) com budgets e consolidação.

### Key Entities *(include if feature involves data)*
- **NudgePolicy**: budgets; quiet hours; intensidade (baixo/médio/alto); status (ativo/pausado).
- **ProactiveMessageLog**: data/hora; tipo (check-in/lembrete/ajuste/MVD/pausa); alvo (tarefa/dia); contado_no_budget (sim/não).
- **DegradationEvent**: data; gatilho (timeout/sem início/seq dias ruins); resultado (Plano B/C/MVD); aceito?; nota curta.
- **AbsenceState**: dias_sem_interação; pausado?; última mensagem final enviada?

## Acceptance Criteria *(mandatory)*
- O sistema nunca excede o budget diário de nudges proativos e consolida mensagens quando há múltiplas pendências.
- O sistema respeita quiet hours e não interrompe execução em progresso.
- Check-in sem resposta degrada automaticamente para Plano B e depois C/MVD com defaults claros.
- Tarefas pendentes recebem no máximo 2 lembretes/dia e depois viram oferta de simplificação/MVD.
- Após 7 dias sem interação, o sistema pausa proatividade até retorno do usuário.
- Tom é não punitivo e reforça consistência mínima como sucesso.

## Business Objectives *(mandatory)*
- Proteger consistência em dias ruins sem burnout (PRD §13; RNF2).
- Reduzir carga cognitiva: o sistema toma a iniciativa com opções simples (RNF1).
- Evitar spam e sustentar confiança/aderência (RNF1/RNF3).
- Manter privacidade por padrão em mensagens e controles (`SPEC-015`).

## Error Handling *(mandatory)*
- **Configurações inválidas** (quiet hours incoerentes / budgets extremos): aplicar defaults seguros e pedir ajuste com 1 pergunta curta.
- **Contestação do usuário** (“eu fiz”): aceitar como relato, mas se for tarefa com gate, pedir evidência mínima conforme `SPEC-003` (sem insistência repetida).
- **Falha de entrega de mensagem**: registrar tentativa e não “compensar” com spam; reavaliar no próximo ciclo.

## Success Criteria *(mandatory)*
### Measurable Outcomes
- **SC-001**: Aumento de dias com execução mínima (Plano C/MVD) em semanas de baixa energia.
- **SC-002**: Redução de abandono após sequências de dias ruins (ex.: após 3 dias).
- **SC-003**: Baixa taxa de reclamação de “spam” (proxy: usuário reduz intensidade/pausa nudges com baixa frequência).
- **SC-004**: Boa recuperação: dias com check-in não respondido ainda resultam em algum passo executado (Plano B/C/MVD).