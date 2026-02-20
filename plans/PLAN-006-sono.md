# Technical Plan: PLAN-006 — Sono: Diário + Rotina Pré-sono + 1 Intervenção simples/semana

**Created**: 2026-02-20  
**Spec**: `specs/SPEC-006-sono.md`  
**PRD Base**: §6.1, §8.3, §5.3, §5.4, §9.1, §14, §11 (RNF1–RNF4), §10 (R2, R6)  
**Related Specs**: `SPEC-002`, `SPEC-003`, `SPEC-007`, `SPEC-011`, `SPEC-015`, `SPEC-016`

## 1) Objetivo do plano
- Implementar captura de **diário mínimo de sono** (≤ 60s) com aceitação de registros parciais.
- Implementar **rotina pré-sono** com poucos passos e versão mínima (dia ruim), registrando aderência de forma leve.
- Implementar **1 intervenção simples por semana** (experimento), com registro de adesão e fechamento semanal.
- Integrar com **quality gates de hábito** (`SPEC-003`) e com tendências/métricas (`SPEC-016`), respeitando privacidade por padrão (`SPEC-015`).

## 2) Non-goals (fora do escopo)
- Não diagnosticar/tratar condições médicas, nem aplicar CBT‑I formal.
- Não exigir wearables/dispositivos.
- Não criar análises avançadas; apenas tendências simples (regularidade, qualidade percebida, energia).
- Não implementar nudges proativos; qualquer lembrete deve seguir `SPEC-011`.

## 3) Assumptions (assunções)
- O “dia” de sono é indexado por `local_date` (timezone do usuário).
- O usuário pode registrar de manhã ou corrigir no mesmo dia; fonte de verdade é o **último registro válido** do dia (default `SPEC-016`).
- Privacidade: dados de sono são C1/C2 (não necessariamente sensíveis), mas podem revelar rotina; manter minimização e controles de apagar (`SPEC-015`).

## 4) Decisões técnicas (Decision log)
- **D-000 — Baseline de plataforma**
  - **Decisão**: adotar o baseline `plans/PLAN-000-platform-baseline.md` para stack/arquitetura/padrões (Go backend + worker/jobs; Postgres/sqlc; Redis; retenção/expiração; redaction).
  - **Motivo**: sono alimenta métricas (`SPEC-016`), revisão semanal (`SPEC-007`) e pode usar nudges (`SPEC-011`); precisa de padrões consistentes.
  - **Alternativas consideradas**: decisões isoladas por feature; descartado.
  - **Impactos/Trade-offs**: baseline central concentra decisões de execução e retenção.

- **D-001 — Gate de sono é leve e orientado a consistência**
  - **Decisão**: “cumpriu sono” = (a) diário mínimo registrado (mesmo parcial) + (b) rotina mínima executada/registrada quando aplicável; em dia ruim, permitir versão mínima.
  - **Motivo**: `SPEC-006` FR-004 + `SPEC-003` (hábito/fundação com fricção proporcional).
  - **Alternativas consideradas**: exigir dados completos todos os dias; descartado (alto atrito).
  - **Impactos/Trade-offs**: mais dados parciais; mitigado por tendências semanais (`SPEC-016`).

- **D-002 — Tendências simples primeiro, métricas derivadas depois**
  - **Decisão**: armazenar horários, qualidade 0–10 e energia 0–10 quando disponíveis; calcular regularidade e médias semanais quando houver ≥3 registros.
  - **Motivo**: alinha com `SPEC-016` (dados suficientes) e mantém baixo custo.
  - **Alternativas consideradas**: estatísticas avançadas; adiado.
  - **Impactos/Trade-offs**: precisão limitada; suficiente para decisões semanais.

- **D-003 — Intervenção semanal como entidade explícita**
  - **Decisão**: modelar `WeeklySleepIntervention` com descrição, regra de adesão e status, para fechamento semanal e histórico.
  - **Motivo**: `SPEC-006` FR-005/FR-006.
  - **Alternativas consideradas**: apenas texto na conversa; descartado (perde rastreabilidade).
  - **Impactos/Trade-offs**: precisa de “week_id” consistente.

## 5) Arquitetura (alto nível)
- **Componentes**
  - **Sleep Domain Service**: captura diário, sugere rotina (normal/mínima), sugere intervenção semanal.
  - **Gate Engine** (`SPEC-003`): gate leve para tarefas de sono (diário + mínimo).
  - **Metrics Aggregator** (`SPEC-016`): regularidade, médias, tendências.
  - **Privacy Service** (`SPEC-015`): retenção, opt-out e apagamento.

- **Fluxos**
  - **Diário (manhã)**: registrar dados (parcial ok) → gate satisfeito para diário → atualizar tendências.
  - **Rotina (noite)**: sugerir 2–4 passos (ou 1–2 no mínimo) → registrar execução (feito/parcial) → compor “cumpriu sono” do dia.
  - **Intervenção semanal**: criar no início da semana → registrar adesão durante a semana → fechar na revisão semanal (`SPEC-007`).

## 6) Contratos e interfaces
- **Comando**: `SubmitSleepDiary(user_id, local_date, slept_at?, woke_at?, quality_0_10?, morning_energy_0_10?, awakenings_note?, timestamp)`
  - **Saída**: `SleepDiaryReceipt(status complete|partial, computed_duration?, regularity_delta?, next_suggestion_short)`

- **Consulta**: `GetSleepToday(user_id, local_date)`
  - **Saída**: resumo curto do diário + status (registrado/parcial/faltando)

- **Comando**: `GetSleepNightRoutine(user_id, local_date, day_context?=normal|bad_day)`
  - **Saída**: `SleepRoutinePlan(steps[], minimum_done_criteria, version normal|minimal)`

- **Comando**: `RecordSleepRoutineResult(user_id, local_date, result done|partial|not_done, note_short?, timestamp)`
  - **Saída**: confirmação curta + sugestão de ajuste mínimo

- **Comando**: `ProposeWeeklySleepIntervention(user_id, week_id, timestamp)`
  - **Saída**: `WeeklyInterventionProposal(description, why_short, adherence_rule)`

- **Comando**: `AcceptOrRejectWeeklyIntervention(user_id, week_id, decision, timestamp)`
  - **Saída**: `WeeklyInterventionStatus`

- **Consulta**: `GetSleepWeekSummary(user_id, week_id)`
  - **Saída**: regularidade + médias (se houver) + 1 sugestão acionável

## 7) Modelo de dados (mínimo)
- **SleepDiaryEntry** (`SPEC-006`)
  - `user_id`, `local_date`
  - `slept_at`, `woke_at` (aprox), `quality_0_10`, `morning_energy_0_10`
  - `computed_duration_min?`, `regularity_delta_min?`
  - `status complete|partial`, `updated_at`
  - **Retenção** (`SPEC-016`): detalhes 90 dias; agregados semanais 12 meses.

- **SleepRoutine**
  - `user_id`, `local_date`, `version normal|minimal`
  - `steps[]` (curtos), `minimum_done_criteria`
  - `result done|partial|not_done`, `note_short?`, `updated_at`
  - **Privacidade**: `note_short` opcional/minimizada.

- **WeeklySleepIntervention**
  - `user_id`, `week_id`
  - `description`, `why_short`, `adherence_rule`
  - `status proposed|accepted|rejected`
  - `adherence_count_done`, `adherence_count_possible?`
  - `closing_outcome worked|not_worked|inconclusive`, `closed_at?`

- **SleepWeeklyAggregates** (`SPEC-016`)
  - `user_id`, `week_id`
  - `avg_quality?`, `avg_energy?`, `avg_regularity_delta?`, `days_with_diary`
  - `trend_vs_prev_week?`

## 8) Regras e defaults
- **Diário mínimo** (P1): aceitar parcial; sempre confirmar com mensagem curta.
- **Gate de hábito** (`SPEC-003` + `SPEC-006`):
  - Diário registrado (mesmo parcial) conta como evidência mínima do diário.
  - “Cumpriu sono” do dia pode exigir: diário + rotina mínima (quando rotina foi solicitada/planejada).
- **Dia ruim**: rotina mínima (1–2 passos) e registro mínimo sem culpa.
- **Tendências** (`SPEC-016`):
  - calcular médias semanais apenas com ≥3 registros no período.
  - comparação com semana anterior só se semana anterior também tiver ≥3 registros.
- **Privacidade** (`SPEC-015`):
  - permitir apagar dados de sono por período/categoria;
  - minimizar notas livres e evitar conteúdo sensível em nudges.

## 9) Observabilidade e métricas
- **Eventos**
  - `sleep_diary_submitted(partial|complete)`
  - `sleep_routine_planned(version)`, `sleep_routine_recorded(result)`
  - `sleep_weekly_intervention_proposed/accepted/rejected/closed`

- **Métricas**
  - Taxa de dias com diário registrado (parcial conta).
  - Tendência de regularidade e qualidade/energia (quando dados suficientes).
  - Aderência à intervenção semanal (done/possible).

## 10) Riscos & mitigação
- **Risco**: diário vira burocracia e o usuário para.  
  **Mitigação**: aceitar parcial; ≤60s; defaults e 1 insight curto.
- **Risco**: dados incompletos impedem tendências.  
  **Mitigação**: marcar “dados insuficientes” e sugerir 1 coleta mínima, sem culpa.
- **Risco**: semana ruim vira “bronca”.  
  **Mitigação**: linguagem protetiva; foco em mínimo e experimento simples.

## 11) Rollout / migração
- **Feature flag**: `sleep_v1`.
- Evolução: permitir “quiet hours” e nudge policy via `SPEC-011` quando habilitar proatividade.

## 12) Plano de testes (como validar)
- **Unit**
  - Validação de ranges (0–10), aceitação parcial, cálculo simples de duração/regularidade.
  - Gate leve: diário parcial satisfaz diário; rotina mínima em dia ruim.
- **Integration**
  - Correção no mesmo dia substitui registro (último válido).
  - Agregados semanais respeitam regra de ≥3 registros.
- **E2E**
  - 7 dias de diários → resumo semanal com tendência.
  - Dia ruim → rotina mínima + registro parcial ainda conta como consistência mínima.
- **Manual / acceptance**
  - Tom não punitivo; mensagens curtas; controle de apagar (`SPEC-015`).

## 13) Task breakdown (execução)
1) **Definir templates de sono (diário + rotina) com gate profile leve**
   - **Entrada**: `SPEC-006` FR-001..FR-004 + `SPEC-003`
   - **Saída**: `TaskTemplate`s e `QualityGateProfile` `sleep_diary_v1`/`sleep_routine_v1`
   - **Critério de pronto**: diário pode ser parcial; rotina tem versão mínima

2) **Modelar entidades e agregados semanais**
   - **Entrada**: `SPEC-016` (sono/energia; retenção)
   - **Saída**: schema lógico + regras de retenção (90d detalhes, 12m agregados)
   - **Critério de pronto**: cálculo de tendências não depende de dados completos

3) **Implementar submissão do diário com cálculo básico**
   - **Entrada**: `SPEC-006` User Story 1
   - **Saída**: `SubmitSleepDiary` + receipt curto + atualização de métricas do dia
   - **Critério de pronto**: ≤1 mensagem de confirmação e aceita parcial

4) **Implementar planejamento/registro de rotina pré-sono**
   - **Entrada**: `SPEC-006` User Story 2
   - **Saída**: `GetSleepNightRoutine` + `RecordSleepRoutineResult`
   - **Critério de pronto**: dia ruim gera versão mínima; falha vira dado + ajuste menor

5) **Implementar intervenção semanal (propor/aceitar/fechar)**
   - **Entrada**: `SPEC-006` FR-005/FR-006
   - **Saída**: `WeeklySleepIntervention` lifecycle + adesão simples
   - **Critério de pronto**: revisão semanal consegue ler “funcionou/não/inconclusivo”

6) **Instrumentar eventos e resumos semanais**
   - **Entrada**: `SPEC-006` SCs + `SPEC-016`
   - **Saída**: eventos `sleep_*` e `GetSleepWeekSummary`
   - **Critério de pronto**: resumo semanal retorna regularidade + médias + 1 ajuste acionável

## 14) Open questions (se existirem)
- (Default adotado) **Quando exigir rotina para “cumpriu sono”**: no MVP, só exigir rotina mínima no gate quando o usuário explicitamente solicitou “rotina da noite” naquele dia; caso contrário, “cumpriu” pode ser apenas diário registrado.

