package main

import (
	"fmt"
	"net/http"

	"github.com/google/go-github/v35/github"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var publishCmd = &cli.Command{
	Name:      "publish",
	Usage:     "publish a release",
	UsageText: "github-release [global-options] publish [options]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "github-sha",
			Aliases: []string{"s"},
			Usage:   "Git SHA of repo to create release on",
			EnvVars: []string{"GITHUB_SHA"},
		},
	},
	Action: actionHandler(publishAction),
}

func publishAction(c *cli.Context, opts *globalOptions) error {
	gh := opts.gh
	repo := opts.repo
	plan := opts.plan

	releaseName := plan.ReleaseName()
	githubSHA := c.String("github-sha")

	release, resp, err := gh.Repositories.GetReleaseByTag(
		c.Context, repo.Owner, repo.Name, releaseName,
	)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			prerelease := true
			release, _, err = gh.Repositories.CreateRelease(
				c.Context, repo.Owner, repo.Name, &github.RepositoryRelease{
					Name:            &releaseName,
					TagName:         &releaseName,
					TargetCommitish: &githubSHA,
					Prerelease:      &prerelease,
				},
			)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	relYAML, _ := yaml.Marshal(release)
	fmt.Printf("release: %+v\n", string(relYAML))

	return nil
}
