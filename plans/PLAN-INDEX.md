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

- `SPEC-001` → `plans/PLAN-001-onboarding.md` — **todo**
- `SPEC-002` → `plans/PLAN-002-rotina-diaria.md` — **todo**
- `SPEC-003` → `plans/PLAN-003-quality-gates.md` — **todo**
- `SPEC-004` → `plans/PLAN-004-ingles-diario.md` — **todo**
- `SPEC-005` → `plans/PLAN-005-java-diario.md` — **todo**
- `SPEC-006` → `plans/PLAN-006-sono.md` — **todo**
- `SPEC-007` → `plans/PLAN-007-revisao-semanal.md` — **todo**

### P2 — Robustez / escala do comportamento

- `SPEC-008` → `plans/PLAN-008-backlog-inteligente.md` — **todo**
- `SPEC-009` → `plans/PLAN-009-falso-progresso.md` — **todo**
- `SPEC-010` → `plans/PLAN-010-governanca-personalizacao.md` — **todo**
- `SPEC-011` → `plans/PLAN-011-nudges-dias-ruins.md` — **todo**

### P3 — Metas paralelas + cross-cutting

- `SPEC-012` → `plans/PLAN-012-vida-saudavel.md` — **todo**
- `SPEC-013` → `plans/PLAN-013-autoestima.md` — **todo**
- `SPEC-014` → `plans/PLAN-014-saas-aposta-semanal.md` — **todo**
- `SPEC-015` → `plans/PLAN-015-privacidade.md` — **todo**
- `SPEC-016` → `plans/PLAN-016-metricas-registros.md` — **todo**

## Dicas de consistência entre PLANs (cross-cutting)

- **Dados e retenção**: todo PLAN que criar/ler dados deve citar `SPEC-015` e alinhar retenção.
- **Gates e evidência**: todo PLAN que marca “concluído” deve alinhar com `SPEC-003`.
- **Métricas e tendências**: todo PLAN que consome métricas deve alinhar com `SPEC-016`.
- **Nudges/proatividade**: todo PLAN que envia proativo deve alinhar com `SPEC-011`.
