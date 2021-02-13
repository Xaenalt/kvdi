# Builders

build-ubuntu-base:
	cd build/desktops/ubuntu && docker build . \
		-f Dockerfile.base \
		-t ${UBUNTU_BASE_IMAGE}

build-arch-base:
	cd build/desktops/arch && docker build . \
		-f Dockerfile.base \
		-t ${ARCH_BASE_IMAGE}

build-app-base:
	cd build/desktops/app-profiles && docker build . \
		-f Dockerfile.base \
		-t ${APP_PROFILE_BASE_IMAGE}

build-qemu:
	cd build/desktops/qemu && docker build . \
		-f Dockerfile.base \
		-t ${QEMU_IMAGE}

build-dosbox:
	cd build/desktops/dosbox && docker build . \
		-f Dockerfile.base \
		-t ${DOSBOX_IMAGE}

build-ubuntu-%: build-ubuntu-base
	cd build/desktops/ubuntu && docker build . \
		-f Dockerfile.desktop \
		--build-arg BASE_IMAGE=${UBUNTU_BASE_IMAGE} \
		--build-arg DESKTOP_PACKAGE=$* \
		-t ${REPO}/${NAME}:ubuntu-$*-${VERSION}

build-arch-%: build-arch-base
	cd build/desktops/arch && docker build . \
		-f Dockerfile.$* \
		--build-arg BASE_IMAGE=${ARCH_BASE_IMAGE} \
		-t ${REPO}/${NAME}:arch-$*-${VERSION}

build-app-%:
	cd build/desktops/app-profiles && docker build . \
		-f Dockerfile.$* \
		--build-arg BASE_IMAGE=${APP_PROFILE_BASE_IMAGE} \
		-t ${REPO}/${NAME}:app-$*-${VERSION}

# Pushers

push-ubuntu-base: build-ubuntu-base
	docker push ${UBUNTU_BASE_IMAGE}

push-arch-base: build-arch-base
	docker push ${ARCH_BASE_IMAGE}

push-dosbox: build-dosbox
	docker push ${DOSBOX_IMAGE}

push-qemu: build-qemu
	docker push ${QEMU_IMAGE}

push-ubuntu-%: build-ubuntu-%
	docker push ${REPO}/${NAME}:ubuntu-$*-${VERSION}

push-arch-%: build-arch-%
	docker push ${REPO}/${NAME}:arch-$*-${VERSION}

push-app-%: build-app-%
	docker push ${REPO}/${NAME}:app-$*-${VERSION}

# Loaders

load-ubuntu-%: $(K3D) build-ubuntu-%
	$(call load_image,${REPO}/${NAME}:ubuntu-$*-${VERSION})

load-arch-%: $(K3D) build-arch-%
	$(call load_image,${REPO}/${NAME}:arch-$*-${VERSION})

load-app-%: $(K3D) build-app-%
	$(call load_image,${REPO}/${NAME}:app-$*-${VERSION})

load-dosbox: $(K3D) build-dosbox
	$(call load_image,${DOSBOX_IMAGE})

load-qemu: $(K3D) build-qemu
	$(call load_image,${QEMU_IMAGE})

#
# For building demo environment

demo-init:
	cd deploy/terraform && terraform init

demo-fmt:
	terraform fmt deploy/terraform

get_ext_ip = $(shell curl https://ifconfig.me 2> /dev/null)
auto_approve = $(shell [[ "$(1)" != "plan" ]] && echo -auto-approve)
TF_ARGS ?=
demo-%:
	cd deploy/terraform && \
		terraform $* \
			-var ext_ip=$(call get_ext_ip) \
			$(TF_ARGS) $(call auto_approve,$*)

demo-state:
	@cd deploy/terraform && terraform state pull

get_demo_host = $(shell cd deploy/terraform && terraform output public_ip)
ssh_demo = 	ssh -o "UserKnownHostsFile /dev/null" -o "StrictHostKeyChecking=no" ec2-user@$(call get_demo_host)
demo-shell:
	$(call ssh_demo)

demo-password:
	$(call ssh_demo) sudo k3s kubectl get secret kvdi-admin-secret -o json | jq .data.password -r | base64 -d && echo
