package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DesktopTemplateSpec defines the desired state of DesktopTemplate
type DesktopTemplateSpec struct {
	// The docker repository and tag to use for desktops booted from this template.
	Image string `json:"image"`
	// The pull policy to use when pulling the container image.
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// Any pull secrets required for pulling the container image.
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
	// Resource requirements to apply to desktops booted from this template.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Configuration options for the instances. This is highly dependant on using
	// the Dockerfiles (or close derivitives) provided in this repository.
	Config *DesktopConfig `json:"config,omitempty"`
	// Arbitrary tags for displaying in the app UI.
	Tags map[string]string `json:"tags,omitempty"`
}

type DesktopConfig struct {
	// A service account to tie to desktops booted from this template.
	// TODO: This should really be per-desktop and by user-grants.
	ServiceAccount string `json:"serviceAccount,omitempty"`
	// Extra system capabilities to add to desktops booted from this template.
	Capabilities []corev1.Capability `json:"capabilities,omitempty"`
	// Whether the sound device should be mounted inside the container. Note that
	// this also requires the image do proper setup if /dev/snd is present.
	EnableSound bool `json:"enableSound,omitempty"`
	// AllowRoot will pass the ENABLE_ROOT envvar to the container. In the Dockerfiles
	// in this repository, this will add the user to the sudo group and ability to
	// sudo with no password.
	AllowRoot bool `json:"allowRoot,omitempty"`
	// The address the VNC server listens on inside the image. This defaults to the
	// UNIX socket /var/run/kvdi/vnc.sock. The novnc-proxy sidecar will forward
	// websockify requests validated by mTLS to this socket.
	// Must be in the format of `tcp://{host}:{port}` or `unix://{path}`.
	SocketAddr string `json:"socketAddr,omitempty"`
}

// DesktopTemplateStatus defines the observed state of DesktopTemplate
type DesktopTemplateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DesktopTemplate is the Schema for the desktoptemplates API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=desktoptemplates,scope=Cluster
type DesktopTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DesktopTemplateSpec   `json:"spec,omitempty"`
	Status DesktopTemplateStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DesktopTemplateList contains a list of DesktopTemplate
type DesktopTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DesktopTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DesktopTemplate{}, &DesktopTemplateList{})
}