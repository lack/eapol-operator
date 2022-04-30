/*
Copyright 2022.

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

package v1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var authenticatorlog = logf.Log.WithName("authenticator-resource")

func (r *Authenticator) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-eapol-eapol-openshift-io-v1-authenticator,mutating=true,failurePolicy=fail,sideEffects=None,groups=eapol.eapol.openshift.io,resources=authenticators,verbs=create;update,versions=v1,name=mauthenticator.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Authenticator{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Authenticator) Default() {
	authenticatorlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-eapol-eapol-openshift-io-v1-authenticator,mutating=false,failurePolicy=fail,sideEffects=None,groups=eapol.eapol.openshift.io,resources=authenticators,verbs=create;update,versions=v1,name=vauthenticator.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Authenticator{}

func assertHasOneAuthMethod(r *Authenticator) error {
	if r.Spec.Authentication.Local == nil && r.Spec.Authentication.Radius == nil {
		return fmt.Errorf("Authenticator must have at least one authentication mechanism configured")
	}
	return nil
}

func assertValidInterfaces(r *Authenticator) error {
	if len(r.Spec.Interfaces) == 0 {
		return fmt.Errorf("Authenticator must have at least one interface defined")
	}
	return nil
}

var assertions = []func(*Authenticator) error{
	assertHasOneAuthMethod,
	assertValidInterfaces,
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Authenticator) ValidateCreate() error {
	authenticatorlog.Info("validate create", "name", r.Name)

	for _, assert := range assertions {
		if err := assert(r); err != nil {
			return err
		}
	}
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Authenticator) ValidateUpdate(old runtime.Object) error {
	authenticatorlog.Info("validate update", "name", r.Name)

	for _, assert := range assertions {
		if err := assert(r); err != nil {
			return err
		}
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Authenticator) ValidateDelete() error {
	authenticatorlog.Info("validate delete", "name", r.Name)

	// No-op
	return nil
}
