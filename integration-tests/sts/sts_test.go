package sts

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openebs/maya/integration-tests/artifacts"
	k8s "github.com/openebs/maya/pkg/client/k8s/v1alpha1"
	csp "github.com/openebs/maya/pkg/cstorpool/v1alpha2"
	cvr "github.com/openebs/maya/pkg/cstorvolumereplica/v1alpha1"
	pvc "github.com/openebs/maya/pkg/kubernetes/persistentvolumeclaim/v1alpha1"
	pod "github.com/openebs/maya/pkg/kubernetes/pod/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// stsYaml holds the yaml spec
	// for statefulset application
	stsYaml artifacts.Artifact = `
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: busybox1
  namespace: default
  labels:
    app: busybox1
spec:
  serviceName: busybox1
  replicas: 3
  selector:
    matchLabels:
      app: busybox1
      openebs.io/replica-anti-affinity: busybox1
  template:
    metadata:
      labels:
        app: busybox1
        openebs.io/replica-anti-affinity: busybox1
    spec:
      containers:
      - name: busybox1
        image: ubuntu
        imagePullPolicy: IfNotPresent
        command:
          - sleep
          - infinity
        volumeMounts:
        - name: busybox1
          mountPath: /busybox1
  volumeClaimTemplates:
  - metadata:
      name: busybox1
    spec:
      accessModes: [ "ReadWriteOnce" ]
      storageClassName: cstor-sts
      resources:
        requests:
          storage: 1Gi`

	// stsSCYaml holds the yaml spe for
	// storageclass required by the statefulset application
	stsSCYaml artifacts.Artifact = `
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: cstor-sts
  annotations:
    openebs.io/cas-type: cstor
    cas.openebs.io/config: |
      - name: ReplicaCount
        value: "1"
      - name: StoragePoolClaim
        value: "cstor-sparse-pool" 
provisioner: openebs.io/provisioner-iscsi
`
)

var _ = Describe("StatefulSet", func() {
	BeforeEach(func() {
		// Extracting storageclass artifacts unstructured
		SCUnstructured, err := artifacts.GetArtifactUnstructured(
			artifacts.Artifact(stsSCYaml),
		)
		Expect(err).ShouldNot(HaveOccurred())

		// Extracting statefulset artifacts unstructured
		STSUnstructured, err := artifacts.GetArtifactUnstructured(stsYaml)
		Expect(err).ShouldNot(HaveOccurred())

		// Extracting statefulset application namespace
		stsNamespace := STSUnstructured.GetNamespace()

		// Generating label selector for stsResources
		stsApplicationLabel := "app=" + STSUnstructured.GetName()

		// Apply sts storageclass
		cu := k8s.CreateOrUpdate(
			k8s.GroupVersionResourceFromGVK(SCUnstructured),
			SCUnstructured.GetNamespace(),
		)
		_, err = cu.Apply(SCUnstructured)
		Expect(err).ShouldNot(HaveOccurred())

		// Apply the sts
		cu = k8s.CreateOrUpdate(
			k8s.GroupVersionResourceFromGVK(STSUnstructured),
			stsNamespace,
		)
		_, err = cu.Apply(STSUnstructured)
		Expect(err).ShouldNot(HaveOccurred())

		// Verify creation of sts instances

		// Check for pvc to get created and bound
		Eventually(func() int {
			pvcs, err := pvc.
				KubeClient(pvc.WithNamespace(stsNamespace)).
				List(metav1.ListOptions{LabelSelector: stsApplicationLabel})
			Expect(err).ShouldNot(HaveOccurred())
			return pvc.
				ListBuilder().
				WithAPIList(pvcs).
				WithFilter(pvc.IsBound()).
				List().
				Len()
		},
			defaultTimeOut, defaultPollingInterval).
			Should(Equal(3), "PVC count should be "+string(3))

		// Check for statefulset pods to get created and running
		Eventually(func() int {
			pods, err := pod.
				KubeClient(pod.WithNamespace(stsNamespace)).
				List(metav1.ListOptions{LabelSelector: stsApplicationLabel})
			Expect(err).ShouldNot(HaveOccurred())
			return pod.
				ListBuilder().
				WithAPIList(pods).
				WithFilter(pod.IsRunning()).
				List().
				Len()
		},
			defaultTimeOut, defaultPollingInterval).
			Should(Equal(3), "Pod count should be "+string(3))
	})

	AfterEach(func() {
		// Extracting storageclass artifacts unstructured
		SCUnstructured, err := artifacts.GetArtifactUnstructured(artifacts.Artifact(stsSCYaml))
		Expect(err).ShouldNot(HaveOccurred())

		// Extracting statefulset artifacts unstructured
		STSUnstructured, err := artifacts.GetArtifactUnstructured(stsYaml)
		Expect(err).ShouldNot(HaveOccurred())

		// Extracting statefulset application namespace
		stsNamespace := STSUnstructured.GetNamespace()

		// Generating label selector for stsResources
		stsApplicationLabel := "app=" + STSUnstructured.GetName()

		// Fetch PVCs to be deleted
		pvcs, err := pvc.KubeClient(pvc.WithNamespace(stsNamespace)).
			List(metav1.ListOptions{LabelSelector: stsApplicationLabel})
		Expect(err).ShouldNot(HaveOccurred())
		// Delete PVCs
		for _, p := range pvcs.Items {
			err = pvc.KubeClient(pvc.WithNamespace(stsNamespace)).
				Delete(p.GetName(), &metav1.DeleteOptions{})
			Expect(err).ShouldNot(HaveOccurred())
		}

		// Delete the sts artifacts
		cu := k8s.DeleteResource(
			k8s.GroupVersionResourceFromGVK(STSUnstructured),
			stsNamespace,
		)
		err = cu.Delete(STSUnstructured)
		Expect(err).ShouldNot(HaveOccurred())

		// Verify deletion of sts instances
		Eventually(func() int {
			pods, err := pod.
				KubeClient(pod.WithNamespace(stsNamespace)).
				List(metav1.ListOptions{LabelSelector: stsApplicationLabel})
			Expect(err).ShouldNot(HaveOccurred())
			return len(pods.Items)
		}, defaultTimeOut, defaultPollingInterval).
			Should(Equal(0), "pod count should be 0")

		// Verify deletion of pvc instances
		Eventually(func() int {
			pvcs, err := pvc.
				KubeClient(pvc.WithNamespace(stsNamespace)).
				List(metav1.ListOptions{LabelSelector: stsApplicationLabel})
			Expect(err).ShouldNot(HaveOccurred())
			return len(pvcs.Items)
		},
			defaultTimeOut, defaultPollingInterval).
			Should(Equal(0), "pvc count should be 0")

		// Delete storageclass
		cu = k8s.DeleteResource(
			k8s.GroupVersionResourceFromGVK(SCUnstructured),
			SCUnstructured.GetNamespace(),
		)
		err = cu.Delete(SCUnstructured)
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("test statefulset application on cstor", func() {
		It("should distribute the cstor volume replicas across pools", func() {
			// Extracting statefulset artifacts unstructured

			STSUnstructured, err := artifacts.GetArtifactUnstructured(stsYaml)
			Expect(err).ShouldNot(HaveOccurred())

			// Extracting statefulset application namespace
			stsNamespace := STSUnstructured.GetNamespace()

			// Generating label selector for stsResources
			stsApplicationLabel := "app=" + STSUnstructured.GetName()
			replicaAntiAffinityLabel := "openebs.io/replica-anti-affinity=" + STSUnstructured.GetName()

			pvcs, err := pvc.
				KubeClient(pvc.WithNamespace(stsNamespace)).
				List(metav1.ListOptions{LabelSelector: stsApplicationLabel})
			Expect(err).ShouldNot(HaveOccurred())
			pvcList := pvc.
				ListBuilder().
				WithAPIList(pvcs).
				WithFilter(pvc.ContainsName(STSUnstructured.GetName())).
				List()
			Expect(pvcList.Len()).Should(Equal(3), "pvc count should be "+string(3))

			cvrs, err := cvr.
				KubeClient().
				List("", metav1.ListOptions{LabelSelector: replicaAntiAffinityLabel})
			Expect(cvrs.Items).Should(HaveLen(3), "cvr count should be "+string(3))

			poolNames := cvr.
				ListBuilder().
				WithAPIList(cvrs).
				List()
			Expect(poolNames.GetUniquePoolNames()).
				Should(HaveLen(3), "pool names count should be "+string(3))

			pools, err := csp.KubeClient().List(metav1.ListOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			nodeNames := csp.ListBuilder().WithAPIList(pools).List()
			Expect(nodeNames.GetPoolUIDs()).
				Should(HaveLen(3), "node names count should be "+string(3))
		})

		PIt("should co-locate the cstor volume targets with application instances", func() {
			// future
		})
	})
})
