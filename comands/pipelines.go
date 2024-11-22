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

func NewPipelinesCommand() *cli.Command {
	commandFlags := []cli.Flag{
		&cli.DurationFlag{
			Name:  "update-interval",
			Value: time.Second * 30,
			Usage: "pipeline status update interval",
		},
	}

	var credentials internal.Credentials
	var gitlabInfo internal.GitlabInfo

	commandFlags = append(commandFlags, credentials.CliFlags()...)
	commandFlags = append(commandFlags, gitlabInfo.CliFlags()...)

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

			err = gitlabInfo.Parse(cCtx)
			if err != nil {
				return err
			}

			updateRequestDuration := cCtx.Duration("update-interval")

			pipeline, err := internal.GetLatestPipeline(gitlabInfo.GitlabHost, gitlabInfo.ProjectName, gitlabInfo.RefName, credentials.Token())
			if err != nil {
				return err
			}

			if _, err := tea.NewProgram(pipelinesApp{
				headers:      []string{},
				rows:         [][]string{},
				pipelineInfo: pipeline,
				gitlabInfo:   gitlabInfo,
				credentials:  credentials,
				buf:          new(strings.Builder),

				timerUpdate: updateRequestDuration,
				timer:       timer.NewWithInterval(time.Millisecond*100, time.Millisecond),
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
	gitlabInfo   internal.GitlabInfo
	credentials  internal.Credentials

	timerUpdate time.Duration
	timer       timer.Model
	buf         *strings.Builder
}

func (m pipelinesApp) Init() tea.Cmd {
	return tea.Batch(m.timer.Init(), tea.SetWindowTitle("gitlab-tools pipelines"))
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
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func requestJobs(m pipelinesApp) (pipelinesApp, error) {
	jobs, err := internal.GetLatestJobs(m.gitlabInfo.GitlabHost, m.gitlabInfo.ProjectName, m.pipelineInfo.Id, m.credentials.Token())
	if err != nil {
		return pipelinesApp{}, err
	}

	triggeredJobs, err := internal.GetLatestTriggeredJobs(m.gitlabInfo.GitlabHost, m.gitlabInfo.ProjectName, m.pipelineInfo.Id, m.credentials.Token())
	if err != nil {
		return pipelinesApp{}, err
	}

	jobs = append(jobs, triggeredJobs...)

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
			)
		}
	}

	return m, nil
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

	m.buf.WriteString(m.gitlabInfo.RefName + " ")
	m.buf.WriteString(internal.EncodeHyperlink(m.pipelineInfo.WebUrl, "Open pipeline"))
	m.buf.WriteString("\n")

	m.buf.WriteString(t.Render())
	m.buf.WriteString("\n")

	m.buf.WriteString("next update in: ")
	m.buf.WriteString(m.timer.View())
	m.buf.WriteString("\n")

	m.buf.WriteString(internal.HelpStyle.Render("q: exit\n"))

	return m.buf.String()
}

type stageInfo struct {
	name  string
	order int
}
