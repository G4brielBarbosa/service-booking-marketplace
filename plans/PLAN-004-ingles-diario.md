# Technical Plan: PLAN-004 — Inglês Diário: Input + Output + Retrieval (rubrica + erros recorrentes)

**Created**: 2026-02-20  
**Spec**: `specs/SPEC-004-ingles-diario.md`  
**PRD Base**: §8.1, §5.3, §5.4, §5.2, §9.1, §14, §10 (R2, R3, R6), §11 (RNF1–RNF4)  
**Related Specs**: `SPEC-002`, `SPEC-003`, `SPEC-009`, `SPEC-011`, `SPEC-015`, `SPEC-016`

## 1) Objetivo do plano
- Implementar o loop diário de inglês como **3 blocos** (input, output/speaking, retrieval) com versões mínimas para dia ruim (Plano C).
- Integrar com **quality gates** (`SPEC-003`): cada bloco só “conta” quando evidência mínima válida é registrada (ou fica explicitamente como tentativa/bloqueado).
- Registrar rubrica de speaking, retrieval e erros/aprendizados de forma minimizada para alimentar tendências e revisão semanal (`SPEC-016`).

## 2) Non-goals (fora do escopo)
- Não escolher plataforma/curso/conteúdo específico; apenas registrar “descrição de alto nível” do input.
- Não implementar transcrição, avaliação automática de pronúncia, scoring por IA ou análise de áudio.
- Não implementar nudges proativos; qualquer lembrete segue `SPEC-011`.
- Não implementar backlog inteligente/seleção ótima de tarefas (isso é `SPEC-008/009`).

## 3) Assumptions (assunções)
- A rotina diária (`SPEC-002`) chamará este domínio via `TaskCatalog`/templates. Este PLAN define o contrato do “conteúdo” dos templates e como registrar evidência.
- A rubrica default do PRD é usada no MVP: 4 dimensões (0–2) total 0–8.
- Privacidade: speaking áudio é C3 por padrão; o usuário pode optar por “não guardar” (processar e descartar) conforme `SPEC-015`.

## 4) Decisões técnicas (Decision log)
- **D-000 — Baseline de plataforma**
  - **Decisão**: adotar o baseline `plans/PLAN-000-platform-baseline.md` (Go backend + Next admin; Postgres/sqlc; Redis/worker; privacidade/retensão) como base de execução e armazenamento.
  - **Motivo**: inglês depende de gates (`SPEC-003`), métricas (`SPEC-016`) e política de C3 para áudio (`SPEC-015`) — padronização evita inconsistência.
  - **Alternativas consideradas**: decidir stack por domínio; descartado.
  - **Impactos/Trade-offs**: baseline central vira referência obrigatória.

- **D-001 — Modelar o loop como três tarefas separadas (com gates próprios)**
  - **Decisão**: no `DailyPlan`, o inglês aparece como 2–3 `PlannedTask`s (input, speaking, retrieval) dependendo do plano A/B/C.
  - **Motivo**: cada bloco tem evidência diferente e permite “execução parcial” sem mentir sobre speaking (alinhado a `SPEC-003` e edge cases).
  - **Alternativas consideradas**: uma tarefa única “Inglês hoje”; descartado (perde granularidade de gate e métricas).
  - **Impactos/Trade-offs**: mais itens no plano; mitigação: no Plano C, limitar e compactar.

- **D-002 — Speaking sem equivalência plena quando áudio não é possível**
  - **Decisão**: speaking exige áudio para `GateResult=satisfied`; sem áudio, registrar `blocked`/`substituted_attempt` e oferecer alternativa de consistência (retrieval curto), sem marcar speaking como concluído.
  - **Motivo**: política global de equivalência do `SPEC-003` (speaking é objetivo de produção oral).
  - **Alternativas consideradas**: aceitar texto como speaking; descartado.
  - **Impactos/Trade-offs**: pode reduzir “dias contados” em speaking; aumenta honestidade de progresso.

- **D-003 — Minimização de dados: armazenar derivados + referências, não conteúdo**
  - **Decisão**: para input/retrieval, armazenar respostas curtas e agregados; para áudio, armazenar por 7 dias (default) ou descartar após validar, mantendo rubrica/resultado.
  - **Motivo**: `SPEC-015` + `SPEC-016` (tendências sem guardar C3).
  - **Alternativas consideradas**: guardar tudo por longos períodos; descartado.
  - **Impactos/Trade-offs**: menos auditabilidade textual; aceitável no MVP.

## 5) Arquitetura (alto nível)
- **Componentes**
  - **English Domain Service**: define templates (A/B/C) e processa submissões (checagem de compreensão, rubrica, retrieval, erro do dia).
  - **Gate Engine** (`SPEC-003`): avalia gates por bloco e produz `GateResult`.
  - **Metrics/Errors Store** (`SPEC-016`): agrega rubricas, consistência, erros recorrentes.
  - **Privacy Service** (`SPEC-015`): aplica opt-out/retensão e disclosure no pedido de áudio.

- **Fluxos**
  - **Input**: registrar descrição curta do input + responder 3 perguntas (ou 1 no Plano C) → gate satisfeito.
  - **Speaking**: enviar áudio + preencher rubrica → gate satisfeito; se áudio impossível → marcar bloqueio e oferecer alternativa mínima (não equivalente).
  - **Retrieval**: responder 5–10 itens (ou 3 no Plano C) → registrar desempenho simples.

## 6) Contratos e interfaces
- **Comando**: `StartEnglishTasks(user_id, local_date, plan_type)`
  - **Saída**: `TaskTemplates[]` para input/speaking/retrieval (com duração, instruções e `gate_profile_id` por bloco)

- **Comando**: `SubmitEnglishInputCheck(user_id, task_id, answers[], timestamp)`
  - **Saída**: `GateResultView` (satisfeito/não) + `next_min_step`

- **Comando**: `SubmitSpeakingAudio(user_id, task_id, audio_payload, timestamp)`
  - **Saída**: `EvidenceReceipt` + `GateStatus(evidence_pending|ready_to_evaluate)`

- **Comando**: `SubmitSpeakingRubric(user_id, task_id, rubric_dimensions, timestamp)`
  - **Saída**: `GateResultView` (ou `evidence_missing` se áudio ausente)

- **Comando**: `SubmitEnglishRetrieval(user_id, task_id, items_answered[], timestamp)`
  - **Saída**: `RetrievalSummary(status=ok|low, targets[])` + `GateResultView` (quando aplicável)

- **Comando**: `LogEnglishErrorOfDay(user_id, local_date, label, note_short?, timestamp)`
  - **Saída**: `ErrorLogRef` + `RecurringStatus(active|recurring|target)`

- **Consulta**: `GetEnglishWeekTrend(user_id, week_id)`
  - **Saída**: `EnglishWeeklyTrend` (consistência, média rubrica, top erros recorrentes) — alimenta `SPEC-007`.

## 7) Modelo de dados (mínimo)
- **EnglishInputSession**
  - `user_id`, `local_date`, `task_id`
  - `duration_est_min`, `content_descriptor` (alto nível)
  - `comprehension_answers[]` (curtas), `status: complete|partial`
  - **Retenção**: C2 moderada (default 90 dias), podendo reter só agregados (C4) quando em modo mínimo.

- **SpeakingEvidence**
  - `evidence_id`, `task_id`, `length_sec`, `stored=kept|discarded`, `retention_days`
  - `content_ref?` (se kept)
  - **Sensibilidade/retensão**: C3; default 7 dias; opt-out “não guardar” → `discarded`.

- **SpeakingRubric**
  - `user_id`, `task_id`, `dimensions(4x0-2)`, `total`, `status complete|partial`
  - **Retenção**: derivados (C4) 12 meses; vínculo com tarefa (C2) 90 dias.

- **EnglishRetrieval**
  - `user_id`, `task_id`, `count_items`, `status ok|low`, `targets[]`
  - **Retenção**: agregados (C4) 12 meses; detalhes mínimos 90 dias.

- **EnglishErrorLogEntry**
  - `user_id`, `local_date`, `label`, `note_short?`, `recurring_group_id?`
  - **Privacidade**: `note_short` opcional e minimizada; redigível; evitar conteúdo sensível.

## 8) Regras e defaults
- **Plano normal (A/B)** (MVP): input 10–30 min + 3 perguntas; speaking 30–180s + rubrica completa; retrieval 5–10 itens.
- **Plano C (dia ruim)**: input 5–10 min + 1 pergunta; speaking 30–60s (se possível) + rubrica mínima; retrieval 3 itens.
- **Gates** (`SPEC-003`):
  - Input: checagem mínima respondida (não vazia) e coerente o suficiente; se falhar, não conclui.
  - Speaking: áudio + rubrica (mínimo definido) → conclui; sem áudio, não conclui.
  - Retrieval: registrar respostas e status; pode concluir com “ok/low” (qualidade entra em métricas).
- **Erros recorrentes** (`SPEC-016` default): recorrente = ≥3 ocorrências em 14 dias (ou ≥3 na semana).
- **Privacidade** (`SPEC-015`): disclosure ao pedir áudio; respeitar opt-out e modo mínimo.
- **Anti-spam**: qualquer follow-up para completar rubrica/evidência só ocorre quando usuário interagir, a menos que `SPEC-011` habilite nudges.

## 9) Observabilidade e métricas
- **Eventos**
  - `english_input_completed(gate_status)`, `english_speaking_completed(gate_status)`, `english_retrieval_completed`
  - `speaking_audio_stored|discarded`, `english_error_logged(label)`
  - `english_gate_failed(reason_code)`

- **Métricas (SCs principais)**
  - Consistência semanal de inglês com gates satisfeitos (`SPEC-016`).
  - Tendência da rubrica (média semanal; comparação com semana anterior).
  - Frequência de erros recorrentes e redução após alvos/reforços (`SPEC-009`).

## 10) Riscos & mitigação
- **Risco**: speaking bloqueado por privacidade/ambiente reduz adesão.  
  **Mitigação**: alternativa mínima não equivalente + transparência; focar em input+retrieval em dias bloqueados; registrar bloqueio para revisão semanal.
- **Risco**: fricção de rubrica/checagem vira burocracia.  
  **Mitigação**: rubrica mínima no Plano C; 1 próximo passo mínimo quando incompleta; medir falhas por “missing”.

## 11) Rollout / migração
- **Feature flag**: `english_daily_v1`.
- **Migração**: greenfield; evoluir campos como opcionais (ex.: `content_descriptor` mais estruturado no futuro).

## 12) Plano de testes (como validar)
- **Unit**
  - Templates A/B/C geram durações e requisitos coerentes.
  - Speaking sem áudio nunca “passa” (equivalência).
  - Recorrência de erro (≥3/14d) é marcada corretamente.
- **Integration**
  - Fluxo por bloco com Gate Engine: input → gate; speaking áudio+rubrica → gate; retrieval → registro.
  - Opt-out C3: áudio é descartado e ainda assim rubrica/resultado persistem.
- **E2E**
  - Um dia bom: concluir 3 blocos com evidência mínima.
  - Dia ruim: concluir versão mínima; speaking bloqueado → registrar substituição sem contar como speaking.
- **Manual / acceptance**
  - Mensagens curtas e não punitivas em falhas de gate e bloqueio por contexto.

## 13) Task breakdown (execução)
1) **Definir templates A/B/C de inglês (TaskCatalog) com gate profiles**
   - **Entrada**: `SPEC-004` FR-001..FR-005 + `SPEC-003`
   - **Saída**: `TaskTemplate`s (input/speaking/retrieval) com instruções e `gate_profile_id`
   - **Critério de pronto**: Plano C cabe em 5–15 min; speaking tem gate explícito

2) **Modelar entidades mínimas e retenção (inclui áudio C3)**
   - **Entrada**: `SPEC-015` + defaults de `SPEC-003/016`
   - **Saída**: schema lógico + políticas keep/discard para áudio
   - **Critério de pronto**: opt-out “não guardar áudio” funciona mantendo derivados

3) **Implementar submissão de checagem de compreensão + avaliação de gate**
   - **Entrada**: `SPEC-004` Scenario Input + `SPEC-003`
   - **Saída**: comando `SubmitEnglishInputCheck` retornando `GateResultView`
   - **Critério de pronto**: falha gera 1 next_min_step (“responda 1 pergunta…” ou “reduza dificuldade”)

4) **Implementar fluxo speaking: evidência de áudio + rubrica**
   - **Entrada**: `SPEC-004` Scenario Output + `SPEC-003/015`
   - **Saída**: intake de áudio (keep/discard) + registro de rubrica + GateResult
   - **Critério de pronto**: sem áudio → not_satisfied; com áudio+rubrica → satisfied

5) **Implementar retrieval e registro de desempenho/targets**
   - **Entrada**: `SPEC-004` FR-003 + `SPEC-016` (targets/erros)
   - **Saída**: `EnglishRetrieval` + status ok/low + itens alvo
   - **Critério de pronto**: Plano C usa 3 itens; dados alimentam tendência semanal

6) **Implementar registro de erro do dia e recorrência**
   - **Entrada**: `SPEC-004` FR-006 + `SPEC-016` limiar
   - **Saída**: `EnglishErrorLogEntry` + atualização de `RecurringError`
   - **Critério de pronto**: após 3 ocorrências em 14 dias vira recorrente; consultável

7) **Instrumentar eventos/métricas de inglês**
   - **Entrada**: `SPEC-004` SCs + `SPEC-016`
   - **Saída**: eventos `english_*` e agregados semanais
   - **Critério de pronto**: revisão semanal consegue ler média de rubrica e consistência

## 14) Open questions (se existirem)
- (Default adotado) **Critério mínimo de rubrica em dia ruim**: permitir rubrica parcial (ao menos 2 dimensões) para registrar tentativa e feedback, mas manter gate de speaking como “não satisfeito” se áudio faltar; quando áudio existir, aceitar rubrica parcial como evidência mínima apenas no Plano C (registrar como `partial`).

