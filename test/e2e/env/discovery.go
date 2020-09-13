package env

import (
	"errors"
	"fmt"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	cacheddiscovery "k8s.io/client-go/discovery/cached"
	fakeDiscoveryF "k8s.io/client-go/discovery/fake"
)

type resourceMapEntry struct {
	list *metav1.APIResourceList
	err  error
}
type fakeDiscovery struct {
	*fakeDiscoveryF.FakeDiscovery

	lock         sync.Mutex
	groupList    *metav1.APIGroupList
	groupListErr error
	resourceMap  map[string]*resourceMapEntry
}

func (c *fakeDiscovery) ServerResourcesForGroupVersion(groupVersion string) (*metav1.APIResourceList, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if rl, ok := c.resourceMap[groupVersion]; ok {
		return rl.list, rl.err
	}
	return nil, errors.New("doesn't exist")
}

func (c *fakeDiscovery) ServerGroups() (*metav1.APIGroupList, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.groupList == nil {
		return nil, errors.New("doesn't exist")
	}
	return c.groupList, c.groupListErr
}

func FakeCachedDiscoveryInterface() discovery.CachedDiscoveryInterface {
	podGVR := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	serviceGVR := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	ingressGVR := schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1beta1", Resource: "ingresses"}
	fooGVR := schema.GroupVersionResource{Group: "samplecontroller.k8s.io", Version: "v1alpha1", Resource: "foos"}
	namespaceGVR := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"}

	fake := &fakeDiscovery{
		groupList: &metav1.APIGroupList{
			Groups: []metav1.APIGroup{
				{
					Name: podGVR.Group,
					Versions: []metav1.GroupVersionForDiscovery{{
						GroupVersion: podGVR.GroupVersion().String(),
						Version:      podGVR.Version,
					}},
				},
				{
					Name: serviceGVR.Group,
					Versions: []metav1.GroupVersionForDiscovery{{
						GroupVersion: serviceGVR.GroupVersion().String(),
						Version:      serviceGVR.Version,
					}},
				},
				{
					Name: ingressGVR.Group,
					Versions: []metav1.GroupVersionForDiscovery{{
						GroupVersion: ingressGVR.GroupVersion().String(),
						Version:      ingressGVR.Version,
					}},
				},
				{
					Name: fooGVR.Group,
					Versions: []metav1.GroupVersionForDiscovery{{
						GroupVersion: fooGVR.GroupVersion().String(),
						Version:      fooGVR.Version,
					}},
				},
				{
					Name: namespaceGVR.Group,
					Versions: []metav1.GroupVersionForDiscovery{{
						GroupVersion: namespaceGVR.GroupVersion().String(),
						Version:      namespaceGVR.Version,
					}},
				},
			},
		},
		resourceMap: map[string]*resourceMapEntry{
			podGVR.GroupVersion().String(): {
				list: &metav1.APIResourceList{
					GroupVersion: podGVR.GroupVersion().String(),
					APIResources: []metav1.APIResource{
						{
							Name:         podGVR.Resource,
							SingularName: "pod",
							Namespaced:   true,
							Kind:         "Pod",
							ShortNames:   []string{"pod"},
						},
						{
							Name:         serviceGVR.Resource,
							SingularName: "service",
							Namespaced:   true,
							Kind:         "Service",
							ShortNames:   []string{"svc"},
						},
						{
							Name:         ingressGVR.Resource,
							SingularName: "ingress",
							Namespaced:   true,
							Kind:         "Ingress",
							ShortNames:   []string{"ingress"},
						},
						{
							Name:         fooGVR.Resource,
							SingularName: "foo",
							Namespaced:   true,
							Kind:         "Foo",
							ShortNames:   []string{"foo"},
						},
						{
							Name:         namespaceGVR.Resource,
							SingularName: "namespace",
							Namespaced:   true,
							Kind:         "Namespace",
							ShortNames:   []string{"ns"},
						},
					},
				},
			},
		},
	}

	discoCacheClient := cacheddiscovery.NewMemCacheClient(fake)
	getRes, err := discoCacheClient.ServerResources()
	if err != nil {
		fmt.Printf("Printing Error: %v\n", err)
	}
	fmt.Printf("ServerResources: %v\n", getRes)

	return discoCacheClient
}
