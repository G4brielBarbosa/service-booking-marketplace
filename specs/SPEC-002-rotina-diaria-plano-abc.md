# Feature Specification: Rotina Diária (Telegram-first) — Check-in + Plano A/B/C + Execução Guiada

**Created**: 2026-02-17  
**PRD Base**: §§5.2, 5.3, 9.1, 10 (R1, R6), 11 (RNF1, RNF2, RNF3), §2 (consultar steps do dia), §6.2, §14

## Caso de uso *(mandatory)*

O usuário quer executar suas metas anuais com consistência, mas tem variação diária de tempo/energia e risco de overload. Esta feature define o fluxo diário (Telegram-first como **UX conversacional**) que:

- Coleta um **check-in mínimo** (tempo disponível e estado/energia) (PRD §9.1).
- Retorna um **Plano A/B/C** executável com prioridades claras, adequado para dias bons e ruins (PRD §5.2; §9.1; R1; RNF2).
- Entrega **execução guiada** com instruções objetivas por tarefa para reduzir fricção e carga cognitiva (PRD §5.3; R6; RNF1).
- Permite **retomar/replanejar** quando o contexto muda (tempo/energia caem, interrupções) sem tom punitivo (PRD §5.2; §13; RNF3).
- Permite consultar os **steps do dia atual** e o estado do plano (PRD §2).

Esta SPEC descreve **o que** a rotina diária deve produzir/permitir observar. Não define stack, comandos, integrações, banco, ou detalhes de implementação.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Check-in mínimo e receber um plano A/B/C executável (Priority: P1)

O usuário quer abrir o Telegram e, em poucos passos, informar quanto tempo e quanta energia tem hoje. Com base nisso, ele quer receber um plano do dia **executável** com:

- 1 prioridade absoluta,
- 1–2 tarefas complementares,
- 1 tarefa de fundação (quando aplicável ao ciclo),

e cada tarefa com instruções objetivas (o que fazer, por quanto tempo e como saber que “está feito”).

**Why this priority**: É o slice mínimo do MVP que materializa planejamento diário com Plano A/B/C e baixa fricção (PRD §14; §§5.2–5.3; §9.1; R1; R6; RNF1–RNF3). Sem isso, o produto não consegue guiar execução diária de forma consistente.

**Independent Test**: Em um ambiente de teste, simular um dia sem depender de revisão semanal:

- executar check-in com (a) muito tempo/energia, (b) tempo/energia médios, (c) pouco tempo/energia;
- verificar que o sistema retorna um Plano A/B/C válido, com a estrutura e instruções exigidas;
- verificar que o plano é consultável e que o usuário consegue iniciar pela prioridade absoluta.

**Acceptance Scenarios**:

1. **Scenario**: Check-in completo e plano recomendado para um dia bom
   - **Given** o usuário informa tempo suficiente e energia adequada (PRD §9.1)
   - **When** realiza o check-in diário
   - **Then** o sistema retorna um plano com prioridade absoluta, 1–2 complementares e uma tarefa de fundação (quando aplicável), com instruções objetivas e estimativa de duração por tarefa (PRD §9.1; §5.3)

2. **Scenario**: Check-in com pouco tempo/energia aciona um Plano C (MVD)
   - **Given** o usuário informa pouco tempo e energia baixa (PRD §9.1; RNF2)
   - **When** realiza o check-in diário
   - **Then** o sistema retorna um Plano C (MVD) que cabe no tempo disponível e mantém consistência, com linguagem não punitiva (PRD RNF2; RNF3)

3. **Scenario**: Check-in mínimo ainda gera um plano válido
   - **Given** o usuário informa apenas tempo disponível e energia (PRD §9.1)
   - **When** realiza o check-in diário
   - **Then** o sistema gera um plano executável sem bloquear por ausência de campos opcionais (PRD RNF1)

4. **Scenario**: O usuário pede “por onde eu começo?”
   - **Given** existe um plano do dia gerado
   - **When** o usuário pede orientação de início
   - **Then** o sistema aponta a prioridade absoluta e descreve o primeiro passo em 1–2 frases (PRD §5.3; RNF1)

---

### User Story 2 - Replanejar no meio do dia e retomar após interrupções (Priority: P2)

O usuário começa com um plano, mas o contexto muda. Ele quer informar o novo tempo/energia (ou o tempo restante) e receber um ajuste do plano para continuar avançando sem frustração.

**Why this priority**: A vida real exige adaptação contínua (PRD §5.2) e o PRD explicitamente assume dias ruins e mudanças de contexto (PRD §4; RNF2). Replanejar sem culpa reduz rigidez e risco de abandono (PRD §13; RNF3).

**Independent Test**: Gerar um Plano A pela manhã; depois simular queda de tempo/energia; validar que o sistema oferece um plano ajustado (B/C) e preserva o registro do que já foi feito.

**Acceptance Scenarios**:

1. **Scenario**: Replanejamento por tempo restante
   - **Given** o usuário recebeu um plano mais extenso no início do dia
   - **When** informa que agora tem menos tempo do que o previsto
   - **Then** o sistema ajusta o plano para caber no novo tempo, explicando a mudança de forma concisa e mantendo a prioridade (PRD §5.2; RNF1)

2. **Scenario**: Retomada após sumiço
   - **Given** o usuário não interagiu por horas após receber o plano
   - **When** volta e pergunta “o que falta hoje?”
   - **Then** o sistema mostra o estado do dia (feito/faltando) e sugere o próximo passo mais viável (PRD §5.3; RNF2)

---

### User Story 3 - Consultar “plano de hoje” e “o que já fiz hoje” (Priority: P2)

O usuário quer um jeito simples de consultar o que estava planejado para o dia e o que já foi realizado, sem depender de memória.

**Why this priority**: Resolve diretamente a necessidade do PRD de consultar “steps do dia atual” (PRD §2) e reduz carga cognitiva (PRD RNF1).

**Independent Test**: Em um dia com ao menos uma tarefa marcada como concluída, validar que o usuário consegue pedir um resumo do plano e um resumo do que já foi feito hoje.

**Acceptance Scenarios**:

1. **Scenario**: Ver o plano do dia
   - **Given** existe um plano gerado para o dia atual
   - **When** o usuário solicita “meu plano de hoje”
   - **Then** o sistema apresenta o plano atual (A/B/C) com prioridade, complementares e fundação, indicando duração e critérios observáveis de “feito” por tarefa (PRD §9.1; §5.3)

2. **Scenario**: Ver os steps concluídos
   - **Given** ao menos uma tarefa do plano foi registrada como concluída
   - **When** o usuário pergunta “o que eu já fiz hoje?”
   - **Then** o sistema lista os steps concluídos e o que falta, de forma curta e clara (PRD §2; RNF1)

### Edge Cases *(mandatory)*

- What happens when o usuário **não informa tempo nem energia** no check-in? (PRD §9.1)
- What happens when o usuário informa **tempo** mas responde energia de forma ambígua (ex.: “tanto faz”, “meh”)? O sistema deve pedir uma única confirmação ou assumir defaults? (PRD RNF1) **[NEEDS CLARIFICATION]**: defaults oficiais para energia.
- What happens when o usuário pede **muitas tarefas** em um dia ruim? O sistema deve impor um limite duro (ex.: apenas Plano C + 1 opcional) ou negociar? (PRD RNF2; §6.2) **[NEEDS CLARIFICATION]**.
- What happens when o usuário tem **0 minutos** e quer “não quebrar a sequência”? Qual é o menor step observável aceitável para contar como consistência do dia? (PRD RNF2; §5.1 MVD) **[NEEDS CLARIFICATION]**.
- What happens when o usuário tenta ativar **mais de 2 metas intensivas** simultaneamente e exige que todas entrem no plano diário? (PRD §6.2) O sistema deve bloquear ou apenas recomendar? **[NEEDS CLARIFICATION]**.
- How does the system handle **restrições de ambiente** (ex.: “não posso gravar áudio”, “sem privacidade agora”) que afetam a executabilidade de tarefas? A rotina diária precisa oferecer substituição equivalente ou apenas marcar como bloqueado? (PRD RNF1) **[NEEDS CLARIFICATION]**.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST coletar no check-in diário, no mínimo: **tempo disponível** e **energia** (PRD §9.1).
- **FR-002**: System MUST permitir check-in rápido (campos mínimos) sem impedir a geração do plano (PRD RNF1).
- **FR-003**: System MUST selecionar e comunicar explicitamente um **Plano A, B ou C** com base nas entradas do check-in (PRD §9.1; R1; RNF2). **[NEEDS CLARIFICATION]**: regras/limiares para classificar A vs B vs C a partir de tempo (5/15/30/60+) e energia (0–10).
- **FR-004**: System MUST retornar um plano diário com a estrutura do PRD: **1 prioridade absoluta**, **1–2 complementares** e **1 tarefa de fundação** (quando aplicável) (PRD §9.1).
- **FR-005**: System MUST descrever cada tarefa do plano com:
  - **o que fazer** (descrição operacional em linguagem natural),
  - **por quanto tempo** (bloco/tempo),
  - **critério observável de “feito”** (PRD §5.3; R6).
- **FR-006**: System MUST suportar replanejamento intradiário quando o usuário atualizar tempo/energia/estado, retornando um plano ajustado (PRD §5.2; §9.1).
- **FR-007**: System MUST permitir retomada após interrupção, mostrando o estado do dia (feito/faltando) e sugerindo o próximo passo viável no tempo restante (PRD §5.3; RNF2).
- **FR-008**: System MUST manter um registro do dia atual com: check-in, plano gerado e status por tarefa (planejada/em progresso/concluída/bloqueada/adiada) para suportar consulta de steps (PRD §2).
- **FR-009**: System MUST permitir ao usuário consultar “meu plano de hoje” e “o que já fiz hoje” no dia atual (PRD §2).
- **FR-010**: System MUST lidar com sobrecarga de metas em paralelo: quando o usuário tentar incluir metas intensivas demais, o sistema MUST sinalizar o limite e orientar uma escolha (PRD §6.2; §13). **[NEEDS CLARIFICATION]**: se a rotina diária pode bloquear a geração do plano até o usuário escolher ou se apenas recomenda.
- **FR-011**: System MUST manter o fluxo diário curto e previsível (PRD RNF1; §9.1).

### Non-Functional Requirements

- **NFR-001**: System MUST manter simplicidade e baixa carga cognitiva: interação curta e previsível (PRD RNF1; §5.3).
- **NFR-002**: System MUST ser robusto a dias ruins: sempre existir um Plano C (MVD) executável no pior cenário (PRD RNF2; §4).
- **NFR-003**: System MUST manter segurança psicológica: feedback firme, orientado a processo e ajuste, sem humilhação/punição (PRD RNF3; §3).
- **NFR-004**: System MUST operar com privacidade por padrão no registro do dia: coletar o mínimo e ser claro sobre o que foi guardado e por quê (PRD RNF4). **[NEEDS CLARIFICATION]**: política de retenção (quais dados do “dia” ficam guardados e por quanto tempo).

### Key Entities *(include if feature involves data)*

- **DailyCheckIn**: data; tempo disponível (ex.: 5/15/30/60+); energia (0–10); humor/estresse (opcional); restrições relevantes quando fornecidas; timestamp (PRD §9.1).
- **DailyPlan**: data; plano (A/B/C); lista de tarefas; prioridade absoluta; complementares; fundação; justificativa curta do plano escolhido; timestamp (PRD §9.1).
- **PlannedTask**: título; meta associada (quando aplicável); duração/bloco; instruções objetivas; critério de “feito”; status (planejada/em progresso/concluída/bloqueada/adiada); observações curtas (PRD §5.3; §2).
- **DailyStepsSummary**: steps concluídos; pendências; breve síntese do estado do dia (PRD §2).

## Acceptance Criteria *(mandatory)*

- Com um check-in mínimo (tempo + energia), o usuário sempre recebe um plano executável (A/B/C) com prioridades claras (PRD §9.1; R1).
- Todo plano diário contém: 1 prioridade absoluta, 1–2 complementares e 1 tarefa de fundação (quando aplicável) (PRD §9.1).
- Cada tarefa do plano inclui instruções objetivas e critério observável de “feito” (PRD §5.3; R6).
- O usuário consegue consultar no dia atual: “meu plano de hoje” e “o que já fiz hoje” (PRD §2).
- Em dia ruim, existe Plano C (MVD) com linguagem não punitiva (PRD RNF2; RNF3).

## Business Objectives *(mandatory)*

- **Consistência com baixa fricção**: tornar o “próximo passo” óbvio e executável diariamente (PRD §§5.2–5.3; §9.1; R6).
- **Adaptação contínua**: permitir replanejamento sem culpa, reduzindo frustração/abandono (PRD §5.2; §13; RNF3).
- **Proteção contra overload**: reforçar limites de metas em paralelo e priorização realista (PRD §6.2; §13).
- **Clareza diária**: permitir consultar steps do dia e estado atual para reduzir carga mental (PRD §2; RNF1).

## Error Handling *(mandatory)*

- **Entrada ausente/ambígua**: se tempo/energia faltarem, o sistema MUST aplicar defaults mínimos e pedir **1** confirmação curta para destravar (PRD RNF1; catálogo do `SPEC-GUIDE`).
- **Dia ruim**: oferecer Plano C (MVD) imediatamente e reforçar que “consistência mínima” é sucesso do dia, sem culpa (PRD RNF2; RNF3).
- **Mudança de contexto**: permitir replanejar por tempo/energia e retornar um plano ajustado conciso (PRD §5.2; §9.1).
- **Usuário some**: ao retornar, mostrar estado (feito/faltando) e sugerir o próximo passo viável (PRD §5.3).
- **Sobrecarga**: reduzir escopo e orientar escolha/pausa de metas intensivas, mantendo tom protetivo (PRD §6.2; RNF3). **[NEEDS CLARIFICATION]**: se há bloqueio obrigatório quando o limite é violado.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Usuário conclui check-in e recebe plano em \( \le \) 2 minutos de interação total (PRD §9.1; RNF1).
- **SC-002**: Em dias de baixa energia/tempo, a taxa de dias com Plano C (MVD) concluído aumenta ao longo de semanas (PRD RNF2).
- **SC-003**: A taxa de dias com “plano consultável” (usuário consegue ver plano e steps do dia) tende a alta e reduz relatos de “me perdi no dia” (PRD §2; RNF1). **[NEEDS CLARIFICATION]**: como o PRD deseja medir “me perdi no dia” (pergunta explícita? proxy?).
