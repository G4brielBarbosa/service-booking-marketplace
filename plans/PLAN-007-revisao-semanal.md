# Technical Plan: PLAN-007 — Revisão Semanal: Painel mínimo + 3 decisões + alvos da semana

**Created**: 2026-02-20  
**Spec**: `specs/SPEC-007-revisao-semanal.md`  
**PRD Base**: §§5.2, 5.5, 9.2, §6.2, §14, §10 (R5, R4, R7), §11 (RNF1–RNF3)  
**Related Specs**: `SPEC-002`, `SPEC-003`, `SPEC-010`, `SPEC-011`, `SPEC-015`, `SPEC-016`, `SPEC-008`, `SPEC-009`

## 1) Objetivo do plano
- Implementar geração de um **painel semanal mínimo** (texto/conversacional) com consistência, qualidade (gates/rubricas), gargalos e sono/energia.
- Implementar o fluxo das **3 decisões** por meta/ciclo (manter/ajustar/pausar) e registrar o resultado.
- Implementar definição e registro de **alvos da semana** (limite) e um resumo final único e acionável.

## 2) Non-goals (fora do escopo)
- Não implementar dashboards gráficos/BI nem análises estatísticas avançadas.
- Não implementar backlog inteligente completo (`SPEC-008`) nem detecção avançada de falso progresso (`SPEC-009`) — apenas **consumir sinais existentes** e registrá-los.
- Não enviar revisão semanal proativa por padrão; no MVP a revisão é **user-initiated**. Se houver proatividade no futuro, deve seguir `SPEC-011`.

## 3) Assumptions (assunções)
- As métricas e registros necessários existem via `SPEC-016` (ou serão criados por seus PLANs): consistência, rubricas, erros recorrentes, sono/energia.
- `week_id` é calculado no timezone do usuário (ex.: ISO week).
- Governança de metas (limite de 2 intensivas) existe e pode ser consultada (`SPEC-010`).

## 4) Decisões técnicas (Decision log)
- **D-001 — Revisão semanal como agregação + conversa guiada**
  - **Decisão**: gerar `WeeklyPanel` a partir de agregados (`SPEC-016`) e conduzir decisões em passos curtos, registrando `WeeklyReviewResult`.
  - **Motivo**: atende FR-001..FR-004 com baixa fricção (10–20 min).
  - **Alternativas consideradas**: revisão como “relatório longo”; descartado (RNF1).
  - **Impactos/Trade-offs**: precisa priorizar informação e limitar texto.

- **D-002 — Dados parciais viram “revisão parcial” explícita**
  - **Decisão**: quando registros forem insuficientes, marcar status `partial`, listar lacunas e propor recuperação mínima (1–2 ações).
  - **Motivo**: `SPEC-007` FR-005 + RNF2/RNF3.
  - **Alternativas consideradas**: “bronca” por dados faltantes; descartado.
  - **Impactos/Trade-offs**: decisões baseadas em menos evidência; deve explicitar incerteza.

- **D-003 — Alvos da semana com limite duro e foco por domínio**
  - **Decisão**: máximo 3 alvos (preferencialmente 1 inglês, 1 java, 1 fundação), com regra “o que NÃO fazer” no resumo final.
  - **Motivo**: evitar overload (`PRD §6.2`) e seguir edge cases da SPEC.
  - **Alternativas consideradas**: alvos ilimitados; descartado.
  - **Impactos/Trade-offs**: restringe ambição; melhora consistência.

## 5) Arquitetura (alto nível)
- **Componentes**
  - **WeeklyReview Service**: gera painel, conduz decisões e registra resultado.
  - **Metrics Query Layer** (`SPEC-016`): fornece agregados semana atual/anterior e sinais (erros recorrentes, rubricas).
  - **Governance Service** (`SPEC-010`): fornece metas ativas/pausadas e aplica limite de 2 intensivas.
  - **Privacy Service** (`SPEC-015`): controla nível de detalhe e expõe “dados que tenho”/apagamento quando solicitado durante revisão.

- **Fluxos**
  - **Iniciar revisão**: usuário pede “revisão semanal” → gerar painel (ou parcial) → apresentar em 1–2 mensagens.
  - **Decisões**: iterar metas relevantes e coletar manter/ajustar/pausar (uma por vez).
  - **Alvos**: sugerir alvos com base em sinais (erros/rubrica/consistência/sono) → usuário escolhe/ajusta → registrar.
  - **Encerrar**: emitir resumo final único com “achados + decisões + alvos + regra de foco”.

## 6) Contratos e interfaces
- **Comando**: `StartWeeklyReview(user_id, week_id, timestamp)`
  - **Saída**: `WeeklyPanelView(status complete|partial, panel_sections[], data_gaps[])`

- **Comando**: `RecordWeeklyDecision(user_id, week_id, goal_domain, decision keep|adjust|pause, note_short?, timestamp)`
  - **Saída**: `DecisionReceipt(impact_preview_short)`
  - **Erros**: `DECISION_INVALID` (ex.: pausar todas intensivas sem manter fundação — sugerir mínimo)

- **Comando**: `SetWeeklyTargets(user_id, week_id, targets[], timestamp)`
  - **Saída**: `TargetsReceipt(rule_of_focus_short)`
  - **Erros**: `TARGETS_TOO_MANY` (aplicar limite)

- **Consulta**: `GetWeeklyReviewSummary(user_id, week_id)`
  - **Saída**: `WeeklyReviewSummary(final_message, decisions, targets)`

- **Consulta**: `GetWeeklyDrilldown(user_id, week_id, dimension)` (opcional MVP)
  - **Saída**: detalhe curto por dia/meta (1 pergunta por vez), sem virar relatório.

## 7) Modelo de dados (mínimo)
- **WeeklyPanel** (`SPEC-007`/`SPEC-016`)
  - `user_id`, `week_id`
  - `consistency_by_goal[]` (dias/semana)
  - `quality_signals[]` (gates pass rate, rubrica médias)
  - `bottlenecks[]` (tempo/energia/atrito — quando disponível)
  - `sleep_energy_trend` (regularidade + médias)
  - `data_gaps[]`

- **WeeklyDecision**
  - `user_id`, `week_id`, `goal_domain`
  - `decision keep|adjust|pause`
  - `justification_short?`
  - `expected_impact_short`

- **WeeklyTargets**
  - `user_id`, `week_id`
  - `targets[]` (cada: domínio, label, critério observável, frequência mínima)
  - `rule_of_focus_short`

- **WeeklyReviewResult**
  - `user_id`, `week_id`
  - `status complete|partial`
  - `panel_ref`, `decisions_ref`, `targets_ref`
  - `final_summary_text` (curto)
  - **Retenção** (`SPEC-016`): agregados 12 meses (C4); detalhes diários 90 dias.

## 8) Regras e defaults
- **Painel mínimo obrigatório** (`SPEC-007` FR-001):
  - consistência por meta
  - qualidade (gates/rubricas)
  - gargalos (tempo/energia/atrito) quando houver
  - sono/energia (tendência simples)
- **Revisão parcial**: quando dados insuficientes, marcar e não comparar com semana anterior (`SPEC-016`).
- **Decisões**:
  - se usuário não decide, aplicar default conservador: manter fundação + 1 meta intensiva; pedir 1 confirmação curta.
  - respeitar limite de 2 metas intensivas ativas (orientar pausa/adiamento) (`SPEC-010`).
- **Alvos**: máximo 3; preferir 1 por domínio; sugerir com base em:
  - erros recorrentes (≥3/14d) e rubrica em queda
  - consistência baixa (ex.: ≤2/7)
  - sono/energia baixo (foco em fundação)
- **Anti-spam**: revisão semanal é iniciada pelo usuário no MVP. Se houver lembrete futuro, seguir budgets/quiet hours (`SPEC-011`).
- **Privacidade**: painel deve evitar expor conteúdo sensível; usar labels neutros e agregados (`SPEC-015`).

## 9) Observabilidade e métricas
- **Eventos**
  - `weekly_review_started`, `weekly_panel_generated(status)`
  - `weekly_decision_recorded(goal, decision)`
  - `weekly_targets_set(count)`
  - `weekly_review_completed(status)`

- **Métricas**
  - **SC-001 (SPEC-007)**: % de semanas com revisão completa ou parcial.
  - Tempo médio para completar revisão (proxy: nº de interações e duração).
  - Taxa de overload detectada/decisões de pausa (alinhado a `SPEC-010`).

## 10) Riscos & mitigação
- **Risco**: revisão vira longa e cansativa.  
  **Mitigação**: painel mínimo, limite de targets, 1 pergunta por vez, resumo final único.
- **Risco**: dados insuficientes geram conclusões erradas.  
  **Mitigação**: status `partial`, explicitar lacunas e sugerir recuperação mínima.
- **Risco**: usuário evita decidir.  
  **Mitigação**: defaults conservadores + 1 confirmação curta.

## 11) Rollout / migração
- **Feature flag**: `weekly_review_v1`.
- Evolução: permitir agendamento/lembrte semanal apenas após `SPEC-011` estar implementada e com opt-out claro (`SPEC-015`).

## 12) Plano de testes (como validar)
- **Unit**
  - Montagem do painel a partir de métricas (`SPEC-016`), incluindo regras de “dados suficientes”.
  - Limite de targets e defaults conservadores.
  - Respeito a governança (máx 2 intensivas).
- **Integration**
  - Semana com dados suficientes → status complete; semana com poucos dados → partial + recuperação mínima.
  - Registro de decisões e targets persistidos e consultáveis.
- **E2E**
  - Fluxo inteiro: iniciar → painel → decisões → targets → resumo final.
- **Manual / acceptance**
  - Tom não punitivo em semana ruim e overload.
  - Painel não expõe conteúdo sensível por padrão.

## 13) Task breakdown (execução)
1) **Definir contrato de agregados consumidos (dependência `SPEC-016`)**
   - **Entrada**: `SPEC-007` FR-001 + `SPEC-016` FR-001..FR-007
   - **Saída**: interface `WeeklyMetricsSnapshot` (sem conteúdo sensível) com campos mínimos
   - **Critério de pronto**: painel mínimo pode ser gerado só com esse snapshot

2) **Modelar entidades `WeeklyPanel/Decision/Targets/Result` e retenção**
   - **Entrada**: `SPEC-007` Key Entities + `SPEC-015/016`
   - **Saída**: schema lógico + política de retenção (12m agregados)
   - **Critério de pronto**: apagar dados por período/categoria é possível sem quebrar o fluxo

3) **Implementar gerador do painel (complete vs partial)**
   - **Entrada**: cenários P1/P2 do `SPEC-007`
   - **Saída**: `StartWeeklyReview` produzindo `WeeklyPanelView` em formato curto
   - **Critério de pronto**: com 1–2 dias de dados, retorna partial e lista lacunas

4) **Implementar registro das 3 decisões por meta e impacto esperado curto**
   - **Entrada**: `SPEC-007` FR-002/FR-006 + `SPEC-010`
   - **Saída**: `RecordWeeklyDecision` + validações (limite intensivas)
   - **Critério de pronto**: decisões ficam consultáveis e influenciáveis por planejamentos futuros

5) **Implementar definição de alvos da semana com limite**
   - **Entrada**: `SPEC-007` FR-003 + edge cases
   - **Saída**: `SetWeeklyTargets` + `rule_of_focus_short`
   - **Critério de pronto**: nunca aceita >3 alvos; oferece recorte

6) **Implementar resumo final único**
   - **Entrada**: `SPEC-007` FR-004
   - **Saída**: `WeeklyReviewSummary.final_message` com achados+decisões+alvos
   - **Critério de pronto**: mensagem única, acionável, sem jargão, concluível em 10–20 min

7) **Instrumentar eventos e métricas de revisão semanal**
   - **Entrada**: `SPEC-007` SCs + `SPEC-016`
   - **Saída**: eventos `weekly_review_*` e contadores de completion
   - **Critério de pronto**: é possível medir % de semanas com revisão e tempo médio de conclusão

## 14) Open questions (se existirem)
- (Default adotado) **User-initiated MVP**: a revisão semanal não é enviada proativamente no MVP. Um lembrete semanal só entra após `SPEC-011` estar implementada com budgets + opt-out.

