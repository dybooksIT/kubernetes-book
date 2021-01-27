# 쿠버네티스 입문 1쇄 독자를 위한 변경 사항 정리
쿠버네티스 입문은 2021년 2월 기준 쿠버네티스 1.19 기준으로 일부 실습 과정 등에 변경이 있었습니다. 다음은 주요 변경 사항을 정리한 것입니다. 2쇄에는 해당 사항이 반영되어 있음을 알려드립니다.

## apiVersion 변경
쿠버네티스 1.16부터 주요 오브젝트의 apiVersion 필드 값이 변경되어 예제 템플릿 전체적으로 이를 반영했습니다.

* 데몬세트(daemonset): extensions/v1beta1 & apps/v1beta → apps/v1
* 디플로이먼트(deployment): extensions/v1beta1, apps/v1beta1, apps/v1beta2 → apps/v1 
* 스테이트풀세트(statefulset): apps/v1beta1, apps/v1beta2 → apps/v1
* 레플리카세트(replicaset): extensions/v1beta1, apps/v1beta1, apps/v1beta2 → apps/v1
* 인그레스(ingress): extensions/v1beta1 -> networking.k8s.io/v1beta1
* 사용자 정의 자원(CustomResourceDefinition): apiextensions.k8s.io/v1beta1 → apiextensions.k8s.io/v1

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

[기존 파일](https://github.com/dybooksIT/kubernetes-book/blob/print1/ingress/ingress-basic.yaml)과 [새 파일](https://github.com/dybooksIT/kubernetes-book/blob/print2/ingress/ingress-basic.yaml)을 비교해서 살펴보기 바랍니다. 경로의 유형을 결정하는 `.pathType`이라는 필드가 새롭게 추가되고 `.backend` 필드의 하위 필드가 정의 방법이 바뀌었습니다.

## 8.2 ingress-nginx 컨트롤러
설정 방법이 바뀌었습니다.

클론한 후 이동할 ingress-nginx/deploy/baremetal 디렉터리가 ingress-nginx/deploy/static/provider/cloud로 바뀌었습니다.

`kubectl create namespace ~` 명령을 실행하지 않고 ingress-nginx/deploy/static/provider/cloud에서 vi deploy.yaml을 실행한 후 콜론([:] 키)을 눌러 286을 입력하고 [Enter] 키를 누릅니다. 해당 행으로 이동하면 [i] 키를 눌러 `type: LoadBalancer`를 `type: NodePort`로 변경합니다. 그리고 :wq로 vi 편집기를 빠져 나옵니다. 이는 NodePort 타입 기반의 서비스를 만들려는 것입니다.

마지막으로 `kubectl applu -k . ` 명령이 아니라 `kubectl apply -k deploy.yaml` 명령을 실행해 ingress-nginx 컨트롤러를 사용할 준비를 마칩니다.

## 8.3 인그레스 SSL 설정하기
코드 8-2 ingress/ssl/ingress-ssl.yaml 파일이 수정되었습니다.

[기존 파일](https://github.com/dybooksIT/kubernetes-book/blob/print1/ingress/ssl/ingress-ssl.yaml)과 [새 파일](https://github.com/dybooksIT/kubernetes-book/blob/print2/ingress/ingress-basic.yaml)을 비교해서 살펴보기 바랍니다. 코드 8-1과 같은 개념의 수정입니다.

## 11.1.2 템플릿으로 시크릿 만들기
시크릿의 타입이 추가되었습니다. 다음 표와 같습니다.

| 시크릿 타입 | 설명 |
|---|---|
| Opaque | 기본값임. 키-값 형식으로 임의의 데이터를 설정할 수 있음 |
| kubernetes.io/service-account-token | 쿠버네티스 인증 토큰을 저장함 |
| kubernetes.io/dockercfg | 도커 저장소 환경 설정 정보를 저장함 |
| kubernetes.io/dockerconfigjson | 도커 저장소 인증 정보를 저장함 |
| kubernetes.io/basic-auth | 기본 인증을 위한 자격 증명을 저장함 |
| kubernetes.io/ssh-auth | SSH 접속을 위한 자격 증명을 저장함 |
| kubernetes.io/tls | TLS 인증서를 저장함 |
| bootstrap.kubernetes.io/token | 부트스트랩 토큰 데이터 정보를 저장함 |

## 13.1.1 kubectl의 config 파일에 있는 TLS 인증 정보 구조 확인하기
311쪽 코드 13-1의 기본 구조가 변경되었습니다. 클러스터 인증에 필요한 해시값을 설정하는 `.clusters[].cluster.certificate-authority-data` 필드가 기본으로 포함됩니다.

## 14.1.2 hostPath | 14.1.3 nfs | 14.3 퍼시스턴트 볼륨 템플릿 | 14.5 레이블로 PVC와 PV 연결하기
코드 14-2 volume/volume-hostpath.yaml, 코드 14-3 volume/volume-nfsserver.yaml의 `.spec.volumes[].hostPath.path` 필드 값에서 macOS Big Sur라면 /tmp라는 경로를 /private/tmp로 수정해야 합니다. 최신 버전의 접근 권한 문제 때문인 것으로 보입니다.

코드 14-5 volume/pv-hostpath.yaml, 코드 14-7 volume/pv-hostpath-label.yaml의 `.spec.hostPath.path` 필드 값에서 macOS Big Sur라면 /tmp/k8s-pv라는 경로를 /private/tmp/k8s-pv로 수정해야 합니다. 최신 버전의 접근 권한 문제 때문인 것으로 보입니다.

## 15.1.2 파드 네트워킹 이해하기
364쪽에서 컨테이너의 NetworkMode를 확인하는 과정이 좀 바뀌었습니다. 2.4.2 Kubespray 기반의 클러스터 환경에서 실행하는 기준입니다.

이제 `kubectl get pods -o wide` 명령을 실행한 후 NODE 항목을 참고해 [코드 15-1]의 파드가 어떤 노드에 있는지 확인합니다. 그리고 `exit` 명령을 실행해 root 계정에서 잠시 로그아웃한 후 `ssh 노드이름 -- 'sudo -i' 'docker ps'` 명령을 실행해 해당 노드에서 [코드 15-1] 파드의 web 혹은 ubuntu의 컨테이너 ID를 확인합니다.

```bash
$ ssh instance-5 -- 'sudo -i' 'docker ps'
```

다음으로 `ssh 노드이름 -- 'sudo -i' 'docker inspect 컨테이너ID | grep Network'` 명령으로 파드에서 사용 중인 컨테이너의 NetworkMode를 살펴보겠습니다.

```bash
$ ssh instance-5 -- 'sudo -i' 'docker inspect c09823cbd0e9 | grep Network'
```

## 15.2 쿠버네티스 서비스 네트워킹
실습을 이어서 진행한 분이라면 368쪽에서 코드 15-2를 저장할 때 `sudo -i` 명령으로 root 계정으로 바꾼 후 `vi pod.yaml` 명령을 실행해야 합니다.

370쪽 마지막은 `ssh` 명령어로 해당 호스트로 접근하기 전에 `exit` 명령으로 root 계정에서 로그아웃한 후 `ssh 노드이름 -- 'sudo iptables -t nat -L'` 명령을 실행합니다.

## 16.3.2 CoreDNS의 질의 구조
381쪽 `kubectl describe configmap coredns -n kube-system` 명령을 실행했을 때의 메시지에서 health 항목에 변화가 있고, ready 항목이 추가되었습니다. upstream 항목이 빠지고 ttl 항목이 추가되면서 그에 관한 설명을 추가했습니다.

* health: http://localhost:8080/health로 CoreDNS의 헬스 체크를 할 수 있습니다. lameduck은 프로세스를 비정상 상태로 만든다는 뜻으로 여기에서는 프로세스가 종료되기 전 5초를 기다리도록 설정했습니다.
* ready: 8181 포트의 HTTP 엔드포인트가 모든 플러그인이 준비되었다는 신호를 보내면 200 OK를 반환합니다.
* ttl: 응답에 대한 사용자 정의 TTL을 지정합니다. 기본값은 5초, 허용되는 최소 TTL은 0초고, 최대값은 3600초입니다. 레코드가 캐싱되지 않는다면 TTL을 0으로 설정합니다.

## 17.1.3 클러스터 레벨 로깅
해당 부분의 실습은 2.4.2 Kubespray 기반의 클러스터 환경에서 진행하기를 권합니다.

393쪽 `docker ps` 명령을 실행하기 전 2.4.2에서 구축한 Kubespray의 instance-1 인스턴스에 SSH로 연결한 후 `sudo -i `명령을 실행해 root 계정으로 접속하기 바랍니다.

394쪽 `kubectl exec -it 파드이름 -- /bin/bash` 명령을 실행하기 전 `kubectl apply -f https://raw.githubusercontent.com/fluent/fluentd-kubernetesdaemonset/master/fluentd-daemonset-syslog.yaml` 명령을 실행해 플루언트디를 설치해두기 바랍니다.

다음으로 `kubectl get pods -n kube-system` 명령을 실행해 플루언트디 파드 이름 하나를 기억해둡니다. 그리고 `kubectl exec -it 기억한파드이름 -- /bin/bash` 명령을 실행해 컨테이너 안에 접속한 후 /var/log/containers 디렉터리로 이동해 `tail -n 1 fluentd-4g5rv_kubesystem_fluentd-XXX~.log` 명령으로 파일 내용을 확인합니다.

도커 런타임에서 로그로테이트 관련 설정을 확인하기 전 `exit` 명령을 실행해 컨테이너 접속을 해제한 후 `ps -ef | grep dockerd` 명령을 실행합니다.

## 17.1.4 플루언트디를 이용해서 로그 수집하기
396쪽 코드 17-3 logging/fluentd-kubernetes-daemonset.yaml에서 `.spec.template` 필드 앞에 다음 필드를 추가해야 합니다.

```yaml
spec:
  selector:
    matchLabels:
      k8s-app: fluentd-logging
      version: v1
```

각주로 알렸던 fluentd-kubernetes-daemonset의 깃허브 저장소가 수정되었습니다.

https://github.com/fluent/fluentd-kubernetes-daemonset/blob/master/docker-image/v1.9/debian-elasticsearch6/Dockerfile

## 17.2 쿠버네티스 대시보드
쿠버네티스 대시보드의 설치 및 로그인 방법이 바뀌었습니다.

403쪽에서 쿠버네티스 대시보드 깃허브에서 recommended.yaml을 다운로드(https://raw.githubusercontent.com/kubernetes/dashboard/master/aio/deploy/recommended.yaml)한 후 dashboard-recommended.yaml로 이름을 바꿉니다. 그리고 코드 17-4처럼 파일 안에 추가 부분을 참고하기 바랍니다.

코드 17-4의 내용이 변경되었습니다. [예제 파일](https://github.com/dybooksIT/kubernetes-book/blob/print2/addon/dashboard-recommended.yml)을 참고하기 바랍니다.

```yaml
apiVersion: v1
 kind: Namespace
metadata:
  name: kubernetes-dashboard
---
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard
  namespace: kubernetes-dashboard
--- # 추가 부분
apiVersion: v1
kind: ServiceAccount
metadata:
  name: admin-user
  namespace: kubernetes-dashboard
---
# 중간 생략
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kubernetes-dashboard
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kubernetes-dashboard
subjects:
- kind: ServiceAccount
  name: kubernetes-dashboard
  namespace: kubernetes-dashboard
--- # 추가 부분
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: admin-user
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: admin-user
  namespace: kubernetes-dashboard
---

# 이후 생략
```

추가한 것은 쿠버네티스 대시보드를 클러스터에 적용한 후 베어러(Bearer) 토큰을 만들어 관리자 로그인하는 데 필요한 서비스 계정과 클러스터롤바인딩 설정입니다. 자세한 내용은 쿠버네티스 대시보드 깃허브의 [Creating sample user](https://github.com/kubernetes/dashboard/blob/master/docs/user/access-control/creating-sample-user.md)를 참고하기 바랍니다.

로그인하는 방법도 다음처럼 바뀌었습니다.

코드 17-4를 `kubectl apply -f dashboard-recommended.yaml` 명령으로 클러스터에 적용한 후 `kubectl proxy` 명령을 실행합니다. 그리고 웹 브라우저에서 http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/ 로 접속하면 대시보드 로그인 화면이 열립니다. 여기에서는 토큰 항목을 선택한 후

`kubectl -n kubernetes-dashboard get secret $(kubectl -n kubernetes-dashboard get sa/admin-user -o jsonpath="{.secrets[0].name}") -o go-template="{{.data.token | base64decode}}"` 명령을 실행해 로그인에 필요한 토큰을 만들고 출력 메시지를 복사해서 대시보드 로그인 화면의 토큰 입력 부분에 붙여넣기한 후 로그인해야 합니다. 자세한 내용은 쿠버네티스 대시보드 깃허브의 [Access](https://github.com/kubernetes/dashboard#access)를 참고하기 바랍니다.

## 17.3.2 힙스터
최신 버전의 쿠버네티스 환경에서는 힙스터를 더는 지원하지 않아 2쇄에서는 내용이 빠졌습니다. 실제로 [깃허브](https://github.com/kubernetes-retired/heapster) 주소도 kubernetes-retired 계정으로 옮겨졌습니다.

## 17.3.3 메트릭 서버
415쪽 설치 및 환경 설정이 수정되었습니다.

메트릭 서버는 Kubernetes Metrics Server 깃허브 저장소(https://github.com/kubernetes-sigs/metrics-server)의 releases 페이지에서 제공하는 components.yaml 파일을 이용해 설치하고 사용합니다. 단, 도커 데스크톱용 쿠버네티스에서 메트릭 서버를 사용하려면 components.yaml의 디플로이먼트 설정에서 `.spec.template.spec.containers[].args[]` 필드에서 다음 부분을 추가합니다.

```yaml
spec:
  template:
    spec:
      containers:
      - args:
        - --cert-dir=/tmp
        - --secure-port=4443
        - --kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname
        - --kubelet-insecure-tls  # 추가해야 하는 부분
        - --kubelet-use-node-status-port
        image: k8s.gcr.io/metrics-server/metrics-server:v0.4.1
        imagePullPolicy: IfNotPresent
```

쿠버네티스 클러스터에서 사용하는 인증서가 공인 인증서가 아니라 사용자 정의 인증서이므로 보안 에러가 발생하지 않도록 무시하는 옵션을 여기에서도 추가하는 것입니다.

도커 데스크톱이 아닌 kubespray 같은 쿠버네티스 클러스터라면 다음 명령을 실행해 메트릭 서버를 설정할 수도 있습니다.

```bash
$ kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

## 17.3.4 프로메테우스
해당 섹션이 17.3.3으로 바뀌었습니다.

프로메테우스 깃허브 저장소(https://github.com/prometheus/prometheus/)의 documentation/examples 디렉터리에 있는 prometheus-kubernetes.yml 파일을 다운로드해 prometheus-kubernetes-config.yaml 바꾸는 것은 같은 데 해당 파일의 내용을 다음처럼 수정해야 합니다.

```yaml
# 이전 생략
# Scrape config for nodes (kubelet).
# 중간 생략
- job_name: 'kubernetes-nodes'

  # Default to scraping over https. If required, just disable this or change to
  # `http`.
  scheme: https

  # 중간 생략
  # <kubernetes_sd_config>.
  tls_config:
    ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    # 중간 생략
    # insecure_skip_verify: true
  bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token

  kubernetes_sd_configs:
  - role: node

  relabel_configs:
  - action: labelmap
    regex: __meta_kubernetes_node_label_(.+)
    # 추가 부분 시작
  - target_label: __address__
    replacement: kubernetes.default.svc:443
  - source_labels: [__meta_kubernetes_node_name]
    regex: (.+)
    target_label: __metrics_path__
    replacement: /api/v1/nodes/${1}/proxy/metrics
    # 추가 부분 끝

# Scrape config for Kubelet cAdvisor.

# 중간 생략

# This job is not necessary and should be removed in Kubernetes 1.6 and
# earlier versions, or it will cause the metrics to be scraped twice.
- job_name: 'kubernetes-cadvisor'

  # Default to scraping over https. If required, just disable this or change to
  # `http`.
  scheme: https

  # 중간 생략
  metrics_path: /metrics/cadvisor

  # 중간 생략
  # <kubernetes_sd_config>.
  tls_config:
    ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    # 중간 생략
    # insecure_skip_verify: true
  bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token

  kubernetes_sd_configs:
  - role: node

  relabel_configs:
  - action: labelmap
    regex: __meta_kubernetes_node_label_(.+)
    # 추가 부분 시작
  - target_label: __address__
    replacement: kubernetes.default.svc:443
  - source_labels: [__meta_kubernetes_node_name]
    regex: (.+)
    target_label: __metrics_path__
    replacement: /api/v1/nodes/${1}/proxy/metrics/cadvisor
    # 추가 부분 끝

# Example scrape config for service endpoints.
# 이후 생략
```

첫 번째 추가 부분은 kubelet과 관련된 메트릭 관련 정보를 가져오는 설정이고 두 번째 추가 부분은 컨테이너 각각의 메트릭을 수집하는 cadvisor의 정보를 가져오는 설정입니다. 해당 설정이 있어야 프로메테우스에서 메트릭 서버와 cadvisor의 메트릭 데이터를 사용할 수 있습니다. 이전에는 포함되었던 해당 정보가 최신 버전 템플릿에서는 제외되어 추가한 것입니다.

코드 17-7 monitoring/prometheus/prometheus-deployment.yaml에서 `.spec.template.spec.containers[].image` 필드 값을 `prom/prometheus:v2.3.2`에서 `prom/prometheus:v2.24.1`로 바꿉니다.

422쪽 맨 아래 kubelet_running_pod_count를 kubelet_running_pods로 바꿔서 입력해야 합니다.

432쪽 코드 17-8 monitoring/prometheus/grafana.yaml에서 `.spec.template` 필드 앞에 다음 필드를 추가해야 합니다.

```yaml
spec:
  selector:
    matchLabels:
      k8s-app: grafana
```

또한 `.spec.template.spec.containers[].image` 필드 값을 `grafana/grafana:5.2.3`에서 `grafana/grafana:7.3.6`으로 바꿉니다.

## 19.4 CRD를 활용한 사용자 정의 컨트롤러 사용하기
코드 19-1이 변경되었습니다. 변경의 주요 부분은 19.5에서 소개하는 자원 유효성 검사 부분의 코드가 기본으로 포함되었다는 점입니다. 이 과정에서 `.spec.validation` 필드가 `.spec.schema` 필드로 변경되었습니다. 하위 필드는 같습니다.

[코드 19-1의 이전 코드](https://github.com/dybooksIT/kubernetes-book/blob/print1/customresourcedefinition/crd-mypod-spec.yaml)와 [새 코드](https://github.com/dybooksIT/kubernetes-book/blob/print2/customresourcedefinition/crd-mypod-spec.yaml), [코드 19-3의 이전 코드](https://github.com/dybooksIT/kubernetes-book/blob/print1/customresourcedefinition/crd-mypod-validation.yaml)와 [새 코드](https://github.com/dybooksIT/kubernetes-book/blob/print2/customresourcedefinition/crd-mypod-validation.yaml),를 비교해서 적용하시기 바랍니다.

## 19.6 사용자 정의 자원의 정보 추가하기
코드 19-4 customresourcedefinition/crd-mypod-spec-additional-partition.yaml에서 `.spec.additionalPrinterColumns.JSONPath` 필드를 `.spec.additionalPrinterColumns.jsonPath` 필드로 바꿔서 사용해야 합니다.

## 19.7 프로메테우스 오퍼레이터 사용하기
450쪽부터 시작하는 설치 과정이 바뀌었습니다.

참고로 설치하기 전에 프로메테우스 오퍼레이터 설치 전 17장에서 설치했던 프로메테우스와 그라파나 설정과 연관되어 문제가 생길 수도 있으므로 `kubectl delete deploy,svc grafana-app prometheus-app prometheus-app-svc` 명령을 실행해 미리 관련 자원을 삭제하기 바랍니다.

먼저 `git clone` 명령을 이용해 kube-prometheus를 클론합니다.

```bash
$ git clone https://github.com/prometheus-operator/kube-prometheus.git
```

kube-prometheus 디렉터리로 이동합니다. 그리고 `kubectl create -f manifests/setup` 명령으로 manifests/setup 디렉터리에 있는 yaml 파일들을 클러스터에 적용합니다.

```bash
$ cd kube-prometheus
$ kubectl apply -f manifests/setup
```

다음 명령을 실행해서 No resources found가 출력되면 클러스터 적용이 완료된 것입니다. 이어서 `kubectl create -f manifests/` 명령을 실행합니다.

```
$ until kubectl get servicemonitors --all-namespaces ; do date; sleep 1; echo ""; done
No resources found
$ kubectl apply -f manifests/
```

다음은 Tip이니 참고하기 바랍니다.

도커 데스크톱에서 프로메테우스 오퍼레이터를 클러스터에 적용시킬 때 간혹 prometheus-k8s-0과 prometheusk8s-1 파드가 Pening 상태로 머무는 상황이 있습니다. 이는 도커 데스크톱에 할당된 메모리 용량이 부족해서 생기는 현상인 상황이 많습니다. [Preference] → [Resources] → [ADVANCED]를 선택한 후 기본 설정되어 있는 2GB를 4GB 정도로 높여서 적용해보기 바랍니다.

## 21.3 헬름 설치하고 사용하기
helm이 1쇄에서는 2 버전이었는데 현재는 3 버전을 사용합니다. 그에 맞춰서 일부 과정이 변경되었습니다. 가장 큰 변화 부분은 틸러 서버가 없어졌다는 것입니다.

478쪽에서 `helm init` 명령을 실행하지 않고 다음 과정을 실행합니다.

`helm repo add stable https://charts.helm.sh/stable` 명령으로 헬름에서 사용할 차트 저장소를 추가합니다.

```bash
$ helm repo add stable https://charts.helm.sh/stable
```

별도의 차트 저장소를 운영하지 않는다면 stable이라는 차트 저장소가 추가됩니다. `helm repo list` 명령으로 저장소 목록을 확인합니다.

```bash
$ helm repo list
```

stable과 로컬 서버에서 저장소를 알 수 있는 도메인을 함께 볼 수 있습니다.

사용할 수 있는 차트들을 확인할 때는 `helm search repo stable` 명령을 실행합니다.

```
$ helm search repo stable
```

여기에서는 차트를 이용해 MySQL을 설치하겠습니다. 먼저 `helm repo update` 명령으로 차트 저장소 정보를 업데이트 합니다.

```bash
$ helm repo update
```

그리고 `helm install my-mysql stable/mysql` 명령을 실행합니다.

```bash
$ helm install my-mysql stable/mysql
```

참고로 `my-mysql`은 차트 이름입니다. 만약 차트 이름을 설정하지 않고 임의의 이름을 만들려고 한다면 `--generate-name`이라는 플래그를 사용합니다. 즉, `helm install stable/mysql --generate-name` 명령을 실행합니다.

## 21.4 헬름 차트의 구조
481쪽에서 차트의 의존성을 명시한 파일 requirements.yaml의 내용이 Chart.yaml에 포함되어 없어졌습니다.

## 21.5 헬름 차트 수정해 사용하기
예제 파일의 메인 디렉터리가 mysql에서 helm으로 바뀌었습니다.

484쪽 `helm install ./mysql` 명령은 `helm install ./mysql ----generate-name`로 실행해야 합니다.

## 21.6 헬름 차트 저장소를 직접 만들어 사용하기
487쪽 `curl http://localhost:8080/api/charts` 명령을 처음 실행했을 때 기존에는 틸러 서버가 있었기에 등록한 차트가 없다고 나오지만 현재는 틸러 서버가 없어서 로컬 서버에 등록된 mysql-0.1.0 관련 차트 정보가 등장합니다. 따라서 새로운 차트를 등록하는 과정을 다음과 같이 변경합니다.

이제 차트를 등록하는 방법도 살펴보겠습니다. [코드 21-2]의 `.name` 필드를 `mysql-addct`로, `.version` 필드를 `0.1.1`로 수정하고 `helm install ./mysql --generate-name` 명령을 실행합니다. 그리고 수정한 mysql-addct 차트를 패키지합니다. 버전 이름을 수정했으므로 mysqladdct-0.1.1.tgz라는 파일로 패키지합니다.

```bash
$ helm package mysql
```

`curl --data-binary "@mysql-0.1.1.tgz" http://localhost:8080/api/charts` 명령을 실행해 차트뮤지엄의 /api/charts api를 이용하도록 요청하면 차트뮤지엄에 차트가 등록됩니다.

```bash
$ curl --data-binary "@mysql-addct-0.1.1.tgz" http://localhost:8080/api/charts
```

다시 `curl http://localhost:8080/api/charts` 명령으로 차트 목록을 확인하면 등록한 차트를 확인할 수 있습니다.

```bash
$ curl http://localhost:8080/api/charts
```

차트뮤지엄의 차트 저장소를 이용해서 mysql-addct 차트를 설치할 때는 다음 명령들을 차례로 실행합니다. 먼저 `helm repo add chartmuseum http://localhost:8080` 명령으로 여러분이 실행한 차트뮤지엄을 새로운 헬름 차트 저장소로 추가합니다.

```bash
$ helm repo add chartmuseum http://localhost:8080
```

`helm search repo chartmuseum` 명령으로 차트 저장소가 제대로 등록됐는지 확인합니다.

```bash
$ helm search repo chartmuseum
```

chartmuseum/mysql과 chartmuseum/mysql-addct라는 차트를 확인할 수 있습니다.

차트를 설치할 때는 `helm install chartmuseum/mysql --generate-name` 명령을 실행합니다.

```bash
$ helm install chartmuseum/mysql-addct --generate-name
```
