package logging

import (
	"fmt"
	"time"

	"cloud.google.com/go/logging"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"
)

type contextLogger struct {
	echo.Context
	logger echo.Logger
}

func (l *contextLogger) Logger() echo.Logger {
	return l.logger
}

// LoggerMiddleware is appengine echo logger middleware.
type LoggerMiddleware struct {
	requestLogger     *logging.Logger
	applicationLogger *logging.Logger

	logLevel log.Lvl

	httpFormat *propagation.HTTPFormat

	moduleID  string
	projectID string
	versionID string
	zone      string
}

// NewLoggerMiddleware returns appengine echo logger middleware.
func NewLoggerMiddleware(client *logging.Client, logLevel log.Lvl, moduleID, projectID, versionID, zone string) *LoggerMiddleware {
	opt := logging.CommonResource(&mrpb.MonitoredResource{
		Type: "gae_app",
		Labels: map[string]string{
			"module_id":  moduleID,
			"project_id": projectID,
			"version_id": versionID,
			"zone":       zone,
		},
	})
	reqLogger := client.Logger(fmt.Sprintf("%s_request", moduleID), opt)
	appLogger := client.Logger(fmt.Sprintf("%s_application", moduleID), opt)
	return &LoggerMiddleware{
		requestLogger:     reqLogger,
		applicationLogger: appLogger,
		logLevel:          logLevel,
		httpFormat:        &propagation.HTTPFormat{},
		moduleID:          moduleID,
		projectID:         projectID,
		versionID:         versionID,
		zone:              zone,
	}
}

// Logger is appengine echo logger middleware.
// Set application logger to echo.Context and write request log.
func (m *LoggerMiddleware) Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		var trace, spanID string
		sc, ok := m.httpFormat.SpanContextFromRequest(req)
		if ok {
			trace = fmt.Sprintf("projects/%s/traces/%s", m.projectID, sc.TraceID)
			spanID = sc.SpanID.String()
		}

		appLogger := New(m.applicationLogger, trace, spanID)
		appLogger.SetLevel(m.logLevel)

		start := time.Now()
		if err := next(&contextLogger{Context: c, logger: appLogger}); err != nil {
			c.Error(err)
		}
		end := time.Now()

		resp := c.Response()
		m.requestLogger.Log(logging.Entry{
			Timestamp: time.Now(),
			Severity:  appLogger.maxSeverity,
			HTTPRequest: &logging.HTTPRequest{
				Request:      req,
				RequestSize:  req.ContentLength,
				Status:       resp.Status,
				ResponseSize: resp.Size,
				Latency:      end.Sub(start),
				RemoteIP:     c.RealIP(),
			},
			Trace:        trace,
			SpanID:       spanID,
			TraceSampled: true,
		})

		return nil
	}
}
