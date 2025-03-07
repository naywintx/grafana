package scope

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"

	scope "github.com/grafana/grafana/pkg/apis/scope/v0alpha1"
	grafanaregistry "github.com/grafana/grafana/pkg/apiserver/registry/generic"
	grafanarest "github.com/grafana/grafana/pkg/apiserver/rest"
	"github.com/grafana/grafana/pkg/services/apiserver/utils"
	apistore "k8s.io/apiserver/pkg/storage"
)

var _ grafanarest.Storage = (*storage)(nil)

type storage struct {
	*genericregistry.Store
}

func newScopeStorage(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*storage, error) {
	strategy := grafanaregistry.NewStrategy(scheme)

	resourceInfo := scope.ScopeResourceInfo
	store := &genericregistry.Store{
		NewFunc:                   resourceInfo.NewFunc,
		NewListFunc:               resourceInfo.NewListFunc,
		PredicateFunc:             Matcher,
		DefaultQualifiedResource:  resourceInfo.GroupResource(),
		SingularQualifiedResource: resourceInfo.SingularGroupResource(),
		TableConvertor: utils.NewTableConverter(
			resourceInfo.GroupResource(),
			[]metav1.TableColumnDefinition{
				{Name: "Name", Type: "string", Format: "name"},
				{Name: "Created At", Type: "date"},
			},
			func(obj any) ([]interface{}, error) {
				m, ok := obj.(*scope.Scope)
				if !ok {
					return nil, fmt.Errorf("expected scope")
				}
				return []interface{}{
					m.Name,
					m.CreationTimestamp.UTC().Format(time.RFC3339),
				}, nil
			},
		),
		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &storage{Store: store}, nil
}

func newScopeDashboardStorage(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*storage, error) {
	strategy := grafanaregistry.NewStrategy(scheme)

	resourceInfo := scope.ScopeDashboardResourceInfo
	store := &genericregistry.Store{
		NewFunc:                   resourceInfo.NewFunc,
		NewListFunc:               resourceInfo.NewListFunc,
		PredicateFunc:             Matcher,
		DefaultQualifiedResource:  resourceInfo.GroupResource(),
		SingularQualifiedResource: resourceInfo.SingularGroupResource(),
		TableConvertor: utils.NewTableConverter(
			resourceInfo.GroupResource(),
			[]metav1.TableColumnDefinition{
				{Name: "Name", Type: "string", Format: "name"},
				{Name: "Created At", Type: "date"},
			},
			func(obj any) ([]interface{}, error) {
				m, ok := obj.(*scope.Scope)
				if !ok {
					return nil, fmt.Errorf("expected scope")
				}
				return []interface{}{
					m.Name,
					m.CreationTimestamp.UTC().Format(time.RFC3339),
				}, nil
			},
		),
		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &storage{Store: store}, nil
}

func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	s, ok := obj.(*scope.Scope)
	if !ok {
		return nil, nil, fmt.Errorf("not a scope")
	}

	return labels.Set(s.Labels), SelectableFields(s), nil
}

// Matcher returns a generic.SelectionPredicate that matches on label and field selectors.
func Matcher(label labels.Selector, field fields.Selector) apistore.SelectionPredicate {
	return apistore.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

func SelectableFields(obj *scope.Scope) fields.Set {
	return generic.MergeFieldsSets(generic.ObjectMetaFieldsSet(&obj.ObjectMeta, false), fields.Set{
		"spec.type": obj.Spec.Type,
	})
}
