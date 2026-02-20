# Feature Specification: Autoestima — registros curtos + revisão de padrões + micro-exposições (“ações de coragem”)

**Created**: 2026-02-19  
**PRD Base**: §8.5, §5.5, §5.3, §9.1, §9.2, §14, §11 (RNF1–RNF3)

## Caso de uso *(mandatory)*

O usuário quer reduzir autocrítica paralisante e aumentar autoeficácia com intervenções pequenas e consistentes, sem virar “terapia no app” nem burocracia. Esta feature define como o sistema:
- captura **registros curtos** quando há autocrítica forte (gatilho → pensamento → resposta alternativa),
- conduz uma **revisão semanal de padrões** para aprender com a semana (não se culpar),
- e promove **micro-exposições** (“ações de coragem”) que constroem competência na prática.

O tom deve ser firme e acolhedor, com segurança psicológica: foco em processo e tendência, não em culpa.

## Scope & Non-goals *(recommended)*

**In scope (MVP slice)**:
- Registro curto “evento de autocrítica”.
- Um fluxo de revisão semanal para identificar padrões e escolher 1 experimento.
- Um mecanismo de micro-exposições pequenas, com definição observável de “feito”.

**Non-goals (agora)**:
- Não diagnosticar/tratar transtornos; não substituir terapia.
- Não coletar narrativas longas; sem journaling extenso.
- Não fazer intervenções profundas/complexas; foco em mínimo viável.

## Definições *(recommended)*

- **Autocrítica forte**: episódio em que o usuário relata crítica interna intensa que atrapalha execução (escala simples opcional).
- **Registro curto**: 3 campos em linguagem simples: gatilho, pensamento automático, resposta alternativa (mais útil/compassiva).
- **Ação de coragem**: micro-ação concreta, pequena e verificável, alinhada a valores/metas, feita apesar do desconforto.

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Registrar autocrítica forte em ≤ 2 minutos e reduzir fricção (Priority: P1)

O usuário percebe um episódio de autocrítica e quer registrar sem esforço para ganhar clareza e reduzir ruminação.

**Why this priority**: O PRD pede registros curtos e intervenção mínima com segurança psicológica (PRD §8.5; RNF1; RNF3).

**Independent Test**:
- Simular 3 episódios em dias diferentes.
- Validar que o sistema coleta os 3 campos, retorna feedback curto e o registro fica consultável na revisão semanal.

**Acceptance Scenarios**:

1. **Scenario**: Registro curto completo
   - **Given** o usuário relata autocrítica forte
   - **When** o sistema conduz o registro curto
   - **Then** o usuário registra: (a) gatilho, (b) pensamento automático, (c) resposta alternativa curta
   - **And** o sistema confirma em 1 mensagem curta e oferece um próximo passo mínimo (opcional)

2. **Scenario**: Usuário com pouca energia → registro mínimo viável
   - **Given** o usuário está muito cansado/sem tempo
   - **When** o sistema pede o registro
   - **Then** o sistema permite registrar só “gatilho + 1 frase de resposta alternativa” (mínimo viável), sem culpa
   - **And** marca como “parcial” e permite completar depois

3. **Scenario**: Usuário pede para “ver o que registrei hoje”
   - **Given** existem registros no dia atual
   - **When** o usuário solicita “meus registros de autoestima hoje”
   - **Then** o sistema lista os registros de forma curta (sem expor detalhes sensíveis desnecessários)

---

### User Story 2 — Quando emoção está alta, oferecer intervenção mínima (Priority: P1)

O usuário está em alta emoção; ele não quer reflexão longa, só quer se estabilizar para não piorar.

**Why this priority**: Robustez a dias ruins e segurança psicológica (RNF2; RNF3).

**Independent Test**:
- Simular um episódio de alta emoção (usuário diz “tô mal/ansioso”).
- Verificar que o sistema oferece 1 intervenção curta e adia reflexão.

**Acceptance Scenarios**:

1. **Scenario**: Alta emoção → intervenção breve e segura
   - **Given** o usuário relata alta emoção/estresse
   - **When** o sistema responde
   - **Then** o sistema oferece uma intervenção breve (ex.: pausa, respiração curta, aterramento simples) e pergunta apenas 1 coisa curta
   - **And** o sistema oferece registrar o episódio depois, quando estiver melhor

2. **Scenario**: Usuário não quer falar
   - **Given** o usuário não quer explicar detalhes
   - **When** o sistema pergunta
   - **Then** o sistema aceita “prefiro não dizer”, registra a escolha e oferece o mínimo viável sem insistência

---

### User Story 3 — Revisão semanal de padrões + escolher 1 experimento (Priority: P2)

O usuário quer identificar padrões e sair com um próximo experimento pequeno para a semana seguinte.

**Why this priority**: Conecta com PRD (revisão semanal e adaptação contínua) e evita repetir os mesmos gatilhos sem aprender (PRD §5.5; §9.2; §8.5).

**Independent Test**:
- Com pelo menos 3 registros na semana, rodar revisão semanal.
- Validar que o sistema identifica 1–2 padrões e define 1 experimento.

**Acceptance Scenarios**:

1. **Scenario**: Revisão semanal com padrões
   - **Given** existem registros da semana
   - **When** o usuário inicia revisão semanal de autoestima
   - **Then** o sistema sintetiza 1–2 padrões (ex.: gatilhos mais comuns) e mantém tom não punitivo (“isso é dado, não culpa”)

2. **Scenario**: Escolher 1 experimento da semana
   - **Given** padrões foram apresentados
   - **When** o usuário escolhe um foco
   - **Then** o sistema define 1 experimento pequeno (ação de coragem) com critério observável de feito e frequência mínima

---

### User Story 4 — Registrar “ações de coragem” e ver tendência (Priority: P2)

O usuário quer construir autoeficácia fazendo micro-exposições e vendo que está avançando.

**Why this priority**: O PRD explicitamente menciona micro-exposições e métricas de “ações de coragem” (PRD §8.5).

**Independent Test**:
- Simular 1 semana com 3 ações registradas.
- Validar que o sistema contabiliza e mostra tendência simples na revisão semanal.

**Acceptance Scenarios**:

1. **Scenario**: Registrar ação de coragem como concluída
   - **Given** existe uma ação de coragem definida para a semana
   - **When** o usuário executa e registra
   - **Then** o sistema marca como feita e registra 1 evidência mínima textual (1 frase do que foi feito)

2. **Scenario**: Ação falha e vira aprendizado, não culpa
   - **Given** o usuário tentou a ação mas não conseguiu concluir
   - **When** relata a tentativa
   - **Then** o sistema registra “tentativa” e propõe uma versão menor, sem punição

## Edge Cases *(mandatory)*

- What happens when o usuário tem zero registros na semana?
  - Sistema faz revisão “mínima”: pergunta 1–2 questões curtas e define 1 ação de coragem mínima para recomeço.
- What happens when o usuário começa a registrar coisas muito sensíveis e pede confidencialidade?
  - Sistema aplica privacidade por padrão (ver `SPEC-015`), incentiva minimização e oferece opção de registrar de forma mais abstrata.
- What happens when o usuário usa o sistema para autoagressão verbal intensa?
  - Sistema mantém tom protetivo, corta loop de culpa, oferece intervenção mínima e sugere buscar apoio profissional quando apropriado (sem alarmismo).
- What happens when o usuário rejeita “ações de coragem” por medo?
  - Sistema reduz para micro-ação ainda menor e valida o medo como dado, não como falha.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST permitir registrar um episódio de autocrítica com um formato curto: gatilho → pensamento automático → resposta alternativa.
- **FR-002**: System MUST suportar registro mínimo viável (parcial) quando tempo/energia forem baixos.
- **FR-003**: System MUST oferecer uma intervenção breve quando emoção estiver alta e postergar reflexão longa.
- **FR-004**: System MUST conduzir uma revisão semanal de autoestima que sintetize padrões a partir dos registros e leve o usuário a escolher 1 experimento para a próxima semana.
- **FR-005**: System MUST suportar definição e registro de “ações de coragem” (micro-exposições) com critério observável de feito.
- **FR-006**: System MUST permitir consulta curta de registros recentes (ex.: “hoje”/“semana”) em formato que preserve privacidade.

### Non-Functional Requirements

- **NFR-001**: System MUST manter baixa fricção: registros concluíveis em ≤ 2 min (PRD RNF1).
- **NFR-002**: System MUST ser robusto a dias ruins: oferecer mínimo viável e intervenção breve (PRD RNF2).
- **NFR-003**: System MUST manter segurança psicológica: tom não punitivo, foco em processo e tendência (PRD RNF3).
- **NFR-004**: System MUST aplicar privacidade por padrão ao lidar com conteúdo sensível (PRD RNF4; ver `SPEC-015`).

### Key Entities *(include if feature involves data)*

- **SelfEsteemRecord**: data; gatilho; pensamento; resposta_alternativa; intensidade_opcional; status (completo/parcial).
- **CourageAction**: descrição; critério_de_feito; status (planejada/feita/tentativa/adiada); evidência_mínima (1 frase); data.
- **SelfEsteemWeeklyReview**: semana; padrões_sintetizados; experimento_escolhido; contagem_de_acoes; tendência simples.

## Acceptance Criteria *(mandatory)*

- O usuário consegue registrar autocrítica forte em formato curto e o sistema aceita versão mínima em dia ruim.
- Em alta emoção, o sistema oferece intervenção breve e não força reflexão longa.
- A revisão semanal sintetiza padrões e termina com 1 experimento/ação de coragem definida com critério observável.
- O usuário consegue registrar ações de coragem e ver contagem/tendência semanal sem burocracia.
- Tom é não punitivo e preserva segurança psicológica.

## Business Objectives *(mandatory)*

- Reduzir autocrítica paralisante e aumentar autoeficácia por micro-ações e reflexão leve (PRD §8.5).
- Sustentar consistência com baixa fricção e robustez a dias ruins (RNF1/RNF2).
- Manter ambiente de aprendizado e ajuste, não de culpa (RNF3).

## Error Handling *(mandatory)*

- **Usuário em alta emoção**: priorizar intervenção breve e mínimo viável; adiar reflexão.
- **Usuário não quer compartilhar detalhes**: aceitar opt-out e oferecer registro abstrato/curto.
- **Registros muito longos**: pedir para resumir em 1–2 frases e registrar o essencial (minimizar fricção).
- **Conteúdo sensível**: incentivar minimização e transparência sobre retenção/visibilidade (ver `SPEC-015`).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Aumento de ações de coragem por semana ao longo de 4–8 semanas (tendência).
- **SC-002**: Redução de frequência/intensidade auto-relatada de autocrítica paralisante ao longo de semanas.
- **SC-003**: Alta taxa de registros curtos concluídos (sem virar journaling longo), mantendo baixa fricção.