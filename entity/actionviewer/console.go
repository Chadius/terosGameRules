package actionviewer

import (
	"fmt"
	"github.com/chadius/terosbattleserver/entity/damagedistribution"
	"github.com/chadius/terosbattleserver/usecase/powerattackforecast"
	"github.com/chadius/terosbattleserver/usecase/powercommit"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"io"
)

// ConsoleActionViewerVerbosity represents options you can use to show how verbose you want the output.
type ConsoleActionViewerVerbosity struct {
	ShowRolls        bool
	ShowTargetStatus bool
}

// ConsoleActionViewer prints the results of actions to the console
type ConsoleActionViewer struct {
	Messages       []string
	IgnorePrinting bool
}

type messagesByPowerUsage struct {
	userAffectsTargetMessages []string
	targetStatusMessages      []string
	rollMessages              []string
	powerResults              []*powercommit.ResultPerTarget
}

// PrintMessages will print all the messages in the buffer, then clear the screen.
func (viewer *ConsoleActionViewer) PrintMessages(output io.Writer) {
	if viewer.IgnorePrinting {
		return
	}

	for _, message := range viewer.Messages {
		io.WriteString(output, message)
		io.WriteString(output, "\n")
	}

	viewer.Messages = []string{}
}

// PrintForecast will generate messages for the given Result and clear the Messages.
func (viewer *ConsoleActionViewer) PrintForecast(powerForecast *powerattackforecast.Forecast, repositories *repositories.RepositoryCollection, output io.Writer) {
	viewer.PrepareForecast(powerForecast, repositories)
	viewer.PrintMessages(output)
}

// PrepareForecast creates messages to show the attack preview.
func (viewer *ConsoleActionViewer) PrepareForecast(powerForecast *powerattackforecast.Forecast, repositories *repositories.RepositoryCollection) {
	for resultIndex, forecast := range powerForecast.ForecastedResultPerTarget {
		if forecast.Attack != nil {
			viewer.createMessagesForAttackOrCounterAttack(forecast, repositories, resultIndex, false)
		}

		if forecast.CounterAttack != nil {
			viewer.createMessagesForAttackOrCounterAttack(forecast, repositories, resultIndex, true)
		}

		if forecast.HealingForecast != nil {
			viewer.createMessagesForHealing(repositories, forecast, resultIndex)
		}
	}
}

func (viewer *ConsoleActionViewer) createMessagesForHealing(repositories *repositories.RepositoryCollection, forecast powerattackforecast.Calculation, resultIndex int) {
	healer := repositories.SquaddieRepo.GetSquaddieByID(forecast.Setup.UserID)
	target := repositories.SquaddieRepo.GetSquaddieByID(forecast.Setup.Targets[0])
	powerToUse := repositories.PowerRepo.GetPowerByID(forecast.Setup.PowerID)

	damageHealedDescription := fmt.Sprintf(", for %d healing", forecast.HealingForecast.RawHitPointsRestored)
	if forecast.HealingForecast.RawHitPointsRestored == 0 {
		damageHealedDescription = " for NO HEALING"
	}

	effectMessage := fmt.Sprintf("%s", damageHealedDescription)

	attackerAndPowerMessage := fmt.Sprintf("%s (%s)", healer.Name(), powerToUse.Name())
	if resultIndex > 0 {
		attackerAndPowerMessage = "- also"
	}

	attackMessage := fmt.Sprintf(
		"%s heals %s%s",
		attackerAndPowerMessage,
		target.Name(),
		effectMessage,
	)

	viewer.Messages = append(viewer.Messages, attackMessage)
}

func (viewer *ConsoleActionViewer) createMessagesForAttackOrCounterAttack(forecast powerattackforecast.Calculation, repositories *repositories.RepositoryCollection, resultIndex int, isACounterAttack bool) {
	attackForecast := forecast.Attack
	attackSetup := forecast.Setup
	if isACounterAttack {
		attackForecast = forecast.CounterAttack
		attackSetup = forecast.CounterAttackSetup
	}

	attackerHitBonus := attackForecast.VersusContext.ToHit
	chanceOutOf36 := getChanceToHitMessageSnippet(attackerHitBonus.ToHitBonus, true)
	effectMessage := getDamageDistributionMessageSnippet(attackForecast.VersusContext.NormalDamage)

	attacker := repositories.SquaddieRepo.GetSquaddieByID(attackSetup.UserID)
	target := repositories.SquaddieRepo.GetSquaddieByID(attackForecast.DefenderContext.TargetID)
	powerToUse := repositories.PowerRepo.GetPowerByID(attackSetup.PowerID)

	attackerAndPowerMessage := fmt.Sprintf("%s (%s)", attacker.Name(), powerToUse.Name())
	if resultIndex > 0 && !isACounterAttack {
		attackerAndPowerMessage = "- also"
	}

	versusMessage := "vs"
	if isACounterAttack {
		versusMessage = "counters"
	}

	attackMessage := fmt.Sprintf(
		"%s %s %s: %+d %s%s",
		attackerAndPowerMessage,
		versusMessage,
		target.Name(),
		attackerHitBonus.ToHitBonus,
		chanceOutOf36,
		effectMessage,
	)

	viewer.Messages = append(viewer.Messages, attackMessage)

	if attackForecast.VersusContext.CanCritical {
		critThreshold := attackForecast.VersusContext.ToHit.ToHitBonus - attackForecast.VersusContext.CriticalHitThreshold

		if critThreshold >= -5 {
			critChanceOutOf36 := getChanceToHitMessageSnippet(critThreshold, false)
			critEffectMessage := getDamageDistributionMessageSnippet(attackForecast.VersusContext.CriticalHitDamage)

			criticalHitAttackMessage := fmt.Sprintf(" crit: %s%s", critChanceOutOf36, critEffectMessage)
			viewer.Messages = append(viewer.Messages, criticalHitAttackMessage)
		}
	}
}

func getDamageDistributionMessageSnippet(damageDistribution *damagedistribution.DamageDistribution) string {
	damageTakenDescription := fmt.Sprintf(", for %d damage", damageDistribution.RawDamageDealt)
	if damageDistribution.RawDamageDealt == 0 {
		damageTakenDescription = " for NO DAMAGE"
	}

	barrierBurnDescription := ""
	if damageDistribution.TotalRawBarrierBurnt > 0 {
		barrierBurnDescription = fmt.Sprintf(" + %d barrier burn", damageDistribution.TotalRawBarrierBurnt)
	}

	if damageDistribution.IsFatalToTarget {
		damageTakenDescription = ", FATAL"
		barrierBurnDescription = ""
	}

	effectMessage := fmt.Sprintf("%s%s", damageTakenDescription, barrierBurnDescription)
	return effectMessage
}

func getChanceToHitMessageSnippet(toHitBonus int, includeParenthesis bool) string {
	toHitLookup := map[int]int{
		-5: 1,
		-4: 3,
		-3: 6,
		-2: 10,
		-1: 15,
		0:  21,
		1:  26,
		2:  30,
		3:  33,
		4:  35,
	}
	chanceOutOf36 := 0
	if toHitBonus > 4 {
		chanceOutOf36 = 36
	} else if toHitBonus < -5 {
		chanceOutOf36 = 0
	} else {
		chanceOutOf36 = toHitLookup[toHitBonus]
	}
	if includeParenthesis {
		return fmt.Sprintf("(%d/36)", chanceOutOf36)
	}
	return fmt.Sprintf("%d/36", chanceOutOf36)
}

// PrintResult will generate messages for the given Result and clear the Messages.
func (viewer *ConsoleActionViewer) PrintResult(powerResult *powercommit.Result, repositories *repositories.RepositoryCollection, verbosity *ConsoleActionViewerVerbosity, output io.Writer) {
	viewer.PrepareResult(powerResult, repositories, verbosity)
	viewer.PrintMessages(output)
}

// PrepareResult creates messages to show the attack result.
func (viewer *ConsoleActionViewer) PrepareResult(powerResult *powercommit.Result, repositories *repositories.RepositoryCollection, verbosity *ConsoleActionViewerVerbosity) {
	messagesPerPowerUsage := viewer.collatePowerResultPerTargetsByResult(powerResult)
	viewer.addUserAffectTargetMessagesByResult(messagesPerPowerUsage, repositories, verbosity)
	if verbosity != nil && verbosity.ShowRolls == true {
		viewer.addRollMessagesByResult(messagesPerPowerUsage, repositories)
	}
	if verbosity != nil && verbosity.ShowTargetStatus == true {
		viewer.addTargetStatusMessagesByResult(messagesPerPowerUsage, repositories)
	}
	viewer.printResultMessagesInOrder(messagesPerPowerUsage)
	viewer.Messages = append(viewer.Messages, "---")
}

func (viewer *ConsoleActionViewer) printResultMessagesInOrder(messagesPerPowerUsage []*messagesByPowerUsage) {
	for _, perGroupMessages := range messagesPerPowerUsage {
		for _, resultMessage := range perGroupMessages.userAffectsTargetMessages {
			viewer.Messages = append(viewer.Messages, resultMessage)
		}

		for _, targetStatusMessage := range perGroupMessages.targetStatusMessages {
			viewer.Messages = append(viewer.Messages, targetStatusMessage)
		}

		for _, rollMessage := range perGroupMessages.rollMessages {
			viewer.Messages = append(viewer.Messages, rollMessage)
		}
	}
}

func (viewer *ConsoleActionViewer) addRollMessagesByResult(messagesPerPowerUsage []*messagesByPowerUsage, repositories *repositories.RepositoryCollection) {
	for _, perGroupMessages := range messagesPerPowerUsage {
		rollMessages := []string{}
		for _, result := range perGroupMessages.powerResults {
			rollMessages = append(rollMessages, viewer.createRollMessage(result, repositories))
		}

		for _, message := range rollMessages {
			perGroupMessages.rollMessages = append(perGroupMessages.rollMessages, message)
		}
	}
}

func (viewer *ConsoleActionViewer) addTargetStatusMessagesByResult(messagesPerPowerUsage []*messagesByPowerUsage, repositories *repositories.RepositoryCollection) {
	for _, perGroupMessages := range messagesPerPowerUsage {
		targetStatusMessages := []string{}
		for _, result := range perGroupMessages.powerResults {
			targetStatusMessages = append(targetStatusMessages, viewer.createTargetStatusMessage(result, repositories))
		}

		for _, message := range targetStatusMessages {
			perGroupMessages.targetStatusMessages = append(perGroupMessages.targetStatusMessages, message)
		}
	}
}

func (viewer *ConsoleActionViewer) addUserAffectTargetMessagesByResult(messagesPerPowerUsage []*messagesByPowerUsage, repositories *repositories.RepositoryCollection, verbosity *ConsoleActionViewerVerbosity) {
	userCausedThePreviousResult := false
	for _, perGroupMessages := range messagesPerPowerUsage {
		userAffectsTargetMessages := []string{}
		for index, result := range perGroupMessages.powerResults {
			userCausedThePreviousResult = index != 0
			userAffectsTargetMessages = append(userAffectsTargetMessages, viewer.createMessageForResultPerTarget(result, repositories, userCausedThePreviousResult, verbosity))
		}

		for _, message := range userAffectsTargetMessages {
			perGroupMessages.userAffectsTargetMessages = append(perGroupMessages.userAffectsTargetMessages, message)
		}
	}
}

func (viewer *ConsoleActionViewer) collatePowerResultPerTargetsByResult(powerResult *powercommit.Result) []*messagesByPowerUsage {
	messagesPerPowerUsage := []*messagesByPowerUsage{}
	previousUserID := ""
	perTargetResultsFromTheSameResult := []*powercommit.ResultPerTarget{}
	for _, result := range powerResult.ResultPerTarget {
		if result.UserID != previousUserID {
			if previousUserID != "" {
				messagesPerPowerUsage = append(messagesPerPowerUsage, &messagesByPowerUsage{powerResults: perTargetResultsFromTheSameResult})
				perTargetResultsFromTheSameResult = nil
			}
		}
		perTargetResultsFromTheSameResult = append(perTargetResultsFromTheSameResult, result)
		previousUserID = result.UserID
	}
	messagesPerPowerUsage = append(messagesPerPowerUsage, &messagesByPowerUsage{powerResults: perTargetResultsFromTheSameResult})
	return messagesPerPowerUsage
}

func (viewer *ConsoleActionViewer) createRollMessage(result *powercommit.ResultPerTarget, repositories *repositories.RepositoryCollection) string {
	squaddieRepo := repositories.SquaddieRepo
	user := squaddieRepo.GetOriginalSquaddieByID(result.UserID)
	target := squaddieRepo.GetOriginalSquaddieByID(result.TargetID)

	if result.Attack != nil {
		return fmt.Sprintf(
			"   %s rolls %d + %d = %d, %s rolls %d + %d = %d",
			user.Name(), result.Attack.AttackRoll, result.Attack.AttackerToHitBonus, result.Attack.AttackerTotal,
			target.Name(), result.Attack.DefendRoll, result.Attack.DefenderToHitPenalty, result.Attack.DefenderTotal,
		)
	}
	return "   Auto-hit"
}

func (viewer *ConsoleActionViewer) createTargetStatusMessage(result *powercommit.ResultPerTarget, repositories *repositories.RepositoryCollection) string {
	squaddieRepo := repositories.SquaddieRepo
	target := squaddieRepo.GetOriginalSquaddieByID(result.TargetID)

	barrierMessage := ""
	if target.MaxBarrier() > 0 {
		barrierMessage = fmt.Sprintf(", %d barrier", target.CurrentBarrier())
	}

	targetStatusMessage := fmt.Sprintf("   %s: %d/%d HP%s",
		target.Name(),
		target.CurrentHitPoints(),
		target.MaxHitPoints(),
		barrierMessage,
	)

	return targetStatusMessage
}

func (viewer *ConsoleActionViewer) createMessageForResultPerTarget(result *powercommit.ResultPerTarget, repositories *repositories.RepositoryCollection, userCausedThePreviousResult bool, verbosity *ConsoleActionViewerVerbosity) string {
	if result.Attack == nil && result.Healing == nil {
		return "Unknown"
	}

	if result.Attack != nil {
		attackResultMessage := viewer.makeMessageForResultPerTargetAttackEffect(result, repositories, userCausedThePreviousResult, verbosity)
		return attackResultMessage
	}
	if result.Healing != nil {
		healingResultMessage := viewer.makeMessageForResultPerTargetHealingEffect(result, repositories, userCausedThePreviousResult, verbosity)
		return healingResultMessage
	}

	return "Unknown"
}

func (viewer *ConsoleActionViewer) makeMessageForResultPerTargetAttackEffect(result *powercommit.ResultPerTarget, repositories *repositories.RepositoryCollection, userCausedThePreviousResult bool, verbosity *ConsoleActionViewerVerbosity) string {
	squaddieRepo := repositories.SquaddieRepo
	target := squaddieRepo.GetOriginalSquaddieByID(result.TargetID)

	hitMessage := "misses"
	effectMessage := ""
	criticalHit := ""
	if result.Attack.CriticallyHitTarget {
		criticalHit = "CRITICALLY "
	}

	if result.Attack.HitTarget {
		if result.Attack.IsCounterAttack {
			hitMessage = "counters"
		} else {
			hitMessage = "hits"
		}
	}

	if result.Attack.HitTarget || result.Attack.CriticallyHitTarget {
		if target.Defense.IsDead() {
			effectMessage = ", felling"
		} else {
			damageTakenDescription := fmt.Sprintf(", for %d damage", result.Attack.Damage.ActualDamageTaken)
			barrierBurnDescription := ""
			if result.Attack.Damage.TotalRawBarrierBurnt > 0 {
				barrierBurnDescription = fmt.Sprintf(" + %d barrier burn", result.Attack.Damage.TotalRawBarrierBurnt)
			}

			effectMessage = fmt.Sprintf("%s%s", damageTakenDescription, barrierBurnDescription)
		}
	}

	userPrefix := viewer.getMessagePrefix(result, repositories, userCausedThePreviousResult, verbosity)

	return fmt.Sprintf("%s %s%s %s%s", userPrefix, criticalHit, hitMessage, target.Name(), effectMessage)
}

func (viewer *ConsoleActionViewer) makeMessageForResultPerTargetHealingEffect(result *powercommit.ResultPerTarget, repositories *repositories.RepositoryCollection, userCausedThePreviousResult bool, verbosity *ConsoleActionViewerVerbosity) string {
	squaddieRepo := repositories.SquaddieRepo
	target := squaddieRepo.GetOriginalSquaddieByID(result.TargetID)

	userPrefix := viewer.getMessagePrefix(result, repositories, userCausedThePreviousResult, verbosity)
	hitPointsRestored := fmt.Sprintf(", for %d healing", result.Healing.HitPointsRestored)
	if result.Healing.HitPointsRestored == 0 {
		hitPointsRestored = " for NO HEALING"
	}

	return fmt.Sprintf("%s heals %s%s", userPrefix, target.Name(), hitPointsRestored)
}

func (viewer *ConsoleActionViewer) getMessagePrefix(result *powercommit.ResultPerTarget, repositories *repositories.RepositoryCollection, userCausedThePreviousResult bool, verbosity *ConsoleActionViewerVerbosity) string {
	squaddieRepo := repositories.SquaddieRepo
	powerRepo := repositories.PowerRepo

	user := squaddieRepo.GetOriginalSquaddieByID(result.UserID)
	powerUsed := powerRepo.GetPowerByID(result.PowerID)

	userPrefix := fmt.Sprintf("%s (%s)", user.Name(), powerUsed.Name())
	if userCausedThePreviousResult {
		userPrefix = "- also"
	}

	return userPrefix
}
