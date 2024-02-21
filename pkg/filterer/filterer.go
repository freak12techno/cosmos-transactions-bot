package filterer

import (
	"fmt"
	configPkg "main/pkg/config"
	configTypes "main/pkg/config/types"
	messagesPkg "main/pkg/messages"
	metricsPkg "main/pkg/metrics"
	"main/pkg/types"
	"strconv"

	"github.com/rs/zerolog"
)

type Filterer struct {
	Logger           zerolog.Logger
	MetricsManager   *metricsPkg.Manager
	Config           *configPkg.AppConfig
	lastBlockHeights map[string]int64
}

func NewFilterer(
	logger *zerolog.Logger,
	config *configPkg.AppConfig,
	metricsManager *metricsPkg.Manager,
) *Filterer {
	return &Filterer{
		Logger: logger.With().
			Str("component", "filterer").
			Logger(),
		MetricsManager:   metricsManager,
		Config:           config,
		lastBlockHeights: map[string]int64{},
	}
}

func (f *Filterer) GetReportableForReporters(
	report types.Report,
) map[string]types.Report {
	reportables := make(map[string]types.Report)

	for _, subscription := range f.Config.Subscriptions {
		chain := f.Config.Chains.FindByName(subscription.Chain)

		reportableFiltered := f.FilterForChainAndSubscription(
			report.Reportable,
			chain,
			subscription,
		)

		if reportableFiltered != nil {
			f.Logger.Info().
				Str("type", report.Reportable.Type()).
				Str("subscription_name", subscription.Name).
				Msg("Got report for subscription")
			reportables[subscription.Reporter] = types.Report{
				Chain:        report.Chain,
				Node:         report.Node,
				Reportable:   reportableFiltered,
				Subscription: subscription,
			}
		}
	}

	return reportables
}

func (f *Filterer) FilterForChainAndSubscription(
	reportable types.Reportable,
	chain *configTypes.Chain,
	subscription *configTypes.Subscription,
) types.Reportable {
	// Filtering out TxError only if chain's log-node-errors = true.
	if _, ok := reportable.(*types.TxError); ok {
		if !subscription.LogNodeErrors {
			f.MetricsManager.LogFilteredEvent(subscription.Name, reportable.Type())
			f.Logger.Debug().Msg("Got transaction error, skipping as node errors logging is disabled")
			return nil
		}

		f.MetricsManager.LogMatchedEvent(subscription.Name, reportable.Type())
		return reportable
	}

	if _, ok := reportable.(*types.NodeConnectError); ok {
		if !subscription.LogNodeErrors {
			f.MetricsManager.LogFilteredEvent(subscription.Name, reportable.Type())
			f.Logger.Debug().Msg("Got node error, skipping as node errors logging is disabled")
			return nil
		}

		f.MetricsManager.LogMatchedEvent(subscription.Name, reportable.Type())
		return reportable
	}

	tx, ok := reportable.(*types.Tx)
	if !ok {
		f.Logger.Error().Str("type", reportable.Type()).Msg("Unsupported reportable type, ignoring.")
		f.MetricsManager.LogFilteredEvent(subscription.Name, reportable.Type())
		return nil
	}

	if !subscription.LogFailedTransactions && tx.Code > 0 {
		f.Logger.Debug().
			Str("hash", tx.GetHash()).
			Msg("Transaction is failed, skipping")
		f.MetricsManager.LogFilteredEvent(subscription.Name, reportable.Type())
		return nil
	}

	txHeight, err := strconv.ParseInt(tx.Height.Value, 10, 64)
	if err != nil {
		f.Logger.Fatal().Err(err).Msg("Error converting height to int64")
	}

	chainLastBlockHeight, ok := f.lastBlockHeights[chain.Name]
	if ok && chainLastBlockHeight > txHeight {
		f.Logger.Debug().
			Str("hash", tx.GetHash()).
			Int64("height", txHeight).
			Int64("last_height", chainLastBlockHeight).
			Msg("Transaction height is less than the last one received, skipping")
		return nil
	}

	if !ok || chainLastBlockHeight < txHeight {
		f.lastBlockHeights[chain.Name] = txHeight
	}

	messages := make([]types.Message, 0)

	for _, message := range tx.Messages {
		filteredMessage := f.FilterMessage(message, subscription, false)
		if filteredMessage != nil {
			messages = append(messages, filteredMessage)
		}
	}

	if len(messages) == 0 {
		f.Logger.Debug().
			Str("hash", tx.GetHash()).
			Msg("All messages in transaction were filtered out, skipping.")
		f.MetricsManager.LogFilteredEvent(subscription.Name, reportable.Type())
		return nil
	}

	tx.Messages = messages
	f.MetricsManager.LogMatchedEvent(subscription.Name, reportable.Type())
	return tx
}

func (f *Filterer) FilterMessage(
	message types.Message,
	subscription *configTypes.Subscription,
	internal bool,
) types.Message {
	if unsupportedMsg, ok := message.(*messagesPkg.MsgUnsupportedMessage); ok {
		if subscription.LogUnknownMessages {
			f.Logger.Error().Str("type", unsupportedMsg.MsgType).Msg("Unsupported message type")
			return message
		} else {
			f.Logger.Debug().Str("type", unsupportedMsg.MsgType).Msg("Unsupported message type")
			return nil
		}
	}

	if unparsedMsg, ok := message.(*messagesPkg.MsgUnparsedMessage); ok {
		if subscription.LogUnparsedMessages {
			f.Logger.Error().Err(unparsedMsg.Error).Str("type", unparsedMsg.MsgType).Msg("Error parsing message")
			return message
		}

		f.Logger.Debug().
			Err(unparsedMsg.Error).
			Str("type", unparsedMsg.MsgType).
			Msg("Not logging unparsed messages, skipping.")
		return nil
	}

	// internal -> filter only if subscription.FilterInternalMessages is true
	// !internal -> filter regardless
	if !internal || subscription.FilterInternalMessages {
		matches, err := subscription.Filters.Matches(message.GetValues())

		f.Logger.Trace().
			Str("type", message.Type()).
			Str("values", fmt.Sprintf("%+v", message.GetValues().ToMap())).
			Str("filters", fmt.Sprintf("%+v", subscription.Filters)).
			Bool("matches", matches).
			Msg("Result of matching message events against filters")

		if err != nil {
			f.Logger.Error().
				Err(err).
				Str("type", message.Type()).
				Msg("Error checking if message matches filters")
		} else if !matches {
			f.Logger.Debug().
				Str("type", message.Type()).
				Msg("Message is ignored by filters.")
			return nil
		}
	}

	if len(message.GetRawMessages()) == 0 {
		return message
	}

	parsedInternalMessages := make([]types.Message, 0)

	// Processing internal messages (such as ones in MsgExec)
	for _, internalMessage := range message.GetParsedMessages() {
		if internalMessageParsed := f.FilterMessage(internalMessage, subscription, true); internalMessageParsed != nil {
			parsedInternalMessages = append(parsedInternalMessages, internalMessageParsed)
		}
	}

	if len(parsedInternalMessages) == 0 {
		f.Logger.Debug().
			Str("type", message.Type()).
			Msg("Message with messages inside has 0 messages after filtering, skipping.")
		return nil
	}

	message.SetParsedMessages(parsedInternalMessages)
	return message
}
