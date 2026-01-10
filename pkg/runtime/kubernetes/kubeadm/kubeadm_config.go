// Copyright © 2021 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubeadm

import (
	"fmt"
	"io"
	"net"
	"strings"

	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta3"

	versionUtils "github.com/sealerio/sealer/utils/version"
	"github.com/sirupsen/logrus"

	"github.com/sealerio/sealer/utils"
	strUtils "github.com/sealerio/sealer/utils/strings"

	"github.com/imdario/mergo"
	"k8s.io/kube-proxy/config/v1alpha1"
	"k8s.io/kubelet/config/v1beta1"
)

// Read config from https://github.com/sealerio/sealer/blob/main/docs/design/clusterfile-v2.md and overwrite default kubeadm.yaml
// Use github.com/imdario/mergo to merge kubeadm config in Clusterfile and the default kubeadm config
// Using a config filter to handle some edge cases

// https://github.com/kubernetes/kubernetes/blob/master/cmd/kubeadm/app/apis/kubeadm/v1beta3/types.go
// Using map to overwrite Kubeadm configs

// nolint
type KubeadmConfig struct {
	v1beta3.InitConfiguration
	v1beta3.ClusterConfiguration
	v1alpha1.KubeProxyConfiguration
	v1beta1.KubeletConfiguration
	v1beta3.JoinConfiguration
}

const (
	EtcdServers = "etcd-servers"
)

const (
	V1991 = "v1.19.1"
	V1992 = "v1.19.2"
	V1150 = "v1.15.0"
	V1200 = "v1.20.0"
	V1230 = "v1.23.0"
	V1330 = "v1.33.0"
	V1310 = "v1.31.0"

	// kubeadm api version
	KubeadmV1beta1 = "kubeadm.k8s.io/v1beta1"
	KubeadmV1beta2 = "kubeadm.k8s.io/v1beta2"
	KubeadmV1beta3 = "kubeadm.k8s.io/v1beta3"
	KubeadmV1beta4 = "kubeadm.k8s.io/v1beta4"
)

// LoadFromClusterfile :Load KubeadmConfig from Clusterfile.
// If it has `KubeadmConfig` in Clusterfile, load every field to each configuration.
// If Kubeadm raw config in Clusterfile, just load it.
func (k *KubeadmConfig) LoadFromClusterfile(kubeadmConfig KubeadmConfig) error {
	k.APIServer.CertSANs = strUtils.RemoveDuplicate(append(k.APIServer.CertSANs, kubeadmConfig.APIServer.CertSANs...))

	return mergo.Merge(k, kubeadmConfig)
}

// Merge Using github.com/imdario/mergo to merge KubeadmConfig to the ClusterImage default kubeadm config, overwrite some field.
// if defaultKubeadmConfig file not exist, use default raw kubeadm config to merge k.KubeConfigSpec empty value
func (k *KubeadmConfig) Merge(kubeadmYamlPath string, decode func(arg string, kind string) (interface{}, error)) error {
	newConfig, err := LoadKubeadmConfigs(kubeadmYamlPath, decode)
	if err != nil {
		return fmt.Errorf("failed to found kubeadm config from %s: %v", kubeadmYamlPath, err)
	}
	k.APIServer.CertSANs = strUtils.RemoveDuplicate(append(k.APIServer.CertSANs, newConfig.APIServer.CertSANs...))

	return mergo.Merge(k, newConfig)
}

func (k *KubeadmConfig) setAPIVersion(apiVersion string) {
	k.InitConfiguration.APIVersion = apiVersion
	k.ClusterConfiguration.APIVersion = apiVersion
	k.JoinConfiguration.APIVersion = apiVersion
}

func (k *KubeadmConfig) setKubeadmAPIVersion() {
	kv := versionUtils.Version(k.KubernetesVersion)
	greaterThanKV1150, err := kv.GreaterThan(V1150)
	if err != nil {
		logrus.Errorf("compare kubernetes version failed: %s", err)
	}
	greaterThanKV1230, err := kv.GreaterThan(V1230)
	if err != nil {
		logrus.Errorf("compare kubernetes version failed: %s", err)
	}
	greaterThanKV1310, err := kv.GreaterThan(V1310)
	if err != nil {
		logrus.Errorf("compare kubernetes version failed: %s", err)
	}
	switch {
	case greaterThanKV1310:
		k.setAPIVersion(KubeadmV1beta4)
	case greaterThanKV1230:
		k.setAPIVersion(KubeadmV1beta3)
	case greaterThanKV1150:
		k.setAPIVersion(KubeadmV1beta2)
	default:
		// Compatible with versions 1.14 and 1.13. but do not recommend.
		k.setAPIVersion(KubeadmV1beta1)
	}
}

func (k *KubeadmConfig) GetCertSANS() []string {
	return k.ClusterConfiguration.APIServer.CertSANs
}

func (k *KubeadmConfig) GetDNSDomain() string {
	return k.ClusterConfiguration.Networking.DNSDomain
}

func (k *KubeadmConfig) GetSvcCIDR() string {
	return k.ClusterConfiguration.Networking.ServiceSubnet
}

func LoadKubeadmConfigs(arg string, decode func(arg string, kind string) (interface{}, error)) (KubeadmConfig, error) {
	kubeadmConfig := KubeadmConfig{}
	initConfig, err := decode(arg, InitConfiguration)
	if err != nil && err != io.EOF {
		return kubeadmConfig, err
	} else if initConfig != nil {
		kubeadmConfig.InitConfiguration = *initConfig.(*v1beta3.InitConfiguration)
	}
	clusterConfig, err := decode(arg, ClusterConfiguration)
	if err != nil && err != io.EOF {
		return kubeadmConfig, err
	} else if clusterConfig != nil {
		kubeadmConfig.ClusterConfiguration = *clusterConfig.(*v1beta3.ClusterConfiguration)
	}
	kubeProxyConfig, err := decode(arg, KubeProxyConfiguration)
	if err != nil && err != io.EOF {
		return kubeadmConfig, err
	} else if kubeProxyConfig != nil {
		kubeadmConfig.KubeProxyConfiguration = *kubeProxyConfig.(*v1alpha1.KubeProxyConfiguration)
	}
	kubeletConfig, err := decode(arg, KubeletConfiguration)
	if err != nil && err != io.EOF {
		return kubeadmConfig, err
	} else if kubeletConfig != nil {
		kubeadmConfig.KubeletConfiguration = *kubeletConfig.(*v1beta1.KubeletConfiguration)
	}
	joinConfig, err := decode(arg, JoinConfiguration)
	if err != nil && err != io.EOF {
		return kubeadmConfig, err
	} else if joinConfig != nil {
		kubeadmConfig.JoinConfiguration = *joinConfig.(*v1beta3.JoinConfiguration)
	}
	return kubeadmConfig, nil
}

func getEtcdEndpointsWithHTTPSPrefix(masters []net.IP) string {
	var tmpSlice []string
	for _, ip := range masters {
		tmpSlice = append(tmpSlice, fmt.Sprintf("https://%s", net.JoinHostPort(ip.String(), "2379")))
	}

	return strings.Join(tmpSlice, ",")
}

func NewKubeadmConfig(fromClusterFile KubeadmConfig, fromFile string, masters []net.IP, apiServerDomain,
	cgroupDriver string, imageRepo string, apiServerVIP net.IP, extraSANs []string) (KubeadmConfig, error) {
	conf := KubeadmConfig{}

	if err := conf.LoadFromClusterfile(fromClusterFile); err != nil {
		return conf, fmt.Errorf("failed to load kubeadm config from clusterfile: %v", err)
	}
	// TODO handle the kubeadm config, like kubeproxy config
	//The configuration set here does not require merge

	conf.InitConfiguration.LocalAPIEndpoint.AdvertiseAddress = masters[0].String()
	conf.ControlPlaneEndpoint = net.JoinHostPort(apiServerDomain, "6443")

	if conf.APIServer.ExtraArgs == nil {
		conf.APIServer.ExtraArgs = make(map[string]string)
	}
	conf.APIServer.ExtraArgs[EtcdServers] = getEtcdEndpointsWithHTTPSPrefix(masters)
	conf.IPVS.ExcludeCIDRs = append(conf.KubeProxyConfiguration.IPVS.ExcludeCIDRs, fmt.Sprintf("%s/32", apiServerVIP))
	conf.KubeletConfiguration.CgroupDriver = cgroupDriver
	conf.ClusterConfiguration.APIServer.CertSANs = []string{"127.0.0.1", apiServerDomain, apiServerVIP.String()}
	conf.ClusterConfiguration.APIServer.CertSANs = append(conf.ClusterConfiguration.APIServer.CertSANs, extraSANs...)
	for _, m := range masters {
		conf.ClusterConfiguration.APIServer.CertSANs = append(conf.ClusterConfiguration.APIServer.CertSANs, m.String())
	}

	if err := conf.Merge(fromFile, utils.DecodeCRDFromFile); err != nil {
		return conf, err
	}

	if err := conf.Merge(DefaultKubeadmConfig, utils.DecodeCRDFromString); err != nil {
		return conf, err
	}

	conf.setKubeadmAPIVersion()

	if conf.ClusterConfiguration.Networking.DNSDomain == "" {
		conf.ClusterConfiguration.Networking.DNSDomain = "cluster.local"
	}
	if conf.JoinConfiguration.Discovery.BootstrapToken == nil {
		conf.JoinConfiguration.Discovery.BootstrapToken = &v1beta3.BootstrapTokenDiscovery{}
	}

	// set cluster image repo,kubeadm will pull container image from this registry.
	if conf.ClusterConfiguration.ImageRepository == "" {
		conf.ClusterConfiguration.ImageRepository = imageRepo
	}
	if conf.ClusterConfiguration.DNS.ImageMeta.ImageRepository == "" {
		conf.ClusterConfiguration.DNS.ImageMeta.ImageRepository = fmt.Sprintf("%s/%s", imageRepo, "coredns")
	}

	return conf, nil
}

type ArgV1beta4 struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

func mapToArgsV1beta4(m map[string]string) []ArgV1beta4 {
	var args []ArgV1beta4
	for k, v := range m {
		args = append(args, ArgV1beta4{Name: k, Value: v})
	}
	return args
}

// GetMarshableConfigs returns configs in a format that can be marshaled based on apiVersion
func (k *KubeadmConfig) GetMarshableConfigs() []interface{} {
	if k.InitConfiguration.APIVersion != KubeadmV1beta4 {
		return []interface{}{
			&k.InitConfiguration,
			&k.ClusterConfiguration,
			&k.KubeletConfiguration,
			&k.KubeProxyConfiguration,
			&k.JoinConfiguration,
		}
	}

	// Handle v1beta4
	initCfg := k.InitConfiguration
	clusterCfg := k.ClusterConfiguration
	joinCfg := k.JoinConfiguration

	// Convert maps to slices for v1beta4
	type InitCfgV1beta4 struct {
		v1beta3.InitConfiguration `json:",inline" yaml:",inline"`
		NodeRegistration          struct {
			v1beta3.NodeRegistrationOptions `json:",inline" yaml:",inline"`
			ExtraArgs                       []ArgV1beta4 `json:"extraArgs,omitempty" yaml:"extraArgs,omitempty"`
		} `json:"nodeRegistration,omitempty" yaml:"nodeRegistration,omitempty"`
	}

	type ClusterCfgV1beta4 struct {
		v1beta3.ClusterConfiguration `json:",inline" yaml:",inline"`
		APIServer                    struct {
			v1beta3.APIServer `json:",inline" yaml:",inline"`
			ExtraArgs         []ArgV1beta4 `json:"extraArgs,omitempty" yaml:"extraArgs,omitempty"`
		} `json:"apiServer,omitempty" yaml:"apiServer,omitempty"`
		ControllerManager struct {
			v1beta3.ControlPlaneComponent `json:",inline" yaml:",inline"`
			ExtraArgs                     []ArgV1beta4 `json:"extraArgs,omitempty" yaml:"extraArgs,omitempty"`
		} `json:"controllerManager,omitempty" yaml:"controllerManager,omitempty"`
		Scheduler struct {
			v1beta3.ControlPlaneComponent `json:",inline" yaml:",inline"`
			ExtraArgs                     []ArgV1beta4 `json:"extraArgs,omitempty" yaml:"extraArgs,omitempty"`
		} `json:"scheduler,omitempty" yaml:"scheduler,omitempty"`
		Etcd struct {
			v1beta3.Etcd `json:",inline" yaml:",inline"`
			Local        *struct {
				v1beta3.LocalEtcd `json:",inline" yaml:",inline"`
				ExtraArgs         []ArgV1beta4 `json:"extraArgs,omitempty" yaml:"extraArgs,omitempty"`
			} `json:"local,omitempty" yaml:"local,omitempty"`
		} `json:"etcd,omitempty" yaml:"etcd,omitempty"`
	}

	type JoinCfgV1beta4 struct {
		v1beta3.JoinConfiguration `json:",inline" yaml:",inline"`
		NodeRegistration          struct {
			v1beta3.NodeRegistrationOptions `json:",inline" yaml:",inline"`
			ExtraArgs                       []ArgV1beta4 `json:"extraArgs,omitempty" yaml:"extraArgs,omitempty"`
		} `json:"nodeRegistration,omitempty" yaml:"nodeRegistration,omitempty"`
	}

	var iV4 InitCfgV1beta4
	iV4.InitConfiguration = initCfg
	iV4.NodeRegistration.NodeRegistrationOptions = initCfg.NodeRegistration
	iV4.NodeRegistration.ExtraArgs = mapToArgsV1beta4(initCfg.NodeRegistration.KubeletExtraArgs)

	var cV4 ClusterCfgV1beta4
	cV4.ClusterConfiguration = clusterCfg
	cV4.APIServer.APIServer = clusterCfg.APIServer
	cV4.APIServer.ExtraArgs = mapToArgsV1beta4(clusterCfg.APIServer.ExtraArgs)
	cV4.ControllerManager.ControlPlaneComponent = clusterCfg.ControllerManager
	cV4.ControllerManager.ExtraArgs = mapToArgsV1beta4(clusterCfg.ControllerManager.ExtraArgs)
	cV4.Scheduler.ControlPlaneComponent = clusterCfg.Scheduler
	cV4.Scheduler.ExtraArgs = mapToArgsV1beta4(clusterCfg.Scheduler.ExtraArgs)
	cV4.Etcd.Etcd = clusterCfg.Etcd
	if clusterCfg.Etcd.Local != nil {
		cV4.Etcd.Local = &struct {
			v1beta3.LocalEtcd `json:",inline" yaml:",inline"`
			ExtraArgs         []ArgV1beta4 `json:"extraArgs,omitempty" yaml:"extraArgs,omitempty"`
		}{}
		cV4.Etcd.Local.LocalEtcd = *clusterCfg.Etcd.Local
		cV4.Etcd.Local.ExtraArgs = mapToArgsV1beta4(clusterCfg.Etcd.Local.ExtraArgs)
	}

	var jV4 JoinCfgV1beta4
	jV4.JoinConfiguration = joinCfg
	jV4.NodeRegistration.NodeRegistrationOptions = joinCfg.NodeRegistration
	jV4.NodeRegistration.ExtraArgs = mapToArgsV1beta4(joinCfg.NodeRegistration.KubeletExtraArgs)

	return []interface{}{
		&iV4,
		&cV4,
		&k.KubeletConfiguration,
		&k.KubeProxyConfiguration,
		&jV4,
	}
}
