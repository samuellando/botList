package list

import "fmt"

type hook struct {
	runner             Runner
	runnerAction       string
	runnerActionConfig string
}

type hookValues struct {
	Runner             *Runner
	RunnerAction       string
	RunnerActionConfig string
}

func createHook(fs ...func(*hookValues)) (hook, error) {
	v := hookValues{}
	for _, f := range fs {
		f(&v)
	}
	if v.Runner == nil  {
		return hook{}, fmt.Errorf("Runner must be set on hook")
	}
	if v.RunnerAction == "" {
		return hook{}, fmt.Errorf("Runner action must be set on hook")
	}
	if v.RunnerActionConfig == "" {
		return hook{}, fmt.Errorf("Runner action config must be set on hook")
	}
	return hook{
		runner:             *v.Runner,
		runnerAction:       v.RunnerAction,
		runnerActionConfig: v.RunnerActionConfig,
	}, nil
}

type Hook interface {
	Runner() Runner
	RunnerAction() string
	RunnerActionConfig() string
}

type ActionHook struct {
	hook         hook
	onActionType []string
}

type ActionHookValues struct {
	Runner             Runner
	RunnerAction       string
	RunnerActionConfig string
	OnActionType       []string
}

func CreateActionHook(fs ...func(*ActionHookValues)) (ActionHook, error) {
	v := ActionHookValues{}
	for _, f := range fs {
		f(&v)
	}
	h, err := createHook(func(hv *hookValues) {
		hv.Runner = &v.Runner
		hv.RunnerAction = v.RunnerAction
		hv.RunnerActionConfig = v.RunnerActionConfig
	})
	if err != nil {
		return ActionHook{}, err
	}
	return ActionHook{
		hook:         h,
		onActionType: v.OnActionType,
	}, nil
}

func (h ActionHook) Runner() Runner {
	return h.hook.runner
}

func (h ActionHook) RunnerAction() string {
	return h.hook.runnerAction
}

func (h ActionHook) RunnerActionConfig() string {
	return h.hook.runnerActionConfig
}

func (h ActionHook) OnActionType() []string {
	return h.onActionType
}

type CronHook struct {
	hook    hook
	cronTab string
}

type CronHookValues struct {
	Runner             Runner
	RunnerAction       string
	RunnerActionConfig string
	CronTab            string
}

func CreateCronHook(fs ...func(*CronHookValues)) (CronHook, error) {
	v := CronHookValues{}
	for _, f := range fs {
		f(&v)
	}
	h, err := createHook(func(hv *hookValues) {
		hv.Runner = &v.Runner
		hv.RunnerAction = v.RunnerAction
		hv.RunnerActionConfig = v.RunnerActionConfig
	})
	if err != nil {
		return CronHook{}, err
	}
	return CronHook{
		hook:    h,
		cronTab: v.CronTab,
	}, nil
}

func (h CronHook) Runner() Runner {
	return h.hook.runner
}

func (h CronHook) RunnerAction() string {
	return h.hook.runnerAction
}

func (h CronHook) RunnerActionConfig() string {
	return h.hook.runnerActionConfig
}

func (h CronHook) CronTab() string {
	return h.cronTab
}
