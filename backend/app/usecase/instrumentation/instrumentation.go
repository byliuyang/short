package instrumentation

import (
	"fmt"

	"github.com/short-d/app/fw"
	"github.com/short-d/short/app/entity"
)

// Instrumentation measures the internal operation of the system.
type Instrumentation struct {
	logger                          fw.Logger
	tracer                          fw.Tracer
	timer                           fw.Timer
	metrics                         fw.Metrics
	analytics                       fw.Analytics
	ctxCh                           chan fw.ExecutionContext
	redirectingAliasToLongLinkCh    chan fw.ExecutionContext
	redirectedAliasToLongLinkCh     chan fw.ExecutionContext
	longLinkRetrievalSucceedCh      chan fw.ExecutionContext
	longLinkRetrievalFailedCh       chan fw.ExecutionContext
	featureToggleRetrievalSucceedCh chan fw.ExecutionContext
	featureToggleRetrievalFailedCh  chan fw.ExecutionContext
	madeFeatureDecisionCh           chan fw.ExecutionContext
	trackCh                         chan fw.ExecutionContext
}

// RedirectingAliasToLongLink tracks RedirectingAliasToLongLink event.
func (i Instrumentation) RedirectingAliasToLongLink(alias string) {
	go func() {
		ctx := <-i.redirectingAliasToLongLinkCh
		userID := i.getUserID(nil, ctx)
		props := map[string]string{
			"request-id": ctx.RequestID,
			"alias":      alias,
		}
		i.analytics.Track("RedirectingAliasToLongLink", props, userID, ctx)
	}()
}

// RedirectedAliasToLongLink tracks RedirectedAliasToLongLink event.
func (i Instrumentation) RedirectedAliasToLongLink(url entity.URL) {
	go func() {
		ctx := <-i.redirectedAliasToLongLinkCh
		userID := i.getUserID(nil, ctx)
		props := map[string]string{
			"request-id": ctx.RequestID,
			"alias":      url.Alias,
			"long-link":  url.OriginalURL,
		}
		i.analytics.Track("RedirectedAliasToLongLink", props, userID, ctx)
	}()
}

// LongLinkRetrievalSucceed tracks the successes when retrieving long links.
func (i Instrumentation) LongLinkRetrievalSucceed() {
	go func() {
		ctx := <-i.longLinkRetrievalSucceedCh
		i.metrics.Count("long-link-retrieval-succeed", 1, 1, ctx)
	}()
}

// LongLinkRetrievalFailed tracks the failures when retrieving long links.
func (i Instrumentation) LongLinkRetrievalFailed(err error) {
	go func() {
		ctx := <-i.longLinkRetrievalFailedCh
		i.logger.Error(err)
		i.metrics.Count("long-link-retrieval-failed", 1, 1, ctx)
	}()
}

// FeatureToggleRetrievalSucceed tracks the successes when retrieving the status
// of the feature toggle.
func (i Instrumentation) FeatureToggleRetrievalSucceed() {
	go func() {
		ctx := <-i.featureToggleRetrievalSucceedCh
		i.metrics.Count("feature-toggle-retrieval-succeed", 1, 1, ctx)
	}()
}

// FeatureToggleRetrievalFailed tracks the failures when retrieving the status
// of the feature toggle.
func (i Instrumentation) FeatureToggleRetrievalFailed(err error) {
	go func() {
		ctx := <-i.featureToggleRetrievalFailedCh
		i.logger.Error(err)
		i.metrics.Count("feature-toggle-retrieval-failed", 1, 1, ctx)
	}()
}

// MadeFeatureDecision tracks MadeFeatureDecision event.
func (i Instrumentation) MadeFeatureDecision(
	featureID string,
	isEnabled bool,
) {
	go func() {
		ctx := <-i.madeFeatureDecisionCh
		userID := i.getUserID(nil, ctx)
		isEnabledStr := fmt.Sprintf("%v", isEnabled)
		props := map[string]string{
			"request-id": ctx.RequestID,
			"feature-id": featureID,
			"is-enabled": isEnabledStr,
		}
		i.analytics.Track("MadeFeatureDecision", props, userID, ctx)
	}()
}

// Track records events happened in the system.
func (i Instrumentation) Track(event string) {
	go func() {
		ctx := <-i.trackCh
		userID := i.getUserID(nil, ctx)
		props := map[string]string{}
		i.analytics.Track(event, props, userID, ctx)
	}()
}

// Done closes all the channels to prevent memory leak.
func (i Instrumentation) Done() {
	close(i.redirectingAliasToLongLinkCh)
	close(i.redirectedAliasToLongLinkCh)
	close(i.longLinkRetrievalSucceedCh)
	close(i.longLinkRetrievalFailedCh)
	close(i.featureToggleRetrievalSucceedCh)
	close(i.featureToggleRetrievalFailedCh)
}

func (i Instrumentation) getUserID(user *entity.User, ctx fw.ExecutionContext) string {
	if user == nil {
		return "anonymous"
	}
	return user.Email
}

// NewInstrumentation initializes instrumentation code.
func NewInstrumentation(logger fw.Logger,
	tracer fw.Tracer,
	timer fw.Timer,
	metrics fw.Metrics,
	analytics fw.Analytics,
	ctxCh chan fw.ExecutionContext,
) Instrumentation {
	redirectingAliasToLongLinkCh := make(chan fw.ExecutionContext)
	redirectedAliasToLongLinkCh := make(chan fw.ExecutionContext)
	longLinkRetrievalSucceedCh := make(chan fw.ExecutionContext)
	longLinkRetrievalFailedCh := make(chan fw.ExecutionContext)
	featureToggleRetrievalSucceedCh := make(chan fw.ExecutionContext)
	featureToggleRetrievalFailedCh := make(chan fw.ExecutionContext)
	madeFeatureDecisionCh := make(chan fw.ExecutionContext)
	trackCh := make(chan fw.ExecutionContext)

	ins := &Instrumentation{
		logger:                          logger,
		tracer:                          tracer,
		timer:                           timer,
		metrics:                         metrics,
		analytics:                       analytics,
		ctxCh:                           ctxCh,
		redirectingAliasToLongLinkCh:    redirectingAliasToLongLinkCh,
		redirectedAliasToLongLinkCh:     redirectedAliasToLongLinkCh,
		longLinkRetrievalSucceedCh:      longLinkRetrievalSucceedCh,
		longLinkRetrievalFailedCh:       longLinkRetrievalFailedCh,
		featureToggleRetrievalSucceedCh: featureToggleRetrievalSucceedCh,
		featureToggleRetrievalFailedCh:  featureToggleRetrievalFailedCh,
		madeFeatureDecisionCh:           madeFeatureDecisionCh,
		trackCh:                         trackCh,
	}
	go func() {
		ctx := <-ctxCh
		go func() { redirectingAliasToLongLinkCh <- ctx }()
		go func() { redirectedAliasToLongLinkCh <- ctx }()
		go func() { longLinkRetrievalSucceedCh <- ctx }()
		go func() { longLinkRetrievalFailedCh <- ctx }()
		go func() { featureToggleRetrievalSucceedCh <- ctx }()
		go func() { featureToggleRetrievalFailedCh <- ctx }()
		go func() { madeFeatureDecisionCh <- ctx }()
		go func() { trackCh <- ctx }()
		close(ctxCh)
	}()
	return *ins
}
