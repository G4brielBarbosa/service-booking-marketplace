# Feature Specification: Personalização progressiva & governança de metas em paralelo (limites de overload)

**Created**: 2026-02-17  
**PRD Base**: §5.1, §5.3, §6.2, 10 (R7), 11 (RNF1, RNF3, RNF4), 13 (riscos)

## Caso de uso *(mandatory)*

O sistema precisa suportar múltiplas metas em paralelo sem causar overload, mantendo consistência e qualidade. O PRD estabelece limites explícitos: no máximo 2 "metas intensivas" por ciclo (ex.: Inglês + Java), enquanto sono/saúde operam em modo "fundação" e SaaS como "aposta semanal" (PRD §6.2). Além disso, o sistema deve começar simples e aumentar sofisticação progressivamente conforme aprende padrões do usuário (PRD R7), evitando frustração por rigidez ou excesso de complexidade inicial (PRD RNF1, RNF3).

**Problema**: Usuários com múltiplas metas podem tentar fazer tudo ao mesmo tempo, resultando em:
- Overload cognitivo e burnout
- Perda de consistência em todas as metas
- Falso progresso (fazendo muitas coisas mal ao invés de poucas bem)
- Frustração com sistema muito complexo desde o início

**Solução**: O sistema impõe limites de metas intensivas ativas, permite pausar/retomar com contexto, detecta sinais de overload e ajusta personalização gradualmente baseado em dados observados (padrões de execução, energia, consistência, qualidade).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Limite de metas intensivas e bloqueio de novas ativações (Priority: P1)

**Why this priority**: Previne overload desde o início e garante que o sistema respeite o limite fundamental do PRD (§6.2). Sem isso, o risco de excesso de metas não é mitigado.

**Independent Test**: Usuário tenta ativar uma terceira meta intensiva quando já tem 2 ativas; sistema bloqueia e oferece opções (pausar uma existente ou manter limite).

**Acceptance Scenarios**:

1. **Scenario**: Tentativa de ativar terceira meta intensiva
   - **Given** usuário tem 2 metas intensivas ativas (ex.: Inglês e Java)
   - **When** usuário tenta ativar uma terceira meta intensiva (ex.: novo curso)
   - **Then** sistema bloqueia ativação e apresenta mensagem explicando limite de 2 metas intensivas por ciclo
   - **And** sistema oferece opções: (a) pausar uma meta existente temporariamente, (b) manter limite atual e focar nas 2 ativas, (c) agendar ativação futura após conclusão/pausa de uma existente

2. **Scenario**: Ativação bem-sucedida dentro do limite
   - **Given** usuário tem 1 meta intensiva ativa (Inglês)
   - **When** usuário ativa segunda meta intensiva (Java)
   - **Then** sistema permite ativação e confirma que ambas estão ativas
   - **And** sistema registra que limite de 2 metas intensivas foi atingido

3. **Scenario**: Metas de fundação não contam para limite intensivo
   - **Given** usuário tem 2 metas intensivas ativas (Inglês e Java)
   - **When** usuário ativa ou mantém ativa meta de fundação (sono ou saúde)
   - **Then** sistema permite sem bloqueio, pois fundação não conta para limite de intensivas
   - **And** sistema mantém distinção clara entre intensivas e fundação na comunicação com o usuário

---

### User Story 2 - Pausar e retomar metas com registro de contexto (Priority: P1)

**Why this priority**: Permite flexibilidade sem perder contexto, essencial para adaptação e segurança psicológica (PRD RNF3). Usuário precisa poder ajustar sem culpa.

**Independent Test**: Usuário pausa uma meta intensiva, sistema registra motivo e data; após período, usuário retoma e sistema oferece plano de retomada baseado no tempo pausado.

**Acceptance Scenarios**:

1. **Scenario**: Pausar meta intensiva com motivo
   - **Given** usuário tem 2 metas intensivas ativas
   - **When** usuário solicita pausar uma meta (ex.: Java) por motivo (ex.: "período de trabalho intenso")
   - **Then** sistema registra: meta pausada, data de pausa, motivo fornecido, estado de progresso no momento da pausa
   - **And** sistema libera slot para nova meta intensiva (se desejado)
   - **And** sistema confirma pausa com mensagem não punitiva (ex.: "Java pausada. Você pode retomar quando estiver pronto.")

2. **Scenario**: Retomar meta pausada
   - **Given** usuário tem meta pausada há X dias/semanas
   - **When** usuário solicita retomar meta pausada
   - **Then** sistema oferece plano de retomada adaptado ao tempo pausado:
     - Se pausa curta (≤ 7 dias): retomada gradual por 1–2 sessões com carga reduzida (versão mínima) e revisão rápida do “último erro/tema”.
     - Se pausa média (8–30 dias): retomada com um diagnóstico leve do domínio (1–2 passos curtos) e 1 semana de carga reduzida antes de voltar ao normal.
     - Se pausa longa (> 30 dias): retomar como “quase nova” para aquele domínio: reiniciar baseline leve do domínio (sem re-onboarding total) e reconstruir alvos/erros recorrentes conforme evidência recente.
   - **And** sistema restaura contexto anterior (erros recorrentes, progresso, preferências) quando relevante

3. **Scenario**: Consultar histórico de pausas
   - **Given** usuário tem histórico de pausas/retomadas
   - **When** usuário consulta histórico de uma meta
   - **Then** sistema exibe: datas de pausa/retomada, motivos registrados, duração de cada pausa
   - **And** sistema permite identificar padrões (ex.: "pausa Java sempre em períodos de trabalho intenso")

---

### User Story 3 - Detecção de sinais de overload e sugestão de ajuste (Priority: P2)

**Why this priority**: Detecta proativamente quando usuário está sobrecarregado mesmo dentro dos limites, permitindo intervenção antes da falha completa. Mitiga risco de rigidez → frustração (PRD §13).

**Independent Test**: Sistema observa padrão de consistência baixa + energia baixa + qualidade reduzida por uma janela suficiente para indicar tendência (default: últimos 7 dias com dados em ≥3 dias); sugere pausar uma meta ou reduzir intensidade.

**Acceptance Scenarios**:

1. **Scenario**: Detecção de overload por múltiplos sinais
   - **Given** sistema observa por uma janela consistente (default: últimos 7 dias): consistência baixa em 2+ metas intensivas (≤2 dias/7), energia média baixa (média ≤3/10), e/ou qualidade (rubricas) em queda (média semanal caiu ≥1 ponto vs semana anterior, quando houver rubrica)
   - **When** sistema detecta padrão de overload (tendência)
   - **Then** sistema apresenta alerta não punitivo: "Notei que você está com menos energia e consistência esta semana. Isso é normal. Quer ajustar algo?"
   - **And** sistema oferece sugestões: (a) pausar temporariamente uma meta intensiva, (b) reduzir carga de uma meta específica, (c) manter e observar mais uma semana
   - **And** sistema registra detecção e resposta do usuário para aprendizado

2. **Scenario**: Overload detectado mas usuário escolhe manter
   - **Given** sistema detectou sinais de overload e sugeriu ajustes
   - **When** usuário escolhe manter metas como estão
   - **Then** sistema aceita escolha sem insistir
   - **And** sistema continua monitorando e pode sugerir novamente no máximo 1x por semana (cooldown default: 7 dias) se o padrão persistir
   - **And** sistema oferece Plano C/MVD mais frequente para dias ruins

---

### User Story 4 - Personalização progressiva baseada em padrões observados (Priority: P2)

**Why this priority**: Implementa PRD R7 (personalização progressiva). Sistema começa simples e aumenta sofisticação conforme aprende, evitando complexidade inicial que causa frustração.

**Independent Test**: Sistema inicia com configurações simples; após dados suficientes para identificar padrões, oferece primeira personalização (ex.: ajuste de horários preferidos, tipos de tarefas que funcionam melhor). Default de “dados suficientes”: ≥14 dias desde onboarding mínimo e ao menos 7 check-ins + 5 sessões concluídas (com gate quando aplicável) no período.

**Acceptance Scenarios**:

1. **Scenario**: Início com configuração simples
   - **Given** usuário está no onboarding ou primeiras semanas
   - **When** sistema apresenta opções de personalização
   - **Then** sistema oferece apenas configurações essenciais e simples (ex.: horários preferidos, nível de lembretes)
   - **And** sistema não apresenta opções avançadas ou complexas inicialmente

2. **Scenario**: Oferta progressiva de personalização após aprendizado
   - **Given** sistema observou padrões por janela mínima (default: 14 dias) com dados suficientes
   - **When** sistema identifica oportunidade de personalização útil baseada nesses padrões
   - **Then** sistema oferece ajuste específico baseado em padrão observado (ex.: "Notei que você tem mais energia às 7h. Quer que eu priorize tarefas intensivas nesse horário?")
   - **And** sistema explica o padrão observado que fundamenta a sugestão
   - **And** usuário pode aceitar, recusar ou ajustar a sugestão

3. **Scenario**: Aumento gradual de sofisticação sem sobrecarregar
   - **Given** usuário já aceitou algumas personalizações anteriores e manteve consistência
   - **When** sistema identifica nova oportunidade de personalização mais sofisticada
   - **Then** sistema oferece uma personalização por vez, não múltiplas simultaneamente
   - **And** sistema valida que usuário está confortável antes de oferecer próxima
   - **And** sistema permite reverter personalizações anteriores se não funcionarem

---

### Edge Cases *(mandatory)*

- **O que acontece se usuário tenta pausar todas as metas intensivas simultaneamente?**  
  Sistema permite pausar todas, mas mantém pelo menos uma meta de fundação ativa (sono ou saúde). Se usuário pausar tudo, sistema oferece MVD mínimo para manter hábito vivo.

- **Como sistema lida com usuário que ativa/pausa metas frequentemente (instabilidade)?**  
  Sistema detecta padrão de instabilidade (default: ≥3 pausas/retomadas em ≤30 dias) e oferece conversa de diagnóstico: "Notei que você está ajustando metas com frequência. Quer revisar sua estratégia geral?" Sugere sessão de planejamento ou redução temporária de ambição.

- **O que acontece se sistema detecta overload mas usuário não responde à sugestão?**  
  Sistema não insiste imediatamente. Default: se não houver resposta em 48h, faz 1 follow-up curto (respeitando políticas anti-spam de `SPEC-011`). Se ainda sem resposta por 7 dias, para de insistir no tema e passa a favorecer MVD/Plano C até haver sinais de recuperação ou decisão explícita.

- **Como sistema diferencia "meta intensiva" de "meta de fundação" e "aposta semanal"?**  
  Usar a classificação default definida em `SPEC-001` (FR-002): intensivas (Inglês/Java), fundação (sono/saúde/autoestima) e aposta semanal (SaaS).

- **O que acontece se usuário quer aumentar limite de metas intensivas além de 2?**  
  Limite é **rígido** (PRD §6.2). Se usuário solicitar aumento, o sistema explica o racional e oferece alternativas dentro do limite (pausar, agendar próximo ciclo, ou corrigir classificação).

- **Como sistema lida com personalização que piora desempenho?**  
  Sistema detecta queda de consistência/qualidade após personalização e oferece reverter. Se usuário confirma reversão, sistema aprende que aquela personalização não funcionou e evita sugerir similar no futuro.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST impor limite máximo de 2 metas intensivas ativas simultaneamente por ciclo (base: PRD §6.2).
- **FR-002**: System MUST bloquear ativação de terceira meta intensiva quando limite já está atingido, oferecendo opções claras (pausar existente, manter limite, agendar futura).
- **FR-003**: System MUST permitir pausar qualquer meta intensiva com registro de: data de pausa, motivo fornecido pelo usuário, estado de progresso no momento da pausa.
- **FR-004**: System MUST permitir retomar meta pausada com plano de retomada adaptado ao tempo de pausa (curta/média/longa).
- **FR-005**: System MUST manter histórico de pausas/retomadas por meta, incluindo datas, motivos e durações.
- **FR-006**: System MUST detectar sinais de overload quando observa por uma janela consistente (default: últimos 7 dias com dados em ≥3 dias): consistência baixa em múltiplas metas intensivas (≤2 dias/7), energia média baixa (média ≤3/10) e/ou qualidade abaixo do padrão anterior (queda de rubrica ≥1 ponto vs semana anterior quando aplicável).
- **FR-007**: System MUST oferecer sugestões de ajuste quando detecta overload: pausar meta, reduzir carga, manter e observar.
- **FR-008**: System MUST começar com configurações simples e aumentar personalização progressivamente conforme aprende padrões do usuário (PRD R7).
- **FR-009**: System MUST oferecer personalizações baseadas em padrões observados (horários, tipos de tarefas, padrões de falha) após dados suficientes para suportar a inferência. Default de suficiência: ≥14 dias desde onboarding mínimo e ao menos 7 check-ins + 5 sessões concluídas no período, e sem sinal forte de overload na última semana.
- **FR-010**: System MUST oferecer uma personalização por vez, validando conforto do usuário antes de oferecer próxima.
- **FR-011**: System MUST permitir reverter personalizações que não funcionaram e aprender com reversões.
- **FR-012**: System MUST distinguir metas intensivas de metas de fundação e apostas semanais para aplicar limites corretos.
- **FR-013**: System MUST permitir consulta do estado atual de limites (quantas metas intensivas ativas, slots disponíveis).

### Non-Functional Requirements

- **NFR-001**: System MUST reduzir overload e frustração (PRD §13), aplicando limites e detectando sinais proativamente.
- **NFR-002**: System MUST manter segurança psicológica ao sugerir ajustes: linguagem não punitiva, foco em processo, sem culpa (PRD RNF3).
- **NFR-003**: System MUST evitar complexidade inicial: começar simples e aumentar sofisticação gradualmente (PRD R7, RNF1).
- **NFR-004**: System MUST aplicar privacidade por padrão: registrar apenas dados necessários para governança e personalização (PRD RNF4).

### Key Entities *(include if feature involves data)*

- **Meta**: Representa uma meta do usuário (ex.: Inglês, Java, sono). Atributos: classificação (intensiva/fundação/aposta semanal) conforme `SPEC-001`, status (ativa/pausada/concluída), data/hora de ativação, data/hora de pausa (se pausada), motivo de pausa (se pausada), histórico de pausas/retomadas.
- **Ciclo de Metas**: Representa um período em que um conjunto de metas está ativo sob regras de limite. Atributos: data de início, metas intensivas ativas (máximo 2), metas de fundação ativas, apostas semanais ativas.
- **Sinal de Overload**: Representa a detecção de um padrão de sobrecarga. Atributos: data/hora de detecção, indicadores observados (consistência, energia, qualidade), sugestões oferecidas, resposta do usuário, data/hora de resolução (se houver).
- **Personalização**: Representa um ajuste personalizado oferecido ao usuário. Atributos: tipo (ex.: horário preferido, tipo de tarefa), padrão observado que fundamenta, data/hora de oferta, status (aceita/recusada/ajustada/revertida), data/hora de reversão (se revertida), resultado observado (melhorou/manteve/piorou desempenho). Default de medição de “resultado”: comparar as 2 semanas antes vs 2 semanas depois (quando houver dados) em (a) consistência (dias/semana), (b) taxa de gates satisfeitos, (c) média de rubrica quando aplicável, e (d) fricção percebida (0–10 na revisão semanal). “Melhorou” se houver melhora em pelo menos 1 métrica sem piora relevante nas outras; “piorou” se cair consistência ≥2 dias/semana ou rubrica ≥1 ponto sem compensação.

## Acceptance Criteria *(mandatory)*

1. **AC-001**: Sistema bloqueia ativação de terceira meta intensiva quando 2 já estão ativas, apresentando mensagem clara e opções acionáveis.
2. **AC-002**: Sistema permite pausar meta intensiva registrando data, motivo e estado de progresso, e libera slot para nova meta se desejado.
3. **AC-003**: Sistema oferece plano de retomada adaptado ao tempo de pausa (curta/média/longa) quando usuário retoma meta pausada.
4. **AC-004**: Sistema detecta padrão de overload após uma janela consistente (default: últimos 7 dias com dados em ≥3 dias), usando limiares default (consistência ≤2/7 em 2+ metas intensivas e/ou energia média ≤3/10 e/ou queda de rubrica ≥1 ponto) e oferece sugestões não punitivas.
4. **AC-004**: Sistema detecta padrão de overload após uma janela consistente (default: últimos 7 dias com dados em ≥3 dias), usando limiares default (consistência ≤2/7 em 2+ metas intensivas e/ou energia média ≤3/10 e/ou queda de rubrica ≥1 ponto) e oferece sugestões não punitivas.
5. **AC-005**: Sistema inicia com configurações simples e oferece primeira personalização baseada em padrões observados após dados suficientes (default: janela mínima 14 dias; ≥7 check-ins + 5 sessões concluídas).
6. **AC-006**: Sistema oferece personalizações uma por vez, validando conforto antes de oferecer próxima.
7. **AC-007**: Sistema permite reverter personalizações e aprende com reversões para evitar sugestões similares futuras.
8. **AC-008**: Sistema mantém histórico consultável de pausas/retomadas por meta, incluindo datas, motivos e durações.
9. **AC-009**: Sistema distingue corretamente metas intensivas de fundação/aposta para aplicar limites apenas em intensivas.
10. **AC-010**: Todas as interações de governança mantêm tom não punitivo e foco em processo, não em culpa (PRD RNF3).

## Business Objectives *(mandatory)*

Esta SPEC suporta os seguintes objetivos do PRD:

- **Prevenção de overload**: Mitiga risco de excesso de metas causando burnout e perda de consistência (PRD §13). Limites explícitos e detecção proativa protegem o usuário.

- **Adaptação contínua**: Permite ajustes sem perder contexto (pausar/retomar) e personalização baseada em evidências reais, não em suposições (PRD §3, R7).

- **Carga cognitiva mínima**: Começa simples e aumenta complexidade gradualmente, evitando frustração inicial com sistema muito sofisticado (PRD §3, RNF1).

- **Segurança psicológica**: Linguagem não punitiva, sem culpa, foco em processo. Usuário pode ajustar metas sem sentir fracasso (PRD RNF3).

- **Sustentabilidade**: Sistema aprende padrões e ajusta, permitindo que usuário mantenha múltiplas metas sem burnout a longo prazo (PRD §3).

- **Evidência > Intuição**: Personalizações baseadas em dados observados, não em suposições (PRD §3).

## Error Handling *(mandatory)*

- **Tentativa de ativar terceira meta intensiva**: Sistema bloqueia e apresenta mensagem explicando limite, oferecendo opções claras (pausar existente, manter limite, agendar futura). Não permite override sem confirmação explícita do usuário sobre entendimento do risco.

- **Pausar todas as metas intensivas**: Sistema permite mas mantém pelo menos uma meta de fundação ativa. Se usuário pausar tudo, oferece MVD mínimo para manter hábito vivo.

- **Instabilidade frequente (múltiplas pausas/retomadas)**: Sistema detecta padrão (default: ≥3 pausas/retomadas em ≤30 dias) e oferece conversa de diagnóstico, sugerindo revisão de estratégia ou redução temporária de ambição.
  Default: instabilidade = ≥3 pausas/retomadas em ≤30 dias; resposta = propor revisar estratégia e reduzir ambição por 1 semana (MVD + 1 meta intensiva).

- **Usuário não responde a sugestão de ajuste por overload**: Sistema não insiste imediatamente. Default: follow-up em 48h (1 mensagem curta, respeitando `SPEC-011`); se sem resposta por 7 dias, parar de insistir no tema e favorecer MVD/Plano C até resposta/recuperação.
  Default: follow-up em 48h (1 mensagem curta, respeitando `SPEC-011`); se sem resposta por 7 dias, parar de insistir no tema e favorecer MVD/Plano C até resposta/recuperação.

- **Personalização piora desempenho**: Sistema detecta queda de consistência/qualidade após personalização e oferece reverter. Se usuário confirma, sistema aprende e evita sugerir similar no futuro.

- **Dados insuficientes para personalização**: Sistema não oferece personalização até ter dados suficientes para identificar padrões com confiança (default: 14 dias; ≥7 check-ins; ≥5 sessões concluídas; sem overload forte recente). Se dados são escassos, mantém configuração simples e informa que personalização virá quando houver padrões claros.
  Default: dados suficientes conforme FR-009 (14 dias; ≥7 check-ins; ≥5 sessões concluídas; sem overload forte recente).

- **Conflito de prioridade ao pausar**: Se usuário precisa pausar mas não sabe qual meta, sistema solicita escolha objetiva (1 prioridade absoluta) e propõe ajustes baseados em dados observados (ex.: "Java teve menor consistência esta semana. Sugestão: pausar Java temporariamente.").

- **Retomar meta após pausa muito longa**: default: pausa longa > 30 dias. Tratar como retomada “quase nova” no domínio: baseline leve do domínio + 1 semana de carga reduzida antes de voltar ao normal.

- **Erro ao registrar pausa/retomada**: Sistema tenta novamente automaticamente. Se falhar, informa usuário e mantém estado anterior até conseguir registrar. Não bloqueia uso do sistema.
- **Erro ao registrar pausa/retomada**: Se o registro falhar, o sistema informa o usuário em mensagem curta e mantém o estado anterior até confirmação. Default de recuperação sem fricção: na próxima interação, o sistema pergunta apenas “Confirmar: Java está pausada? (sim/não)” e reconcilia o estado; não envia múltiplas mensagens proativas sobre isso.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Sistema impede ativação de terceira meta intensiva quando limite já está atingido (meta: 100% de bloqueio).

- **SC-002**: % das pausas de metas registradas com motivo fornecido pelo usuário (meta inicial: ≥80%).

- **SC-003**: Sistema detecta padrão de overload quando o usuário apresenta sinais consistentes na janela mínima (default: 7 dias; dados em ≥3 dias) e oferece ajuste no máximo 1x/semana (meta inicial: detectar em ≥80% dos casos em que critérios são atendidos e o usuário concorda que houve semana ruim).

- **SC-004**: Primeira personalização baseada em padrões é oferecida após dados suficientes (default: 14 dias; ≥7 check-ins; ≥5 sessões concluídas; sem overload forte recente). Meta: oferecer no máximo 1 personalização/semana.

- **SC-005**: Taxa de aceitação de personalizações oferecidas (meta inicial: ≥50%).

- **SC-006**: Taxa de reversão de personalizações (meta inicial: ≤30%).

- **SC-007**: Consistência média com 2 metas intensivas ativas (meta inicial: não cair abaixo de 3/7 dias em cada meta na média de 4 semanas, quando o usuário reporta energia média ≥4).

- **SC-008**: Padrão de instabilidade (default: ≥3 pausas/retomadas em ≤30 dias) recebe sugestão de diagnóstico em 100% dos casos após detecção.

- **SC-009**: Distinção correta entre metas intensivas/fundação/aposta (medir por casos de teste de classificação e auditoria de mudanças de estado; meta inicial: 100% de aplicação correta do limite apenas em intensivas).

- **SC-010**: Todas as mensagens de governança mantêm tom não punitivo (avaliação qualitativa: revisão de amostras de mensagens).
