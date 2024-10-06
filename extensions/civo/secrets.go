/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package civo

import (
	"fmt"

	providerConfig "github.com/mgp-systems/internal-fabrik8-api/pkg/providerConfigs"
	pkgtypes "github.com/mgp-systems/internal-fabrik8-api/pkg/types"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
)

func BootstrapCivoMgmtCluster(clientset kubernetes.Interface, cl *pkgtypes.Cluster, destinationGitopsRepoURL string) error {
	opts := providerConfig.BootstrapOptions{
		GitUser:                  cl.GitAuth.User,
		DestinationGitopsRepoURL: destinationGitopsRepoURL,
		GitProtocol:              cl.GitProtocol,
		CloudflareAPIToken:       cl.CloudflareAuth.APIToken,
		CloudAuth:                cl.CivoAuth.Token,
		DNSProvider:              cl.DNSProvider,
		CloudProvider:            cl.CloudProvider,
		HTTPSPassword:            cl.GitAuth.Token,
		SSHToken:                 cl.GitAuth.PrivateKey,
	}

	if err := providerConfig.BootstrapMgmtCluster(clientset, opts); err != nil {
		log.Error().Msgf("unable to bootstrap management cluster: %s", err)
		return fmt.Errorf("unable to bootstrap management cluster: %w", err)
	}

	// Create secrets
	if err := providerConfig.BootstrapSecrets(clientset, cl); err != nil {
		log.Error().Msgf("unable to bootstrap secrets: %s", err)
		return fmt.Errorf("unable to bootstrap secrets: %w", err)
	}

	return nil
}
