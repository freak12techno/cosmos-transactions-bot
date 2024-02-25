package telegram

import (
	"main/pkg/constants"

	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) HandleGetConfig(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got get config query")

	reporter.MetricsManager.LogReporterQuery(reporter.Name(), constants.ReporterQueryGetConfig)

	configString := reporter.Config.GetConfigAsString()
	return reporter.BotReply(c, configString)
}
