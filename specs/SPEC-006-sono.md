# Feature Specification: Sono — Diário + Rotina Pré-sono + 1 Intervenção simples/semana

**Created**: 2026-02-19  
**PRD Base**: §6.1, §8.3, §5.3, §5.4, §9.1, §14, §11 (RNF1–RNF4), §10 (R2, R6)

## Caso de uso *(mandatory)*

O usuário quer melhorar sono de forma sustentável porque sono/energia são “infraestrutura” para consistência e aprendizagem (PRD §6.1). Esta feature define um fluxo Telegram-first que:
- captura um **diário mínimo** (rápido) para gerar tendência,
- ajuda o usuário a executar uma **rotina pré-sono realista** (pequena, sem perfeccionismo),
- propõe **1 intervenção simples por semana** (pequeno experimento),
- mantém robustez a dias ruins (Plano C / mínimo viável),
- e aplica um **gate leve de hábito**: “cumpriu sono” não significa “dormi perfeito”, e sim “registrei e executei o mínimo combinando” (alinhado a `SPEC-003`).

Esta SPEC descreve comportamentos observáveis; não define protocolo clínico (CBT-I completo), nem diagnóstico/tratamento médico.

## Scope & Non-goals *(recommended)*

**In scope (MVP slice)**:
- Diário do sono (manhã) + registro mínimo.
- Rotina pré-sono (noite) com passos pequenos e ajustáveis por tempo/energia.
- Uma intervenção semanal simples com adesão rastreável.
- Feedback curto e não punitivo baseado em tendências (não em “culpa”).

**Non-goals (agora)**:
- Não fazer prescrição médica, não diagnosticar insônia/OSA, não executar CBT‑I formal.
- Não exigir dispositivos/wearables.
- Não otimizar por métricas avançadas (p95, modelos, etc.).

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Registrar diário do sono em ≤ 60s e receber feedback curto (Priority: P1)

O usuário acorda e quer registrar rapidamente o essencial do sono para criar tendência sem burocracia.

**Why this priority**: Sem diário mínimo não há base para adaptação por evidência (PRD §6.1; §8.3; RNF1).

**Independent Test**:
- Simular 7 manhãs com registros completos e parciais.
- Validar que o sistema aceita registro rápido, calcula tendência simples e produz feedback curto e consultável.

**Acceptance Scenarios**:

1. **Scenario**: Diário completo ao acordar
   - **Given** é manhã e o usuário está começando o dia
   - **When** fornece horário aproximado que dormiu, horário que acordou, qualidade percebida (0–10) e energia pela manhã (0–10)
   - **Then** o sistema registra o diário, confirma em uma mensagem curta e retorna 1 insight simples (ex.: tendência ou próximo foco)

2. **Scenario**: Diário parcial não bloqueia (dia corrido)
   - **Given** o usuário tem pouco tempo
   - **When** fornece apenas horários (dormiu/acordou) ou apenas qualidade/energia
   - **Then** o sistema registra o que foi possível, marca como parcial e não usa linguagem punitiva

3. **Scenario**: Usuário pede para ver “como foi meu sono essa semana?”
   - **Given** existem registros de sono na semana
   - **When** o usuário solicita resumo semanal de sono
   - **Then** o sistema retorna um resumo curto com regularidade aproximada + médias de qualidade/energia (se houver), e 1 sugestão de próximo experimento

---

### User Story 2 — Executar uma rotina pré-sono pequena e realista (Priority: P1)

À noite, o usuário quer reduzir fricção e ter um “próximo passo” simples que aumente probabilidade de dormir melhor.

**Why this priority**: O PRD pede baixa fricção e robustez a dias ruins; a rotina pré-sono é o mecanismo diário mais direto (PRD §8.3; RNF1; RNF2).

**Independent Test**:
- Simular noite “normal” e noite “dia ruim”.
- Verificar que há um plano curto e executável, com critério observável do mínimo feito.

**Acceptance Scenarios**:

1. **Scenario**: Rotina pré-sono em dia normal
   - **Given** o usuário tem um horário-alvo (mesmo que aproximado)
   - **When** solicita a rotina da noite
   - **Then** o sistema fornece 2–4 passos pequenos (executáveis) e um critério claro do mínimo aceitável hoje

2. **Scenario**: Dia ruim — rotina mínima (Plano C) sem culpa
   - **Given** o usuário reporta baixa energia, estresse alto ou pouco tempo
   - **When** solicita “mínimo viável para dormir melhor hoje”
   - **Then** o sistema oferece uma rotina mínima (1–2 passos) e reforça que “mínimo feito” é sucesso do dia

3. **Scenario**: Usuário não consegue cumprir a rotina planejada
   - **Given** havia uma rotina sugerida
   - **When** o usuário diz que não conseguiu
   - **Then** o sistema registra a falha como dado (sem culpa), pergunta 1 coisa curta sobre o obstáculo e propõe um ajuste menor para a próxima noite

---

### User Story 3 — Fazer 1 intervenção simples por semana (experimento) e acompanhar adesão (Priority: P2)

O usuário quer testar uma mudança pequena por semana e ver se melhora regularidade/qualidade/energia.

**Why this priority**: Evita “lista de dicas” e cria aprendizado incremental (PRD §8.3; §5.5).

**Independent Test**:
- Simular 2 semanas com uma intervenção em cada semana.
- Validar que existe definição da intervenção, registro de adesão e reflexão curta no fim da semana.

**Acceptance Scenarios**:

1. **Scenario**: Intervenção semanal proposta com justificativa curta
   - **Given** existem dados mínimos da semana (mesmo parciais)
   - **When** inicia uma nova semana (ou solicita “experimento da semana”)
   - **Then** o sistema sugere 1 intervenção simples, explica o “por quê” em 1–2 frases e define como medir adesão

2. **Scenario**: Usuário rejeita intervenção sugerida
   - **Given** uma intervenção foi sugerida
   - **When** o usuário diz que não faz sentido
   - **Then** o sistema oferece alternativa mais simples (ou adia), mantendo tom não punitivo

3. **Scenario**: Fechamento semanal do experimento
   - **Given** a semana passou e houve alguma adesão (ou não)
   - **When** o usuário faz a revisão semanal (ou o sistema pergunta)
   - **Then** o sistema registra “funcionou / não funcionou / inconclusivo” e escolhe o próximo ajuste mínimo

## Edge Cases *(mandatory)*

- What happens when o usuário esquece o diário por 2+ dias?
  - O sistema oferece retomada com defaults e registra lacuna sem penalizar; sugere o próximo passo mínimo.
- What happens when horários são muito irregulares (ex.: turnos/viagem)?
  - O sistema evita moralismo, foca em tendência e propõe intervenção compatível com a restrição.
- What happens when o usuário reporta sofrimento significativo (ex.: ansiedade intensa, muitos despertares) e pede “tratamento”?
  - O sistema mantém escopo de produto, sugere buscar ajuda profissional quando apropriado e oferece apenas intervenções leves e seguras (sem protocolo clínico).
- What happens when o usuário não quer registrar dados por privacidade?
  - O sistema permite operar no mínimo com menos dados, explica impacto e aplica privacidade por padrão (ver `SPEC-015`).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST suportar registro de diário mínimo do sono com campos de produto (ex.: dormiu/acordou, qualidade percebida 0–10, energia pela manhã 0–10; despertares/opcional).
- **FR-002**: System MUST permitir diário parcial e registrar lacunas sem “penalizar” com linguagem punitiva.
- **FR-003**: System MUST oferecer uma rotina pré-sono com poucos passos, adaptável por tempo/energia (incluindo versão mínima para dia ruim).
- **FR-004**: System MUST definir um **gate leve** do que conta como “cumpriu sono hoje” (hábito/fundação), alinhado a `SPEC-003` (ex.: registro mínimo + execução do mínimo acordado).
- **FR-005**: System MUST propor 1 intervenção simples por semana e registrar adesão (feito/não feito/parcial) de forma curta.
- **FR-006**: System MUST fornecer um resumo semanal curto do sono (regularidade + qualidade/energia quando houver) e 1 recomendação acionável.
- **FR-007**: System MUST permitir consulta rápida do estado atual (ex.: “sono hoje”, “tendência da semana”), sem exigir múltiplas mensagens.

### Non-Functional Requirements

- **NFR-001**: System MUST manter fricção mínima: diário em ≤ 60s e rotinas em poucos passos (PRD RNF1).
- **NFR-002**: System MUST ser robusto a dias ruins: sempre existir “mínimo viável” (PRD RNF2).
- **NFR-003**: System MUST manter segurança psicológica: feedback firme, não punitivo, orientado a ajuste (PRD RNF3).
- **NFR-004**: System MUST aplicar privacidade por padrão: coletar o mínimo necessário e ser transparente sobre uso (PRD RNF4; ver `SPEC-015`).

### Key Entities *(include if feature involves data)*

- **SleepDiaryEntry**: data; dormiu_aprox; acordou_aprox; qualidade_percebida_0_10; energia_manha_0_10; despertares_opcional; status (completo/parcial).
- **SleepRoutine**: data/noite; passos_minimos; versão (normal/minima); status (planejada/feita/parcial).
- **WeeklySleepIntervention**: semana; descrição; objetivo; regra de adesão; status (aceita/recusada); adesão (contagem simples).
- **SleepWeeklySummary**: semana; regularidade_aprox; medias_qualidade_energia (se houver); principal_observacao; próximo_experimento.

## Acceptance Criteria *(mandatory)*

- O usuário consegue registrar diário do sono de forma rápida (≤ 60s) e o sistema aceita registros parciais.
- Em dia ruim, o sistema oferece rotina mínima e “gate mínimo” sem culpa.
- Existe uma definição observável do que conta como “cumpriu sono” (gate de hábito), alinhada com `SPEC-003`.
- O usuário consegue obter resumo semanal curto e sair com 1 experimento/intervenção para a semana.

## Business Objectives *(mandatory)*

- Tratar sono/energia como fundação que sustenta execução e aprendizagem (PRD §6.1).
- Reduzir fricção e aumentar consistência com planos mínimos para dias ruins (PRD RNF1/RNF2).
- Promover adaptação contínua baseada em evidência (PRD §5.5).

## Error Handling *(mandatory)*

- **Ausência de dados**: registrar lacuna, oferecer retomada e um próximo passo mínimo; sem punição.
- **Dados incoerentes/ambíguos**: pedir 1 confirmação curta e seguir com defaults conservadores.
- **Usuário some**: ao retornar, mostrar estado (o que falta hoje) e sugerir o menor passo viável.
- **Privacidade**: se usuário expressar desconforto, explicar minimização e permitir operar com menos dados (ver `SPEC-015`).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Aumento de regularidade do sono ao longo de semanas (menor variação de horários).
- **SC-002**: Melhora de tendência em qualidade percebida e/ou energia pela manhã ao longo de 2–4 semanas.
- **SC-003**: Alta taxa de dias com diário mínimo registrado (mesmo parcial) sem aumento de fricção percebida.
- **SC-004**: Aumento da taxa de “mínimo viável” cumprido em dias ruins, sem linguagem punitiva.