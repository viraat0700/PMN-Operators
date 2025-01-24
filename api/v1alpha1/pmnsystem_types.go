/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EnvironmentVariables struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
type EnvironmentVariablesNMSMagmaLte struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
type EnvironmentVariablesOrc8rNotifier struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
type EnvironmentVariablesOrc8rNginx struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
type EnvironmentVariablesDirectoryD struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
type EnvironmentVariablesPrometheusKafka struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
type EnvironmentVariablesNMSMagmaLteDevEnv struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Image struct {
	Repository string `json:"repository,omitempty"`
	Tag        string `json:"tag,omitempty"`
}

type Persistent struct {
	PvcClaimName     string `json:"pvcClaimName,omitempty"`
	StorageClassName string `json:"storageClassName,omitempty"`
}

type PersistenForStatefulSet struct {
	PvcClaimName     string `json:"pvcClaimName,omitempty"`
	StorageClassName string `json:"storageClassName,omitempty"`
}

type SecretConfig struct {
	SecretName    string   `json:"secretName,omitempty"`
	RequiredFiles []string `json:"requiredFiles,omitempty"`
}

type AlertmanagerConfigurer struct {
	Replica int32 `json:"replica,omitempty"`
	// NodeSelector                string                      `json:"nodeSelector,omitempty"`
	// Toleration                  string                      `json:"toleration,omitempty"`
	Affinity                    string                      `json:"affinity,omitempty"`
	ImageAlertmanagerConfigurer ImageAlertmanagerConfigurer `json:"imageAlertManagerConfigurer,omitempty"`
	AlertManagerConfigPort      int32                       `json:"alertManagerConfigPort,omitempty"`
	AlertManagerConfPath        string                      `json:"alertManagerConfPath,omitempty"`
	AlertmanagerURL             string                      `json:"alertmanagerURL,omitempty"`
	ServiceSpec                 ServiceSpec                 `json:"serviceSpec,omitempty"`
}

type PortSpec struct {
	Name       string `json:"name,omitempty"`
	Port       int32  `json:"port,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
	TargetPort int32  `json:"targetPort,omitempty"`
}

type ServiceSpec struct {
	Type     string     `json:"type,omitempty"`
	PortSpec []PortSpec `json:"portSpec,omitempty"`
}
type ImageAlertmanagerConfigurer struct {
	Repository      string `json:"repository,omitempty"`
	Tag             string `json:"tag,omitempty"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
}

type Metrics struct {
	Volumes Volumes `json:"volumes,omitempty"`
}
type Volumes struct {
	PrometheusConfig PrometheusConfig `json:"prometheusConfig,omitempty"`
}

type PrometheusConfig struct {
	VolumeSpec string `json:"volumeSpec,omitempty"`
}
type Prometheus struct {
	Replicas int32 `json:"replicas,omitempty"`
}

// ================alertManager====================
type ImageAlertmanager struct {
	Repository      string `json:"repository,omitempty"`
	Tag             string `json:"tag,omitempty"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
}
type AlertManager struct {
	// NodeSelector      string            `json:"nodeSelector,omitempty"`
	// Toleration        string            `json:"toleration,omitempty"`
	ImageAlertmanager ImageAlertmanager `json:"imageAlertmanager,omitempty"`
	ServiceSpec       ServiceSpec       `json:"serviceSpec,omitempty"`
}

// ===================User Grafana ====================
type ImageUserGrafana struct {
	Repository      string `json:"repository,omitempty"`
	Tag             string `json:"tag,omitempty"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
}
type UserGrafana struct {
	Replica int32 `json:"replicas,omitempty"`
	// NodeSelector       string               `json:"nodeSelector,omitempty"`
	// Toleration         string               `json:"toleration,omitempty"`
	ImageUserGrafana   ImageUserGrafana     `json:"imageUserGrafana,omitempty"`
	VolumesUserGrafana []VolumesUserGrafana `json:"volumesUserGrafana,omitempty"`
	ServiceSpec        ServiceSpec          `json:"serviceSpec,omitempty"`
}
type VolumesUserGrafana struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

// ====================orc8r-prometheus-nginx-prox=========================
type ServiceOrc8rSpec struct {
	Type          string          `json:"type,omitempty"`
	PortOrc8rSpec []PortOrc8rSpec `json:"portOrc8rSpec,omitempty"`
}
type PortOrc8rSpec struct {
	Name       string `json:"name,omitempty"`
	Port       int32  `json:"port,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
	TargetPort int32  `json:"targetPort,omitempty"`
	NodePort   int32  `json:"nodePort,omitempty"`
}
type PrometheusNginxProxy struct {
	Nginx Nginx `json:"nginx,omitempty"`
}
type Nginx struct {
	ServiceOrc8rSpec          ServiceOrc8rSpec          `json:"serviceOrc8rSpec,omitempty"`
	Replica                   int32                     `json:"replicas,omitempty"`
	ImagePrometheusNginxProxy ImagePrometheusNginxProxy `json:"imagePrometheusNginxProxy,omitempty"`
	VolumeMountPath           VolumeMountPath           `json:"volumeMountPath,omitempty"`
	SecretName                string                    `json:"secretName,omitempty"`
}
type ImagePrometheusNginxProxy struct {
	Repository      string `json:"repository,omitempty"`
	Tag             string `json:"tag,omitempty"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
}
type VolumeMountPath struct {
	MountPath []string `json:"mountPath,omitempty"`
	SubPath   []string `json:"subPath,omitempty"`
}

// ========================== Prometheus Kafka Adapter===============
type PrometheusKafkaAdapter struct {
	Replicas                              int32                                 `json:"replicas,omitempty"`
	ImagePrometheusKafkaAdapter           ImagePrometheusKafkaAdapter           `json:"imagePrometheusKafkaAdapter,omitempty"`
	VolumeMountPathPrometheusKafkaAdapter VolumeMountPathPrometheusKafkaAdapter `json:"volumeMountPathPrometheusKafkaAdapter,omitempty"`
	// Toleration                            string                                `json:"toleration,omitempty"`
	// NodeSelector                          string                                `json:"nodeSelector,omitempty"`
	ServiceSpecPrometheusKafkaAdapter ServiceSpecPrometheusKafkaAdapter `json:"serviceSpecPrometheusKafkaAdapter,omitempty"`
}
type ImagePrometheusKafkaAdapter struct {
	Repository      string `json:"repository,omitempty"`
	Tag             string `json:"tag,omitempty"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
}
type ServiceSpecPrometheusKafkaAdapter struct {
	Type                           string                           `json:"type,omitempty"`
	PortSpecPrometheusKafkaAdapter []PortSpecPrometheusKafkaAdapter `json:"portSpecPrometheusKafkaAdapter,omitempty"`
}
type PortSpecPrometheusKafkaAdapter struct {
	Name       string `json:"name,omitempty"`
	Port       int32  `json:"port,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
	TargetPort string `json:"targetPort,omitempty"`
}

type VolumeMountPathPrometheusKafkaAdapter struct {
	MountPath  []string `json:"mountPath,omitempty"`
	SecretName string   `json:"secretName,omitempty"`
}

// ==================== Prometheus Configurer ====================
type PrometheusConfigurer struct {
	Replicas int32 `json:"replicas,omitempty"`
	// Toleration                string                    `json:"toleration,omitempty"`
	// NodeSelector              string                    `json:"nodeSelector,omitempty"`
	ImagePrometheusConfigurer ImagePrometheusConfigurer `json:"imagePrometheusConfigurer,omitempty"`
	Volume                    Volume                    `json:"volume,omitempty"`
	Args                      []string                  `json:"args,omitempty"`
	ServiceSpec               ServiceSpec               `json:"serviceSpec,omitempty"`
}

type ImagePrometheusConfigurer struct {
	Repository      string `json:"repository,omitempty"`
	Tag             string `json:"tag,omitempty"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
}

type Volume struct {
	VolumeClaimName string `json:"volumeClaimName,omitempty"`
	VolumeMountPath string `json:"volumeMountPath,omitempty"`
}

// =================== Prometheus Cache ====================
type PrometheusCache struct {
	Replicas int32 `json:"replicas,omitempty"`
	// Toleration           string               `json:"toleration,omitempty"`
	// NodeSelector         string               `json:"nodeSelector,omitempty"`
	ImagePrometheusCache ImagePrometheusCache `json:"imagePrometheusCache,omitempty"`
	Args                 []string             `json:"args,omitempty"`
	ServiceSpec          ServiceSpec          `json:"serviceSpec,omitempty"`
}

type ImagePrometheusCache struct {
	Repository      string `json:"repository,omitempty"`
	Tag             string `json:"tag,omitempty"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
}

// ====================NMS-MagmaLte ========================
type NmsMagmaLte struct {
	Replicas int32 `json:"replicas,omitempty"`
	// Toleration             string                 `json:"toleration,omitempty"`
	// NodeSelector           string                 `json:"nodeSelector,omitempty"`
	ImageMagmaLte          ImageMagmaLte          `json:"imageMagmaLte,omitempty"`
	VolumeMountNmsMagmaLte VolumeMountNmsMagmaLte `json:"volumeMountNmsMagmaLte,omitempty"`
	VolumesNmsMagmaLte     VolumesNmsMagmaLte     `json:"volumesNmsMagmaLte,omitempty"`
	ServiceSpec            ServiceSpec            `json:"serviceSpec,omitempty"`
}
type ImageMagmaLte struct {
	Repository      string `json:"repository,omitempty"`
	Tag             string `json:"tag,omitempty"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
}
type VolumeMountNmsMagmaLte struct {
	VolumeMountPath []string `json:"volumeMountPath,omitempty"`
	VolumeSubPath   []string `json:"volumeSubPath,omitempty"`
}
type VolumesNmsMagmaLte struct {
	Secretname []string `json:"secretName,omitempty"`
}

// =======================Orc8r Notifier Deployment ======================
type Orc8rNotifier struct {
	ImageOrc8rNotifier ImageOrc8rNotifier `json:"imageOrc8rNotifier,omitempty"`
	Args               []string           `json:"args,omitempty"`
	LivenessProbePort  int32              `json:"livenessProbe,omitempty"`
	ServiceSpec        ServiceSpec        `json:"serviceSpec,omitempty"`
	ReadinessProbePort int32              `json:"readinessProbe,omitempty"`
	PortDeployment     int32              `json:"portDeployment,omitempty"`
}

type ImageOrc8rNotifier struct {
	Repository      string `json:"repository,omitempty"`
	Tag             string `json:"tag,omitempty"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
}

// ============orc8r Nginx Deployment ======================
type Orc8rNginxDeployment struct {
	Replicas        int32           `json:"replicas,omitempty"`
	ImageOrc8rNginx ImageOrc8rNginx `json:"imageOrc8rNginx,omitempty"`
	// Toleration                 string                     `json:"toleration,omitempty"`
	// NodeSelector               string                     `json:"nodeSelector,omitempty"`
	VolumesOrc8rNginx          VolumesOrc8rNginx          `json:"volumesOrc8rNginx,omitempty"`
	PortOrc8rNginx             PortOrc8rNginx             `json:"portOrc8rNginx,omitempty"`
	VolumesMountPathOrc8rNginx VolumesMountPathOrc8rNginx `json:"volumesMountPathOrc8rNginx,omitempty"`
}
type ImageOrc8rNginx struct {
	Repository      string `json:"repository,omitempty"`
	Tag             string `json:"tag,omitempty"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
}
type VolumesOrc8rNginx struct {
	SecretName []string `json:"secretName,omitempty"`
}
type VolumesMountPathOrc8rNginx struct {
	MountPath []string `json:"mountPath,omitempty"`
}

type PortOrc8rNginx struct {
	Port []int32 `json:"port,omitempty"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PmnsystemSpec defines the desired state of Pmnsystem
type PmnsystemSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Pmnsystem. Edit pmnsystem_types.go to remove/update
	ReplicaCount                  int32                                   `json:"replicaCount,omitempty"`
	NginxImage                    string                                  `json:"nginxImage,omitempty"`
	NotifierImage                 string                                  `json:"notifierImage,omitempty"`
	PullPolicy                    string                                  `json:"pullPolicy,omitempty"`
	Persistent                    Persistent                              `json:"persistent,omitempty"`
	PersistentForStatefulSet      PersistenForStatefulSet                 `json:"persistentForStatefulSet,omitempty"`
	NameSpace                     string                                  `json:"nameSpace,omitempty"`
	Image                         Image                                   `json:"image,omitempty"`
	EnvVariables                  []EnvironmentVariables                  `json:"envVariables,omitempty"`
	EnvVariablesDirectoryD        []EnvironmentVariablesDirectoryD        `json:"envVariablesDirectoryD,omitempty"`
	EnvVariablesOrc8rNginx        []EnvironmentVariablesOrc8rNginx        `json:"envVariablesOrc8rNginx,omitempty"`
	EnvVariablesOrc8rNotifier     []EnvironmentVariablesOrc8rNotifier     `json:"envVariablesOrc8rNotifier,omitempty"`
	EnvVariablesNMSMagmaLte       []EnvironmentVariablesNMSMagmaLte       `json:"envVariablesNMSMagmaLte,omitempty"`
	EnvVariablesNMSMagmaLteDevEnv []EnvironmentVariablesNMSMagmaLteDevEnv `json:"envVariablesNMSMagmaLteDevEnv,omitempty"`
	EnvVariablesPrometheusKafka   []EnvironmentVariablesPrometheusKafka   `json:"envVariablesPrometheusKafka,omitempty"`
	ImagePullSecrets              string                                  `json:"imagePullSecrets,omitempty"`
	DevEnvironment                bool                                    `json:"devEnvironment,omitempty"`
	CloudEnvironment              bool                                    `json:"cloudEnvironment,omitempty"`
	CertDir                       string                                  `json:"certDir,omitempty"`
	RepoPath                      string                                  `json:"repoPath,omitempty"`
	Secrets                       []SecretConfig                          `json:"secrets,omitempty"`
	AlertmanagerConfigurer        AlertmanagerConfigurer                  `json:"alertmanagerConfigurer,omitempty"`
	Metrics                       Metrics                                 `json:"metrics,omitempty"`
	Prometheus                    Prometheus                              `json:"prometheus,omitempty"`
	AlertManager                  AlertManager                            `json:"alertManager,omitempty"`
	UserGrafana                   UserGrafana                             `json:"userGrafana,omitempty"`
	PrometheusNginxProxy          PrometheusNginxProxy                    `json:"prometheusNginxProxy,omitempty"`
	PrometheusKafkaAdapter        PrometheusKafkaAdapter                  `json:"prometheusKafkaAdapter,omitempty"`
	PrometheusConfigurer          PrometheusConfigurer                    `json:"prometheusConfigurer,omitempty"`
	PrometheusCache               PrometheusCache                         `json:"prometheusCache,omitempty"`
	NmsMagmaLte                   NmsMagmaLte                             `json:"nmsMagmaLte,omitempty"`
	Orc8rNotifier                 Orc8rNotifier                           `json:"orc8rNotifier,omitempty"`
	Orc8rNginxDeployment          Orc8rNginxDeployment                    `json:"orc8rNginxDeployment,omitempty"`
}

// PmnsystemStatus defines the observed state of Pmnsystem
type PmnsystemStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Pmnsystem is the Schema for the pmnsystems API
type Pmnsystem struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PmnsystemSpec   `json:"spec,omitempty"`
	Status PmnsystemStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PmnsystemList contains a list of Pmnsystem
type PmnsystemList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pmnsystem `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pmnsystem{}, &PmnsystemList{})
}
