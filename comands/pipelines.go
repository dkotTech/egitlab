package comands

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dkotTech/egitlab/internal"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/urfave/cli/v2"
)

const (
	defaultLogsCount = 3
)

func NewPipelinesCommand() *cli.Command {
	commandFlags := []cli.Flag{
		&cli.DurationFlag{
			Name:  "update-interval",
			Value: time.Second * 30,
			Usage: "pipeline status update interval",
		},
		&cli.StringSliceFlag{
			Name:  "create-job",
			Value: nil,
			Usage: "create a job on update by name, if job in manual state",
		},
		&cli.StringSliceFlag{
			Name:  "retry-job",
			Value: nil,
			Usage: "retry a job on update by name, if job in failed or success state",
		},
	}

	var credentials internal.Credentials
	var serviceInfo internal.ServiceInfo
	var httpSettings internal.HttpSettings

	commandFlags = append(commandFlags, credentials.CliFlags()...)
	commandFlags = append(commandFlags, serviceInfo.CliFlags()...)
	commandFlags = append(commandFlags, httpSettings.CliFlags()...)

	return &cli.Command{
		Name:    "pipelines",
		Aliases: []string{"p"},
		Usage:   "get a pipeline status",
		Flags:   commandFlags,
		Action: func(cCtx *cli.Context) error {
			err := credentials.Parse(cCtx)
			if err != nil {
				return err
			}

			err = serviceInfo.Parse(cCtx)
			if err != nil {
				return err
			}

			err = httpSettings.Parse(cCtx)
			if err != nil {
				return err
			}

			updateRequestDuration := cCtx.Duration("update-interval")
			createJobs := cCtx.StringSlice("create-job")
			retryJobs := cCtx.StringSlice("retry-job")

			pipeline, err := internal.GetLatestPipeline(serviceInfo.Host, serviceInfo.ProjectName, serviceInfo.RefName, credentials)
			if err != nil {
				return err
			}

			if _, err := tea.NewProgram(pipelinesApp{
				headers:      []string{},
				rows:         [][]string{},
				pipelineInfo: pipeline,
				serviceInfo:  serviceInfo,
				credentials:  credentials,
				buf:          new(strings.Builder),

				timerUpdate: updateRequestDuration,
				timer:       timer.NewWithInterval(time.Millisecond*100, time.Millisecond),

				extraButtons: make(map[string]func(pipelinesApp) (pipelinesApp, tea.Cmd)),

				createJobs: createJobs,
				retryJobs:  retryJobs,
			}, tea.WithAltScreen()).Run(); err != nil {
				return err
			}

			return nil
		},
	}
}

type pipelinesApp struct {
	headers []string
	rows    [][]string

	pipelineInfo internal.GitlabPipeline
	serviceInfo  internal.ServiceInfo
	credentials  internal.Credentials

	timerUpdate time.Duration
	timer       timer.Model
	buf         *strings.Builder

	lastRequestPipelineTime time.Time

	extraButtons map[string]func(pipelinesApp) (pipelinesApp, tea.Cmd)

	createJobs []string
	retryJobs  []string

	jobCreatedLogs []string
}

func (m pipelinesApp) Init() tea.Cmd {
	return tea.Batch(m.timer.Init(), tea.SetWindowTitle(m.credentials.Service.String()+" pipelines"))
}

func (m pipelinesApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd
	case timer.TimeoutMsg:
		m.timer.Timeout = m.timerUpdate

		mNew, err := requestJobs(m)
		if err != nil {
			panic(err)
		}

		return mNew, mNew.timer.Start()
	case tea.KeyMsg:
		key := msg.String()

		switch key {
		case "ctrl+c", "q", "й":
			return m, tea.Quit
		case "u", "г":
			m.timer.Timeout = time.Millisecond
			return m, m.timer.Start()
		}

		m.jobCreatedLogs = append(m.jobCreatedLogs, key)

		if call, ok := m.extraButtons[key]; ok {
			m.jobCreatedLogs = append(m.jobCreatedLogs, "ok")
			return call(m)
		}

	}

	return m, nil
}

func requestJobs(m pipelinesApp) (pipelinesApp, error) {
	jobs, err := getGitlabJobs(m)
	if err != nil {
		return pipelinesApp{}, err
	}

	runnableJobsList, err := collectJobsForManualRun(m, jobs)
	if err != nil {
		return pipelinesApp{}, err
	}

	if runJobsNow(m, runnableJobsList) {
		// get update from gitlab after jobs run
		jobs, err = getGitlabJobs(m)
		if err != nil {
			return pipelinesApp{}, err
		}

		runnableJobsList, err = collectJobsForManualRun(m, jobs)
		if err != nil {
			return pipelinesApp{}, err
		}
	}

	stages := make(map[string][]internal.GitlabJob)
	stagesList := make([]stageInfo, 0)

	for i, jb := range jobs {
		if _, ok := stages[jb.Stage]; !ok {
			stagesList = append(stagesList, stageInfo{jb.Stage, i})
		}
		stages[jb.Stage] = append(stages[jb.Stage], jb)
	}

	m.headers = func() []string {
		sort.Slice(stagesList, func(i, j int) bool {
			return stagesList[i].order > stagesList[j].order
		})

		list := make([]string, 0, len(stagesList))
		for _, info := range stagesList {
			list = append(list, info.name)
		}

		return list
	}()

	maxJobsLen := 0
	for _, jbs := range stages {
		if len(jbs) > maxJobsLen {
			maxJobsLen = len(jbs)
		}
	}

	headersLen := len(m.headers)

	m.rows = make([][]string, 0, maxJobsLen)
	for j := 0; j < maxJobsLen; j++ {
		m.rows = append(m.rows, make([]string, headersLen))
	}

	for i, header := range m.headers {
		jbs, ok := stages[header]
		if !ok {
			continue
		}

		for j, jb := range jbs {
			var dur time.Duration

			if jb.Duration != nil {
				dur, err = time.ParseDuration(fmt.Sprintf("%fs", *jb.Duration))
				if err != nil {
					return pipelinesApp{}, err
				}
			}

			jbUrl := internal.EncodeHyperlink(jb.WebUrl, "Open")

			m.rows[j][i] = lipgloss.JoinVertical(
				lipgloss.Left,
				fmt.Sprint("> ", jb.Name),
				internal.StatusStyles[jb.Status].Emoji+" "+internal.StatusStyles[jb.Status].Style.Render(jb.Status),
				jbUrl,
				dur.Round(time.Second).String(),
				runnableJobsList.GetFormatButton(jb.Status, jb.Name),
			)
		}
	}

	m.lastRequestPipelineTime = time.Now()

	m = addCreateJobsExtraButtons(m, runnableJobsList)

	return m, nil
}

func getGitlabJobs(m pipelinesApp) ([]internal.GitlabJob, error) {
	jobs, err := internal.GetLatestJobs(m.serviceInfo.Host, m.serviceInfo.ProjectName, m.pipelineInfo.Id, m.credentials)
	if err != nil {
		return nil, err
	}

	triggeredJobs, err := internal.GetLatestTriggeredJobs(m.serviceInfo.Host, m.serviceInfo.ProjectName, m.pipelineInfo.Id, m.credentials)
	if err != nil {
		return nil, err
	}

	jobs = append(jobs, triggeredJobs...)

	return jobs, nil
}

func (m pipelinesApp) View() string {
	t := table.New().
		Border(lipgloss.ThickBorder()).
		BorderStyle(internal.TableBorderStyle).
		Headers(m.headers...).
		BorderRow(true).
		BorderColumn(true).
		Rows(m.rows...)

	m.buf.Reset()

	m.buf.WriteString(m.serviceInfo.RefName + " ")
	m.buf.WriteString(internal.EncodeHyperlink(m.pipelineInfo.WebUrl, "Open pipeline"))
	m.buf.WriteString("\n")

	m.buf.WriteString(t.Render())
	m.buf.WriteString("\n")

	m.buf.WriteString("next update in: ")
	m.buf.WriteString(m.timer.View())
	m.buf.WriteString("\n")

	m.buf.WriteString("last update: ")
	m.buf.WriteString(m.lastRequestPipelineTime.Format(time.RFC3339))
	m.buf.WriteString("\n")

	logsCounts := 0
	for i := len(m.jobCreatedLogs) - 1; i >= 0; i-- {
		if logsCounts == defaultLogsCount {
			break
		}

		m.buf.WriteString(internal.HelpStyle.Render(m.jobCreatedLogs[i]))
		m.buf.WriteString("\n")

		logsCounts++
	}

	m.buf.WriteString(internal.HelpStyle.Render("q: exit | u: update now"))
	m.buf.WriteString("\n")

	return m.buf.String()
}

func addCreateJobsExtraButtons(m pipelinesApp, runnableJobsList runnableJobs) pipelinesApp {
	for _, jb := range runnableJobsList.manuals {
		m.extraButtons[jb.Button] = func(app pipelinesApp) (pipelinesApp, tea.Cmd) {
			app.timer.Timeout = time.Millisecond
			app.jobCreatedLogs = append(app.jobCreatedLogs, jb.MsgLog(jb.Run()))
			return app, app.timer.Start()
		}
	}

	for _, jb := range runnableJobsList.failed {
		m.extraButtons[jb.Button] = func(app pipelinesApp) (pipelinesApp, tea.Cmd) {
			app.timer.Timeout = time.Millisecond
			app.jobCreatedLogs = append(app.jobCreatedLogs, jb.MsgLog(jb.Run()))
			return app, app.timer.Start()
		}
	}

	for _, jb := range runnableJobsList.success {
		m.extraButtons[jb.Button] = func(app pipelinesApp) (pipelinesApp, tea.Cmd) {
			app.timer.Timeout = time.Millisecond
			app.jobCreatedLogs = append(app.jobCreatedLogs, jb.MsgLog(jb.Run()))
			return app, app.timer.Start()
		}
	}

	return m
}

func runJobsNow(m pipelinesApp, runnableJobsList runnableJobs) bool {
	var update bool

	for _, jobName := range m.createJobs {
		if found, ok := runnableJobsList.manuals[jobName]; ok {
			update = true
			m.jobCreatedLogs = append(m.jobCreatedLogs, found.MsgLog(found.Run()))
		}
	}

	m.createJobs = m.createJobs[:0]

	for _, jobName := range m.retryJobs {
		if found, ok := runnableJobsList.failed[jobName]; ok {
			update = true
			m.jobCreatedLogs = append(m.jobCreatedLogs, found.MsgLog(found.Run()))
		}

		if found, ok := runnableJobsList.success[jobName]; ok {
			update = true
			m.jobCreatedLogs = append(m.jobCreatedLogs, found.MsgLog(found.Run()))
		}
	}

	m.retryJobs = m.retryJobs[:0]

	return update
}

type stageInfo struct {
	name  string
	order int
}

type jobInfo[T any] struct {
	Run    func() (T, error)
	Button string

	MsgLog func(T, error) string
}

func collectJobsForManualRun(m pipelinesApp, jobs []internal.GitlabJob) (runnableJobs, error) {
	runableJobsInfo := runnableJobs{
		manuals: make(jobsByName[internal.RunJobResponse]),
		failed:  make(jobsByName[internal.RetryJobResponse]),
		success: make(jobsByName[internal.RetryJobResponse]),
	}

	altKey := 1

	for _, jb := range jobs {
		switch jb.Status {
		case "manual":
			runableJobsInfo.manuals[jb.Name] = jobInfo[internal.RunJobResponse]{
				Run: func() (internal.RunJobResponse, error) {
					return internal.RunJob(m.serviceInfo.Host, m.serviceInfo.ProjectName, jb.ID, m.credentials)
				},

				Button: fmt.Sprintf("alt+%d", altKey),
				MsgLog: func(r internal.RunJobResponse, err error) string {
					if err != nil {
						return fmt.Sprintf("%s: job create err: %s", time.Now().String(), err.Error())
					}
					return fmt.Sprintf("%s: job created: %s", time.Now().String(), r.Name)
				},
			}
		case "failed":
			runableJobsInfo.failed[jb.Name] = jobInfo[internal.RetryJobResponse]{
				Run: func() (internal.RetryJobResponse, error) {
					return internal.RetryJob(m.serviceInfo.Host, m.serviceInfo.ProjectName, jb.ID, m.credentials)
				},

				Button: fmt.Sprintf("alt+%d", altKey),
				MsgLog: func(r internal.RetryJobResponse, err error) string {
					if err != nil {
						return fmt.Sprintf("%s: failed job retry err: %s", time.Now().String(), err.Error())
					}
					return fmt.Sprintf("%s: failed job retry: %s", time.Now().String(), r.Name)
				},
			}
		case "success":
			runableJobsInfo.success[jb.Name] = jobInfo[internal.RetryJobResponse]{
				Run: func() (internal.RetryJobResponse, error) {
					return internal.RetryJob(m.serviceInfo.Host, m.serviceInfo.ProjectName, jb.ID, m.credentials)
				},

				Button: fmt.Sprintf("alt+%d", altKey),
				MsgLog: func(r internal.RetryJobResponse, err error) string {
					if err != nil {
						return fmt.Sprintf("%s: success job retry err: %s", time.Now().String(), err.Error())
					}
					return fmt.Sprintf("%s: success job retry: %s", time.Now().String(), r.Name)
				},
			}
		default:
			continue
		}

		altKey++
	}

	return runableJobsInfo, nil
}

type runnableJobs struct {
	manuals jobsByName[internal.RunJobResponse]
	failed  jobsByName[internal.RetryJobResponse]
	success jobsByName[internal.RetryJobResponse]
}

func (r runnableJobs) GetFormatButton(jobStatus, jobName string) string {

	switch jobStatus {
	case "manual":
		return internal.HelpStyle.Render(fmt.Sprintf("%s: run job", r.manuals[jobName].Button))
	case "failed":
		return internal.HelpStyle.Render(fmt.Sprintf("%s: retry job", r.failed[jobName].Button))
	case "success":
		return internal.HelpStyle.Render(fmt.Sprintf("%s: run again job", r.success[jobName].Button))
	}

	return ""
}

type jobsByName[T any] map[string]jobInfo[T]
