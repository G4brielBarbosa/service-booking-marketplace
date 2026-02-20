# Technical Plan: PLAN-013 — Autoestima: registros curtos + revisão de padrões + micro-exposições (“ações de coragem”)

**Created**: 2026-02-20  
**Spec**: `specs/SPEC-013-autoestima.md`  
**PRD Base**: §8.5, §5.5, §5.3, §9.1, §9.2, §14, §11 (RNF1–RNF3)  
**Related Specs**: `SPEC-002`, `SPEC-007`, `SPEC-011`, `SPEC-015`, `SPEC-016`, `SPEC-003`

## 1) Objetivo do plano
- Implementar captura de **registro curto** de autocrítica (gatilho → pensamento → resposta alternativa) com versão mínima (dia ruim/alta emoção).
- Implementar **intervenção breve** quando emoção estiver alta (estabilização) e permitir “prefiro não dizer”.
- Implementar **ações de coragem** (micro-exposições) com critério observável de feito e registro de tentativa vs feito.
- Alimentar revisão semanal (`SPEC-007`) com padrões e contagem de ações de coragem via métricas (`SPEC-016`), respeitando privacidade por padrão (`SPEC-015`).

## 2) Non-goals (fora do escopo)
- Não diagnosticar/tratar transtornos; não substituir terapia.
- Não coletar journaling longo nem narrativas detalhadas por padrão.
- Não enviar nudges proativos de conteúdo emocional no MVP; qualquer proatividade segue `SPEC-011`.

## 3) Assumptions (assunções)
- Conteúdo pode ser altamente sensível (C3). MVP deve funcionar com **minimização e modo mínimo**.
- Os registros podem ser “abstratos” (sem detalhes) e ainda úteis para padrão/contagem.
- “Ação de coragem” é tratada como meta de fundação (não intensiva), mas deve ser observável e pequena.

## 4) Decisões técnicas (Decision log)
- **D-000 — Baseline de plataforma**
  - **Decisão**: adotar o baseline `plans/PLAN-000-platform-baseline.md` para execução/armazenamento, com ênfase em privacidade C3, redaction e retenção/expiração.
  - **Motivo**: autoestima lida com conteúdo sensível; precisa de padrões consistentes de opt-out/modo mínimo e logs neutros (`SPEC-015`).
  - **Alternativas consideradas**: tratar privacidade “caso a caso”; descartado.
  - **Impactos/Trade-offs**: baseline vira referência obrigatória para política de dados.

- **D-001 — Conteúdo sensível com defaults conservadores**
  - **Decisão**: por padrão, armazenar registros de autoestima como **C3** com retenção curta e opção “não guardar” (apenas contar/derivar).
  - **Motivo**: `SPEC-015` (confiança) + natureza sensível do domínio.
  - **Alternativas consideradas**: guardar indefinidamente; descartado.
  - **Impactos/Trade-offs**: menos histórico textual; manter agregados e padrões em C4.

- **D-002 — Separar “conteúdo” de “metadados úteis”**
  - **Decisão**: persistir metadados mínimos não sensíveis (timestamp, intensidade opcional, tags abstratas) mesmo quando o usuário opta por não guardar texto.
  - **Motivo**: permite tendências e revisão semanal sem conteúdo sensível.
  - **Alternativas consideradas**: apagar tudo; reduz utilidade.
  - **Impactos/Trade-offs**: precisa disclosure claro (“vou guardar só contagem/metadata”).

- **D-003 — Gate leve para ações de coragem (não para conteúdo emocional)**
  - **Decisão**: “feito” para ação de coragem exige evidência mínima de 1 frase (ou checkbox) e critério observável; registro de autocrítica não deve virar gate pesado (apenas “registrado/parcial”).
  - **Motivo**: manter baixa fricção e segurança psicológica; alinhar a `SPEC-003` (hábito/fundação).
  - **Alternativas consideradas**: exigir detalhes completos; descartado.

## 5) Arquitetura (alto nível)
- **Componentes**
  - **SelfEsteem Domain Service**: captura registros, sugere intervenção breve, gerencia ações de coragem e revisão de padrões.
  - **Privacy Service** (`SPEC-015`): opt-out C3, retenção curta, modo mínimo e apagamento.
  - **Metrics Aggregator** (`SPEC-016`): contagem de registros e ações; tendência simples.
  - **Weekly Review Integration** (`SPEC-007`): injeta padrões e 1 experimento da semana.

## 6) Contratos e interfaces
- **Comando**: `StartSelfEsteemRecord(user_id, local_date, context=normal|high_emotion|low_energy, timestamp)`
  - **Saída**: `RecordPrompt(fields_required, minimal_option, privacy_disclosure_short)`

- **Comando**: `SubmitSelfEsteemRecord(user_id, record_payload, timestamp, storage_mode keep|minimal|discard_content)`
  - **Saída**: `RecordReceipt(status complete|partial, next_step_short)`

- **Comando**: `StartHighEmotionIntervention(user_id, timestamp)`
  - **Saída**: `InterventionPlan(steps_short, one_question_optional)`

- **Comando**: `DefineCourageAction(user_id, week_id, description_short, done_criteria, frequency_min, timestamp)`
  - **Saída**: `CourageActionView(status=planned)`

- **Comando**: `RecordCourageActionOutcome(user_id, action_id, outcome done|attempt|defer, evidence_short?, timestamp)`
  - **Saída**: `OutcomeReceipt(next_variant_smaller?)`

- **Consulta**: `GetSelfEsteemToday(user_id, local_date)`
  - **Saída**: lista curta (resumo abstrato) de registros e ações do dia (sem conteúdo sensível por padrão)

- **Consulta**: `GetSelfEsteemWeekSummary(user_id, week_id)`
  - **Saída**: padrões (tags/gatilhos abstratos), contagem de ações, 1 sugestão de experimento

## 7) Modelo de dados (mínimo)
- **SelfEsteemRecord**
  - `record_id`, `user_id`, `local_date`, `timestamp`
  - `trigger_text?`, `thought_text?`, `alternative_response_text?` (C3; opcionais)
  - `intensity_0_10?`
  - `tags_abstract[]` (ex.: `work`, `self_doubt`, `comparison`) (C4)
  - `status complete|partial`
  - `storage_mode` (kept|minimal|discard_content)
  - **Retenção**: conteúdo C3 curto (default 30 dias) e opt-out; agregados/tags C4 12 meses.

- **CourageAction**
  - `action_id`, `user_id`, `week_id`
  - `description_short`, `done_criteria`, `frequency_min`
  - `status planned|done|attempt|deferred`
  - `evidence_short?` (1 frase; evitar detalhes sensíveis)

- **SelfEsteemWeeklyAggregates** (`SPEC-016`)
  - `user_id`, `week_id`
  - `records_count`, `avg_intensity?`
  - `courage_actions_done_count`, `attempt_count`
  - `top_tags_abstract[]`

## 8) Regras e defaults
- Registro curto padrão: 3 campos; versão mínima: gatilho + resposta alternativa (1 frase).
- Alta emoção: oferecer intervenção breve; permitir “prefiro não dizer”; registrar apenas metadata se desejado.
- Ação de coragem: critério observável e pequeno; falha vira `attempt` e sugere versão menor.
- Privacidade: default “modo mínimo” disponível; sempre evitar mostrar conteúdo sensível em resumos.
- Anti-spam: nenhum lembrete proativo emocional no MVP; se houver, budgets/quiet hours e títulos neutros (`SPEC-011`).

## 9) Observabilidade e métricas
- **Eventos**: `self_esteem_recorded(partial|complete, storage_mode)`, `high_emotion_intervention_started`, `courage_action_defined`, `courage_action_outcome`
- **Métricas**: contagem semanal de ações de coragem; tendência de intensidade (quando coletada); taxa de registros concluídos sem virar journaling.

## 10) Riscos & mitigação
- **Risco**: coleta sensível reduz confiança. Mitigar com opt-out, modo mínimo, retenção curta e disclosures.
- **Risco**: usuário em alta emoção se sente pressionado. Mitigar com intervenção breve e opção “prefiro não dizer”.
- **Risco**: ações de coragem viram metas grandes. Mitigar com recorte e critério de feito pequeno.

## 11) Rollout / migração
- **Feature flag**: `self_esteem_v1`.

## 12) Plano de testes (como validar)
- **Unit**: versões mínima/normal; storage_mode e retenção; resumos sem conteúdo sensível.
- **Integration**: opt-out C3 preserva agregados; revisão semanal consome padrões/contagens.
- **E2E**: episódio alta emoção → intervenção breve; ação de coragem definida → outcome → tendência semanal.
- **Manual**: tom protetivo; apagamento por período/categoria (`SPEC-015`).

## 13) Task breakdown (execução)
1) Definir schema `SelfEsteemRecord`/`CourageAction` + retenção/opt-out  
   - **Entrada**: `SPEC-013` FR-001..FR-006 + `SPEC-015`  
   - **Saída**: modelo lógico e políticas (C3 curto, C4 agregados)  
   - **Critério de pronto**: modo mínimo funciona sem perder contagens

2) Implementar fluxo de registro curto (normal e mínimo) + “prefiro não dizer”  
   - **Entrada**: User Story 1/2  
   - **Saída**: comandos e receipts curtos  
   - **Critério de pronto**: registro ≤2 min e sem insistência

3) Implementar ações de coragem com critério de feito e outcomes  
   - **Entrada**: FR-005 + cenários  
   - **Saída**: definição/registro e sugestão de versão menor após tentativa  
   - **Critério de pronto**: contagem semanal e status correto (done/attempt)

4) Implementar agregados semanais e integração com revisão semanal  
   - **Entrada**: `SPEC-016` + `SPEC-007`  
   - **Saída**: `SelfEsteemWeekSummary` consumível na revisão  
   - **Critério de pronto**: revisão semanal sugere 1 experimento baseado em padrões/contagens

## 14) Open questions (se existirem)
- (Default adotado) **Retenção de conteúdo C3 de autoestima**: 30 dias (mais curta que C1/C2), mantendo apenas agregados/tags por 12 meses. Ajustável via `SPEC-015`.

