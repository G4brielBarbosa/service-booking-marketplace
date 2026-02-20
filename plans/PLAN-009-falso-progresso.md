# Technical Plan: PLAN-009 — Detecção de falhas reais (“falso progresso”) e reforço automático

**Created**: 2026-02-20  
**Spec**: `specs/SPEC-009-deteccao-falso-progresso.md`  
**PRD Base**: §5.4, §5.5, §9.1, §9.2, §10 (R3, R2), §13, §11 (RNF1–RNF3)  
**Related Specs**: `SPEC-003`, `SPEC-016`, `SPEC-008`, `SPEC-007`, `SPEC-002`, `SPEC-011`, `SPEC-015`

## 1) Objetivo do plano
- Implementar detecção MVP de **sinais observáveis** de falso progresso (erro recorrente, queda de rubrica, retrieval fraco, gates falhando por ausência de evidência).
- Implementar acionamento de **um reforço curto e verificável** (com evidência mínima e versão mínima em dia ruim), registrando resultado.
- Integrar sinais e reforços com backlog (`SPEC-008`), revisão semanal (`SPEC-007`) e planejamento diário (`SPEC-002`) sem causar overload.

## 2) Non-goals (fora do escopo)
- Não fazer diagnóstico educacional/psicológico nem inferências “profundas”.
- Não otimizar com ML; MVP usa janelas/limiares determinísticos.
- Não acionar múltiplos reforços por dia; limitar intervenções para evitar burocracia.
- Não enviar nudges proativos por padrão; se houver, seguir `SPEC-011`.

## 3) Assumptions (assunções)
- Os dados necessários existem via `GateResult`, rubricas, retrieval e erros recorrentes (`SPEC-003` + `SPEC-016`).
- Reforço é uma micro-tarefa com `definition_of_done` e, se aplicável, gate (`SPEC-003`).
- O sistema consegue diferenciar “dia ruim” (Plano C) e reduzir reforço para versão mínima.

## 4) Decisões técnicas (Decision log)
- **D-001 — Sinais determinísticos com janela global**
  - **Decisão**: definir uma janela global MVP para detecção: `last_14d` (ou “últimas 5 sessões do domínio”, quando existir).
  - **Motivo**: consistência e testabilidade (FR-002).
  - **Alternativas consideradas**: janelas por domínio; adiado (complexidade).
  - **Impactos/Trade-offs**: pode ser menos sensível em alguns domínios; aceitável.

- **D-002 — “Um foco por vez” e limite diário de reforços**
  - **Decisão**: no máximo 1 reforço principal por dia (FR-004), priorizando o sinal de maior impacto.
  - **Motivo**: evitar overload e manter fricção baixa (RNF1).
  - **Alternativas consideradas**: múltiplos reforços; descartado.
  - **Impactos/Trade-offs**: alguns sinais ficam para depois; backlog captura.

- **D-003 — Reforço como entidade explícita com versão mínima**
  - **Decisão**: modelar `ReinforcementAction` com `definition_of_done` e `minimal_variant`.
  - **Motivo**: suportar dias ruins (RNF2) sem perder verificabilidade.
  - **Alternativas consideradas**: reforço ad hoc via texto; descartado (não auditável).
  - **Impactos/Trade-offs**: precisa de catálogo mínimo de reforços por tipo de sinal.

## 5) Arquitetura (alto nível)
- **Componentes**
  - **Signal Detector**: calcula `FalseProgressSignal` a partir de agregados e resultados recentes.
  - **Reinforcement Catalog**: mapeia tipo de sinal → reforço recomendado (normal/minimal).
  - **Reinforcement Orchestrator**: decide se aciona agora (limites) e registra outcome.
  - **Backlog Integration** (`SPEC-008`): cria/atualiza item quando sinal persistir ou reforço não for feito.
  - **Metrics Sink** (`SPEC-016`): registra eventos e contagens.

## 6) Contratos e interfaces
- **Comando**: `DetectSignals(user_id, now, window=14d)`
  - **Saída**: `SignalList(signals[])` (ordenados por severidade)

- **Comando**: `ProposeReinforcement(user_id, signal_id, day_context, timestamp)`
  - **Saída**: `ReinforcementProposal(description_short, definition_of_done, minimal_variant, why_short)`

- **Comando**: `RecordReinforcementOutcome(user_id, reinforcement_id, outcome done|not_done|attempt, timestamp, note_short?)`
  - **Saída**: `OutcomeReceipt(next_step_short)`

- **Comando**: `MaybeEnqueueBacklogItemFromSignal(user_id, signal_id)`
  - **Saída**: `BacklogItemRef?` (criado/atualizado)

- **Consulta**: `GetWeeklySignalsSummary(user_id, week_id)`
  - **Saída**: lista curta de sinais e ações tomadas (para `SPEC-007`).

## 7) Modelo de dados (mínimo)
- **FalseProgressSignal** (`SPEC-009`)
  - `signal_id`, `user_id`, `timestamp`
  - `type`: `recurring_error | rubric_drop | retrieval_low | gate_missing_evidence`
  - `severity`: `low|medium|high`
  - `window`: `14d`
  - `evidence_refs[]` (IDs de agregados/GateResults; sem C3)
  - `status`: `open|addressed|snoozed|resolved`
  - **Retenção**: agregados 12 meses; sinais podem expirar (ex.: 90 dias) mantendo contagens.

- **ReinforcementAction**
  - `reinforcement_id`, `user_id`, `signal_id`
  - `domain`, `description_short`, `definition_of_done`
  - `minimal_variant`
  - `status`: `proposed|done|not_done|attempt`
  - `gate_profile_id?` (quando aplicável; `SPEC-003`)
  - `created_at`, `updated_at`

- **SignalOutcome**
  - `outcome_id`, `signal_id`, `reinforcement_id?`
  - `result`: `done|not_done|attempt`
  - `note_short?`
  - `observed_effect?` (opcional MVP; ex.: “rubrica estabilizou”)

## 8) Regras e defaults
- **Detecção (MVP)**
  - `recurring_error`: erro recorrente ≥3/14d (`SPEC-016`).
  - `rubric_drop`: queda ≥1 ponto vs semana anterior com dados suficientes (`SPEC-016`).
  - `retrieval_low`: ≥2 ocorrências `low` na semana (ou 3/14d).
  - `gate_missing_evidence`: 3+ falhas por “missing evidence” na semana (`SPEC-003` SC-005).
- **Escolha do sinal principal**
  - priorizar: recurring_error/high > rubric_drop > retrieval_low > gate_missing_evidence (ajustável).
- **Limites**
  - máximo 1 reforço principal/dia; máximo 1 foco por vez.
  - cooldown para re-propor o mesmo reforço: 3 dias (exceto severidade high).
- **Dia ruim**
  - usar `minimal_variant` (≤2–5 min) e aceitar `attempt` como registro, sem “virar bronca”.
- **Privacidade**
  - reforços devem evitar conteúdo sensível; evidência bruta segue `SPEC-015`.

## 9) Observabilidade e métricas
- **Eventos**
  - `signal_detected(type,severity)`, `reinforcement_proposed`, `reinforcement_outcome(result)`
  - `signal_enqueued_to_backlog`, `signal_resolved`
- **Métricas**
  - Tempo até “primeiro reforço” após sinal high.
  - % de reforços concluídos (normal/minimal) e efeito em erro/rubrica nas semanas seguintes.
  - Taxa de intervenção (deve ser baixa e não burocrática).

## 10) Riscos & mitigação
- **Risco**: sistema parece “acusatório”.  
  **Mitigação**: linguagem protetiva, foco em processo, sinais como “dados” e não culpa.
- **Risco**: intervenção excessiva aumenta fricção.  
  **Mitigação**: limite 1/dia + minimal_variant + cooldown.
- **Risco**: dados insuficientes para sinal.  
  **Mitigação**: marcar “sinal fraco” e pedir 1 dado mínimo no próximo ciclo (sem acusar).

## 11) Rollout / migração
- **Feature flag**: `false_progress_v1`.
- Backfill: opcional — gerar sinais iniciais a partir das últimas 2 semanas quando houver dados.

## 12) Plano de testes (como validar)
- **Unit**
  - regras de detecção por limiar/janela
  - seleção de sinal principal e limites (1/dia, cooldown)
  - escolha do reforço (normal vs minimal)
- **Integration**
  - sinal → proposta → outcome → backlog item quando persistente
  - integração com gate profiles quando reforço exigir evidência mínima
- **E2E**
  - simular 3 ocorrências de erro → sinal high → reforço → outcome registrado → aparece na revisão semanal
- **Manual / acceptance**
  - tom não punitivo; explicação 1–2 frases; reforço pequeno

## 13) Task breakdown (execução)
1) **Definir schema de `FalseProgressSignal` e `ReinforcementAction`**
   - **Entrada**: `SPEC-009` FR-001..FR-005 + `SPEC-016`
   - **Saída**: modelo lógico + retenção/mínimos de privacidade
   - **Critério de pronto**: sinais não dependem de C3; evidência_refs são agregados

2) **Implementar detecção MVP (4 sinais) com janela 14d**
   - **Entrada**: `SPEC-009` FR-002 + defaults
   - **Saída**: `Signal Detector` determinístico
   - **Critério de pronto**: testes unit cobrem limiares e casos “dados insuficientes”

3) **Criar catálogo mínimo de reforços (normal/minimal) por tipo de sinal**
   - **Entrada**: `SPEC-009` FR-003/FR-004
   - **Saída**: mapeamento sinal → reforço + definition_of_done
   - **Critério de pronto**: minimal_variant ≤5 min e observável

4) **Implementar orquestração com limites (1/dia) e cooldown**
   - **Entrada**: FR-004 + edge cases
   - **Saída**: `Reinforcement Orchestrator`
   - **Critério de pronto**: nunca propõe 2 reforços no mesmo dia

5) **Integrar com backlog e revisão semanal**
   - **Entrada**: `SPEC-008` + `SPEC-007`
   - **Saída**: item de backlog quando sinal persistir ou reforço for rejeitado repetidamente; resumo semanal de sinais
   - **Critério de pronto**: revisão semanal mostra 1–2 sinais principais e ação tomada

6) **Instrumentar eventos/métricas**
   - **Entrada**: `SPEC-009` SCs
   - **Saída**: eventos `signal_*` e `reinforcement_*`
   - **Critério de pronto**: possível medir taxa de reforços e impacto em tendências

## 14) Open questions (se existirem)
- (Default adotado) **Sem proatividade no MVP**: sinais/reforços aparecem quando usuário gera/consulta plano, finaliza tarefa, ou entra na revisão semanal. Nudges proativos ficam para `SPEC-011`.

