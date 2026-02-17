# Feature Specification: Métricas & Registros — consistência, rubricas, tendências, erros recorrentes

**Created**: 2026-02-17  
**PRD Base**: §5.4, §9.2, §§8.1–8.5 (métricas), 10 (R2, R5, R6), 11 (RNF1, RNF4)

## Caso de uso *(mandatory)*

O sistema precisa manter métricas e registros mínimos para suportar:
- **Planejamento adaptativo**: identificar gargalos (tempo/energia/ansiedade/dificuldade) e ajustar planos diários/semanais
- **Quality gates**: validar que tarefas de aprendizagem foram concluídas com qualidade mínima (rubricas)
- **Revisão semanal**: fornecer síntese de tendências, consistência e erros recorrentes para decisões acionáveis (manter/ajustar/pausar)
- **Consulta do progresso**: permitir que o usuário consulte rapidamente os steps executados e seu desenvolvimento no dia atual

O desafio é capturar dados suficientes para decisões baseadas em evidência, mantendo **baixa fricção** (mínima digitação) e **privacidade por padrão** (coletar apenas o necessário).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Consulta rápida do progresso do dia atual (Priority: P1)

O usuário quer verificar rapidamente quais steps executou hoje e como está seu desenvolvimento em cada meta, sem precisar de múltiplas mensagens/consultas.

**Why this priority**: Suporta a necessidade explícita do PRD (§2) de "consultar os steps que fiz e meu desenvolvimento no dia atual". É pré-requisito para confiança no sistema e para revisões semanais.

**Independent Test**: Usuário executa 2-3 tarefas em diferentes metas durante o dia, depois solicita resumo do progresso do dia. Sistema retorna lista consolidada com status de cada meta.

**Acceptance Scenarios**:

1. **Scenario**: Consulta progresso do dia após executar tarefas
   - **Given** usuário executou 1 tarefa de inglês (input + rubrica), 1 tarefa de Java (prática + recall) e preencheu diário do sono hoje
   - **When** usuário solicita "progresso hoje" ou "como estou hoje"
   - **Then** sistema retorna resumo consolidado mostrando: inglês (1/1 concluído, rubrica média X), Java (1/1 concluído), sono (diário preenchido), consistência parcial do dia

2. **Scenario**: Consulta progresso do dia sem tarefas executadas
   - **Given** é início do dia e nenhuma tarefa foi executada ainda
   - **When** usuário solicita progresso do dia
   - **Then** sistema retorna "Dia iniciado - nenhuma tarefa executada ainda" e oferece opção de ver plano do dia

3. **Scenario**: Consulta progresso com dados parciais
   - **Given** usuário executou tarefa de inglês mas não completou rubrica ainda
   - **When** usuário solicita progresso do dia
   - **Then** sistema retorna inglês como "em progresso" (tarefa iniciada, rubrica pendente) e não marca como concluído até quality gate ser atendido

---

### User Story 2 - Registro automático de consistência por meta (Priority: P1)

O sistema registra automaticamente a consistência (dias/semana) de cada meta quando tarefas são concluídas, permitindo consulta rápida e uso em revisões semanais.

**Why this priority**: Base para revisão semanal (PRD §9.2) e detecção de gargalos. Suporta decisões de manter/ajustar/pausar.

**Independent Test**: Usuário executa tarefas de inglês em 4 dias da semana. Sistema registra consistência 4/7 e disponibiliza em revisão semanal.

**Acceptance Scenarios**:

1. **Scenario**: Registro automático ao concluir tarefa com quality gate
   - **Given** usuário conclui tarefa de inglês com rubrica preenchida (quality gate atendido)
   - **When** sistema processa conclusão da tarefa
   - **Then** sistema registra automaticamente "inglês executado" para o dia atual na métrica de consistência semanal

2. **Scenario**: Consistência não conta tarefa sem quality gate
   - **Given** usuário marca tarefa como "feita" mas não preencheu rubrica (quality gate não atendido)
   - **When** sistema verifica consistência
   - **Then** sistema não conta este dia na consistência semanal e mantém tarefa como "pendente de evidência"

3. **Scenario**: Consulta de consistência semanal
   - **Given** semana atual tem 5 dias decorridos e usuário executou inglês em 3 dias
   - **When** usuário solicita consistência de inglês ou sistema gera revisão semanal
   - **Then** sistema retorna "3/5 dias" (ou "3/7" se considerar semana completa) e tendência comparada à semana anterior

---

### User Story 3 - Registro e consolidação de rubricas de qualidade (Priority: P1)

O sistema registra rubricas de qualidade (quando aplicável) e calcula tendências semanais para identificar melhora/piora em aprendizagem.

**Why this priority**: Suporta quality gates (PRD §5.4, R2) e detecção de "falso progresso" (PRD R3). Base para ajustes de estratégia.

**Independent Test**: Usuário executa 5 tarefas de inglês na semana com rubricas. Sistema calcula média semanal e compara com semana anterior, mostrando tendência.

**Acceptance Scenarios**:

1. **Scenario**: Registro de rubrica ao concluir tarefa de aprendizagem
   - **Given** usuário conclui tarefa de speaking em inglês e preenche rubrica (clareza: 2, fluidez: 1, correção: 2, vocabulário: 1)
   - **When** sistema processa conclusão com rubrica
   - **Then** sistema registra rubrica completa (total: 6/8) associada à tarefa e à data, e disponibiliza para cálculo de tendências

2. **Scenario**: Cálculo de tendência semanal de rubricas
   - **Given** semana atual tem 4 rubricas de inglês (médias: 5, 6, 6, 7) e semana anterior tinha média 5.2
   - **When** sistema gera revisão semanal ou usuário solicita tendência
   - **Then** sistema retorna média semanal atual (6.0) e indica "melhora" comparada à semana anterior (+0.8)

3. **Scenario**: Rubrica incompleta não bloqueia registro parcial
   - **Given** usuário preenche apenas 2 dos 4 critérios da rubrica (clareza e fluidez)
   - **When** sistema processa rubrica parcial
   - **Then** sistema registra valores preenchidos, marca como "parcial" e solicita complemento opcional (não bloqueia conclusão se critérios mínimos forem atendidos conforme definição da meta)

---

### User Story 4 - Registro e consolidação de erros recorrentes (Priority: P1)

O sistema identifica, registra e consolida erros recorrentes por domínio (inglês/Java) e acompanha tendência de redução ao longo do tempo.

**Why this priority**: Suporta reforço automático (PRD R3) e "alvo da semana" (PRD §9.2). Base para backlog inteligente.

**Independent Test**: Usuário comete mesmo erro de gramática em inglês 3 vezes na semana. Sistema identifica como recorrente, registra e sugere como "alvo da semana" na revisão.

**Acceptance Scenarios**:

1. **Scenario**: Identificação automática de erro recorrente
   - **Given** usuário comete erro "forgot + gerund" em speaking de inglês pela 3ª vez na semana
   - **When** sistema processa registro de erro da tarefa
   - **Then** sistema identifica como "erro recorrente" (≥3 ocorrências), consolida na lista de erros recorrentes de inglês e marca frequência

2. **Scenario**: Tendência de redução de erro recorrente
   - **Given** erro "forgot + gerund" apareceu 5 vezes na semana anterior e 2 vezes nesta semana
   - **When** sistema gera revisão semanal ou usuário consulta erros recorrentes
   - **Then** sistema mostra erro com frequência atual (2x) e indica "redução" comparada à semana anterior (-3 ocorrências)

3. **Scenario**: Erro recorrente vira "alvo da semana"
   - **Given** sistema identifica 2 erros recorrentes em inglês (frequência alta) e 1 em Java
   - **When** sistema gera revisão semanal
   - **Then** sistema sugere 1 erro de inglês e 1 de Java como "alvos da semana" para reforço, baseado em frequência e impacto

---

### User Story 5 - Registro de sono/energia com tendências (Priority: P1)

O sistema registra métricas de sono (regularidade, qualidade percebida, energia pela manhã) e calcula tendências para identificar padrões que impactam execução.

**Why this priority**: Sono é "infraestrutura" (PRD §6.1) e impacta diretamente consistência e aprendizagem. Base para ajustes de planejamento.

**Independent Test**: Usuário preenche diário do sono por 7 dias. Sistema calcula regularidade (diferença entre horários), média de qualidade e energia, e mostra tendência.

**Acceptance Scenarios**:

1. **Scenario**: Registro diário de sono ao acordar
   - **Given** usuário preenche diário do sono ao acordar (horário dormiu: 23:30, horário acordou: 07:00, qualidade: 7/10, energia manhã: 6/10)
   - **When** sistema processa diário do sono
   - **Then** sistema registra todos os valores, calcula duração (7h30min) e regularidade comparada ao dia anterior, e disponibiliza para tendências

2. **Scenario**: Cálculo de tendência semanal de sono
   - **Given** semana atual tem regularidade média de 1h30min de variação e qualidade média 6.5/10, enquanto semana anterior tinha 2h de variação e qualidade 5.8/10
   - **When** sistema gera revisão semanal ou usuário solicita tendência de sono
   - **Then** sistema retorna melhor regularidade (-30min) e melhora de qualidade (+0.7) comparada à semana anterior, e associa ao impacto em energia/execução

3. **Scenario**: Sono incompleto não bloqueia registro parcial
   - **Given** usuário preenche apenas horários de sono mas não qualidade/energia
   - **When** sistema processa diário parcial
   - **Then** sistema registra valores disponíveis (regularidade calculável), marca como "parcial" e não penaliza consistência (sono ainda conta como registrado)

---

### User Story 6 - Revisão semanal com síntese de métricas (Priority: P1)

O sistema consolida todas as métricas da semana (consistência, rubricas, erros recorrentes, sono/energia) em painel simples para revisão e decisões acionáveis.

**Why this priority**: Suporta revisão semanal explícita do PRD (§9.2, R5). Base para decisões de manter/ajustar/pausar e escolha de "alvos da semana".

**Independent Test**: Usuário completa semana de atividades. Sistema gera painel com consistência por meta, tendências de rubricas, erros recorrentes consolidados e tendência de sono/energia, permitindo 3 decisões.

**Acceptance Scenarios**:

1. **Scenario**: Geração automática de revisão semanal
   - **Given** semana atual está completa (7 dias) e usuário executou tarefas em múltiplas metas
   - **When** sistema detecta fim da semana ou usuário solicita revisão semanal
   - **Then** sistema gera painel consolidado mostrando: consistência por meta (ex.: inglês 5/7, Java 4/7), tendências de rubricas (melhora/piora), erros recorrentes consolidados, tendência de sono/energia, e oferece 3 decisões (manter/ajustar/pausar) + sugestão de alvos da semana

2. **Scenario**: Revisão semanal com dados parciais
   - **Given** semana atual tem apenas 3 dias decorridos mas usuário solicita revisão antecipada
   - **When** sistema gera revisão parcial
   - **Then** sistema mostra métricas disponíveis (3 dias), marca como "revisão parcial" e oferece opção de revisão completa ao fim da semana

3. **Scenario**: Decisões acionáveis na revisão semanal
   - **Given** revisão semanal mostra inglês com consistência baixa (3/7) e tendência de rubricas piorando
   - **When** usuário escolhe "ajustar" para inglês na revisão
   - **Then** sistema registra decisão, sugere ajustes específicos (ex.: reduzir carga, focar em input) e atualiza planejamento para próxima semana

---

### Edge Cases *(mandatory)*

- **Dados faltantes em múltiplos dias**: Como sistema lida com semana com 2-3 dias sem registros? Sistema deve calcular consistência apenas sobre dias com dados disponíveis e marcar lacunas sem penalizar, oferecendo recuperação de baseline mínima se necessário.

- **Rubricas inconsistentes entre metas**: Como sistema lida quando inglês usa rubrica 4 critérios (0-2 cada) e Java usa rubrica diferente? Sistema deve normalizar para comparação de tendências ou manter separado por domínio conforme definição de cada meta.

- **Erros recorrentes com nomenclatura diferente**: Como sistema identifica que "forgot + gerund" e "esqueceu de usar gerund" são o mesmo erro? Sistema deve usar categorização/agrupamento semântico ou permitir marcação manual de equivalência.

- **Sono com múltiplos registros no mesmo dia**: Como sistema lida se usuário preenche diário do sono mais de uma vez no mesmo dia (correção)? Sistema deve manter último registro válido e deixar claro que houve correção/atualização. **[NEEDS CLARIFICATION: precisamos manter histórico de correções?]**

- **Consulta de progresso durante execução de tarefa**: Como sistema lida quando usuário solicita progresso do dia enquanto está executando uma tarefa (rubrica pendente)? Sistema deve mostrar estado atual (em progresso) e não contar como concluído até quality gate.

- **Revisão semanal sem dados suficientes**: Como sistema gera revisão na primeira semana (sem baseline anterior)? Sistema deve mostrar apenas dados disponíveis, marcar como "semana inicial" e não comparar com semana anterior inexistente.

- **Métricas de metas pausadas**: Como sistema trata métricas de metas que foram pausadas no meio da semana? Sistema deve excluir da consistência ativa mas manter histórico para retomada futura.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST registrar automaticamente consistência por meta (dias/semana) quando tarefas são concluídas com quality gate atendido.

- **FR-002**: System MUST permitir consulta rápida do progresso do dia atual, mostrando steps executados e desenvolvimento por meta em formato consolidado.

- **FR-003**: System MUST registrar rubricas de qualidade (quando aplicável) associadas a tarefas e datas, e calcular tendências semanais (média, comparação com semana anterior).

- **FR-004**: System MUST identificar, registrar e consolidar erros recorrentes por domínio (inglês/Java) com frequência, e calcular tendência de redução ao longo do tempo.

- **FR-005**: System MUST registrar métricas de sono (regularidade calculada a partir de horários, qualidade percebida 0-10, energia pela manhã 0-10) e calcular tendências semanais.

- **FR-006**: System MUST gerar revisão semanal consolidada com: consistência por meta, tendências de rubricas, erros recorrentes consolidados, tendência de sono/energia, e oferecer 3 decisões (manter/ajustar/pausar) + sugestão de alvos da semana.

- **FR-007**: System MUST manter histórico de métricas suficiente para comparações semanais (mínimo: semana atual + semana anterior).

- **FR-008**: System MUST aplicar quality gates antes de contar tarefas na consistência (tarefas sem evidência mínima não contam como concluídas).

### Non-Functional Requirements

- **NFR-001**: System MUST manter captura de dados com mínima digitação, usando templates e valores padrão quando possível (PRD R6, RNF1).

- **NFR-002**: System MUST aplicar privacidade por padrão: coletar apenas métricas necessárias para funcionalidades descritas, com clareza do que é guardado e por quê (PRD RNF4).

- **NFR-003**: System MUST garantir que consultas de progresso e revisões semanais sejam rápidas o suficiente para manter a interação curta e previsível (PRD RNF1). **[NEEDS CLARIFICATION: existe um target de tempo?]**

- **NFR-004**: System MUST lidar com dados faltantes/parciais sem penalizar usuário, oferecendo recuperação de baseline mínima quando necessário.

### Key Entities *(include if feature involves data)*

- **Métrica de Consistência**: Representa execução de tarefas por meta em um período (ex.: semana). Atributos: meta (referência), período (ex.: semana/ano), dias com execução válida, total de dias do período, tendência (comparação com período anterior).

- **Registro de Rubrica**: Representa avaliação de qualidade de uma tarefa de aprendizagem. Atributos: tarefa (referência), data, critérios (ex.: clareza, fluidez, correção aceitável, vocabulário/variedade — valores 0–2 cada conforme PRD §8.1), total, domínio (inglês/Java), status (completa/parcial).

- **Erro Recorrente**: Representa um padrão de erro identificado como frequente. Atributos: domínio (inglês/Java), descrição/categoria, frequência no período atual, frequência no período anterior, tendência (redução/aumento/estável), status (ativo/alvo da semana/resolvido).

- **Registro de Sono**: Representa entrada diária do diário do sono. Atributos: data, horário_dormiu, horário_acordou, duração_calculada, qualidade_percebida (0-10), energia_manhã (0-10), regularidade_calculada (diferença vs dia anterior), status (completo/parcial).

- **Revisão Semanal**: Representa a consolidação de métricas do período semanal. Atributos: semana/ano, consistência por meta (lista), tendências de rubricas (lista), erros recorrentes consolidados (lista), tendência de sono/energia (conjunto de indicadores), decisões oferecidas (manter/ajustar/pausar), alvos da semana sugeridos (lista), status (completa/parcial).

- **Progresso do Dia**: Representa o estado consolidado de execução no dia atual. Atributos: data, metas com tarefas (lista com status: concluído/em progresso/pendente), rubricas preenchidas (contagem), consistência parcial do dia (calculada sobre tarefas concluídas conforme regras aplicáveis).

## Acceptance Criteria *(mandatory)*

- **AC-001**: Usuário consegue consultar progresso do dia atual e ver lista consolidada de steps executados e desenvolvimento por meta de forma rápida o suficiente para manter a interação curta. **[NEEDS CLARIFICATION: existe target de tempo?]**

- **AC-002**: Sistema registra automaticamente consistência quando tarefa é concluída com quality gate, e usuário consegue ver consistência semanal por meta na revisão.

- **AC-003**: Sistema calcula e exibe tendências de rubricas (média semanal, comparação com semana anterior) na revisão semanal.

- **AC-004**: Sistema identifica erros recorrentes (≥3 ocorrências) automaticamente, consolida por domínio e mostra tendência de redução na revisão semanal.

- **AC-005**: Sistema registra métricas de sono diariamente, calcula regularidade e tendências semanais (qualidade, energia, regularidade) disponíveis na revisão.

- **AC-006**: Sistema gera revisão semanal consolidada com todas as métricas (consistência, rubricas, erros, sono) e oferece 3 decisões acionáveis + alvos da semana.

- **AC-007**: Sistema lida com dados faltantes/parciais sem penalizar consistência, oferecendo recuperação quando necessário.

- **AC-008**: Tarefas sem quality gate atendido não contam na consistência e permanecem como "pendentes de evidência".

## Business Objectives *(mandatory)*

- **Evidência > Intuição**: Métricas e registros permitem decisões baseadas em dados reais (desempenho, energia, tempo) em vez de intuição, suportando adaptação contínua do plano (PRD §3, §5.2).

- **Qualidade antes de quantidade**: Rubricas e quality gates garantem que progresso seja mensurável e real, evitando "ilusões de progresso" (PRD §1, R2, R3).

- **Adaptação contínua**: Tendências e erros recorrentes identificam gargalos (tempo/energia/ansiedade/dificuldade) para ajustes de estratégia (PRD §5.2, §5.5).

- **Carga cognitiva mínima**: Captura automática e consultas rápidas reduzem esforço manual, mantendo foco em execução (PRD R6, RNF1).

- **Segurança psicológica**: Tratamento de dados faltantes sem penalização e foco em tendências (não em culpa) mantém ambiente de aprendizado (PRD RNF3).

- **Privacidade por padrão**: Coleta mínima necessária com clareza do propósito suporta confiança no sistema (PRD RNF4).

## Error Handling *(mandatory)*

- **Dados faltantes em múltiplos dias**: Sistema calcula consistência apenas sobre dias com dados disponíveis, marca lacunas sem penalizar e oferece recuperação de baseline mínima (ex.: "quer preencher dados faltantes?") quando o gap for significativo. **[NEEDS CLARIFICATION: qual critério de “gap significativo”?]**

- **Rubricas incompletas ou inválidas**: Sistema aceita rubricas parciais se critérios mínimos forem atendidos (conforme definição da meta), registra como "parcial" e não bloqueia conclusão. Se rubrica estiver completamente ausente em tarefa que exige quality gate, sistema bloqueia conclusão e oferece caminho mais curto para gerar evidência.

- **Erros recorrentes com nomenclatura ambígua**: Sistema usa categorização semântica quando possível, mas permite marcação manual de equivalência pelo usuário se agrupamento automático falhar. Erros não agrupados são tratados separadamente até confirmação.

- **Sono com múltiplos registros no mesmo dia**: Sistema mantém o último registro válido como fonte de verdade e deixa claro para o usuário que houve uma correção/atualização no mesmo dia. **[NEEDS CLARIFICATION: precisamos manter histórico de correções?]**

- **Consulta de progresso durante execução**: Sistema mostra estado atual (tarefas em progresso com rubricas pendentes) e diferencia claramente "concluído" de "em progresso", não contando como consistência até quality gate.

- **Revisão semanal sem baseline**: Sistema mostra apenas dados disponíveis da semana atual, marca como "semana inicial" e não tenta comparar com semana anterior inexistente. Oferece baseline manual opcional se usuário quiser.

- **Métricas de metas pausadas**: Sistema exclui metas pausadas da consistência ativa e revisão semanal, mas mantém histórico completo para retomada futura. Mostra claramente "meta pausada em [data]" na consulta.

- **Sobrecarga de dados históricos**: Sistema mantém histórico suficiente para as comparações previstas (ex.: semana atual vs semana anterior) e garante que consultas continuem rápidas (PRD RNF1). **[NEEDS CLARIFICATION: qual política de histórico além do mínimo?]**

- **Entrada ambígua em consultas**: Se usuário solicita "progresso" sem especificar dia/semana, sistema assume "dia atual" por padrão e pede confirmação com 1 pergunta curta se contexto for ambíguo.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Tempo de resposta para consulta de progresso do dia (médio e p95). **[NEEDS CLARIFICATION: target de tempo?]**

- **SC-002**: Sistema registra automaticamente consistência em 100% das tarefas concluídas com quality gate (sem necessidade de ação manual do usuário).

- **SC-003**: Tempo de geração da revisão semanal (médio e p95) para consolidação de métricas (consistência, rubricas, erros, sono). **[NEEDS CLARIFICATION: target de tempo?]**

- **SC-004**: Taxa de identificação de erros recorrentes com limiar definido (ex.: “3 ocorrências” como exemplo do PRD) sem necessidade de marcação manual. **[NEEDS CLARIFICATION: qual limiar e meta/threshold?]**

- **SC-005**: Tendências de rubricas e sono são calculadas corretamente (média do período, comparação com período anterior) quando há dados suficientes. **[NEEDS CLARIFICATION: critério de “dados suficientes”?]**

- **SC-006**: Fricção de registro: tempo médio diário gasto em registros manuais (excluindo execução de tarefas). **[NEEDS CLARIFICATION: target de tempo?]**

- **SC-007**: Sistema lida com dados faltantes/parciais sem penalizar consistência e oferece recuperação quando gap é significativo. **[NEEDS CLARIFICATION: critério de gap e meta/threshold?]**

- **SC-008**: Quality gates são aplicados corretamente: 0% de tarefas sem evidência contam na consistência.
