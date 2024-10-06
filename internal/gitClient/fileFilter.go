/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package gitClient //nolint:revive,stylecheck // allowed temporarily during code reorg

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	pkg "github.com/mgp-systems/internal-fabrik8-api/internal"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// AppendFile verify if a file must be appended to committed gitops
// meant to help exclude undesired state files to be pushed to gitops
func AppendFile(cloudType string, reponame string, filename string) bool {
	// result := true
	// TODO: make this to be loaced by Arrays of exclusion rules
	// TODO: Make this a bit more fancier
	// Once we have some critical mass of rules, this will be improved
	if cloudType == pkg.CloudAws {
		if strings.Contains(reponame, "gitops") {
			if filename == "terraform/base/kubeconfig" {
				// https://github.com/konstructio/kubefirst/issues/926
				log.Debug().Msgf("file not included on commit[#926]: '%s'", filename)
				return false
			}
		}
	}
	if cloudType == pkg.CloudK3d {
		if strings.Contains(reponame, "gitops") {
			if strings.HasPrefix(filename, "argo-workflows") {
				// https://github.com/konstructio/kubefirst/issues/959
				log.Debug().Msgf("file not included on commit[#959]: '%s'", filename)
				return false
			}
		}
	}

	return true
}

// GitAddWithFilter Check workdir for files to commit
// filter out the undersired ones based on context
func GitAddWithFilter(w *git.Worktree) error {
	status, err := w.Status()
	if err != nil {
		log.Debug().Msgf("error getting worktree status: %s", err)
		return fmt.Errorf("error getting worktree status: %w", err)
	}

	for file, s := range status {
		log.Printf("the file is %s the status is %v", file, s.Worktree)
		if AppendFile(viper.GetString("cloud"), "gitops", file) {
			_, err = w.Add(file)
			if err != nil {
				log.Error().Err(err).Msgf("error getting worktree status: %s", err)
				return fmt.Errorf("error adding file %s: %w", file, err)
			}
		}
	}

	return nil
}
