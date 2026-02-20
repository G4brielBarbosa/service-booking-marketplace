# Technical Plan: PLAN-008 — Backlog Inteligente e Priorização baseada em lacunas observadas

**Created**: 2026-02-20  
**Spec**: `specs/SPEC-008-backlog-inteligente.md`  
**PRD Base**: §5.2, §5.5, §9.1, §10 (R4), §6.2, §11 (RNF1–RNF3), §13  
**Related Specs**: `SPEC-002`, `SPEC-003`, `SPEC-009`, `SPEC-010`, `SPEC-011`, `SPEC-015`, `SPEC-016`, `SPEC-007`

## 1) Objetivo do plano
- Implementar um **backlog de itens derivados de lacunas observadas** (erros recorrentes, queda de rubrica, gates falhando, baixa consistência, sinais de overload), com explicação curta “por que agora”.
- Implementar **priorização e limite de itens ativos** para evitar overload (shortlist acionável).
- Permitir ações do usuário: **aceitar / adiar / rejeitar / consultar backlog**, registrando decisões e reduzindo repetição excessiva.

## 2) Non-goals (fora do escopo)
- Não otimizar priorização com ML; MVP usa regras determinísticas e observáveis.
- Não virar gerenciador de tarefas genérico; backlog é para progresso das metas do PRD.
- Não implementar proatividade de sugestões via nudges; recomendações são mostradas **quando o usuário interagir** (proatividade segue `SPEC-011`).

## 3) Assumptions (assunções)
- Sinais são derivados de métricas/registro (`SPEC-016`) e resultados de gates (`SPEC-003`).
- Governança/limites de metas intensivas (`SPEC-010`) está disponível para que backlog não proponha expansão de escopo indevida.
- O backlog alimenta `SPEC-002` (planejamento diário) como **candidatos opcionais**, não como obrigações automáticas no MVP.

## 4) Decisões técnicas (Decision log)
- **D-001 — BacklogItem como entidade de domínio, com “evidence_ref”**
  - **Decisão**: cada item guarda: origem/sinal, evidência resumida (sem conteúdo sensível), critério observável de feito e prioridade.
  - **Motivo**: transparência (“por que agora”) e auditabilidade (`SPEC-008` FR-001/FR-006).
  - **Alternativas consideradas**: apenas texto livre; descartado (não rastreável).
  - **Impactos/Trade-offs**: exige modelagem e uma camada de geração de sinais.

- **D-002 — Limite de itens ativos com política simples**
  - **Decisão**: manter uma fila “não ativa” e uma shortlist “ativa”. Default MVP:
    - `active_limit = 5` (faixa recomendada 3–7).
    - no máximo 2 itens ativos para metas intensivas ao mesmo tempo.
  - **Motivo**: evitar overload e spam cognitivo (PRD §6.2; `SPEC-008` FR-004).
  - **Alternativas consideradas**: mostrar tudo; descartado.
  - **Impactos/Trade-offs**: itens podem demorar a aparecer; aceitável.

- **D-003 — Cooldown por rejeição/adiamento**
  - **Decisão**: quando o usuário adia/rejeita, aplicar cooldown (default: 7 dias) antes de reativar item similar.
  - **Motivo**: reduzir repetição excessiva (FR-005) e respeitar segurança psicológica.
  - **Alternativas consideradas**: re-sugerir sempre; descartado.
  - **Impactos/Trade-offs**: pode atrasar intervenção; mitigado por sinais fortes (severidade alta).

## 5) Arquitetura (alto nível)
- **Componentes**
  - **Signal Extractor**: transforma métricas/estados em `BacklogSignal`s (regras MVP).
  - **Backlog Generator**: cria/atualiza `BacklogItem`s com critério de feito alinhado a gates (`SPEC-003`).
  - **Prioritizer**: seleciona itens ativos (limites, severidade, recência, governança).
  - **Backlog Store**: persistência dos itens e decisões do usuário.
  - **Planner Integration**: fornece “itens aceitos/ativos” como candidatos para o plano diário (`SPEC-002`) e revisão semanal (`SPEC-007`).

## 6) Contratos e interfaces
- **Comando**: `RefreshBacklog(user_id, timestamp)` (user-initiated no MVP)
  - **Saída**: `BacklogSnapshot(active_items[], inactive_count, explanation_short?)`

- **Consulta**: `GetBacklog(user_id)`
  - **Saída**: `BacklogSnapshot` (ativos primeiro; cada item com “por que agora” + critério de feito)

- **Comando**: `ActOnBacklogItem(user_id, item_id, action accept|defer|reject, reason_short?, timestamp)`
  - **Saída**: `BacklogActionReceipt(next_suggestion_short)`

- **Consulta**: `GetBacklogSignals(user_id, limit?)` (debug/observabilidade; opcional)
  - **Saída**: lista curta de sinais e severidade (sem conteúdo sensível)

## 7) Modelo de dados (mínimo)
- **BacklogSignal**
  - `signal_id`, `user_id`
  - `type`: `recurring_error | rubric_drop | low_consistency | gate_failing | low_energy | overload`
  - `severity`: `low | medium | high`
  - `window`: ex.: `last_7d`, `last_14d`, `last_5_sessions`
  - `evidence_refs[]` (IDs de agregados/gates; sem C3)
  - `created_at`, `updated_at`
  - **Retenção**: agregados (C4) 12 meses; sem conteúdo sensível.

- **BacklogItem**
  - `item_id`, `user_id`
  - `title_short`
  - `domain` (english/java/sleep/health/self_esteem/saas)
  - `origin_signal_id`
  - `why_now_short` (1 frase)
  - `definition_of_done` (critério observável; se envolver evidência, referenciar `gate_profile_id`)
  - `priority_score` (derivado), `status`: `active | accepted | deferred | rejected | completed | archived`
  - `cooldown_until?`, `created_at`, `updated_at`
  - **Privacidade**: texto neutro; sem detalhes sensíveis.

- **BacklogDecision**
  - `decision_id`, `user_id`, `item_id`
  - `action`, `reason_short?`, `timestamp`

## 8) Regras e defaults
- **Regras de geração de sinais (MVP)**
  - `recurring_error`: erro recorrente (≥3/14d) ativo em inglês/java (`SPEC-016`).
  - `rubric_drop`: queda ≥1 ponto vs semana anterior (quando dados suficientes) (`SPEC-009`/`SPEC-016`).
  - `low_consistency`: ≤2 dias/7 em meta intensiva (janela semanal) (`SPEC-010`).
  - `gate_failing`: 3+ falhas por “missing evidence” na semana (fricção) (`SPEC-003` SC-005).
  - `overload`: combinação conforme `SPEC-010` (consistência baixa + energia baixa + rubrica em queda).

- **Priorização (MVP)**
  - ordenar por `severity` → recência → impacto estimado.
  - nunca ativar itens que expandam escopo em múltiplas metas intensivas simultaneamente (governança).
  - se `overload=high`, priorizar itens de **redução de escopo/MVD** em vez de “mais coisas”.

- **Limites**
  - `active_limit=5`, `accepted_limit=3` (itens aceitos “planejáveis”).
  - cooldown default: 7 dias para `defer/reject` (exceto severidade high).

- **Privacidade/retensão** (`SPEC-015`)
  - itens guardam apenas referências e agregados; nunca guardar conteúdo sensível bruto (C3).

## 9) Observabilidade e métricas
- **Eventos**
  - `backlog_refreshed`, `signal_detected(type,severity)`
  - `backlog_item_created`, `backlog_item_activated`, `backlog_item_actioned(action)`
  - `backlog_item_completed` (quando critério de feito satisfaz gate, se aplicável)

- **Métricas**
  - Taxa de itens aceitos → concluídos (com gate quando aplicável).
  - Tamanho médio do backlog ativo (deve permanecer ≤ limite).
  - Razões de rejeição/adiamento (para ajustar regras).

## 10) Riscos & mitigação
- **Risco**: backlog vira lista infinita.  
  **Mitigação**: limites + inativos + arquivamento automático de baixo impacto.
- **Risco**: recomendações parecem aleatórias.  
  **Mitigação**: `why_now_short` + `evidence_refs` (auditável).
- **Risco**: backlog incentiva overload.  
  **Mitigação**: regras de governança e foco em redução quando sinais de overload surgem.

## 11) Rollout / migração
- **Feature flag**: `backlog_v1`.
- Backfill: opcional — gerar itens iniciais a partir de semana atual/últimos 14 dias quando `SPEC-016` estiver disponível.

## 12) Plano de testes (como validar)
- **Unit**
  - Regras de sinal (limiares, janelas) e priorização.
  - Limites de itens ativos e cooldown.
- **Integration**
  - Dado um conjunto de agregados (`SPEC-016`), `RefreshBacklog` cria itens com “por que agora”.
  - `ActOnBacklogItem` atualiza status e impede re-sugestão imediata.
- **E2E**
  - Erro recorrente em Java → item criado → usuário aceita → item aparece como candidato no plano diário (sem auto-enfileirar).
- **Manual / acceptance**
  - Linguagem não punitiva; backlog curto e acionável.

## 13) Task breakdown (execução)
1) **Definir schema `BacklogSignal/BacklogItem/BacklogDecision` + retenção**
   - **Entrada**: `SPEC-008` Key Entities + `SPEC-015/016`
   - **Saída**: modelo lógico e política de minimização
   - **Critério de pronto**: nenhum campo requer C3; itens são neutros e redigíveis

2) **Implementar regras MVP de extração de sinais**
   - **Entrada**: `SPEC-008` FR-001 + `SPEC-016` agregados
   - **Saída**: `Signal Extractor` com 5–6 sinais básicos
   - **Critério de pronto**: sinais são determinísticos e testáveis (janelas/limiares explícitos)

3) **Implementar geração/atualização de itens com `why_now_short`**
   - **Entrada**: `SPEC-008` FR-002/FR-006
   - **Saída**: itens com critério de feito alinhado a gates (`SPEC-003`)
   - **Critério de pronto**: todo item tem 1 frase “por que agora” e 1 critério observável

4) **Implementar priorização + limites + cooldown**
   - **Entrada**: `SPEC-008` FR-004/FR-005 + `SPEC-010`
   - **Saída**: shortlist ativa (≤5) e regras de recirculação
   - **Critério de pronto**: backlog ativo nunca excede limite; rejeitados não voltam antes do cooldown

5) **Implementar APIs de consulta e ação do usuário**
   - **Entrada**: `SPEC-008` FR-005/FR-006
   - **Saída**: `GetBacklog` + `ActOnBacklogItem`
   - **Critério de pronto**: aceitar/adiar/rejeitar funciona e é refletido na próxima consulta

6) **Integrar com planejamento diário e revisão semanal (read-only)**
   - **Entrada**: `SPEC-002` e `SPEC-007`
   - **Saída**: “itens aceitos” elegíveis como candidatos (não obrigatórios)
   - **Critério de pronto**: plano do dia pode incluir 0–1 item de backlog sem aumentar overload

## 14) Open questions (se existirem)
- (Default adotado) **Backlog é user-initiated**: não envia sugestões proativas no MVP; proatividade só depois de `SPEC-011` com budgets/opt-out.

