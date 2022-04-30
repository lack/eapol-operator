package v1

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func makeA11r() *Authenticator {
	a11r := &Authenticator{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "eapol.eapol.openshift.io/v1",
			Kind:       "Authenticator",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "authenticator",
			Namespace: "default",
		},
		Spec: AuthenticatorSpec{
			Enabled:    true,
			Interfaces: []string{"eth0"},
			Authentication: Auth{
				Radius: &Radius{},
			},
		},
	}
	return a11r
}

var _ = Describe("Ensure exactly one authentication mechanism is defined", func() {
	var a11r *Authenticator
	BeforeEach(func() {
		a11r = makeA11r()
		a11r.Spec.Authentication.Local = nil
		a11r.Spec.Authentication.Radius = nil
	})
	It("Should fail if no configuration is present", func() {
		Expect(assertHasOneAuthMethod(a11r)).NotTo(Succeed())
	})
	It("Should succeed if local authentication is configured", func() {
		a11r.Spec.Authentication.Local = &Local{}
		Expect(assertHasOneAuthMethod(a11r)).To(Succeed())
	})
	It("Should succeed if RADIUS authentication is configured", func() {
		a11r.Spec.Authentication.Radius = &Radius{}
		Expect(assertHasOneAuthMethod(a11r)).To(Succeed())
	})
})

var _ = Describe("Ensure Interface list is not empty", func() {
	var a11r *Authenticator
	BeforeEach(func() {
		a11r = makeA11r()
	})
	It("Should fail if the interface list is empty", func() {
		a11r.Spec.Interfaces = []string{}
		Expect(assertValidInterfaces(a11r)).NotTo(Succeed())
	})
	It("Should succeed if the interface list has one entry", func() {
		a11r.Spec.Interfaces = []string{"one"}
		Expect(assertValidInterfaces(a11r)).To(Succeed())
	})
	It("Should succeed if the interface list has more than one entry", func() {
		a11r.Spec.Interfaces = []string{"one", "two"}
		Expect(assertValidInterfaces(a11r)).To(Succeed())
	})
})
