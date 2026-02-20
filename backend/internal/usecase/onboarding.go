package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/abriesouza/super-assistente/internal/domain"
	"github.com/abriesouza/super-assistente/internal/port"
	"github.com/google/uuid"
)

type OnboardingUseCase struct {
	users      port.UserRepository
	sessions   port.OnboardingRepository
	goals      port.GoalCycleRepository
	privacy    port.PrivacyPolicyRepository
	baselines  port.BaselineRepository
	mvds       port.MVDRepository
	events     port.EventRepository
	idempotent port.IdempotencyStore
	log        *slog.Logger
}

func NewOnboardingUseCase(
	users port.UserRepository,
	sessions port.OnboardingRepository,
	goals port.GoalCycleRepository,
	privacy port.PrivacyPolicyRepository,
	baselines port.BaselineRepository,
	mvds port.MVDRepository,
	events port.EventRepository,
	idempotent port.IdempotencyStore,
	log *slog.Logger,
) *OnboardingUseCase {
	return &OnboardingUseCase{
		users:      users,
		sessions:   sessions,
		goals:      goals,
		privacy:    privacy,
		baselines:  baselines,
		mvds:       mvds,
		events:     events,
		idempotent: idempotent,
		log:        log,
	}
}

type OnboardingPrompt struct {
	NextStep    domain.StepID  `json:"next_step"`
	Message     string         `json:"message"`
	Choices     []string       `json:"choices,omitempty"`
	Disclosure  string         `json:"disclosure,omitempty"`
}

type OnboardingProgress struct {
	Status       domain.OnboardingStatus `json:"status"`
	NextStep     domain.StepID           `json:"next_step"`
	Message      string                  `json:"message"`
	PendingItems []domain.PendingItem    `json:"pending_items,omitempty"`
}

func (uc *OnboardingUseCase) StartOnboarding(ctx context.Context, telegramUserID, chatID int64) (*OnboardingPrompt, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		user_ := domain.NewUserProfile(telegramUserID, chatID)
		user = &user_
		if err := uc.users.Save(ctx, user); err != nil {
			return nil, fmt.Errorf("creating user: %w", err)
		}

		pp := domain.NewDefaultPrivacyPolicy(user.UserID)
		if err := uc.privacy.Save(ctx, &pp); err != nil {
			return nil, fmt.Errorf("creating privacy policy: %w", err)
		}
	}

	existing, err := uc.sessions.FindByUserID(ctx, user.UserID)
	if err == nil && existing != nil {
		return uc.buildPromptForStep(existing.CurrentStepID, existing), nil
	}

	session := domain.NewOnboardingSession(user.UserID)
	if err := uc.sessions.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("creating onboarding session: %w", err)
	}

	loc := user.Location()
	evt := domain.NewEvent(user.UserID, domain.EventOnboardingStarted, domain.SensitivityC5, nil, loc)
	_ = uc.events.Append(ctx, evt)

	return uc.buildPromptForStep(domain.StepWelcome, session), nil
}

func (uc *OnboardingUseCase) SubmitAnswer(ctx context.Context, telegramUserID int64, stepID domain.StepID, value map[string]any) (*OnboardingProgress, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "Inicie o onboarding com /start")
	}

	session, err := uc.sessions.FindByUserID(ctx, user.UserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Sessão de onboarding não encontrada", "Inicie com /start")
	}

	// Validate goals limit before submitting
	if stepID == domain.StepSelectGoals {
		if domErr := uc.validateGoalSelection(value); domErr != nil {
			return nil, domErr
		}
	}

	if domErr := session.SubmitAnswer(stepID, value); domErr != nil {
		return nil, domErr
	}

	if err := uc.sessions.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("updating session: %w", err)
	}

	// Side effects per step
	loc := user.Location()
	if err := uc.processStepSideEffects(ctx, user, session, stepID, value, loc); err != nil {
		uc.log.Error("step side effect failed", "step", stepID, "error", err)
	}

	evt := domain.NewEvent(user.UserID, domain.EventOnboardingStepCompleted, domain.SensitivityC5,
		map[string]any{"step_id": string(stepID)}, loc)
	_ = uc.events.Append(ctx, evt)

	if session.Status == domain.OnboardingMinimumCompleted {
		evt := domain.NewEvent(user.UserID, domain.EventOnboardingMinCompleted, domain.SensitivityC5, nil, loc)
		_ = uc.events.Append(ctx, evt)
	}

	return &OnboardingProgress{
		Status:       session.Status,
		NextStep:     session.CurrentStepID,
		Message:      uc.messageForStep(session.CurrentStepID, session),
		PendingItems: session.PendingItems,
	}, nil
}

func (uc *OnboardingUseCase) ResumeOnboarding(ctx context.Context, telegramUserID int64) (*OnboardingPrompt, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	session, err := uc.sessions.FindByUserID(ctx, user.UserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Nenhum onboarding em progresso", "Inicie com /start")
	}

	loc := user.Location()
	evt := domain.NewEvent(user.UserID, domain.EventOnboardingResumed, domain.SensitivityC5, nil, loc)
	_ = uc.events.Append(ctx, evt)

	return uc.buildPromptForStep(session.CurrentStepID, session), nil
}

func (uc *OnboardingUseCase) GetSummary(ctx context.Context, telegramUserID int64) (*domain.OnboardingSummary, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	session, err := uc.sessions.FindByUserID(ctx, user.UserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Onboarding não iniciado", "")
	}

	cycle, _ := uc.goals.FindByUserID(ctx, user.UserID)
	pp, _ := uc.privacy.FindByUserID(ctx, user.UserID)
	baselines, _ := uc.baselines.FindAllByUser(ctx, user.UserID)
	mvd, _ := uc.mvds.FindByUserID(ctx, user.UserID)

	summary := &domain.OnboardingSummary{
		PendingItems: session.PendingItems,
		Baselines:    baselines,
		MVD:          mvd,
		PrivacyPolicy: pp,
	}

	if cycle != nil {
		summary.ActiveGoals = cycle.ActiveGoals
		summary.PausedGoals = cycle.PausedGoals
	}

	restrictions := session.GetAnswer(domain.StepRestrictions)
	if restrictions != nil {
		summary.Restrictions = restrictions.Value
	}

	return summary, nil
}

func (uc *OnboardingUseCase) SetPrivacyPolicy(ctx context.Context, telegramUserID int64, optOutCategories []domain.SensitivityLevel, minimalMode bool) error {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return domain.NewNotFoundError("Usuário não encontrado", "")
	}

	pp, err := uc.privacy.FindByUserID(ctx, user.UserID)
	if err != nil {
		p := domain.NewDefaultPrivacyPolicy(user.UserID)
		pp = &p
	}

	pp.OptOutCategories = optOutCategories
	pp.MinimalMode = minimalMode
	pp.UpdatedAt = time.Now()

	if err := uc.privacy.Update(ctx, pp); err != nil {
		if err := uc.privacy.Save(ctx, pp); err != nil {
			return fmt.Errorf("saving privacy policy: %w", err)
		}
	}

	loc := user.Location()
	evt := domain.NewEvent(user.UserID, domain.EventPrivacyPolicySet, domain.SensitivityC5,
		map[string]any{"minimal_mode": minimalMode, "opt_out_count": len(optOutCategories)}, loc)
	_ = uc.events.Append(ctx, evt)

	return nil
}

func (uc *OnboardingUseCase) validateGoalSelection(value map[string]any) *domain.DomainError {
	goalsRaw, ok := value["goals"].([]any)
	if !ok {
		return domain.NewValidationError("Seleção de metas inválida", "Envie a lista de metas")
	}

	intensiveCount := 0
	for _, g := range goalsRaw {
		gStr, ok := g.(string)
		if !ok {
			continue
		}
		if domain.GoalClassifications[domain.GoalID(gStr)] == domain.GoalClassIntensive {
			intensiveCount++
		}
	}

	if intensiveCount > domain.MaxIntensiveGoals {
		return domain.NewValidationError(
			"Limite de metas intensivas excedido",
			fmt.Sprintf("Você selecionou %d metas intensivas, mas o máximo é %d. Escolha quais manter e quais pausar para proteger sua consistência.",
				intensiveCount, domain.MaxIntensiveGoals),
		)
	}

	return nil
}

func (uc *OnboardingUseCase) processStepSideEffects(ctx context.Context, user *domain.UserProfile, session *domain.OnboardingSession, stepID domain.StepID, value map[string]any, loc *time.Location) error {
	switch stepID {
	case domain.StepSelectGoals:
		return uc.saveGoalCycle(ctx, user, value, loc)
	case domain.StepSleepBaseline:
		return uc.saveBaseline(ctx, user, domain.GoalSleep, value, domain.BaselineMinimum)
	case domain.StepEnglishBase:
		return uc.saveBaseline(ctx, user, domain.GoalEnglish, value, domain.BaselineMinimum)
	case domain.StepJavaBaseline:
		return uc.saveBaseline(ctx, user, domain.GoalJava, value, domain.BaselineMinimum)
	case domain.StepMVD:
		return uc.saveMVD(ctx, user, value)
	}
	return nil
}

func (uc *OnboardingUseCase) saveGoalCycle(ctx context.Context, user *domain.UserProfile, value map[string]any, loc *time.Location) error {
	goalsRaw, _ := value["goals"].([]any)
	var entries []domain.GoalEntry
	for _, g := range goalsRaw {
		gStr, _ := g.(string)
		entries = append(entries, domain.GoalEntry{ID: domain.GoalID(gStr)})
	}

	cycle, domErr := domain.NewGoalCycle(user.UserID, entries)
	if domErr != nil {
		return domErr
	}

	if err := uc.goals.Save(ctx, cycle); err != nil {
		return err
	}

	evt := domain.NewEvent(user.UserID, domain.EventGoalCycleSet, domain.SensitivityC5,
		map[string]any{"active_count": len(entries), "intensive_count": len(cycle.IntensiveGoals())}, loc)
	return uc.events.Append(ctx, evt)
}

func (uc *OnboardingUseCase) saveBaseline(ctx context.Context, user *domain.UserProfile, goalDomain domain.GoalID, value map[string]any, completeness domain.BaselineCompleteness) error {
	bs := &domain.BaselineSnapshot{
		BaselineID:   uuid.New(),
		UserID:       user.UserID,
		Domain:       goalDomain,
		Data:         value,
		Completeness: completeness,
		CapturedAt:   time.Now(),
		UpdatedAt:    time.Now(),
	}
	return uc.baselines.Save(ctx, bs)
}

func (uc *OnboardingUseCase) saveMVD(ctx context.Context, user *domain.UserProfile, value map[string]any) error {
	itemsRaw, _ := value["items"].([]any)
	var items []domain.MVDItem
	for _, item := range itemsRaw {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		items = append(items, domain.MVDItem{
			Domain:   domain.GoalID(getString(m, "domain")),
			Action:   getString(m, "action"),
			Duration: getString(m, "duration"),
			Criteria: getString(m, "criteria"),
		})
	}

	whenToUse, _ := value["when_to_use"].(string)
	if whenToUse == "" {
		whenToUse = "Quando tempo ou energia estiverem baixos — o objetivo é manter o hábito vivo, não render muito."
	}

	mvd := &domain.MinimumViableDaily{
		MVDID:     uuid.New(),
		UserID:    user.UserID,
		Items:     items,
		WhenToUse: whenToUse,
		UpdatedAt: time.Now(),
	}

	loc := user.Location()
	evt := domain.NewEvent(user.UserID, domain.EventMVDDefined, domain.SensitivityC5,
		map[string]any{"items_count": len(items)}, loc)
	_ = uc.events.Append(ctx, evt)

	return uc.mvds.Save(ctx, mvd)
}

func getString(m map[string]any, key string) string {
	v, _ := m[key].(string)
	return v
}

func (uc *OnboardingUseCase) buildPromptForStep(step domain.StepID, session *domain.OnboardingSession) *OnboardingPrompt {
	return &OnboardingPrompt{
		NextStep:   step,
		Message:    uc.messageForStep(step, session),
		Choices:    uc.choicesForStep(step),
		Disclosure: uc.disclosureForStep(step),
	}
}

func (uc *OnboardingUseCase) messageForStep(step domain.StepID, session *domain.OnboardingSession) string {
	switch step {
	case domain.StepWelcome:
		return "Olá! 👋 Sou seu assistente pessoal para metas anuais.\n\n" +
			"Vou fazer algumas perguntas rápidas para entender suas metas e criar um plano personalizado. " +
			"Leva poucos minutos e você pode parar e retomar quando quiser.\n\n" +
			"Vamos começar?"
	case domain.StepSelectGoals:
		return "Quais metas você quer trabalhar neste ciclo?\n\n" +
			"🎯 *Metas intensivas* (máx. 2):\n" +
			"• `english` — Fluência em inglês\n" +
			"• `java` — Evolução em Java\n\n" +
			"🏗 *Fundação* (sem limite):\n" +
			"• `sleep` — Dormir melhor\n" +
			"• `health` — Vida saudável\n" +
			"• `self_esteem` — Autoestima\n\n" +
			"📦 *Aposta semanal*:\n" +
			"• `saas` — Avançar no SaaS\n\n" +
			"Envie as metas separadas por vírgula (ex: english, java, sleep)"
	case domain.StepRestrictions:
		return "Qual é sua principal restrição de agenda?\n\n" +
			"Exemplos: \"trabalho 8h-18h\", \"tempo limitado à noite\", \"sem horário fixo\"\n\n" +
			"Se não souber agora, envie \"não sei\" que seguimos."
	case domain.StepSleepBaseline:
		return "Vamos registrar sua baseline de sono.\n\n" +
			"Responda brevemente:\n" +
			"1. Que horas dormiu ontem? (ex: 23:30)\n" +
			"2. Que horas acordou? (ex: 07:00)\n" +
			"3. Qualidade do sono (0-10)?\n" +
			"4. Energia pela manhã (0-10)?"
	case domain.StepEnglishBase:
		return "Baseline de inglês:\n\n" +
			"1. Como você avalia seu speaking? (0-10)\n" +
			"2. Como você avalia sua compreensão? (0-10)\n\n" +
			"(Se quiser, pode enviar um áudio curto em inglês para referência futura.)"
	case domain.StepJavaBaseline:
		return "Baseline de Java:\n\n" +
			"1. Como você avalia seu nível geral? (0-10)\n" +
			"2. Tópicos que domina? (ex: OOP, Collections)\n" +
			"3. Tópicos que quer melhorar? (ex: Streams, Testes)"
	case domain.StepMVD:
		return "Agora vamos definir seu *Mínimo Viável Diário (MVD)* — o plano para dias ruins.\n\n" +
			"O MVD é o mínimo que mantém seus hábitos vivos mesmo com pouca energia.\n\n" +
			"Vou sugerir um MVD baseado nas suas metas. Pode aceitar ou ajustar:\n\n" +
			uc.suggestMVD(session)
	case domain.StepPrivacy:
		return "🔒 *Privacidade*\n\n" +
			"Guardo apenas o mínimo para funcionar:\n" +
			"• Check-ins e planos (90 dias)\n" +
			"• Resultados de aprendizagem (90 dias)\n" +
			"• Conteúdo sensível como áudios (7 dias — ou pode desativar)\n" +
			"• Métricas agregadas (12 meses)\n\n" +
			"Você pode apagar qualquer dado a qualquer momento.\n\n" +
			"Quer desativar armazenamento de conteúdo sensível (áudios/textos pessoais)?\n" +
			"Responda: `sim` (modo mínimo) ou `não` (padrão)"
	case domain.StepSummary:
		return "Confirmação necessária — revise o resumo e envie `ok` para concluir."
	default:
		if session != nil && session.Status == domain.OnboardingMinimumCompleted {
			return "✅ Onboarding mínimo concluído! Você já pode usar a rotina diária.\n\n" +
				"Use /resumo para ver seu perfil e /rotina para começar."
		}
		return "Pronto para o próximo passo."
	}
}

func (uc *OnboardingUseCase) suggestMVD(session *domain.OnboardingSession) string {
	goalsAnswer := session.GetAnswer(domain.StepSelectGoals)
	if goalsAnswer == nil {
		return "• 5 min de leitura/listening em inglês\n• 10 min de prática de Java\n• Registrar sono"
	}

	suggestion := ""
	goalsRaw, _ := goalsAnswer.Value["goals"].([]any)
	for _, g := range goalsRaw {
		gStr, _ := g.(string)
		switch domain.GoalID(gStr) {
		case domain.GoalEnglish:
			suggestion += "• 5 min de listening em inglês\n"
		case domain.GoalJava:
			suggestion += "• 10 min de leitura/prática Java\n"
		case domain.GoalSleep:
			suggestion += "• Registrar sono (30s)\n"
		case domain.GoalHealth:
			suggestion += "• 5 min de caminhada ou alongamento\n"
		case domain.GoalSelfEsteem:
			suggestion += "• 1 registro rápido de gratidão\n"
		case domain.GoalSaaS:
			suggestion += "• Anotar 1 próximo passo do SaaS\n"
		}
	}
	return suggestion + "\nEnvie `aceitar` ou descreva seu MVD personalizado."
}

func (uc *OnboardingUseCase) choicesForStep(step domain.StepID) []string {
	switch step {
	case domain.StepWelcome:
		return []string{"Vamos lá!"}
	case domain.StepPrivacy:
		return []string{"sim", "não"}
	case domain.StepSummary:
		return []string{"ok"}
	default:
		return nil
	}
}

func (uc *OnboardingUseCase) disclosureForStep(step domain.StepID) string {
	switch step {
	case domain.StepSleepBaseline:
		return "Seus dados de sono são usados para calibrar planos e medir tendências. Retenção: 90 dias."
	case domain.StepEnglishBase:
		return "Sua autoavaliação é usada como referência inicial. Áudios (se enviados) ficam no máximo 7 dias."
	default:
		return ""
	}
}
