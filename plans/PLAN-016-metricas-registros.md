# Technical Plan: PLAN-016 — Métricas & Registros: consistência, rubricas, tendências, erros recorrentes

**Created**: 2026-02-20  
**Spec**: `specs/SPEC-016-metricas-e-registros.md`  
**PRD Base**: §5.4, §9.2, §§8.1–8.5, §10 (R2, R5, R6), §11 (RNF1, RNF4)  
**Related Specs**: `SPEC-002`, `SPEC-003`, `SPEC-007`, `SPEC-015`, `SPEC-011`, `SPEC-004`, `SPEC-005`, `SPEC-006`, `SPEC-012`, `SPEC-013`, `SPEC-014`

## 1) Objetivo do plano
- Implementar o modelo mínimo de **registros diários** e **agregados semanais** para suportar: consulta do dia (PRD §2), quality gates, revisão semanal e detecção de sinais.
- Implementar cálculo de **consistência** (dias/semana), **rubricas** (média/tendência), **erros recorrentes** (≥3/14d) e **sono/energia** (regularidade e médias).
- Implementar políticas de **retenção/performance** (targets de resposta) e integração com privacidade (`SPEC-015`) e nudges (`SPEC-011`) para evitar spam na recuperação de gaps.

## 2) Non-goals (fora do escopo)
- Não construir dashboards gráficos/BI.
- Não armazenar conteúdo sensível bruto em agregados (C4 não contém C3).
- Não implementar análises estatísticas avançadas; apenas médias/comparações e contagens.

## 3) Assumptions (assunções)
- Todo domínio registra eventos mínimos (ex.: gate satisf/failed, rubrica totals, minutes, etc.).
- O sistema tem timezone do usuário para definir `local_date` e `week_id`.
- Quando dados faltarem, o sistema não penaliza; marca lacunas e sugere coleta mínima (sem nudges proativos fora de `SPEC-011`).

## 4) Decisões técnicas (Decision log)
- **D-001 — Separar “detalhes diários” de “agregados semanais”**
  - **Decisão**: manter tabelas/coleções de registros diários (C1/C2) com retenção moderada, e agregados semanais (C4) com retenção longa.
  - **Motivo**: performance e privacidade; suporta comparações semanais mesmo quando detalhes expiram.
  - **Alternativas consideradas**: calcular tudo on-the-fly; descartado (latência).
  - **Impactos/Trade-offs**: exige job de agregação incremental.

- **D-002 — Consistência depende de gate satisfeito (aprendizagem)**
  - **Decisão**: contar um dia para meta intensiva somente quando `GateResult=satisfied` (FR-008).
  - **Motivo**: qualidade > quantidade; evita falso progresso.
  - **Alternativas consideradas**: contar tentativa; descartado para intensivas (pode contar como micro-consistência separada).

- **D-003 — Normalização de rubricas por domínio**
  - **Decisão**: armazenar rubricas em estrutura flexível (`dimensions` + `total` + `max_total`) para permitir domínios diferentes, mantendo comparação por domínio.
  - **Motivo**: inglês tem 0–8; outros podem variar.
  - **Alternativas consideradas**: fixar schema rígido; descartado.

## 5) Arquitetura (alto nível)
- **Componentes**
  - **Event Ingestor**: recebe eventos de domínio (tarefas, gates, rubricas, sono, saúde, autoestima, SaaS).
  - **Daily Snapshot Builder**: constrói `DailyProgressSnapshot` para consulta rápida do dia.
  - **Weekly Aggregator**: calcula agregados semana atual e anterior (quando dados suficientes).
  - **Recurring Error Tracker**: atualiza contagens por janela rolante (14 dias) e status (ativo/recorrente/alvo).
  - **Query API**: expõe consultas rápidas (dia e semana) com targets de latência.
  - **Privacy/Retention** (`SPEC-015`): aplica expiração e garante que agregados não contém C3.

## 6) Contratos e interfaces
- **Consulta**: `GetTodayProgress(user_id, local_date)`
  - **Saída**: `DailyProgressView` (≤1 mensagem): por meta/domínio, status (done/in_progress/pending), rubrica média do dia quando existir, e link para plano do dia.

- **Consulta**: `GetWeeklyMetrics(user_id, week_id)`
  - **Saída**: `WeeklyMetricsView` (consistência por meta, rubricas, erros recorrentes, sono/energia, saúde, autoestima, SaaS)

- **Comando**: `IngestDomainEvent(user_id, event_type, payload_min, timestamp, sensitivity=C1|C2|C4)`
  - **Saída**: `IngestReceipt`

- **Comando**: `RecomputeWeeklyAggregates(user_id, week_id)` (admin/job)
  - **Saída**: `AggregateReceipt`

## 7) Modelo de dados (mínimo)
- **DomainEventLog** (C1/C2; sem C3 por padrão)
  - `event_id`, `user_id`, `timestamp`, `local_date`, `week_id`
  - `type`, `payload_min` (IDs, totals, counts; sem conteúdo sensível)
  - **Retenção**: 90 dias (detalhes) ou conforme policy.

- **DailyProgressSnapshot** (C1)
  - `user_id`, `local_date`
  - `plan_ref`, `tasks_status_summary` (counts e IDs)
  - `by_domain` (ex.: inglês: blocks done/pending; java: done/attempt; sono: diary recorded)
  - `updated_at`
  - **Retenção**: 90 dias.

- **ConsistencyMetric** (C4)
  - `user_id`, `week_id`, `domain`
  - `days_done`, `days_total`, `trend_vs_prev_week?`
  - **Retenção**: 12 meses.

- **RubricRecord** (C2→C4)
  - `user_id`, `local_date`, `week_id`, `domain`, `task_id`
  - `dimensions` (map), `total`, `max_total`, `status complete|partial`
  - **Retenção**: detalhe 90 dias; agregados semanais 12 meses.

- **RubricWeeklyAggregate** (C4)
  - `user_id`, `week_id`, `domain`
  - `avg_total`, `count`, `trend_vs_prev_week?`

- **RecurringErrorAggregate** (C4)
  - `user_id`, `domain`, `label`
  - `count_14d`, `count_week`, `trend_vs_prev_week?`, `status active|recurring|target|archived`, `last_seen_at`

- **SleepWeeklyAggregate** (C4)
  - `user_id`, `week_id`
  - `avg_quality?`, `avg_energy?`, `avg_regularity_delta?`, `days_with_diary`

- **HealthWeeklyAggregate** (C4)
  - `user_id`, `week_id`
  - `total_min`, `days_active`, `signals_summary`

- **SelfEsteemWeeklyAggregate** (C4)
  - `user_id`, `week_id`
  - `records_count`, `avg_intensity?`, `courage_done_count`, `top_tags_abstract[]`

- **SaasWeeklyAggregate** (C4)
  - `user_id`, `week_id`
  - `bet_status`, `weeks_completed`, `weeks_attempted`, `weeks_zero`

## 8) Regras e defaults
- **Erro recorrente (default)**: ≥3 ocorrências em janela rolante de 14 dias (ou ≥3 na mesma semana).
- **Dados suficientes (default)**:
  - para médias semanais: ≥3 registros no período
  - para comparação com semana anterior: semana anterior também ≥3
- **Gaps significativos**: ≥2 dias consecutivos sem registros ou ≥3 dias na semana sem registros (edge cases).
- **Targets de performance**:
  - progresso do dia: 1 mensagem; médio ≤5s; p95 ≤15s
  - revisão/semana: ≤2 mensagens; médio ≤10s; p95 ≤30s
- **Privacidade/retensão** (`SPEC-015`):
  - detalhes diários: 90 dias
  - agregados: 12 meses
  - C3 bruto: 7 dias ou “não guardar” (fora do escopo desta SPEC, mas deve ser compatível)

## 9) Observabilidade e métricas
- **Eventos**
  - `daily_snapshot_updated`, `weekly_aggregate_updated`
  - `recurring_error_marked`, `rubric_trend_computed`
  - `gap_detected(significant=true|false)` (sem proatividade automática)
- **Métricas**
  - latência de `GetTodayProgress` e `GetWeeklyMetrics` (médio/p95)
  - fricção de registro (proxy: nº de prompts manuais/dia; target ≤2 min/dia)

## 10) Riscos & mitigação
- **Risco**: queries lentas com histórico grande.  
  **Mitigação**: agregados semanais + retenção + índices por (user_id, week_id, local_date).
- **Risco**: dados expirados quebram revisão.  
  **Mitigação**: usar agregados; marcar parcial; sugerir coleta mínima.
- **Risco**: agregados vazam sensibilidade.  
  **Mitigação**: C4 não contém C3; labels neutros; redaction.

## 11) Rollout / migração
- **Feature flag**: `metrics_v1`.
- Backfill: opcional — recomputar agregados de semana atual/anterior ao habilitar.

## 12) Plano de testes (como validar)
- **Unit**
  - cálculo de consistência condicionado a gate
  - regra de erro recorrente 14d
  - “dados suficientes” e comparação vs semana anterior
- **Integration**
  - ingest eventos → snapshots diários → agregados semanais
  - apagamento/expiração reduz detalhes mas mantém agregados
- **E2E**
  - dia com tarefas em progresso → progresso do dia marca corretamente sem contar como concluído
  - semana inicial sem baseline anterior → revisão parcial
- **Manual**
  - respostas curtas e dentro de targets (1–2 mensagens)

## 13) Task breakdown (execução)
1) Definir schemas de snapshots e agregados (C4) + índices lógicos  
   - **Entrada**: `SPEC-016` FR-001..FR-007  
   - **Saída**: modelo lógico completo  
   - **Critério de pronto**: cada agregação tem regra de “dados suficientes” explícita

2) Implementar ingestão de eventos mínimos por domínio  
   - **Entrada**: planos de domínio + `SPEC-003`  
   - **Saída**: `IngestDomainEvent` com payload mínimo sem C3  
   - **Critério de pronto**: eventos suficientes para calcular consistência/rubricas/erros/sono

3) Implementar builder de progresso do dia (snapshot)  
   - **Entrada**: User Story 1 + AC-001  
   - **Saída**: `GetTodayProgress` baseado em snapshot  
   - **Critério de pronto**: 1 mensagem; diferencia `in_progress` vs `completed` vs `evidence_pending`

4) Implementar agregador semanal (semana atual vs anterior)  
   - **Entrada**: FR-006/FR-007 + defaults de suficiência  
   - **Saída**: `GetWeeklyMetrics` e job incremental  
   - **Critério de pronto**: trends corretas e “dados insuficientes” quando aplicável

5) Implementar tracker de erros recorrentes (14d) e status  
   - **Entrada**: FR-004 + SC-004  
   - **Saída**: `RecurringErrorAggregate` atualizado por evento  
   - **Critério de pronto**: 100% detecção para labels idênticos; agrupamento semântico pode ficar “later”

6) Integrar retenção/expiração com `SPEC-015`  
   - **Entrada**: política de retenção  
   - **Saída**: expiração automática e comportamento em dados faltantes  
   - **Critério de pronto**: queries continuam rápidas; revisão semanal funciona como parcial

## 14) Open questions (se existirem)
- (Default adotado) **Agrupamento semântico de erros**: MVP só garante 100% para labels idênticos; equivalência manual pode ser adicionada depois (como descrito na SPEC).

