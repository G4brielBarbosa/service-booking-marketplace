# SPEC Index — Super Assistente Pessoal (Telegram-first)

**Fonte de verdade**: `PRD.md` (este índice só organiza o trabalho em SPECS “O QUE”).

## Como usar este índice

- Cada item abaixo é uma SPEC **independentemente testável**.
- Para escrever uma SPEC, siga o guia em `specs/SPEC-GUIDE.md`.
- Cada SPEC deve referenciar explicitamente as seções do PRD que a fundamentam.

## Ordem sugerida (P1 → P2 → P3)

### P1 — MVP (base: PRD §14)

1. `SPEC-001` — Onboarding e Diagnóstico Leve (1–2 semanas)  
   - **Base PRD**: §5.1, §14

2. `SPEC-002` — Rotina Diária (Telegram-first): Check-in + Plano A/B/C + Execução Guiada  
   - **Base PRD**: §§5.2, 5.3, 9.1, 10 (R1, R6), 11 (RNF1, RNF2)

3. `SPEC-003` — Quality Gates & Evidência Mínima (aprendizagem e hábitos)  
   - **Base PRD**: §5.4, §§8.1–8.3, 10 (R2, R3), 13 (riscos)

4. `SPEC-004` — Inglês Diário: Input + Output + Retrieval (rubrica + erros recorrentes)  
   - **Base PRD**: §8.1, §14

5. `SPEC-005` — Java Diário: Prática Deliberada + Retrieval + Revisão de Erros  
   - **Base PRD**: §8.2, §14

6. `SPEC-006` — Sono: Diário + Rotina Pré-sono + 1 Intervenção simples/semana  
   - **Base PRD**: §6.1, §8.3, §14

7. `SPEC-007` — Revisão Semanal: Painel + 3 decisões (manter/ajustar/pausar) + alvos da semana  
   - **Base PRD**: §§5.2, 5.5, 9.2, 10 (R5), §14

### P2 — Robustez / escala do comportamento

8. `SPEC-008` — Backlog Inteligente e Priorização baseada em lacunas observadas  
   - **Base PRD**: §5.2, 10 (R4)

9. `SPEC-009` — Detecção de falhas reais (“falso progresso”) e reforço automático  
   - **Base PRD**: 10 (R3), 13 (riscos)

10. `SPEC-010` — Personalização progressiva e governança de metas em paralelo (limites de overload)  
   - **Base PRD**: §6.2, 10 (R7), 13 (riscos)

11. `SPEC-011` — Nudges/Lembretes “sem spam” + robustez a dias ruins (Planos B/C + MVD)  
   - **Base PRD**: §5.3, 11 (RNF1, RNF2, RNF3)

### P3 — Metas paralelas + cross-cutting

12. `SPEC-012` — Vida saudável (planejamento semanal de atividade) + métricas + limites por dor/energia  
   - **Base PRD**: §8.4

13. `SPEC-013` — Autoestima: registros curtos + revisão de padrões + micro-exposições  
   - **Base PRD**: §8.5

14. `SPEC-014` — SaaS (aposta semanal): bloco profundo + microblocos  
   - **Base PRD**: §8.6

15. `SPEC-015` — Privacidade por padrão (dados mínimos; clareza do que é guardado e por quê)  
   - **Base PRD**: 11 (RNF4)

16. `SPEC-016` — Métricas & Registros (consistência, rubricas, tendências, erros recorrentes)  
   - **Base PRD**: §9.2, §§8.1–8.5 (métricas), 10 (R5)

