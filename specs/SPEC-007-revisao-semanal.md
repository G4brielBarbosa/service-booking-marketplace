# Feature Specification: Revisão Semanal — Painel mínimo + 3 decisões + alvos da semana

**Created**: 2026-02-19  
**PRD Base**: §§5.2, 5.5, 9.2, §6.2, §14, §10 (R5, R4, R7), §11 (RNF1–RNF3)

## Caso de uso *(mandatory)*

O usuário quer um ritual semanal curto (10–20 min) que transforme dados do dia a dia em **decisões acionáveis**, evitando “tentar tudo ao mesmo tempo” e evitando ilusões de progresso. Esta feature define um fluxo Telegram-first que:
- compila um **painel mínimo** com consistência, qualidade (rubricas/gates), gargalos e sono/energia,
- conduz o usuário a **3 decisões por meta/ciclo**: **manter / ajustar / pausar**,
- define **alvos da semana** (ex.: 1 de inglês, 1 de Java, 1 de sono/saúde) para foco,
- e gera um resultado claro: “o que muda na próxima semana” (sem virar burocracia).

A revisão semanal deve manter **segurança psicológica** (tom não punitivo), ser robusta a semanas com poucos dados, e respeitar limites de metas em paralelo (PRD §6.2).

> Importante: esta SPEC define o fluxo e os outputs observáveis; não desenha dashboards técnicos e não define armazenamento/arquitetura. As métricas/registro base vêm de `SPEC-016`.

## Scope & Non-goals *(recommended)*

**In scope (MVP slice)**:
- Painel semanal mínimo (texto/conversacional) e síntese.
- Condução das 3 decisões e registro do resultado.
- Seleção de alvos da semana e compromisso.
- Saída em formato “plano semanal” de alto nível (sem detalhar plano diário; isso é `SPEC-002`).

**Non-goals (agora)**:
- Não implementar backlog inteligente completo (`SPEC-008`) nem detecção avançada de falso progresso (`SPEC-009`) — apenas mostrar sinais observáveis disponíveis.
- Não criar análises complexas/estatísticas; foco em “insight suficiente para decidir”.
- Não impor psicometria nem linguagem clínica.

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Completar revisão semanal e sair com decisões e alvos (Priority: P1)

O usuário quer concluir a revisão semanal em 10–20 minutos e terminar sabendo:
- como foi a semana (consistência/qualidade/sono-energia),
- o que manter/ajustar/pausar,
- e quais são os alvos claros da próxima semana.

**Why this priority**: É parte do MVP (PRD §14) e concretiza adaptação contínua baseada em evidência (PRD §5.5; §9.2; R5).

**Independent Test**:
- Simular uma semana com dados: check-ins, planos, gates, rubricas (quando houver), sono/energia.
- Rodar o fluxo e validar que o output final inclui: painel mínimo + decisões registradas + alvos definidos + mudanças explícitas para a próxima semana.

**Acceptance Scenarios**:

1. **Scenario**: Revisão semanal com dados suficientes
   - **Given** existe histórico da semana (mesmo que com alguns dias faltantes)
   - **When** o usuário inicia “revisão semanal”
   - **Then** o sistema apresenta um painel mínimo com: consistência por meta, qualidade (gates/rubricas), gargalos principais e sono/energia (tendência simples)

2. **Scenario**: Usuário toma 3 decisões (manter/ajustar/pausar)
   - **Given** o painel foi apresentado
   - **When** o usuário escolhe para cada meta relevante “manter”, “ajustar” ou “pausar”
   - **Then** o sistema registra as decisões e explica, de forma curta, o efeito prático de cada decisão na próxima semana

3. **Scenario**: Definir alvos da semana (foco)
   - **Given** decisões foram registradas
   - **When** o usuário define os alvos da semana
   - **Then** o sistema registra 1–3 alvos máximos (ex.: 1 inglês, 1 Java, 1 sono/saúde) e descreve como eles se manifestam no dia a dia (em linguagem de produto, sem tarefas técnicas)

4. **Scenario**: Encerramento com “resumo executável”
   - **Given** a revisão terminou
   - **When** o usuário pede “resumo final”
   - **Then** o sistema retorna uma mensagem única com: principais aprendizados da semana + decisões + alvos + uma regra simples de foco (“o que não fazer”/limite)

---

### User Story 2 — Semana com poucos dados (usuário ausente) não vira bronca (Priority: P1)

O usuário teve uma semana caótica e quase não registrou nada. Ele quer que o sistema:
- não culpe,
- mostre o que dá para inferir,
- e saia com um plano mínimo para recuperar baseline.

**Why this priority**: Robustez a falhas reais e segurança psicológica sustentam consistência (PRD RNF2; RNF3).

**Independent Test**:
- Simular semana com 1–2 dias de dados.
- Verificar que o sistema executa revisão “parcial” e propõe uma recuperação mínima.

**Acceptance Scenarios**:

1. **Scenario**: Revisão semanal parcial explícita
   - **Given** a semana tem poucos registros
   - **When** o usuário inicia revisão semanal
   - **Then** o sistema marca como “revisão parcial”, lista o que está desconhecido e evita conclusões fortes

2. **Scenario**: Recuperação mínima para próxima semana
   - **Given** revisão parcial
   - **When** o usuário quer “voltar pros trilhos”
   - **Then** o sistema define 1–2 ações mínimas de coleta/execução para a próxima semana (ex.: check-in diário mínimo + 1 meta intensiva ativa), sem aumentar escopo

---

### User Story 3 — Detectar sinais de overload e sugerir redução de escopo (Priority: P2)

O usuário pode ter tentado metas demais. Na revisão, ele quer ajuda para reduzir escopo sem sensação de fracasso.

**Why this priority**: O PRD impõe limites e prevê overload como risco (PRD §6.2; §13).

**Independent Test**:
- Simular semana com múltiplas metas e baixa consistência/energia.
- Verificar que o sistema sugere pausa/ajuste e respeita limite de 2 metas intensivas.

**Acceptance Scenarios**:

1. **Scenario**: Sinais de overload aparecem no painel e viram sugestão
   - **Given** consistência baixa em múltiplas metas e energia média baixa na semana
   - **When** o sistema apresenta o painel
   - **Then** ele sugere reduzir escopo (pausar uma meta intensiva ou reduzir carga), com linguagem protetiva e não punitiva

2. **Scenario**: Usuário escolhe pausar e o sistema preserva identidade mínima
   - **Given** o usuário escolhe “pausar” uma meta intensiva
   - **When** confirma a decisão
   - **Then** o sistema registra pausa e garante que permanece ao menos uma meta de fundação ativa (sono/saúde) com mínimo viável

## Edge Cases *(mandatory)*

- What happens when o usuário tenta definir muitos alvos da semana?
  - Sistema impõe limite (ex.: máximo 3; preferencialmente 1 por domínio) e explica que foco protege consistência.
- What happens when há conflito entre dados e percepção do usuário (“parece que fui bem, mas dados mostram pouca consistência”)?
  - Sistema valida a percepção, mostra os dados com cuidado e propõe um experimento simples para a próxima semana.
- What happens when metas foram pausadas no meio da semana?
  - Sistema mostra claramente o período ativo vs pausado e não “penaliza” consistência fora do período ativo.
- What happens when o usuário quer rever “o que eu já fiz” da semana em detalhe?
  - Sistema oferece um resumo por dia/por meta em forma curta (sem virar relatório longo) e permite drill-down em 1 pergunta por vez.
- What happens when o usuário não escolhe decisões (evita decidir)?
  - Sistema propõe defaults conservadores (ex.: manter fundação + 1 meta intensiva) e pede apenas 1 confirmação curta.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST compilar um **painel semanal mínimo** contendo, no mínimo:
  - consistência por meta (dias/semana),
  - qualidade (gates/rubricas quando aplicável),
  - gargalos principais (tempo/energia/atrito),
  - sono/energia (tendência simples).
  (Fonte de métricas: `SPEC-016`)

- **FR-002**: System MUST conduzir o usuário a **3 decisões** por meta/ciclo: **manter / ajustar / pausar**, registrando a escolha.
- **FR-003**: System MUST permitir definir **alvos da semana** (limite máximo) e registrar compromisso.
- **FR-004**: System MUST produzir um **resumo final único** com: principais achados + decisões + alvos + regra de foco para a semana.
- **FR-005**: System MUST suportar revisão semanal **parcial** quando dados forem insuficientes, explicitando limitações e definindo recuperação mínima.
- **FR-006**: System MUST respeitar governança de metas em paralelo (PRD §6.2): no máximo 2 metas intensivas ativas por ciclo; quando violado, a revisão deve orientar escolha (pausar/adiar).
- **FR-007**: System MUST registrar recomendações/ajustes propostos para que possam influenciar o planejamento da próxima semana (sem definir como isso é implementado).

### Non-Functional Requirements

- **NFR-001**: System MUST manter a revisão concluível em 10–20 minutos, com interação curta e previsível (PRD §9.2; RNF1).
- **NFR-002**: System MUST manter segurança psicológica: tom não punitivo, foco em tendências e ajustes (PRD RNF3).
- **NFR-003**: System MUST ser robusto a semanas ruins com poucos dados, evitando “bronca” e propondo mínimo viável (PRD RNF2).

### Key Entities *(include if feature involves data)*

- **WeeklyPanel**: período; consistência_por_meta; sinais_de_qualidade; gargalos; sono_energia_tendencia; avisos_de_dados_faltantes.
- **WeeklyDecision**: meta; decisão (manter/ajustar/pausar); justificativa curta; impacto esperado (texto curto).
- **WeeklyTargets**: lista de alvos (limitada); descrição observável; motivo.
- **WeeklyReviewResult**: painel; decisões; alvos; resumo_final; status (completa/parcial).

## Acceptance Criteria *(mandatory)*

- A revisão semanal sempre resulta em:
  - painel mínimo,
  - decisões registradas (manter/ajustar/pausar),
  - alvos da semana definidos,
  - e um resumo final único e acionável.
- A revisão funciona mesmo com poucos dados (marca como parcial e define recuperação mínima).
- O fluxo respeita limites de metas intensivas e reduz escopo quando há sinais de overload.
- Tom não punitivo e fricção baixa são mantidos.

## Business Objectives *(mandatory)*

- Transformar dados em decisões acionáveis (adaptação contínua) (PRD §5.5; §9.2).
- Proteger consistência evitando overload e reforçando foco (PRD §6.2; §13).
- Sustentar motivação com feedback firme e seguro (PRD RNF3).

## Error Handling *(mandatory)*

- **Poucos dados**: marcar revisão parcial; evitar inferências fortes; propor 1–2 passos mínimos para recuperar baseline.
- **Entrada ambígua**: fazer apenas 1 pergunta de clarificação por vez e seguir com defaults conservadores.
- **Overload**: orientar redução (pausar/ajustar) com linguagem protetiva; não aumentar escopo.
- **Usuário some**: permitir retomar a revisão do ponto em que parou, exibindo o próximo passo único.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Percentual de semanas em que a revisão é completada (completa ou parcial).
- **SC-002**: Aumento de consistência/qualidade após 2–4 ciclos em metas-alvo definidas.
- **SC-003**: Redução de semanas com overload percebido após adoção de limites e alvos da semana.