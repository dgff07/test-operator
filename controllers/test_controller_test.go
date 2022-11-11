package controllers

import (
	"context"
	"time"

	demov1alpha1 "github.com/dgff07/test-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Test controller", func() {
	Context("Test controller test", func() {

		const TestResourceName = "test-resource-example"

		ctx := context.Background()

		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      TestResourceName,
				Namespace: TestResourceName,
			},
		}

		typeNamespaceName := types.NamespacedName{Name: TestResourceName, Namespace: TestResourceName}

		BeforeEach(func() {
			By("Creating the Namespace to perform the tests")
			err := k8sClient.Create(ctx, namespace)
			Expect(err).To(Not(HaveOccurred()))

			By("Override stdout to capture the logs in the tests")

		})

		AfterEach(func() {
			// TODO(user): Attention if you improve this code by adding other context test you MUST
			// be aware of the current delete namespace limitations. More info: https://book.kubebuilder.io/reference/envtest.html#testing-considerations
			By("Deleting the Namespace to perform the tests")
			_ = k8sClient.Delete(ctx, namespace)

		})

		It("should successfully reconcile a custom resource for Test", func() {

			By("Creating the custom resource for the Kind Test")
			testResource := &demov1alpha1.Test{
				ObjectMeta: metav1.ObjectMeta{
					Name:      TestResourceName,
					Namespace: namespace.Name,
				},
				Spec: demov1alpha1.TestSpec{
					Size: 1,
				},
			}
			err := k8sClient.Create(ctx, testResource)
			Expect(err).To(Not(HaveOccurred()))

			By("Checking if the custom resource was successfully created")
			Eventually(func() error {
				found := &demov1alpha1.Test{}
				return k8sClient.Get(ctx, typeNamespaceName, found)
			}, time.Second*3, time.Second).Should(Succeed())

			By("Checking if the namespace is created with the same name as the Test resource")
			Eventually(func() error {
				found := &v1.Namespace{}
				return k8sClient.Get(ctx, typeNamespaceName, found)
			}, time.Second*3, time.Second).Should(Succeed())

		})
	})
})
