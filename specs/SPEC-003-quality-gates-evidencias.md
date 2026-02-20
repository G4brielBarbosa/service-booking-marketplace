# Feature Specification: Quality Gates & Evidência Mínima (aprendizagem e hábitos)

**Created**: 2026-02-17  
**PRD Base**: §5.4, §§8.1–8.3, §§5.2–5.3, §9.1, §10 (R2, R3), §11 (RNF1–RNF4), §13, §14

## Caso de uso *(mandatory)*

O usuário quer avançar em metas de **aprendizagem/competência** (Inglês e Java) e **hábitos/fundação** (sono/saúde) sem cair em “falso progresso” (PRD §1, §5.4; R3). Para isso, o produto precisa definir, de forma **objetiva e testável**, quando uma tarefa:

- **conta como concluída** (há evidência mínima adequada ao objetivo),
- **não conta como concluída** (faltou evidência, evidência inválida ou incompleta),
- e qual é o **caminho mais curto** para “consertar” a conclusão (sem burocracia) (PRD RNF1; §5.4).

Esta SPEC define um mecanismo cross-cutting de **Quality Gates** (PRD §5.4; R2) que:

- aplica **fricção proporcional** (aprendizagem exige evidência mais forte; hábitos exigem registro mínimo) (PRD §5.4; RNF1),
- é **robusto a dias ruins** com alternativas mínimas (MVD/Plano C) sem “passar pano” para aprendizagem (PRD RNF2; §9.1),
- mantém **segurança psicológica**: feedback firme, não punitivo, orientado a processo e ajuste (PRD RNF3; §3),
- respeita **privacidade por padrão**, deixando claro o que é guardado e por quê (PRD RNF4).

> Importante: esta SPEC não define stack, arquitetura, banco, integração, formatos de arquivo, nem mecanismos técnicos. Ela define apenas comportamentos e resultados observáveis.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Concluir tarefa de aprendizagem somente com evidência mínima (Priority: P1)

O usuário quer “marcar como feito” uma tarefa de aprendizagem (Inglês/Java), mas o produto deve impedir conclusão sem evidência mínima, para evitar “falso progresso” (PRD §5.4; R2).

**Why this priority**: É o coração do MVP para qualidade (PRD §14) e a principal mitigação do risco “fiz mas não aprendi” (PRD §13; R3). Sem isso, o produto vira lista de tarefas e perde o diferencial.

**Independent Test**: Criar uma tarefa de aprendizagem com gate definido, tentar concluí-la sem evidência, com evidência parcial e com evidência válida, verificando (a) bloqueio, (b) orientações de correção, (c) registro do resultado do gate, sem depender de revisão semanal.

**Acceptance Scenarios**:

1. **Scenario**: Bloqueio de conclusão sem evidência mínima
   - **Given** existe uma tarefa de aprendizagem do dia com um quality gate definido (PRD §5.4)
   - **When** o usuário tenta marcar a tarefa como concluída sem fornecer evidência
   - **Then** o sistema não aceita a conclusão e oferece o caminho mais curto para produzir evidência mínima, em linguagem curta e não punitiva (PRD R2; RNF1; RNF3)

2. **Scenario**: Evidência enviada, porém inválida/ilegível/inaudível
   - **Given** a tarefa requer evidência (PRD §5.4)
   - **When** o usuário envia evidência vazia, ilegível ou incompreensível
   - **Then** o sistema não aceita a conclusão, explica por que a evidência não é válida e solicita reenvio ou alternativa equivalente **se** a SPEC da feature permitir (PRD RNF3)

3. **Scenario**: Evidência parcial não satisfaz o gate
   - **Given** a tarefa exige um conjunto mínimo de evidências (ex.: múltiplos itens) (PRD §5.4)
   - **When** o usuário fornece apenas parte do mínimo
   - **Then** o sistema mantém o status como “não concluída”, destaca o que falta e oferece o menor passo adicional para completar o gate (PRD RNF1)

4. **Scenario**: Evidência válida satisfaz o gate e a tarefa conta como concluída
   - **Given** a tarefa de aprendizagem está pendente e possui um gate definido (PRD §5.4)
   - **When** o usuário fornece evidência mínima válida conforme o gate
   - **Then** o sistema aceita a conclusão e registra que o gate foi satisfeito (PRD R2)

---

### User Story 2 - Concluir tarefa de hábito/fundação com fricção mínima e sem culpa (Priority: P1)

O usuário quer manter consistência em hábitos de fundação (ex.: sono) mesmo em dia ruim. Ele precisa de um gate “leve”, baseado em registro mínimo e aderência proporcional, que não o faça desistir por burocracia.

**Why this priority**: O PRD prioriza sono/energia como infraestrutura (PRD §6.1) e exige robustez a dias ruins (RNF2). Sem gates proporcionais para hábitos, o usuário pode abandonar por excesso de fricção.

**Independent Test**: Simular um dia ruim e um dia normal com uma tarefa de hábito; verificar que o sistema aceita conclusão com evidência mínima apropriada (registro curto), sem exigir evidência “pesada”, e com feedback protetivo.

**Acceptance Scenarios**:

1. **Scenario**: Gate mínimo para hábito em dia normal
   - **Given** existe uma tarefa de hábito/fundação com gate definido como “registro mínimo” (PRD §5.4; §8.3)
   - **When** o usuário fornece o registro mínimo esperado
   - **Then** o sistema aceita a conclusão e registra o essencial para tendência, mantendo a interação curta (PRD RNF1)

2. **Scenario**: Dia ruim aciona versão mínima do gate (MVD) para hábito
   - **Given** o usuário reportou pouco tempo/energia (PRD §9.1; RNF2)
   - **When** tenta cumprir a tarefa de fundação
   - **Then** o sistema oferece um gate mínimo (MVD) que preserve consistência e aceita a conclusão quando o mínimo é atendido, com linguagem não punitiva (PRD RNF2; RNF3)

---

### User Story 3 - Garantir transparência do “por que não contou” e como recuperar (Priority: P2)

Quando o gate bloqueia a conclusão, o usuário quer entender rapidamente o motivo e o passo mais curto para resolver, sem se sentir punido ou confuso.

**Why this priority**: Reduz frustração e rigidez (PRD §13) e sustenta simplicidade (RNF1) e segurança psicológica (RNF3).

**Independent Test**: Forçar falha de gate em diferentes motivos (ausência, evidência inválida, parcial) e validar que o sistema explica claramente e oferece uma recuperação curta.

**Acceptance Scenarios**:

1. **Scenario**: Mensagem de bloqueio é curta e acionável
   - **Given** o gate falhou por ausência ou incompletude de evidência
   - **When** o usuário tenta concluir
   - **Then** o sistema explica o motivo em 1–2 frases e dá um próximo passo único e simples (PRD RNF1)

2. **Scenario**: Feedback firme sem humilhação
   - **Given** o usuário falhou o gate repetidas vezes
   - **When** tenta concluir novamente
   - **Then** o sistema mantém tom firme e de aprendizado (não punitivo), reforçando que o objetivo é progresso real e oferecendo um caminho mínimo viável (PRD RNF3)

---

### User Story 4 - Detectar “falso progresso” e acionar reforço observável (Priority: P2)

O usuário quer que o sistema identifique sinais de “falso progresso” (ex.: repetição do mesmo erro ou queda persistente de compreensão) e reaja com reforço pequeno e verificável, em vez de só aceitar conclusões.

**Why this priority**: Implementa R3 (detecção de falhas reais) e mitiga risco central do PRD (PRD §1; §13).

**Independent Test**: Com um conjunto de tarefas concluídas e registros de erro recorrente, verificar que o sistema aciona um reforço e exige evidência mínima desse reforço.

**Acceptance Scenarios**:

1. **Scenario**: Erro recorrente dispara reforço
   - **Given** o usuário registrou o mesmo erro recorrente em múltiplas ocasiões recentes (PRD §5.4; R3; §8.1–§8.2)
   - **When** conclui uma tarefa relacionada ao domínio
   - **Then** o sistema solicita um reforço curto ligado ao erro e registra evidência mínima da tentativa de reforço (PRD R3)

2. **Scenario**: Reforço não vira burocracia
   - **Given** o usuário está com pouco tempo
   - **When** o reforço é solicitado
   - **Then** o sistema oferece uma versão mínima do reforço que ainda seja observável e útil, evitando sobrecarga (PRD RNF1; RNF2)

---

### User Story 5 - Privacidade por padrão na evidência (Priority: P3)

O usuário quer sentir segurança ao fornecer evidências (especialmente áudio/texto sensível). Ele quer clareza do que é guardado, por quanto tempo e como pedir remoção/opt-out, sem perder usabilidade.

**Why this priority**: Privacidade por padrão é requisito explícito (PRD RNF4). Evidências podem ser sensíveis.

**Independent Test**: No fluxo de envio de evidência, validar que o sistema explicita o que guarda e por quê, e oferece alternativa quando aplicável.

**Acceptance Scenarios**:

1. **Scenario**: Transparência de coleta e armazenamento mínimo
   - **Given** uma tarefa exige evidência
   - **When** o usuário pergunta “o que você guarda?”
   - **Then** o sistema explica claramente quais dados mínimos são guardados para medir progresso e por quê (PRD RNF4)

2. **Scenario**: Evidência sensível tem alternativa quando definida
   - **Given** o usuário não pode compartilhar evidência sensível no momento
   - **When** pede alternativa
   - **Then** o sistema oferece alternativa equivalente **apenas** quando a alternativa preserva o objetivo de aprendizagem/hábito (equivalência global definida nesta SPEC); caso contrário, marca a situação como **bloqueada por contexto** e sugere o mínimo possível para manter consistência (PRD RNF4; RNF2).

### Edge Cases *(mandatory)*

- What happens when o usuário quer “marcar como feito mesmo assim” (bypass explícito do gate)? (PRD §5.4; RNF3)
- How does system handle evidência parcial por falta de tempo (ex.: 1/3 itens)?
  - Para **aprendizagem/competência**: registrar como **tentativa** (não concluída), manter o gate como “não satisfeito”, e indicar exatamente o que faltou. Tentativa **não conta** como consistência da meta intensiva, mas pode contar como “micro-consistência do dia” quando aplicável (ver `SPEC-002`).
  - Para **hábito/fundação**: se o registro mínimo do hábito foi atendido, pode contar como concluído (gate proporcional).
- What happens when o usuário tem **0 minutos** e só quer “não quebrar a sequência”?
  - Definir **micro-consistência** (≤ 30–60s) como: 1 micro-step observável.
  - Para **aprendizagem**: micro-consistência pode ser 1 item de retrieval (1 pergunta/1 item), registrado como **MVD micro/tentativa**, sem contar como tarefa concluída com qualidade.
  - Para **hábito/fundação**: o menor registro (ex.: diário mínimo) pode satisfazer o gate proporcional e contar como concluído.
- What happens when o usuário não pode enviar áudio por ambiente/privacidade?
  - Para tarefas cujo objetivo é **speaking**: **não existe equivalência plena** ao áudio por padrão; o sistema marca a tarefa como **bloqueada por contexto** (sem culpa) e oferece alternativa que preserve consistência (ex.: retrieval curto), registrada como tentativa/substituição (não como conclusão de speaking).
  - Para tarefas cujo objetivo é **compreensão/retrieval**: respostas em texto são evidência equivalente válida.
- How does system handle evidência fraudável (ex.: “sim, fiz” sem comprovação)? O sistema deve manter a exigência de evidência mínima e oferecer alternativa observável (PRD §5.4).
- What happens when o usuário some no meio do gate (abandona conversa) e volta depois? O sistema deve permitir retomar a partir do que faltava (PRD §5.3; RNF2).
- How does system handle quando o usuário pede para “consultar os steps do dia atual” e quer ver o que foi aceito/rejeitado pelos gates (PRD §2)?
- What happens when uma SPEC de feature define gates que geram fricção excessiva para o usuário (contradiz RNF1)?
  - O sistema deve tratar como “gate pesado demais” quando houver repetidas falhas por falta de evidência **ou** quando o usuário indicar fricção alta (ver SC-005).
  - Resposta: oferecer o **caminho mínimo** que ainda preserve o objetivo (reduzir itens, reduzir tamanho da evidência, manter 1 único critério) **ou** marcar como bloqueado e levar ajuste para revisão semanal (`SPEC-007`) quando não houver como preservar objetivo sem descaracterizar.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST classificar tarefas em, no mínimo, **(a) aprendizagem/competência** e **(b) hábito/fundação**, pois o nível de evidência exigida é diferente (PRD §5.4; §§8.1–8.3).
- **FR-002**: System MUST definir, para cada tarefa, um **Quality Gate** explícito composto por: (a) evidência mínima requerida, (b) critérios de validade, e (c) resultado observável (aceitar/rejeitar conclusão) (PRD §5.4; R2).
- **FR-003**: System MUST impedir que tarefas de aprendizagem sejam consideradas concluídas sem evidência mínima válida (PRD §5.4; R2).
- **FR-004**: System MUST permitir que tarefas de hábito/fundação sejam concluídas com evidência mínima proporcional (registro mínimo), evitando burocracia (PRD RNF1; §8.3).
- **FR-005**: System MUST, quando um gate falhar, comunicar: (a) **por que** falhou e (b) **o menor próximo passo** para satisfazer o gate (PRD RNF1; RNF3).
- **FR-006**: System MUST aceitar conclusão quando o gate for satisfeito e registrar o **resultado do gate** (satisfeito/não satisfeito) junto da tarefa (PRD §5.4).
- **FR-007**: System MUST lidar com evidência inválida (vazia, ilegível, inaudível) solicitando reenvio ou alternativa equivalente **quando definida** (PRD RNF3).
- **FR-008**: System MUST suportar evidência alternativa apenas quando a alternativa preservar o objetivo do gate. Política global de equivalência (defaults):
  - **Compreensão (input)**: respostas em texto às perguntas de compreensão são equivalência válida.
  - **Retrieval**: respostas em texto são equivalência válida.
  - **Speaking (objetivo: produção oral)**: **áudio é a evidência padrão**; quando áudio for impossível/opt-out ativo, o sistema MUST marcar como bloqueado/substituído e oferecer alternativa de consistência (não equivalente) registrada como tentativa.
- **FR-009**: System MUST registrar, quando aplicável, **rubricas de qualidade** (ex.: speaking) conforme definido pelas SPECS de domínio (PRD §8.1) e associá-las a uma evidência do dia.
- **FR-010**: System MUST registrar e manter “erros recorrentes” e permitir que eles sejam usados para acionar reforços com evidência mínima (PRD §5.4; R3; §§8.1–8.2).
- **FR-011**: System MUST permitir ao usuário consultar o estado do dia: tarefas concluídas/pendentes e quais passaram/falharam gates (PRD §2; §5.3).
- **FR-012**: System MUST operar com privacidade por padrão: antes de coletar evidências potencialmente sensíveis, explicar o mínimo guardado e o propósito (PRD RNF4; ver `SPEC-015`) e aplicar defaults de retenção:
  - **Conteúdo sensível bruto (ex.: áudio)**: retenção curta por padrão (default: **7 dias**) e opção de **não guardar** (processar/validar e descartar).
  - **Derivados não sensíveis**: manter por mais tempo (ex.: resultado do gate, rubrica total, contagens/tendências em `SPEC-016`).
  - O usuário MUST conseguir apagar evidências a qualquer momento (conforme `SPEC-015`), entendendo o impacto em tendências/revisões.

### Non-Functional Requirements

- **NFR-001**: System MUST manter fricção proporcional e evitar burocracia: o gate deve ser o **mínimo** necessário para reduzir falso progresso (PRD RNF1; §5.4).
- **NFR-002**: System MUST ser robusto a dias ruins: sempre que fizer sentido, oferecer uma versão mínima viável do gate (MVD/Plano C) sem humilhar o usuário (PRD RNF2; RNF3).
- **NFR-003**: System MUST manter segurança psicológica: feedback firme, focado em processo e ajuste, sem punição/humilhação (PRD RNF3; §3).
- **NFR-004**: System MUST aplicar privacidade por padrão: coletar o mínimo necessário e ser transparente sobre uso/retensão (PRD RNF4).

### Key Entities *(include if feature involves data)*

- **QualityGate**: definição do gate para uma tarefa (tipo de tarefa; evidência mínima; critérios de validade; caminho mínimo de recuperação).
- **Evidence**: registro de evidência fornecida pelo usuário (tipo; descrição curta; se é sensível; status de validade; timestamp).
- **GateResult**: resultado do gate (satisfeito/não satisfeito); motivo resumido; “próximo passo mínimo” quando falhou.
- **RubricScore**: pontuação por dimensões (quando aplicável) + observações curtas (PRD §8.1; §8.2 quando houver rubrica).
- **RecurringError**: descrição; exemplos; contagem; tendência; último visto; status (ativo/arquivado) (PRD R3).
- **ReinforcementAttempt**: reforço curto solicitado; evidência mínima de tentativa; resultado (feito/não feito).

## Acceptance Criteria *(mandatory)*

- Tarefas de aprendizagem só contam como concluídas quando o Quality Gate correspondente for satisfeito por evidência mínima válida (PRD §5.4; R2).
- Quando o gate falha, o sistema explica o motivo e oferece o caminho mais curto para correção, mantendo interação curta e tom não punitivo (PRD RNF1; RNF3).
- Tarefas de hábitos/fundação têm gates proporcionais (registro mínimo) e comportamentos para dia ruim (MVD/Plano C) (PRD RNF2; §6.1; §8.3).
- Evidência inválida não é aceita; o sistema solicita reenvio ou alternativa equivalente quando definida (PRD RNF3).
- O sistema suporta detecção de “falso progresso” via sinais observáveis (ex.: erros recorrentes) e aciona reforço curto com evidência mínima (PRD R3; §13).
- O usuário consegue consultar o estado do dia e quais tarefas passaram/falharam gates (PRD §2).
- O sistema mantém privacidade por padrão e transparência do que guarda e por quê (PRD RNF4).

## Business Objectives *(mandatory)*

- **Progresso real (qualidade > quantidade)**: reduzir conclusões “sem aprender” por meio de evidência mínima e rubricas quando aplicável (PRD §3; §5.4; R2).
- **Reduzir falso progresso**: identificar padrões de falha e acionar reforço prático (PRD R3; §13).
- **Baixa fricção sustentável**: aplicar o mínimo de evidência necessário, sem transformar o assistente em burocracia (PRD RNF1; R6).
- **Robustez e consistência**: manter ação mínima em dias ruins com gates mínimos apropriados (PRD RNF2; §4).
- **Segurança psicológica**: manter o usuário em movimento com feedback firme e não punitivo (PRD RNF3).
- **Privacidade por padrão**: evidência coletada com transparência e minimização (PRD RNF4).

## Error Handling *(mandatory)*

- **Entrada ausente/ambígua**: se o usuário tenta concluir sem enviar evidência, o sistema deve pedir apenas o mínimo necessário (1 passo) e evitar múltiplas perguntas longas (PRD RNF1).
- **Dia ruim / pouca energia**: oferecer versão mínima do gate quando fizer sentido (especialmente em hábitos/fundação) e reconhecer esforço sem culpa (PRD RNF2; RNF3).
- **Evidência inválida/ilegível/inaudível**: não aceitar conclusão; pedir reenvio. Se o usuário não puder enviar a evidência padrão, aplicar política global de equivalência (FR-008): oferecer equivalência quando preserva objetivo; caso contrário marcar bloqueio/substituição e sugerir mínimo viável (RNF2/RNF3).
- **Evidência parcial**: manter como “não concluído” e indicar exatamente o que falta; registrar como **tentativa** (não concluída) para fins de histórico e micro-consistência, sem contar como consistência da meta intensiva.
- **Tentativa de bypass**: explicar brevemente que a tarefa não conta sem evidência por design (progresso real) e oferecer a rota mínima; manter tom firme e respeitoso (PRD §5.4; RNF3).
- **Usuário some**: permitir retomar o gate do ponto em que parou e lembrar o que faltava, sem penalizar (PRD §5.3; RNF2).
- **Sobrecarga**: se o gate estiver pesado demais para o contexto, orientar a escolha de um plano mínimo (MVD) ou adiar a tarefa, sem “validar” aprendizagem sem evidência (PRD RNF2; §6.2; §13).
- **Privacidade**: quando o usuário mostrar desconforto, explicar minimização e oferecer alternativa/adiamento quando aplicável (PRD RNF4).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Tendência a zero de tarefas de aprendizagem registradas como “concluídas” sem evidência mínima válida (PRD §5.4; R2).
- **SC-002**: Aumento da taxa de tarefas de aprendizagem com evidência válida (ex.: respostas de checagem, rubrica preenchida, etc.) por semana (PRD §8.1–§8.2).
- **SC-003**: Redução na recorrência de erros-alvo ao longo de semanas quando reforços são acionados (PRD R3; §9.2).
- **SC-004**: Em semanas com baixa energia/tempo, aumento da taxa de conclusão de MVD/Plano C para hábitos/fundação (PRD RNF2; §6.1; §8.3).
- **SC-005**: Manter baixa fricção percebida: medir fricção via auto-relato na revisão semanal (`SPEC-007`): pergunta “Quão pesado foi evidenciar esta semana? (0–10)”. Usar também um proxy simples: % de falhas de gate por “evidência ausente” repetida.

