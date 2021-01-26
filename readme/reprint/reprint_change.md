# 쿠버네티스 입문 버전 이슈 정리
쿠버네티스 입문은 2021년 2월 기준 쿠버네티스 1.19 기준으로 일부 실습 과정 등에 변경이 있었습니다. 다음은 주요 변경 사항을 정리한 것입니다. 2쇄에는 해당 사항이 반영되어 있음을 알려드립니다.

## apiVersion 변경
쿠버네티스 1.16부터 주요 오브젝트의 apiVersion 필드 값이 변경되어 예제 템플릿 전체적으로 이를 반영했습니다.

* 데몬세트(daemonset): extensions/v1beta1 & apps/v1beta → apps/v1
* 디플로이먼트(deployment): extensions/v1beta1, apps/v1beta1, apps/v1beta2 → apps/v1 
* 스테이트풀세트(statefulset): apps/v1beta1, apps/v1beta2 -> apps/v1
* 레플리카세트(replicaset): extensions/v1beta1, apps/v1beta1, apps/v1beta2 -> apps/v1
* 인그레스(ingress): extensions/v1beta1 -> networking.k8s.io/v1beta1
* 사용자 정의 자원(CustomResourceDefinition): apiextensions.k8s.io/v1beta1 -> apiextensions.k8s.io/v1

## 2.2 도커 데스크톱을 이용한 쿠버네티스 설치
윈도우 10과 macOS에 공통으로 적용되는 사항입니다.

* 도커허브에 회원 가입해 로그인하지 않아도 도커 데스크톱을 다운로드해 사용할 수 있습니다.
* 안정(stable)과 실험(edge)로 나눠서 버전을 관리하던 도커 데스크톱이 하나로 통합되어 해당 부분을 삭제했습니다.

## 2.3.1 구글 쿠버네티스 엔진
구글 클라우드 플랫폼에서 서울 리전이 생기면서 해당 내용을 바꾸었습니다.

## 2.4.2 Kubespray
버전 2.11.0에서 버전 2.14.2로 변경했습니다. 그에 따라 주요 환경의 버전도 업그레이드되었습니다(쿠버네티스 버전은 1.19.5입니다). 단, 설치하는 방법은 1쇄와 다른 점은 없습니다.

Kubespray가 지원하는 네트워크 플러그인이 일부 바뀌었습니다. 기존 8개에서 10개로 늘었습니다. 자세한 내용은 [이곳](https://github.com/kubernetes-sigs/kubespray#network-plugins)을 눌러 참고할 수 있습니다.

## 3.1.1 kubectl 설치

**085쪽 우분투 리눅스 설치 부분**

```bash
$ sudo apt-get update && sudo apt-get install -y apt-transport-https

→

$ sudo apt-get update && sudo apt-get install -y apt-transport-https gnupg2 curl
```

## 3.1.2 기본 사용법
`kubectl run`의 `--generator` 플래그는 더 이상 사용할 수 없습니다. 따라서 087쪽, 096쪽의 명령어에서도 삭제했습니다.

`kubectl` 관련 명령에서 리눅스나 macOS 셸 명령을 연계할 때는 터미널 명령 앞에 `--` 플래그를 붙여야 합니다.

## 3.2.1 kubectl run으로 컨테이너 실행하기 | 7.3 서비스 사용하기
디플로이먼트를 생성하면서 파드를 실행할 때는 더 이상 `kubectl run` 명령을 사용할 수 없습니다.

097쪽의 `kubectl run nginx-app --image nginx --port=80` 명령은 `kubectl create deployment nginx-app --image nginx --port=80`로 수정해서 실행해야 합니다.
200쪽의 `kubectl run nginx-for-service --image=nginx --replicas=2 --port=80 --labels="app=nginx-for-svc"` 역시 `kubectl create deployment nginx-for-svc --image=nginx --replicas=2 --port=80`로 변경해 실행합니다. 한편 `--lables` 플래그가 없어지면서 레이블 이름은 파드 이름인 nginx-for-service로 자동 설정됩니다.

## 4.3.1 네임스페이스
도커 데스크톱 쿠버네티스 클러스터에서 `kubectl get namespace` 명령을 실행했을 때 더는 docker 네임스페이스를 만들지 않습니다.

윈도우용 도커 데스크톱의 네임스페이스 컨텍스트는 docker-for-desktop, macOS는 docker-desktop이었는데 docker-desktop으로 통일되었습니다. 윈도우에서 해당하는 부분에 명령을 입력할 때 모두 docker-desktop으로 바꾸시기 바랍니다.

## 5.4 kubelet으로 컨테이너 진단하기
startupProbe 내용이 추가되었습니다.

컨테이너 안 애플리케이션이 시작되었는지 나타냅니다. 스타트업 프로브는 진단이 성공할 때까지 다른 나머지 프로브는 활성화되지 않으며, 진단이 실패하면 kubelet이 컨테이너를 종료시키고, 컨테이너를 재시작 정책에 따라 처리합니다. 컨테이너에 스타트업 프로브가 없으면 기본 상태 값은 Success입니다.

## 5.7 스태틱 파드
137쪽 macOS 도커 데스크톱의 리눅스 가상 머신에 접속하는 방법이 달라졌습니다.

`screen ~/Library/Containers/com.docker.docker/Data/vms/0/tty` → `docker run -it --privileged --pid=host debian nsenter -t 1 -m -u -n -i sh`

## 6.3.4 디플로이먼트 배포 정리, 배포 재개, 재시작하기
배포 재시작하는 명령과 예를 마지막에 추가했습니다.

```bash
$ kubectl rollout restart deployment/nginx-deployment
deployment.apps/nginx-deployment restarted
$ kubectl rollout history deploy/nginx-deployment
deployment.apps/nginx-deployment
REVISION CHANGE-CAUSE
4 <none>
5 <none>
6 version 1.10.1
7 version 1.11
8 version 1.11
```

## 6.7.1 크론잡 사용하기
191쪽 첫 번째 명령에서 `kubectl run`을 `kubectl create cronjob`으로 변경해 실행해야 합니다.

## 8.1 인그레스의 개념
코드 8-1 ingress/ingress-basic.yaml 파일이 수정되었습니다.

[기존 파일](https://github.com/dybooksIT/kubernetes-book/blob/print1/ingress/ingress-basic.yaml)과 [새 파일](https://github.com/dybooksIT/kubernetes-book/blob/print2/ingress/ingress-basic.yaml)을 비교해서 살펴보기 바랍니다. 경로의 유형을 결정하는 `.pathType`이라는 필드가 새롭게 생기고 `.backend` 필드의 하위 필드 설정이 바뀌었습니다.

## 8.2 ingress-nginx 컨트롤러
설정 방법이 바뀌었습니다.

클론한 후 이동할 ingress-nginx/deploy/baremetal 디렉터리가 ingress-nginx/deploy/static/provider/cloud로 바뀌었습니다.

`kubectl create namespace ~` 명령을 실행하지 않고 ingress-nginx/deploy/static/provider/cloud에서 vi deploy.yaml을 실행한 후 콜론([:] 키)을 눌러 286을 입력하고 [Enter] 키를 누릅니다. 해당 행으로 이동하면 [i] 키를 눌러 `type: LoadBalancer`를 `type: NodePort`로 변경합니다. 그리고 :wq로 vi 편집기를 빠져 나옵니다. 이는 NodePort 타입 기반의 서비스를 만들려는 것입니다.

마지막으로 `kubectl applu -k . ` 명령이 아니라 `kubectl apply -k deploy.yaml` 명령을 실행해 ingress-nginx 컨트롤러를 사용할 준비를 마칩니다.

## 8.3 인그레스 SSL 설정하기
코드 8-2 ingress/ssl/ingress-ssl.yaml 파일이 수정되었습니다.

[기존 파일](https://github.com/dybooksIT/kubernetes-book/blob/print1/ingress/ssl/ingress-ssl.yaml)과 [새 파일](https://github.com/dybooksIT/kubernetes-book/blob/print2/ingress/ingress-basic.yaml)을 비교해서 살펴보기 바랍니다. 코드 8-1과 같은 개념의 수정입니다.

















