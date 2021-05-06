package main

import (
	"fmt"
	"os"

	"github.com/google/go-github/v35/github"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var app = &cli.App{
	Name:  "github-release",
	Usage: "manage GitHub releases",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "github-token",
			Aliases: []string{"t"},
			Usage:   "GitHub API Token (read from GITHUB_TOKEN if set)",
		},
		&cli.StringFlag{
			Name:    "repo",
			Aliases: []string{"r"},
			Usage:   "Owner/name of GitHub repo to publish releases to",
			EnvVars: []string{"GITHUB_REPOSITORY"},
			Value:   "jimeh/release-tests",
		},
		&cli.PathFlag{
			Name:      "plan",
			Aliases:   []string{"p"},
			Usage:     "Load plan from `FILE`",
			EnvVars:   []string{"BUILD_PLAN"},
			Required:  true,
			TakesFile: true,
		},
	},
	Commands: []*cli.Command{
		checkCmd,
		publishCmd,
	},
}

type globalOptions struct {
	gh   *github.Client
	repo *Repo
	plan *Plan
}

func actionHandler(
	f func(*cli.Context, *globalOptions) error,
) func(*cli.Context) error {
	return func(c *cli.Context) error {
		b, err := os.ReadFile(c.Path("plan"))
		if err != nil {
			return err
		}

		plan := &Plan{}
		err = yaml.Unmarshal(b, plan)
		if err != nil {
			return err
		}

		token := c.String("github-token")
		if t := os.Getenv("GITHUB_TOKEN"); t != "" {
			token = t
		}

		opts := &globalOptions{
			gh:   NewGitHubClient(c.Context, token),
			repo: NewRepo(c.String("repo")),
			plan: plan,
		}

		return f(c, opts)
	}
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}
