# Install Kubebuilder
curl -L -o kubebuilder "https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)"\nchmod +x kubebuilder && sudo mv kubebuilder /usr/local/bin/

# Create Project
kubebuilder init \ 
    --domain hector.dev \ 
    --repo github.com/Teachh/K8S-Shared-PVC

# Create API with CRD and a Controller
kubebuilder create api \ 
    --group crd \ 
    --version v1 \ 
    --kind SharedPVC \ 
    --resource true \ 
    --controller true
