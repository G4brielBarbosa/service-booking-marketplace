# Technical Plan: PLAN-012 — Vida saudável: planejamento semanal + métricas + limites por dor/energia

**Created**: 2026-02-20  
**Spec**: `specs/SPEC-012-vida-saudavel.md`  
**PRD Base**: §8.4, §6.1, §5.2, §5.3, §9.1, §9.2, §14, §11 (RNF1–RNF4)  
**Related Specs**: `SPEC-002`, `SPEC-003`, `SPEC-007`, `SPEC-011`, `SPEC-015`, `SPEC-016`, `SPEC-010`

## 1) Objetivo do plano
- Implementar criação de **plano semanal simples** de saúde (força + atividade moderada) com versão mínima viável.
- Implementar adaptação diária por **tempo/energia** (Plano A/B/C) e por **dor/fadiga** (segurança).
- Implementar registro mínimo (minutos/sessões + sinais de dor/fadiga) e resumo semanal consumível pela revisão semanal (`SPEC-007`) via métricas (`SPEC-016`).

## 2) Non-goals (fora do escopo)
- Não prescrever treino clínico, reabilitação ou orientação médica específica.
- Não exigir wearables, biometria avançada, contagem de calorias/macros.
- Não criar planilhas/dashboards avançados; foco em “insight suficiente para ajustar”.
- Não implementar nudges proativos no MVP; se houver, seguir `SPEC-011`.

## 3) Assumptions (assunções)
- “Saúde” é tratada como meta de fundação (não conta no limite de intensivas) (`SPEC-010`).
- O usuário pode ter semanas com zero dados; o sistema deve propor recomeço mínimo sem culpa.
- Privacidade: registros de dor/fadiga podem ser sensíveis; manter minimização e permitir apagar (`SPEC-015`).

## 4) Decisões técnicas (Decision log)
- **D-000 — Baseline de plataforma**
  - **Decisão**: adotar o baseline `plans/PLAN-000-platform-baseline.md` (persistência Postgres/sqlc, agregados `SPEC-016`, worker/jobs, retenção/privacidade) como base.
  - **Motivo**: saúde alimenta revisão semanal e métricas; precisa dos mesmos padrões de retenção, redaction e performance.
  - **Alternativas consideradas**: stack por feature; descartado.
  - **Impactos/Trade-offs**: baseline central concentra decisões de execução/armazenamento.

- **D-001 — Plano semanal como entidade explícita com versão mínima**
  - **Decisão**: modelar `HealthWeeklyPlan` com sessões planejadas e “versão mínima” para semanas ruins.
  - **Motivo**: suporta FR-001/FR-002 e reduz “tudo ou nada”.
  - **Alternativas consideradas**: apenas recomendações diárias sem plano; descartado (menos direção).
  - **Impactos/Trade-offs**: precisa de `week_id` e persistência.

- **D-002 — Sinais de segurança (dor/fadiga) como primeira classe**
  - **Decisão**: registrar `HealthSignal` (dor/fadiga, intensidade simples) e ajustar recomendações diárias quando presentes.
  - **Motivo**: FR-005 e sustentabilidade (evitar excesso).
  - **Alternativas consideradas**: ignorar sinais; descartado (risco).
  - **Impactos/Trade-offs**: exige linguagem cuidadosa “não é prescrição médica”.

- **D-003 — Gate de saúde é leve (hábito/fundação)**
  - **Decisão**: usar gate proporcional (`SPEC-003`): conclusão = registro mínimo de execução (feito/parcial) + minutos aproximados; em dia ruim aceitar alternativa leve como consistência mínima.
  - **Motivo**: RNF1/RNF2; saúde precisa ser aderível.
  - **Alternativas consideradas**: exigir detalhes de treino; descartado.
  - **Impactos/Trade-offs**: menos granularidade; suficiente para tendência semanal.

## 5) Arquitetura (alto nível)
- **Componentes**
  - **Health Domain Service**: cria plano semanal, sugere atividade do dia, registra execução e sinais.
  - **Gate Engine** (`SPEC-003`): gate leve para tarefas de saúde.
  - **Metrics Aggregator** (`SPEC-016`): minutos/sessões, consistência semanal, sinais.
  - **Weekly Review Integration** (`SPEC-007`): consome `HealthWeeklySummary`.
  - **Privacy Service** (`SPEC-015`): minimização, retenção e apagamento.

- **Fluxos**
  - Planejar semana → registrar plano + versão mínima.
  - No dia: sugerir sessão do plano ou alternativa curta (A/B/C) + ajuste por dor/fadiga.
  - Registrar execução (feito/parcial) + sinal (se houver) → alimentar agregados semanais.

## 6) Contratos e interfaces
- **Comando**: `CreateHealthWeeklyPlan(user_id, week_id, constraints?, timestamp)`
  - **Saída**: `HealthWeeklyPlanView(plan, minimal_plan, why_short)`

- **Consulta**: `GetHealthTodaySuggestion(user_id, local_date, day_context, pain_fatigue? , timestamp)`
  - **Saída**: `HealthDaySuggestion(type, duration_est_min, done_criteria_short, safety_note_short)`

- **Comando**: `RecordHealthActivity(user_id, local_date, activity_type, duration_min?, result done|partial|not_done, timestamp)`
  - **Saída**: `ActivityReceipt(consistency_counted?, next_suggestion_short)`

- **Comando**: `RecordHealthSignal(user_id, local_date, signal_type pain|fatigue, intensity_0_10, note_short?, timestamp)`
  - **Saída**: `SignalReceipt(adaptation_applied_short)`

- **Consulta**: `GetHealthWeekSummary(user_id, week_id)`
  - **Saída**: `HealthWeeklySummary(total_min, days_active, signals, adjustment_recommended)`

## 7) Modelo de dados (mínimo)
- **HealthWeeklyPlan**
  - `user_id`, `week_id`
  - `planned_sessions[]` (force/moderate; dia opcional)
  - `minimal_plan[]` (versão mínima)
  - `status active|archived`

- **HealthActivityLog**
  - `user_id`, `local_date`
  - `type force|moderate|light_alternative`
  - `duration_min?`, `result done|partial|not_done`
  - **Retenção**: detalhes 90 dias; agregados 12 meses (`SPEC-016`).

- **HealthSignal**
  - `user_id`, `local_date`
  - `type pain|fatigue`, `intensity_0_10`, `note_short?`
  - **Privacidade**: `note_short` opcional/minimizada; redigível/apagável (`SPEC-015`).

- **HealthWeeklyAggregates** (`SPEC-016`)
  - `user_id`, `week_id`
  - `total_min`, `days_active`, `signals_count_by_type`, `avg_signal_intensity?`

## 8) Regras e defaults
- Plano semanal default: 1–2 força + 2–4 moderada; versão mínima: 1 força curta + 2 caminhadas curtas (ajustável).
- Dia ruim (Plano C): alternativa leve e segura (ex.: caminhada curta/mobilidade leve) + registro mínimo.
- Dor/fadiga acima do normal: reduzir intensidade/substituir por alternativa leve; se persistente, sugerir “semana de redução”.
- Privacidade: não coletar detalhes excessivos; permitir apagar sinais e logs.

## 9) Observabilidade e métricas
- **Eventos**: `health_plan_created`, `health_activity_recorded`, `health_signal_recorded`, `health_week_summary_generated`
- **Métricas**: dias ativos/semana, minutos totais/semana, sinais de excesso (tendência), semanas “zero”.

## 10) Riscos & mitigação
- **Risco**: virar “tudo ou nada”. Mitigar com versão mínima e linguagem protetiva.
- **Risco**: aconselhamento médico indevido. Mitigar com disclaimers e alternativas seguras genéricas.
- **Risco**: dados insuficientes. Mitigar com resumo parcial e pedido de 1 registro mínimo.

## 11) Rollout / migração
- **Feature flag**: `health_v1`.

## 12) Plano de testes (como validar)
- **Unit**: geração do plano e versão mínima; adaptação por dor/fadiga; gate leve.
- **Integration**: logs → agregados semanais; resumo semanal consumível pela revisão semanal.
- **E2E**: semana com dados vs semana sem dados; dia ruim com alternativa leve.
- **Manual**: tom não punitivo; privacidade (apagamento de sinais).

## 13) Task breakdown (execução)
1) Definir `HealthWeeklyPlan` + contratos de criação/consulta  
   - **Entrada**: `SPEC-012` FR-001/FR-002  
   - **Saída**: schema e handlers  
   - **Critério de pronto**: plano e versão mínima gerados em poucos passos

2) Implementar sugestão diária com adaptação por contexto e dor/fadiga  
   - **Entrada**: FR-003/FR-005  
   - **Saída**: `GetHealthTodaySuggestion` com alternativas seguras  
   - **Critério de pronto**: dia ruim sempre tem passo executável

3) Implementar registro mínimo de execução e sinais + agregados semanais  
   - **Entrada**: FR-004/FR-006 + `SPEC-016`  
   - **Saída**: logs + `HealthWeeklySummary`  
   - **Critério de pronto**: revisão semanal lê total/consistência/sinais sem burocracia

4) Instrumentar eventos e retenção/apagamento  
   - **Entrada**: `SPEC-015` + `SPEC-016`  
   - **Saída**: eventos e política de retenção (90d detalhes, 12m agregados)  
   - **Critério de pronto**: usuário pode apagar período/categoria sem quebrar tendências futuras

## 14) Open questions (se existirem)
- (Default adotado) **Intensidade de dor/fadiga “alta”**: usar limiar simples ≥7/10 para acionar substituição por alternativa leve e registrar aviso “semana de redução” se repetido ≥3x/semana.

