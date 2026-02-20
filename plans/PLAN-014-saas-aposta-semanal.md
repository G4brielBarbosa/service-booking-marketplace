# Technical Plan: PLAN-014 — SaaS: “aposta semanal” (bloco profundo) + microblocos

**Created**: 2026-02-20  
**Spec**: `specs/SPEC-014-saas-aposta-semanal.md`  
**PRD Base**: §8.6, §6.2, §5.2, §5.3, §9.1, §9.2, §14, §11 (RNF1–RNF3)  
**Related Specs**: `SPEC-002`, `SPEC-003`, `SPEC-007`, `SPEC-010`, `SPEC-011`, `SPEC-015`, `SPEC-016`, `SPEC-008`

## 1) Objetivo do plano
- Implementar definição semanal de **aposta** (foco único) com **definição de pronto** e resultado esperado observável.
- Implementar registro de execução do **bloco profundo** (30–120 min) e 1–2 **microblocos** opcionais, distinguindo `concluído` vs `tentativa`.
- Implementar downgrade em semana ruim (resultado mínimo viável) e integração com revisão semanal (`SPEC-007`) e métricas (`SPEC-016`) sem canibalizar fundamentos (`SPEC-010`).

## 2) Non-goals (fora do escopo)
- Não definir stack, repositório, board, issues, CI/CD nem ferramentas do SaaS.
- Não transformar SaaS em meta intensiva diária; permanece aposta semanal + microblocos.
- Não fazer backlog completo do SaaS; apenas 1 foco semanal (conecta com `SPEC-008` se necessário).
- Não enviar nudges proativos no MVP; se houver, seguir `SPEC-011`.

## 3) Assumptions (assunções)
- SaaS é “aposta semanal” (não conta como intensiva diária) conforme PRD §6.2.
- Evidência mínima do SaaS é **textual e curta** (1–3 bullets do resultado/entrega/decisão), evitando logs sensíveis.
- O sistema pode operar mesmo com semana ruim: permitir recorte do pronto e registrar “mínimo viável”.

## 4) Decisões técnicas (Decision log)
- **D-000 — Baseline de plataforma**
  - **Decisão**: adotar o baseline `plans/PLAN-000-platform-baseline.md` (event log/agregados, worker/jobs, retenção e redaction) como base.
  - **Motivo**: SaaS precisa registrar “pronto/resultado” com minimização e alimentar métricas semanais sem expor dados sensíveis.
  - **Alternativas consideradas**: registros ad hoc por conversa; descartado.
  - **Impactos/Trade-offs**: baseline define retenção default e padrões de storage.

- **D-001 — Definição de pronto como gate leve de “entrega observável”**
  - **Decisão**: tratar “pronto” como um `GateProfile` de domínio (não técnico): para concluir o bloco profundo, o usuário deve registrar `resultado_observado` compatível com o pronto.
  - **Motivo**: evitar ilusão de progresso (User Story 2) e alinhar com `SPEC-003` (sem burocracia).
  - **Alternativas consideradas**: contar “trabalhei um pouco”; descartado (falso progresso).
  - **Impactos/Trade-offs**: depende de descrição honesta; mitigado com templates de pronto.

- **D-002 — Downgrade explícito em semana ruim**
  - **Decisão**: manter um campo `minimal_outcome` para a aposta, ativado quando sinais de semana ruim existirem (energia baixa/overload) (`SPEC-010`).
  - **Motivo**: reduzir culpa e preservar consistência mínima.
  - **Alternativas consideradas**: “pular a semana”; aceita, mas preferir mínimo viável.
  - **Impactos/Trade-offs**: menos ambição; mais aderência.

- **D-003 — Microblocos opcionais com pronto simples**
  - **Decisão**: microblocos são opcionais e nunca são exigidos em semana ruim; cada microbloco tem pronto simples e evidência mínima de 1 frase.
  - **Motivo**: RNF1/RNF2 e evitar canibalizar fundação.
  - **Alternativas consideradas**: microblocos sempre; descartado.

## 5) Arquitetura (alto nível)
- **Componentes**
  - **SaaS Domain Service**: define aposta semanal, microblocos e captura de resultado.
  - **Gate Engine** (`SPEC-003`): gate leve para “bloco profundo concluído” e “microbloco concluído”.
  - **Metrics Aggregator** (`SPEC-016`): contagem de semanas com aposta concluída/tentativa e “semanas zero”.
  - **Governance/Overload** (`SPEC-010`): sinal para downgrade.
  - **Privacy Service** (`SPEC-015`): minimização e retenção do texto de resultado.

## 6) Contratos e interfaces
- **Comando**: `DefineSaasWeeklyBet(user_id, week_id, focus, expected_outcome, definition_of_done, deep_block_minutes, timestamp)`
  - **Saída**: `SaasBetView(status planned, minimal_outcome_suggestion?)`

- **Comando**: `DefineSaasMicroBlocks(user_id, week_id, blocks[], timestamp)`
  - **Saída**: `MicroBlocksView(count, optional=true)`

- **Comando**: `RecordSaasDeepWorkOutcome(user_id, week_id, result_observed, done_check yes|no, timestamp)`
  - **Saída**: `GateResultView(satisfied|not_satisfied, next_min_step)`

- **Comando**: `RecordSaasMicroBlockOutcome(user_id, micro_block_id, evidence_short, done_check yes|no, timestamp)`
  - **Saída**: `GateResultView`

- **Comando**: `DowngradeSaasBet(user_id, week_id, minimal_outcome, timestamp)`
  - **Saída**: `SaasBetView(updated)`

- **Consulta**: `GetSaasWeekSummary(user_id, week_id)`
  - **Saída**: status da aposta, resultado, microblocos, e 1 aprendizado curto.

## 7) Modelo de dados (mínimo)
- **SaasWeeklyBet**
  - `user_id`, `week_id`
  - `focus`, `expected_outcome`, `definition_of_done`
  - `deep_block_minutes`
  - `minimal_outcome?` (para semana ruim)
  - `status planned|completed|attempt|deferred`
  - `result_observed?` (1–3 bullets)
  - **Retenção**: C2/C4 moderada/longa; texto deve ser curto e redigível; permitir apagar por período (`SPEC-015`).

- **DeepWorkBlock**
  - `user_id`, `week_id`, `minutes_planned`, `minutes_done?`, `status`

- **MicroBlock**
  - `micro_block_id`, `user_id`, `week_id`
  - `definition_of_done_simple`, `evidence_short?`, `status`

- **SaasWeeklyAggregates** (`SPEC-016`)
  - `weeks_completed`, `weeks_attempted`, `weeks_zero`, `avg_deep_minutes?`

## 8) Regras e defaults
- 1 aposta semanal por semana; foco único; se foco grande demais, recortar para menor entrega possível.
- Gate de conclusão: `done_check=yes` + `result_observed` compatível com pronto; caso contrário `attempt` + next_min_step.
- Semana ruim: permitir downgrade para `minimal_outcome` (ex.: “decidir X”, “validar 1 hipótese”, “escrever 3 bullets”).
- Microblocos: opcionais (1–2); em semana ruim podem ser omitidos sem culpa.
- Privacidade: resultado deve ser neutro e sem detalhes sensíveis; permitir apagar.

## 9) Observabilidade e métricas
- **Eventos**: `saas_bet_defined`, `saas_bet_downgraded`, `saas_deep_block_outcome`, `saas_micro_block_outcome`
- **Métricas**: % de semanas com aposta concluída ou mínimo viável concluído; redução de semanas “zero” (SCs da SPEC).

## 10) Riscos & mitigação
- **Risco**: pronto vago gera falso progresso. Mitigar com templates e recortes.
- **Risco**: canibaliza fundamentos. Mitigar com governança e downgrade automático sugerido em overload.
- **Risco**: texto do resultado vira longo/sensível. Mitigar com limite de 1–3 bullets e modo mínimo.

## 11) Rollout / migração
- **Feature flag**: `saas_weekly_bet_v1`.

## 12) Plano de testes (como validar)
- **Unit**: recorte de pronto; gate de conclusão vs tentativa; downgrade em semana ruim.
- **Integration**: revisão semanal consome resumo do SaaS; métricas agregadas atualizam.
- **E2E**: semana normal vs semana ruim; microblocos opcionais.
- **Manual**: linguagem não punitiva; privacidade/retensão.

## 13) Task breakdown (execução)
1) Definir schema `SaasWeeklyBet/DeepWorkBlock/MicroBlock` + retenção  
   - **Entrada**: `SPEC-014` FR-001..FR-004 + `SPEC-015`  
   - **Saída**: modelo lógico e políticas  
   - **Critério de pronto**: resultado é curto, redigível e apagável

2) Implementar criação da aposta semanal com pronto e recorte  
   - **Entrada**: User Story 1  
   - **Saída**: `DefineSaasWeeklyBet` + templates de pronto  
   - **Critério de pronto**: foco grande é recortado; 1 resultado esperado

3) Implementar registro do bloco profundo (concluído vs tentativa) + next_min_step  
   - **Entrada**: User Story 2  
   - **Saída**: gate leve e status correto  
   - **Critério de pronto**: sem pronto satisfeito não marca concluído

4) Implementar downgrade para mínimo viável em semana ruim  
   - **Entrada**: User Story 4 + `SPEC-010`  
   - **Saída**: `DowngradeSaasBet`  
   - **Critério de pronto**: reduz escopo sem culpa e preserva consistência mínima

5) Implementar microblocos opcionais e registro simples  
   - **Entrada**: User Story 3  
   - **Saída**: microblocos com pronto simples + evidência mínima  
   - **Critério de pronto**: omitidos em semana ruim sem penalizar

6) Instrumentar eventos e agregados semanais  
   - **Entrada**: SCs + `SPEC-016`  
   - **Saída**: métricas de semanas concluídas/tentativa/zero  
   - **Critério de pronto**: revisão semanal mostra progresso do SaaS sem relatório longo

## 14) Open questions (se existirem)
- (Default adotado) **Retenção do texto `result_observed`**: tratar como C2 moderado (90 dias) e manter apenas contagens/estados (C4) por 12 meses, para minimizar dados sobre o trabalho.

