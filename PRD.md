## PRD — Super Assistente Pessoal (Telegram-first) para Metas Anuais

### 1) Contexto e problema
Você quer um **sistema pessoal** que funcione como um “super assistente” para **criar, acompanhar e refinar** tarefas/hábitos até você alcançar suas metas anuais, com **interação ideal por Telegram** e **máxima automação**, reduzindo esforço manual.

O desafio central não é “saber o que fazer”, e sim:
- **Executar consistentemente** em dias bons e ruins.
- **Evitar ilusões de progresso** (ex.: “fiz” mas fiz mal; “estudei” mas não aprendi).
- **Adaptar o plano** com base em evidências reais (desempenho, energia, tempo, contexto).
- **Sustentar mudanças** sem burnout e sem burocracia.

Metas anuais suportadas:
- Fluência em inglês (ênfase em **speaking** e **comprehensible input**)
- Evolução em Java (aprendizado contínuo e mensurável)
- Dormir melhor
- Vida saudável
- Melhorar autoestima
- Em paralelo: avançar no seu SaaS (multicálculo e gestão para corretores de seguros)

### 2) Visão (o “mundo perfeito”)
Um agente que:
- **Entende suas metas como sistemas** (rotinas, restrições, crenças, ambiente).
- **Planeja, executa e audita**: cria planos diários, guia a execução e valida qualidade.
- **Aprende com seus dados**: identifica gargalos (tempo/energia/ansiedade/dificuldade), ajusta o plano e escolhe intervenções com maior probabilidade de funcionar.
- **Faz coaching pragmático**: reduz fricção, cria compromisso, remove obstáculos e mantém você em movimento.
- **Não aceita “marquei como feito”** sem critérios mínimos de qualidade quando o objetivo é aprendizagem/competência.
- **Precisa ter algum meio de consultar os steps que fiz e meu desenvolvimento no dia atual por exemplo**

### 3) Princípios do produto
- **Evidência > Intuição**: decisões guiadas por ciência de aprendizagem e mudança de comportamento (com margem para personalização).
- **Qualidade antes de quantidade**: tarefas com rubricas claras e checagem de qualidade (“quality gates”).
- **Adaptação contínua**: o plano muda conforme desempenho, tempo disponível e energia.
- **Carga cognitiva mínima**: interação curta, objetiva, automatizada; o sistema faz o “trabalho de organização”.
- **Pequenas vitórias + progressão**: começa simples, aumenta dificuldade gradualmente.
- **Segurança psicológica**: feedback firme, mas não punitivo; foco em processo, não em culpa.

### 4) Persona e premissas
Usuário: você, com ambição alta e múltiplas metas simultâneas, que valoriza automação e quer **resultados mensuráveis**.

Premissas operacionais:
- Há dias com pouco tempo/energia; o sistema precisa ter **planos A/B/C**.
- Você terá progresso real se houver: **(1) consistência**, **(2) feedback**, **(3) ajuste**.
- “Sono e energia” são **alavancas** que afetam execução e aprendizagem; portanto, o produto precisa tratar isso como fundamento.

### 5) Escopo (o que o produto faz)
#### 5.1 Onboarding e diagnóstico (1–2 semanas)
Objetivo: criar linha de base e calibrar o sistema ao seu contexto.
- Coleta inicial: metas anuais, restrições (horários, trabalho, saúde), preferências, ambiente, gatilhos de falha, motivadores.
- Linha de base:
  - Sono: diário simples (horário de dormir/levantar, despertares, qualidade percebida).
  - Saúde: atividade física atual, dores/limitações, alimentação (auto-percepção).
  - Inglês: amostras de fala e compreensão (curtas).
  - Java: diagnóstico por exercícios curtos (conceitos + código).
  - Autoestima: padrões de autocrítica, situações gatilho, rotinas de autocuidado.
  - SaaS: definição de 1–2 “apostas” trimestrais (resultados esperados).
- Definição do “mínimo viável diário” (MVD): um conjunto pequeno de ações sustentáveis mesmo em dia ruim.

#### 5.2 Planejamento adaptativo (diário e semanal)
- **Planejamento diário**: 3–5 blocos curtos, com prioridade clara e alternativa se faltar tempo.
- **Planejamento semanal**: revisão de métricas, identificação de gargalos e ajuste de estratégia.
- **Backlog inteligente**: tarefas propostas pelo sistema com base em lacunas (ex.: “você errou X 3x — vamos reforçar com Y”).

#### 5.3 Execução guiada (Telegram-first)
Interação conversacional com:
- Check-in rápido (energia/tempo/estado mental) para escolher Plano A/B/C.
- Instruções objetivas (o que fazer, por quanto tempo, qual critério de qualidade).
- Lembretes e “nudges” quando necessário (sem spam).

#### 5.4 Auditoria de qualidade (“não foi feito até estar bem feito”)
O sistema define **critérios** e pede evidência proporcional ao objetivo:
- Aprendizagem (Inglês/Java): micro-testes, produção (fala/código), e revisão de erros e notas de conversacao e identificar erros para melhoria.
- Hábitos (sono/saúde): diário mínimo + consistência + sinais de melhora.

#### 5.5 Feedback, reflexão e ajuste
- Feedback imediato pós-tarefa (o que funcionou, o que ajustar).
- Revisões semanais com foco em:
  - Tendências (melhora/piora).
  - Motivos de falhas (tempo/energia/ansiedade/ambiente).
  - Próximos experimentos.

### 6) Priorização estratégica (por que e como)
#### 6.1 Fundamentos primeiro: sono e energia
Racional: sono insuficiente e/ou de baixa qualidade está associado a piora de funções executivas e desempenho cognitivo; isso impacta diretamente consistência, autocontrole e aprendizagem. Por isso, **o produto prioriza sono como “infraestrutura”** e evita depender apenas de “força de vontade”.  
Referências: revisões e evidências recentes sobre sono e cognição/executivo (ver seção 12).

#### 6.2 Paralelismo com limites
O produto suporta metas em paralelo, mas impõe limites:
- **No máximo 2 “metas intensivas” por ciclo** (ex.: Inglês + Java), enquanto sono/saúde ficam em modo “fundação”.
- O SaaS entra como “aposta semanal” (1 bloco maior) para evitar canibalizar o resto.

### 7) Metodologias incorporadas (sem “nome bonito”, com execução)
#### 7.1 Mudança de comportamento
- **Planos “Se–Então” (Implementation Intentions)**: transformar intenção em ação automática (“Se for 21:30, então preparo rotina do sono”; “Se eu abrir o Telegram às 07:30, então faço o check-in”). Evidência de efetividade em meta-análises (ver seção 12).
- **MCII (Mental Contrasting + Implementation Intentions)**: visualizar objetivo + obstáculos reais + plano se–então (especialmente útil para consistência).
- **Fricção mínima**: reduzir passos; preparar ambiente; templates; ações de 2 minutos quando necessário.

#### 7.2 Aprendizagem eficaz (Inglês e Java)
O produto se apoia fortemente em:
- **Prática de recuperação (retrieval/practice testing)**: lembrar sem olhar, com feedback.
- **Prática espaçada (distributed practice)**: rever em intervalos, não em maratona.
- **Interleaving (mistura de tópicos)**: alternar tipos de problemas para melhorar transferência.
Essas técnicas têm forte suporte em revisões clássicas e meta-análises (ver seção 12).

#### 7.3 Sono: tratar como competência, não “higiene”
O produto usa princípios de **CBT-I** (Terapia Cognitivo-Comportamental para Insônia) como referência de alta evidência para insônia crônica, com ênfase em componentes mais eficazes (ex.: controle de estímulos, restrição do sono, reestruturação cognitiva) e com cuidado para não virar “lista de dicas” (ver seção 12).

### 8) Requisitos por objetivo (o que significa “progresso real”)
#### 8.1 Inglês (speaking + comprehensible input)
##### Definição de sucesso
- Você consegue **se expressar com clareza** em temas cotidianos e profissionais, com progresso mensurável em fluência, vocabulário ativo e compreensão.

##### Loop diário (exemplo)
- **Input compreensível (10–30 min)**: conteúdo levemente acima do nível, com contexto (áudio/vídeo/texto), visando compreensão global.
- **Output guiado (5–15 min)**:
  - monólogo de 1–3 min gravado, ou
  - shadowing curto, ou
  - conversa guiada por prompts.
- **Retrieval (3–7 min)**: recall de 5–10 itens (expressões/padrões) sem olhar.

##### Quality gates (critérios mínimos)
- Input: responder 3 perguntas de compreensão (sem “chute fácil”).
- Output: gravar e fazer uma autoavaliação rápida por rubrica:
  - clareza (0–2), fluidez (0–2), correção aceitável (0–2), vocabulário/variedade (0–2).
- Erros recorrentes: virar “alvo da semana” (repetição deliberada).

##### Métricas
- Minutos de input (semanal), consistência (dias/semana).
- Amostras semanais de speaking (tempo + rubrica).
- Lista de “erros recorrentes” e taxa de redução.

#### 8.2 Java (aprendizado mensurável)
##### Definição de sucesso
- Evolução contínua em competências (fundamentos → padrões → projeto), com evidência por código e avaliações curtas.

##### Loop diário (exemplo)
- **Prática deliberada (20–45 min)**: kata/exercício com restrição clara (ex.: usar streams; escrever testes; refatorar).
- **Retrieval (5–10 min)**: explicar um conceito “de cabeça” (ou responder mini-quiz).
- **Revisão de erros (5–10 min)**: catalogar 1 erro e a correção (conceitual ou de implementação).

##### Quality gates
- Toda tarefa de código precisa de:
  - definição de “feito” (ex.: passa em testes; cobre casos; refatoração mínima),
  - uma breve explicação do raciocínio (2–5 linhas) e
  - registro do principal aprendizado/erro.

##### Métricas
- Tempo em prática deliberada (semanal).
- Taxa de acerto em quizzes/recall por tópico.
- “Erros recorrentes” por categoria (ex.: collections, OOP, testes, streams).
- Entregas pequenas (microprojetos) por mês.

#### 8.3 Dormir melhor
##### Definição de sucesso
- Melhor consistência de horário, melhor qualidade percebida, redução de latência/despertares (quando aplicável), com melhoria funcional (energia, foco).

##### Loop diário (exemplo)
- Check-in: horário alvo, cafeína, sonecas, estresse.
- Rotina pré-sono pequena e realista.
- Diário de sono (30–60s) ao acordar.

##### Quality gates
- “Cumpriu” sono = diário preenchido + aderência ao plano do dia (mesmo que parcial) + revisão se algo falhou.

##### Métricas
- Regularidade (diferença entre horários).
- Qualidade percebida (0–10).
- Energia pela manhã (0–10).

#### 8.4 Vida saudável
##### Definição de sucesso
- Aumento sustentável de atividade física e hábitos de saúde, alinhado às diretrizes gerais e às suas limitações.

##### Requisitos de base (diretrizes)
- Referência: recomendações da OMS para adultos (150 min/semana moderada ou equivalente; fortalecimento 2x/semana).  
O produto trata isso como **direção**, não como “tudo ou nada”.

##### Loop semanal (exemplo)
- 2 sessões de força (curtas).
- 2–4 sessões de cardio/atividade moderada.
- Ajuste por energia, dor e agenda.

##### Métricas
- Minutos totais/semana e consistência.
- Dor/fadiga (para evitar excesso).

#### 8.5 Autoestima
##### Definição de sucesso
- Menos autocrítica paralisante, mais autoeficácia percebida e estabilidade emocional em situações desafiadoras.

##### Intervenções-base
- **Auto-compaixão** (práticas estruturadas curtas) e exercícios inspirados em CBT para reestruturação de pensamentos automáticos.
- Exposição gradual a ações que constroem competência (não só “afirmações”).

##### Quality gates
- Se houve autocrítica forte no dia: registrar gatilho + pensamento + resposta alternativa (curta).
- Revisão semanal: padrões e um experimento para a semana seguinte.

##### Métricas
- Frequência/intensidade de autocrítica (auto-nota 0–10).
- Ações de coragem (micro-exposições) por semana.

#### 8.6 SaaS (avanço em paralelo)
##### Definição de sucesso
- Progresso consistente em entrega/validação (mesmo que pequeno), sem sabotar fundamentos.

##### Formato recomendado
- 1 “bloco profundo” semanal + 1–2 microblocos de manutenção.
- O sistema exige clareza do resultado do bloco (definição de pronto).

### 9) Fluxos essenciais (jornada do usuário)
#### 9.1 Rotina diária (2–7 minutos de interação total)
- Check-in: “tempo disponível hoje (5/15/30/60+)”, “energia (0–10)”, “humor/estresse (0–10)”
- O sistema escolhe Plano A/B/C e envia:
  - 1 prioridade absoluta
  - 1–2 complementares
  - 1 tarefa de fundação (sono/saúde)
- Pós-tarefa: evidência mínima + rubrica + ajuste para amanhã.

#### 9.2 Revisão semanal (10–20 minutos)
- Painel simples:
  - Consistência por meta
  - Qualidade (rubricas) e gargalos
  - “Erros recorrentes” (Inglês/Java)
  - Sono/energia (tendência)
- Decisões:
  - manter / ajustar / pausar
  - escolher “alvo da semana” (1 em inglês, 1 em Java, 1 de sono/saúde)

### 10) Requisitos funcionais (produto)
- **R1 — Planejamento adaptativo**: gerar tarefas diárias com Plano A/B/C a partir de tempo/energia/contexto.
- **R2 — Quality gates**: tarefas de aprendizagem só contam como concluídas com evidência mínima.
- **R3 — Detecção de falhas reais**: identificar quando você “fez” mas não aprendeu (ex.: errando o mesmo padrão) e reagir com reforço.
- **R4 — Backlog e priorização**: manter lista de próximos passos baseada em lacunas observadas, não só em “vontade”.
- **R5 — Revisão semanal**: síntese de métricas + recomendações acionáveis para a semana.
- **R6 — Automação e baixa fricção**: capturar dados com o mínimo de digitação; usar templates; reduzir escolha.
- **R7 — Personalização progressiva**: começar simples e aumentar sofisticação conforme o sistema aprende seu padrão.

### 11) Requisitos não-funcionais (produto)
- **RNF1 — Simplicidade**: interação curta e previsível; nada de “projeto” para usar o assistente.
- **RNF2 — Robustez a dias ruins**: sempre existir um plano mínimo que mantém a identidade/hábito vivo.
- **RNF3 — Segurança psicológica**: feedback honesto sem humilhação; foco em aprendizado e ajustes.
- **RNF4 — Privacidade por padrão**: coletar o mínimo necessário; clareza do que é guardado e por quê.

### 12) Referências (base de evidências para decisões do PRD)
#### Sono (CBT-I, componentes eficazes, formatos)
- American Academy of Sleep Medicine — guideline (tratamentos comportamentais/psicológicos para insônia crônica): `https://pmc.ncbi.nlm.nih.gov/articles/PMC7853203/`
- Component Network Meta-Analysis (CBT-I; componentes como restrição do sono e controle de estímulos): `https://jamanetwork.com/journals/jamapsychiatry/fullarticle/2814164`
- VA/DoD Clinical Practice Guideline 2025 (insônia crônica/OSA; evidência revisada até 2024): `https://www.healthquality.va.gov/guidelines/CD/insomnia/I-OSA-CPG_2025-Guildeline_final_20250422.pdf`
- Meta-análise 2025 (digital CBT-I automatizado): `https://www.nature.com/articles/s41746-025-01514-4`
- Overview de revisões/meta-análises (distúrbios do sono e cognição em adultos): `https://pmc.ncbi.nlm.nih.gov/articles/PMC10900040/`
- Revisão/meta-análise (restrição do sono e consequências neurocognitivas): `https://www.sciencedirect.com/science/article/abs/pii/S0149763417301641`

#### Aprendizagem (retrieval practice, distributed practice)
- Dunlosky et al. (2013) — revisão de técnicas de aprendizagem (practice testing e distributed practice como altamente efetivas): `https://pubmed.ncbi.nlm.nih.gov/26173288`
- Revisão sistemática aplicada (retrieval practice em escolas/salas de aula): `https://link.springer.com/article/10.1007/s10648-021-09595-9`
- Meta-análise de técnicas (inclui prática de teste e prática distribuída): `https://www.frontiersin.org/articles/10.3389/feduc.2021.581216/full`

#### Aquisição de linguagem (input, listening/reading, interação e feedback)
- Extensive Reading — meta-análise (efeitos em múltiplos domínios, incluindo proficiência e oralidade): `https://link.springer.com/article/10.1007/s10648-025-10068-6`
- Meta-análise (conhecimento de vocabulário e compreensão de leitura/escuta): `https://eric.ed.gov/?id=EJ1342379`
- Listening strategy instruction — meta-análise (efeito em compreensão auditiva): `https://experts.nau.edu/en/publications/the-effectiveness-of-second-language-listening-strategy-instructi`
- Revisão (interação em SLA instruída — papel de negociação/feedback/output): `https://www.cambridge.org/core/journals/language-teaching/article/interaction-and-instructed-second-language-acquisition/78A156EE200F744F5978F99BFB073DBE`

#### Mudança de comportamento (planos se–então, MCII)
- Meta-análise (MCII e alcance de metas): `https://www.ncbi.nlm.nih.gov/pmc/articles/PMC8149892/`
- Revisão/meta-análise (implementation intentions e atividade física): `https://journals.plos.org/plosone/article?id=10.1371/journal.pone.0206294`

#### Saúde (diretrizes de atividade física)
- OMS 2020 (Physical Activity and Sedentary Behaviour — guideline PDF): `https://iris.who.int/bitstream/handle/10665/336656/9789240015128-eng.pdf`

#### Autoestima / saúde mental (auto-compaixão)
- Meta-análise 2023 (intervenções de auto-compaixão; depressão/ansiedade/estresse): `https://link.springer.com/content/pdf/10.1007/s12671-023-02148-x.pdf`

### 13) Riscos e mitigações (produto)
- **Risco: excesso de metas → overload**  
  Mitigação: limite de metas intensivas por ciclo + Plano B/C + MVD.
- **Risco: “falso progresso” (fiz mas não aprendi)**  
  Mitigação: quality gates, micro-testes e reforço de erros recorrentes.
- **Risco: rigidez → frustração**  
  Mitigação: adaptação por energia/tempo; linguagem não punitiva; foco em tendência.
- **Risco: perfeccionismo/autocrítica**  
  Mitigação: rubricas “suficientemente bom”, reforço de processo, práticas de auto-compaixão.

### 14) MVP (primeira versão do produto)
Entregas mínimas para gerar valor rápido:
- Onboarding + diagnóstico leve (sono, inglês, Java, energia/tempo)
- Planejamento diário com Plano A/B/C
- Inglês: input + speaking curto com rubrica
- Java: prática deliberada + mini-quiz/recall
- Diário do sono + 1 intervenção simples por semana
- Revisão semanal com 3 decisões (manter/ajustar/pausar)

