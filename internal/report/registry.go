package report

import "fmt"

func NewByName(name string) (Analyzer, error) {
	switch name {
	case "level":
		return NewLevelReport(), nil
	case "service":
		return NewServiceReport(), nil
	case "message":
		return NewMessageReport(), nil
	case "status":
		return NewStatusReport(), nil
	case "method":
		return NewMethodReport(), nil
	case "path":
		return NewPathReport(), nil
	case "host":
		return NewHostReport(), nil
	case "user":
		return NewUserReport(), nil
	case "region":
		return NewRegionReport(), nil
	case "environment":
		return NewEnvironmentReport(), nil
	case "trace":
		return NewTraceReport(), nil
	case "requestid":
		return NewRequestIDReport(), nil
	case "errorcode":
		return NewErrorCodeReport(), nil
	case "source":
		return NewSourceReport(), nil
	case "component":
		return NewComponentReport(), nil
	case "version":
		return NewVersionReport(), nil
	case "operation":
		return NewOperationReport(), nil
	case "tenant":
		return NewTenantReport(), nil
	case "device":
		return NewDeviceReport(), nil
	case "browser":
		return NewBrowserReport(), nil
	case "platform":
		return NewPlatformReport(), nil
	case "outcome":
		return NewOutcomeReport(), nil
	case "category":
		return NewCategoryReport(), nil
	case "queue":
		return NewQueueReport(), nil
	case "worker":
		return NewWorkerReport(), nil
	default:
		return nil, fmt.Errorf("unknown group: %s", name)
	}
}
