# Technical Plan: PLAN-010 — Personalização progressiva & governança de metas em paralelo (limites de overload)

**Created**: 2026-02-20  
**Spec**: `specs/SPEC-010-governanca-metas-paralelo.md`  
**PRD Base**: §5.1, §5.3, §6.2, 10 (R7), 11 (RNF1, RNF3, RNF4), 13 (riscos)  
**Related Specs**: `SPEC-001`, `SPEC-002`, `SPEC-007`, `SPEC-011`, `SPEC-015`, `SPEC-016`, `SPEC-003`, `SPEC-009`

## 1) Objetivo do plano
- Implementar governança de metas em paralelo com **limite rígido** de 2 metas intensivas ativas por ciclo, distinguindo intensivas/fundação/aposta semanal.
- Implementar **pausar/retomar** metas intensivas com registro de contexto e plano de retomada (curta/média/longa pausa).
- Implementar detecção de **sinais de overload** e sugestão de ajuste com cooldown (no máximo 1x/semana), sem tom punitivo.
- Implementar **personalização progressiva**: começar simples e oferecer ajustes baseados em padrões observados após dados suficientes, 1 por vez e reversível.

## 2) Non-goals (fora do escopo)
- Não permitir “aumentar o limite” acima de 2 metas intensivas (é rígido).
- Não implementar otimização/IA avançada para personalização; MVP usa regras e evidência simples.
- Não enviar nudges proativos de governança no MVP sem `SPEC-011` (budgets/quiet hours/opt-out).

## 3) Assumptions (assunções)
- A classificação default de metas já é definida no onboarding (`SPEC-001`).
- Métricas necessárias para overload/personalização vêm de `SPEC-016` (consistência, energia, rubricas, gates).
- As recomendações/alertas são exibidas quando o usuário interagir (check-in, revisão semanal), e proatividade fica para `SPEC-011`.

## 4) Decisões técnicas (Decision log)
- **D-001 — Estado de “ciclo” separado do histórico**
  - **Decisão**: modelar `GoalCycle` (estado atual) e registrar mudanças via `GoalChangeEvent`.
  - **Motivo**: facilita consulta de “slots” e pausas sem reescrever passado, e alimenta métricas (`SPEC-016`).
  - **Alternativas consideradas**: apenas flags em `UserProfile`; descartado (perde histórico).
  - **Impactos/Trade-offs**: mais entidades; melhora auditabilidade.

- **D-002 — Overload como sinal composto com janela fixa**
  - **Decisão**: aplicar a janela e limiares da SPEC (últimos 7 dias com dados em ≥3 dias) para gerar `OverloadSignal`.
  - **Motivo**: consistência com a SPEC e testabilidade.
  - **Alternativas consideradas**: janelas adaptativas; adiado.
  - **Impactos/Trade-offs**: pode falhar em semanas atípicas; mitigado por “tom protetivo” e opção de ignorar.

- **D-003 — Personalização versionada, 1 por vez, reversível**
  - **Decisão**: modelar `PersonalizationSuggestion` com status (offered/accepted/reverted) e “padrão observado” como justificativa.
  - **Motivo**: R7 exige progressividade e reversão.
  - **Alternativas consideradas**: aplicar mudanças automaticamente; descartado (risco de frustração).
  - **Impactos/Trade-offs**: exige armazenamento de “before/after” para comparar impacto; usar métricas simples.

## 5) Arquitetura (alto nível)
- **Componentes**
  - **Governance Service**: enforce limite, pause/resume, consulta de slots.
  - **Overload Detector**: calcula sinal composto e aplica cooldown.
  - **Personalization Engine (MVP rules)**: identifica oportunidade após dados suficientes e cria 1 sugestão por vez.
  - **Metrics Layer** (`SPEC-016`): fornece inputs (consistência/energia/rubrica/gates).
  - **Planner Integration** (`SPEC-002`): aplica decisões/pausas ao gerar plano diário (não agenda metas pausadas).
  - **Privacy Service** (`SPEC-015`): garante minimização em registros de motivos e preferências.

## 6) Contratos e interfaces
- **Comando**: `SetActiveGoals(user_id, intensive_goals[], foundation_goals[], weekly_bets[], timestamp)`
  - **Saída**: `GoalCycleView(active, paused, intensive_count, slots_available)`
  - **Erros**: `GOAL_LIMIT_EXCEEDED` (requer escolha/pausa)

- **Comando**: `PauseGoal(user_id, goal, reason_short?, timestamp)`
  - **Saída**: `GoalStatusView(goal, status=paused, slots_available)`

- **Comando**: `ResumeGoal(user_id, goal, timestamp)`
  - **Saída**: `ResumePlanView(pause_duration_class, suggested_load_reduction, next_steps_short)`

- **Consulta**: `GetGoalGovernanceStatus(user_id)`
  - **Saída**: intensivas ativas, slots, histórico de pausas

- **Comando**: `DetectOverload(user_id, now)`
  - **Saída**: `OverloadSignalView(severity, indicators, suggestions[], cooldown_remaining?)`

- **Comando**: `RecordOverloadResponse(user_id, signal_id, choice pause|reduce|keep, timestamp)`
  - **Saída**: confirmação + efeito esperado curto

- **Comando**: `MaybeOfferPersonalization(user_id, now)`
  - **Saída**: `PersonalizationSuggestionView?` (uma por vez) com padrão observado e opção aceitar/recusar/ajustar

- **Comando**: `ActOnPersonalization(user_id, suggestion_id, action accept|reject|revert|adjust, timestamp)`
  - **Saída**: confirmação + “como vamos medir se funcionou”

## 7) Modelo de dados (mínimo)
- **GoalCycle**
  - `cycle_id`, `user_id`, `started_at`, `updated_at`
  - `intensive_active[]` (max 2), `foundation_active[]`, `weekly_bets_active[]`
  - `paused_goals[]` (goal, paused_at, reason_short?)

- **GoalChangeEvent**
  - `event_id`, `user_id`, `timestamp`
  - `type`: `activated|paused|resumed|scheduled`
  - `goal`, `reason_short?`, `pause_duration_days?`
  - **Privacidade**: `reason_short` opcional/minimizado.

- **OverloadSignal** (`SPEC-010`)
  - `signal_id`, `user_id`, `detected_at`
  - `window=7d`, `data_days_count`
  - `indicators` (consistency_low?, energy_avg?, rubric_drop?)
  - `severity`, `status open|closed`
  - `cooldown_until` (default 7 dias)

- **PersonalizationSuggestion**
  - `suggestion_id`, `user_id`, `offered_at`
  - `type` (ex.: preferred_time, reminder_intensity, task_mix)
  - `observed_pattern_short` (1 frase)
  - `proposal` (estrutura simples)
  - `status offered|accepted|rejected|reverted|adjusted`
  - `measurement_plan_short` (comparar 2 semanas antes vs 2 depois; métricas de `SPEC-016`)

## 8) Regras e defaults
- **Limite rígido**: `max_intensive_active=2` (PRD §6.2).
- **Classificação default**: intensivas (Inglês/Java), fundação (sono/saúde/autoestima), aposta semanal (SaaS) (`SPEC-001`).
- **Retomada por duração**
  - curta ≤7d: 1–2 sessões reduzidas + revisão rápida do último erro/tema
  - média 8–30d: diagnóstico leve do domínio (1–2 passos) + 1 semana de carga reduzida
  - longa >30d: “quase nova”: baseline leve do domínio + 1 semana reduzida
- **Overload (default)**: janela 7d com dados em ≥3 dias; consistência ≤2/7 em 2+ intensivas e/ou energia média ≤3/10 e/ou queda de rubrica ≥1 ponto vs semana anterior (quando aplicável).
- **Cooldown overload**: no máximo 1 sugestão por semana (7 dias).
- **Dados suficientes para personalização**: ≥14 dias desde onboarding mínimo + ≥7 check-ins + ≥5 sessões concluídas no período, e sem overload forte na última semana.
- **Uma sugestão por vez**: nunca oferecer múltiplas personalizações simultaneamente.
- **Privacidade/retensão** (`SPEC-015`): histórico de governança é C5 (preferências/estado); reter enquanto usuário usar o sistema; permitir apagar.

## 9) Observabilidade e métricas
- **Eventos**
  - `goal_limit_blocked`, `goal_paused`, `goal_resumed`
  - `overload_detected`, `overload_choice_recorded`
  - `personalization_offered/accepted/rejected/reverted`
- **Métricas**
  - % de tentativas bloqueadas corretamente (limite).
  - Taxa de pausas com motivo (meta inicial ≥80%).
  - Taxa de aceitação/reversão de personalizações e impacto (comparar 2 semanas antes/depois via `SPEC-016`).

## 10) Riscos & mitigação
- **Risco**: usuário sente “controle demais”.  
  **Mitigação**: linguagem protetiva, opções claras, aceitar “manter” sem insistir; limitar frequência de sugestões.
- **Risco**: personalização piora desempenho.  
  **Mitigação**: reversível + medição simples; oferecer reverter ao detectar piora.
- **Risco**: poucos dados impedem inferência.  
  **Mitigação**: não oferecer personalização; manter configuração simples.

## 11) Rollout / migração
- **Feature flag**: `governance_v1`.
- Migração: greenfield; compatível com onboarding existente (derivar classificação default).

## 12) Plano de testes (como validar)
- **Unit**
  - limite de 2 intensivas e mensagens/erros de domínio
  - classes de retomada por duração
  - overload detector (janela/limiares) + cooldown
  - gating de “dados suficientes” para personalização
- **Integration**
  - pausa/resume afeta geração do plano diário (`SPEC-002`)
  - personalização aceita/revertida altera preferências usadas pelo planejador
- **E2E**
  - ativar 3ª intensiva → bloqueia e oferece opções
  - semana com sinais → sugestão de reduzir/pausar registrada e refletida no ciclo
- **Manual / acceptance**
  - tom não punitivo e baixa fricção

## 13) Task breakdown (execução)
1) **Modelar `GoalCycle` + `GoalChangeEvent` e contratos de pause/resume**
   - **Entrada**: `SPEC-010` FR-001..FR-005
   - **Saída**: schema lógico + handlers
   - **Critério de pronto**: consulta mostra slots e histórico de pausas

2) **Implementar enforcement do limite de 2 intensivas**
   - **Entrada**: FR-001/FR-002
   - **Saída**: bloqueio determinístico com opções acionáveis
   - **Critério de pronto**: 100% das tentativas de 3ª intensiva são bloqueadas

3) **Implementar retomada com classes curta/média/longa**
   - **Entrada**: FR-004
   - **Saída**: `ResumePlanView` com carga reduzida e próximos passos
   - **Critério de pronto**: cada classe retorna plano coerente e curto

4) **Implementar detector de overload + cooldown**
   - **Entrada**: FR-006/FR-007
   - **Saída**: `OverloadSignal` + sugestão `pause/reduce/keep`
   - **Critério de pronto**: não sugere mais de 1x/semana; registra escolha do usuário

5) **Implementar personalização progressiva MVP**
   - **Entrada**: FR-008..FR-011
   - **Saída**: 1 sugestão por vez, com reversão e medição simples
   - **Critério de pronto**: só oferta após dados suficientes; aceitar/reverter altera preferências

6) **Instrumentar eventos/métricas de governança**
   - **Entrada**: SC-001..SC-010
   - **Saída**: eventos e contadores
   - **Critério de pronto**: possível auditar bloqueios, pausas e impacto de personalizações

## 14) Open questions (se existirem)
- (Default adotado) **Sem proatividade no MVP**: overload/personalização são apresentados no check-in e na revisão semanal; nudges proativos ficam para `SPEC-011`.

