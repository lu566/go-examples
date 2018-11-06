package v1alpha1

import (
	"github.com/hidevopsio/hiboot/pkg/system"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	CODEPATH   string = "codepath"
	CLONE      string = "clone"
	COMPILE    string = "compile"
	BUILDIMAGE string = "buildImage"
	PUSHIMAGE  string = "pushImage"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Foo is a specification for a Foo resource
type Pipeline struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Specification of the desired behavior of the Deployment.
	// +optional
	Spec PipelineSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Most recently observed status of the Deployment.
	// +optional
	Status PipelineStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type PipelineSpec struct {
	App               string            `json:"app"  protobuf:"bytes,1,opt,name=app"`
	Profile           string            `json:"profile"  protobuf:"bytes,2,opt,name=profile"`
	Project           string            `json:"project"  protobuf:"bytes,3,opt,name=project"`
	Cluster           string            `json:"cluster"  protobuf:"bytes,4,opt,name=cluster"`
	Namespace         string            `json:"namespace"  protobuf:"bytes,5,opt,name=namespace"`
	Scm               Scm               `json:"scm"  protobuf:"bytes,6,opt,name=scm"`
	Version           string            `json:"version"  protobuf:"bytes,7,opt,name=version"`
	DockerRegistry    string            `json:"dockerRegistry"  protobuf:"bytes,8,opt,name=dockerRegistry"`
	NodeSelector      string            `json:"nodeSelector"  protobuf:"bytes,9,opt,name=nodeSelector"`
	Identifiers       []string          `json:"identifiers"  protobuf:"bytes,10,opt,name=identifiers"`
	ConfigFiles       []string          `json:"configFiles"  protobuf:"bytes,11,opt,name=configFiles"`
	Ports             []Ports           `json:"ports" protobuf:"bytes,12,opt,name=ports"`
	BuildConfigs      BuildConfigs      `json:"buildConfigs" protobuf:"bytes,13,opt,name=buildConfigs"`
	DeploymentConfigs DeploymentConfigs `json:"deploymentConfigs" protobuf:"bytes,14,opt,name=deploymentConfigs"`
	GatewayConfigs    GatewayConfigs    `json:"gatewayConfigs" protobuf:"bytes,15,opt,name=gatewayConfigs"`
	Events            []Events          `json:"events" protobuf:"bytes,16,opt,name=events"`
}

type Events struct {
	Name       string   `json:"name" protobuf:"bytes,1,opt,name=name"`
	EventTypes []string `json:"eventTypes" protobuf:"bytes,2,opt,name=eventTypes"`
}

type Ports struct {
	Name          string `json:"name" protobuf:"bytes,1,opt,name=name"`
	Port          int32  `json:"port"  protobuf:"bytes,2,opt,name=port"`
	ContainerPort int32  `json:"containerPort" protobuf:"bytes,3,opt,name=containerPort"`
	Protocol      string `json:"protocol,omitempty" protobuf:"bytes,4,opt,name=protocol"`
}

type Scm struct {
	Type     string `json:"type" protobuf:"bytes,1,opt,name=type"`
	Url      string `json:"url" protobuf:"bytes,2,opt,name=url"`
	Ref      string `json:"ref" protobuf:"bytes,3,opt,name=ref"`
	UserName string `json:"userName" protobuf:"bytes,4,opt,name=userName"`
	Password string `json:"password" protobuf:"bytes,4,opt,name=password"`
}

type DeploymentConfigs struct {
	HealthEndPoint string       `json:"healthEndPoint" protobuf:"bytes,1,opt,name=healthEndPoint"`
	Replicas       int32        `json:"replicas" protobuf:"bytes,2,opt,name=replicas"`
	Env            []system.Env `json:"env" protobuf:"bytes,3,opt,name=env"`
	Labels         Labels       `json:"labels" protobuf:"bytes,4,opt,name=labels"`
	Project        string       `json:"project" protobuf:"bytes,5,opt,name=project"`
}

type Labels struct {
	App     string `json:"app" protobuf:"bytes,1,opt,name=app"`
	Version string `json:"version" protobuf:"bytes,2,opt,name=version"`
	Cluster string `json:"cluster" protobuf:"bytes,3,opt,name=cluster"`
}

type BuildConfigs struct {
	TagFrom     string       `json:"tagFrom" protobuf:"bytes,1,opt,name=tagFrom"`
	ImageStream string       `json:"imageStream" protobuf:"bytes,2,opt,name=imageStream"`
	Env         []system.Env `json:"env" protobuf:"bytes,3,opt,name=env"`
	Project     string       `json:"project" protobuf:"bytes,4,opt,name=project"`
	BaseImage   string       `json:"baseImage" protobuf:"bytes,6,opt,name=baseImage"`
	Tags        []string     `json:"tags" protobuf:"bytes,7,opt,name=tags"`
	DockerFile  []string     `json:"dockerFile" protobuf:"bytes,8,opt,name=dockerFile"`
	DstDir      string       `json:"dstDir" protobuf:"bytes,9,opt,name=dstDir"`
	CloneType   string       `json:"cloneType" protobuf:"bytes,9,opt,name=cloneType"`
	CompileCmd  []CompileCmd `json:"compileCmd" protobuf:"bytes,10,opt,name=compileCmd"`
	DeployData  DeployData   `json:"deployData" protobuf:"bytes,11,opt,name=deployData"`
}

type GatewayConfigs struct {
	Uri         string `json:"uri" protobuf:"bytes,1,opt,name=uri"`
	UpstreamUrl string `json:"upstreamUrl" protobuf:"bytes,2,opt,name=upstreamUrl"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FooList is a list of Foo resources
type PipelineList struct {
	metav1.TypeMeta `json:",inline" protobuf:"bytes,1,opt,name=kind"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,2,opt,name=kind"`
	Items           []Pipeline `json:"items" protobuf:"bytes,3,opt,name=kind"`
}

type PipelineStatus struct {
	Kind           string           `json:"kind" protobuf:"bytes,1,opt,name=kind"`
	Name           string           `json:"name" protobuf:"bytes,2,opt,name=name"`
	Namespace      string           `json:"namespace" protobuf:"bytes,3,opt,name=namespace"`
	Phase          string           `json:"phase"  protobuf:"bytes,4,opt,name=phase"`
	Stages         []PipelineStages `json:"stages" protobuf:"bytes,5,opt,name=stages"`
	StartTimestamp metav1.Time      `json:"startTimestamp" protobuf:"bytes,6,opt,name=startTimestamp"`
}

type PipelineStages struct {
	Name                 string `json:"name" protobuf:"bytes,1,opt,name=name"`
	StartTime            int64  `json:"startTime" protobuf:"bytes,2,opt,name=startTime"`
	DurationMilliseconds int64  `json:"durationMilliseconds" protobuf:"bytes,3,opt,name=durationMilliseconds"`
}
