package tests_test

import (
	"sigs.k8s.io/kustomize/v3/k8sdeps/kunstruct"
	"sigs.k8s.io/kustomize/v3/k8sdeps/transformer"
	"sigs.k8s.io/kustomize/v3/pkg/fs"
	"sigs.k8s.io/kustomize/v3/pkg/loader"
	"sigs.k8s.io/kustomize/v3/pkg/plugins"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/kustomize/v3/pkg/resource"
	"sigs.k8s.io/kustomize/v3/pkg/target"
	"sigs.k8s.io/kustomize/v3/pkg/validators"
	"testing"
)

func writeJupyterWebAppBaseApplication(th *KustTestHarness) {
	th.writeF("/manifests/jupyter/jupyter-web-app/base/application/application.yaml", `
apiVersion: app.k8s.io/v1beta1
kind: Application
metadata:
  name: jupyter-web-app
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: jupyter-web-app
      app.kubernetes.io/instance: jupyter-web-app-v0.7.0
      app.kubernetes.io/managed-by: kfctl
      app.kubernetes.io/component: jupyter-web-app
      app.kubernetes.io/part-of: kubeflow
      app.kubernetes.io/version: v0.7.0
  componentKinds:
  - group: core
    kind: ConfigMap
  - group: apps
    kind: Deployment
  - group: rbac.authorization.k8s.io
    kind: RoleBinding
  - group: rbac.authorization.k8s.io
    kind: Role
  - group: core
    kind: ServiceAccount
  - group: core
    kind: Service
  - group: networking.istio.io
    kind: VirtualService
  descriptor:
    type: jupyter-web-app
    version: v1beta1
    description: Provides a UI which allows the user to create/conect/delete jupyter notebooks.
    maintainers:
    - name: Kimonas Sotirchos
      email: kimwnasptd@arrikto.com
    owners:
    - name: Kimonas Sotirchos
      email: kimwnasptd@arrikto.com
    keywords:
     - jupyterhub
     - jupyter ui
     - notebooks  
    links:
    - description: About
      url: https://github.com/kubeflow/kubeflow/tree/master/components/jupyter-web-app
    - description: Docs
      url: https://www.kubeflow.org/docs/notebooks 
  addOwnerRef: true

`)
	th.writeK("/manifests/jupyter/jupyter-web-app/base/application", `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- application.yaml
`)
	th.writeK("/manifests/jupyter/jupyter-web-app/base", `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
- application
- core
- istio
commonLabels:
  app.kubernetes.io/name: jupyter-web-app
  app.kubernetes.io/instance: jupyter-web-app-v0.7.0
  app.kubernetes.io/managed-by: kfctl
  app.kubernetes.io/component: jupyter-web-app
  app.kubernetes.io/part-of: kubeflow
  app.kubernetes.io/version: v0.7.0
`)
}

func TestJupyterWebAppBaseApplication(t *testing.T) {
	th := NewKustTestHarness(t, "/manifests/jupyter/jupyter-web-app/base/application")
	writeJupyterWebAppBaseApplication(th)
	m, err := th.makeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	expected, err := m.AsYaml()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	targetPath := "../jupyter/jupyter-web-app/base/application"
	fsys := fs.MakeRealFS()
	lrc := loader.RestrictionRootOnly
	_loader, loaderErr := loader.NewLoader(lrc, validators.MakeFakeValidator(), targetPath, fsys)
	if loaderErr != nil {
		t.Fatalf("could not load kustomize loader: %v", loaderErr)
	}
	rf := resmap.NewFactory(resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl()), transformer.NewFactoryImpl())
	pc := plugins.DefaultPluginConfig()
	kt, err := target.NewKustTarget(_loader, rf, transformer.NewFactoryImpl(), plugins.NewLoader(pc, rf))
	if err != nil {
		th.t.Fatalf("Unexpected construction error %v", err)
	}
	actual, err := kt.MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	th.assertActualEqualsExpected(actual, string(expected))
}
