package localgloo_e2e_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	controlplanebootstrap "github.com/solo-io/gloo/internal/control-plane/bootstrap"
	functiondiscoveryopts "github.com/solo-io/gloo/internal/function-discovery/options"
	upstreamdiscbootstrap "github.com/solo-io/gloo/internal/upstream-discovery/bootstrap"
	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/pkg/localgloo"
)

var _ = Describe("LocalglooE2e", func() {
	var (
		tmpDir string
		err    error
		xdsPort = 8081

		baseOpts              bootstrap.Options
		controlPlaneOpts      = controlplanebootstrap.Options{
			IngressOptions: controlplanebootstrap.IngressOptions{
				BindAddress: "::",
				Port: uint32(xdsPort),
				SecurePort: uint32(xdsPort+1),
			},
		}
		upstreamDiscoveryOpts upstreamdiscbootstrap.Options
		functionDiscoveryOpts = functiondiscoveryopts.DiscoveryOptions{
			//AutoDiscoverSwagger:   true,
			//AutoDiscoverNATS:      true,
			//AutoDiscoverFaaS:      true,
			//AutoDiscoverFission:   true,
			//AutoDiscoverProjectFn: true,
			//AutoDiscoverGRPC:      true,
		}
	)
	BeforeEach(func() {
		tmpDir, err = ioutil.TempDir("", "localgloo-test")
		Expect(err).To(BeNil())
		configDir := tmpDir
		filesDir := filepath.Join(tmpDir, "files")
		secretsDir := filepath.Join(tmpDir, "secrets")
		os.MkdirAll(configDir, 0755)
		os.MkdirAll(filesDir, 0755)
		os.MkdirAll(secretsDir, 0755)
		baseOpts.ConfigStorageOptions.Type = "file"
		baseOpts.FileStorageOptions.Type = "file"
		baseOpts.SecretStorageOptions.Type = "file"
		baseOpts.FileOptions.ConfigDir = configDir
		baseOpts.FileOptions.SecretDir = secretsDir
		baseOpts.FileOptions.FilesDir = filesDir
		baseOpts.ConfigStorageOptions.SyncFrequency = time.Second
		baseOpts.FileStorageOptions.SyncFrequency = time.Second
		baseOpts.SecretStorageOptions.SyncFrequency = time.Second
		controlPlaneOpts.Options = baseOpts
		upstreamDiscoveryOpts.Options = baseOpts
	})
	AfterEach(func() {
		os.RemoveAll(tmpDir)
	})
	It("runs everything", func() {
		stop := make(chan struct{})
		localgloo.Run(stop, xdsPort, baseOpts, controlPlaneOpts, upstreamDiscoveryOpts, functionDiscoveryOpts)

		Expect(err).To(BeNil())
	})
})