deploy-minikube:
	minikube image build -t email-svc:v1 .
	cd deployment && pulumi up -y