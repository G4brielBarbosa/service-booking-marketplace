# Technical Plan: PLAN-005 — Java Diário: Prática Deliberada + Retrieval + Revisão de Erros

**Created**: 2026-02-20  
**Spec**: `specs/SPEC-005-java-diario.md`  
**PRD Base**: §8.2, §5.4, §5.3, §9.1, §14, §10 (R2, R3, R6), §11 (RNF1–RNF3)  
**Related Specs**: `SPEC-002`, `SPEC-003`, `SPEC-008`, `SPEC-009`, `SPEC-011`, `SPEC-015`, `SPEC-016`

## 1) Objetivo do plano
- Implementar o loop diário de Java com 3 componentes: **prática → retrieval → registro de erro/aprendizado** (normal e mínimo para dia ruim).
- Integrar com **quality gates** (`SPEC-003`) para que a sessão só conte como concluída com evidência mínima válida; parcial vira `attempt`.
- Consolidar erros recorrentes e métricas mínimas para tendências e revisão semanal (`SPEC-016`).

## 2) Non-goals (fora do escopo)
- Não definir editor/IDE, plataforma de exercícios, repositório, stack de testes, ou pipeline CI.
- Não implementar backlog inteligente/seleção ótima de exercícios (isso é `SPEC-008/009`).
- Não implementar nudges proativos (seguir `SPEC-011` quando existir).

## 3) Assumptions (assunções)
- A rotina diária (`SPEC-002`) consome templates de Java do `TaskCatalog`; este PLAN define como templates/gates e registros se conectam.
- Evidências de Java podem ser textuais no MVP (descrição curta do que foi feito, respostas de retrieval, e registro de erro), reduzindo dependência de conteúdo sensível.
- Recorrência de erros segue default do `SPEC-016` (≥3 em 14 dias).

## 4) Decisões técnicas (Decision log)
- **D-001 — Separar evidência de prática, retrieval e aprendizado**
  - **Decisão**: modelar evidência mínima como três artefatos pequenos: `JavaPracticeEvidence`, `RetrievalResult`, `LearningLogEntry`.
  - **Motivo**: mantém baixa fricção e permite detecção de sinais (retrieval fraco/erros recorrentes) (`SPEC-016`/`SPEC-009`).
  - **Alternativas consideradas**: exigir código/PR/arquivo; adiado (alto atrito no MVP).
  - **Impactos/Trade-offs**: evidência é mais “auto-report”; mitigada por retrieval e repetição de erro.

- **D-002 — Gate de Java é “mínimo observável” (não burocrático)**
  - **Decisão**: para concluir: (a) 1 registro de prática (curto), (b) retrieval (1–3 itens no mínimo), (c) 1 erro/aprendizado (curto). Se faltar algo, vira `attempt` e retorna `next_min_step`.
  - **Motivo**: `SPEC-003` (aprendizagem exige evidência) + `SPEC-005` (baixa fricção).
  - **Alternativas consideradas**: gates pesados (testes, cobertura); adiado.
  - **Impactos/Trade-offs**: pode aceitar “evidência fraca”; métricas de falha/recorrência sinalizam necessidade de reforço.

- **D-003 — Retrieval pode disparar reforço mínimo**
  - **Decisão**: quando `RetrievalResult.status=low`, solicitar 1 reforço mínimo observável (ex.: reexplicar em 2–4 frases) e registrar como `ReinforcementAttempt` (ponte para `SPEC-009`).
  - **Motivo**: evitar “fiz mas não aprendi” sem virar burocracia (PRD R3).
  - **Alternativas consideradas**: ignorar retrieval; descartado.
  - **Impactos/Trade-offs**: adiciona 1 micro-pass adicional em dias ruins; deve respeitar Plano C e limites.

## 5) Arquitetura (alto nível)
- **Componentes**
  - **Java Domain Service**: templates A/B/C, captura de evidência mínima e registro de erros.
  - **Gate Engine** (`SPEC-003`): avalia gate do “Java session task”.
  - **Metrics/Recurring Errors** (`SPEC-016`): consolida consistência, retrieval, erros recorrentes.
  - **Backlog/Signals (futuro)**: `SPEC-008/009` consomem sinais (não implementado aqui).

- **Fluxos**
  - **Sessão Java (A/B)**: prática (descrição curta do que fez) → retrieval → erro/aprendizado → gate satisfeito.
  - **Sessão Java (C)**: micro-exercício + 1 retrieval + 1 aprendizado → gate mínimo; se incompleto, `attempt`.

## 6) Contratos e interfaces
- **Comando**: `StartJavaSessionTasks(user_id, local_date, plan_type)`
  - **Saída**: `TaskTemplate` da sessão (com `gate_profile_id`)

- **Comando**: `SubmitJavaPracticeEvidence(user_id, task_id, evidence_short, timestamp)`
  - **Saída**: `EvidenceReceipt` + `GateStatus(evidence_pending|ready_to_evaluate)`

- **Comando**: `SubmitJavaRetrieval(user_id, task_id, answers[], timestamp)`
  - **Saída**: `RetrievalSummary(status ok|low, targets[])`

- **Comando**: `SubmitJavaLearningLog(user_id, task_id, error_or_learning, fix_or_note, category?, timestamp)`
  - **Saída**: `LearningLogRef` + `RecurringStatus(active|recurring|target)`

- **Comando**: `EvaluateJavaGate(user_id, task_id, timestamp)`
  - **Saída**: `GateResultView(reason_short, next_min_step)`

- **Consulta**: `GetJavaWeekTrend(user_id, week_id)`
  - **Saída**: consistência, retrieval_ok_rate, top erros recorrentes

## 7) Modelo de dados (mínimo)
- **JavaSession**
  - `user_id`, `local_date`, `task_id`
  - `objective_constraint` (texto curto), `planned_min`, `status`

- **JavaPracticeEvidence**
  - `evidence_id`, `task_id`, `evidence_short` (1–3 linhas), `validity`
  - **Sensibilidade**: C2; deve evitar incluir código longo/segredos; pode ser redigível.

- **RetrievalResult**
  - `task_id`, `count_items`, `status ok|low`, `targets[]`
  - **Retenção**: agregados (C4) 12 meses; detalhes mínimos 90 dias.

- **LearningLogEntry**
  - `task_id`, `category?`, `error_or_learning`, `fix_or_note`
  - **Privacidade**: conteúdo curto e minimizado; redigível.

- **RecurringError** (`SPEC-016`)
  - `user_id`, `domain=java`, `label/category`, `count_14d`, `last_seen_at`, `status`

- **GateResult** (`SPEC-003`)
  - referência ao perfil `java_session_v1` e aos evidence ids.

## 8) Regras e defaults
- **Plano A/B (dia normal)**: prática 20–45 min + retrieval 5–10 min (2–5 itens) + erro/aprendizado 1 item.
- **Plano C (dia ruim)**: micro-exercício 5–10 min + retrieval 1 item + 1 erro/aprendizado.
- **Gate mínimo (aprendizagem)** (`SPEC-003`): sem os 3 componentes mínimos → não conclui; vira `attempt` com `next_min_step` único.
- **Erro recorrente**: ≥3 ocorrências em 14 dias (default `SPEC-016`).
- **Privacidade** (`SPEC-015`): não coletar conteúdo desnecessário; permitir apagar logs; retenção padrão 90 dias (detalhes) + 12 meses (agregados).

## 9) Observabilidade e métricas
- **Eventos**
  - `java_practice_submitted`, `java_retrieval_submitted(status)`, `java_learning_logged(category)`
  - `java_gate_evaluated(satisfied|not, reason_code)`
  - `java_recurring_error_marked`

- **Métricas**
  - Consistência semanal de Java com gates satisfeitos.
  - Taxa de retrieval `ok` vs `low` por semana.
  - Frequência de erros recorrentes e tendência de redução.

## 10) Riscos & mitigação
- **Risco**: evidência “fraca” vira autoengano.  
  **Mitigação**: retrieval obrigatório mínimo + recorrência de erros + reforços (`SPEC-009`) quando sinais aparecerem.
- **Risco**: fricção alta em dias ruins.  
  **Mitigação**: Plano C com 1 item de retrieval e 1 aprendizado; 1 próximo passo único se faltar evidência.

## 11) Rollout / migração
- **Feature flag**: `java_daily_v1`.
- Evolução: permitir anexar referência externa (link/commit) no futuro sem exigir no MVP.

## 12) Plano de testes (como validar)
- **Unit**
  - Gate mínimo para Java (3 componentes) e `attempt` quando incompleto.
  - Recorrência (≥3/14d) funciona.
  - Retrieval low dispara sugestão de reforço mínimo (registrável).
- **Integration**
  - Fluxo completo: prática → retrieval → learning log → gate satisfied.
  - Fluxo parcial: prática apenas → gate not_satisfied + next_min_step.
- **E2E**
  - Dia normal vs dia ruim, com consulta de “o que já fiz hoje” refletindo `attempt/evidence_pending`.
- **Manual / acceptance**
  - Linguagem firme e não punitiva em falhas repetidas (RNF3).

## 13) Task breakdown (execução)
1) **Definir template A/B/C de Java + gate profile `java_session_v1`**
   - **Entrada**: `SPEC-005` FR-001..FR-004 + `SPEC-003`
   - **Saída**: `TaskTemplate` com instruções + critérios de feito + `gate_profile_id`
   - **Critério de pronto**: Plano C cabe em 5–15 min e é observável

2) **Modelar entidades mínimas e retenção (C2/C4)**
   - **Entrada**: `SPEC-015` + `SPEC-016`
   - **Saída**: schema lógico para evidências/retrieval/logs e agregados
   - **Critério de pronto**: detalhes expiram (90d) e agregados persistem (12m)

3) **Implementar submissão de evidência de prática + retrieval**
   - **Entrada**: cenários P1 do `SPEC-005`
   - **Saída**: comandos `SubmitJavaPracticeEvidence` e `SubmitJavaRetrieval`
   - **Critério de pronto**: validações básicas (não vazio) e registro de status ok/low

4) **Implementar registro de erro/aprendizado e atualização de recorrência**
   - **Entrada**: `SPEC-005` FR-006/FR-007 + `SPEC-016`
   - **Saída**: `LearningLogEntry` + `RecurringError` atualizado
   - **Critério de pronto**: 3 ocorrências em 14 dias marcam como recorrente

5) **Integrar avaliação do gate e atualização do status da tarefa**
   - **Entrada**: `SPEC-003` + integração com `SPEC-002`
   - **Saída**: `GateResultView` e transição `evidence_pending/attempt/completed`
   - **Critério de pronto**: nunca marcar `completed` sem gate satisfeito

6) **Instrumentar eventos e agregados semanais**
   - **Entrada**: `SPEC-016`
   - **Saída**: eventos `java_*` + métricas de consistência/retrieval
   - **Critério de pronto**: revisão semanal consegue consumir tendências

## 14) Open questions (se existirem)
- (Default adotado) **Objetivo/restrição da sessão**: no MVP, o `objective_constraint` pode ser preenchido com um template genérico (“praticar X por Y min + explicar conceito Z”) e refinado depois por backlog inteligente (`SPEC-008`).

