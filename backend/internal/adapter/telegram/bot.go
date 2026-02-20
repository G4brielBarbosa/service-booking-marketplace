package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/abriesouza/super-assistente/internal/domain"
	"github.com/abriesouza/super-assistente/internal/port"
	"github.com/abriesouza/super-assistente/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

type Bot struct {
	api          *tgbotapi.BotAPI
	onboarding   *usecase.OnboardingUseCase
	dailyRoutine *usecase.DailyRoutineUseCase
	gates        *usecase.GateUseCase
	idempotent   port.IdempotencyStore
	featureOn    bool
	dailyOn      bool
	gatesOn      bool
	log          *slog.Logger

	pendingEvidence map[int64]pendingEvidenceCtx
}

type pendingEvidenceCtx struct {
	TaskID    string
	ProfileID string
}

func NewBot(
	token string,
	onboarding *usecase.OnboardingUseCase,
	dailyRoutine *usecase.DailyRoutineUseCase,
	gates *usecase.GateUseCase,
	idempotent port.IdempotencyStore,
	featureOn bool,
	dailyOn bool,
	gatesOn bool,
	log *slog.Logger,
) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("creating telegram bot: %w", err)
	}

	log.Info("telegram bot authorized", "username", api.Self.UserName)

	return &Bot{
		api:             api,
		onboarding:      onboarding,
		dailyRoutine:    dailyRoutine,
		gates:           gates,
		idempotent:      idempotent,
		featureOn:       featureOn,
		dailyOn:         dailyOn,
		gatesOn:         gatesOn,
		log:             log,
		pendingEvidence: make(map[int64]pendingEvidenceCtx),
	}, nil
}

func (b *Bot) SendText(ctx context.Context, chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	_, err := b.api.Send(msg)
	return err
}

func (b *Bot) SendTextWithKeyboard(ctx context.Context, chatID int64, text string, options [][]string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"

	if len(options) > 0 {
		var rows [][]tgbotapi.KeyboardButton
		for _, row := range options {
			var buttons []tgbotapi.KeyboardButton
			for _, opt := range row {
				buttons = append(buttons, tgbotapi.NewKeyboardButton(opt))
			}
			rows = append(rows, buttons)
		}
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(rows...)
	}

	_, err := b.api.Send(msg)
	return err
}

func (b *Bot) StartPolling(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := b.api.GetUpdatesChan(u)

	b.log.Info("telegram long polling started")

	for {
		select {
		case <-ctx.Done():
			b.log.Info("telegram polling stopped")
			return
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			go b.handleUpdate(ctx, update)
		}
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	msg := update.Message
	chatID := msg.Chat.ID
	userID := msg.From.ID
	text := strings.TrimSpace(msg.Text)

	// Idempotency check per telegram update_id (PLAN-000 D-007)
	idempKey := fmt.Sprintf("tg_update_%d", update.UpdateID)
	existing, _ := b.idempotent.Check(ctx, idempKey)
	if existing != nil {
		b.log.Debug("duplicate update skipped", "update_id", update.UpdateID)
		return
	}

	b.log.Info("message received",
		"chat_id", chatID,
		"user_id", userID,
		"text_length", len(text),
	)

	if !b.featureOn {
		_ = b.SendText(ctx, chatID, "O assistente está em manutenção. Tente novamente em breve.")
		return
	}

	b.routeMessage(ctx, userID, chatID, text)

	// Store idempotency record (24h TTL)
	rec := domain.NewIdempotencyRecord(idempKey, "processed", 24*time.Hour)
	_ = b.idempotent.Store(ctx, rec)
}

func (b *Bot) routeMessage(ctx context.Context, userID, chatID int64, text string) {
	lower := strings.ToLower(text)

	switch {
	case lower == "/start" || lower == "vamos lá!":
		b.handleStart(ctx, userID, chatID)
	case lower == "/resumo" || lower == "/summary":
		b.handleSummary(ctx, userID, chatID)
	case lower == "/retomar" || lower == "/resume":
		b.handleResume(ctx, userID, chatID)
	case lower == "/privacidade" || lower == "/privacy":
		b.handlePrivacyInfo(ctx, userID, chatID)

	// Daily routine commands (PLAN-002)
	case b.dailyOn && (lower == "/checkin" || lower == "check-in"):
		b.handleCheckInStart(ctx, userID, chatID)
	case b.dailyOn && (lower == "/plano" || lower == "meu plano"):
		b.handleGetPlan(ctx, userID, chatID)
	case b.dailyOn && (lower == "/feito" || lower == "o que fiz"):
		b.handleGetSteps(ctx, userID, chatID)
	case b.dailyOn && lower == "/replanejar":
		b.handleReplan(ctx, userID, chatID, text)
	case b.dailyOn && strings.HasPrefix(lower, "iniciar "):
		b.handleTaskAction(ctx, userID, chatID, domain.ActionStart, strings.TrimPrefix(lower, "iniciar "))
	case b.dailyOn && strings.HasPrefix(lower, "bloquear "):
		b.handleTaskAction(ctx, userID, chatID, domain.ActionBlock, strings.TrimPrefix(lower, "bloquear "))
	case b.dailyOn && strings.HasPrefix(lower, "feito "):
		if b.gatesOn {
			b.handleGateAwareMarkDone(ctx, userID, chatID, strings.TrimPrefix(lower, "feito "))
		} else {
			b.handleTaskAction(ctx, userID, chatID, domain.ActionMarkDoneReq, strings.TrimPrefix(lower, "feito "))
		}
	case b.dailyOn && strings.HasPrefix(lower, "adiar "):
		b.handleTaskAction(ctx, userID, chatID, domain.ActionDefer, strings.TrimPrefix(lower, "adiar "))

	case b.gatesOn && (lower == "/gates" || lower == "gates"):
		b.handleGateSummary(ctx, userID, chatID)

	default:
		if b.gatesOn && b.tryHandleEvidenceResponse(ctx, userID, chatID, text) {
			return
		}
		if b.dailyOn && b.tryParseCheckInResponse(ctx, userID, chatID, text) {
			return
		}
		b.handleOnboardingAnswer(ctx, userID, chatID, text)
	}
}

func (b *Bot) handleStart(ctx context.Context, userID, chatID int64) {
	prompt, err := b.onboarding.StartOnboarding(ctx, userID, chatID)
	if err != nil {
		b.log.Error("start onboarding failed", "error", err)
		_ = b.SendText(ctx, chatID, "Erro ao iniciar. Tente novamente.")
		return
	}

	b.sendPrompt(ctx, chatID, prompt)
}

func (b *Bot) handleSummary(ctx context.Context, userID, chatID int64) {
	summary, err := b.onboarding.GetSummary(ctx, userID)
	if err != nil {
		_ = b.SendText(ctx, chatID, "Nenhum onboarding encontrado. Use /start para começar.")
		return
	}

	text := formatSummary(summary)
	_ = b.SendText(ctx, chatID, text)
}

func (b *Bot) handleResume(ctx context.Context, userID, chatID int64) {
	prompt, err := b.onboarding.ResumeOnboarding(ctx, userID)
	if err != nil {
		_ = b.SendText(ctx, chatID, "Nenhum onboarding em progresso. Use /start para começar.")
		return
	}

	b.sendPrompt(ctx, chatID, prompt)
}

func (b *Bot) handlePrivacyInfo(ctx context.Context, userID, chatID int64) {
	text := "🔒 *O que guardo e por quê:*\n\n" +
		"• *Check-ins/planos* — para calibrar seus dias (90 dias)\n" +
		"• *Resultados de aprendizagem* — para medir progresso (90 dias)\n" +
		"• *Conteúdo sensível* (áudios) — validação rápida (7 dias, desativável)\n" +
		"• *Métricas agregadas* — para tendências semanais (12 meses)\n" +
		"• *Configurações* — enquanto usar o sistema\n\n" +
		"Você pode apagar qualquer dado a qualquer momento.\n" +
		"Para desativar conteúdo sensível: /start e responder na etapa de privacidade."
	_ = b.SendText(ctx, chatID, text)
}

func (b *Bot) handleOnboardingAnswer(ctx context.Context, userID, chatID int64, text string) {
	// Determine current step from session
	user, err := b.onboarding.GetSummary(ctx, userID)
	if err != nil {
		_ = b.SendText(ctx, chatID, "Use /start para começar o onboarding.")
		return
	}
	_ = user // summary available but we need the session step

	// Try to resume to get current step
	prompt, err := b.onboarding.ResumeOnboarding(ctx, userID)
	if err != nil {
		_ = b.SendText(ctx, chatID, "Use /start para começar.")
		return
	}

	value := b.parseAnswerForStep(prompt.NextStep, text)
	progress, err := b.onboarding.SubmitAnswer(ctx, userID, prompt.NextStep, value)
	if err != nil {
		if domErr, ok := err.(*domain.DomainError); ok {
			_ = b.SendText(ctx, chatID, domErr.Message+"\n\n"+domErr.Detail)
		} else {
			b.log.Error("submit answer failed", "error", err)
			_ = b.SendText(ctx, chatID, "Erro ao processar resposta. Tente novamente.")
		}
		return
	}

	if progress.Status == domain.OnboardingMinimumCompleted {
		_ = b.SendText(ctx, chatID, "✅ *Onboarding mínimo concluído!*\n\n"+
			"Você já pode usar a rotina diária.\n"+
			"Use /resumo para ver seu perfil.")
		return
	}

	nextPrompt := &usecase.OnboardingPrompt{
		NextStep: progress.NextStep,
		Message:  progress.Message,
	}
	b.sendPrompt(ctx, chatID, nextPrompt)
}

func (b *Bot) parseAnswerForStep(step domain.StepID, text string) map[string]any {
	switch step {
	case domain.StepWelcome:
		return map[string]any{"confirmed": true}

	case domain.StepSelectGoals:
		parts := strings.Split(text, ",")
		var goals []any
		for _, p := range parts {
			g := strings.TrimSpace(strings.ToLower(p))
			if g != "" {
				goals = append(goals, g)
			}
		}
		return map[string]any{"goals": goals}

	case domain.StepRestrictions:
		return map[string]any{"main_restriction": text}

	case domain.StepSleepBaseline:
		return parseSleepBaseline(text)

	case domain.StepEnglishBase:
		return parseEnglishBaseline(text)

	case domain.StepJavaBaseline:
		return parseJavaBaseline(text)

	case domain.StepMVD:
		lower := strings.ToLower(text)
		if lower == "aceitar" || lower == "ok" {
			return map[string]any{"accepted_default": true}
		}
		return map[string]any{"custom_mvd": text}

	case domain.StepPrivacy:
		lower := strings.ToLower(text)
		return map[string]any{"opt_out_c3": lower == "sim"}

	case domain.StepSummary:
		return map[string]any{"confirmed": true}

	default:
		return map[string]any{"raw": text}
	}
}

func parseSleepBaseline(text string) map[string]any {
	lines := strings.Split(text, "\n")
	result := map[string]any{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Try to parse numbered answers
		if len(line) > 2 && (line[0] == '1' || line[0] == '2' || line[0] == '3' || line[0] == '4') {
			val := strings.TrimLeft(line[1:], ".): ")
			switch line[0] {
			case '1':
				result["sleep_time"] = val
			case '2':
				result["wake_time"] = val
			case '3':
				result["quality"] = val
			case '4':
				result["morning_energy"] = val
			}
		}
	}

	// Fallback: treat single-line as comma separated
	if len(result) == 0 {
		parts := strings.Split(text, ",")
		if len(parts) >= 2 {
			result["sleep_time"] = strings.TrimSpace(parts[0])
			result["wake_time"] = strings.TrimSpace(parts[1])
		}
		if len(parts) >= 3 {
			result["quality"] = strings.TrimSpace(parts[2])
		}
		if len(parts) >= 4 {
			result["morning_energy"] = strings.TrimSpace(parts[3])
		}
	}

	return result
}

func parseEnglishBaseline(text string) map[string]any {
	result := map[string]any{}
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 2 && (line[0] == '1' || line[0] == '2') {
			val := strings.TrimLeft(line[1:], ".): ")
			switch line[0] {
			case '1':
				result["speaking_self_eval"] = val
			case '2':
				result["comprehension_self_eval"] = val
			}
		}
	}

	if len(result) == 0 {
		parts := strings.Split(text, ",")
		if len(parts) >= 1 {
			result["speaking_self_eval"] = strings.TrimSpace(parts[0])
		}
		if len(parts) >= 2 {
			result["comprehension_self_eval"] = strings.TrimSpace(parts[1])
		}
	}

	return result
}

func parseJavaBaseline(text string) map[string]any {
	result := map[string]any{}
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 2 && (line[0] == '1' || line[0] == '2' || line[0] == '3') {
			val := strings.TrimLeft(line[1:], ".): ")
			switch line[0] {
			case '1':
				result["level_self_eval"] = val
			case '2':
				result["topics_strong"] = val
			case '3':
				result["topics_improve"] = val
			}
		}
	}

	if len(result) == 0 {
		result["raw_answer"] = text
	}

	return result
}

// --- Daily Routine handlers ---

func (b *Bot) handleCheckInStart(ctx context.Context, userID, chatID int64) {
	_ = b.SendTextWithKeyboard(ctx, chatID,
		"Bom dia! Vamos planejar o dia.\n\n"+
			"Quanto tempo você tem disponível hoje (minutos)?\n"+
			"E como está sua energia (0-10)?\n\n"+
			"Responda no formato: `tempo energia`\n"+
			"Exemplo: `60 7` (60 min, energia 7)",
		[][]string{{"15 3", "30 5", "60 7", "90 8"}},
	)
}

func (b *Bot) tryParseCheckInResponse(ctx context.Context, userID, chatID int64, text string) bool {
	parts := strings.Fields(text)
	if len(parts) < 2 {
		return false
	}

	timeMin, err := strconv.Atoi(parts[0])
	if err != nil {
		return false
	}

	energy := domain.NormalizeEnergy(parts[1])

	user, err := b.onboarding.GetSummary(ctx, userID)
	_ = user

	now := time.Now()
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	localDate := now.In(loc).Format("2006-01-02")

	view, err := b.dailyRoutine.SubmitDailyCheckIn(ctx, userID, localDate, timeMin, energy, nil, nil)
	if err != nil {
		if domErr, ok := err.(*domain.DomainError); ok {
			_ = b.SendText(ctx, chatID, domErr.Message)
		} else {
			b.log.Error("check-in failed", "error", err)
			_ = b.SendText(ctx, chatID, "Erro no check-in. Tente novamente.")
		}
		return true
	}

	_ = b.SendText(ctx, chatID, formatPlanView(view))
	return true
}

func (b *Bot) handleGetPlan(ctx context.Context, userID, chatID int64) {
	now := time.Now()
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	localDate := now.In(loc).Format("2006-01-02")

	view, err := b.dailyRoutine.GetTodayPlan(ctx, userID, localDate)
	if err != nil {
		_ = b.SendText(ctx, chatID, "Nenhum plano para hoje. Use /checkin para começar.")
		return
	}

	_ = b.SendText(ctx, chatID, formatPlanView(view))
}

func (b *Bot) handleGetSteps(ctx context.Context, userID, chatID int64) {
	now := time.Now()
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	localDate := now.In(loc).Format("2006-01-02")

	summary, err := b.dailyRoutine.GetTodayStepsSummary(ctx, userID, localDate)
	if err != nil {
		_ = b.SendText(ctx, chatID, "Nenhum registro para hoje. Use /checkin para começar.")
		return
	}

	_ = b.SendText(ctx, chatID, formatStepsSummary(summary))
}

func (b *Bot) handleReplan(ctx context.Context, userID, chatID int64, text string) {
	now := time.Now()
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	localDate := now.In(loc).Format("2006-01-02")

	parts := strings.Fields(text)
	var newTime, newEnergy *int
	if len(parts) >= 2 {
		if t, err := strconv.Atoi(parts[1]); err == nil {
			newTime = &t
		}
	}
	if len(parts) >= 3 {
		e := domain.NormalizeEnergy(parts[2])
		newEnergy = &e
	}

	view, explanation, err := b.dailyRoutine.ReplanDay(ctx, userID, localDate, newTime, newEnergy, nil)
	if err != nil {
		if domErr, ok := err.(*domain.DomainError); ok {
			_ = b.SendText(ctx, chatID, domErr.Message)
		} else {
			_ = b.SendText(ctx, chatID, "Erro ao replanejar. Tente /checkin.")
		}
		return
	}

	_ = b.SendText(ctx, chatID, explanation+"\n\n"+formatPlanView(view))
}

func (b *Bot) handleTaskAction(ctx context.Context, userID, chatID int64, action domain.TaskAction, taskQuery string) {
	now := time.Now()
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	localDate := now.In(loc).Format("2006-01-02")

	task, err := b.dailyRoutine.FindTaskByTitle(ctx, userID, localDate, strings.TrimSpace(taskQuery))
	if err != nil {
		_ = b.SendText(ctx, chatID, "Tarefa não encontrada. Use /plano para ver suas tarefas.")
		return
	}

	view, err := b.dailyRoutine.UpdateTaskStatus(ctx, userID, localDate, task.TaskID, action, "")
	if err != nil {
		if domErr, ok := err.(*domain.DomainError); ok {
			_ = b.SendText(ctx, chatID, domErr.Message+"\n"+domErr.Detail)
		} else {
			_ = b.SendText(ctx, chatID, "Erro ao atualizar tarefa.")
		}
		return
	}

	statusEmoji := statusIcon(view.Status)
	msg := fmt.Sprintf("%s *%s*\nStatus: %s %s", statusEmoji, task.Title, view.Status, statusEmoji)
	if view.NextStep != "" {
		msg += "\n\n" + view.NextStep
	}
	_ = b.SendText(ctx, chatID, msg)
}

func formatPlanView(view *domain.DailyPlanView) string {
	if view == nil {
		return "Nenhum plano gerado."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📋 *Plano %s do dia* (~%d min)\n", view.PlanType, view.TotalEstimatedMin))
	sb.WriteString(fmt.Sprintf("_%s_\n\n", view.Rationale))

	sb.WriteString(fmt.Sprintf("🎯 *Prioridade:* %s (%d min)\n", view.PriorityTask.Title, view.PriorityTask.EstimatedMin))
	sb.WriteString(fmt.Sprintf("   %s\n   ✅ %s\n", view.PriorityTask.Instructions, view.PriorityTask.DoneCriteria))

	if len(view.ComplementaryTasks) > 0 {
		sb.WriteString("\n📌 *Complementares:*\n")
		for _, t := range view.ComplementaryTasks {
			sb.WriteString(fmt.Sprintf("• %s (%d min) %s\n", t.Title, t.EstimatedMin, statusIcon(t.Status)))
		}
	}

	if view.FoundationTask != nil {
		sb.WriteString(fmt.Sprintf("\n🏗 *Fundação:* %s (%d min)\n", view.FoundationTask.Title, view.FoundationTask.EstimatedMin))
	}

	sb.WriteString("\n💡 Comece pela prioridade. Use `iniciar [tarefa]` e `feito [tarefa]`.")
	return sb.String()
}

func formatStepsSummary(s *domain.DailyStepsSummary) string {
	var sb strings.Builder
	sb.WriteString("📊 *Progresso de hoje*\n\n")

	if len(s.Done) > 0 {
		sb.WriteString("✅ *Concluídos:*\n")
		for _, t := range s.Done {
			sb.WriteString(fmt.Sprintf("  • %s\n", t.Title))
		}
		sb.WriteString("\n")
	}

	if len(s.InProgress) > 0 {
		sb.WriteString("▶️ *Em andamento:*\n")
		for _, t := range s.InProgress {
			sb.WriteString(fmt.Sprintf("  • %s\n", t.Title))
		}
		sb.WriteString("\n")
	}

	if len(s.Pending) > 0 {
		sb.WriteString("⏳ *Pendentes:*\n")
		for _, t := range s.Pending {
			sb.WriteString(fmt.Sprintf("  • %s (%d min)\n", t.Title, t.EstimatedMin))
		}
		sb.WriteString("\n")
	}

	if len(s.Blocked) > 0 {
		sb.WriteString("🚫 *Bloqueados:*\n")
		for _, t := range s.Blocked {
			reason := ""
			if t.BlockReason != nil {
				reason = " — " + *t.BlockReason
			}
			sb.WriteString(fmt.Sprintf("  • %s%s\n", t.Title, reason))
		}
	}

	total := len(s.Done) + len(s.InProgress) + len(s.Pending) + len(s.Blocked)
	if total == 0 {
		sb.WriteString("Nenhuma tarefa registrada. Use /checkin para começar.")
	}

	return sb.String()
}

func statusIcon(status domain.TaskStatus) string {
	switch status {
	case domain.TaskCompleted:
		return "✅"
	case domain.TaskInProgress:
		return "▶️"
	case domain.TaskBlocked:
		return "🚫"
	case domain.TaskDeferred:
		return "⏭"
	case domain.TaskEvidencePending:
		return "🔍"
	case domain.TaskAttempt:
		return "🔄"
	default:
		return "⏳"
	}
}

func (b *Bot) sendPrompt(ctx context.Context, chatID int64, prompt *usecase.OnboardingPrompt) {
	if len(prompt.Choices) > 0 {
		options := [][]string{prompt.Choices}
		_ = b.SendTextWithKeyboard(ctx, chatID, prompt.Message, options)
	} else {
		_ = b.SendText(ctx, chatID, prompt.Message)
	}

	if prompt.Disclosure != "" {
		_ = b.SendText(ctx, chatID, "ℹ️ "+prompt.Disclosure)
	}
}

// --- Quality Gates handlers (PLAN-003) ---

func (b *Bot) handleGateAwareMarkDone(ctx context.Context, userID, chatID int64, taskQuery string) {
	now := time.Now()
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	localDate := now.In(loc).Format("2006-01-02")

	task, err := b.dailyRoutine.FindTaskByTitle(ctx, userID, localDate, strings.TrimSpace(taskQuery))
	if err != nil {
		_ = b.SendText(ctx, chatID, "Tarefa não encontrada. Use /plano para ver suas tarefas.")
		return
	}

	result, err := b.gates.RequestTaskCompletion(ctx, userID, localDate, task.TaskID)
	if err != nil {
		if domErr, ok := err.(*domain.DomainError); ok {
			_ = b.SendText(ctx, chatID, domErr.Message)
		} else {
			b.log.Error("gate request completion failed", "error", err)
			_ = b.SendText(ctx, chatID, "Erro ao processar conclusão.")
		}
		return
	}

	switch result.Status {
	case domain.CompletionCompleted:
		_ = b.SendText(ctx, chatID, fmt.Sprintf("✅ *%s* concluída!", task.Title))

	case domain.CompletionAlreadyCompleted:
		_ = b.SendText(ctx, chatID, fmt.Sprintf("✅ *%s* já está concluída.", task.Title))

	case domain.CompletionEvidenceRequired:
		b.pendingEvidence[userID] = pendingEvidenceCtx{
			TaskID:    task.TaskID.String(),
			ProfileID: result.EvidenceRequest.ProfileID,
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("🔍 *%s* precisa de evidência para concluir.\n\n", task.Title))
		sb.WriteString("*O que preciso:*\n")
		for _, req := range result.EvidenceRequest.Requirements {
			sb.WriteString(fmt.Sprintf("• %s\n", req.Description))
		}
		sb.WriteString(fmt.Sprintf("\nℹ️ _%s_", result.EvidenceRequest.PrivacyDisclosure))
		sb.WriteString("\n\nEnvie sua evidência agora (texto ou áudio).")

		_ = b.SendText(ctx, chatID, sb.String())
	}
}

func (b *Bot) tryHandleEvidenceResponse(ctx context.Context, userID, chatID int64, text string) bool {
	pending, ok := b.pendingEvidence[userID]
	if !ok {
		return false
	}

	taskID, err := uuid.Parse(pending.TaskID)
	if err != nil {
		delete(b.pendingEvidence, userID)
		return false
	}

	profile := domain.LookupGateProfile(pending.ProfileID)
	sensitivity := domain.SensitivityC2
	kind := domain.EvidenceText
	if profile != nil && profile.EquivalencePolicy == domain.EquivAudioRequiredNoEquivalent {
		sensitivity = domain.SensitivityC3
	}

	summary := text
	if len(summary) > 100 {
		summary = summary[:100] + "..."
	}

	receipt, err := b.gates.SubmitEvidence(ctx, userID, taskID, kind, sensitivity, summary, &text)
	if err != nil {
		b.log.Error("submit evidence failed", "error", err)
		_ = b.SendText(ctx, chatID, "Erro ao processar evidência. Tente novamente.")
		return true
	}

	if !receipt.Valid {
		_ = b.SendText(ctx, chatID, "Evidência inválida. "+receipt.Reason+"\nTente enviar novamente.")
		return true
	}

	now := time.Now()
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	localDate := now.In(loc).Format("2006-01-02")

	grView, err := b.gates.EvaluateGate(ctx, userID, localDate, taskID)
	if err != nil {
		b.log.Error("evaluate gate failed", "error", err)
		_ = b.SendText(ctx, chatID, "Erro ao avaliar evidência.")
		delete(b.pendingEvidence, userID)
		return true
	}

	delete(b.pendingEvidence, userID)

	if grView.GateStatus == domain.GateSatisfied {
		_ = b.SendText(ctx, chatID, "✅ *Concluído!* Gate satisfeito.")
	} else {
		msg := fmt.Sprintf("⚠️ %s\n\n💡 *Próximo passo:* %s", grView.ReasonShort, grView.NextMinStep)
		_ = b.SendText(ctx, chatID, msg)

		taskIDStr := taskID.String()
		b.pendingEvidence[userID] = pendingEvidenceCtx{
			TaskID:    taskIDStr,
			ProfileID: pending.ProfileID,
		}
	}

	return true
}

func (b *Bot) handleGateSummary(ctx context.Context, userID, chatID int64) {
	now := time.Now()
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	localDate := now.In(loc).Format("2006-01-02")

	items, err := b.gates.GetTodayGateSummary(ctx, userID, localDate)
	if err != nil {
		_ = b.SendText(ctx, chatID, "Nenhum gate para hoje. Use /checkin para começar.")
		return
	}

	if len(items) == 0 {
		_ = b.SendText(ctx, chatID, "Nenhuma tarefa com gate hoje.")
		return
	}

	var sb strings.Builder
	sb.WriteString("🔍 *Gates de hoje*\n\n")
	for _, item := range items {
		icon := "⏳"
		if item.GateStatus == domain.GateSatisfied {
			icon = "✅"
		} else if item.GateStatus == domain.GateNotSatisfied && item.ReasonShort != "Aguardando evidência." {
			icon = "⚠️"
		}
		sb.WriteString(fmt.Sprintf("%s *%s*: %s\n", icon, item.TaskTitle, item.ReasonShort))
		if item.NextMinStep != "" && item.GateStatus != domain.GateSatisfied {
			sb.WriteString(fmt.Sprintf("   💡 %s\n", item.NextMinStep))
		}
	}

	_ = b.SendText(ctx, chatID, sb.String())
}

func formatSummary(s *domain.OnboardingSummary) string {
	var sb strings.Builder
	sb.WriteString("📋 *Resumo do Onboarding*\n\n")

	if len(s.ActiveGoals) > 0 {
		sb.WriteString("🎯 *Metas ativas:*\n")
		for _, g := range s.ActiveGoals {
			class := domain.GoalClassifications[g.ID]
			sb.WriteString(fmt.Sprintf("  • %s (%s)\n", g.ID, class))
		}
		sb.WriteString("\n")
	}

	if len(s.PausedGoals) > 0 {
		sb.WriteString("⏸ *Metas pausadas:*\n")
		for _, g := range s.PausedGoals {
			sb.WriteString(fmt.Sprintf("  • %s\n", g.ID))
		}
		sb.WriteString("\n")
	}

	if s.Restrictions != nil {
		if r, ok := s.Restrictions["main_restriction"].(string); ok && r != "" {
			sb.WriteString(fmt.Sprintf("⏰ *Restrição:* %s\n\n", r))
		}
	}

	if len(s.Baselines) > 0 {
		sb.WriteString("📊 *Baselines:*\n")
		for _, b := range s.Baselines {
			sb.WriteString(fmt.Sprintf("  • %s: %s\n", b.Domain, b.Completeness))
		}
		sb.WriteString("\n")
	}

	if s.MVD != nil && len(s.MVD.Items) > 0 {
		sb.WriteString("🔋 *MVD (dia ruim):*\n")
		for _, item := range s.MVD.Items {
			sb.WriteString(fmt.Sprintf("  • %s (%s)\n", item.Action, item.Duration))
		}
		sb.WriteString("\n")
	}

	if len(s.PendingItems) > 0 {
		sb.WriteString("⏳ *Pendências:*\n")
		for _, p := range s.PendingItems {
			sb.WriteString(fmt.Sprintf("  • %s\n", p.Description))
		}
		sb.WriteString("\n")
	}

	if s.PrivacyPolicy != nil {
		mode := "padrão"
		if s.PrivacyPolicy.MinimalMode {
			mode = "mínimo (sem conteúdo sensível)"
		}
		sb.WriteString(fmt.Sprintf("🔒 *Privacidade:* modo %s\n", mode))
	}

	return sb.String()
}
