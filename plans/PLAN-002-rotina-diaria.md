# Technical Plan: PLAN-002 — Rotina Diária (Telegram-first) — Check-in + Plano A/B/C + Execução Guiada

**Created**: 2026-02-20  
**Spec**: `specs/SPEC-002-rotina-diaria-plano-abc.md`  
**PRD Base**: §§5.2, 5.3, 9.1, 10 (R1, R6), 11 (RNF1, RNF2, RNF3), §2 (consultar steps do dia), §6.2, §14  
**Related Specs**: `SPEC-001`, `SPEC-003`, `SPEC-010`, `SPEC-011`, `SPEC-015`, `SPEC-016`

## 1) Objetivo do plano
- Implementar o fluxo Telegram-first de **check-in mínimo** (tempo + energia) e seleção explícita de **Plano A/B/C** com justificativa curta.
- Implementar **geração de plano diário executável** (1 prioridade absoluta, 1–2 complementares, 1 fundação quando aplicável), com instruções objetivas e critério observável de “feito” por tarefa.
- Implementar **registro do dia atual**: check-in, plano e status por tarefa para consulta “meu plano de hoje” e “o que já fiz hoje”, com replanejamento e retomada após interrupções.

## 2) Non-goals (fora do escopo)
- Não implementar nudges proativos/anti-spam (budgets, timeouts e quiet hours) — isso é `SPEC-011`. Aqui cobre apenas comportamentos **quando o usuário interage**.
- Não implementar o motor completo de quality gates/evidência (validação e bloqueio de conclusão) — isso é `SPEC-003`. Aqui apenas:
  - carrega/expõe o “critério observável de feito”,
  - registra estados (incluindo “pendente de evidência”) e
  - integra por contratos com o gate engine.
- Não implementar a lógica completa de tarefas por domínio (Inglês/Java/Sono/...) — isso é coberto pelos PLANs das SPECS de domínio. Este PLAN define a **orquestração diária** e o contrato para que tarefas “pluguem” no plano.
- Não implementar dashboards/revisão semanal — isso é `SPEC-007`/`SPEC-016`.

## 3) Assumptions (assunções)
- Existe um `UserProfile` e preferências mínimas vindas do onboarding (`SPEC-001`); se o usuário ainda não fez onboarding, a rotina diária aplica defaults protetivos e oferece o “mínimo para destravar” (sem bloquear).
- “Dia atual” é definido por `timezone` do usuário. A chave do estado diário é `(user_id, local_date)`.
- O sistema mantém um catálogo de “templates de tarefas” por meta, mas no MVP ele pode começar simples (tarefas genéricas por domínio) e evoluir com `SPEC-004/005/006/...`.
- Privacidade/retensão/opt-out seguem `SPEC-015` e devem ser aplicadas ao que for registrado no estado diário (especialmente se alguma tarefa exigir evidência sensível no futuro via `SPEC-003`).

## 4) Decisões técnicas (Decision log)
- **D-000 — Baseline de plataforma**
  - **Decisão**: adotar o baseline `plans/PLAN-000-platform-baseline.md` como fonte de verdade para stack/arquitetura/padrões cross-cutting (Go backend + Next admin; Postgres/sqlc; Redis/worker; idempotência; privacidade C1–C5).
  - **Motivo**: garantir consistência entre rotina diária, jobs/timeouts (`SPEC-011`), métricas (`SPEC-016`) e privacidade (`SPEC-015`) sem duplicação.
  - **Alternativas consideradas**: escolher stack por feature; descartado.
  - **Impactos/Trade-offs**: ajustes de stack exigem atualização do baseline.

- **D-001 — “Estado do dia” como agregação persistida**
  - **Decisão**: persistir `DailyState` contendo `DailyCheckIn`, `DailyPlan` e lista de `PlannedTask` com status e timestamps.
  - **Motivo**: habilita consulta rápida de steps (PRD §2; `SPEC-002` FR-008/FR-009) e retomada.
  - **Alternativas consideradas**: recomputar plano sempre sem persistência; descartado por perder rastreabilidade e “o que já fiz hoje”.
  - **Impactos/Trade-offs**: precisa resolver idempotência e concorrência (mensagens repetidas).

- **D-002 — Seleção de plano A/B/C determinística por limiares default (MVP)**
  - **Decisão**: aplicar os limiares da SPEC como regra determinística (tempo/energia) para escolher A/B/C e registrar justificativa.
  - **Motivo**: previsibilidade, testabilidade e baixa fricção (RNF1) no MVP.
  - **Alternativas consideradas**: inferência “inteligente” por histórico; adiado (R7, `SPEC-010`).
  - **Impactos/Trade-offs**: pode parecer simplista no início; é aceitável e evoluível.

- **D-003 — Geração de plano via “PlanGenerator” + “TaskCatalog” plugável**
  - **Decisão**: separar orquestração (rotina diária) do conteúdo das tarefas (domínios). `TaskCatalog` fornece candidatos de tarefas; `PlanGenerator` compõe prioridade/complementares/fundação.
  - **Motivo**: evita acoplamento com SPECS de domínio e permite evolução incremental.
  - **Alternativas consideradas**: hardcode por domínio; descartado (explode manutenção).
  - **Impactos/Trade-offs**: exige contratos estáveis de `TaskTemplate`/`PlannedTask`.

- **D-004 — Estados de tarefa compatíveis com gates (sem implementá-los aqui)**
  - **Decisão**: usar status de tarefa: `planned | in_progress | completed | blocked | deferred | evidence_pending | attempt`.
  - **Motivo**: alinha com `SPEC-003` (concluído só com evidência) e edge cases do `SPEC-002` (micro-consistência, bloqueios por contexto).
  - **Alternativas consideradas**: apenas `todo/done`; descartado (insuficiente para qualidade e dias ruins).
  - **Impactos/Trade-offs**: mais estados; requer mensagens curtas para não confundir (RNF1).

## 5) Arquitetura (alto nível)
- **Componentes**
  - **Telegram Adapter**: normaliza mensagens para intents (check-in, replan, consult, task update).
  - **DailyRoutine Service**: coordena check-in, geração de plano, persistência do estado diário e consultas.
  - **PlanGenerator**: compõe Plano A/B/C e seleciona tarefas (prioridade/complementares/fundação) respeitando regras.
  - **TaskCatalog (interface)**: retorna `TaskTemplate`s para metas ativas e fundação (implementações futuras por domínio).
  - **GateIntegration (interface)**: registra “gate ref” e atualiza `evidence_pending/completed` quando `SPEC-003` validar evidência.
  - **Storage Layer**: repositórios de `DailyState` e event log (`SPEC-016`).
  - **Privacy Service**: aplica política do usuário para minimização/redação/retensão (`SPEC-015`).

- **Fluxos**
  - **Fluxo A — Check-in → Plano**: receber tempo+energia → escolher A/B/C → gerar tarefas → persistir `DailyState` → responder com plano em 1 mensagem.
  - **Fluxo B — Consultas**: “meu plano de hoje” / “o que já fiz hoje?” → ler `DailyState` → render curto.
  - **Fluxo C — Replanejamento**: receber novo tempo/energia/tempo restante → gerar plano ajustado (preservando status do que já foi feito) → persistir versão e responder mudança concisa.

## 6) Contratos e interfaces
Contratos de domínio (independentes de HTTP/DB) para suportar Telegram-first.

- **Comando**: `SubmitDailyCheckIn(user_id, local_date, time_available_min, energy_0_10, mood_stress_0_10?, constraints_text?, timestamp)`
  - **Saída**: `DailyPlanView(plan_type, rationale, priority_task, complementary_tasks[], foundation_task?, total_estimated_min)`
  - **Erros**:
    - `CHK_MISSING_TIME_OR_ENERGY` (aplicar 1 pergunta curta para destravar)
    - `CHK_AMBIGUOUS_ENERGY` (normalizar via mapeamento rápido; se ainda ambíguo, assumir 5)

- **Consulta**: `GetTodayPlan(user_id, local_date)`
  - **Saída**: `DailyPlanView` (curto; inclui status por tarefa)

- **Consulta**: `GetTodayStepsSummary(user_id, local_date)`
  - **Saída**: `DailyStepsSummary(done[], pending[], blocked[], in_progress[], note?)`

- **Comando**: `ReplanDay(user_id, local_date, new_time_available_min?, new_energy_0_10?, time_remaining_min?, timestamp)`
  - **Saída**: `DailyPlanView` + `PlanChangeExplanation(1-2 frases)`

- **Comando**: `UpdateTaskStatus(user_id, local_date, task_id, action, context?)`
  - `action`: `start | block(reason) | defer | mark_done_request | add_note`
  - **Saída**: `TaskStatusView(task_id, status, next_step?)`
  - **Integração gate**:
    - `mark_done_request` vira `evidence_pending` quando o task template exige gate (conforme `SPEC-003`), e só vai para `completed` após `GateResult=satisfied`.

## 7) Modelo de dados (mínimo)
- **DailyState**
  - `user_id`, `local_date` (chave composta)
  - `check_in_id?`, `plan_id?`
  - `tasks[]` (lista de `PlannedTask`)
  - `created_at`, `updated_at`
  - **Retenção**: default 90 dias (`SPEC-002` NFR-004; `SPEC-016` política de histórico; governado por `SPEC-015`).

- **DailyCheckIn** (`SPEC-002` Key Entities)
  - `check_in_id`, `user_id`, `local_date`
  - `time_available_min` (normalizado), `energy_0_10`
  - `mood_stress_0_10?`, `constraints_text?` (minimizar; opcional)
  - `created_at`
  - **Sensibilidade**: C1; `constraints_text` pode conter dados pessoais → tratar como potencialmente sensível e redigir/limitar.

- **DailyPlan**
  - `plan_id`, `user_id`, `local_date`
  - `plan_type`: `A | B | C`
  - `rationale` (curto; ex.: “Plano C porque tempo=10 e energia=2”)
  - `priority_task_id`, `complementary_task_ids[]`, `foundation_task_id?`
  - `version` (incrementa a cada replan), `created_at`

- **PlannedTask**
  - `task_id`, `user_id`, `local_date`
  - `title` (neutro; sem conteúdo sensível), `goal_domain` (english/java/sleep/health/self_esteem/saas)
  - `estimated_min`, `instructions` (curtas), `done_criteria` (curto)
  - `status`: `planned | in_progress | completed | blocked | deferred | evidence_pending | attempt`
  - `block_reason?` (curto), `note?` (curto; potencialmente sensível)
  - `gate_profile?` / `gate_ref?` (ponte para `SPEC-003`)
  - **Retenção/sensibilidade**: C1/C2; evitar armazenar evidência bruta aqui (isso fica em `SPEC-003` com política `SPEC-015`).

- **DomainEventLog** (`SPEC-016`)
  - `daily_check_in_submitted`, `daily_plan_generated`, `day_replanned`
  - `task_started`, `task_blocked`, `task_done_requested`, `task_completed`
  - **Privacidade**: payloads com títulos neutros e IDs; sem notas sensíveis por padrão.

## 8) Regras e defaults
- **Limiar A/B/C (default)** (`SPEC-002` FR-003):
  - Energia baixa: 0–3; média: 4–6; alta: 7–10.
  - Plano C (MVD): se `time_available_min ≤ 15` **OU** `energy_0_10 ≤ 3`.
  - Plano A: se `time_available_min ≥ 60` **E** `energy_0_10 ≥ 7`.
  - Plano B: demais casos.
  - Sempre registrar justificativa curta.

- **Estrutura do plano** (`SPEC-002` FR-004):
  - 1 prioridade absoluta
  - 1–2 complementares
  - 1 fundação quando aplicável (ex.: sono/saúde), respeitando metas ativas do ciclo.

- **Overload e limites**
  - Nunca agendar tarefas de mais de 2 metas intensivas no mesmo dia (alinhado a `SPEC-010` e `SPEC-002` FR-010).
  - Em Plano C: limite duro anti-overload (default): **máx 2 itens** (prioridade mínima + fundação mínima quando aplicável) (`SPEC-002` Edge Cases).

- **Micro-consistência (0 minutos)** (`SPEC-002` Edge Cases + `SPEC-003`)
  - Registrar como `attempt`/“MVD micro” apenas se houver 1 micro-step observável; não marcar aprendizagem como `completed` sem gate satisfeito.

- **Bloqueio por contexto**
  - Se restrição de ambiente impedir (ex.: sem privacidade), permitir `blocked` com `block_reason` neutro e oferecer alternativa mínima via TaskCatalog (sem equivalência plena para speaking; ver `SPEC-003`).

- **Privacidade/retensão** (`SPEC-015` + `SPEC-016`)
  - `DailyState` (C1) reter 90 dias por padrão; agregados (C4) 12 meses.
  - Não incluir conteúdo sensível em títulos/instruções; notas são opcionais e devem ser minimizadas/redigíveis.

## 9) Observabilidade e métricas
- **Logs/eventos mínimos**
  - `daily_check_in_submitted` (tempo, energia, plano escolhido)
  - `daily_plan_generated` (plan_type, counts, intensive_goals_count)
  - `day_replanned` (versão; motivo resumido)
  - `today_plan_viewed`, `today_steps_viewed`
  - `task_status_changed` (sem notas sensíveis)

- **Métricas (targets iniciais)**
  - **SC-001 (SPEC-002)**: tempo de interação até retornar o plano ≤ 2 min (proxy: número de mensagens até plano; duração entre primeira mensagem e `daily_plan_generated`).
  - **SC-003 (SPEC-002 + SPEC-016)**: taxa de consultas “plano/steps” e auto-relato “me perdi no dia” (capturado na revisão semanal) tende a cair.
  - Distribuição A/B/C por semana (para validar robustez a dias ruins e calibrar defaults).

## 10) Riscos & mitigação
- **Risco**: entradas ambíguas (energia “meh”, tempo não numérico) quebram fluxo.  
  **Mitigação**: normalização rápida + 1 pergunta curta no máximo; defaults conservadores (energia=5).
- **Risco**: plano longo demais em dia ruim aumenta culpa/abandono.  
  **Mitigação**: limite duro no Plano C; linguagem não punitiva (RNF2/RNF3).
- **Risco**: inconsistência por timezone (dia atual errado).  
  **Mitigação**: definir timezone no onboarding e permitir ajuste; logar “timezone missing”.
- **Risco**: tarefas marcadas como concluídas sem gate (falso progresso).  
  **Mitigação**: status `evidence_pending` + integração obrigatória com `SPEC-003` antes de virar `completed`.

## 11) Rollout / migração
- **Feature flag**: `daily_routine_v1`.
- **Migração**: greenfield; evoluir schema com campos opcionais (`version`, `gate_ref`).
- **Compatibilidade**: se `DailyState` não existir para o dia, gerar sob demanda no primeiro check-in.

## 12) Plano de testes (como validar)
- **Unit**
  - Seleção A/B/C pelos limiares (inclui edge cases).
  - Composição do plano: prioridade + complementares + fundação; limite de 2 intensivas.
  - Transições de status de tarefa (inclui `evidence_pending`).
  - Renderização curta de “plano de hoje” e “o que já fiz hoje”.
- **Integration**
  - Persistir `DailyState` e replanejar preservando histórico de versões e status já feitos.
  - Consultas em momentos diferentes do dia (antes/depois de status updates).
  - Aplicação de política de privacidade (redação de notas; retenção configurada) via `SPEC-015`.
- **E2E**
  - Check-in mínimo → plano → iniciar tarefa → marcar “feito” (vira `evidence_pending`) → (simulado) gate satisfeito → `completed`.
  - Retomada após “sumiço”: consulta “o que falta hoje?” retorna estado correto.
- **Manual / acceptance**
  - Linguagem não punitiva em Plano C e bloqueios.
  - Consulta do dia em ≤ 1 mensagem (quando possível), alinhado a `SPEC-016` AC-001.

## 13) Task breakdown (execução)
Tarefas pequenas (1–4h) ordenadas por dependências.

1) **Definir contratos do domínio de rotina diária**
   - **Entrada**: `SPEC-002` FR-001..FR-011; `SPEC-016` User Story 1
   - **Saída**: interfaces (commands/queries), erros e eventos mínimos
   - **Critério de pronto**: cobre check-in, replan, consultas e update de status sem depender de nudges

2) **Modelar `DailyState` + `DailyCheckIn` + `DailyPlan` + `PlannedTask`**
   - **Entrada**: Key Entities do `SPEC-002` + retenção de `SPEC-015`/`SPEC-016`
   - **Saída**: schema lógico + tabela de sensibilidade/retensão por campo
   - **Critério de pronto**: entidades suportam consulta de steps e integração com gates sem armazenar evidência sensível bruta

3) **Implementar seleção A/B/C e justificativa curta**
   - **Entrada**: `SPEC-002` FR-003 + edge cases de energia ambígua
   - **Saída**: função determinística `selectPlanType(time, energy)` + normalização
   - **Critério de pronto**: testes unit cobrem todos limiares e mapeamentos (“meh/ok/bem/péssimo”)

4) **Definir `TaskTemplate` e interface `TaskCatalog` (plugável)**
   - **Entrada**: `SPEC-002` FR-005 + integração futura com SPECS de domínio
   - **Saída**: contrato para templates (título, duração, instruções, done_criteria, gate_ref opcional)
   - **Critério de pronto**: `PlanGenerator` pode operar com um catálogo “mock” e gerar plano válido

5) **Implementar `PlanGenerator` (prioridade/complementares/fundação)**
   - **Entrada**: `SPEC-002` FR-004/FR-010 + limite Plano C (máx 2 itens)
   - **Saída**: geração de `DailyPlan` respeitando: 2 intensivas máximo e fundação quando aplicável
   - **Critério de pronto**: casos de teste cobrem (A, B, C) e overload; plano sempre “executável”

6) **Persistir `DailyState` com versionamento de replan**
   - **Entrada**: `SPEC-002` FR-006/FR-007/FR-008
   - **Saída**: repositório + estratégia de idempotência para check-in repetido
   - **Critério de pronto**: replan preserva tarefas concluídas e incrementa `DailyPlan.version`

7) **Implementar consultas “meu plano de hoje” e “o que já fiz hoje”**
   - **Entrada**: `SPEC-002` FR-009 + `SPEC-016` AC-001
   - **Saída**: renderização curta baseada em `DailyState`
   - **Critério de pronto**: retorna em ≤ 1 mensagem na maioria dos casos; diferencia `completed` vs `in_progress` vs `evidence_pending`

8) **Implementar comandos de status (start/block/defer/done_request)**
   - **Entrada**: `SPEC-002` FR-008 + `SPEC-003` (integração de gate)
   - **Saída**: atualizações de status e registro de eventos
   - **Critério de pronto**: `done_request` nunca vira `completed` sem gate quando `gate_ref` existir

9) **Integrar PrivacyPolicy na rotina diária (minimização/redação)**
   - **Entrada**: `SPEC-015` + `SPEC-002` NFR-004
   - **Saída**: regras para títulos neutros, notas opcionais, e comportamento em modo mínimo
   - **Critério de pronto**: opt-out C3 não quebra o plano; nenhuma resposta expõe conteúdo sensível por padrão

10) **Instrumentar eventos e métricas mínimas da rotina diária**
   - **Entrada**: `SPEC-016` targets (SC-001/SC-003) + `SPEC-002` SC-001
   - **Saída**: event log + contadores/latências básicas
   - **Critério de pronto**: possível medir tempo até plano e frequência de consultas sem leitura de conteúdo sensível

## 14) Open questions (se existirem)
- (Default adotado) **Sem onboarding prévio**: gerar um Plano C mínimo (check-in + 1 tarefa neutra de fundação + “próximo passo: fazer onboarding mínimo”), sem bloquear; registrar `onboarding_missing=true` para métricas.

