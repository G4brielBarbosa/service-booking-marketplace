# Technical Plan: PLAN-003 — Quality Gates & Evidência Mínima (aprendizagem e hábitos)

**Created**: 2026-02-20  
**Spec**: `specs/SPEC-003-quality-gates-evidencias.md`  
**PRD Base**: §5.4, §§8.1–8.3, §§5.2–5.3, §9.1, §10 (R2, R3), §11 (RNF1–RNF4), §13, §14  
**Related Specs**: `SPEC-002`, `SPEC-004`, `SPEC-005`, `SPEC-006`, `SPEC-009`, `SPEC-011`, `SPEC-015`, `SPEC-016`

## 1) Objetivo do plano
- Implementar um mecanismo cross-cutting de **Quality Gates** que controle quando uma tarefa pode virar **concluída** (aprendizagem só com evidência mínima válida; hábitos com gate proporcional).
- Implementar o ciclo de vida de **evidência → validação → GateResult** com transparência (“por que não contou” + “próximo passo mínimo”).
- Implementar integração com a rotina diária (`SPEC-002`) e com métricas/registros (`SPEC-016`) sem armazenar conteúdo sensível além do necessário (privacidade por padrão `SPEC-015`).

## 2) Non-goals (fora do escopo)
- Não implementar scoring automático, STT, análise de pronúncia, detecção perfeita de fraude ou ML.
- Não criar um “bypass” que marque aprendizagem como concluída sem evidência (a SPEC pede bloquear).
- Não definir o conteúdo específico de gates de cada domínio além do que já está nas SPECS de domínio (ex.: rubrica de speaking do PRD/SPEC-004).
- Não implementar nudges proativos; qualquer proatividade deve seguir `SPEC-011` e pode ficar para depois.

## 3) Assumptions (assunções)
- As tarefas do dia (de `SPEC-002`) incluem metadados suficientes para mapear para um gate (“gate profile”) quando aplicável.
- A política de privacidade do usuário (`PrivacyPolicy`, `SPEC-015`) está disponível no momento de solicitar/armazenar evidências.
- O sistema consegue separar **conteúdo sensível bruto** (C3) de **derivados não sensíveis** (C4) e aplicar retenção diferente.

## 4) Decisões técnicas (Decision log)
- **D-001 — Gate como artefato explícito por tarefa**
  - **Decisão**: cada `PlannedTask` pode carregar `gate_profile` (tipo + requisitos mínimos) e o sistema registra `GateResult` para aquela execução.
  - **Motivo**: satisfaz FR-002/FR-006 e permite consulta “passou/falhou” (FR-011) e métricas (`SPEC-016`).
  - **Alternativas consideradas**: gate implícito por domínio sem registro; descartado (não auditável).
  - **Impactos/Trade-offs**: exige disciplina para manter o gate “mínimo necessário” (NFR-001).

- **D-002 — Separar armazenamento de evidência bruta vs derivados**
  - **Decisão**: armazenar **Evidence** e aplicar políticas por categoria: C3 (curta/opt-out) vs derivados (resultado do gate, rubricas, contagens).
  - **Motivo**: `SPEC-015` + FR-012 (retenção curta de áudio, opção “não guardar”).
  - **Alternativas consideradas**: armazenar tudo no mesmo registro; descartado (piora privacidade e retenção).
  - **Impactos/Trade-offs**: mais entidades e regras; melhora confiança e conformidade.

- **D-003 — Estados de tarefa compatíveis com “tentativa”**
  - **Decisão**: quando evidência for parcial/ausente, registrar como `attempt` e manter `GateResult=not_satisfied` com `next_min_step`.
  - **Motivo**: edge cases do `SPEC-003` e `SPEC-002` (micro-consistência vs conclusão real).
  - **Alternativas consideradas**: contar parcial como concluído; descartado (falso progresso).
  - **Impactos/Trade-offs**: precisa comunicação clara em 1–2 frases (RNF1/RNF3).

- **D-004 — Política global de equivalência (MVP) centralizada**
  - **Decisão**: implementar uma política global de equivalência conforme FR-008, referenciada pelas SPECS de domínio:
    - Compreensão/retrieval: texto é evidência equivalente.
    - Speaking: áudio é padrão; sem áudio → `blocked/substituted` (não equivalente).
  - **Motivo**: consistência cross-cutting, evita decisões ad hoc.
  - **Alternativas consideradas**: cada domínio inventar equivalência; descartado (inconsistente).
  - **Impactos/Trade-offs**: speaking fica “bloqueável” em ambientes sem privacidade; é intencional (transparente e não punitivo).

## 5) Arquitetura (alto nível)
- **Componentes**
  - **Gate Engine**: avalia se uma tarefa pode concluir com base no `gate_profile` + evidências recebidas.
  - **Evidence Intake**: normaliza evidências (texto/áudio/metadado) e aplica `PrivacyPolicy` (guardar vs processar-e-descartar).
  - **Gate Result Store**: persiste `GateResult` e referências mínimas para auditoria.
  - **Domain Policy Registry**: catálogo de `gate_profile` por tipo de tarefa/domínio (ex.: inglês input, speaking, java prática, sono diário).
  - **Metrics Sink**: emite eventos/contagens para `SPEC-016` (sem conteúdo sensível).

- **Fluxos**
  - **Fluxo A — Tentativa de concluir tarefa**: usuário pede “feito” → se tarefa exige gate → status vira `evidence_pending` → sistema solicita evidência mínima → Gate Engine avalia → aceita (`completed`) ou rejeita (`attempt`/`evidence_pending`) com `next_min_step`.
  - **Fluxo B — Evidência inválida**: intake marca inválida → GateResult rejeitado com motivo claro + pedido de reenvio/alternativa (se equivalente).
  - **Fluxo C — Consulta**: “o que contou hoje?” → retorna lista com `GateResult` (satisfeito/não satisfeito) e motivo resumido.

## 6) Contratos e interfaces
- **Comando**: `RequestTaskCompletion(user_id, local_date, task_id, timestamp)`
  - **Saída**: `CompletionRequestResult`
    - `status`: `completed | evidence_required | already_completed`
    - `evidence_request?`: `EvidenceRequest(profile, items[], validity_rules, privacy_disclosure_short)`
  - **Erros**: `TASK_NOT_FOUND`, `TASK_NOT_COMPLETABLE` (ex.: bloqueada)

- **Comando**: `SubmitEvidence(user_id, task_id, evidence_payload, timestamp)`
  - **Saída**: `EvidenceReceipt(evidence_id, validity=valid|invalid, reason?, stored=kept|discarded)`
  - **Erros**: `EVIDENCE_TOO_LARGE`, `EVIDENCE_UNSUPPORTED`, `EVIDENCE_REJECTED_BY_POLICY` (ex.: modo mínimo proíbe guardar; ainda pode processar e descartar)

- **Comando**: `EvaluateGate(user_id, task_id, timestamp)`
  - **Saída**: `GateResultView`
    - `gate_status`: `satisfied | not_satisfied`
    - `reason_short` (1 frase)
    - `next_min_step` (1 passo)
    - `result_artifacts`: `rubric?`, `retrieval_summary?`, `error_log_ref?`

- **Consulta**: `GetGateStatus(user_id, task_id)`
  - **Saída**: `GateResultView` (inclui histórico mínimo: último resultado + timestamps)

- **Consulta**: `GetTodayGateSummary(user_id, local_date)` (`SPEC-016` / PRD §2)
  - **Saída**: lista curta de tarefas com `passed/failed/pending` e `next_min_step` para pendências.

## 7) Modelo de dados (mínimo)
- **QualityGateProfile**
  - `profile_id`, `domain`, `task_type` (ex.: `english_speaking`, `java_retrieval`, `sleep_diary`)
  - `task_class`: `learning | habit`
  - `requirements[]` (itens mínimos: ex.: “3 perguntas”, “rubrica completa”, “registro mínimo”)
  - `validity_rules[]` (ex.: “não vazio”, “áudio audível”, “rubrica 4 dimensões”)
  - `equivalence_policy` (referência à política global + overrides por domínio quando permitido)

- **Evidence**
  - `evidence_id`, `user_id`, `task_id`, `timestamp`
  - `kind`: `text_answer | rubric | audio | metadata`
  - `sensitivity`: `C2 | C3`
  - `storage_policy_applied`: `kept_7d | kept_custom | discarded_after_processing`
  - `content_ref?` (ponte para storage, se `kept`)
  - `summary` (curta; sem conteúdo sensível por padrão; ex.: “áudio 45s”, “3 respostas”)
  - **Retenção**: C3 default 7 dias; C2 conforme `SPEC-015` (moderada) e necessidades de consulta (`SPEC-016`).

- **GateResult**
  - `gate_result_id`, `user_id`, `task_id`, `timestamp`
  - `gate_status`: `satisfied | not_satisfied`
  - `failure_reason_code?`: `missing | invalid | partial | not_equivalent | other`
  - `reason_short`, `next_min_step`
  - `evidence_ids[]` (referências)
  - `derived_metrics` (ex.: `rubric_total`, `retrieval_ok=false`)
  - **Retenção**: moderada (ex.: 90 dias) e agregados semanais derivados (12 meses) via `SPEC-016`.

- **RubricScore** (derivado; sem conteúdo bruto)
  - `rubric_id`, `user_id`, `task_id`, `domain`
  - `dimensions` (ex.: 4 dims 0–2), `total`, `status: complete|partial`
  - **Retenção**: conforme C2/C4 (guardar derivados por mais tempo).

- **RecurringError** (cross-cutting com `SPEC-016`)
  - `error_id`, `user_id`, `domain`, `label`, `count_14d`, `count_week`, `last_seen_at`, `status: active|archived|target`
  - **Retenção**: enquanto útil (agregado); sem exemplos sensíveis por padrão.

- **ReinforcementAttempt** (para `SPEC-009`)
  - `attempt_id`, `user_id`, `domain`, `signal_ref`, `requested_at`, `done_at?`, `status`

## 8) Regras e defaults
- **Classificação de tarefa**: `learning` vs `habit` (FR-001).
- **Blocking rules (aprendizagem)**: sem evidência mínima válida → **não** conclui (FR-003).
- **Gates leves (hábito/fundação)**: registro mínimo + mínimo acordado; aceitar parcial quando apropriado (FR-004; alinhado a `SPEC-006`).
- **Equivalência global (MVP)** (FR-008):
  - Compreensão/retrieval: texto é equivalente.
  - Speaking: áudio é padrão; sem áudio → `blocked/substituted` e registrar como tentativa, não conclusão.
- **Micro-consistência**: permitido registrar 1 micro-step observável como `attempt` (≤60s), sem virar `completed` em aprendizagem.
- **Privacidade defaults** (FR-012; `SPEC-015`):
  - C3 bruto: 7 dias, com opção “não guardar” (processar e descartar).
  - Derivados não sensíveis: guardar para tendências (`SPEC-016`).
- **Anti-spam**: qualquer pedido de evidência deve ser 1 passo por vez; proatividade (lembretes) segue `SPEC-011`.

## 9) Observabilidade e métricas
- **Eventos**
  - `gate_required`, `evidence_requested`, `evidence_received(valid|invalid, stored|discarded)`
  - `gate_evaluated(satisfied|not_satisfied, reason_code)`
  - `task_blocked_by_policy` (ex.: speaking sem áudio)
  - `gate_friction_signal` (ex.: repetidas falhas por “missing evidence”)

- **Métricas (targets iniciais)**
  - % de tarefas de aprendizagem `completed` com `GateResult=satisfied` (target: ~100%; SC-001).
  - Distribuição de falhas por motivo (`missing/invalid/partial/not_equivalent`) para reduzir fricção (SC-005).
  - Latência média de avaliação de gate (operacional; suportar `SPEC-016` NFR-003).

## 10) Riscos & mitigação
- **Risco**: gates pesados demais → abandono.  
  **Mitigação**: “fricção proporcional”, 1 próximo passo mínimo, e medir fricção (SC-005).
- **Risco**: privacidade impede evidência e o usuário sente punição.  
  **Mitigação**: modo mínimo + explicação clara do trade-off; registrar bloqueio sem culpa.
- **Risco**: inconsistência entre domínios sobre equivalência.  
  **Mitigação**: política global central (D-004) e referência explícita em PLANs de domínio.

## 11) Rollout / migração
- **Feature flag**: `quality_gates_v1`.
- **Migração**: greenfield; futuras mudanças de `QualityGateProfile` devem ser versionadas (`profile_version`) para não reavaliar histórico.

## 12) Plano de testes (como validar)
- **Unit**
  - Avaliação determinística de gate por `QualityGateProfile` (missing/invalid/partial).
  - Política de equivalência: speaking sem áudio nunca “passa”.
  - Geração de `reason_short` + `next_min_step` (1 passo).
- **Integration**
  - Fluxo `done_request` → `evidence_pending` → `SubmitEvidence` → `EvaluateGate` → atualiza tarefa (via integração com `SPEC-002`).
  - Opt-out C3: “processar e descartar” funciona e GateResult mantém derivados.
- **E2E**
  - Tarefa de aprendizagem: tentar concluir sem evidência → bloqueia e pede mínimo; enviar evidência válida → conclui.
  - Evidência inválida (vazia/inaudível) → rejeita com pedido de reenvio.
- **Manual / acceptance**
  - Mensagens de bloqueio firmes e não punitivas (RNF3).
  - “O que você guarda?” responde consistente com `SPEC-015`.

## 13) Task breakdown (execução)
1) **Definir `QualityGateProfile` MVP e política de equivalência**
   - **Entrada**: `SPEC-003` FR-001..FR-010 + FR-008
   - **Saída**: catálogo mínimo de perfis (learning/habit) + regras de equivalência centralizadas
   - **Critério de pronto**: speaking requer áudio; compreensão/retrieval aceitam texto; perfis são “mínimos” e testáveis

2) **Modelar entidades `Evidence` e `GateResult` com retenção/sensibilidade**
   - **Entrada**: `SPEC-003` FR-012 + `SPEC-015` defaults
   - **Saída**: schema lógico + tabela de retenção por categoria (C2/C3/C4)
   - **Critério de pronto**: C3 tem 7 dias e opt-out; derivados preservam utilidade para `SPEC-016`

3) **Implementar “intake” de evidências com modo mínimo**
   - **Entrada**: `SPEC-015` opt-out + `SPEC-003` FR-007/FR-012
   - **Saída**: pipeline de validação básica (vazio/ilegível/inaudível) e decisão keep vs discard
   - **Critério de pronto**: evidência inválida não passa; opt-out não quebra fluxo (GateResult ainda existe)

4) **Implementar avaliação de gate (missing/invalid/partial) + `next_min_step`**
   - **Entrada**: `SPEC-003` cenários P1
   - **Saída**: função determinística que produz `GateResultView` curto e acionável
   - **Critério de pronto**: toda falha produz 1 motivo + 1 próximo passo mínimo; linguagem curta

5) **Integrar com `DailyState` (transição para `evidence_pending`/`completed`)**
   - **Entrada**: `SPEC-002` FR-008 + `SPEC-003` FR-006
   - **Saída**: contrato/handler que impede `completed` sem gate quando aplicável
   - **Critério de pronto**: tarefas de aprendizagem nunca contam como concluídas sem `GateResult=satisfied`

6) **Instrumentar eventos/métricas de fricção e qualidade**
   - **Entrada**: `SPEC-003` SC-001..SC-005 + `SPEC-016`
   - **Saída**: eventos `gate_*` e contadores por motivo de falha
   - **Critério de pronto**: possível medir “falhas por missing evidence” e fricção percebida proxy

## 14) Open questions (se existirem)
- (Default adotado) **Processar e descartar**: quando opt-out C3 estiver ativo, o sistema ainda pode validar evidência no momento e guardar apenas derivados (GateResult + rubric totals), desde que isso seja explicado no disclosure (`SPEC-015`).

