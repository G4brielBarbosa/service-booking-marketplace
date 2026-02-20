# PLAN Guide — Como escrever PLANs (camada “HOW”) a partir do PRD + SPECS

Este guia existe para que agentes (e você) consigam criar planos técnicos **consistentes**, **revisáveis** e **executáveis**, mantendo a separação:

- **PRD/SPECS** = “O QUÊ / POR QUÊ” (produto, comportamento, critérios)
- **PLAN** = “COMO” (decisões técnicas, arquitetura, contratos, riscos, rollout)

O PLAN pode mencionar stack, integrações e decisões de arquitetura. Evite grandes blocos de código; o objetivo é orientar a implementação com clareza.

## Regras de ouro

- **1 PLAN por SPEC**: cada arquivo `plans/PLAN-XXX-*.md` deve implementar exatamente 1 SPEC (slice).
- **Link explícito**: todo PLAN deve referenciar a SPEC e as seções do PRD-base.
- **Decisões explícitas**: registrar escolhas e trade-offs (e por que).
- **Sem escopo implícito**: incluir “Non-goals” e “Out of scope”.
- **MVP-first**: descrever o mínimo para entregar o comportamento da SPEC; o resto vai para “Later”.
- **Privacidade por padrão**: respeitar `SPEC-015` (minimização, retenção, opt-out, modo mínimo).

## Template obrigatório (copie e preencha)

```markdown
# Technical Plan: PLAN-XXX — [SPEC TITLE]

**Created**: YYYY-MM-DD
**Spec**: `specs/SPEC-XXX-...md`
**PRD Base**: [listar seções do PRD]
**Related Specs**: [ex.: SPEC-003, SPEC-015, SPEC-016]

## 1) Objetivo do plano
[1–3 bullets: o que vai ser implementado para cumprir a SPEC.]

## 2) Non-goals (fora do escopo)
- ...

## 3) Assumptions (assunções)
- ...

## 4) Decisões técnicas (Decision log)
Liste decisões como “D-001…”, cada uma com:
- **Decisão**:
- **Motivo**:
- **Alternativas consideradas**:
- **Impactos/Trade-offs**:

## 5) Arquitetura (alto nível)
Descreva componentes e responsabilidades.
- **Componentes**: [ex.: BotTelegram, Scheduler, Storage, Analytics]
- **Fluxos**: [1–2 fluxos principais em bullets; opcional: mermaid]

## 6) Contratos e interfaces
Defina contratos em nível de engenharia (sem implementação):
- **Entradas/Saídas**: payloads e campos principais
- **Estados**: ex.: `pending/in_progress/completed/blocked`
- **Erros**: códigos/mensagens de domínio (não HTTP)

## 7) Modelo de dados (mínimo)
Liste entidades/tabelas/coleções (alto nível) + campos essenciais.
Inclua retenção e sensibilidade (referenciar `SPEC-015`).

## 8) Regras e defaults
Transcreva defaults que a implementação precisa seguir (ex.: limiares, janelas, budgets).
Se o default já está em uma SPEC (ex.: `SPEC-011`, `SPEC-003`), apenas referencie.

## 9) Observabilidade e métricas
- **Logs/eventos** mínimos necessários
- **Métricas** para validar Success Criteria (quando aplicável)

## 10) Riscos & mitigação
- ...

## 11) Rollout / migração
- **Feature flags** (se necessário)
- **Backfill / migração de dados** (se necessário)

## 12) Plano de testes (como validar)
- **Unit**:
- **Integration**:
- **E2E**:
- **Manual / acceptance**:

## 13) Task breakdown (execução)
Quebrar em tarefas pequenas (ideal 1–4h cada), ordenadas por dependências.
Cada tarefa deve indicar:
- **Entrada** (qual requisito/trecho da SPEC cobre)
- **Saída** (artefato/resultado observável)
- **Critério de pronto**

## 14) Open questions (se existirem)
Se sobrar alguma dúvida, ela deve ser pequena e explícita.
```

## Convenção de nomes

- `plans/PLAN-001-onboarding.md` (para `SPEC-001`)
- `plans/PLAN-002-rotina-diaria.md` (para `SPEC-002`)
- etc.

## Checklist final (antes de considerar o PLAN “pronto”)

- O PLAN implementa **apenas** a SPEC alvo e cita `SPEC-015` quando tocar dados sensíveis.
- Existe “Non-goals” para conter escopo.
- Decisões técnicas estão explícitas (Decision log).
- Há um modelo mínimo de dados/contratos suficiente para implementar.
- Há um task breakdown ordenado e testável.
```
