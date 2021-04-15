package powerusage

import (
	"github.com/cserrant/terosBattleServer/entity/powerusagecontext"
)

// CalculatePowerForecast applies logic to the context and shows the expected forecast of its effect.
func CalculatePowerForecast(context *powerusagecontext.PowerUsageContext) *powerusagecontext.PowerForecast {
	summary := &powerusagecontext.PowerForecast{
		UserSquaddieID:      context.ActingSquaddieID,
		PowerID:             context.PowerID,
		AttackPowerForecast: []*powerusagecontext.AttackingPowerForecast{},
	}

	for _, targetID := range context.TargetSquaddieIDs {
		summary.AttackPowerForecast = append(summary.AttackPowerForecast, GetExpectedDamage(context, &powerusagecontext.AttackContext{
			PowerID:           context.PowerID,
			AttackerID:        context.ActingSquaddieID,
			TargetID:          targetID,
			IsCounterAttack: false,
		}))
	}
	return summary
}
