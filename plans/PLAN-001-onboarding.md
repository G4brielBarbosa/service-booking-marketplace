# Technical Plan: PLAN-001 — Onboarding e Diagnóstico Leve (1–2 semanas)

**Created**: 2026-02-20  
**Spec**: `specs/SPEC-001-onboarding-diagnostico.md`  
**PRD Base**: §5.1, §§5.2–5.3, §6.2, §9.1, §10 (R1, R6, R7), §11 (RNF1–RNF4), §14  
**Related Specs**: `SPEC-002`, `SPEC-003`, `SPEC-010`, `SPEC-011`, `SPEC-015`, `SPEC-016`

## 1) Objetivo do plano
- Implementar o **onboarding mínimo** para destravar a rotina diária: metas ativas do ciclo (com limite de 2 intensivas), restrições mínimas, baseline mínima relevante, MVD, e **resumo consultável**.
- Implementar **diagnóstico leve progressivo** (1–2 semanas) para completar lacunas por domínio em passos únicos, sem burocracia.
- Implementar **retomada e revisão** do onboarding (sem “resetar tudo”), com política “vale daqui pra frente”.

## 2) Non-goals (fora do escopo)
- Não implementar geração completa de plano diário nem execução guiada (`SPEC-002`) além do “pronto para começar”.
- Não implementar o motor completo de **quality gates/evidência** (`SPEC-003`) — apenas armazenar referências e estados necessários e garantir que o onboarding não “finja evidência”.
- Não implementar nudges proativos/anti-spam (budgets, timeouts, quiet hours) além de armazenar configurações iniciais (`SPEC-011`).
- Não implementar dashboards/relatórios avançados — apenas registrar eventos e métricas mínimas necessárias (`SPEC-016`).
- Não definir/implementar escolhas de conteúdo para inglês/java, plataformas, STT, scoring automático, etc. (ver SPECS de domínio).

## 3) Assumptions (assunções)
- **Identidade**: um usuário é identificado de forma estável por `telegram_user_id` (e opcionalmente `chat_id`). MVP assume 1 “perfil” por usuário, com suporte natural a multiusuário.
- **Fuso horário**: derivado do onboarding (pergunta simples) ou default para o timezone do servidor se ausente; usado para “dia atual”, semanas e quiet hours.
- **Metas suportadas**: conjunto fechado do PRD (Inglês, Java, Sono, Vida saudável, Autoestima, SaaS).
- **Baseline mínima** é “suficiente para operar” (como na SPEC): se áudio não for possível, registrar substitutos e marcar speaking baseline como pendente (alinhado a `SPEC-003` + `SPEC-015`).
- **Plataforma/stack**: seguir o baseline em `plans/PLAN-000-platform-baseline.md` (Go backend + Next admin + Postgres/sqlc + Redis + Docker Compose; idempotência, event log, retenção C1–C5).

## 4) Decisões técnicas (Decision log)
- **D-000 — Baseline de plataforma**
  - **Decisão**: adotar o baseline `plans/PLAN-000-platform-baseline.md` como fonte de verdade para stack/arquitetura/padrões cross-cutting.
  - **Motivo**: evitar duplicação e conflitos entre PLANs; garantir consistência em idempotência, jobs, eventos e privacidade.
  - **Alternativas consideradas**: definir stack por feature; descartado (inconsistente).
  - **Impactos/Trade-offs**: mudanças de stack passam a exigir atualização do baseline.

- **D-001 — Onboarding como máquina de estados persistida**
  - **Decisão**: modelar onboarding como uma `OnboardingSession` com `status` + `current_step` + `pending_items` persistidos.
  - **Motivo**: suporta interrupções/retomada (FR-012) e mantém interação curta (RNF1).
  - **Alternativas consideradas**: wizard stateless baseado apenas em mensagens; descartado por perder retomada e auditabilidade.
  - **Impactos/Trade-offs**: exige persistência e idempotência por mensagem; reduz bugs de “perdi onde estava”.

- **D-002 — Separar “configuração” de “registros históricos”**
  - **Decisão**: mudanças em metas/restrições/MVD/quiet hours/opt-outs são “configuração” e passam a valer **daqui pra frente**; registros históricos não são reescritos (podem receber anotação de mudança).
  - **Motivo**: alinha com `SPEC-001` (política default) e evita inconsistência em métricas (`SPEC-016`).
  - **Alternativas consideradas**: versionamento completo de histórico; adiado (complexidade).
  - **Impactos/Trade-offs**: simplifica; limita auditoria detalhada de revisões (aceitável no MVP).

- **D-003 — Privacidade como entidade de primeira classe (modo mínimo)**
  - **Decisão**: introduzir `PrivacyPolicy` e `SensitiveEvidenceHandling` desde o onboarding (opt-out por categoria; retenção por categoria), aplicável a todos os fluxos.
  - **Motivo**: cross-cutting mandatória (`SPEC-015`) e reduz bloqueios de adoção (RNF4).
  - **Alternativas consideradas**: tratar privacidade “mais tarde”; descartado (RNF4).
  - **Impactos/Trade-offs**: adiciona complexidade inicial, mas evita retrabalho e decisões perigosas depois.

- **D-004 — Eventos de domínio para evidência e métricas**
  - **Decisão**: registrar eventos (ex.: `onboarding_step_completed`, `privacy_opt_out_set`) para alimentar métricas e auditoria mínima.
  - **Motivo**: `SPEC-016` precisa de registros para tendências/targets; `SPEC-003` e `SPEC-011` exigem evidência e controle de fricção/spam.
  - **Alternativas consideradas**: calcular tudo “na hora” sem eventos; descartado por perda de rastreabilidade.
  - **Impactos/Trade-offs**: exige schema de eventos e disciplina de logging.

## 5) Arquitetura (alto nível)
- **Componentes**
  - **Telegram Adapter**: recebe mensagens/updates, normaliza para comandos de domínio e envia respostas.
  - **Conversation Orchestrator**: roteia intenção (onboarding/ajustes/consulta) e mantém o “próximo passo único”.
  - **Onboarding Service**: lógica de passos, validação leve, retomada, revisão e geração de resumo.
  - **Privacy Service**: leitura/escrita de políticas (opt-out, retenção, modo mínimo) e geração de “disclosure” (explicação curta).
  - **Storage Layer**: repositórios para entidades do onboarding + event log.
  - **Metrics/Event Sink**: grava eventos e contadores mínimos (alinhado a `SPEC-016`).

- **Fluxos**
  - **Fluxo A — Onboarding mínimo**: usuário inicia → coleta metas ativas e restrições mínimas → baseline mínima por domínio relevante → define MVD → configura privacidade básica → gera resumo consultável.
  - **Fluxo B — Diagnóstico progressivo**: após onboarding mínimo, a cada aceitação do usuário → coletar 1 pendência → atualizar resumo e pendências restantes.

## 6) Contratos e interfaces
Definir contratos de domínio (independentes de HTTP/DB), usados pelo adaptador do Telegram.

- **Comando**: `StartOnboarding(user_id, timestamp, locale?, timezone?)`
  - **Saída**: `OnboardingPrompt(next_step, choices?, disclosure?)`

- **Comando**: `SubmitOnboardingAnswer(session_id, step_id, answer, timestamp)`
  - **Saída**: `OnboardingProgress(status, next_step, summary_delta?, pending_items?)`
  - **Erros de domínio**:
    - `ONB_ANSWER_INVALID` (formato fora do esperado)
    - `ONB_CONFLICT_OVERLOAD` (mais de 2 metas intensivas selecionadas; requer escolha)
    - `ONB_STEP_OUT_OF_ORDER` (answer para step não atual; resolver por idempotência/retomada)

- **Comando**: `ResumeOnboarding(user_id)`
  - **Saída**: `OnboardingPrompt(next_step)`

- **Comando**: `ReviseOnboardingField(user_id, field_key, new_value)`
  - **Saída**: `OnboardingSummary(updated_summary, effective_from="now")`
  - **Erros**: `ONB_FIELD_NOT_EDITABLE` (tenta editar registro histórico)

- **Consulta**: `GetOnboardingSummary(user_id)`
  - **Saída**: resumo curto com metas ativas/pausadas, restrições principais, baseline mínima coletada (por domínio), MVD e privacidade/opt-outs.

- **Comando**: `SetPrivacyPolicy(user_id, opt_out_categories, retention_overrides?, minimal_mode?)`
  - **Saída**: confirmação + efeitos práticos (1–2 frases), alinhado a `SPEC-015`.

## 7) Modelo de dados (mínimo)
Entidades e campos essenciais (nomes ilustrativos; podem virar tabelas/coleções).

- **UserProfile**
  - `user_id` (chave)
  - `telegram_user_id`, `primary_chat_id`
  - `timezone`, `locale`
  - `created_at`, `updated_at`

- **PrivacyPolicy** (`SPEC-015`)
  - `user_id` (chave)
  - `opt_out_categories` (C1–C5; em especial C3 “conteúdo sensível”)
  - `retention_days_by_category` (defaults + overrides)
  - `minimal_mode_enabled` (bool)
  - **Sensibilidade/retensão**: define comportamento e expiração; não armazena conteúdo sensível em si.

- **OnboardingSession**
  - `session_id`, `user_id`
  - `status`: `new | in_progress | minimum_completed | completed`
  - `current_step_id`
  - `answers` (map por `step_id` → valor normalizado)
  - `pending_items` (lista curta de pendências de baseline)
  - `started_at`, `last_interaction_at`, `completed_at?`
  - **Retenção**: enquanto o usuário usar o sistema; ao “apagar tudo” (`SPEC-015`), remover e voltar a “usuário novo”.

- **ActiveGoalCycle**
  - `cycle_id`, `user_id`
  - `active_goals` (lista)
  - `paused_goals` (lista + motivo opcional)
  - `intensive_goals` (derivado; no máximo 2)
  - `started_at`, `updated_at`
  - **Retenção**: enquanto ativo + histórico mínimo de mudanças (via event log).

- **BaselineSnapshot**
  - `baseline_id`, `user_id`
  - `domain`: `sleep | english | java | self_esteem | context`
  - `data` (campos mínimos do domínio, conforme SPEC-001)
  - `completeness`: `minimum | partial | complete`
  - `captured_at`, `updated_at`
  - **Sensibilidade**: se houver “conteúdo bruto sensível” (C3), guardar só referência e aplicar política (7 dias ou não guardar).

- **MinimumViableDaily (MVD)**
  - `mvd_id`, `user_id`
  - `items` (lista curta; cada item: domínio/meta, duração estimada, critério observável)
  - `when_to_use` (texto curto)
  - `updated_at`
  - **Retenção**: configuração ativa (com histórico opcional via eventos).

- **DomainEventLog** (`SPEC-016`)
  - `event_id`, `user_id`, `timestamp`
  - `event_type`, `payload_summary` (sem conteúdo sensível por padrão), `sensitivity_level`
  - **Retenção**: alinhar com `SPEC-015` (eventos sem C3 podem reter 12 meses; eventos com C3 devem ser minimizados/redigidos).

## 8) Regras e defaults
- **Limite de metas intensivas**: no máximo 2 ativas por ciclo (`SPEC-001` FR-002; alinhado a `SPEC-010`).
- **Classificação default** (`SPEC-001`): intensivas (Inglês, Java); fundação (sono/saúde/autoestima); aposta semanal (SaaS).
- **Onboarding mínimo obrigatório** (`SPEC-001` FR-004): metas ativas do ciclo + MVD + 1 restrição principal + baseline mínima de sono + baseline mínima das metas intensivas ativas + privacidade/opt-out básico (`SPEC-015`).
- **Privacidade defaults** (`SPEC-015` + `SPEC-016`):
  - C1/C2 (check-ins/planos; evidências não sensíveis): retenção moderada (default 90 dias quando aplicável).
  - C3 (conteúdo sensível bruto): retenção curta (default 7 dias) e opção “não guardar”.
  - C4 (agregados semanais): retenção longa (default 12 meses).
- **Quiet hours default** (capturar no onboarding): 22:00–07:00 (`SPEC-011`), ajustável.
- **Alternativa sem áudio**: permitir onboarding mínimo sem speaking áudio, registrando substituto e marcando speaking baseline como pendente (`SPEC-001` + política de equivalência `SPEC-003` + privacidade `SPEC-015`).

## 9) Observabilidade e métricas
Eventos/logs mínimos (sem conteúdo sensível por padrão) e métricas para validar SCs.

- **Eventos**
  - `onboarding_started`, `onboarding_minimum_completed`, `onboarding_completed`
  - `onboarding_step_completed(step_id)`, `onboarding_resumed`, `onboarding_field_revised(field_key)`
  - `goal_cycle_set(active_goals, intensive_count)`, `mvd_defined(items_count)`
  - `privacy_policy_set(opt_out_categories, minimal_mode)`

- **Métricas (targets iniciais)**
  - **SC-001 (SPEC-001)**: tempo total até `onboarding_minimum_completed` ≤ 10 min (mediana) e drop-off por step.
  - % de usuários/ciclos respeitando limite de 2 metas intensivas (auditoria de `goal_cycle_set`).
  - Taxa de pendências concluídas no diagnóstico progressivo em 14 dias (com “pendente por privacidade” separado).

## 10) Riscos & mitigação
- **Risco**: onboarding longo → abandono.  
  **Mitigação**: máquina de estados + passo único + “onboarding mínimo” explícito e curto; pendências viram diagnóstico progressivo.
- **Risco**: coleta de dados sensíveis reduz confiança.  
  **Mitigação**: disclosure curto + opt-out C3 + modo mínimo por padrão (`SPEC-015`).
- **Risco**: usuário tenta ativar metas demais e quebra consistência.  
  **Mitigação**: limite rígido e linguagem protetiva (`SPEC-001`/`SPEC-010`).
- **Risco**: revisões “reescrevendo passado” quebram métricas.  
  **Mitigação**: D-002; registrar mudanças via eventos (`SPEC-016`).

## 11) Rollout / migração
- **Feature flag**: `onboarding_v1` (permite iterar fluxo sem quebrar usuários).
- **Migração**: greenfield; nenhuma migração inicial. Evoluções de schema devem ser backward compatible (novos campos opcionais).

## 12) Plano de testes (como validar)
- **Unit**
  - Máquina de estados do onboarding (transições; retomada; idempotência por `step_id`).
  - Validação de limites (2 metas intensivas) e defaults.
  - Políticas de privacidade: aplicação de opt-out e retenção (simulada).
- **Integration**
  - Persistência completa: iniciar → responder → interromper → retomar → revisar campo → obter resumo.
  - Redação/minimização: garantir que logs/eventos não carregam C3 por padrão (`SPEC-015`).
- **E2E**
  - Fluxo P1 completo (onboarding mínimo) via Telegram adapter (mensagens sequenciais).
  - Fluxo parcial (tempo curto) e diagnóstico progressivo (1 pendência/dia).
- **Manual / acceptance**
  - Perguntas “o que você guarda?” e opt-out operando com explicação curta (RNF4).
  - Tom não punitivo em overload e recusas (RNF3).

## 13) Task breakdown (execução)
Tarefas pequenas (1–4h) ordenadas por dependências.

1) **Definir contratos de domínio do onboarding**
   - **Entrada**: `SPEC-001` FR-001..FR-013; `SPEC-015` FR-001..FR-006
   - **Saída**: documento/artefato de interfaces (commands/queries/events) com erros de domínio e payloads mínimos
   - **Critério de pronto**: revisão interna confirma cobertura de todos os cenários P1 e edge cases relevantes sem depender de outras SPECS

2) **Modelar entidades mínimas e retenção/sensibilidade**
   - **Entrada**: `SPEC-001` Key Entities + FR-013; `SPEC-015` categorias/retensão; `SPEC-016` histórico mínimo
   - **Saída**: schema lógico (entidades + campos + índices lógicos) e tabela de retenção por categoria
   - **Critério de pronto**: cada entidade declara sensibilidade (C1–C5) e política de retenção/opt-out; não há coleta “extra”

3) **Implementar persistência de `OnboardingSession` + idempotência por step**
   - **Entrada**: cenários de retomada e out-of-order (`SPEC-001`)
   - **Saída**: repositório/DAO + estratégia de idempotência para respostas repetidas
   - **Critério de pronto**: testes de integração passam para “interrompe e retoma” e “responde duas vezes o mesmo step”

4) **Implementar máquina de estados do onboarding mínimo (P1)**
   - **Entrada**: `SPEC-001` User Story 1 + ACs
   - **Saída**: orquestração dos passos mínimos e geração de `OnboardingSummary`
   - **Critério de pronto**: fluxo completo cria metas ativas, baseline mínima, MVD e resumo consultável

5) **Aplicar governança: limite de 2 metas intensivas**
   - **Entrada**: `SPEC-001` FR-002 + Edge Cases; `SPEC-010`
   - **Saída**: regra aplicada na seleção de metas com 1 pergunta de escolha e default protetivo
   - **Critério de pronto**: tentativa de 3ª meta intensiva nunca conclui sem resolver; defaults batem a SPEC

6) **Implementar “baseline suficiente sem áudio” e marcação de pendência**
   - **Entrada**: `SPEC-001` Scenario 7; `SPEC-003` equivalência; `SPEC-015` modo mínimo
   - **Saída**: coleta alternativa + flags `speaking_baseline_pending`
   - **Critério de pronto**: onboarding mínimo conclui com substituto e resumo deixa explícito o que ficou pendente

7) **Implementar diagnóstico progressivo (P2) com 1 pendência por vez**
   - **Entrada**: `SPEC-001` User Story 2
   - **Saída**: mecanismo de `pending_items` e coleta incremental
   - **Critério de pronto**: após onboarding mínimo, usuário consegue completar pendências uma a uma sem “virar questionário”

8) **Implementar revisão de campos de configuração (vale daqui pra frente)**
   - **Entrada**: `SPEC-001` User Story 3 + política default
   - **Saída**: comando `ReviseOnboardingField` e atualização consistente do resumo
   - **Critério de pronto**: editar restrição/quiet hours/MVD altera apenas configuração; registros históricos não mudam

9) **Implementar `PrivacyPolicy` (opt-out + retenção overrides) e disclosures**
   - **Entrada**: `SPEC-015` FR-001..FR-006
   - **Saída**: leitura/escrita de política + respostas curtas “o que guardo e por quê?”
   - **Critério de pronto**: opt-out C3 impede armazenamento de conteúdo sensível bruto; explicações são consistentes entre fluxos

10) **Registrar eventos e métricas mínimas do onboarding**
   - **Entrada**: `SPEC-016` (registros/targets) + SC-001 de `SPEC-001`
   - **Saída**: event log com eventos do onboarding + contadores de tempo/conclusão
   - **Critério de pronto**: é possível medir tempo até onboarding mínimo e drop-off por step sem ler conteúdo sensível

11) **Implementar comandos de consulta (“meu resumo de onboarding”)**
   - **Entrada**: `SPEC-001` Scenario 5 + FR-010
   - **Saída**: endpoint/handler que retorna resumo curto e revisável
   - **Critério de pronto**: funciona em ≤ 1 mensagem e sem depender de outras features

## 14) Open questions (se existirem)
- (Default adotado) **Timezone**: perguntar no onboarding se faltar; se ainda ausente, usar timezone do servidor e permitir ajuste posterior em “configurações”.

