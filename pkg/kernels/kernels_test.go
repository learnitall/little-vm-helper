package kernels

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDir(t *testing.T) {
	configs := []*Conf{
		nil,
		{
			Kernels: []KernelConf{
				{
					Name: "bpf-next",
					URL:  "git://git.kernel.org/pub/scm/linux/kernel/git/bpf/bpf-next.git",
					Conf: []ConfigOption{
						{"--enable", "CONFIG_DEBUG_INFO"},
						{"--disable", "CONFIG_DEBUG_KERNEL"},
						{"--enable CONFIG_BPF"},
						{"--enable CONFIG_BPF_SYSCALL"},
					},
				}, {
					Name: "5.18.8",
					URL:  "git://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git#v5.18.8",
					Conf: []ConfigOption{
						{"--enable", "CONFIG_DEBUG_INFO"},
						{"--disable", "CONFIG_DEBUG_KERNEL"},
						{"--enable CONFIG_BPF"},
						{"--enable CONFIG_BPF_SYSCALL"},
					},
				},
			},
		},
	}

	for _, conf := range configs {
		// NB: anonymous function so that os.RemoveAll() is called in all iterations
		func() {
			dir, err := ioutil.TempDir("", "test_kernel")
			assert.Nil(t, err)
			defer os.RemoveAll(dir)
			err = InitDir(dir, conf)
			assert.Nil(t, err)

			if conf == nil {
				conf = &Conf{
					Kernels: make([]KernelConf, 0),
				}
			}

			var xconf Conf
			data, err := os.ReadFile(path.Join(dir, ConfigFname))
			assert.Nil(t, err)
			err = json.Unmarshal(data, &xconf)
			assert.Nil(t, err)
			assert.Equal(t, &xconf, conf)

		}()
	}
}