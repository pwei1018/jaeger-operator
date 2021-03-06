package storage

import (
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/jaegertracing/jaeger-operator/pkg/apis/jaegertracing/v1"
	"github.com/jaegertracing/jaeger-operator/pkg/util"
)

// ESRole returns the role to be created for Elasticsearch
func ESRole(jaeger *v1.Jaeger) rbacv1.Role {
	return rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{rbacv1.AutoUpdateAnnotationKey: "true"},
			Labels: map[string]string{
				"app":                          "jaeger",
				"app.kubernetes.io/name":       fmt.Sprintf("%s-elasticsearch", jaeger.Name),
				"app.kubernetes.io/instance":   jaeger.Name,
				"app.kubernetes.io/component":  "es-role",
				"app.kubernetes.io/part-of":    "jaeger",
				"app.kubernetes.io/managed-by": "jaeger-operator",
			},
			Name:            fmt.Sprintf("%s-elasticsearch", jaeger.Name),
			Namespace:       jaeger.Namespace,
			OwnerReferences: []metav1.OwnerReference{util.AsOwner(jaeger)},
		},
		Rules: []rbacv1.PolicyRule{
			{
				// These values are virtual and defined in SearchGuard sg_config.yml under subjectAccessReviews
				// The SG invokes this API to allow the request
				// TOKEN=$(oc serviceaccounts get-token jaeger-simple-prod)
				// curl -k -v -XPOST  -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" https://127.0.0.1:8443/apis/authorization.k8s.io/v1/selfsubjectaccessreviews -d '{"kind":"SelfSubjectAccessReview","apiVersion":"authorization.k8s.io/v1","spec":{"resourceAttributes":{"group":"jaeger.openshift.io","verb":"get","resource":"jaeger"}}}'
				APIGroups: []string{"elasticsearch.jaegertracing.io"},
				Resources: []string{"jaeger"},
				Verbs:     []string{"get"},
			},
		},
	}
}

// ESRoleBinding returns the Elasticsearch role bindings to be created for the given subjects
func ESRoleBinding(jaeger *v1.Jaeger, sas ...string) rbacv1.RoleBinding {
	rb := rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-elasticsearch", jaeger.Name),
			Namespace: jaeger.Namespace,
			Labels: map[string]string{
				"app":                          "jaeger",
				"app.kubernetes.io/name":       fmt.Sprintf("%s-elasticsearch", jaeger.Name),
				"app.kubernetes.io/instance":   jaeger.Name,
				"app.kubernetes.io/component":  "es-rolebinding",
				"app.kubernetes.io/part-of":    "jaeger",
				"app.kubernetes.io/managed-by": "jaeger-operator",
			},
			OwnerReferences: []metav1.OwnerReference{util.AsOwner(jaeger)},
		},
		RoleRef: rbacv1.RoleRef{
			Kind: "Role",
			Name: fmt.Sprintf("%s-elasticsearch", jaeger.Name),
		},
	}
	for _, sa := range sas {
		sb := rbacv1.Subject{
			Kind:      rbacv1.ServiceAccountKind,
			Namespace: jaeger.Namespace,
			Name:      sa,
		}
		rb.Subjects = append(rb.Subjects, sb)
	}

	return rb
}
