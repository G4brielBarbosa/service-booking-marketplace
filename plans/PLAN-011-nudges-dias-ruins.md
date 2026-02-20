# Technical Plan: PLAN-011 — Nudges/Lembretes “sem spam” + Robustez a dias ruins (degradação A→B→C→MVD)

**Created**: 2026-02-20  
**Spec**: `specs/SPEC-011-nudges-dias-ruins.md`  
**PRD Base**: §5.1, §5.2, §5.3, §6.2, §9.1, §10 (R6), §11 (RNF1–RNF4), §13, §14  
**Related Specs**: `SPEC-002`, `SPEC-003`, `SPEC-015`, `SPEC-016`, `SPEC-010`, `SPEC-007`

## 1) Objetivo do plano
- Implementar **política anti-spam** para mensagens proativas (budgets, intervalo mínimo, consolidação, quiet hours).
- Implementar **timeouts** e **degradação progressiva**: check-in sem resposta (B em 90min; C/MVD em 6h) e tarefa sem início (lembrete em 4h; simplificação após +6h).
- Implementar comportamento de **ausência prolongada**: 7 dias sem interação → 1 mensagem final e depois pausar proatividade.
- Garantir privacidade por padrão em nudges (mensagens neutras; sem conteúdo sensível) e instrumentar métricas/eventos (`SPEC-016`).

## 2) Non-goals (fora do escopo)
- Não implementar notificações fora do Telegram.
- Não implementar otimização por ML para “melhor horário”; MVP usa regras default e quiet hours.
- Não tentar “compensar falhas de entrega” com spam.

## 3) Assumptions (assunções)
- Existe um `NudgePolicy` por usuário (com defaults) e um log de mensagens proativas para contabilizar budgets.
- O sistema consegue identificar “sinal de início” de uma tarefa (ex.: status `in_progress`, evidência parcial, ou interação com o task).
- Timeouts são executados por um **scheduler/worker** (cron, filas ou job runner) — sem definir stack.

## 4) Decisões técnicas (Decision log)
- **D-001 — Budgets e logs como mecanismo de qualidade**
  - **Decisão**: toda mensagem proativa que conta para budget registra `ProactiveMessageLog` com `counted_in_budget=true/false` e tipo.
  - **Motivo**: impedir regressão para spam e permitir auditoria/observabilidade.
  - **Alternativas consideradas**: apenas “best effort” sem contagem; descartado.
  - **Impactos/Trade-offs**: exige persistência e checagem antes de enviar.

- **D-002 — Consolidação first**
  - **Decisão**: quando houver múltiplas pendências elegíveis, enviar 1 mensagem consolidada com 1 escolha curta (“topa o mínimo hoje?”).
  - **Motivo**: RNF1 e FR-001/FR-004; reduz interrupções.
  - **Alternativas consideradas**: 1 mensagem por tarefa; descartado.
  - **Impactos/Trade-offs**: menos granularidade; melhora UX.

- **D-003 — Proatividade pausável e user-controlled**
  - **Decisão**: `NudgePolicy.status=active|paused` e preferência “não me lembre” aplicável (alinhado a `SPEC-015`).
  - **Motivo**: respeito ao usuário; reduz churn e reclamações de spam.
  - **Alternativas consideradas**: sempre enviar; descartado.
  - **Impactos/Trade-offs**: menos intervenção; mas preserva confiança.

## 5) Arquitetura (alto nível)
- **Componentes**
  - **Eligibility Evaluator**: decide se um nudge é elegível (budgets, quiet hours, intervalo mínimo, tarefa em progresso, ausência).
  - **Nudge Scheduler**: agenda e dispara checks de timeout (90min/6h/4h/+6h) e ausência prolongada.
  - **Nudge Composer**: gera mensagens neutras e consolidadas, com escolhas pequenas.
  - **ProactiveMessageLog Store**: registra envios, tentativas e falhas de entrega.
  - **Absence Tracker**: calcula dias sem interação e controla pausa automática.
  - **Privacy/Redaction** (`SPEC-015`): garante que nudges não vazem conteúdo sensível.

## 6) Contratos e interfaces
- **Comando**: `EvaluateDailyCheckInTimeout(user_id, local_date, now)`
  - **Saída**: `NudgeDecision(send|skip, reason_code, message?)`

- **Comando**: `EvaluateTaskNoStartTimeout(user_id, local_date, task_id, now)`
  - **Saída**: `NudgeDecision(send|skip, message?)`

- **Comando**: `EvaluateAbsence(user_id, now)`
  - **Saída**: `AbsenceDecision(nudge_final_sent?|paused?)`

- **Comando**: `SendProactiveMessage(user_id, message_type, payload, now)`
  - **Saída**: `SendReceipt(sent|skipped, budget_remaining, log_id?)`

- **Comando**: `SetNudgePolicy(user_id, quiet_hours?, intensity?, paused?, timestamp)`
  - **Saída**: `PolicyReceipt(effective_now)`

- **Consulta**: `GetNudgePolicy(user_id)`
  - **Saída**: budgets, quiet hours, status, última alteração

## 7) Modelo de dados (mínimo)
- **NudgePolicy**
  - `user_id`
  - `daily_budget_max` (default 3)
  - `per_task_budget_max` (default 2)
  - `min_interval_hours` (default 3)
  - `quiet_hours_start/end` (default 22:00–07:00; timezone do usuário)
  - `status active|paused`
  - `intensity low|medium|high` (mapeia para quão agressivo é o downgrade; MVP pode manter `medium`)
  - **Retenção**: enquanto usuário usar; parte de C5 (preferências).

- **ProactiveMessageLog**
  - `log_id`, `user_id`, `timestamp`
  - `type`: `checkin|reminder|simplify|mvd_offer|absence_final|followup`
  - `target`: `day|task_id`
  - `counted_in_budget` (bool), `budget_snapshot` (remaining)
  - `delivery_status`: `sent|failed|skipped`
  - **Privacidade**: não armazenar conteúdo sensível; `message_preview` opcional e neutro.

- **DegradationEvent**
  - `event_id`, `user_id`, `local_date`, `trigger`
  - `result_plan`: `B|C|MVD`
  - `accepted?`, `timestamp`

- **AbsenceState**
  - `user_id`
  - `last_user_interaction_at`
  - `days_without_interaction`
  - `final_message_sent_at?`
  - `nudges_paused_until_user_returns` (bool)

## 8) Regras e defaults
Defaults são os da SPEC (e devem ser configuráveis depois):

- **Budgets**
  - diário: max 3 nudges proativos/dia
  - por tarefa: max 2/dia
  - intervalo mínimo: 3h entre nudges (exceto 1 check-in do dia)
  - consolidação: 2+ pendências → 1 mensagem

- **Quiet hours**
  - default 22:00–07:00; nunca enviar dentro da janela

- **Timeouts**
  - check-in sem resposta: 90 min → Plano B leve
  - sem resposta até 6h após check-in: oferta única Plano C/MVD (se houver budget)
  - tarefa sem início: 4h → 1 lembrete curto
  - sem avanço após lembrete: +6h → versão menor ou MVD (respeitando budget)

- **Sequência de dias ruins**
  - 3 dias seguidos com execução muito baixa → sugerir começar direto pelo MVD

- **Ausência prolongada**
  - 7 dias sem interação → 1 mensagem final + pausar nudges até retorno

- **Privacidade**
  - títulos neutros (“Inglês: loop mínimo”, “Sono: diário”); nunca incluir trechos de evidência ou texto emocional por padrão (`SPEC-015`).

## 9) Observabilidade e métricas
- **Eventos**
  - `nudge_eligible_evaluated(reason_code)`
  - `nudge_sent(type)`, `nudge_skipped(reason_code)`, `nudge_failed_delivery`
  - `degradation_applied(trigger, result_plan)`
  - `nudges_paused_due_to_absence`, `nudges_resumed_on_return`

- **Métricas**
  - taxa de “budget hit” (dias em que budget impediu enviar) — deve ser baixa
  - taxa de reclamação de spam (proxy: usuário pausa/baixa intensidade)
  - % de dias com “algum passo executado” após check-in não respondido (SC-004)

## 10) Riscos & mitigação
- **Risco**: scheduler falha e dispara mensagens atrasadas (fora do contexto).  
  **Mitigação**: revalidar elegibilidade no envio; se não elegível, `skip` sem tentar compensar.
- **Risco**: mensagem proativa em momento sensível.  
  **Mitigação**: quiet hours + política “mensagem neutra”; opt-out/pause imediato.
- **Risco**: excesso de nudges em semana ruim.  
  **Mitigação**: budgets + consolidação + “pausar após 7 dias sem interação”.

## 11) Rollout / migração
- **Feature flag**: `nudges_v1`.
- Migração: ao habilitar, criar `NudgePolicy` default para usuários existentes.

## 12) Plano de testes (como validar)
- **Unit**
  - budgets e contagem correta (diário/por tarefa)
  - quiet hours e intervalo mínimo
  - timeouts 90min/6h/4h/+6h e consolidação
  - ausência 7 dias → mensagem final + pause
- **Integration**
  - end-to-end com `DailyState`: check-in enviado → sem resposta → B → sem resposta → C/MVD (respeitando budgets)
  - tarefa em progresso impede lembrete
  - logs gravados com payload neutro (sem C3)
- **E2E**
  - simular um dia com múltiplas pendências → 1 nudge consolidado
  - simular 7 dias sem interação → 1 mensagem final e depois silêncio
- **Manual / acceptance**
  - tom não punitivo; escolhas pequenas; opt-out funcionando

## 13) Task breakdown (execução)
1) **Modelar `NudgePolicy` + `ProactiveMessageLog` + `AbsenceState`**
   - **Entrada**: `SPEC-011` Key Entities + FR-001/FR-007/FR-008
   - **Saída**: schema lógico + retenção e regras de redaction
   - **Critério de pronto**: logs não guardam conteúdo sensível; policy é configurável

2) **Implementar evaluator de elegibilidade (budgets/quiet hours/intervalo)**
   - **Entrada**: `SPEC-011` FR-001/FR-002
   - **Saída**: função determinística `isEligible(nudge)` com reason codes
   - **Critério de pronto**: sempre explica por que pulou; nunca excede budgets

3) **Implementar timeouts de check-in e tarefa sem início**
   - **Entrada**: `SPEC-011` FR-003/FR-004
   - **Saída**: jobs que avaliam timeouts e disparam mensagens quando elegíveis
   - **Critério de pronto**: 90min → B; 6h → C/MVD; 4h lembrete; +6h simplificação

4) **Implementar consolidação e mensagens neutras**
   - **Entrada**: FR-001/FR-008 + política de privacidade `SPEC-015`
   - **Saída**: `Nudge Composer` com 1 mensagem para múltiplas pendências
   - **Critério de pronto**: sem trechos sensíveis; escolhas máximas 2

5) **Implementar ausência prolongada e pausa automática**
   - **Entrada**: FR-007
   - **Saída**: job `EvaluateAbsence` + transição `paused_until_return`
   - **Critério de pronto**: após 7 dias, envia 1 vez e para; retoma só após interação

6) **Instrumentar eventos e métricas**
   - **Entrada**: `SPEC-011` SCs + `SPEC-016`
   - **Saída**: eventos `nudge_*` e agregados de spam/recuperação
   - **Critério de pronto**: possível medir budget hits e recuperação após degradação

## 14) Open questions (se existirem)
- (Default adotado) **“Check-in do dia” conta no budget**: conta sim (para manter limite real). Caso você queira que não conte, isso deve ser configurável via `NudgePolicy`.

