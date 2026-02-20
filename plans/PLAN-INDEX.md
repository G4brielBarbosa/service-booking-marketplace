# PLAN Index — Mapa de PLANs (1 por SPEC)

Este índice organiza a camada “HOW” (PLANs) e mantém o vínculo explícito com SPECS/PRD.

- **Fonte de verdade do “O QUÊ”**: `PRD.md` + `specs/`
- **Fonte de verdade do “COMO”**: `plans/`
- **Regra**: 1 PLAN por SPEC (um slice implementável)

## Como usar

- Para criar um novo PLAN, siga `plans/PLAN-GUIDE.md`.
- Ao concluir um PLAN, atualize este índice marcando como “pronto”.

## Convenção

- `SPEC-XXX` → `plans/PLAN-XXX-<slug>.md`
- Um PLAN deve referenciar explicitamente a SPEC e o PRD Base.

## Ordem sugerida (MVP primeiro)

Baseado em `specs/SPEC-INDEX.md` (P1 → P2 → P3).

### P1 — MVP

- `SPEC-001` → `plans/PLAN-001-onboarding.md` — **pronto**
- `SPEC-002` → `plans/PLAN-002-rotina-diaria.md` — **pronto**
- `SPEC-003` → `plans/PLAN-003-quality-gates.md` — **pronto**
- `SPEC-004` → `plans/PLAN-004-ingles-diario.md` — **pronto**
- `SPEC-005` → `plans/PLAN-005-java-diario.md` — **pronto**
- `SPEC-006` → `plans/PLAN-006-sono.md` — **pronto**
- `SPEC-007` → `plans/PLAN-007-revisao-semanal.md` — **pronto**

### P2 — Robustez / escala do comportamento

- `SPEC-008` → `plans/PLAN-008-backlog-inteligente.md` — **pronto**
- `SPEC-009` → `plans/PLAN-009-falso-progresso.md` — **pronto**
- `SPEC-010` → `plans/PLAN-010-governanca-personalizacao.md` — **pronto**
- `SPEC-011` → `plans/PLAN-011-nudges-dias-ruins.md` — **pronto**

### P3 — Metas paralelas + cross-cutting

- `SPEC-012` → `plans/PLAN-012-vida-saudavel.md` — **pronto**
- `SPEC-013` → `plans/PLAN-013-autoestima.md` — **pronto**
- `SPEC-014` → `plans/PLAN-014-saas-aposta-semanal.md` — **pronto**
- `SPEC-015` → `plans/PLAN-015-privacidade.md` — **pronto**
- `SPEC-016` → `plans/PLAN-016-metricas-registros.md` — **pronto**

## Dicas de consistência entre PLANs (cross-cutting)

- **Dados e retenção**: todo PLAN que criar/ler dados deve citar `SPEC-015` e alinhar retenção.
- **Gates e evidência**: todo PLAN que marca “concluído” deve alinhar com `SPEC-003`.
- **Métricas e tendências**: todo PLAN que consome métricas deve alinhar com `SPEC-016`.
- **Nudges/proatividade**: todo PLAN que envia proativo deve alinhar com `SPEC-011`.
