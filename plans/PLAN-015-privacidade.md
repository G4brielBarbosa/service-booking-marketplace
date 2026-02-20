# Technical Plan: PLAN-015 — Privacidade por padrão: dados mínimos, transparência e controle do usuário

**Created**: 2026-02-20  
**Spec**: `specs/SPEC-015-privacidade-por-padrao.md`  
**PRD Base**: §11 (RNF4), §5.1, §5.3, §5.4, §9.1, §9.2, §10 (R6), §2, §14  
**Related Specs**: `SPEC-001`, `SPEC-002`, `SPEC-003`, `SPEC-011`, `SPEC-016`, `SPEC-007`, `SPEC-013`

## 1) Objetivo do plano
- Implementar **PrivacyPolicy** como componente cross-cutting (categorias C1–C5) com defaults de retenção, opt-out por categoria (especialmente C3) e **modo mínimo**.
- Implementar transparência: responder “o que você guarda e por quê?” por fluxo/categoria em linguagem curta.
- Implementar controle do usuário: **inventário de dados**, **apagar por categoria/período**, ajustar retenção, e “apagar tudo”.
- Garantir que ausência/remoção de dados **não vira punição**: sistema explicita limitações e ajusta recomendações.

## 2) Non-goals (fora do escopo)
- Não detalhar criptografia/infra/segurança operacional (fora do nível SPEC/PLAN aqui).
- Não prometer anonimização perfeita; foco em minimização, controle e retenção.
- Não implementar export completo de dados no MVP (pode entrar depois), mas garantir inventário e apagamento.

## 3) Assumptions (assunções)
- Todo dado armazenado é classificado em uma categoria C1–C5.
- Existe um event log/métricas agregadas (`SPEC-016`) que podem sobreviver à remoção de detalhes, quando configurado pelo usuário.
- Fluxos que coletam conteúdo sensível (ex.: áudio, textos emocionais) podem operar em “processar e descartar” mantendo derivados não sensíveis.

## 4) Decisões técnicas (Decision log)
- **D-000 — Baseline de plataforma**
  - **Decisão**: adotar o baseline `plans/PLAN-000-platform-baseline.md` como referência de stack/execução (Go API + worker/jobs; Postgres/sqlc; Redis) e como ponto único para políticas e defaults cross-cutting.
  - **Motivo**: privacidade governa todos os fluxos e precisa ser aplicada de forma uniforme (intake, redaction, expiração, deleção).
  - **Alternativas consideradas**: cada feature implementar suas próprias regras; descartado.
  - **Impactos/Trade-offs**: mudanças de defaults passam a exigir atualizar baseline + este PLAN.

- **D-001 — Classificação e políticas por categoria**
  - **Decisão**: centralizar política em `PrivacyPolicy` com:
    - opt-out por categoria (especialmente C3)
    - retenção (dias) por categoria
    - modo mínimo (bool) + preferências de redaction
  - **Motivo**: consistência cross-cutting e previsibilidade.
  - **Alternativas consideradas**: regras por feature isoladas; descartado.
  - **Impactos/Trade-offs**: exige que todos os fluxos “declarem categoria”; melhora governança.

- **D-002 — Derivados > conteúdo bruto**
  - **Decisão**: preferir armazenar derivados não sensíveis (C4) e resultados (GateResult, contagens) em vez de conteúdo bruto (C3).
  - **Motivo**: reduzir risco e manter utilidade.
  - **Alternativas consideradas**: guardar tudo e pedir desculpas depois; descartado.

- **D-003 — Apagamento é irreversível e auditável**
  - **Decisão**: `DeletionRequest` com confirmação e registro mínimo (sem conteúdo) para auditoria (“apagou categoria X em período Y”).
  - **Motivo**: clareza e confiança; evita acidentes.
  - **Alternativas consideradas**: apagar sem confirmação; descartado.

## 5) Arquitetura (alto nível)
- **Componentes**
  - **Privacy Service**: aplica policy no momento de coleta/armazenamento, e fornece APIs de transparência/controle.
  - **Data Classification Registry**: mapeia entidades/campos → C1..C5 e sensibilidade (C3).
  - **Retention/Expiry Worker**: expira dados automaticamente por categoria.
  - **Deletion Engine**: executa apagamento por categoria/período e “apagar tudo”.
  - **Redaction Layer**: garante que respostas (plano, nudges, resumos) não vazem C3.

## 6) Contratos e interfaces
- **Consulta**: `ExplainDataPolicy(user_id, context_flow?)`
  - **Saída**: `PolicyExplanation(categories[], why_short, retention_short, controls_short)`

- **Consulta**: `GetDataInventorySummary(user_id)`
  - **Saída**: `DataInventorySummary(categories[{name, approx_range, counts, sensitivity_note}], default_view=summary)`

- **Comando**: `SetPrivacyPolicy(user_id, opt_out_categories?, retention_overrides?, minimal_mode?, timestamp)`
  - **Saída**: `PolicyReceipt(effects_short)`

- **Comando**: `RequestDeletion(user_id, category, period?, confirm=false, timestamp)`
  - **Saída**: `DeletionPreview(what_will_be_deleted, impact_short)` quando `confirm=false`
  - **Saída**: `DeletionReceipt(done, impact_short)` quando `confirm=true`

- **Comando**: `DeleteAllUserData(user_id, confirm=false, timestamp)`
  - **Saída**: preview/receipt similar

## 7) Modelo de dados (mínimo)
- **PrivacyPolicy** (`SPEC-015`)
  - `user_id`
  - `opt_out_categories[]` (C1..C5; foco em C3)
  - `retention_days_by_category` (defaults + overrides)
  - `minimal_mode_enabled` (bool)
  - `updated_at`

- **DataInventorySummary** (materializado ou computado)
  - `user_id`
  - `counts_by_category`, `date_range_by_category`
  - `last_updated_at`

- **DeletionRequest**
  - `request_id`, `user_id`
  - `category`, `period_start?`, `period_end?`
  - `confirmed_at?`, `executed_at?`, `status`
  - `impact_short`

## 8) Regras e defaults
- **Categorias** (como na SPEC):
  - C1 check-ins/planos
  - C2 evidências de aprendizagem (não sensíveis)
  - C3 conteúdo sensível bruto (áudio, textos emocionais)
  - C4 métricas agregadas
  - C5 governança/preferências
- **Retenção default (MVP)** (alinhada a `SPEC-016` e SPECS):
  - C1/C2 (detalhes diários): 90 dias
  - C3 (bruto): 7 dias (ou menos; autoestima pode ser 30 dias configurável) e opt-out “não guardar”
  - C4 (agregados): 12 meses
  - C5: enquanto usuário usar, com opção de apagar tudo
- **Modo mínimo**
  - preferir `discard_content` para C3 (processar e descartar) e guardar só derivados (GateResult, contagens).
- **Transparência**
  - sempre que pedir C3, mostrar disclosure 1–2 frases (o mínimo guardado + por quanto tempo + como não guardar).
- **Sem punição**
  - se dados expirarem/apagarem, consultas e revisões mostram “não tenho X porque expirou/apagou” e pedem 1 coleta mínima se necessário.

## 9) Observabilidade e métricas
- **Eventos**
  - `privacy_policy_set`, `opt_out_enabled(category)`, `minimal_mode_enabled`
  - `retention_expired(category,count)`
  - `deletion_requested(category,period)`, `deletion_executed`
- **Métricas**
  - % de usuários com opt-out C3 ativo (proxy de confiança/necessidade)
  - taxa de deleções e impactos em funcionalidades (para ajustar defaults)

## 10) Riscos & mitigação
- **Risco**: inconsistência entre fluxos (um guarda C3 sem disclosure).  
  **Mitigação**: registry de classificação + testes de contrato.
- **Risco**: apagamento quebra revisão semanal.  
  **Mitigação**: usar agregados quando possível; explicitamente parcial; sugerir coleta mínima.
- **Risco**: usuário escreve dados pessoais acidentalmente.  
  **Mitigação**: oferecer “não guardar esta mensagem” e redaction por padrão.

## 11) Rollout / migração
- **Feature flag**: `privacy_v1`.
- Ao habilitar, criar `PrivacyPolicy` default para usuários existentes.

## 12) Plano de testes (como validar)
- **Unit**
  - policy defaults e overrides
  - opt-out C3 implica “discard after processing”
  - retention expiry por categoria
- **Integration**
  - onboarding/quality gates/nudges respeitam policy e redaction
  - apagar por período atualiza inventário e consultas futuras
- **E2E**
  - “o que você guarda?” consistente em onboarding, evidência e revisão semanal
  - “apagar tudo” retorna ao estado de usuário novo
- **Manual**
  - linguagem simples e curta; confirmação irreversível

## 13) Task breakdown (execução)
1) Definir registry C1–C5 e mapear entidades/campos  
   - **Entrada**: `SPEC-015` + planos existentes  
   - **Saída**: tabela de classificação e sensibilidade  
   - **Critério de pronto**: toda entidade tem categoria e regra de retenção

2) Implementar `PrivacyPolicy` + APIs de set/get + disclosures  
   - **Entrada**: FR-001/FR-005/FR-006  
   - **Saída**: policy aplicada e mensagens curtas  
   - **Critério de pronto**: opt-out C3 e modo mínimo funcionando

3) Implementar inventário e previews de apagamento  
   - **Entrada**: FR-003/FR-004  
   - **Saída**: `DataInventorySummary` + `DeletionPreview`  
   - **Critério de pronto**: usuário entende o que será apagado e impacto

4) Implementar execução de apagamento + expiração automática  
   - **Entrada**: FR-004/FR-006  
   - **Saída**: Deletion Engine + Retention Worker  
   - **Critério de pronto**: dados expiram e deletam sem quebrar o sistema (apenas reduz utilidade)

## 14) Open questions (se existirem)
- (Default adotado) **Retenção de C3 para autoestima**: manter default 30 dias (mais curta que C1/C2, mais longa que áudio 7 dias) e permitir ajustar para 7 dias ou “não guardar” via policy.

