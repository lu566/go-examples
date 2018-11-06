package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Foo is a specification for a Foo resource
type Build struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Specification of the desired behavior of the Deployment.
	// +optional
	Spec BuildSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Most recently observed status of the Deployment.
	// +optional
	Status BuildStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type BuildStatus struct {
	Kind           string        `json:"kind" protobuf:"bytes,1,opt,name=kind"`
	Name           string        `json:"name" protobuf:"bytes,2,opt,name=name"`
	Namespace      string        `json:"namespace" protobuf:"bytes,3,opt,name=namespace"`
	Phase          string        `json:"phase"  protobuf:"bytes,4,opt,name=phase"`
	Stages         []BuildStages `json:"stages" protobuf:"bytes,5,opt,name=stages"`
	StartTimestamp int64         `json:"startTimestamp" protobuf:"bytes,6,opt,name=startTimestamp"`
	EventType      []string      `json:"eventType" protobuf:"bytes,7,opt,name=eventType"`
}

type BuildSpec struct {
	CloneConfig    BuildCloneConfig `json:"cloneConfig" protobuf:"bytes,1,opt,name=cloneConfig"`
	App            string           `json:"app" protobuf:"bytes,2,opt,name=app"`
	CodeType       string           `json:"codeType" protobuf:"bytes,3,opt,name=codeType"` //
	CompileCmd     []CompileCmd     `json:"compileCmd" protobuf:"bytes,4,opt,name=compileCmd"`
	CloneType      string           `json:"cloneType" protobuf:"bytes,5,opt,name=cloneType"` //基础镜像包
	BaseImage      string           `json:"baseImage" protobuf:"bytes,6,opt,name=baseImage"`
	Tags           []string         `json:"tags" protobuf:"bytes,7,opt,name=tags"`
	DockerFile     []string         `json:"dockerFile" protobuf:"bytes,8,opt,name=dockerFile"`
	DockerRegistry string           `json:"dockerRegistry" protobuf:"bytes,9,opt,name=dockerRegistry"`
	Events         []string         `json:"events" protobuf:"bytes,12,opt,name=events"`
	NodeService    string           `json:"nodeService" protobuf:"bytes,13,opt,name=nodeService"`
	DeployData     DeployData       `json:"deployData" protobuf:"bytes,14,opt,name=deployData"`
}

type DeployData struct {
	Replicas       int32
	Labels         map[string]string
	Ports          []int
	Envs           map[string]string
	HostPathVolume map[string]string
}

type BuildCloneConfig struct {
	// http://gitlab.vpclb.cn   http://
	Url      string `json:"url"  protobuf:"bytes,1,opt,name=url"`
	Branch   string `json:"branch"  protobuf:"bytes,2,opt,name=branch"`
	DstDir   string `json:"dstDir"  protobuf:"bytes,3,opt,name=dstDir"`
	Depth    int32  `json:"depth,omitempty" protobuf:"varint,4,opt,name=depth,proto3"`
	Username string `json:"username"  protobuf:"bytes,5,opt,name=username"`
	Password string `json:"password"  protobuf:"bytes,6,opt,name=password"`
}

type CompileType string

const (
	Command CompileType = "command"
	Script  CompileType = "script"
)

type CompileCmd struct {
	ExecType      string   `protobuf:"bytes,1,opt,name=execType,proto3" json:"execType,omitempty"`
	Script        string   `protobuf:"bytes,2,opt,name=Script,proto3" json:"Script,omitempty"`
	CommandName   string   `protobuf:"bytes,3,opt,name=commandName,proto3" json:"commandName,omitempty"`
	CommandParams []string `protobuf:"bytes,4,rep,name=commandParams,proto3" json:"commandParams,omitempty"`
}

type BuildStages struct {
	Name                 string `json:"name" protobuf:"bytes,1,opt,name=name"`
	StartTime            int64  `json:"startTime" protobuf:"bytes,2,opt,name=startTime"`
	DurationMilliseconds int64  `json:"durationMilliseconds" protobuf:"bytes,3,opt,name=durationMilliseconds"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FooList is a list of Foo resources
type BuildList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Build `json:"items"`
}
