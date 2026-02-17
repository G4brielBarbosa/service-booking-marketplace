# Feature Specification: Nudges/Lembretes "sem spam" + Robustez a dias ruins (Planos B/C + MVD)

**Created**: 2026-02-17  
**PRD Base**: §5.1, §5.3, §6.2, 11 (RNF1, RNF2, RNF3)

## Caso de uso *(mandatory)*

O sistema precisa manter o usuário engajado e em movimento mesmo em dias ruins (baixa energia, pouco tempo, estresse), sem gerar spam ou frustração. Quando o usuário não responde ao check-in ou não executa tarefas, o sistema deve oferecer alternativas progressivamente mais simples (Plano A → B → C → MVD) e enviar lembretes estratégicos que respeitam limites de frequência e contexto.

**Problema**: Dias ruins são inevitáveis e podem quebrar consistência. O sistema precisa:
- Detectar quando o usuário está em dificuldade (não responde, não executa, energia baixa)
- Oferecer planos alternativos automaticamente sem exigir decisão complexa
- Enviar lembretes úteis sem ser intrusivo
- Manter a identidade do hábito vivo mesmo com execução mínima

**Fluxos principais**:
1. Check-in não respondido → sistema oferece Plano B/C automaticamente após um período configurável de ausência de resposta
2. Tarefa não executada → sistema sugere alternativa mais simples ou MVD
3. Múltiplas falhas consecutivas → sistema reduz expectativa e foca em manutenção de hábito
4. Lembretes contextuais → sistema envia no momento certo, respeitando limites de frequência

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Check-in não respondido: oferta automática de Plano B/C (Priority: P1)

O usuário não responde ao check-in matinal dentro de um período razoável (configurável). O sistema detecta a ausência e, em vez de esperar indefinidamente, oferece automaticamente um Plano B ou C baseado no histórico recente e padrões observados, reduzindo a barreira de entrada.

**Why this priority**: Evita que um dia ruim vire uma falha completa. Mantém o usuário em movimento com mínimo esforço cognitivo. É fundamental para robustez (RNF2).

**Independent Test**: Simular ausência de resposta ao check-in e verificar se o sistema oferece alternativas automaticamente após timeout configurável.

**Acceptance Scenarios**:

1. **Scenario**: Check-in não respondido, sistema oferece Plano B automaticamente
   - **Given** usuário recebeu check-in em um horário combinado e não respondeu dentro do período configurado
   - **When** sistema detecta expiração do período configurável de resposta ao check-in **[NEEDS CLARIFICATION: qual período padrão?]**
   - **Then** sistema envia mensagem: "Não recebi seu check-in. Vou assumir um dia corrido e sugerir um Plano B leve. Responda 'ok' para aceitar ou envie seu check-in para personalizar."
   - **And** sistema apresenta Plano B pré-configurado (1 tarefa prioridade + 1 fundação mínima)
   - **And** se usuário continuar sem responder após um segundo período configurável **[NEEDS CLARIFICATION: qual comportamento/intervalo?]**, sistema oferece Plano C automaticamente

2. **Scenario**: Check-in não respondido, histórico indica energia baixa recorrente
   - **Given** usuário não respondeu check-in e histórico mostra energia baixa recorrente nos últimos dias **[NEEDS CLARIFICATION: qual limiar e janela?]**
   - **When** sistema detecta timeout e analisa padrão
   - **Then** sistema oferece diretamente Plano C (MVD) com mensagem empática: "Parece que você está com energia baixa. Vamos manter o essencial hoje?"
   - **And** Plano C contém apenas 1 tarefa mínima de cada meta ativa (ex.: 5 min input inglês, 1 exercício Java simples, diário sono)

3. **Scenario**: Usuário responde após timeout mas antes de execução automática
   - **Given** sistema já ofereceu Plano B automaticamente após timeout
   - **When** usuário envia check-in completo (tempo/energia/humor)
   - **Then** sistema descarta Plano B automático e gera Plano personalizado baseado no check-in real
   - **And** sistema registra que houve atraso mas não penaliza consistência

---

### User Story 2 - Tarefa não executada: sugestão progressiva de alternativas (Priority: P1)

O usuário recebeu uma tarefa do Plano A mas não a executou dentro de um prazo esperado (configurável). O sistema detecta a não execução e oferece alternativas progressivamente mais simples, sempre mantendo a identidade do hábito.

**Why this priority**: Evita que uma tarefa difícil bloqueie todo o progresso. Permite degradação graciosa mantendo consistência. Essencial para RNF2 (robustez a dias ruins).

**Independent Test**: Simular tarefa não executada e verificar se sistema oferece alternativas em cascata (tarefa original → simplificada → mínima → MVD).

**Acceptance Scenarios**:

1. **Scenario**: Tarefa de inglês não executada, sistema oferece versão simplificada
   - **Given** usuário recebeu uma tarefa no Plano A e não executou dentro do prazo configurado **[NEEDS CLARIFICATION: como definir prazo por tipo de tarefa?]**
   - **When** sistema detecta não execução no prazo configurado
   - **Then** sistema envia: "Não vi execução da tarefa de inglês. Quer tentar uma versão mais rápida? Opção 1: Input 10 min + Speaking 3 min. Opção 2: Apenas input 15 min. Responda 1 ou 2."
   - **And** se usuário escolher opção, sistema atualiza tarefa e reinicia timer
   - **And** se usuário não responder dentro de um período configurável **[NEEDS CLARIFICATION]**, sistema oferece MVD (apenas a ação mínima daquela meta)

2. **Scenario**: Tarefa de Java não executada, sistema oferece recall mínimo
   - **Given** usuário recebeu tarefa "Prática deliberada 30 min + Quiz" e não executou
   - **When** sistema detecta não execução após prazo
   - **Then** sistema oferece: "Java não executado. Vamos manter o hábito com um recall rápido de 5 min? Responda 'sim' ou 'pular hoje'."
   - **And** se usuário escolher "sim", sistema envia 1 pergunta de recall simples
   - **And** se usuário escolher "pular hoje", sistema registra como exceção sem culpa e não conta como falha de consistência

3. **Scenario**: Múltiplas tarefas não executadas no mesmo dia
   - **Given** usuário não executou 2+ tarefas do Plano A no mesmo dia
   - **When** sistema detecta padrão de não execução
   - **Then** sistema envia mensagem consolidada: "Vejo que hoje está difícil executar as tarefas. Vamos focar no mínimo essencial? Proponho: [lista MVD de cada meta ativa]"
   - **And** sistema oferece executar tudo de uma vez em bloco curto (ex.: 15 min total)
   - **And** se usuário aceitar MVD, sistema marca todas as tarefas originais como "substituídas por MVD" sem penalizar

---

### User Story 3 - Lembretes contextuais sem spam (Priority: P1)

O sistema envia lembretes estratégicos quando necessário, respeitando limites de frequência (configuráveis) e contexto (não enviar durante horários de sono, não enviar se usuário já está executando).

**Why this priority**: Lembretes são necessários para manter engajamento, mas spam gera frustração e desengajamento. Balanceamento crítico para RNF1 (simplicidade) e RNF3 (segurança psicológica).

**Independent Test**: Simular diferentes cenários de lembretes e verificar se limites de frequência e contexto são respeitados.

**Acceptance Scenarios**:

1. **Scenario**: Lembrete inicial após período configurável sem início de execução
   - **Given** usuário recebeu tarefa e não iniciou execução dentro do prazo configurado para iniciar **[NEEDS CLARIFICATION]**
   - **When** sistema verifica status da tarefa após esse prazo
   - **Then** sistema envia primeiro lembrete: "Lembrete: você tem a tarefa [nome] pendente. Quer começar agora ou prefere ajustar?"
   - **And** sistema registra data/hora do lembrete
   - **And** sistema não envia outro lembrete da mesma tarefa antes do intervalo mínimo configurado **[NEEDS CLARIFICATION]**

2. **Scenario**: Lembrete não enviado durante horário de sono
   - **Given** tarefa não executada e horário atual está dentro da janela de sono configurada pelo usuário
   - **When** sistema verifica se deve enviar lembrete
   - **Then** sistema não envia lembrete
   - **And** sistema agenda lembrete para após o fim da janela de sono
   - **And** sistema registra que lembrete foi adiado por horário de sono

3. **Scenario**: Lembrete não enviado se usuário já está executando
   - **Given** usuário iniciou execução de tarefa (enviou evidência parcial ou marcou "em progresso")
   - **When** sistema verifica se deve enviar lembrete para essa tarefa
   - **Then** sistema não envia lembrete
   - **And** sistema cancela qualquer lembrete agendado para essa tarefa

4. **Scenario**: Limite máximo de lembretes atingido
   - **Given** sistema já enviou a quantidade máxima de lembretes configurada para a mesma tarefa **[NEEDS CLARIFICATION: qual limite padrão?]**
   - **When** sistema verifica se deve enviar mais um lembrete
   - **Then** sistema não envia mais lembretes para essa tarefa enquanto o limite permanecer atingido
   - **And** sistema assume que tarefa não será executada e oferece alternativa (Plano B/C ou MVD) na próxima interação

5. **Scenario**: Lembrete consolidado para múltiplas tarefas pendentes
   - **Given** usuário tem 3 tarefas pendentes e sistema pode enviar lembretes para todas
   - **When** sistema verifica lembretes pendentes
   - **Then** sistema consolida em 1 mensagem: "Você tem 3 tarefas pendentes hoje: [lista]. Quer ajustar o plano ou prefere focar no mínimo essencial?"
   - **And** sistema conta como 1 lembrete por tarefa (não 3 lembretes separados)
   - **And** sistema respeita intervalo mínimo configurado antes de enviar outro lembrete consolidado **[NEEDS CLARIFICATION]**

---

### User Story 4 - MVD (Mínimo Viável Diário) como rede de segurança (Priority: P1)

Quando o usuário está em dificuldade extrema (múltiplas falhas, energia muito baixa, check-in não respondido por muito tempo), o sistema oferece o **MVD** como última alternativa para manter a identidade do hábito vivo sem exigir esforço significativo. O MVD é um conjunto pequeno de ações sustentáveis mesmo em dia ruim (PRD §5.1) e pode ser oferecido independentemente do Plano A/B/C.

**Why this priority**: MVD é a rede de segurança que evita que dias ruins virem semanas ruins. Mantém consistência mesmo em condições adversas. Fundamental para RNF2 e RNF3.

**Independent Test**: Simular condições de dificuldade extrema e verificar se sistema oferece MVD apropriado.

**Acceptance Scenarios**:

1. **Scenario**: MVD oferecido após 3+ falhas consecutivas
   - **Given** usuário não executou tarefas por vários dias consecutivos **[NEEDS CLARIFICATION: qual janela caracteriza “sequência de falhas”?]**
   - **When** sistema detecta padrão de falhas consecutivas
   - **Then** sistema envia mensagem empática: "Vejo que você está passando por uma fase difícil. Vamos manter o essencial? Proponho o MVD de hoje: [lista de ações mínimas, ~10 min total]"
   - **And** sistema não pressiona por execução completa, apenas oferece
   - **And** se usuário executar MVD, sistema celebra como vitória e não penaliza dias anteriores

2. **Scenario**: MVD automático quando energia reportada ≤3
   - **Given** usuário respondeu check-in com energia ≤3
   - **When** sistema processa check-in
   - **Then** sistema oferece diretamente MVD sem apresentar Plano A/B
   - **And** mensagem: "Energia baixa detectada. Vamos com o mínimo essencial hoje?"
   - **And** MVD contém apenas ações de 2-5 minutos por meta ativa

3. **Scenario**: MVD como alternativa quando usuário rejeita Plano B/C
   - **Given** sistema ofereceu Plano B e usuário respondeu "muito difícil" ou "não consigo hoje"
   - **When** sistema recebe rejeição
   - **Then** sistema oferece MVD imediatamente: "Sem problemas. Que tal apenas [ação mínima única]? Só 5 minutos."
   - **And** se usuário aceitar MVD, sistema marca como sucesso e não registra como falha

---

### Edge Cases *(mandatory)*

- **Usuário some completamente (sem resposta por ausência prolongada)**: Sistema envia 1 mensagem final oferecendo MVD e perguntando se está tudo bem. Se não houver resposta após novo período prolongado **[NEEDS CLARIFICATION: qual política/intervalos?]**, sistema pausa lembretes automáticos e aguarda retorno do usuário. Não penaliza consistência durante ausência.

- **Usuário responde mas não executa (ciclo de promessas não cumpridas)**: Após repetidos ciclos de "vou fazer" sem execução **[NEEDS CLARIFICATION: quando o sistema considera que virou padrão?]**, sistema reduz expectativa automaticamente e oferece apenas MVD por um curto período, sem pressionar por planos maiores.

- **Múltiplos lembretes de tarefas diferentes no mesmo horário**: Sistema consolida todos os lembretes em 1 mensagem única para evitar spam, respeitando limite de frequência por tarefa individual.

- **Usuário executa parcialmente (ex.: fez input mas não speaking)**: Sistema celebra a parte executada, oferece completar a parte faltante de forma simplificada, mas não exige. Se usuário não completar, sistema marca como "parcial" e não conta como falha completa.

- **Horário de sono configurado incorretamente ou muda**: Sistema permite ajuste de janela de sono a qualquer momento. Lembretes respeitam janela atual, não histórica.

- **Energia/tempo reportados inconsistentes (ex.: energia 2 mas tempo 60min)**: Sistema prioriza energia sobre tempo para escolher plano. Se inconsistência for extrema, sistema pede confirmação: "Você reportou energia 2 mas tempo 60min. Está se sentindo bem mas sem tempo, ou vice-versa?"

- **Tarefa executada mas evidência inválida/ausente**: Sistema pede reenvio de evidência uma vez. Se não receber dentro de um prazo configurável, oferece alternativa simplificada que não requer evidência complexa (ex.: recall verbal em vez de gravação). **[NEEDS CLARIFICATION: qual prazo?]**

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST detectar ausência de resposta ao check-in após um timeout configurável **[NEEDS CLARIFICATION: timeout padrão ideal?]** e oferecer automaticamente Plano B ou C baseado em histórico e padrões observados.

- **FR-002**: System MUST detectar não execução de tarefas após prazo esperado (configurável) **[NEEDS CLARIFICATION: prazo padrão ideal por tipo de tarefa?]** e oferecer alternativas progressivamente mais simples (tarefa original → simplificada → mínima → MVD).

- **FR-003**: System MUST enviar lembretes contextuais respeitando limites configuráveis: máximo de lembretes por tarefa **[NEEDS CLARIFICATION: limite ideal?]**, intervalo mínimo entre lembretes da mesma tarefa **[NEEDS CLARIFICATION: intervalo ideal?]**, não enviar durante horário de sono configurado, não enviar se usuário já está executando.

- **FR-004**: System MUST consolidar múltiplos lembretes pendentes em 1 mensagem quando possível para evitar spam.

- **FR-005**: System MUST oferecer MVD automaticamente quando detectar condições de dificuldade extrema (energia ≤3, 3+ falhas consecutivas, rejeição de Planos B/C).

- **FR-006**: System MUST permitir que usuário ajuste timeout de check-in, limites de lembretes e prazos esperados de execução (com limites razoáveis) **[NEEDS CLARIFICATION: quais limites?]**.

- **FR-007**: System MUST permitir que usuário configure janela de sono (horário início/fim) para respeitar em lembretes.

- **FR-008**: System MUST registrar exceções (dias ruins, MVD executado, tarefas não executadas) sem penalizar consistência quando apropriado (ex.: MVD conta como sucesso, não como falha).

- **FR-009**: System MUST pausar lembretes automáticos após ausência prolongada e aguardar retorno explícito do usuário **[NEEDS CLARIFICATION: qual política/intervalos?]**.

- **FR-010**: System MUST celebrar execuções parciais e oferecer completar de forma simplificada, mas não exigir conclusão.

### Non-Functional Requirements

- **NFR-001**: System MUST manter simplicidade na interação: mensagens curtas, objetivas, sem jargão técnico (PRD RNF1).

- **NFR-002**: System MUST ser robusto a dias ruins: sempre existir um plano mínimo (MVD) que mantém identidade/hábito vivo mesmo em condições adversas (PRD RNF2).

- **NFR-003**: System MUST manter segurança psicológica: linguagem empática, não punitiva, foco em processo e pequenas vitórias, não em culpa ou falhas (PRD RNF3).

- **NFR-004**: System MUST evitar spam: consolidar lembretes, respeitar limites de frequência, não enviar durante sono ou execução ativa.

- **NFR-005**: System MUST ser adaptativo: ajustar oferta de planos baseado em padrões observados (energia baixa recorrente → oferecer MVD diretamente).

### Key Entities *(include if feature involves data)*

- **Check-in**: Representa a troca de mensagens do check-in diário. Atributos: data/hora de envio, data/hora de resposta (se houver), respostas (tempo disponível, energia 0–10, humor/estresse 0–10), janela de resposta configurada, status (pendente/respondido/expirado).

- **Tarefa**: Representa uma ação a ser executada pelo usuário. Atributos: identificador, tipo (inglês/java/sono/saúde/etc.), plano de origem (A/B/C ou MVD), esforço/complexidade percebida, janela de execução esperada (configurável), status (pendente/em progresso/executada/substituída/cancelada), contagem de lembretes enviados, data/hora do último lembrete.

- **Lembrete**: Representa uma mensagem de lembrete enviada ao usuário sobre tarefa pendente. Atributos: identificador, referência à tarefa, tipo (inicial/consolidado/final), data/hora de envio, conteúdo, contexto aplicado (respeitou janela de sono, usuário estava executando, consolidação com outras tarefas).

- **MVD (Mínimo Viável Diário)**: Representa o conjunto mínimo de ações para manter o hábito vivo em dia ruim. Atributos: identificador, metas ativas cobertas, ações mínimas propostas (curtas e sustentáveis), data/hora de oferta, data/hora de execução (se houver), status (oferecido/aceito/executado/rejeitado). **[NEEDS CLARIFICATION: definição exata do conteúdo do MVD por meta?]**

- **Configuração de Lembretes**: Representa preferências do usuário para lembretes. Atributos: timeout de resposta ao check-in, prazo esperado de execução por tipo de tarefa, janela de sono (início/fim), máximo de lembretes por tarefa, intervalo mínimo entre lembretes.

- **Padrão de Dificuldade**: Representa a detecção de condições adversas. Atributos: tipo (energia baixa recorrente/falhas consecutivas/rejeição de planos), contagem, data/hora da detecção, resposta sugerida (ex.: oferecer MVD, oferecer Plano C, pausar lembretes).

## Acceptance Criteria *(mandatory)*

1. **AC-001**: Quando check-in não é respondido dentro do timeout configurável, sistema oferece automaticamente Plano B ou C sem exigir ação do usuário.

2. **AC-002**: Quando tarefa não é executada dentro do prazo configurável, sistema oferece alternativa simplificada. Se não houver resposta dentro de um período configurável, oferece versão mínima ou MVD.

3. **AC-003**: Lembretes são enviados respeitando: máximo por tarefa (configurável), intervalo mínimo (configurável), não durante sono, não durante execução ativa.

4. **AC-004**: Múltiplos lembretes pendentes são consolidados em 1 mensagem quando possível.

5. **AC-005**: MVD é oferecido automaticamente quando energia ≤3, 3+ falhas consecutivas, ou rejeição de Planos B/C.

6. **AC-006**: Execuções parciais são celebradas e sistema oferece completar de forma simplificada, mas não exige.

7. **AC-007**: Exceções (MVD, dias ruins) são registradas sem penalizar consistência quando apropriado.

8. **AC-008**: Após ausência prolongada (conforme política definida), sistema pausa lembretes automáticos e aguarda retorno explícito. **[NEEDS CLARIFICATION: qual política?]**

9. **AC-009**: Usuário pode configurar timeout de check-in, prazo de execução, e janela de sono com limites razoáveis.

10. **AC-010**: Mensagens mantêm tom empático, não punitivo, focado em processo e pequenas vitórias.

## Business Objectives *(mandatory)*

Esta SPEC suporta diretamente os seguintes objetivos do PRD:

- **Consistência**: Mantém usuário em movimento mesmo em dias ruins através de planos alternativos e MVD, evitando que uma falha vire uma semana ruim.

- **Adaptação contínua**: Sistema detecta padrões de dificuldade e ajusta oferta automaticamente (ex.: energia baixa recorrente → MVD direto).

- **Carga cognitiva mínima**: Lembretes consolidados, ofertas automáticas de planos alternativos, e MVD reduzem necessidade de decisão complexa em momentos de baixa energia.

- **Segurança psicológica**: Linguagem empática, celebração de execuções parciais, registro de exceções sem culpa, e foco em processo (não em falhas) mantêm usuário engajado sem gerar frustração ou autocrítica.

- **Robustez a dias ruins**: Sempre existe um plano mínimo (MVD) que mantém identidade do hábito vivo, cumprindo RNF2 do PRD.

- **Simplicidade**: Mensagens curtas, consolidação de lembretes, e ofertas automáticas reduzem fricção e complexidade da interação.

## Error Handling *(mandatory)*

- **Check-in ausente/ambiguidade**: Sistema aplica um timeout configurável e oferece Plano B automaticamente quando esse tempo expira. **[NEEDS CLARIFICATION: timeout padrão ideal?]** Se o check-in chegar após a oferta automática, sistema prioriza o check-in real e recalibra o plano do dia.

- **Tarefa não executada sem evidência clara**: Sistema assume não execução após prazo esperado e oferece alternativa. Se usuário contestar ("eu fiz mas não marquei"), sistema aceita contestação e pede evidência mínima ou marca como executada com nota.

- **Usuário some completamente (ausência prolongada)**: Sistema envia 1 mensagem final oferecendo MVD e perguntando se está tudo bem. Se não houver resposta após novo período prolongado, sistema pausa lembretes automáticos e aguarda retorno explícito. **[NEEDS CLARIFICATION: qual política/intervalos?]** Não penaliza consistência durante ausência.

- **Sobrecarga de lembretes (múltiplas tarefas pendentes)**: Sistema consolida lembretes em 1 mensagem única. Respeita limites configuráveis de lembretes por tarefa e de intervalo mínimo entre lembretes. **[NEEDS CLARIFICATION: limites/intervalos padrão?]**

- **Configuração inválida (timeout muito curto/longo, janela de sono inconsistente)**: Sistema valida limites razoáveis para evitar spam e frustração e pede ajuste quando a configuração for inválida. **[NEEDS CLARIFICATION: quais limites e defaults?]**

- **Energia/tempo reportados inconsistentes**: Sistema prioriza energia sobre tempo para escolher plano. Se inconsistência for extrema (ex.: energia 2 mas tempo 60min), pede confirmação: "Você reportou energia 2 mas tempo 60min. Está se sentindo bem mas sem tempo, ou vice-versa?"

- **Evidência inválida ou ausente após execução parcial**: Sistema celebra parte executada, pede evidência da parte faltante uma vez. Se não receber dentro de um prazo configurável, oferece alternativa simplificada que não requer evidência complexa (ex.: recall verbal em vez de gravação). **[NEEDS CLARIFICATION: qual prazo?]**

- **Ciclo de promessas não cumpridas (usuário diz "vou fazer" mas não executa)**: Após padrão de promessas não cumpridas ser detectado, sistema reduz expectativa e oferece apenas MVD por um período curto, sem pressionar por planos maiores. **[NEEDS CLARIFICATION: qual regra/limiar e por quanto tempo reduzir a expectativa?]**

- **Horário de sono muda ou configuração incorreta**: Sistema permite ajuste a qualquer momento. Lembretes respeitam a janela atual. Se usuário reportar que recebeu lembrete durante sono, sistema sugere revisar a janela de sono e pede confirmação.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Taxa de recuperação em dias ruins: % de dias em que check-in não foi respondido inicialmente e ainda assim ocorre execução de Plano B/C ou MVD (medido semanalmente). **[NEEDS CLARIFICATION: meta/threshold desejado?]**

- **SC-002**: Abandono após sequência de dias ruins: % de usuários que param de usar após 3+ dias ruins consecutivos (medido mensalmente). **[NEEDS CLARIFICATION: meta/threshold desejado?]**

- **SC-003**: Efetividade de lembretes: % de tarefas com lembrete enviado que são executadas dentro de um intervalo definido após o lembrete (medido semanalmente). **[NEEDS CLARIFICATION: qual intervalo e meta?]**

- **SC-004**: Adoção de MVD: % dos MVDs oferecidos que são aceitos e executados (medido semanalmente). **[NEEDS CLARIFICATION: meta/threshold desejado?]**

- **SC-005**: Percepção de spam: % de usuários que reportam lembretes como "intrusivos" ou "spam" (medido mensalmente via feedback opcional). **[NEEDS CLARIFICATION: meta/threshold desejado?]**

- **SC-006**: Diferença de consistência entre dias normais e dias ruins: comparar consistência em dias com energia baixa vs dias com energia normal (medido semanalmente). **[NEEDS CLARIFICATION: definição de corte de energia e meta/threshold?]**

- **SC-007**: Latência entre expiração do check-in e oferta automática de Plano B/C (medido como tempo médio). **[NEEDS CLARIFICATION: meta/threshold desejado?]**

- **SC-008**: Taxa de follow-up em execuções parciais: % de execuções parciais em que o sistema oferece completar de forma simplificada e o usuário tenta completar (aceita ou rejeita explicitamente) (medido semanalmente). **[NEEDS CLARIFICATION: meta/threshold desejado?]**
